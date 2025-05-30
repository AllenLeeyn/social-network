// src/components/groups/GroupDetail.js

import React, { useState } from "react";
import { toast } from 'react-toastify';

import Modal from "../Modal";
import CreatePostForm from "./CreatePostForm";
import CreateEventForm from "./CreateEventForm";

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

    const handleEventSubmit = (eventData) => {
        // TODO: Add event to group (API or state update)
        toast.success(`Event created: ${eventData.title}`);
        setShowEventModal(false);
    };


    return (
        <div>
            <div className="group-detail-header">
                <h2>Welcome to {group.title}</h2>
                <p>{group.description}</p>
                {/* Action buttons for members */}
                {isMember && (
                    <div className="group-detail-actions">
                        <button onClick={() => setShowPostModal(true)}>
                            Create Post
                        </button>
                        <button onClick={() => setShowEventModal(true)}>
                            Create Event
                        </button>
                    </div>
                )}
            </div>
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

