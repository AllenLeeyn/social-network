"use client";

import { useWebsocketContext } from '../contexts/WebSocketContext';
import { useActiveChat } from '../contexts/ActiveChatContext';
import { useMemo } from 'react';

export default function DirectMessagesList() {
    const { userList, activeDM, setActiveDM } = useWebsocketContext();
    const { activeChat, setActiveChat } = useActiveChat();

    const activeChatUsers = useMemo(
        () => userList.filter(u => activeDM.includes(u.id)),
        [userList, activeDM]
    );

    return (
        <ul className='direct'>
            {activeChatUsers.map(user => (
                <li
                key={user.id}
                className={`conversation-item${activeChat?.id === user.id ? ' active' : ''}${user.unread ? ' unread' : ''}`}
                onClick={() => {
                    setActiveChat(user);
                        if (!activeDM.includes(user.id)) {
                            setActiveDM(prev => [...prev, user.id]);
                        }
                    }}
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
