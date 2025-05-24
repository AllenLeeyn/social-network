'use client';

import React from 'react';
import SidebarSection from '../../components/SidebarSection';
import UsersList from '../../components/UsersList';
import MessagesChatbox from '../../components/messages-Chatbox';
import DirectMessagesList from '../../components/messages-DirectMessageList';
import GroupChatsList from '../../components/messages-GroupChatsList';
import '../../styles/globals.css';

export default function MessagePage() {
    return (
        <main>
            <div className='message-page-layout'>
                {/* Left Sidebar */}
                <aside className='sidebar left-sidebar'>
                    <SidebarSection title='Direct'>
                        <DirectMessagesList />
                    </SidebarSection>
                    <SidebarSection title='Group'>
                        <GroupChatsList />
                    </SidebarSection>
                </aside>
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
