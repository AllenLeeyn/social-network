'use client';

import React, { useState, useMemo } from 'react';
import SidebarSection from '../../components/SidebarSection';
import UsersList from '../../components/UsersList';
import '../../styles/globals.css';
import './notification.css';


import { mockNotifications } from '../../data/mockData'

const notificationFilters = [
    { key: 'all', label: 'All' },
    { key: 'follow_request', label: 'Follow Requests' },
    { key: 'group_invitation', label: 'Group Invitations' },
    { key: 'group_join_request', label: 'Join Requests' },
    { key: 'group_event', label: 'Group Events' },
    { key: 'unread', label: 'Unread' }
];


export default function NotificationPage() {

    const [selectedFilter, setSelectedFilter] = useState('all');

  // Filtering logic
    const filteredNotifications = useMemo(() => {
        if (selectedFilter === 'all') return mockNotifications;
        if (selectedFilter === 'unread') return mockNotifications.filter(n => !n.isRead);
        return mockNotifications.filter(n => n.type === selectedFilter);
    }, [selectedFilter]);



    return (
        <main>
            <div className='notification-page-layout'>
                {/* Left Sidebar */}
                <aside className='sidebar left-sidebar'>
                    <SidebarSection title='Notifications'>
                                <ul className="notification-filter-list">
                                {notificationFilters.map(filter => (
                                    <li
                                    key={filter.key}
                                    className={selectedFilter === filter.key ? 'active' : ''}
                                    onClick={() => setSelectedFilter(filter.key)}
                                    style={{ cursor: 'pointer', padding: '0.5em 0' }}
                                    >
                                    {filter.label}
                                    </li>
                                ))}
                                </ul>
                    </SidebarSection>
                </aside>
                {/* Main Notification Feed */}
                <section className='main-feed notification-section'>
                    <h2>Notifications</h2>
                    <div className="notification-list">
                        {mockNotifications.length === 0 ? (
                            <p>Your notifications will appear here.</p>
                        ) : (
                            mockNotifications.map(notification => (
                                <div
                                    key={notification.id}
                                    className={`notification-card${notification.isRead ? '' : ' unread'}`}
                                >
                                    <div className="notification-content">
                                        {notification.fromUser && (
                                            <img
                                                src={notification.fromUser.avatar}
                                                alt={notification.fromUser.name}
                                                className="notification-avatar"
                                            />
                                        )}
                                        <div>
                                            <div className="notification-message">
                                                {notification.message}
                                            </div>
                                            <div className="notification-timestamp">
                                                {new Date(notification.timestamp).toLocaleString()}
                                            </div>
                                            <div className="notification-actions">
                                                {notification.actions.includes('accept') && (
                                                    <button>Accept</button>
                                                )}
                                                {notification.actions.includes('decline') && (
                                                    <button>Decline</button>
                                                )}
                                                {notification.actions.includes('view') && (
                                                    <button>View</button>
                                                )}
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            ))
                        )}
                    </div>
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
