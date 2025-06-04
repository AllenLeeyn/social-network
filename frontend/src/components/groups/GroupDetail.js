// src/components/groups/GroupDetail.js
// groups/uuid/main-view, main-feed

import React, { useState, useEffect } from "react";
import { toast } from 'react-toastify';

import Modal from "../Modal";
import CreatePostForm from "./[uuid]/posts/CreatePostForm";
import CreateEventForm from "./[uuid]/events/CreateEventForm";
import GroupHeader from "./[uuid]/GroupHeader";
import EventCard from "./[uuid]/events/EventCard";
import { createPost } from "../../lib/apiPosts";
import PostCard from '../PostCard'

import { formatDate } from '../../utils/formatDate';

import "../../styles/groups/GroupDetail.css"; 


export default function GroupDetail({ group, onBack, onRequestJoin }) {
    console.log(group.id)
    if (!group) return null;

    // Modal state
    const [showPostModal, setShowPostModal] = useState(false);
    const [showEventModal, setShowEventModal] = useState(false);

    // Events state
    const [groupEvents, setGroupEvents] = useState([]);
    const [loadingEvents, setLoadingEvents] = useState(true);

    // Event modal state
    const [selectedEvent, setSelectedEvent] = useState(null);
    const [currentRSVP, setCurrentRSVP] = useState(null);

    // const eventDate = "2025-06-09T15:04:05Z";

    const isMember = group.status === "accepted";

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

    const handlePostSubmit = async (postData) => {                
        try {
            const data = await createPost(postData);
            if (data) {
                window.location.href = `/post/${data.data}`;
            } else {
                toast.error(data.message || "Failed to create post");
            }
        } catch (err) {
            toast.error(err.message || "Error creating post");
        }
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
                <section className="member-section">
                    <strong>Members:</strong>
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
                            groupID={group.id}
                            onSubmit={handlePostSubmit}
                            onClose={() => setShowPostModal(false)}
                        />
                    </Modal>
                )}
                {showEventModal && (
                    <Modal onClose={() => setShowEventModal(false)}>
                        <CreateEventForm
                            groupId={group.id}
                            onSubmit={handleEventSubmit}
                            onClose={() => setShowEventModal(false)}
                        />
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

