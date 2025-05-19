'use client';

import { createContext, useContext, useEffect, useRef, useState } from "react";

const WebSocketContext = createContext();

export function WebSocketProvider( {children} ) {
    const [userList, setUserList] = useState([]);
    const [isConnected, setIsConnected] = useState(false);
    const [connectionParams, setConnectionParams] = useState({
        sessionId: null,
        userId: null
    });
    const ws = useRef(null);

    // Main useEffect for ws lifecycle
    useEffect(() => {
        if (!connectionParams.sessionId || !connectionParams.userId) return;

        const connectWebSocket = () => {

            // if socket already open, return,
            if (ws.current?.readyState === WebSocket.OPEN) return;

            // Create new connection
            ws.current = new WebSocket(
                `ws://localhost:8080/ws?session=${connectionParams.sessionId}&user=${connectionParams.userId}`
            );

            // event handling
            ws.current.onopen = () => {
                setIsConnected(true);
                console.log('Global Websocket connected')
            };

            // actions
            ws.current.onmessage = (event) => {
                const data = JSON.parse(event.data);
                // Actions
                if (data.action === 'userListUpdate') {
                    setUserList(data.users);
                }
                // add more here
            };

            // when connection is closed, try to reestablish
            ws.current.onclose = () => {
                setIsConnected(false);
                console.log('Attempting reconnect...');
                setTimeout(connectWebSocket, 5000);
            };

            ws.current.onerror = (error) => {
                console.error('Websocket error:', error)
            };
        };
        
        
        connectWebSocket();
        
        return () => {
            if (ws.current?.readyState === WebSocket.OPEN) {
                ws.current.close();
            }
        };
    }, [connectionParams.sessionId, connectionParams.userId]); // reconnect when any of these changes

    // Connection Handler
    const connect = (sessionId, userId) => {
        setConnectionParams({ sessionId, userId });
    };

    const sendMessage = (message) => {
        if (ws.current?.readyState === WebSocket.OPEN) {
            ws.current.send(JSON.stringify(message));
        }
    };

    return (
        <WebSocketContext.Provider value={{
            userList,
            isConnected,
            connect,
            sendMessage,
        }}>
            {children}
        </WebSocketContext.Provider>
    );
}

export const useWebsocket = () => useContext(WebSocketContext);