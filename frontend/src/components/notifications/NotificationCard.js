'use client';

import React from 'react';

export default function NotificationCard({ notification }) {
    const acceptRejectActions = ['follow_request', 'group_invite', 'group_request'];
    const viewActions = ['follow_request_accepted', 'group_event'];
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
                {acceptRejectActions.includes(notification.target_detailed_type) && <button>Accept</button>}
                {acceptRejectActions.includes(notification.target_detailed_type) && <button>Decline</button>}
                {viewActions.includes(notification.target_detailed_type) && <button>View</button>}
            </div>
            </div>
        </div>
        </div>
    );
}
