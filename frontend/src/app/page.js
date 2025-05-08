// This is our homepage
import React from 'react';

export default function HomePage() {

// Posts
const samplePosts = [
{ id: 1, title: 'Post Title 1', author: 'UserA', snippet: 'This is a snippet of the first post...' },
{ id: 2, title: 'Post Title 2', author: 'UserB', snippet: 'This is a snippet of the second post...', online: true, unread: true },
{ id: 3, title: 'Post Title 3', author: 'UserC', snippet: 'This is a snippet of the third post...' },
];

// Categories
const sampleCategories = [
{ id: 1, name: "Technology" },
{ id: 2, name: "Health" },
{ id: 3, name: "Travel" },
];

// Users
const sampleUsers = [
{ id: 1, username: "UserA", fullName: "Alice Anderson", avatar: "/avatars/alice.png", online: true },
{ id: 2, username: "UserB", fullName: "Bob Brown", avatar: "/avatars/bob.png", online: true, unread: true },
{ id: 3, username: "UserC", fullName: "Charlie Clark", avatar: "/avatars/charlie.png" },
];

return (
    <main>
        <h1>Welcome to grit:Hub!</h1>
        <p>Your place to connect and share.</p>
        <div style={{ marginBottom: '1rem' }}>
            <button style={{ marginRight: '1rem' }}>Sign Up</button>
            <button>Log In</button>
        </div>
        
        <div className="homepage-sections"> 
            <section className="sidebar categories-section">
                <h2>Categories</h2>
                {/* Replace with <Categories categories={sampleCategories} /> when ready */}
                <ul className="categories">
                {sampleCategories.map(cat => (
                    <li key={cat.id} className={`category-item${cat.active ? " active" : ""}`}>
                        <strong>{cat.name}</strong>
                    </li>
                ))}
                </ul>
            </section>
            <section className="main-feed post-list-section">
            <h2>Latest Posts</h2>
                {/* Replace with <PostList posts={samplePosts} /> when ready */}
                <div>
                    {samplePosts.map(post => (
                    <div key={post.id} className="post-item">
                        <h3>{post.title}</h3>
                        <p><em>by {post.author}</em></p>
                        <p>{post.snippet}</p>
                    </div>
                    ))}
                </div>
            </section>
            <section className="sidebar users-section">
                <h2>Active Users</h2>
                {/* Replace with <UserList users={sampleUsers} /> when ready */}
                <ul className="users">
                {sampleUsers.map(user => (
                    <li key={user.id} className={`user-item${user.online ? " online" : ""}${user.unread ? " unread" : ""}`}>
                        <img src={user.avatar} alt={user.username} />
                        <span>{user.fullName} ({user.username})</span>
                    </li>
                ))}
                </ul>   
            </section>
        </div>
    </main>
    );
}

