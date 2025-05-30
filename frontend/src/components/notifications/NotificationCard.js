'use client';

import React from 'react';
import { TimeAgo } from '../../utils/formatDate';
import { readNotification } from "../../lib/apiNotifications";
import { submitFollowResponse } from "../../lib/apiFollow";
import { submitGroupRequestOrInviteResponse } from "../../lib/apiGroups";

export default function NotificationCard({ notification }) {
    const acceptRejectActions = ['follow_request', 'group_invite', 'group_request'];
    const viewActions = ['follow_request_accepted', 'group_event'];

    const handleNotificationFeedback = async (status) => {
        try {
            await readNotification({ id: notification.id, is_read: 1 });

            if (notification.target_detailed_type === 'follow_request') {
                await submitFollowResponse({ follower_uuid: notification.target_uuid, status });
            } else if (notification.target_detailed_type === 'group_invite' || notification.target_detailed_type === 'group_request') {
                await submitGroupRequestOrInviteResponse({ follower_uuid: notification.target_uuid, group_uuid: notification.target_uuid, status });
            } else if (notification.target_detailed_type === 'follow_request_accepted') {
                //todo go to the user's profile /profile/uuid
            } else if (notification.target_detailed_type === 'group_event') {
                //todo go to the group event page /group/event/uuid
            }
        } catch (err) {
            toast.error(err.message || "Failed to submit feedback");
        }
    };
    

    return (
        <div className={`notification-card${notification.is_read ? '' : ' unread'}`}>
        <div className="notification-content">
            {notification.from_user && (
                <div className="notification-user">
                    {notification.from_user.avatar && (
                        <img
                            src={notification.from_user.avatar}
                            alt={notification.from_user.nick_name}
                            className="notification-avatar"
                        />
                    )}
                    <span className="notification-nickname">
                        {notification.from_user.nick_name}
                    </span>
                </div>
            )}

            <div>
            <div className="notification-message">{notification.message}</div>
            <div className="notification-timestamp">
                {TimeAgo(notification.created_at)}
            </div>
            <div className="notification-actions">
                {!notification.is_read && acceptRejectActions.includes(notification.target_detailed_type) && <button onClick={() => handleNotificationFeedback('accepted')}>Accept</button>}
                {!notification.is_read && acceptRejectActions.includes(notification.target_detailed_type) && <button onClick={() => handleNotificationFeedback('declined')}>Decline</button>}
                {!notification.is_read && viewActions.includes(notification.target_detailed_type) && <button onClick={() => handleNotificationFeedback('view')}>View</button>}
            </div>
            </div>
        </div>
        </div>
    );
}
