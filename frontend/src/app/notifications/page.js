'use client';

import { useRouter } from "next/navigation";
import { useEffect, useState, useMemo } from "react";
import SidebarSection from '../../components/SidebarSection';
import UsersList from '../../components/UsersList';
import NotificationFilterList from '../../components/notifications/NotificationFilterList';
import NotificationList from '../../components/notifications/NotificationList';
import '../../styles/globals.css';
import './notification.css';
import '../../styles/notifications/FilterList.css'

import { mockNotifications } from '../../data/mockData';
import { fetchNotifications } from "../../lib/apiNotifications";

const notificationFilters = [
    { key: 'all', label: 'All' },
    { key: 'follow_request', label: 'Follow Requests' },
    { key: 'group_invitation', label: 'Group Invitations' },
    { key: 'group_join_request', label: 'Join Requests' },
    { key: 'group_event', label: 'Group Events' },
    { key: 'unread', label: 'Unread' }
];

export default function NotificationPage() {
    const [notifications, setNotifications] = useState([]);
    const [loading, setLoading] = useState(true);
    const [selectedFilter, setSelectedFilter] = useState('all');
    useEffect(() => {
        async function fetchData() {
            try {
                const notificationData = await fetchNotifications();
                setNotifications(notificationData.data);
            } catch (err) {
                setError(err.message);
            } finally {
                setLoading(false);
            }
        }
        fetchData();
    }, []);

    const filteredNotifications = useMemo(() => {
        if (!notifications) return [];
        if (selectedFilter === 'all') return notifications;
        if (selectedFilter === 'unread') return notifications.filter(n => !n.isRead);
        return notifications.filter(n => n.type === selectedFilter);
    }, [notifications, selectedFilter]);

    return (
        <main>
        <div className='notification-page-layout'>
            {/* Left Sidebar */}
            <aside className='sidebar left-sidebar'>
            <SidebarSection title='Notifications'>
                <NotificationFilterList
                filters={notificationFilters}
                selectedFilter={selectedFilter}
                onSelect={setSelectedFilter}
                />
            </SidebarSection>
            </aside>
            {/* Main Notification Feed */}
            <section className='main-feed notification-section'>
            <h2>Notifications</h2>
            <NotificationList notifications={filteredNotifications} />
            </section>
            {/* Right Sidebar */}
            <aside className="sidebar right-sidebar">
            <SidebarSection title="All Users">
                <UsersList />
            </SidebarSection>
            </aside>
        </div>
        </main>
    );
}
