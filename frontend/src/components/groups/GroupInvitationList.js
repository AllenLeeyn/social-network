
// src/components/groups/GroupInvitationList.js

import React, { useState } from "react";
import { mockInvitations } from "../../data/mockData";

// Props:
// - currentUser: string (the logged-in user)

export default function GroupInvitationList({ currentUser = "frank" }) {
    // Filter invitations for the current user and pending status
    const [invitations, setInvitations] = useState(
        mockInvitations.filter(inv => inv.toUser === currentUser && inv.status === "pending")
    );

    // Handler for accepting/declining
    const handleRespond = (id, response) => {
        setInvitations(prev =>
        prev.map(inv =>
            inv.id === id ? { ...inv, status: response } : inv
        )
        );
        // Here you would also call your backend to update invitation status
    };

    if (invitations.length === 0) {
        return <div>No group invitations at this time.</div>;
    }

    return (
        <div>
        <h3>Group Invitations</h3>
        <ul style={{ listStyle: "none", padding: 0 }}>
            {invitations.map(inv => (
            <li key={inv.id} style={{ marginBottom: 16, padding: 12, background: "#f7fafc", borderRadius: 6 }}>
                <div>
                <b>{inv.fromUser}</b> invited you to join <b>{inv.groupTitle}</b>
                </div>
                <div style={{ marginTop: 8 }}>
                {inv.status === "pending" ? (
                    <>
                    <button
                        onClick={() => handleRespond(inv.id, "accepted")}
                        style={{ marginRight: 8 }}
                    >
                        Accept
                    </button>
                    <button
                        onClick={() => handleRespond(inv.id, "declined")}
                    >
                        Decline
                    </button>
                    </>
                ) : (
                    <span style={{ color: inv.status === "accepted" ? "green" : "red" }}>
                    {inv.status === "accepted" ? "Accepted" : "Declined"}
                    </span>
                )}
                </div>
            </li>
            ))}
        </ul>
        </div>
    );
}
