// This is our homepage
import React from 'react';
import SidebarSection from '../components/SidebarSection/SidebarSection';
import PostList from '../components/PostList';

import {
    samplePosts,
    sampleCategories,
    sampleUsers,
    sampleGroups,
    sampleConnections
} from '../data/mockData';


export default function HomePage() {

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
                <section className="main-feed post-list-section">
                    <h2>Latest Posts</h2>
                    <PostList posts={samplePosts} />
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


