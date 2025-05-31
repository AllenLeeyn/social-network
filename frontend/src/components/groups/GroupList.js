import React, { useState, useEffect } from "react";
import GroupCard from "./GroupCard"; 
import { toast } from "react-toastify";

// Props:
// - filter: "my_groups" | "discover"
// - onSelectGroup: function(group) => void

export default function GroupList({ filter, onSelectGroup }) {
    const [groups, setGroups] = useState([]);
    const [loading, setLoading] = useState(true);

    useEffect(() => {
        fetch('/frontend-api/groups')
        .then(res => res.json())
        .then(data => {
            setGroups(data.data); 
            setLoading(false);
        });
    }, []);

    // req 
    let filteredGroups = groups;
    if (filter === 'my_groups') {
        filteredGroups = groups.filter(g => g.status === "accepted");
    } else if (filter === 'discover') {
        filteredGroups = groups.filter(g => g.status !== "accepted");
    }


    function handleRequestJoin(group) {
        // TODO: Implement request to join logic
    toast.success(`Request to join "${group.title}" sent!`);
    }


    return (
        <div className="group-list">
        {loading && <p>Loading groups...</p>}
        {!loading && (!filteredGroups || filteredGroups.length === 0)  && <p>No groups found.</p>}
        {!loading && Array.isArray(filteredGroups) && filteredGroups.map(group => (
            <GroupCard
                key={group.uuid}
                group={group}
                onRequestJoin={handleRequestJoin}
            />
        ))}
        </div>
    );
}
