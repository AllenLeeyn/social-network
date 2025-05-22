"use client";
import { useWebsocketContext } from '../contexts/WebSocketContext';
import "../styles/globals.css";

export default function UsersList( { onUserSelect, activeConversation } ) {

    const { userList, isConnected } = useWebsocketContext();

    return (
        <div className='sidebar-section'>
            <h3>({isConnected ? '✅ Connected' : '❌ Disconnected'})</h3>
            <ul className='users'>
                {userList.map(user => (
                    <li 
                        key={user.id}
                        className={`user-item ${user.online ? 'online' : ''} ${user.unread ? 'unread' : ''} ${activeConversation?.id === user.id ? 'active' : ''}`}
                        onClick={() => onUserSelect && onUserSelect(user.id, user.name)}
                    >
                        {user.name}
                    </li>
                ))}
            </ul>
        </div>
    );
}


