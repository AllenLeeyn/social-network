'use client';

import React from 'react';

export default function NotificationCard({ notification }) {
    return (
        <div className={`notification-card${notification.isRead ? '' : ' unread'}`}>
        <div className="notification-content">
            {notification.fromUser && (
            <img
                src={notification.fromUser.avatar}
                alt={notification.fromUser.name}
                className="notification-avatar"
            />
            )}
            <div>
            <div className="notification-message">{notification.message}</div>
            <div className="notification-timestamp">
                {new Date(notification.timestamp).toLocaleString()}
            </div>
            <div className="notification-actions">
                {notification.actions.includes('accept') && <button>Accept</button>}
                {notification.actions.includes('decline') && <button>Decline</button>}
                {notification.actions.includes('view') && <button>View</button>}
            </div>
            </div>
        </div>
        </div>
    );
}
