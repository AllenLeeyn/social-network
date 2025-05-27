"use client";

import { useWebsocketContext } from '../contexts/WebSocketContext';
import { useActiveChat } from '../contexts/ActiveChatContext';
import { useRouter } from 'next/navigation'

import "../styles/globals.css";

export default function UsersList( { activeConversation } ) {

    const { userList, isConnected, setUserList } = useWebsocketContext();
    const { setActiveChat } = useActiveChat();
    const router = useRouter();

    // Sorting: unread first, then online, then alphabetical
    const sortedUsers = [...userList].sort((a, b) => {
        // Unread messages first
        if (a.unread && !b.unread) return -1;
        if (!a.unread && b.unread) return 1;
        // Online users next
        if (a.online && !b.online) return -1;
        if (!a.online && b.online) return 1;
        // Alphabetical order
        return a.name.localeCompare(b.name);
    });

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
                {sortedUsers.map(user => (
                    <li
                        key={user.id}
                        role="button"
                        tabIndex={0}
                        className={[
                            'user-item',
                            user.online ? 'online' : '',
                            user.unread ? 'unread' : '',
                            activeConversation?.id === user.id ? 'active' : ''
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
        </div>
    );
    // return (
    //     <div className='sidebar-section'>
    //         <h3>({isConnected ? '✅ Connected' : '❌ Disconnected'})</h3>
    //         <ul className='users'>
    //             {userList.map(user => (
    //                 <li 
    //                     key={user.id}
    //                     className={`user-item ${user.online ? 'online' : ''} ${user.unread ? 'unread' : ''} ${activeConversation?.id === user.id ? 'active' : ''}`}
    //                     onClick={() => handleUserClick(user)}
    //                 >
    //                     {user.name}
    //                 </li>
    //             ))}
    //         </ul>
    //     </div>
    // );
}


