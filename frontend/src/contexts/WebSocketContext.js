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
    });

    // Fetch session ID from server on mount and store in localStorage
    useEffect(() => {
        async function syncSession() {
            try {
                const res = await fetch('/frontend-api/session');
                const { sessionId } = await res.json();
                if (sessionId) {
                    localStorage.setItem('session-id', sessionId);
                    console.log(sessionId)
                    setSessionId(sessionId);

                }
            } catch (error) {
                console.error('Session sync failed:', error);
                // Fallback to localStorage if endpoint fails
                const storedSession = localStorage.getItem('session-id');
                if (storedSession) setSessionId(storedSession);
            }
        }
        syncSession();
    }, []);

    
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

    // Handle manual connection (e.g., after login)
    const connect = useCallback((newSessionId) => {
        localStorage.setItem('session-id', newSessionId);
        setSessionId(newSessionId);
    }, []);


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
        sessionId ? `ws://localhost:8080/api/ws?session=${sessionId}` : null,
        onMessage,
        {
            initialDelay: 1000,
            maxDelay: 30000,
            backoffFactor: 2,
            maxAttempts: null,
        }
    );


    return (
        <WebSocketContext.Provider value={{ isConnected, userList, messages, sendAction, connect }}>
            {children}
        </WebSocketContext.Provider>
    );
}

export const useWebsocketContext = () => useContext(WebSocketContext);

