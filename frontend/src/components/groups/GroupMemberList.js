// src/components/groups/GroupMembersList.js
import React from "react";

export default function GroupMembersList({ members = [], requests = [] }) {
    const accepted = members.filter(m => m.status === 'accepted');
    const invited = members.filter(m => m.status === 'invited');

    if (accepted.length === 0 && invited.length === 0 && requests.length === 0) {
        return <div>No members yet.</div>;
    }

    return (
        <div>
            <h4>Current Members</h4>
            <ul>
                {accepted.map((member, idx) => (
                    <li key={member.follower_uuid || idx}>{member.follower_name}</li>
                ))}
            </ul>
            <h4>Pending Requests</h4>
            <ul>
                {requests.map((member, idx) => (
                    <li key={member.follower_uuid || idx}>{member.follower_name}</li>
                ))}
            </ul>
        </div>
    );
}
