"use client";

import { useWebsocketContext } from '../contexts/WebSocketContext';
import { useActiveChat } from '../contexts/ActiveChatContext';
import { useState, useMemo, useEffect } from 'react';

export default function MessagesChatbox() {
    // const { messages, sendAction, userList, isTyping } = useWebsocketContext();
    const { userList, messages, sendAction, isConnected, isTyping, currentChatId, setCurrentChatId } = useWebsocketContext();
    const { activeChat } = useActiveChat();
    const [inputMessage, setInputMessage] = useState('');

     // Get your own UUID from context or localStorage
    const userUuid = typeof window !== 'undefined' ? localStorage.getItem('user-uuid') : null;

    // Keep currentChatId in sync with activeChat
    useEffect(() => {
        if (activeChat && activeChat.id !== currentChatId) {
            setCurrentChatId(activeChat.id);
        }
    }, [activeChat, currentChatId, setCurrentChatId]);

    // Filter messages for the active chat
    const filteredMessages = useMemo(() => {
        if (!activeChat) return [];
        return messages.filter(
            m => (m.senderUUID === activeChat.id && m.receiverUUID === userUuid) || 
            (m.senderUUID === userUuid && m.receiverUUID === activeChat.id)
        );
    }, [messages, activeChat, userUuid]);

    // Handle typing and sending
    const handleInputChange = (e) => {
        setInputMessage(e.target.value);
        if (activeChat && userUuid) {
            sendAction({ 
                action: 'typing', 
                senderUUID: userUuid,
                receiverUUID: activeChat.id 
            });
        }
    };

    const handleSendMessage = (e) => {
        e.preventDefault();
        if (!inputMessage.trim() || !activeChat || !userUuid || !activeChat.id) return;
        console.log("Sending message", {
        senderUUID: userUuid,
        receiverUUID: activeChat?.id,
        content: inputMessage,
        });
        sendAction({
            action: 'message',
            senderUUID: userUuid,
            receiverUUID: activeChat.id,
            content: inputMessage,
            createdAt: new Date().toISOString() 
        });
        setInputMessage('');
    };

    return (
        <div className='chat-component'>
            <h2>
                {activeChat ? activeChat.name : 'Select a user to chat'}
                {!isConnected && <span style={{color: 'red', marginLeft: '1em'}}>Disconnected</span>}
            </h2>
            <div className='messages-list'>
                {filteredMessages.length === 0 ? (
                    <div className="no-messages">
                    {activeChat
                        ? "No messages yet.  Start the conversation!"
                        : "Select a user to view messages."}
                    </div>
                ) : (
                    filteredMessages.map((msg, index) => {
                    const sender = userList.find(u => u.id === msg.senderUUID) || { name: 'You' };
                        return (
                            <div key={index} className='message-item'>
                            <span><strong>{sender.name}</strong>: {msg.content}</span>
                            <span className='timestamp'>{new Date(msg.createdAt).toLocaleTimeString()}</span>
                            </div>
                        );
                    })
                )}
            </div>
            {isTyping && activeChat && (
                <div className="typing-indicator">{activeChat.name} is typing...</div>
            )}
            <form onSubmit={handleSendMessage} className='message-input-form'>
                <input
                    type='text'
                    value={inputMessage}
                    onChange={handleInputChange}
                    placeholder='Type a message...'
                    disabled={!activeChat}
                    className='message-input'
                />
            </form>
        </div>
    );
}
