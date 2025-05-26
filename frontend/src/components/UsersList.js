"use client";

import { useWebsocketContext } from '../contexts/WebSocketContext';
import { useActiveChat } from '../contexts/ActiveChatContext';
import { useRouter } from 'next/navigation'

import "../styles/globals.css";

export default function UsersList( { activeConversation } ) {

    const { userList, isConnected, setUserList } = useWebsocketContext();
    const { setActiveChat } = useActiveChat();
    const router = useRouter();

    const handleUserClick = (user) => {
        setActiveChat({ id: user.id, name: user.name });
        setUserList(prev => prev.map(u => u.id === user.id ? { ...u, unread: false } : u));
        console.log("User clicked:", user);
        router.push('/messages');
    }

    return (
        <div className='sidebar-section'>
            <h3>({isConnected ? '✅ Connected' : '❌ Disconnected'})</h3>
            <ul className='users'>
                {userList.map(user => (
                    <li 
                        key={user.id}
                        className={`user-item ${user.online ? 'online' : ''} ${user.unread ? 'unread' : ''} ${activeConversation?.id === user.id ? 'active' : ''}`}
                        onClick={() => handleUserClick(user)}
                    >
                        {user.name}
                    </li>
                ))}
            </ul>
        </div>
    );
}


