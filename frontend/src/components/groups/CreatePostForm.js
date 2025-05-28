import React, { useState } from "react";

export default function CreatePostForm({ groupId, onSubmit, onClose }) {
    const [title, setTitle] = useState("");
    const [content, setContent] = useState("");

    const handleSubmit = (e) => {
        e.preventDefault();
        // You can add validation here if needed
        onSubmit && onSubmit({ groupId, title, content });
        onClose && onClose();
    };

    return (
        <form className="create-post-form" onSubmit={handleSubmit}>
        <h3>Create Post</h3>
        <label>
            Title:
            <input
            type="text"
            value={title}
            onChange={e => setTitle(e.target.value)}
            required
            />
        </label>
        <label>
            Content:
            <textarea
            value={content}
            onChange={e => setContent(e.target.value)}
            required
            />
        </label>
        <div className="form-actions">
            <button type="submit">Post</button>
            <button type="button" onClick={onClose}>Cancel</button>
        </div>
        </form>
    );
}
