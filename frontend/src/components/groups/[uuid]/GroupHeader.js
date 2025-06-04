// src/components/groups/GroupHeader.js

import React from "react";

export default function GroupHeader({
    group,
    isMember,
    onShowPostModal,
    onShowEventModal,
}) {
    if (!group) return null;

    return (
        <div className="group-detail-header">
        <h2>Welcome to {group.title}</h2>
        <p>{group.description}</p>
        <div className="group-detail-actions">
            {isMember ? (
            <>
                <button onClick={onShowPostModal}>Create Post</button>
                <button onClick={onShowEventModal}>Create Event</button>
            </>
            ) : null}
        </div>
        </div>
    );
}
