// src/components/groups/GroupDetail.js

import React from "react";
import { mockEvents } from "../../data/mockData";

// Props:
// - group: the group object to display (required)
// - onBack: function to call when the user wants to go back to the list (optional)
// - currentUser: string (for permissions, optional)

export default function GroupDetail({ group, onBack, currentUser = "alice" }) {
    if (!group) return null;

    // Filter events for this group
    const groupEvents = mockEvents.filter(e => e.groupId === group.id);

    // Determine if current user is a member
    const isMember = group.members.includes(currentUser);

    return (
        <div className="group-detail">
        <button onClick={onBack} style={{ marginBottom: 16 }}>
            ← Back to Groups
        </button>
        <h2>{group.title}</h2>
        <p>{group.description}</p>

        <section style={{ margin: "16px 0" }}>
            <strong>Members:</strong>
            <ul>
            {group.members.map(member => (
                <li key={member}>{member}</li>
            ))}
            </ul>
            {!isMember && (
            <button style={{ marginTop: 8 }}>
                Request to Join
            </button>
            )}
            {isMember && (
            <button style={{ marginTop: 8 }}>
                Invite User
            </button>
            )}
        </section>

        <section>
            <strong>Upcoming Events:</strong>
            {groupEvents.length === 0 ? (
            <div>No events yet.</div>
            ) : (
            <ul>
                {groupEvents.map(event => (
                <li key={event.id} style={{ marginBottom: 12 }}>
                    <div>
                    <b>{event.title}</b> — {new Date(event.dateTime).toLocaleString()}
                    </div>
                    <div style={{ fontSize: "0.95em" }}>{event.description}</div>
                    {isMember && (
                    <div style={{ marginTop: 4 }}>
                        <button style={{ marginRight: 8 }}>Going</button>
                        <button>Not Going</button>
                    </div>
                    )}
                </li>
                ))}
            </ul>
            )}
            {isMember && (
            <button style={{ marginTop: 16 }}>
                Create Event
            </button>
            )}
        </section>
        </div>
    );
}
