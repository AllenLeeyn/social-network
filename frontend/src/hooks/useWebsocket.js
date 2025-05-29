'use client';
import { useEffect, useRef, useCallback, useState } from "react";

// Add reconnection logic
export function useWebsocket(
    url,
    onMessage,
    {
        initialDelay = 1000,         // 1 second delay
        maxDelay = 30000,            // 30 seconds
        backoffFactor = 2,           // Double each time it tries to reconnect
        maxAttempts = null,          // unlimited attempts by default
    } = {}
) {
    // Refs for persistent values between renders
    const ws = useRef(null);
    const reconnectDelay = useRef(initialDelay);
    const attempts = useRef(0);
    const reconnectTimeout = useRef(null);
    const isMounted = useRef(true);
    const [isConnected, setIsConnected] = useState(false);

    // Core conn / re-conn logic
    const connectWebSocket = useCallback(() => {
        // clear pending reconn
        clearTimeout(reconnectTimeout.current);


        // Close old connection if it exists and is not already closed
        if (ws.current && ws.current.readyState !== WebSocket.CLOSED) {
            ws.current.close();
            ws.current = null;
        }

        // no url, return
        if (!url) {
            console.warn('Skipping WebSocket connection - invalid URL');
            return;
        }

        // max retries limit 
        if (maxAttempts !== null && attempts.current >= maxAttempts) {
            console.warn('Max reconect attempts reached');
            return;
        }
        
        attempts.current += 1; // track conn attempts
        // setting new ws conn
        ws.current = new WebSocket(url);



        // conn established handler
        ws.current.onopen = () => {
            if (!isMounted.current) return;
            setIsConnected(true);
            attempts.current = 0;
            reconnectDelay.current = initialDelay;
            console.log('Websocket Connected');
        }

        // close conn handler
        ws.current.onclose = () => {
            if (!isMounted.current) return;
            setIsConnected(false);
            const jitter = Math.random() * 1000;
            const nextDelay = Math.min(reconnectDelay.current * backoffFactor, maxDelay) + jitter;
            reconnectDelay.current = nextDelay;

            console.log(`Reconnecting in ${Math.round(nextDelay)}ms...`);
            reconnectTimeout.current = setTimeout(connectWebSocket, nextDelay);
        };

        // err handler
        ws.current.onerror = (err) => {
            if (!isMounted.current) return;
            console.error('WebSocket error:', err);
            ws.current?.close();        // triggers onclose and goes into re-conn
        }

        // msg handler
        ws.current.onmessage = (event) => {
            try {
                const data = JSON.parse(event.data);
                onMessage?.(data); // Pass raw data to context
            } catch (e) {
                console.error('Invalid JSON', e);
            }
        };
    }, [ url, onMessage, initialDelay, maxDelay, backoffFactor, maxAttempts ]);

    // Start connection and handling clean up
    // Manage conn lifecycle
    // useEffect / onRender

    useEffect(() => {
        isMounted.current = true;
        connectWebSocket();

        // cleanup func
        return () => {
            isMounted.current = false;
            clearTimeout(reconnectTimeout.current);
            // Only close if connection is open/connecting
            if (ws.current?.readyState !== WebSocket.CLOSED  ) {
                ws.current?.close();
            }
        };
    }, [ connectWebSocket ]);


    // message sending func
    const sendAction = useCallback((msg) =>  {
        if (ws.current?.readyState === WebSocket.OPEN) {
            ws.current.send(JSON.stringify(msg));
        } else {
            console.warn('Cannot send message - WebSocket not open');
        }
    }, []);

    return { isConnected, sendAction };
}

