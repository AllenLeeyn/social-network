"use client";

import React, { useState, useEffect } from "react";
import { fetchMyPosts } from "../../../lib/apiPosts";
import "./profile.css";
import SidebarSection from "../../../components/SidebarSection";
import PostCard from "../../../components/PostCard";
import ProfileCard from "../../../components/ProfileCard";
import UsersList from "../../../components/UsersList";
import {
    myPosts,
    myActivity,
    sampleGroups,
    sampleFollowers,
    sampleFollowing,
} from "../../../data/mockData";

export default function ProfilePage({ params }) {
    const [showFollowers, setShowFollowers] = useState(false);
    const [showFollowing, setShowFollowing] = useState(false);
    const [isPrivateProfile, setIsPrivateProfile] = useState(false);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    const [myPosts, setMyPosts] = useState([]);

    // Optionally, fetch posts only if not private
    useEffect(() => {
    if (isPrivateProfile) {
        setMyPosts([]);
        setLoading(false);
        return;
    }
    async function loadMyData() {
        try {
        const [myPostsData] = await Promise.all([fetchMyPosts()]);
        setMyPosts(myPostsData.data);
        } catch (err) {
        setError(err.message);
        } finally {
        setLoading(false);
        }
    }
    loadMyData();
    }, [isPrivateProfile]);

    return (
    <main>
        <div className="homepage-layout">
        {/* Left Sidebar */}
        <aside className="sidebar left-sidebar">
            <SidebarSection title="Privacy">
            <div className="privacy-toggle">
                <label htmlFor="privacySwitch">Visibility:</label>
                <select
                id="privacySwitch"
                value={isPrivateProfile ? "private" : "public"}
                onChange={(e) => setIsPrivateProfile(e.target.value === "private")}
                >
                <option value="public">Public</option>
                <option value="private">Private</option>
                </select>
            </div>
            </SidebarSection>
            <SidebarSection title="My Activity">
            <ul className="categories">
                {myActivity.map((cat) => (
                <li key={cat.id} className="category-item">
                    <strong>{cat.name}</strong>
                </li>
                ))}
            </ul>
            </SidebarSection>
            <SidebarSection title="My Groups">
            <ul className="groups">
                {sampleGroups.map((group) => (
                <li key={group.id} className="group-item">
                    <strong>{group.name}</strong>
                </li>
                ))}
            </ul>
            </SidebarSection>
        </aside>

        {/* Profile Content */}
        <section className="main-post-section">
            <ProfileCard uuid={params.uuid} setPrivateProfile={setIsPrivateProfile} />

          {/* Main Post Content */}
            {!isPrivateProfile && (
            <div className="user-posts-section">
                <h3>Posts</h3>
                {myPosts.length > 0 ? (
                myPosts.map((post) => <PostCard key={post.id} post={post} />)
                ) : (
                <p>No posts yet.</p>
                )}
            </div>
            )}

            {/* Followers Modal */}
            {showFollowers && (
            <div className="modal">
                <div className="modal-content">
                <h3>Followers</h3>
                <button
                    className="close"
                    onClick={() => setShowFollowers(false)}
                >
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
                <button
                    className="close"
                    onClick={() => setShowFollowing(false)}
                >
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
        </section>

        {/* Right Sidebar */}
        <aside className="sidebar right-sidebar">
            <UsersList />
        </aside>
        </div>
    </main>
    );
}
