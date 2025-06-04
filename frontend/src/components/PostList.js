"use client";
import "../styles/PostList.css";
import PostCard from './PostCard';

export default function PostsPage({ posts = [] }) {
    return (
        <div className="post-list">
            <ul>
                {Array.isArray(posts) && posts.length > 0 ? (
                    posts.map((post) => (
                        <li className="post-item" key={post.id}>
                            <PostCard post={post} />
                        </li>
                    ))
                ) : (
                    <li>No posts to display.</li>
                )}
            </ul>
        </div>
    );
}