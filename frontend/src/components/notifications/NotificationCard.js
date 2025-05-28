'use client';

import React from 'react';

export default function NotificationCard({ notification }) {
    return (
        <div className={`notification-card${notification.is_read ? '' : ' unread'}`}>
        <div className="notification-content">
            {notification.from_user && (
            <img
                src={notification.from_user.avatar}
                alt={notification.from_user.first_name}
                className="notification-avatar"
            />
            )}
            <div>
            <div className="notification-message">{notification.message}</div>
            <div className="notification-timestamp">
                {new Date(notification.created_at).toLocaleString()}
            </div>
            <div className="notification-actions">
                {/*notification.actions.includes('accept') &&*/ <button>Accept</button>}
                {/*notification.actions.includes('decline') &&*/ <button>Decline</button>}
                {/*notification.actions.includes('view') &&*/ <button>View</button>}
            </div>
            </div>
        </div>
        </div>
    );
}
