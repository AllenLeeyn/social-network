// src/components/groups/GroupDetail.js
// groups/uuid/main-view, main-feed

import React, { useState } from "react";
import { toast } from 'react-toastify';

import Modal from "../Modal";
import CreatePostForm from "./posts/CreatePostForm";
import CreateEventForm from "./events/CreateEventForm";
import GroupHeader from "./GroupHeader";

import { formatDate } from '../../utils/formatDate';

import "../../styles/groups/GroupDetail.css"; 


export default function GroupDetail({ group, onBack }) {

    // Modal state
    const [showPostModal, setShowPostModal] = useState(false);
    const [showEventModal, setShowEventModal] = useState(false);

    const eventDate = "2025-06-09T15:04:05Z";

    if (!group) return null;

    const isMember = group.status === "accepted";

    // PH
    const groupEvents = [];

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
                // Optionally: refresh event list here
            } else {
                toast.error('Failed to create event.');
            }
        } catch (err) {
            toast.error('Network error.');
        }
        setShowEventModal(false);
    };


    return (
        <div>
            <GroupHeader
                group={group}
                isMember={isMember}
                onShowPostModal={() => setShowPostModal(true)}
                onShowEventModal={() => setShowEventModal(true)}
                onRequestJoin={() => {/* handle join logic */}}
                onInviteUser={() => {/* handle invite logic */}}
            />
            <div className="group-detail">
                {onBack && (
                    <button onClick={onBack} className="group-detail-back-btn">
                        ← Back to Groups
                    </button>
                )}
                <section className="group-detail-members-section">
                    <strong>Members:</strong> {group.members_count}
                    {/* You can add more member info if backend provides it */}
                    {!isMember && (
                        <button className="group-detail-join-btn">
                            Request to Join
                        </button>
                    )}
                    {isMember && (
                        <button className="group-detail-invite-btn">
                            Invite User
                        </button>
                    )}
                </section>

                <section>
                    <strong>Upcoming Events:</strong>
                    {groupEvents.length === 0 ? (
                        <div>No events yet.</div>
                    ) : (
                        <ul className="group-detail-events-list">
                            {groupEvents.map(event => (
                                <li key={event.id} className="group-detail-event-item">
                                    <div>
                                        <b>{event.title}</b> — {formatDate(event.start_time || event.dateTime)}
                                    </div>
                                    <div className="group-detail-event-desc">{event.description}</div>
                                    {isMember && (
                                        <div className="group-detail-event-actions">
                                            <button>Going</button>
                                            <button>Not Going</button>
                                        </div>
                                    )}
                                </li>
                            ))}
                        </ul>
                    )}
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
            </div>
                        
        </div>
    );
}

