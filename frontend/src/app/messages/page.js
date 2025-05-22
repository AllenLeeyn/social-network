"use client"; 

import React, { useEffect, useState } from 'react';
import SidebarSection from '../../components/SidebarSection';
import UsersList from '../../components/UsersList';
import { useWebsocketContext } from '../../contexts/WebSocketContext';
import './messages.css';


import {
    sampleUsers,
    sampleConversations,
    sampleMessages,
} from '../../data/mockData';



export default function MessagePage() {
    const { messages: contextMessages, sendAction, userList } = useWebsocketContext();
    const [activeConversation, setActiveConversation] = useState(null);
    const [inputMessage, setInputMessage] = useState('');

    // const messages = sampleMessages;
    const messages = contextMessages;
    
    const handleSendMessage = (e) => {
        e.preventDefault();
        if (!inputMessage.trim() || !activeConversation) return;
    
        
        const newMessage = {
            action: 'message',
            content: inputMessage,
            receiverID: Number(activeConversation.id)
        };
        
        sendAction(newMessage);
        setInputMessage('');
    };

    const handleUserSelect = (userId, userName) => {
        setActiveConversation({ id: userId, name: userName });
    };



    return (
        <main>
            <div className='message-page-layout'>
                {/* Left Sidebar */}
                <aside className='sidebar left-sidebar'>
                    <SidebarSection title='Individuals'>
                        <ul className='individuals'>
                            {sampleConversations.filter(c => c.type === 'individual').map(convo => (
                                <li
                                    key={convo.id}
                                    className={`conversation-item${activeConversation === convo.id ? 'active' : ''}${convo.unread ? 'unread' : ''}`}
                                    onClick={() => setActiveConversation(convo.id)}
                                >
                                    {convo.name}
                                    {convo.unread > 0 && <span className='unread-count'>{convo.unread}</span>}
                                </li>
                            ))}
                        </ul>
                    </SidebarSection>
                    <SidebarSection title="Group">
                        <ul className='groups'>
                            {sampleConversations.filter( c => c.type === 'group').map(convo => (
                                <li
                                key={convo.id}
                                className={`conversation-item ${activeConversation === convo.id ? 'active' : ''} ${convo.unread ? 'unread' : ''}`}
                                onClick={() => setActiveConversation(convo.id)}>
                                    {convo.name}
                                    {convo.unread > 0 && <span className='unread-count'>{convo.unread}</span>}
                                </li>
                            ))}
                        </ul>
                    </SidebarSection>
                </aside>
                 {/* Main Chat Area */}
                <section className='main-feed message-list-section'>
                    {/* Chat Header + Messages + Typing Indicator */}
                    <h2>
                        {activeConversation ? activeConversation.name : 'Latest Messages'}
                    </h2>
                    <div className="message-view">
                        <h2>
                        {activeConversation?sampleConversations.find(c => c.id === activeConversation)?.name : 'Latest Messages'}
                        </h2>
                        <div className='message-list'>
                            {messages
                                .filter(m => activeConversation && m.conversationId === activeConversation)
                                .map(message => {
                                    const sender = sampleUsers.find(u => u.id === message.senderId);
                                    return (
                                        <div key={message.id} className='message-item'>
                                            <img src={sender?.avatar} alt={sender?.username} className='avatar' />
                                            <div className='message-content'>
                                                <div className='message-header'>
                                                    <span className='sender-name'>{sender?.username ?? 'You'}</span>
                                                    <span className='timestamp'>{new Date(message.timestamp).toLocaleTimeString()}</span>
                                                </div>
                                                <p>{message.content}</p>
                                            </div>
                                        </div>
                                    );
                                })}
                        </div>
                    </div>
                    
                    <form onSubmit={handleSendMessage} className='message-input-form'>
                        <input
                            type='text'
                            value={inputMessage}
                            onChange={(e) => setInputMessage(e.target.value)}
                            placeholder='Type a message...'
                            disabled={!activeConversation}
                            className='message-input'
                        />
                    </form>
                </section>
                {/* Right Sidebar */}
                <aside className="sidebar right-sidebar">
                    <SidebarSection title="Active Users">
                        <UsersList 
                            onUserSelect={handleUserSelect}  
                            activeConversation={activeConversation?.id}
                        />
                    </SidebarSection>
                </aside>
            </div>
        </main>
    );
}

