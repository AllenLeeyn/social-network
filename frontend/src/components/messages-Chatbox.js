"use client";

import { useWebsocketContext } from '../contexts/WebSocketContext';
import { useActiveChat } from '../contexts/ActiveChatContext';
import { useState, useMemo, useEffect, useRef, useCallback } from 'react';


export default function MessagesChatbox() {

    const { 
        userList, 
        messages, 
        sendAction, 
        isConnected, 
        isTyping, 
        currentChatId, 
        setCurrentChatId,
        isLoadingMore, 
        setIsLoadingMore,
        hasMore,
    } = useWebsocketContext();
    
    const { activeChat } = useActiveChat();
    const [inputMessage, setInputMessage] = useState('');
    const messagesEndRef = useRef(null);
    const messagesContainerRef = useRef(null);
    const prevScrollHeight = useRef(0);
    const scrollThrottleRef = useRef(false);


     // Get your own UUID from context or localStorage
    const userUuid = typeof window !== 'undefined' ? localStorage.getItem('user-uuid') : null;

    // Keep currentChatId in sync with activeChat
    useEffect(() => {
        if (activeChat && activeChat.id !== currentChatId) {
            setCurrentChatId(activeChat.id);
        }
    }, [activeChat, currentChatId]);

    // Filter messages for the active chat
    const filteredMessages = useMemo(() => {
        if (!activeChat) return [];
        return messages.filter(m => 
            (m.senderUUID === activeChat.id && m.receiverUUID === userUuid) || 
            (m.senderUUID === userUuid && m.receiverUUID === activeChat.id)
        );
    }, [messages, activeChat, userUuid]);

    // Auto-scroll to bottom on new messages
    const prevMessagesLength = useRef(filteredMessages.length);

    useEffect(() => {
        // Only scroll if a new message was added at the end (not when loading older messages)
        if (
            !isLoadingMore &&
            filteredMessages.length > prevMessagesLength.current
        ) {
            messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
        }
        prevMessagesLength.current = filteredMessages.length;
    }, [filteredMessages, isLoadingMore]);

    const handleScroll = useCallback(() => {
        if (scrollThrottleRef.current || !messagesContainerRef.current || isLoadingMore || !hasMore) return;

    const { scrollTop } = messagesContainerRef.current;
        if (scrollTop < 100 && filteredMessages.length > 0) {
            scrollThrottleRef.current = true;
            setTimeout(() => {
                scrollThrottleRef.current = false;
            }, 300); // 300ms throttle

            setIsLoadingMore(true);
            prevScrollHeight.current = messagesContainerRef.current.scrollHeight;
            loadPreviousMessages();
        }
    }, [isLoadingMore, hasMore, filteredMessages, activeChat?.id]);

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


    const loadPreviousMessages = () => {
    if (!filteredMessages.length) {
        setIsLoadingMore(false);
        return;
    }
    setIsLoadingMore(true);
    console.log("Loading more, isLoadingMore:", isLoadingMore);
    const oldestMsg = filteredMessages[0]; // assuming messages are oldest-to-newest
        sendAction({
            action: "messageReq",
            receiverUUID: activeChat.id,
            content: oldestMsg.ID.toString(), // send timestamp as cursor
        });
    }; 



    // console.log("messages:", messages);
    // console.log("activeChat.id:", activeChat?.id, "userUuid:", userUuid);
    // console.log("filteredMessages:", filteredMessages);

    return (
        <div className='chat-component'>
            <h2>
                {activeChat ? activeChat.name : 'Select a user to chat'}
                {!isConnected && <span style={{color: 'red', marginLeft: '1em'}}>Disconnected</span>}
            </h2>
            <div className='messages-list' ref={messagesContainerRef}>
                {filteredMessages.length === 0 ? (
                    <div className="no-messages">
                    {activeChat
                        ? "No messages yet.  Start the conversation!"
                        : "Select a user to view messages."}
                    </div>
                ) : (
                    filteredMessages.map((msg, index) => {
                        const isSent = msg.senderUUID === userUuid;
                        const sender = userList.find(u => u.id === msg.senderUUID) ;
                        return (
                            <div key={index} className={`message-item ${isSent ? 'sent' : 'received'}`}>
                                <div className="message-bubble">
                                    <strong>{isSent ? "You" : (sender ? sender.name : msg.senderUUID)}</strong>
                                    <div className="message-content">{msg.content}</div>
                                    <span className='timestamp'>{new Date(msg.createdAt).toLocaleTimeString()}</span>
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
            </form>
        </div>
    );
}
