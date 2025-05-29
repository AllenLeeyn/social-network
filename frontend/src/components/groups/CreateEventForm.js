import React, { useState } from "react";

export default function CreateEventForm({ groupId, onSubmit, onClose }) {
    const [title, setTitle] = useState("");
    const [description, setDescription] = useState("");
    const [dateTime, setDateTime] = useState("");

    const handleSubmit = (e) => {
        e.preventDefault();
        onSubmit && onSubmit({ groupId, title, description, dateTime });
        onClose && onClose();
    };

    return (
        <form className="create-event-form" onSubmit={handleSubmit}>
        <h3>Create Event</h3>
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
            Description:
            <textarea
            value={description}
            onChange={e => setDescription(e.target.value)}
            required
            />
        </label>
        <label>
            Date & Time:
            <input
            type="datetime-local"
            value={dateTime}
            onChange={e => setDateTime(e.target.value)}
            required
            />
        </label>
        <div className="form-actions">
            <button type="submit">Create Event</button>
            <button type="button" onClick={onClose}>Cancel</button>
        </div>
        </form>
    );
}
