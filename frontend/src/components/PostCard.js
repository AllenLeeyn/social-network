"use client";

import { useState } from "react";
import Link from "next/link";
import { FaThumbsUp, FaThumbsDown, FaCommentAlt } from "react-icons/fa";
import { TimeAgo } from "../utils/TimeAgo";
import DynamicImage from "./DynamicImage";
import { toast } from "react-toastify";
import { fetchPostById, submitPostFeedback } from "../lib/apiPosts";
import { FaUserCircle } from "react-icons/fa";
import Image from "next/image";
import "../styles/PostList.css";



export default function PostCard({ post }) {
  const [liked, setLiked] = useState(post.liked);
  const [disliked, setDisliked] = useState(post.disliked);
  const [likeCount, setLikeCount] = useState(post.like_count);
  const [dislikeCount, setDislikeCount] = useState(post.dislike_count);

  const handlePostFeedback = async (rating) => {
    try {
      await submitPostFeedback({ parent_id: post.id, rating });
    } catch (err) {
      toast.error(err.message || "Failed to submit feedback");
    }
  };

  const handleLike = () => {
    handlePostFeedback(1);
    if (liked) {
      setLiked(false);
      setLikeCount((c) => c - 1);
    } else {
      setLiked(true);
      setLikeCount((c) => c + 1);
      if (disliked) {
        setDisliked(false);
        setDislikeCount((c) => c - 1);
      }
    }
  };

  const handleDislike = () => {
    handlePostFeedback(-1);
    if (disliked) {
      setDisliked(false);
      setDislikeCount((c) => c - 1);
    } else {
      setDisliked(true);
      setDislikeCount((c) => c + 1);
      if (liked) {
        setLiked(false);
        setLikeCount((c) => c - 1);
      }
    }
  };

  return (
    <div>
      <h2>
        <Link href={`/post/${post.uuid}`}>{post.title}</Link>
      </h2>

      <div className="user-info">
        <div className="user-avatar">
          {post.user.profile_image ? (
            <Image
              src={`/frontend-api/image/${post.user.profile_image}`}
              alt="User Avatar"
              width={40}
              height={40}
            />
          ) : (
            <FaUserCircle size={50} color="#aaa" />
          )}
        </div>
        <div className="user-details">
          <Link href={`/profile/${post.user.uuid}`} className="user-name">
            {post.user.nick_name}
          </Link>
          <div className="timestamp">{TimeAgo(post.created_at)}</div>
        </div>
      </div>

      <pre>{post.content}</pre>

      {post.post_files && post.post_files.length > 0 && (
        <div className="post-images">
          {post.post_files.map((file, index) => (
            <DynamicImage
              key={index}
              src={`/frontend-api/image/${file.file_uploaded_name}`}
              alt={`Post attachment ${index + 1}`}
            />
          ))}
        </div>
      )}

      <div className="post-stats">
        <button
          className={`stat-btn ${liked ? "liked" : ""}`}
          onClick={handleLike}
        >
          <FaThumbsUp /> {likeCount}
        </button>

        <button
          className={`stat-btn ${disliked ? "disliked" : ""}`}
          onClick={handleDislike}
        >
          <FaThumbsDown /> {dislikeCount}
        </button>

        <button className="stat-btn" aria-disabled="true">
          <FaCommentAlt /> {post.comment_count}
        </button>
      </div>

      <small className="category-tags">
        {post.categories.map((cat) => (
          <span key={cat.id || cat.name} className="category-badge">
            {cat.name}
          </span>
        ))}
        {post.group?.title && (
          <span className="category-badge">group: {post.group.title}</span>
        )}
      </small>
    </div>
  );
}
