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
      {connections.map((conn) => {
        const isLeader = conn.leader_name === currentUser;
        const nameToShow = isLeader ? conn.follower_name : conn.leader_name;
        const uuidToLink = isLeader ? conn.follower_uuid : conn.leader_uuid;

        return (
        <li key={uuidToLink} className="connection-item">
          <span>
            <strong>
              {nameToShow}
            </strong>
          </span>
        </li>
      )
      })}
    </ul>
  );
}
