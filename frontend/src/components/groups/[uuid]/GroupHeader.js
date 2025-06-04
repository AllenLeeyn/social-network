// src/components/groups/GroupHeader.js

import React from "react";

export default function GroupHeader({
    group,
    isMember,
    onShowPostModal,
    onShowEventModal,
    onRequestJoin,
}) {
    if (!group) return null;


    const isPending = group.status === "requested";

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
            ) : isPending ? (
                <button disabled>Pending</button>
            ) : (
                <button onClick={onRequestJoin}>Request to Join</button>
            )}
        </div>
        </div>
    );
}
