"use client";
import "../styles/globals.css";


export default function PostsPage( {posts}) {

    return (
        <div className="post-list">
            <h1>Posts</h1>
            <ul>
                {posts.map((post) => (
                <li key={post.id}>
                    <h2>{post.title}</h2>
                    <p>{post.content}</p>
                    <small>By {post.userName}</small>
                </li>
                ))}
            </ul>
        </div>
        );
}
