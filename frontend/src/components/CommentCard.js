'use client';

import { useState } from 'react';
import Link from 'next/link';
import { FaThumbsUp, FaThumbsDown } from 'react-icons/fa';
import { TimeAgo } from '../utils/TimeAgo';
import { toast } from 'react-toastify';
import DynamicImage from './DynamicImage';
import { submitCommentFeedback } from "../lib/apiPosts";
import { FaUserCircle } from 'react-icons/fa';
import Image from 'next/image';

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
      handleCommentFeedback(1)
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
      handleCommentFeedback(-1)
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
      <pre>{comment.content}</pre>
      {comment.attached_image && (
        <div className="post-images">
            <DynamicImage
              key={comment.attached_image}
              src={`/frontend-api/image/${comment.attached_image}`}
              alt={`comment attachment`}
            />
        </div>
      )}

      <div className="user-info">
        <div className="user-avatar">
        {comment.user.profile_image ? (
          <Image
            src={`/frontend-api/image/${comment.user.profile_image}`}
            alt="User Avatar"
            width={40}
            height={40}
          />
        ) : (
          <FaUserCircle size={30} color="#aaa"/>
        )}
        </div>
        <div className="user-details">
          <Link href={`/user/${comment.user.uuid}`} className="user-name">
            {comment.user.nick_name}
          </Link>
          <div className="timestamp">{TimeAgo(comment.created_at)}</div>
        </div>
      </div>

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
