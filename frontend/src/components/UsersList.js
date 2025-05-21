"use client";
import { useWebsocketContext } from '../contexts/WebSocketContext';
import "../styles/globals.css";

export default function UsersList({ users }) {

    const { isConnected } = useWebsocketContext();
    // const { userList, isConnected} = useWebsocketContext();

    // For testing, you can override userList with mock data
    const userList = [
      { id: 1, name: 'Alice', online: true },
      { id: 2, name: 'Bob', online: false },
    ];

    return (
        <div>
            <h3>({isConnected ? '✅ Connected' : '❌ Disconnected'})</h3>
            <ul>
                {userList.map(user => (
                <li key={user.id}>{user.name}</li>
                ))}
            </ul>
        </div>
    );
}

//                 <span>{user.fullName} ({user.username})</span>
