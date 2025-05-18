"use client"; 

// we just need the post.id to use this method
import { useParams } from 'next/navigation';
import { useEffect, useState } from "react";
import SidebarSection from '../../../components/SidebarSection';
import CommentsSection from '../../../components/CommentSection';
import "./post.css";

import {
    samplePosts,
    sampleCategories,
    sampleUsers,
    sampleGroups,
    sampleConnections,
    sampleComments,
} from '../../../data/mockData';


export default function PostPage() {
const { id } = useParams();
    const [post, setPost] = useState(null);
    const [comments, setComments] = useState([]);
    const [categories, setCategories] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    useEffect(() => {
        async function fetchData() {
            try {
                // fetch the individual post
                const postRes = await fetch(`/api/posts/${id}`);
                if (!postRes.ok) throw new Error("Post not found");
                const postData = await postRes.json();
                setPost(postData);

                // fetch comments for this post
                const commentsRes = await fetch(`/api/comments?postId=${id}`);
                if (commentsRes.ok) {
                    setComments(await commentsRes.json());
                }

                // Optionally, fetch categories
                const catRes = await fetch("/api/posts");
                if (catRes.ok) {
                    setCategories(await catRes.json());
                }
            } catch (err) {
                setError(err.message);
            } finally {
                setLoading(false);
            }
        }
        fetchData();
    }, [id]);

    if (loading) return <div>Loading...</div>;
    if (error) return <div>Error: {error}</div>;

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
