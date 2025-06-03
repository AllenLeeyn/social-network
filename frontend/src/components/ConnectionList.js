"use client";
import { useEffect, useState } from 'react';

export default function ConnectionList({ connections, loading, error }) {
  const [currentUser, setCurrentUser] = useState(null);

  useEffect(() => {
    if (typeof window !== 'undefined') {
      const nickname = localStorage.getItem('user-nick_name');
      setCurrentUser(nickname);
    }
  }, []);

  if (!currentUser) return null;
  if (loading) return <div>Loading...</div>;
  if (error) return <div>Error: {error}</div>;
  if (!connections || connections.length === 0)
    return <div>No connections found.</div>;

  return (
    <ul className="connections">
      {connections
        .filter(conn => conn.follower_name === currentUser)
        .map(conn => (
          <li key={conn.leader_uuid} className="connection-item">
            <span>
              <strong>{conn.leader_name}</strong>
            </span>
          </li>
        ))}
    </ul>
  );
}
