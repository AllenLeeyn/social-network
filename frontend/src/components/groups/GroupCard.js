'use client'

import '../../styles/groups/GroupCard.css'

export default function GroupCard({ group, onInvite, onRequestJoin }) {

    return (
        <div className="group-card">
        <h3>{group.title}</h3>
        <p>{group.description}</p>
        <p>Members: {group.members_count}</p>
        <p>Creator: {group.creator_name}</p>
        {/* Conditional Button */}
        {group.status === "accepted" ? (
            <button onClick={() => onInvite(group)}>Invite</button>
        ) : (
            <button onClick={() => onRequestJoin(group)}>Request to Join</button>
        )}
        </div>
    );
}
