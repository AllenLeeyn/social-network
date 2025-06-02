'use client';

import Image from "next/image";
import DOMPurify from 'dompurify';
import { FaUserCircle } from "react-icons/fa";
import { TimeAgo } from '../../utils/TimeAgo';
import { useNotifications } from '../../contexts/NotificationsContext';
import { readNotification } from "../../lib/apiNotifications";
import { submitFollowResponse } from "../../lib/apiFollow";
import { submitGroupRequestOrInviteResponse } from "../../lib/apiGroups";
import toast from "react-toastify";
import { useRouter } from "next/navigation";

export default function NotificationCard({ notification }) {
    const { refreshNotifications } = useNotifications();
    const router = useRouter();

    const acceptRejectActions = ['follow_request', 'group_invite', 'group_request'];

    const handleNotificationRespond = async (status) => {
        try {

            if (notification.target_detailed_type === 'follow_request') {
                await submitFollowResponse({ follower_uuid: notification.target_uuid, status });
            } else if (notification.target_detailed_type === 'group_invite') {
                await submitGroupRequestOrInviteResponse({ follower_uuid: notification.to_user_uuid, group_uuid: notification.target_uuid, status });
            } else if (notification.target_detailed_type === 'group_request') {
                await submitGroupRequestOrInviteResponse({ follower_uuid: notification.from_user.uuid, group_uuid: notification.target_uuid, status });
            }

           handleViewClick();
        } catch (err) {
            toast.error(err.message || "Failed to respond notification");
        }
    };

    const handleViewClick = async (href) => {
        if (!notification.is_read) {
            await readNotification({ id: notification.id, is_read: 1 });
        }
        await refreshNotifications();
        if (!href) {
            return;
        } else {
            router.push(href);
        }
    };

    return (
        <div className={`notification-card${notification.is_read ? '' : ' unread'}`}>
            <div className="notification-content">
                {notification.from_user && (
                    <div className="notification-user">
                        {notification.from_user.profile_image ? (
                            <Image
                                src={`/frontend-api/image/${notification.from_user.profile_image}`}
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
                    <div className="notification-message" dangerouslySetInnerHTML={{ __html: DOMPurify.sanitize(notification.message) }}></div>
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
                                <a
                                    href="#"
                                    className="link-btn"
                                    onClick={e => {
                                        e.preventDefault();
                                        const href =
                                            notification.target_detailed_type === 'follow_request'
                                                ? `/profile/${notification.target_uuid}`
                                                : notification.target_detailed_type === 'group_invite' || notification.target_detailed_type === 'group_request'
                                                    ? `/groups/${notification.target_uuid}`
                                                    : "#";
                                        handleViewClick(href);
                                    }}
                                >
                                    View
                                </a>
                            )
                            
                        )}
                        {(['follow_request_responded', 'group_request_responded', 'group_event'].includes(notification.target_detailed_type)) && (
                            <a
                                href="#"
                                className="link-btn"
                                onClick={e => {
                                    e.preventDefault();
                                    const href =
                                        notification.target_detailed_type === 'follow_request_responded'
                                            ? `/profile/${notification.target_uuid}`
                                            : `/groups/${notification.target_uuid}`;
                                    handleViewClick(href);
                                }}
                            >
                                View
                            </a>
                        )}
                        
                    </div>
                </div>
            </div>
        </div>
    );
}