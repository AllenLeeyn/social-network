"use client";

import { useWebsocketContext } from '../contexts/WebSocketContext';
import { useActiveChat } from '../contexts/ActiveChatContext';
import { useState, useMemo, useEffect } from 'react';

export default function MessagesChatbox() {
    // const { messages, sendAction, userList, isTyping } = useWebsocketContext();
    const { userList, messages, sendAction, isConnected, isTyping, currentChatId, setCurrentChatId } = useWebsocketContext();
    const { activeChat } = useActiveChat();
    const [inputMessage, setInputMessage] = useState('');


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
            m => m.senderId === activeChat.id || m.receiverId === activeChat.id
        );
    }, [messages, activeChat]);

    // Handle typing and sending
    const handleInputChange = (e) => {
        setInputMessage(e.target.value);
        if (activeChat) {
            sendAction({ action: 'typing', receiverID: activeChat.id });
        }
    };

    const handleSendMessage = (e) => {
        e.preventDefault();
        if (!inputMessage.trim() || !activeChat) return;
        sendAction({
            action: 'message',
            content: inputMessage,
            receiverID: Number(activeChat.id)
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
                    const sender = userList.find(u => u.id === msg.senderId) || { name: 'You' };
                        return (
                            <div key={index} className='message-item'>
                            <span><strong>{sender.name}</strong>: {msg.content}</span>
                            <span className='timestamp'>{new Date(msg.timestamp).toLocaleTimeString()}</span>
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



{/*             <div className='messages-list'>
                {filteredMessages.map((msg, index) => {
                    const sender = userList.find(u => u.id === msg.senderId) || { name: 'You' };
                    return (
                        <div key={index} className='message-item'>
                            <span><strong>{sender.name}</strong>: {msg.content}</span>
                            <span className='timestamp'>{new Date(msg.timestamp).toLocaleTimeString()}</span>
                        </div>
                    );
                })}
            </div> */}