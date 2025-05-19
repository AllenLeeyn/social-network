'use client';

import React,  { useState } from 'react';
import SidebarSection from '../components/SidebarSection';
import PostList from '../components/PostList';
import CreatePost from '../components/CreatePost'
import Modal from '../components/Modal'


import {
    samplePosts,
    sampleCategories,
    sampleUsers,
    sampleGroups,
    sampleConnections
} from '../data/mockData';


export default function HomePage() {

    const [showModal, setShowModal] = useState(false);

return (
        <main>
            <div className="homepage-layout">
                {/* Left Sidebar */}
                <aside className="sidebar left-sidebar">
                    <SidebarSection title="Categories">
                        <ul className="categories">
                            {sampleCategories.map(cat => (
                            <li key={cat.id} className="category-item">
                                <strong>{cat.name}</strong>
                            </li>
                            ))}
                        </ul>
                    </SidebarSection>
                    <SidebarSection title="Groups">
                        <ul className="groups">
                        {sampleGroups.map(group => (
                            <li key={group.id} className="group-item">
                            <strong>{group.name}</strong>
                            </li>
                        ))}
                        </ul>
                    </SidebarSection>
                    <SidebarSection title="Connections">
                        <ul className="connections">
                        {sampleConnections.map(conn => (
                            <li key={conn.id} className="connection-item">
                            <span><strong>{conn.fullName} ({conn.username})</strong></span>
                            </li>
                        ))}
                        </ul>
                    </SidebarSection>
                </aside>

                {/* Center / Main View */}
                <section className="main-feed post-list-section">
                    <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
                        <h2>Latest Posts</h2>
                        <button
                            className="create-post-btn"
                            onClick={() => setShowModal(true)}
                            aria-label="Create a new post"
                        >
                            + Create Post
                        </button>
                    </div>
                        <PostList posts={samplePosts} />

                        {showModal && (
                            <Modal onClose={() => setShowModal(false)} title="Create Post">
                            {/* Post creation form goes here */}
                            <CreatePost categories={sampleCategories} onClose={() => setShowModal(false)} />
                            </Modal>
                        )}
                </section>

                {/* Right Sidebar */}
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


