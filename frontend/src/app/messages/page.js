'use client';

import React from 'react';
import SidebarSection from '../../components/SidebarSection';
import UsersList from '../../components/UsersList';
import MessagesChatbox from '../../components/messages/Chatbox';
import GroupChatsList from '../../components/messages/GroupChatList';
import '../../styles/globals.css'
import './messages.css';

export default function MessagePage() {
    return (
        <main>
            <div className='message-page-layout'>
                {/* Main Chat Area */}
                <section className='main-feed message-list-section'>
                    <MessagesChatbox />
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
