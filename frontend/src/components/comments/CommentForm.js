// Handles new comment submissions
'use clients';
import { useState } from "react";

import "./Comments.css"

export default function CommentForm() {
    const [author, setAuthor] = useState("");
    const [content, setContent] = useState("");

    const handleSubmit = (e) => {
        e.preventDefault();
        // For now, just log the comment. In a real app, you'd update state or call an API.
        console.log("New comment:", { author, content, timestamp: new Date().toISOString() });
        setAuthor("");
        setContent("");
    };

    return (
    <form className="comment-form" onSubmit={handleSubmit}>
        <h2>Feel free to comment!</h2>
        <textarea
            placeholder="Add a comment..."
            value={content}
            onChange={e => setContent(e.target.value)}
            required
        />
        <button type="submit">Post Comment</button>
        </form>
    );
}
