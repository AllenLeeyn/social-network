'use client';
import { createContext, useContext, useEffect, useState } from 'react';
import { fetchNotifications } from '../lib/apiNotifications';

const NotificationsContext = createContext();

export function NotificationsProvider({ children }) {
  const [notifications, setNotifications] = useState([]);
  const [loadingNotification, setLoadingNotification] = useState(true);

  const fetchData = async () => {
    try {
      const data = await fetchNotifications();
      setNotifications(data.data);
    } catch (err) {
      console.error('Error fetching notifications:', err);
    } finally {
      setLoadingNotification(false);
    }
  };
  // update the notifications every 5 minutes
  useEffect(() => {
    fetchData();

    const interval = setInterval(fetchData, 300000); // 5 minutes
    return () => clearInterval(interval);
  }, []);

  return (
    <NotificationsContext.Provider
      value={{
        notifications,
        loadingNotification,
        refreshNotifications: fetchData, 
        setNotifications, // optional, for local updates
      }}
    >
      {children}
    </NotificationsContext.Provider>
  );
}

export const useNotifications = () => useContext(NotificationsContext);
