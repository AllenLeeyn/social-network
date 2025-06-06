"use client";

import { useWebsocketContext } from '../../contexts/WebSocketContext';
import { useActiveChat } from '../../contexts/ActiveChatContext';
import { useState, useEffect, useRef, useCallback } from 'react';
import Link from 'next/link';
import EmojiPicker from "emoji-picker-react";

export default function MessagesChatbox() {

    const {
        messages,
        sendAction, 
        isConnected, 
        isTyping, 
        currentChatUUID, 
        setCurrentChatUUID,
        currentGroupUUID, 
        setCurrentGroupUUID,
        isLoadingMore, 
        setIsLoadingMore,
        hasMore,
    } = useWebsocketContext();
    
    const { activeChat } = useActiveChat();
    const [inputMessage, setInputMessage] = useState('');
    const [showPicker, setShowPicker] = useState(false);
    const messagesEndRef = useRef(null);
    const messagesContainerRef = useRef(null);
    const prevScrollHeight = useRef(0);
    const scrollThrottleRef = useRef(false);

    // Get your own UUID from context or localStorage
    const userUUID = typeof window !== 'undefined' ? localStorage.getItem('user-uuid') : null;

    const handleEmojiClick = (emojiData) => {
        setInputMessage((prev) => prev + emojiData.emoji);
    };

    // Keep currentChatId in sync with activeChat
    useEffect(() => {
        if (activeChat && activeChat.uuid !== currentChatUUID &&
            activeChat.groupUUID != currentGroupUUID
        ) {
            setCurrentChatUUID(activeChat.uuid);
            setCurrentGroupUUID(activeChat.groupUUID);
        }
    }, [activeChat, currentChatUUID, currentGroupUUID]);

    const handleScroll = useCallback(() => {
        if (scrollThrottleRef.current || 
            !messagesContainerRef.current || 
            isLoadingMore || 
            !hasMore) return;

        const { scrollTop } = messagesContainerRef.current;
            if (scrollTop < 100 && messages.length > 0) {
                scrollThrottleRef.current = true;
                setTimeout(() => {
                    scrollThrottleRef.current = false;
                }, 300); 

            setIsLoadingMore(true);
            prevScrollHeight.current = messagesContainerRef.current.scrollHeight;
            loadPreviousMessages();
        }
    }, [isLoadingMore, hasMore, messages, activeChat?.uuid]);

    useEffect(() => {
        const container = messagesContainerRef.current;
        if (!container) return;

        const throttledScroll = () => {
            if (isLoadingMore) return;
            handleScroll();
        };

        container.addEventListener('scroll', throttledScroll);
        return () => container.removeEventListener('scroll', throttledScroll);
    }, [handleScroll, isLoadingMore]);

    // Maintain scroll position after prepending
    useEffect(() => {
        if (isLoadingMore || !messagesContainerRef.current) return;
        const newScrollHeight = messagesContainerRef.current.scrollHeight;
        messagesContainerRef.current.scrollTop = newScrollHeight - prevScrollHeight.current;
    }, [messages, isLoadingMore]);

    // Handle typing and sending
    const handleInputChange = (e) => {
        setInputMessage(e.target.value);
        if (activeChat && userUUID) {
            sendAction({ 
                action: 'typing', 
                senderUUID: userUUID,
                receiverUUID: activeChat.receiverUUID 
            });
        }
    };

    const handleSendMessage = (e) => {
        e.preventDefault();
        if (!inputMessage.trim() || !activeChat || !userUUID || !activeChat.uuid) return;
        console.log(activeChat.groupUUID); 
        sendAction({
            action: 'message',
            receiverUUID: activeChat.receiverUUID,
            groupUUID: activeChat.groupUUID,
            content: inputMessage,
        });
        setInputMessage('');
    };

    const loadPreviousMessages = () => {
        if (!messages.length) {
            setIsLoadingMore(false);
            return;
        }
        setIsLoadingMore(true);
        console.log("Loading more, isLoadingMore:", isLoadingMore);
        const oldestMsg = messages[0]; // assuming messages are oldest-to-newest
            sendAction({
                action: "messageReq",
                receiverUUID: activeChat.receiverUUID,
                groupUUID: activeChat.groupUUID,
                content: oldestMsg.ID.toString(), // send timestamp as cursor
            });
    };

    return (
        <div className='chat-component'>
            <h2>
                {activeChat ? (
                    <Link href={`/profile/${activeChat.uuid}`} style={{ fontWeight: 'bold', textDecoration: 'none' }}>
                        {activeChat.name}
                    </Link>
                    ) : 'Select a user to chat'}
                {!isConnected && <span style={{color: 'red', marginLeft: '1em'}}>Disconnected</span>}
            </h2>
            <div className='messages-list' ref={messagesContainerRef}>
                {messages.length === 0 ? (
                    <div className="no-messages">
                    {activeChat
                        ? "No messages yet.  Start the conversation!"
                        : "Select a user to view messages."}
                    </div>
                ) : (
                    messages.map((msg) => {
                        const isSent = msg.senderUUID === userUUID;
                        return (
                            <div key={msg.ID} className={`message-item ${isSent ? 'sent' : 'received'}`}>
                                <div className="message-bubble">
                                    <pre className="message-content">{msg.content}</pre>
                                    <span className='timestamp'>
                                        {isSent ? "You " : (msg.senderName || "unknown ")}
                                         {`[${new Date(msg.createdAt).toLocaleTimeString()}]`}
                                    </span>
                                </div>
                            </div>
                        );
                    })
                )}
                <div ref={messagesEndRef} />
            </div>
            <div className="typing-indicator-area">
                {isTyping && activeChat && (
                    <div className="typing-indicator">{activeChat.name} is typing...</div>
                )}
            </div>

            <form onSubmit={handleSendMessage} className='message-input-form'>
                <input
                    type='text'
                    value={inputMessage}
                    onChange={handleInputChange}
                    placeholder='Type a message...'
                    disabled={!activeChat}
                    className='message-input'
                />

                <button
                    type='button'
                    onClick={() => setShowPicker((prev) => !prev)}
                    title="Insert emoji"
                    >
                    😊
                </button>
            </form>

            {showPicker && (
                <div>
                    <EmojiPicker onEmojiClick={handleEmojiClick} />
                </div>
            )}
        </div>
    );
}
