"use client";
import Link from "next/link";
import "../styles/PostList.css";


export default function PostsPage({ posts = [] }) {
    return (
        <div className="post-list">
            <ul>
                {posts.map((post) => (
                    <li className="post-item" key={post.id}>
                        <h2>
                        <Link href={`/post?id=${post.uuid}`}>{post.title}</Link>
                        </h2>
                        <p>{post.content}</p>
                        <small>By {post.user.nick_name}</small>
                    </li>
                ))}
            </ul>
        </div>
    );
}