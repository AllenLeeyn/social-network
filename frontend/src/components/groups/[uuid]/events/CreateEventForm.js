import React, { useState } from "react";
import { getMinDateTime } from '../../../../utils/getMinDateTime'; // adjust path as needed

export default function CreateEventForm({ groupUUID, onSubmit, onClose }) {
    const [title, setTitle] = useState("");
    const [description, setDescription] = useState("");
    const [dateTime, setDateTime] = useState("");
    const [location, setLocation] = useState("");
    const [duration, setDuration] = useState(60); // default 60 minutes

    const handleSubmit = (e) => {
        e.preventDefault();

        // convert dateTime to ISO string for backend
        const start_time = new Date(dateTime).toISOString();

        onSubmit && onSubmit({ 
            group_uuid: groupUUID, 
            title, 
            description,
            location,
            duration_minutes: Number(duration), 
            start_time, 
        });
        onClose && onClose();
    };

    return (
        <form className="create-event-form" onSubmit={handleSubmit}>
            <h3>Create Event</h3>
            <label>
                <h4>Title</h4>
                <input
                    type="text"
                    value={title}
                    onChange={e => setTitle(e.target.value)}
                    required
                />
            </label>
            <label>
                <h4>Description</h4>
                <textarea
                    value={description}
                    onChange={e => setDescription(e.target.value)}
                    required
                />
            </label>
            <label>
                <h4>Location</h4>
                <input
                    type="text"
                    value={location}
                    onChange={e => setLocation(e.target.value)}
                    required
                />
            </label>
            <label>
                <h4>Date & Time</h4>
                <input
                    type="datetime-local"
                    value={dateTime}
                    onChange={e => setDateTime(e.target.value)}
                    required
                    min={getMinDateTime()}
                />
            </label>
            <label>
                <h4>Duration (minutes)</h4>
                <input
                    type="number"
                    min={1}
                    value={duration}
                    onChange={e => setDuration(e.target.value)}
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
