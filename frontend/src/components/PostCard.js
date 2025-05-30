'use client';

import { useState } from 'react';
import Link from 'next/link';
import { FaThumbsUp, FaThumbsDown, FaCommentAlt } from 'react-icons/fa';
import { TimeAgo } from '../utils/TimeAgo';
import { toast } from 'react-toastify';
import { fetchPostById, submitPostFeedback } from "../lib/apiPosts";

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
    if (liked) {
      setLiked(false);
      handlePostFeedback(0)
      setLikeCount((c) => c - 1);
    } else {
      setLiked(true);
      handlePostFeedback(1)
      setLikeCount((c) => c + 1);
      if (disliked) {
        setDisliked(false);
        setDislikeCount((c) => c - 1);
      }
      // optionally call API to like
    }
  };

  const handleDislike = () => {
    if (disliked) {
      setDisliked(false);
      handlePostFeedback(0)
      setDislikeCount((c) => c - 1);
      // optionally call API to undo dislike
    } else {
      setDisliked(true);
      handlePostFeedback(-1)
      setDislikeCount((c) => c + 1);
      if (liked) {
        setLiked(false);
        setLikeCount((c) => c - 1);
      }
      // optionally call API to dislike
    }
  };

  return (
    <div>
      <h2>
        <Link href={`/post?id=${post.uuid}`}>{post.title}</Link>
      </h2>

      <small>
        By{' '}
        <Link href={`/user/${post.user.uuid}`}>
          {post.user.nick_name}
        </Link>{' '}
        [{TimeAgo(post.created_at)}]
      </small>

      <pre>{post.content}</pre>

      <div className="post-stats">
        <button
          className={`stat-btn ${liked ? 'liked' : ''}`}
          onClick={handleLike}
        >
          <FaThumbsUp /> {likeCount}
        </button>

        <button
          className={`stat-btn ${disliked ? 'disliked' : ''}`}
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
      </small>
    </div>
  );
}
