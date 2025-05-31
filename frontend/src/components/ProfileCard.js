"use client";

import { useEffect, useState } from 'react';

export default function ProfileCard() {
    const [user, setUser] = useState(null);

    useEffect(() => {
        async function fetchProfile() {
        try {
            const res = await fetch('/api/profile', {
            credentials: 'include', // if cookies/session are used
            });
            if (!res.ok) throw new Error('Failed to fetch user data');
            const data = await res.json();
            setUser(data);
        } catch (err) {
            console.error(err);
        }
    }

    fetchProfile();
    }, []);

    if (!user) return <div>Loading profile...</div>;

    return (
        <div className="profile-header">
                <img
                    src={currentUser.avatar}
                    alt={currentUser.username}
                    className="profile-avatar"
                />
                <div className="profile-info">
                    <h2>{currentUser.username}</h2>
                    <p>
                    <strong>{currentUser.fullName}</strong>
                    </p>
                    <p>{currentUser.email}</p>
                    <p>
                    <span>Date of Birth:</span> {currentUser.dateOfBirth}
                    </p>
                    <p className="bio">{currentUser.bio}</p>
                    <div className="connection-buttons">
                    <button onClick={() => setShowFollowers(true)}>
                        Followers ({sampleFollowers.length})
                    </button>
                    <button onClick={() => setShowFollowing(true)}>
                        Following ({sampleFollowing.length})
                    </button>
                    </div>
                </div>
                </div>
    );
}
