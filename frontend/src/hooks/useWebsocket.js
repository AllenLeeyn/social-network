'use client';
import { useEffect, useRef, useCallback, useState } from "react";

// Add reconnection logic
export function useWebsocket(url, onMessage) {
    const ws = useRef(null);
    const [isConnected, setIsConnected] = useState(false);

    useEffect(() => {
        if (!url) return;

        ws.current = new window.WebSocket(url);

        ws.current.onopen = () => {
            setIsConnected(true);
            console.log('Websocket Connected');
        };

        ws.current.onclose = () => {
            setIsConnected(false);
            console.log('Websocket disconnected');
        };

        ws.current.onerror = (err) => {
            setIsConnected(false);
            console.error('WebSocket error', err);
        };
        
        ws.current.onmessage = (event) => {
            try {
                const data = JSON.parse(event.data);
                onMessage?.(data);
            } catch (e) {
                console.error('Invalid JSON', e);
            }
        };

        return () => {
            ws.current?.close();
        };
    }, [url, onMessage]);

    const sendMessage = useCallback((msg) => {
        if (ws.current?.readyState === WebSocket.OPEN) {
            ws.current.send(JSON.stringify(msg));
        }
    }, []);

    return { isConnected, sendMessage };
}