// Handles new comment submissions
"use client";
import { useRef, useState } from "react";
import { submitComment } from "../lib/apiPosts";
import { handleImage } from "../lib/handleImage"; 

import "../styles/Comments.css";

export default function CommentForm({ postUUID, onCommentSubmitted }) {
  const [content, setContent] = useState("");
  const [postImage, setImage] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const fileInputRef = useRef(null);
  const [previewUrl, setPreviewUrl] = useState(null);

  const handleFileChange = (e) => {
    const file = e.target.files[0];
    if (!file) return;
    setImage(e.target.files[0])

    const url = URL.createObjectURL(file);
    setPreviewUrl(url);
  };

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

      await submitComment({ post_uuid: postUUID, content,
        attached_image: imageUUID ? Object.values(imageUUID)[0] : null });
      setContent("");

      if (fileInputRef.current) {
        fileInputRef.current.value = '';
      }
      setImage(null)
      setPreviewUrl(null);
      if (onCommentSubmitted) onCommentSubmitted();

    } catch (err) {
      setError(err.message || "Failed to submit comment");
    } finally {
      setLoading(false);
    }
  };

  return (
    <form className="comment-form" onSubmit={handleSubmit}>
      <h3>Feel free to comment!</h3>
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
        onChange={handleFileChange}
        ref={fileInputRef}
      />
      {previewUrl && (
        <img
          src={previewUrl}
          alt="Image Preview"
          style={{ width: 100, height: 100, objectFit: "cover", marginTop: 10 }}
        />
      )}
      <button type="submit" disabled={loading}>
        {loading ? "Posting..." : "Post Comment"}
      </button>
      {error && <div className="error">{error}</div>}
    </form>
  );
}
