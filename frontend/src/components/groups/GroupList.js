// src/components/groups/GroupList.js

import React from "react";
import { mockGroups } from "../../data/mockData";

// Props:
// - type: "my_groups" | "discover"
// - onSelectGroup: function(group) => void
// - currentUser: string (optional, for filtering "my groups")

export default function GroupList({ type = "my_groups", onSelectGroup, currentUser = "alice" }) {
    // Filter groups based on type
    let groupsToShow = mockGroups;
    if (type === "my_groups") {
        groupsToShow = mockGroups.filter(g => g.members.includes(currentUser));
    } else if (type === "discover") {
        groupsToShow = mockGroups.filter(g => !g.members.includes(currentUser));
    }

    if (groupsToShow.length === 0) {
        return <div className="group-list-empty">No groups to display.</div>;
    }

    return (
        <ul className="group-list">
        {groupsToShow.map(group => (
            <li
            key={group.id}
            className="group-list-item"
            onClick={() => onSelectGroup(group)}
            tabIndex={0}
            role="button"
            style={{
                cursor: "pointer",
                padding: "12px",
                borderBottom: "1px solid #ececec",
                background: "#fff",
                borderRadius: "4px",
                marginBottom: "8px",
                boxShadow: "0 1px 2px rgba(0,0,0,0.02)"
            }}
            >
            <div>
                <strong>{group.title}</strong>
            </div>
            <div style={{ color: "#666", fontSize: "0.95em" }}>{group.description}</div>
            <div style={{ color: "#888", fontSize: "0.85em" }}>
                Members: {group.members.length}
            </div>
            </li>
        ))}
        </ul>
    );
}
