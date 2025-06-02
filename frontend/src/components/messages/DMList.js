"use client";

import { useWebsocketContext } from '../../contexts/WebSocketContext';
import { useActiveChat } from '../../contexts/ActiveChatContext';
import { useMemo } from 'react';

import '/src/styles/messages/DMList.css'

export default function DirectMessagesList() {
/*     const { userList, activeDM, setActiveDM } = useWebsocketContext();
    const { activeChat, setActiveChat } = useActiveChat();

    const activeChatUsers = useMemo(
        () => userList.filter(u => activeDM.includes(u.id)),
        [userList, activeDM]
    );

    return (
        <ul className='direct'>
            {activeChatUsers.map(user => {
                const classes = [
                    "dm-item",
                    activeChat?.id === user.id && "active",
                    user.unread && "unread"
                ]
                .filter(Boolean)
                .join(" ");

                return (
                    <li
                        key={user.id}
                        className={classes}
                        onClick={() => {
                            setActiveChat(user);
                            if (!activeDM.includes(user.id)) {
                                setActiveDM(prev => [...prev, user.id]);
                            }
                        }}
                    >
                        {user.name}
                        {user.online && <span className="dot online" />}
                        {user.unread > 0 && <span className="unread-count">{user.unread}</span>}
                    </li>
                );
            })}
            {activeChatUsers.length === 0 && (
                <li className="no-active-chats">No active chats yet</li>
            )}
        </ul>
    ); */
}
