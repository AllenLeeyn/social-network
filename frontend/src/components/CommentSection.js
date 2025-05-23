// Renders the full comment section
// src/components/comments/CommentSection.js

import CommentForm from "./CommentForm";
import "../styles/Comments.css";
import { submitCommentFeedback } from "../lib/apiPosts";

export default function CommentSection({
  comments = [],
  title = "Comments",
  postId,
  onCommentSubmitted,
}) {
  // Sort comments by id (ascending)
  const sortedComments = [...comments].sort((a, b) => a.id - b.id);

  // Add handlers for like/dislike
  const handleFeedback = async (commentId, rating) => {
    try {
      await submitCommentFeedback({ parent_id: commentId, rating });
      if (onCommentSubmitted) onCommentSubmitted(); // Refresh comments
    } catch (err) {
      alert(err.message || "Failed to submit feedback");
    }
  };

  return (
    <div className="comments-section">
      <h4>{title}</h4>
      <ul className="comments-list">
        {sortedComments.length === 0 && <li>No comments yet.</li>}
        {sortedComments.map((comment) => (
          <li key={comment.id} className="comment-item">
            <div className="comment-author">
              <strong>{comment.user.nick_name}</strong>
            </div>
            <div className="comment-content">{comment.content}</div>
            <div className="comment-timestamp">
              {new Date(comment.created_at).toLocaleString()}
            </div>
            <div className="comment-actions">
              <button
                onClick={() => handleFeedback(comment.id, 1)}
                aria-label="Like"
                disabled={comment.liked}
              >
                👍 Like {comment.like_count}
              </button>
              <button
                onClick={() => handleFeedback(comment.id, -1)}
                aria-label="Dislike"
                style={{ marginLeft: "1em" }}
                disabled={comment.disliked}
              >
                👎 Dislike {comment.dislike_count}
              </button>
            </div>
          </li>
        ))}
      </ul>
      <CommentForm postId={postId} onCommentSubmitted={onCommentSubmitted} />
    </div>
  );
}
