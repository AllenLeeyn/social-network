"use client"; 

import React, { useEffect, useState } from 'react';
import SidebarSection from '../../components/SidebarSection';
import './messages.css';


import {
    sampleUsers,
    sampleConversations,
    sampleMessages,
} from '../../data/mockData';


export default function MessagePage() {
    
    const [activeConversation, setActiveConversation] = useState(null);
    const [messages, setMessages] = useState(sampleMessages);
    const [inputMessage, setInputMessage] = useState('');

    useEffect(() => {
        const ws = new WebSocket('ws://localhost:8080/ws');

        ws.onopen = () => {
        };

        ws.onmessage = (event) => {
            // Handle incoming messages from the backend
            const newMessage = JSON.parse(event.data);
            setMessages(prev => [...prev, newMessage]);
        };
        return () => ws.close();
    }, []);
    
    const handleSendMessage = (e) => {
        e.preventDefault();
        if (!inputMessage.trim() || !activeConversation) return;
    
        // Temporary mock send until backend is connected
        const newMessage = {
            id: Date.now().toString(),
            content: inputMessage,
            senderId: 'current-user',
            timestamp: new Date().toISOString(),
            conversationId: activeConversation,
        };
        
        setMessages(prev => [...prev, newMessage]);
        setInputMessage('');
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
                        <ul className="users">
                            {sampleUsers.map(user => (
                                <li key={user.id} className={`user-item${user.online ? " online" : ""}${user.unread ? " unread" : ""}`}>
                                    <img src={user.avatar} alt={user.username} />
                                    <span>{user.fullName} ({user.username})</span>
                                </li>
                            ))}
                        </ul>
                    </SidebarSection>
                </aside>
            </div>
        </main>
    );
}
