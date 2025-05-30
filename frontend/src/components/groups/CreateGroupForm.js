'use client';
import React, { useState } from 'react';
import { toast } from 'react-toastify';
import '../../styles/groups/CreateGroupForm.css'

export default function CreateGroupForm({ onSuccess }) {
    const [title, setTitle] = useState('');
    const [description, setDescription] = useState('');
    const [loading, setLoading] = useState(false);


    const handleSubmit = async (e) => {
        e.preventDefault();
        setLoading(true);
        try {
            const res = await fetch('/frontend-api/groups/create', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ title, description }),
            });
            setLoading(false);
            if (res.ok) {
                toast.success('Group created!');
                if (onSuccess) onSuccess();
            } else {
                toast.error('Failed to create group.');
            }
            } catch (err) {
            setLoading(false);
            toast.error('Network error.');
        }
    };

    return (
        <form onSubmit={handleSubmit}>
            <div className="form-group">
                <label>
                Group Title
                <input
                    type="text"
                    value={title}
                    onChange={e => setTitle(e.target.value)}
                    required
                />
                </label>
            </div>
            <div className="form-group">
                <label>
                Group Description
                <textarea
                    value={description}
                    onChange={e => setDescription(e.target.value)}
                    required
                />
                </label>
            </div>
            <button type="submit" disabled={loading}>
                {loading ? 'Creating...' : 'Create'}
            </button>
        </form>
        );
}
