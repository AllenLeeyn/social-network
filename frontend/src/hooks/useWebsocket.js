'use client';
import { useEffect, useRef, useCallback } from "react";

export function useWebsocket(url, onMessage) {
    const ws = useRef(null);

    useEffect(() => {
        ws.current = new window.WebSocket(url);

        ws.current.onopen = () => console.log('Websocket Connected');
        ws.current.onclose = () => console.log('Websocket disconnected');
        ws.current.onerror = (err) => console.error('WebSocket error', err);
        
        ws.current.onmessage = (event) => {
            try {
                const data = JSON.parse(event.data);
                onMessage?.(data);
            } catch (e) {
                console.error('Invalid JSON', e);
            }
        };

        return () => ws.current?.close();
    }, [url, onMessage]);

    const sendMessage = useCallback((msg) => {
        if (ws.current?.readyState === WebSocket.OPEN) {
            ws.current.send(JSON.stringify(msg));
        }
    }, []);

    return { sendMessage }
}