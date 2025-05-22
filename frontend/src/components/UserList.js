"use client";
import "../styles/globals.css";

export default function UsersList({ users }) {
    return (
        <ul className="users">
        {users.map((user) => (
            <li key={user.id} className={`user-item${user.online ? " online" : ""}${user.unread ? " unread" : ""}`}>
            <img src={user.avatar} alt={user.username} />
            <span>{user.fullName} ({user.username})</span>
            </li>
        ))}
        </ul>
    );
}