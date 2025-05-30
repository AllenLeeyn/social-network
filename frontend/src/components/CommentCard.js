'use client';

import { useState } from 'react';
import Link from 'next/link';
import { FaThumbsUp, FaThumbsDown } from 'react-icons/fa';
import { TimeAgo } from '../utils/formatDate';
import { toast } from 'react-toastify';
import { submitCommentFeedback } from "../lib/apiPosts";

export default function CommentCard({ comment }) {
  const [liked, setLiked] = useState(comment.liked);
  const [disliked, setDisliked] = useState(comment.disliked);
  const [likeCount, setLikeCount] = useState(comment.like_count);
  const [dislikeCount, setDislikeCount] = useState(comment.dislike_count);

  const handleCommentFeedback = async (rating) => {
    try {
      await submitCommentFeedback({ parent_id: comment.id, rating });
    } catch (err) {
      toast.error(err.message || "Failed to submit feedback");
    }
  };

  const handleLike = () => {
    if (liked) {
      setLiked(false);
      handleCommentFeedback(0)
      setLikeCount((c) => c - 1);
    } else {
      setLiked(true);
      handleCommentFeedback(1)
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
      handleCommentFeedback(0)
      setDislikeCount((c) => c - 1);
      // optionally call API to undo dislike
    } else {
      setDisliked(true);
      handleCommentFeedback(-1)
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
      <pre>{comment.content}</pre>

      <small>
        By{' '}
        <Link href={`/user/${comment.user.uuid}`}>
          {comment.user.nick_name}
        </Link>{' '}
        [{TimeAgo(comment.created_at)}]
      </small>

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
      </div>
    </div>
  );
}
