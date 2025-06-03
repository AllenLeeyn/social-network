// components/UserListPanel.js
'use client';
import { useState } from 'react';
import Link from 'next/link';

export default function FollowingsList({ title, users, displayProperty, linkProperty }) {
    console.log(users)
  const [isOpen, setIsOpen] = useState(false);
    const count = Array.isArray(users) ? users.length : 0;

  return (
    <div className="user-list-panel">
      <button
        className="toggle-btn"
        onClick={() => setIsOpen(!isOpen)}
        aria-expanded={isOpen}
      >
        {`${title} (${count})`}
      </button>

      {isOpen && (
        <div className="panel-content">
          {isOpen && (
            <div className="panel-content">
                {count === 0 ? (
                <p>No {title.toLowerCase()} found.</p>
                ) : (
                <ul className="users">
                    {users.map((user, index) => (
                    <li key={`${title}-${index}`} className="user-item">
                        <Link href={`/profile/${user[linkProperty]}`} className="user-link">
                        <span>{user[displayProperty]}</span>
                        </Link>
                    </li>
                    ))}
                </ul>
                )}
            </div>
            )}
        </div>
      )}
    </div>
  );
}
