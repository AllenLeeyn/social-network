"use client";
import { usePosts } from "../hooks/usePosts";

import "../styles/globals.css";


export default function PostsPage() {
    const { posts, error, loading } = usePosts();

    if (loading) return <div>Loading...</div>;
    if (error) return <div>Error: {error}</div>;

    return (
        <div className="post-list">
            <h1>Posts</h1>
            <ul>
                {posts.map((post) => (
                <li key={post.id}>
                    <h2>{post.title}</h2>
                    <p>{post.content || post.snippet}</p>
                    <small>By {post.userName || post.author}</small>
                </li>
                ))}
            </ul>
        </div>
        );
}
