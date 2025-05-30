import { useState } from 'react';
import NotificationList from './NotificationList';
import '../../app/notifications/notification.css';

export default function NotificationBell({ notifications }) {
  const [open, setOpen] = useState(false);

const unreadNotifications = (notifications || []).filter(n => !n.is_read);
const unreadCount = unreadNotifications.length;

  return (
    <div className="notif-wrapper">
      <a className="bell-btn" onClick={() => setOpen(!open)}>
        🔔
        {unreadCount > 0 && (
          <span className="notif-count">{unreadCount}</span>
        )}
      </a>

      {open && (
        <div className="notif-dropdown">
          <NotificationList notifications={unreadNotifications} />
            <a href="/notifications" className="view-all-btn" onClick={() => setOpen(false)}>
            Show All Notifications
            </a>
        </div>
      )}
    </div>
  );
}
