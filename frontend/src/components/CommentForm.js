// Handles new comment submissions
"use client";
import { useState } from "react";
import { submitComment } from "../lib/apiPosts";

import "../styles/Comments.css";

export default function CommentForm({ postId, onCommentSubmitted }) {
  const [content, setContent] = useState("");
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError(null);
    try {
      await submitComment({ post_id: Number(postId), content });
      setContent("");
      if (onCommentSubmitted) onCommentSubmitted();
    } catch (err) {
      setError(err.message || "Failed to submit comment");
    } finally {
      setLoading(false);
    }
  };

  return (
    <form className="comment-form" onSubmit={handleSubmit}>
      <h2>Feel free to comment!</h2>
      <textarea
        placeholder="Add a comment..."
        value={content}
        onChange={(e) => setContent(e.target.value)}
        required
        disabled={loading}
      />
      <button type="submit" disabled={loading}>
        {loading ? "Posting..." : "Post Comment"}
      </button>
      {error && <div className="error">{error}</div>}
    </form>
  );
}
