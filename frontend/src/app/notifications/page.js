'use client';

import React from 'react';
import SidebarSection from '../../components/SidebarSection';
import UsersList from '../../components/UsersList';
import '../../styles/globals.css';

import './notification.css';

export default function NotificationPage() {
    return (
        <main>
            <div className='notification-page-layout'>
                {/* Left Sidebar */}
                <aside className='sidebar left-sidebar'>
                    <SidebarSection title='Notifications' />
                    {/* No inner content yet */}
                </aside>
                {/* Main Notification Feed */}
                <section className='main-feed notification-section'>
                    {/* Placeholder for Notification List */}
                    <div className="notification-card">
                        <h2>Notifications</h2>
                        <p>Your notifications will appear here.</p>
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
