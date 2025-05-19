import { useState } from "react";
import { useWebsocket } from "../hooks/useWebsocket";

export default function ChatComponent({ sessionId }) {
    const [messages, setMessages] = useState([]);
    const [inputMessage, setInputMessage] = useState("");

    const handleAction = (data) => {
        // Handle different actions. message, typing, and other we implement
        if (data.action === 'message') {
            setMessages((msgs) => [...msgs,data]);
        }
        // Handle other actions
    };

    const { sendMessage } = useWebsocket(
        `ws://localhost:8080/ws?session=${sessionId}`, handleAction
    );

    const handleSendMessage = (e) => {
        e.preventDefault();
        if (!inputMessage.trim()) return;

        const newMessage = {
            action: "message",
            content: inputMessage,
            timestamp: new Date().toISOString(),
            // can add more keys: value for data
        };

        sendMessage(JSON.stringify(newMessage));
        setMessages((prev) => [...prev, newMessage]);
        setInputMessage("");

    };

    return (
        <div className='chat-component'>
            <div className='messages-list'>
                {messages.map((msg, index) => {
                    return (
                    <div key={index} className='message-item'>
                        <span><strong>{msg.senderId}</strong>:{msg.content}</span>
                        <span className='timestamp'>{new Date(msg.timestamp).toLocaleTimeString()}</span>
                    </div>
                    )
                })}
            </div>      

            <form onSubmit={handleSendMessage} className='message-input-form'>
                <input
                    type='text'
                    value={inputMessage}
                    onChange={(e) => setInputMessage(e.target.value)}
                    placeholder='Type a message...'
                    className='message-input'
                />
            </form>
        </div>
    );
}