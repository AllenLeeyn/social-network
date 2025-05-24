"use client";

import { useWebsocketContext } from '../contexts/WebSocketContext';
import { useActiveChat } from '../contexts/ActiveChatContext';
import { useMemo } from 'react';

export default function DirectMessagesList() {
    const { messages, userList } = useWebsocketContext();
    const { activeChat, setActiveChat } = useActiveChat();

    // Find all user IDs you've chatted with
    const activeChatUserIds = useMemo(() => {
        const ids = new Set();
        messages.forEach(m => {
            if (m.senderId !== undefined) ids.add(m.senderId);
            if (m.receiverId !== undefined) ids.add(m.receiverId);
        });
        return Array.from(ids);
    }, [messages]);

    const activeChatUsers = useMemo(
        () => userList.filter(u => activeChatUserIds.includes(u.id)),
        [userList, activeChatUserIds]
    );

    return (
        <ul className='direct'>
            {activeChatUsers.map(user => (
                <li
                key={user.id}
                className={`conversation-item${activeChat?.id === user.id ? ' active' : ''}${user.unread ? ' unread' : ''}`}
                onClick={() => setActiveChat(user)}
                >
                {user.name}
                {user.online && <span className="dot online" />}
                {user.unread > 0 && <span className='unread-count'>{user.unread}</span>}
                </li>
            ))}
            {activeChatUsers.length === 0 && (
                <li className="no-active-chats">No active chats yet</li>
            )}
        </ul>
    );
}
