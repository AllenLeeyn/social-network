"use client";

export default function ConnectionList({ connections, loading, error }) {
  if (loading) return <div>Loading...</div>;
  if (error) return <div>Error: {error}</div>;
  if (!connections || connections.length === 0)
    return <div>No connections found.</div>;

  return (
    <ul className="connections">
      {connections.map((conn, idx) => (
        <li key={conn.leader_uuid || idx} className="connection-item">
          <span>
            <strong>
              {conn.leader_name}
            </strong>
          </span>
        </li>
      ))}
    </ul>
  );
}
