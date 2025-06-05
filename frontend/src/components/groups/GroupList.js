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
    fetch("/frontend-api/groups")
      .then((res) => res.json())
      .then((data) => {
        setGroups(data.data);
        setLoading(false);
      });
  }, []);

  if (!groups || groups.length === 0) return <div>No groups found.</div>;

  // req
  let filteredGroups = groups || [];
  if (filter === "my_groups") {
    filteredGroups = groups.filter((g) => g.status === "accepted");
  } else if (filter === "discover") {
    filteredGroups = groups.filter((g) => g.status !== "accepted");
  }

  function handleRequestJoin(group) {
    fetch("/frontend-api/groups/join", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ group_uuid: group.uuid }),
    })
      .then((res) => res.json())
      .then((data) => {
        if (data.success) {
          setGroups((prevGroups) =>
            prevGroups.map((g) =>
              g.uuid === group.uuid ? { ...g, status: "requested" } : g
            )
          );

          toast.success("Request sent!");
        } else {
          toast.error(data.error || "Request failed.");
        }
      })
      .catch(() => toast.error("Network error."));
    toast.success(`Request to join "${group.title}" sent!`);
  }

  return (
    <div className="group-list">
      {loading && <p>Loading groups...</p>}
      {!loading && (!filteredGroups || filteredGroups.length === 0) && (
        <p>No groups found.</p>
      )}
      {!loading &&
        Array.isArray(filteredGroups) &&
        filteredGroups.map((group) => (
          <GroupCard
            key={group.uuid}
            group={group}
            onRequestJoin={handleRequestJoin}
          />
        ))}
    </div>
  );
}
