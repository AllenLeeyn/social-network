// Renders the full comment section
// src/components/comments/CommentSection.js

import CommentForm from "./CommentForm";
import CommentCard from "./commentCard";
import "../styles/Comments.css";

export default function CommentSection({
  comments = [],
  title = "Comments",
  postId,
  onCommentSubmitted,
}) {
  // Sort comments by id (ascending)
  const sortedComments = [...comments].sort((a, b) => a.id - b.id);


  return (
    <div className="comments-section">
      <h4>{title}</h4>
      <ul className="comments-list">
        {sortedComments.length === 0 && <li>No comments yet.</li>}
        {sortedComments.map((comment) => (
          <li key={comment.id}>
            <CommentCard comment={comment} />
          </li>
        ))}
      </ul>
      <CommentForm postId={postId} onCommentSubmitted={onCommentSubmitted} />
    </div>
  );
}
