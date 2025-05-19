"use client";

import React, { useState } from "react";
import "./profile.css";
import SidebarSection from "../../components/SidebarSection";
import { sampleUsers } from "../../data/mockData";

export default function ProfilePage() {
    const [showFollowers, setShowFollowers] = useState(false);
    const [showFollowing, setShowFollowing] = useState(false);

    const currentUser = {
    username: "UserZ",
    fullName: "Zoe Zimmerman",
    avatar: "/avatars/zoe.png",
    bio: "Frontend developer and design enthusiast.",
};

    const sampleFollowers = [
    {id: 1, username: "UserA", fullName: "Alice Anderson", avatar: "/avatars/alice.png"},
    {id: 2, username: "UserB", fullName: "Bob Brown", avatar: "/avatars/bob.png"},
    {id: 3, username: "UserC", fullName: "Charlie Clark", avatar: "/avatars/charlie.png"},
    ];

    const sampleFollowing = [
    {id: 4, username: "UserD", fullName: "David Davis", avatar: "/avatars/david.png"},
    {id: 5, username: "UserE", fullName: "Emma Evans", avatar: "/avatars/emma.png"},
    {id: 6, username: "UserF", fullName: "Frank Foster", avatar: "/avatars/frank.png"},
    ];


    return (
        <main>
        <div className="profile-header">
            <img
                src={currentUser.avatar}
                alt={currentUser.username}
                className="profile-avatar"
            />
                <div className="profile-info">
                    <h2>{currentUser.username}</h2>
                    <p>{currentUser.fullName}</p>
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

            {/* Followers Modal */}
            {showFollowers && (
            <div className="modal">
                <div className="modal-content">
                <h3>Followers</h3>
                <button className="close" onClick={() => setShowFollowers(false)}>
                    ✖
                </button>
                    <ul className="users">
                        {sampleFollowers.map((user) => (
                        <li key={user.id} className="user-item">
                            <img src={user.avatar} alt={user.username} />
                            <span>
                            {user.fullName} ({user.username})
                            </span>
                        </li>
                        ))}
                    </ul>
                </div>
            </div>
            )}

            {/* Following Modal */}
            {showFollowing && (
            <div className="modal">
                <div className="modal-content">
                <h3>Following</h3>
                <button className="close" onClick={() => setShowFollowing(false)}>
                    ✖
                </button>
                    <ul className="users">
                        {sampleFollowing.map((user) => (
                        <li key={user.id} className="user-item">
                            <img src={user.avatar} alt={user.username} />
                            <span>
                            {user.fullName} ({user.username})
                            </span>
                        </li>
                        ))}
                    </ul>
                </div>
            </div>
            )}

            <div>
                <aside className="sidebar right-sidebar">
                    <SidebarSection title="Active Users">
                        <ul className="users">
                        {sampleUsers.map(user => (
                            <li key={user.id} className={`user-item${user.online ? " online" : ""}${user.unread ? " unread" : ""}`}>
                            <img src={user.avatar} alt={user.username} />
                            <span>{user.fullName} ({user.username})</span>
                            </li>
                        ))}
                        </ul>
                    </SidebarSection>
                </aside>
            </div>
        </main>
    );
}
