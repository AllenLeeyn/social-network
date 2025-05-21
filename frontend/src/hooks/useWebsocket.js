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
    const [isConnected, setIsConnected] = useState(false);

    // Core conn / re-conn logic
    const connectWebSocket = useCallback(() => {
        // no url, return
        if (!url || url.includes('session=undefined')) return;

        // max retries limit 
        if (maxAttempts !== null && attempts.current >= maxAttempts) {
            console.warn('Max reconect attempts reached');
            return;
        }

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
            attempts.current += 1;

            // calc new delay(exponential backoff with max cap)
            const nextDelay = Math.min(reconnectDelay.current * backoffFactor, maxDelay);
            reconnectDelay.current = Math.min(nextDelay, maxDelay);

            console.log(`Websocket Disconnected. Reconnecting in ${reconnectDelay.current}ms...`);
            reconnectTimeout.current = setTimeout(connectWebSocket, reconnectDelay.current);
        };

        // err handler
        ws.current.onerror = (err) => {
            if (!isMounted.current) return;
            setIsConnected(false);
            // console.error('Websocket error', err);
            if (err && Object.keys(err).length > 0) {
                    console.error('Websocket error:', err);
                }
            ws.current?.close();        // triggers onclose and goes into re-conn
        }

        // msg handler
        ws.current.onmessage = (event) => {
            try {
                const data = JSON.parse(event.data);
                onMessage?.(data);
            } catch (e) {
                console.error('Invalid JSON', e);
            }
        };
    }, [ url, onMessage, initialDelay, maxDelay, backoffFactor, maxAttempts ]);

    // Start connection and handling clean up
    // Manage conn lifecycle
    // useEffect / onRender

    const isMounted = useRef(true);

    useEffect(() => {
        isMounted.current = true;
        connectWebSocket();

        // cleanup func
        return () => {
            isMounted.current = false;
            clearTimeout(reconnectTimeout.current);
            ws.current?.close();
        };
    }, [ connectWebSocket ]);

    // message sending func
    const sendAction = useCallback((msg) =>  {
        if (ws.current?.readyState === WebSocket.OPEN) {
            ws.current.send(JSON.stringify(msg));
        }
    }, []);

    return { isConnected, sendAction };
}

