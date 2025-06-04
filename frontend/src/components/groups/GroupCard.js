'use client'

import Link from 'next/link';
import '../../styles/groups/GroupCard.css'

export default function GroupCard({ group, onRequestJoin }) {

    // Helper to prevent navigation when button is clicked
    const handleButtonClick = (e, handler) => {
        e.preventDefault();
        e.stopPropagation();
        onRequestJoin(group);
    };

    let actionTaken = null;
    if (group.status === "invited") {
        actionTaken = <span className="group-card-invited">Invited</span>;
    } else if (group.status === "requested") {
        actionTaken = <span className="group-card-pending">Pending</span>;
    } else if (group.status === "") {
        actionTaken = (
            <button onClick={handleButtonClick}>
                Request to Join
            </button>
        );
    }

    return (
        <Link href={`/groups/${group.uuid}`} className="group-card-link" style={{ textDecoration: 'none', color: 'inherit' }}>
            <div className="group-card">
            <h3>{group.title}</h3>
            <p>{group.description}</p>
            <p>Members: {group.members_count}</p>
            <p>Creator: {group.creator_name}</p>
            {actionTaken}
            </div>
        </Link>
    );
}
