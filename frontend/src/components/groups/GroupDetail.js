// src/components/groups/GroupDetail.js

import React, { useState } from "react";
import { mockEvents } from "../../data/mockData";

import Modal from "../Modal";
import CreatePostForm from "./CreatePostForm";
import CreateEventForm from "./CreateEventForm";

import "../../styles/groups/GroupDetail.css"; 


export default function GroupDetail({ group, onBack, currentUser = "alice" }) {

    // Modal state
    const [showPostModal, setShowPostModal] = useState(false);
    const [showEventModal, setShowEventModal] = useState(false);

    if (!group) return null;

    const groupEvents = mockEvents.filter(e => e.groupId === group.id);
    const isMember = group.members.includes(currentUser);

    const handlePostSubmit = (postData) => {
        // TODO: Add post to group (API or state update)
        alert(`Post created: ${postData.title}`);
        setShowPostModal(false);
    };

    const handleEventSubmit = (eventData) => {
        // TODO: Add event to group (API or state update)
        alert(`Event created: ${eventData.title}`);
        setShowEventModal(false);
    };


    return (
            <div className="group-detail">
                <button onClick={onBack} className="group-detail-back-btn">
                    ← Back to Groups
                </button>
                <h2>{group.title}</h2>
                <p>{group.description}</p>

                <section className="group-detail-members-section">
                    <strong>Members:</strong>
                    <ul className="group-detail-members-list">
                        {group.members.map(member => (
                            <li key={member}>{member}</li>
                        ))}
                    </ul>
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

                <section>
                    <strong>Upcoming Events:</strong>
                    {groupEvents.length === 0 ? (
                        <div>No events yet.</div>
                    ) : (
                        <ul className="group-detail-events-list">
                            {groupEvents.map(event => (
                                <li key={event.id} className="group-detail-event-item">
                                    <div>
                                        <b>{event.title}</b> — {new Date(event.dateTime).toLocaleString()}
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
                            groupId={group.id}
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
            </div>
        );
}

