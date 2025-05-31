import React from "react";

export default function EventCard({ event }) {
    return (
        <div className="event-card">
        <h3>{event.title}</h3>
        <div>
            {event.location} — {event.start_time}
        </div>
        <div>{event.description}</div>
        </div>
    );
}