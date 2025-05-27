'use client';

import React from 'react';
import NotificationCard from './NotificationCard';

export default function NotificationList({ notifications }) {
    if (!notifications.length) {
        return <p>No notifications found.</p>;
    }
    return (
        <div className="notification-list">
        {notifications.map(notification => (
            <NotificationCard key={notification.id} notification={notification} />
        ))}
        </div>
    );
}
