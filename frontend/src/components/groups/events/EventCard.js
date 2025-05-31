// src/components/groups/events/EventCard.js

import React from "react";
import { formatDate } from '../../../utils/formatDate';

export default function EventCard({ event, onClick }) {
    return (
        <div className="event-card" onClick={onClick} style={{ cursor: 'pointer' }}>
            <b>{event.title}</b>
            <span> — {formatDate(event.start_time)}</span>
        </div>
    );
}
