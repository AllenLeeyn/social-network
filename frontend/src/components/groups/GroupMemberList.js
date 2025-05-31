// src/components/groups/GroupMembersList.js
import React from "react";

export default function GroupMembersList({ members }) {
    if (!members || members.length === 0) {
        return <div>No members yet.</div>;
    }
    return (
        <ul className="group-members-list">
        {members.map(member => (
            <li key={member.uuid}>
            {member.name || member.username}
            </li>
        ))}
        </ul>
    );
}
