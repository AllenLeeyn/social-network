// src/components/groups/GroupMembersList.js
import React from "react";

export default function GroupMembersList({ 
    members = [], 
    requests = [], 
    groupUuid, 
    onApproveRequest,
    onDenyRequest, 
}) {
    const accepted = members.filter(m => m.status === 'accepted');
    const invited = members.filter(m => m.status === 'invited');

    if (accepted.length === 0 && invited.length === 0 && requests.length === 0) {
        return <div>No members yet.</div>;
    }

    return (
        <div>
            <ul>
                {accepted.map((member, idx) => (
                    <li key={member.follower_uuid || idx}>
                        {member.follower_name}
                    </li>
                ))}
            </ul>
            {requests.length > 0 && <h4>Pending Requests</h4>}
            <ul>
                {requests.map((member, idx) => (
                    <li key={member.follower_uuid || idx}>
                        {member.follower_name}
                        <button
                            style={{ marginLeft: '1rem' }}
                            onClick={() => onApproveRequest(member.follower_uuid)}
                        >
                            Approve
                        </button>
                        <button
                            style={{ marginLeft: '0.5rem', color: 'red' }}
                            onClick={() => onDenyRequest(member.follower_uuid)}
                        >
                            Deny
                        </button>
                    </li>
                ))}
            </ul>
        </div>
    );
}
