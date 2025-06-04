'use client';

import React from 'react';
import NotificationCard from './NotificationCard';

export default function NotificationList({ notifications }) {
    if (!notifications.length) {
        return (
            <div className='notification-list-empty'>
                <p>No notifications found.</p>
            </div>
        );
    }
    return (
        <div className="notification-list">
        {notifications.map(notification => (
            <NotificationCard 
                key={notification.id} 
                notification={notification}
            />
        ))}
        </div>
    );
}
