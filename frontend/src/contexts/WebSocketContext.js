'use client';

import { createContext, useContext, useCallback, useState, useEffect, useRef } from "react";
import { useWebsocket } from "../hooks/useWebsocket";

const WebSocketContext = createContext();

export function WebSocketProvider( { children } ) {
    const [userList, setUserList] = useState([]);
    const [currentChatId, setCurrentChatId] = useState(null); // consider UUID for chatID
    const [messages, setMessages] = useState([]);
    const [isTyping, setIsTyping] = useState(false);
    const [sessionId, setSessionId] = useState(() => {
        if (typeof window !== 'undefined') {
            return localStorage.getItem('session-id') || null;
        }
        return null;
    });
        
    const [isLoadingMore, setIsLoadingMore] = useState(false);
    const [hasMore, setHasMore] = useState(true);


    const userUuid = typeof window !== 'undefined' ? localStorage.getItem('user-uuid') : null;

    // use a ref for currentChatId to avoid unnecessary re-renders ---
    const currentChatIdRef = useRef(currentChatId);
    useEffect(() => {
        currentChatIdRef.current = currentChatId;
    }, [currentChatId]);

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
        console.log("WebSocket received:", data);
        switch (data.action) {
            case 'userList':
                setUserList(
                    data.allClients.map((name, index) => ({
                        name,
                        id: data.clientIDs[index],
                        online: data.onlineClients.includes(data.clientIDs[index]),
                        unread: data.unreadMsgClients?.includes(data.clientIDs[index]) || false,
                    }))
                );
                break;
            case 'online':
                setUserList(prev => {
                        if (prev.some(user => user.id === data.id)) {
                            // Update online status for existing user
                            return prev.map(user =>
                                user.id === data.id ? { ...user, online: true } : user
                            );
                        }
                        // Add new user to the list
                        return [...prev, { id: data.id, name: data.name || 'New User', online: true, unread: false }];
                    });
                break;
            case 'offline':
                setUserList(prev =>
                    prev.map(user =>
                        user.id === data.id ? { ...user, online: false } : user
                    )
                );
                break;
            case 'messageSendOK':
                // Optionally handle message send confirmation
                break;
            case 'messageHistory':
                if (isLoadingMore) {
                        setMessages(prev => [...data.content, ...prev]);
                        setHasMore(data.content.length === 10); // If less than 10, no more to load
                    } else {
                        setMessages(data.content);
                        setHasMore(data.content.length === 10);
                    }
                    setIsLoadingMore(false);
                    break;
            case 'message':
                setMessages(prev => [...prev, data]);
                sendAction({ action: 'userListReq' });
                if (data.senderUUID !== currentChatIdRef.current) {
                    setUserList(prev =>
                        prev.map(user =>
                            user.id === data.senderUUID ? { ...user, unread: true } : user
                        )
                    );
                }
                break;
            case 'messageAck':
                setUserList(prev =>
                    prev.map(user =>
                        user.id === data.senderUUID ? { ...user, unread: false } : user
                    )
                );
                break;
            case 'typing':
                if (data.senderUUID === currentChatIdRef.current) {
                    setIsTyping(true);
                    setTimeout(() => setIsTyping(false), 3000);
                }
                break;
            default:
                console.log('Unknown WebSocket action received:', data.action, data);
                break;
        }
    }, []);

    // Declare useWebsocket AFTER onMessage is defined, BEFORE any useEffect that uses sendAction
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

    
    // Reset typing state when changing chat
    useEffect(() => {
        setIsTyping(false);
    }, [currentChatId]);

    useEffect(() => {
        if (currentChatId && userUuid) {
            // console.log("Sending messageReq for chat history:", currentChatId, userUuid);
            sendAction({
                action: "messageReq",
                receiverUUID: currentChatId,
                content: "-1"
            });
            sendAction({
                action: 'messageAck',
                receiverUUID: currentChatId,
                senderUUID: userUuid
            });
        }
    }, [currentChatId, userUuid, sendAction]);




    return (
        <WebSocketContext.Provider value={{ 
            isConnected, 
            userList,
            setUserList,
            messages, 
            sendAction, 
            connect,
            isTyping,
            currentChatId,
            setCurrentChatId,
            isLoadingMore,
            setIsLoadingMore,
            hasMore,
        }}>
            {children}
        </WebSocketContext.Provider>
    );
}

export const useWebsocketContext = () => useContext(WebSocketContext);

