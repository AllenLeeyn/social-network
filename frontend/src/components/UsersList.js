"use client";

import { useWebsocketContext } from '../contexts/WebSocketContext';
import { useActiveChat } from '../contexts/ActiveChatContext';
import { useRouter } from 'next/navigation'

import "../styles/globals.css";

export default function UsersList( { activeConversation } ) {

    const userUUID = typeof window !== 'undefined' ? localStorage.getItem('user-uuid') : null;
    const { userList, isConnected, setUserList, setMessages, sendAction } = useWebsocketContext();
    const { activeChat, setActiveChat } = useActiveChat();
    const router = useRouter();
    
    const handleUserClick = (user) => {
        setActiveChat({ 
            type: user.type,
            uuid: user.uuid,
            name: user.name,
            receiverUUID: user.receiverUUID,
            groupUUID: user.groupUUID,
        });

        setMessages([]);
        sendAction({
            action: "messageReq",
            receiverUUID: user.receiverUUID,
            groupUUID: user.groupUUID,
            content: "-1"
        });
        sendAction({
            action: 'messageAck',
            receiverUUID: userUUID,
            senderUUID: user.receiverUUID
        });
        setUserList(prev => prev.map(u => u.uuid === user.uuid ? { ...u, unread: false } : u));
        console.log("User clicked:", user);
        router.push('/messages');
    }

    const users = userList.filter(user => user.type === "user");
    const groups = userList.filter(user => user.type === "group");
    
    return (
        <div className='sidebar-section'>
            <h3>({isConnected ? '✅ Connected' : '❌ Disconnected'})</h3>
            <h4>Users</h4>
            <ul className='users'>
                {users.map(user => (
                    <li
                        key={user.uuid}
                        role="button"
                        tabIndex={0}
                        className={[
                            'user-item',
                            user.online ? 'online' : '',
                            user.unread ? 'unread' : '',
                            activeConversation?.uuid === user.uuid ? 'active' : ''
                        ].filter(Boolean).join(' ')}
                        onClick={() => handleUserClick(user)}
                        onKeyDown={e => {
                            if (e.key === 'Enter' || e.key === ' ') handleUserClick(user);
                        }}
                    >
                        {user.name}
                        {user.unread > 0 && <span className='unread-count'>{user.unread}</span>}
                        {user.online && <span className="dot online" />}
                    </li>
                ))}
            </ul>

            <h4>Groups</h4>
            <ul className='groups'>
                {groups.map(group => (
                    <li
                        key={group.uuid}
                        role="button"
                        tabIndex={0}
                        className={'group-item'}
                        onClick={() => handleUserClick(group)}
                        onKeyDown={e => {
                            if (e.key === 'Enter' || e.key === ' ') handleUserClick(group);
                        }}
                    >
                        {group.name}
                    </li>
                ))}
            </ul>
        </div>
    );
}


