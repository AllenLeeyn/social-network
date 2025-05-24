"use client";

import { useActiveChat } from '../contexts/ActiveChatContext';
import { useMemo } from 'react';
import { sampleConversations } from '../data/mockData';

export default function GroupChatsList() {
    const { activeChat, setActiveChat } = useActiveChat();

    const groups = useMemo(
        () => sampleConversations.filter(c => c.type === 'group'),
        []
    );  

    return (
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
    );
}
