'use client';

import { createContext, useContext, useCallback, useState, useEffect } from "react";
import { useWebsocket } from "../hooks/useWebsocket";


const WebSocketContext = createContext();

export function WebSocketProvider( { children } ) {
    const [userList, setUserList] = useState([]);
    const [currentChatId, setCurrentChatId] = useState(null);
    const [messages, setMessages] = useState([]);
    const [isTyping, setIsTyping] = useState(false);
    const [sessionId, setSessionId] = useState(() => {
        if (typeof window !== 'undefined') {
            const id = localStorage.getItem('session-id');
            return id || null;
        }
    })

    // memoize the onMessage handler with useCallback
    // if this doesnt change, it will pass and wont re-render; WS connection will always render when new data gets passed
    const onMessage = useCallback((data) => {
        if (data.action === 'start') {
            setUserList(
                data.allClients.map((name,index) => ({
                    name,
                    id: data.clientIDs[index],
                    online: data.onlineClients.includes(data.clientIDs[index]),
                    unread: data.unreadMsgClients?.includes(data.clientIDs[index]) || false,
                }))
            );
        } else if (data.action === 'online') {
            setUserList(prev =>
                prev.map(user =>
                    user.id === data.id ? { ...user, online: true} : user
                )
            );
        } else if (data.action === 'offline') {
            setUserList(prev => 
                prev.map(user =>
                    user.id === data.id ? { ...user, online:false } : user
                )
            );
        } else if (data.action === 'messageSendOK') {
            // 
        } else if (data.action === 'messageHistory') {
            setMessages(data.content);
        } else if (data.action === 'message') {
            setMessages(prev => [...prev, data]);
            if (data.senderId !== currentChatId) {
                setUserList(prev => 
                    prev.map(user =>
                        user.id === data.senderId ? { ...user, unread: true } : user
                    )
                );
            }
        } else if (data.action === 'typing') {
            if (data.senderId === currentChatId) {
                setIsTyping(true);
                setTimeout(() => setIsTyping(false), 800);
            }
        }
    }, [currentChatId]);

    const { isConnected, sendAction } = useWebsocket(
        sessionId ? `ws://localhost:8080/ws?session=${sessionId}` : null,
        onMessage,
        {
            initialDelay: 1000,
            maxDelay: 30000,
            backoffFactor: 2,
            maxAttempts: null,
        }
    );

// insert  copy

    // Sync session ID with localStorage
    useEffect(() => {
        const handleStorage = () => {
        const newSession = localStorage.getItem('session-id');
        if (newSession !== sessionId) {
            setSessionId(newSession);
        }
        };
        window.addEventListener('storage', handleStorage);
        return () => window.removeEventListener('storage', handleStorage);
    }, [sessionId]);

    // conn handling
    const connect = useCallback((newSessionId) => {
        localStorage.setItem('session-id', newSessionId);
        setSessionId(newSessionId);
    }, []);

    return (
        <WebSocketContext.Provider value={{ isConnected, userList, messages, sendAction, connect }}>
            {children}
        </WebSocketContext.Provider>
    );
}

export const useWebsocketContext = () => useContext(WebSocketContext);



    // // using hook for connection logic
    // const { isConnected, sendAction } = useWebsocket(
    //     sessionId ? `ws://localhost:8080/ws?session=${sessionId}` : null,
    //     (data) => {
    //             // "start": Initialize user list
    //             if (data.action === "start") {
    //                 setUserList(
    //                     data.allClients.map((name, index) => ({
    //                     name,
    //                     id: data.clientIDs[index],
    //                     online: data.onlineClients.includes(data.clientIDs[index]),
    //                     unread: data.unreadMsgClients?.includes(data.clientIDs[index]) || false,
    //                     }))
    //                 );
    //             }
    //             // "online": Mark user as online
    //             else if (data.action === "online") {
    //                 setUserList(prev =>
    //                     prev.map(user =>
    //                     user.id === data.id ? { ...user, online: true } : user
    //                     )
    //                 );
    //             }
    //             // "offline": Mark user as offline
    //             else if (data.action === "offline") {
    //                 setUserList(prev =>
    //                     prev.map(user =>
    //                     user.id === data.id ? { ...user, online: false } : user
    //                     )
    //                 );
    //             }
    //             // "messageSendOK": Clear input (handled in component)
    //             else if (data.action === "messageSendOK") {
    //                 // You can set a flag or call a callback here if needed
    //                 // Clearing input is usually handled in the component
    //             }
    //             // "messageHistory": Load chat history
    //             else if (data.action === "messageHistory") {
    //                 setMessages(data.content); // Adjust based on your data structure
    //             }
    //             // "message": Add new message
    //             else if (data.action === "message") {
    //                 setMessages(prev => [...prev, data]);
    //                 // Mark message as read/ack if needed
    //                 if (data.senderId !== currentChatId) {
    //                     // Mark as unread in user list
    //                     setUserList(prev =>
    //                     prev.map(user =>
    //                         user.id === data.senderId ? { ...user, unread: true } : user
    //                     )
    //                     );
    //             }
    //             }
    //             // "typing": Show typing indicator
    //             else if (data.action === "typing") {
    //                 if (data.senderId === currentChatId) {
    //                     setIsTyping(true);
    //                     setTimeout(() => setIsTyping(false), 800); // Hide after delay
    //                 }
    //             }
    //         // Add more msg handler 
    //     },
    //     {
    //         initialDelay: 1000,
    //         maxDelay: 30000,
    //         backoffFactor: 2,
    //         maxAttempts: null
    //     }
    // );
    // // console.log('sending Session ID', sessionId)

