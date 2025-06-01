// src/components/groups/GroupMembersList.js
import React from "react";

export default function GroupMembersList({ members }) {
    if (!members || members.length === 0) {
        return <div>No members yet.</div>;
    }
    const accepted = members.filter(m => m.status === 'accepted');
    const invited = members.filter(m => m.status === 'invited');
    const requested = members.filter(m => m.status === 'requested');

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
                {requested.map((member, idx) => (
                    <li key={member.follower_uuid || idx}>{member.follower_name}</li>
                ))}
            </ul>
            <h4>Invited</h4>
            <ul>
                {invited.map((member, idx) => (
                    <li key={member.follower_uuid || idx}>{member.follower_name}</li>
                ))}
            </ul>
        </div>
    );
}   
