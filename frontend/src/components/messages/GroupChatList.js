"use client";

import { useActiveChat } from '../../contexts/ActiveChatContext';
import { useMemo } from 'react';

export default function GroupChatsList() {
    const { activeChat, setActiveChat } = useActiveChat();

    const groups = useMemo(
        () => {},
        []
    );  

    return /* (
        <ul className='groups'>
            {groups.map(convo => (
                <li
                key={convo.id}
                className={`conversation-item${activeChat?.id === convo.id ? ' active' : ''}${convo.unread ? ' unread' : ''}`}
                onClick={() => setActiveChat(convo)}
                >
                    {convo.name}
                    {convo.unread > 0 && <span className='unread-count'>{convo.unread}</span>}
                </li>
            ))}
        </ul>
    ); */
}
