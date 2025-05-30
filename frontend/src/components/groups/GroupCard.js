'use client'

import Link from 'next/link';
import '../../styles/groups/GroupCard.css'

export default function GroupCard({ group, onInvite, onRequestJoin }) {

    // Helper to prevent navigation when button is clicked
    const handleButtonClick = (e, handler) => {
        e.preventDefault();
        e.stopPropagation();
        handler(group);
    };

    return (
        <Link href={`/groups/${group.uuid}`} className="group-card-link" style={{ textDecoration: 'none', color: 'inherit' }}>
            <div className="group-card">
            <h3>{group.title}</h3>
            <p>{group.description}</p>
            <p>Members: {group.members_count}</p>
            <p>Creator: {group.creator_name}</p>
            {/* Conditional Button */}
            {group.status === "accepted" ? (
                    <button onClick={e => handleButtonClick(e, onInvite)}>Invite</button>
                ) : (
                    <button onClick={e => handleButtonClick(e, onRequestJoin)}>Request to Join</button>
                )}
            </div>
        </Link>
    );
}
