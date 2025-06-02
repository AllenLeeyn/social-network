"use client";

import React, { useState, useEffect } from "react";
import { fetchMyPosts } from "../../lib/apiPosts";
import "./profile.css";
import SidebarSection from "../../components/SidebarSection";
import PostCard from "../../components/PostCard";
import {
  myPosts,
  myActivity,
  sampleGroups,
  sampleFollowers,
  sampleFollowing,
} from "../../data/mockData";

export default function ProfilePage() {
  const [showFollowers, setShowFollowers] = useState(false);
  const [showFollowing, setShowFollowing] = useState(false);
  const [isPrivate, setIsPrivate] = useState(false);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  const [myPosts, setMyPosts] = useState([]);

  const currentUser = {
    id: 99,
    username: "UserA",
    fullName: "Allen Lee",
    email: "allen.lee@grytlab.sg",
    dateOfBirth: "2003-01-01",
    avatar: "/avatars/allen.png",
    bio: "Backend developer and cartographer.",
    email: "allen.lee@grytlab.sg",
  };

  useEffect(() => {
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
  }, []);

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
                value={isPrivate ? "private" : "public"}
                onChange={(e) => setIsPrivate(e.target.value === "private")}
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

          {/* Main Post Content */}
          <div className="user-posts-section">
            <h3>{currentUser.username}'s Posts</h3>
            {myPosts.length > 0 ? (
              myPosts.map((post) => <PostCard key={post.id} post={post} />) 
            ) : (
              <p>No posts yet.</p>
            )}
          </div>

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
          <SidebarSection title="Active Users">
            <ul className="users">
              {sampleFollowers.map((user) => (
                <li
                  key={user.id}
                  className={`user-item${user.online ? " online" : ""}${
                    user.unread ? " unread" : ""
                  }`}
                >
                  <img src={user.avatar} alt={user.username} />
                  <span>
                    {user.fullName} ({user.username})
                  </span>
                </li>
              ))}
            </ul>
          </SidebarSection>
        </aside>
      </div>
    </main>
  );
}
