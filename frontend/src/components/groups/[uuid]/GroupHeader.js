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

    const handleButtonClick = (e) => {
        e.preventDefault();
        if (onRequestJoin) onRequestJoin(group);
    };

    let actionTaken = null;
    if (group.status === "requested") {
        actionTaken = <span className="group-detail-pending">Pending</span>;
    } else if (group.status === "invited") {
        actionTaken = <span className="group-detail-invited">Invited</span>;
    } else if (!isMember) {
        actionTaken = (
            <button onClick={handleButtonClick}>Request to Join</button>
        );
    }

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
                ) : (
                    actionTaken
                )}
            </div>
        </div>
    );
}
