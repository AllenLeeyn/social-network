// src/components/groups/GroupDetail.js
// groups/uuid/main-view, main-feed

import React, { useState, useEffect } from "react";
import { toast } from 'react-toastify';

import Modal from "../Modal";
import CreatePostForm from "./[uuid]/posts/CreatePostForm";
import CreateEventForm from "./[uuid]/events/CreateEventForm";
import GroupHeader from "./[uuid]/GroupHeader";
import EventCard from "./[uuid]/events/EventCard";

import { formatDate } from '../../utils/formatDate';

import "../../styles/groups/GroupDetail.css"; 


export default function GroupDetail(
    {   group, 
        members, 
        requests,
        allUsers,
        loadingUsers,
        onBack, 
        onRequestJoin,
        onApproveRequest,   
        onDenyRequest,
        onInviteUser,
    }) {
        
    if (!group) return null;

    // Modal state
    const [showPostModal, setShowPostModal] = useState(false);
    const [showEventModal, setShowEventModal] = useState(false);
    const [showMembersModal, setShowMembersModal] = useState(false);

    const handleOpenMembersModal = () => setShowMembersModal(true);
    const handleCloseMembersModal = () => setShowMembersModal(false);

    // Events state
    const [groupEvents, setGroupEvents] = useState([]);
    const [loadingEvents, setLoadingEvents] = useState(true);

    // Event modal state
    const [selectedEvent, setSelectedEvent] = useState(null);
    const [currentRSVP, setCurrentRSVP] = useState(null);

    // const eventDate = "2025-06-09T15:04:05Z";

    // checks for req. to join, pending, invited, accepted
    const isMember = group.status === "accepted";
    const isPending = group.status === "requested" ||  group.status === "invited";
    // const isInvited = group.status === "invited";

    // checks for the member info for modal Member Section
    const acceptedMembers = members.filter(m => m.status === 'accepted');
    const pendingRequests = requests.filter(r => r.status === 'requested' || r.status === 'invited');

    
    const currentUserUuid = localStorage.getItem('user-uuid');

    const memberUuids = new Set(members.map(m => m.follower_uuid || m.uuid));
    const nonMembers = allUsers.filter(u => !memberUuids.has(u.uuid));
    const filteredNonMembers = nonMembers.filter(u => u.uuid !== currentUserUuid);

    const refreshEvents = () => {
        fetch(`/frontend-api/groups/events/${group.uuid}`)
            .then(res => {
            if (!res.ok) {
                return res.text().then(text => {
                throw new Error(`HTTP error! Status: ${res.status}. Response: ${text}`);
                });
            }
            return res.json();
        })
        .then(data => setGroupEvents(data.data || []))
        .catch(err => {
            console.error("Failed to refresh events:", err);
            toast.error("Failed to load events");
        });
    };


    // Fetch events for this group
    useEffect(() => {
        if (!group?.uuid) return;
        setLoadingEvents(true);
        fetch(`/frontend-api/groups/events/${group.uuid}`)
            .then(res => res.json())
            .then(data => {
                setGroupEvents(data.data || []);
                setLoadingEvents(false);
            })
            .catch(() => setLoadingEvents(false));
    }, [group?.uuid]);

    const handlePostSubmit = (postData) => {
        // TODO: Add post to group (API or state update)
        toast.success(`Post created: ${postData.title}`);
        setShowPostModal(false);
    };

    const handleEventSubmit = async (eventData) => {
        console.log("Submitting eventData:", eventData);
        try {
            const res = await fetch('/frontend-api/events/create', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(eventData),
            });
            if (res.ok) {
                toast.success('Event created!');
                refreshEvents();
            } else {
                toast.error('Failed to create event.');
            }
        } catch (err) {
            toast.error('Network error.');
        }
        setShowEventModal(false);
    };

    // RSVP handler for modal
    const handleRSVP = async (status) => {
        console.log(selectedEvent)
        if (!selectedEvent) return;
        try {
            const res = await fetch('/frontend-api/groups/events/response', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    event_uuid: selectedEvent.uuid,
                    response: status,
                }),
            });
            let data = {};
            try {
                data = await res.json();
            } catch (e) {
                // If response is not JSON, data stays empty
            }
            if (res.ok) {
                setCurrentRSVP(status);
                toast.success(`RSVP updated: ${status === "accepted" ? "Going" : "Not Going"}`);
                refreshEvents();
            } else {
                console.error("Backend error:", data);
                toast.error('Failed to update RSVP.');
            }
        } catch (err) {
            console.error("Network error:", err);
            toast.error('Network error.');
        }
    };


    return (
        <div>
            <GroupHeader
                group={group}
                isMember={isMember}
                onShowPostModal={() => setShowPostModal(true)}
                onShowEventModal={() => setShowEventModal(true)}
                onRequestJoin={onRequestJoin}
            />
            <div className="group-detail">
                {onBack && (
                    <button onClick={onBack} className="group-detail-back-btn">
                        ← Back to Groups
                    </button>
                )}
                <section className="group-detail-members-section">
                    <div
                        className="group-members-label group-members-label--interactive"
                        onClick={handleOpenMembersModal}
                        tabIndex={0}
                        role="button"
                        aria-label="Show members"
                        onKeyDown={e => {
                            if (e.key === "Enter" || e.key === " ") handleOpenMembersModal();
                        }}
                        >
                        <strong>Members:</strong> {group.members_count}
                    </div>
                    {/* You can add more member info if backend provides it */}
                    <div className="group-members-actions">
                        {isMember ? (
                            <button onClick={handleOpenMembersModal}>Invite User</button>
                        ) : isPending ? (
                            <button disabled>Pending</button>
                        ) : (
                            <button onClick={onRequestJoin}>Request to Join</button>
                        )}
                    </div>
                </section>

                {/* Events Section */}
                <section className="event-section">
                    <strong>Upcoming Events:</strong>
                    {loadingEvents ? (
                        <div>Loading events...</div>
                    ) : groupEvents.length === 0 ? (
                        <div>No events yet.</div>
                    ) : (
                        <ul className="group-detail-events-list">
                            {groupEvents.map(event => (
                                <li key={event.uuid} className="group-detail-event-item">
                                    <EventCard
                                        event={event}
                                        isMember={isMember}
                                        onClick={() => {
                                            // OPEN THE MODAL with event details!
                                            setSelectedEvent(event);
                                            setCurrentRSVP(event.status); // If you have RSVP info in event
                                        }}
                                    />
                                </li>
                            ))}
                        </ul>
                    )}
                </section>
                <section className="post-section">
                    <strong>Posts:</strong>
                    {/* <PostCard>

                    </PostCard> */}
                </section>

                {/* Modals */}
                {showPostModal && (
                    <Modal onClose={() => setShowPostModal(false)}>
                        <CreatePostForm
                            groupId={group.uuid}
                            onSubmit={handlePostSubmit}
                            onClose={() => setShowPostModal(false)}
                        />
                    </Modal>
                )}
                {showEventModal && (
                    <Modal onClose={() => setShowEventModal(false)}>
                        <CreateEventForm
                            groupId={group.uuid}
                            onSubmit={handleEventSubmit}
                            onClose={() => setShowEventModal(false)}
                        />
                    </Modal>
                )}
                {/* THIS IS YOUR MEMBER DETAIL MODAL! */}
                {showMembersModal && (
                    <Modal title="Group Members" onClose={handleCloseMembersModal}>
                        <section>
                            <h3>Members</h3>
                            <ul>
                            {acceptedMembers.length > 0 ? (
                                    acceptedMembers.map(m => (
                                    <li key={m.follower_uuid || m.uuid}>{m.follower_name || m.name}</li>
                                    ))
                                ) : (
                                    <li>No members</li>
                                )}
                            </ul>
                        </section>
                        {isMember && (
                        <section>
                            <h3>Pending Members</h3>
                            <ul>
                            {pendingRequests.length > 0 ? (
                                pendingRequests.map(r => (
                                <li key={r.follower_uuid || r.uuid || r.user_uuid}>
                                    {r.follower_name || r.name}
                                    <button
                                    style={{ marginLeft: '1rem' }}
                                    onClick={() => onApproveRequest(r.follower_uuid || r.uuid || r.user_uuid)}
                                    >
                                    Approve
                                    </button>
                                    <button
                                    style={{ marginLeft: '0.5rem', color: 'red' }}
                                    onClick={() => onDenyRequest(r.follower_uuid || r.uuid || r.user_uuid)}
                                    >
                                    Deny
                                    </button>
                                </li>
                                ))
                            ) : (
                                <li>No pending members</li>
                            )}
                            </ul>
                        </section>
                        )}
                        {isMember && (
                        <section>
                        <h3>Non-Members</h3>
                        {loadingUsers ? (
                            <div>Loading users...</div>
                        ) : (
                            <ul>
                            {filteredNonMembers.length > 0 ? (
                                [...filteredNonMembers]
                                .sort((a, b) => a.nick_name.localeCompare(b.nick_name))
                                .map(u => (
                                    <li key={u.uuid}>
                                    {u.nick_name}
                                    <button
                                        style={{ marginLeft: '1rem' }}
                                        onClick={() => onInviteUser(u.uuid)}
                                    >
                                        Invite
                                    </button>
                                    </li>
                                ))
                            ) : (
                                <li>Everyone is a member!</li>
                            )}
                            </ul>
                        )}
                        </section>
                        )}
                    </Modal>
                )} 

                {/* THIS IS YOUR EVENT DETAIL MODAL! */}
                {selectedEvent && (
                    <Modal
                        onClose={() => {
                            setSelectedEvent(null);
                            setCurrentRSVP(null);
                        }}
                        title={selectedEvent.title}
                    >
                        <div>
                            <p><strong>Date/Time:</strong> {formatDate(selectedEvent.start_time)}</p>
                            {selectedEvent.location && (
                                <p><strong>Location:</strong> {selectedEvent.location}</p>
                            )}
                            <p><strong>Description:</strong> {selectedEvent.description}</p>
                            <p><strong>Attending:</strong> {selectedEvent.attend_count} going</p>
                            {isMember && (
                                <div style={{ marginTop: "1rem" }}>
                                    <button
                                        onClick={() => handleRSVP("accepted")}
                                        className={currentRSVP === "accepted" ? "active" : ""}
                                    >
                                        Going
                                    </button>
                                    <button
                                        onClick={() => handleRSVP("declined")}
                                        className={currentRSVP === "declined" ? "active" : ""}
                                    >
                                        Not Going
                                    </button>
                                </div>
                            )}
                        </div>
                    </Modal>
                )}
            </div>
        </div>
    );
}

