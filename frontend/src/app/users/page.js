"use client";

import { useState, useEffect } from "react";
import { fetchUsers } from "../../lib/apiAuth";
import UserCard from "../../components/UserCard";
import "./users.css";

export default function UsersPage() {
  const [users, setUsers] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    async function loadUsers() {
      try {
        const data = await fetchUsers();
        setUsers(data || []);
      } catch (err) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    }
    loadUsers();
  }, []);

  if (loading) return <div className="loading">Loading users...</div>;
  if (error) return <div className="error">Error: {error}</div>;

  return (
    <main className="users-page">
      <div className="users-header">
        <h2>All Users ({users.length})</h2>
      </div>
      <div className="users-grid">
        {users.map((user) => (
          <UserCard key={user.uuid} user={user} />
        ))}
      </div>
    </main>
  );
}
