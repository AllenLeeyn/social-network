// Handles new comment submissions
"use client";
import { useState } from "react";
import { submitComment } from "../lib/apiPosts";
import { handleImage } from "../lib/handleImage"; 

import "../styles/Comments.css";

export default function CommentForm({ postId, onCommentSubmitted }) {
  const [content, setContent] = useState("");
  const [postImage, setImage] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError(null);
    try {

      let imageUUID = null;
      if (postImage) {
        try {
          imageUUID = await handleImage([postImage]);
        } catch (err) {
          setError("Image upload failed: " + err.message);
          return;
        }
      }

      await submitComment({ post_id: Number(postId), content,
        attached_image: imageUUID ? Object.values(imageUUID)[0] : null });
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
      <input
        type="file"
        accept="image/*"
        onChange={(e) => setImage(e.target.files[0])}
      />
      <button type="submit" disabled={loading}>
        {loading ? "Posting..." : "Post Comment"}
      </button>
      {error && <div className="error">{error}</div>}
    </form>
  );
}
