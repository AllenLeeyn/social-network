"use client"; // Required to use hooks like useParams

import { useParams } from 'next/navigation';
import SidebarSection from '../components/sidebarSection/SidebarSection';
import CommentsSection from '../components/comments/CommentSection';
import "../styles/globals.css";

import {
    samplePosts,
    sampleCategories,
    sampleUsers,
    sampleGroups,
    sampleConnections,
    sampleComments,
} from '../data/mockData';


export default function PostPage() {
    const params = useParams();
    const postId = params.id;

    // Find the post by ID (convert postId to number)
    // const post = samplePosts.find(p => p.id === Number(postId));
    const post = postId
    ? samplePosts.find(p => p.id === Number(postId))
    : samplePosts[1];

    const commentsForThisPost = post
    ? sampleComments.filter(comment => comment.postId === post.id)
    : [];


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

                {/* Main Post Content */}
                <section className="main-post-section">
                    {post ? (
                        <div key={post.id} className="post-item">
                            <h3>{post.title}</h3>
                            <p><em>by {post.author}</em></p>
                            <p>{post.snippet}</p>
                        </div>
                    ) : (
                        <div className="post-item">
                            <h3>Post not found</h3>
                        </div>
                    )}

                <CommentsSection title="Comments" comments={commentsForThisPost}/>


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
