"use client";
import { useWebsocketContext } from '../contexts/WebSocketContext';
import "../styles/globals.css";

export default function UsersList({ users }) {


    const { userList, isConnected} = useWebsocketContext();


    return (
        <ul className="users">
            {userList?.map((user) => (
                <li key={user.id} className={`user-item${user.online ? " online" : ""}${user.unread ? " unread" : ""}`}>
                </li>
            ))}
        </ul>
    );
}

//                 <span>{user.fullName} ({user.username})</span>
