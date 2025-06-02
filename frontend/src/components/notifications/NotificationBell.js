import { useState, useEffect, useRef } from 'react';
import NotificationList from './NotificationList';
import '../../app/notifications/notification.css';

export default function NotificationBell({ notifications }) {
  const [open, setOpen] = useState(false);
  const wrapperRef = useRef(null); // for detecting outside click

  const unreadNotifications = (notifications || []).filter(n => !n.is_read);
  const unreadCount = unreadNotifications.length;

  // Close dropdown when clicking outside
  useEffect(() => {
    const handleClickOutside = (event) => {
      if (wrapperRef.current && !wrapperRef.current.contains(event.target)) {
        setOpen(false);
      }
    };

    document.addEventListener('mousedown', handleClickOutside);
    return () => {
      document.removeEventListener('mousedown', handleClickOutside);
    };
  }, []);

  return (
    <div className="notif-wrapper" ref={wrapperRef}>
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
