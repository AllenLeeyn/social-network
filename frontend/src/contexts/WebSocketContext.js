'use client';

import { createContext, useContext, useCallback, useState } from "react";
import { useWebsocket } from "../hooks/useWebsocket";


const WebSocketContext = createContext();

export function WebSocketProvider( { children } ) {
    const [userList, setUserList] = useState([]);
    const [sessionId, setSessionId] = useState(null);

    // using hook for connection logic
    const { isConnected, sendAction } = useWebsocket(
        sessionId ? `ws://localhost:8080/ws?session=${sessionId}` : null,
        (data) => {
            if (data.action === 'userListUpdate') {
                setUserList(data.users)
            }
            // Add more msg handler 
        },
        {
            initialDelay: 1000,
            maxDelay: 30000,
            backoffFactor: 2,
            maxAttempts: null
        }
    );
    // console.log('sending Session ID', sessionId)

    // conn handling
    const connect = useCallback((newSessionId) => {
        setSessionId(newSessionId);
    }, []);

    return (
        <WebSocketContext.Provider value ={{
            userList,
            isConnected,
            connect,
            sendAction
        }}>
            {children}
        </WebSocketContext.Provider>
    );
}

export const useWebsocketContext = () => useContext(WebSocketContext);

