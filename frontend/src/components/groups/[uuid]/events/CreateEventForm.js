import React, { useState } from "react";
import { getMinDateTime } from '../../../../utils/getMinDateTime'; // adjust path as needed

export default function CreateEventForm({ groupId, onSubmit, onClose }) {
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
            group_uuid: groupId, 
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
                Location:
                <input
                    type="text"
                    value={location}
                    onChange={e => setLocation(e.target.value)}
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
                    min={getMinDateTime()}
                />
            </label>
            <label>
                Duration (minutes):
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
