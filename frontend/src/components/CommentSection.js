// Renders the full comment section
// src/components/comments/CommentSection.js

import CommentForm from './CommentForm';
import '../styles/Comments.css'; 

export default function CommentSection({ comments = [], title = "Comments" }) {
    return (
        <div className="comments-section">
        <h4>{title}</h4>
        <ul className="comments-list">
            {comments.length === 0 && <li>No comments yet.</li>}
            {comments.map(comment => (
            <li key={comment.id} className="comment-item">
                <div className="comment-author"><strong>{comment.user.nick_name}</strong></div>
                <div className="comment-content">{comment.content}</div>
                <div className="comment-timestamp">{new Date(comment.created_at).toLocaleString()}</div>
                <div className="comment-actions">
                    <label>
                        <input
                            type="radio"
                            name={`like-dislike-${comment.id}`}
                            checked={comment.liked}
                            readOnly
                        />
                        👍 Like {comment.like_count}
                    </label>
                    <label style={{ marginLeft: '1em' }}>
                        <input
                            type="radio"
                            name={`like-dislike-${comment.id}`}
                            checked={comment.disliked}
                            readOnly
                        />
                        👎 Dislike {comment.dislike_count}
                    </label>
                </div>
            </li>
            ))}
        </ul>
        <CommentForm />
        </div>
    );
}
