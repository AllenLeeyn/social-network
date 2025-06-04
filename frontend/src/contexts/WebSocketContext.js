'use client';

import { createContext, useContext, useCallback, useState, useEffect, useRef } from "react";
import { useWebsocket } from "../hooks/useWebsocket";
import { usePathname } from 'next/navigation';
import { useActiveChat } from './ActiveChatContext'; 

const WebSocketContext = createContext();

export function WebSocketProvider( { children } ) {
    const [userList, setUserList] = useState([]);
    const [currentChatUUID, setCurrentChatUUID] = useState(null);
    const [currentGroupUUID, setCurrentGroupUUID] = useState(null);
    const [messages, setMessages] = useState([]);
    const [isTyping, setIsTyping] = useState(false);
    const { activeChat } = useActiveChat();

    const [isLoadingMore, setIsLoadingMore] = useState(false);
    const isLoadingMoreRef = useRef(isLoadingMore);
    const [hasMore, setHasMore] = useState(true);

    useEffect(() => {
        isLoadingMoreRef.current = isLoadingMore;
    }, [isLoadingMore]);

    const userUUID = typeof window !== 'undefined' ? localStorage.getItem('user-uuid') : null;

    // use a ref for currentChatId to avoid unnecessary re-renders ---
    const currentChatUUIDRef = useRef(currentChatUUID);
    useEffect(() => {
        currentChatUUIDRef.current = currentChatUUID;
    }, [currentChatUUID]);

    const pathname = usePathname();
    useEffect(() => {
    if (pathname.startsWith('/messages')) {
        currentChatUUIDRef.current = activeChat?.uuid || null;
    } else {
        currentChatUUIDRef.current = null;
    }
    }, [pathname, activeChat]);

    const onMessage = useCallback((data) => {
        console.log("WebSocket received:", data);
        switch (data.action) {
            case 'userList':
                const followingsName = data.followingsName ?? [];
                const followingsUUID = data.FollowingsUUID ?? [];
                const onlineFollowings = data.onlineFollowings ?? [];
                const unreadMsgFollowings = data.unreadMsgFollowings ?? [];

                const usersAndGroups = [
                ...followingsName.map((name, index) => ({
                    type: "user",
                    name,
                    uuid: followingsUUID[index],
                    groupUUID: '00000000-0000-0000-0000-000000000000',
                    online: onlineFollowings.includes(followingsUUID[index]),
                    unread: unreadMsgFollowings?.includes(followingsUUID[index]) || false,
                    receiverUUID: followingsUUID[index],
                })),
                ...(data.groupList?.map(group => ({
                    type: 'group',
                    name: group.title,
                    uuid: group.uuid,
                    groupUUID: group.uuid,
                    online: true,
                    unread: false,
                    receiverUUID: group.creator_uuid,
                })) || [])
                ];

                const uniqueUserList = Array.from(
                new Map(usersAndGroups.map(item => [item.uuid, item])).values()
                );

                setUserList(uniqueUserList);
                break;

            case 'online':
                setUserList(prev => {
                        if (prev.some(user => user.uuid === data.uuid)) {
                            // Update online status for existing user
                            return prev.map(user =>
                                user.uuid === data.uuid ? { ...user, online: true } : user
                            );
                        }
                        // Add new user to the list
                        return [...prev, { 
                            type: "user",
                            name: data.name || 'New User',
                            uuid: data.uuid, 
                            groupUUID: '00000000-0000-0000-0000-000000000000',
                            online: true,
                            unread: false,
                            receiverUUID: data.uuid }];
                    });
                break;

            case 'offline':
                setUserList(prev =>
                    prev.map(user =>
                        user.uuid === data.uuid ? { ...user, online: false } : user
                    )
                );
                break;

            case 'messageSendOK':
                // Optionally handle message send confirmation
                break;

            case 'messageHistory':
                const newMessages = Array.isArray(data.content) ? [...data.content].reverse() : [];
                console.log('isLoadingMore:', isLoadingMore, 'new batch:', newMessages.map(m => m.ID));

                setMessages(prev => {
                    const existingIDs = new Set(prev.map(m => m.ID));
                    const merged = [...newMessages.filter(m => !existingIDs.has(m.ID)), ...prev];
                    merged.sort((a, b) => a.ID - b.ID);
                    return merged;
                });

                setIsLoadingMore(false);
                setHasMore(newMessages.length === 10); 
                break;

            case 'message':
                if (data.senderUUID === currentChatUUIDRef.current || 
                    data.receiverUUID === currentChatUUIDRef.current ||
                    data.groupUUID === currentChatUUIDRef.current) {
                    setMessages(prev => [...prev, data]);
                }

                if (data.groupUUID !== '00000000-0000-0000-0000-000000000000') return;
                if (data.senderUUID !== currentChatUUIDRef.current) {
                    setUserList(prev =>{
                        const index = prev.findIndex(user => 
                            user.uuid === data.senderUUID && 
                            user.type === "user");
                        if (index === -1) return prev;

                        const bumpedUser = { ...prev[index], unread: true };
                        const updatedList = [bumpedUser, ...prev.slice(0, index), ...prev.slice(index + 1)];
                        return updatedList;
                    });
                }
                break;

            case 'typing':
                if (data.senderUUID === currentChatUUIDRef.current) {
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
    const baseURL = process.env.NEXT_PUBLIC_BACKEND_URL;
    const { isConnected, sendAction } = userUUID ? 
        useWebsocket(
        `${baseURL}/api/ws`,
        onMessage,
        {
            initialDelay: 1000,
            maxDelay: 1000,
            backoffFactor: 2,
            maxAttempts: 5,
        }
        ) : 
        {isConnected: false, sendAction: () => {}};


    useEffect(() => {
        setIsTyping(false);
    }, [currentChatUUID]);

    return (
        <WebSocketContext.Provider value={{ 
            isConnected, 
            userList, setUserList,
            messages, setMessages,
            sendAction, 
            isTyping,
            currentChatUUID, setCurrentChatUUID,
            currentGroupUUID, setCurrentGroupUUID,
            isLoadingMore, setIsLoadingMore,
            hasMore
        }}>
            {children}
        </WebSocketContext.Provider>
    );
}

export const useWebsocketContext = () => useContext(WebSocketContext);
