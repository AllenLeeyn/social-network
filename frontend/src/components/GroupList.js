"use client";

export default function GroupList({ groups, loading, error }) {
  if (loading) return <div>Loading...</div>;
  if (error) return <div>Error: {error}</div>;
  if (!groups || groups.length === 0)
    return <div>No groups found.</div>;

  return (
    <ul className="connections">
      {groups.map((g, idx) => (
        <li key={g.uuid || idx} className="connection-item">
          <span>
            <strong>
              {g.title}
            </strong>
          </span>
        </li>
      ))}
    </ul>
  );
}
