// src/components/groups/GroupMembersList.js
import React from "react";

export default function GroupMembersList({ members }) {
    console.log(members.map(m => m.uuid));
    if (!members || members.length === 0) {
        return <div>No members yet.</div>;
    }
    return (
        <ul className="group-members-list">
        {members.map((member, index) => (
            <li key={member.follower_uuid || index}>
                {member.follower_name}
            </li>
        ))}
        </ul>
    );
}   
