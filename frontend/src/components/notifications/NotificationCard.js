'use client';

import React from 'react';
import Link from "next/link";
import Image from "next/image";
import { FaUserCircle } from "react-icons/fa";
import { TimeAgo } from '../../utils/TimeAgo';
import { readNotification } from "../../lib/apiNotifications";
import { submitFollowResponse } from "../../lib/apiFollow";
import { submitGroupRequestOrInviteResponse } from "../../lib/apiGroups";

export default function NotificationCard({ notification }) {
    const acceptRejectActions = ['follow_request', 'group_invite', 'group_request'];
    const viewActions = ['follow_request_accepted', 'group_event'];

    const handleNotificationRespond = async (status) => {
        try {
            await readNotification({ id: notification.id, is_read: 1 });

            if (notification.target_detailed_type === 'follow_request') {
                await submitFollowResponse({ follower_uuid: notification.target_uuid, status });
            } else if (notification.target_detailed_type === 'group_invite') {
                await submitGroupRequestOrInviteResponse({ follower_uuid: notification.to_user_uuid, group_uuid: notification.target_uuid, status });
            } else if (notification.target_detailed_type === 'group_request') {
                await submitGroupRequestOrInviteResponse({ follower_uuid: notification.from_user.uuid, group_uuid: notification.target_uuid, status });
            }
        } catch (err) {
            toast.error(err.message || "Failed to respond notification");
        }
    };
    

    return (
        <div className={`notification-card${notification.is_read ? '' : ' unread'}`}>
        <div className="notification-content">
            {notification.from_user && (
                <div className="notification-user">
                    {notification.from_user.avatar ? (
                        <Image
                            src={notification.from_user.avatar}
                            alt={notification.from_user.nick_name}
                            width={50}
                            height={50}
                            style={{ borderRadius: "50%" }}
                            />
                    ) : (
                        <FaUserCircle
                        size={50}
                        color="#aaa"
                        style={{ verticalAlign: "middle" }}
                        />
                    )}
                    <br />
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
                { acceptRejectActions.includes(notification.target_detailed_type) && (
                    !notification.is_read ? (
                        <div>
                            <button onClick={() => handleNotificationRespond('accepted')}>Accept</button>
                            <button onClick={() => handleNotificationRespond('declined')}>Decline</button>
                        </div>    
                    ) : (
                        notification.target_detailed_type === 'follow_request' ?
                        (<Link href={`/profile/${notification.to_user_uuid}`} className="link-btn">View</Link>) :
                        notification.target_detailed_type === 'group_invite' ?
                        (<Link href={`/groups/${notification.target_uuid}`} className="link-btn">View</Link>) :
                        notification.target_detailed_type === 'group_request' ?
                        (<Link href={`/groups/${notification.target_uuid}`} className="link-btn">View</Link>
                        ) : null
                    )
                    
                )}
                {notification.target_detailed_type === 'follow_request_accepted' && (
                    <Link href={`/profile/${notification.to_user_uuid}`} className="link-btn">View</Link>
                )}
                {notification.target_detailed_type === 'group_event' && (
                    <Link href={`/groups/${group.uuid}`} className="link-btn">View</Link>
                )}
                
            </div>
            </div>
        </div>
        </div>
    );
}
