'use client';
import React, { useState } from 'react';
import { toast } from 'react-toastify';
import '../../styles/groups/CreateGroupForm.css'

export default function CreateGroupForm({ onSuccess }) {
    const [title, setTitle] = useState('');
    const [description, setDescription] = useState('');
    const [loading, setLoading] = useState(false);

    const [titleError, setTitleError] = useState('');
    const [descError, setDescError] = useState('');

    const validateTitle = (value) => {
        if (value.length < 3) return 'Title must be at least 3 characters.';
        if (value.length > 100) return 'Title must be at most 100 characters.';
        return '';
    };
    const validateDesc = (value) => {
        if (value.length < 10) return 'Description must be at least 10 characters.';
        if (value.length > 1000) return 'Description must be at most 1000 characters.';
        return '';
    };

    const handleTitleChange = (e) => {
        const value = e.target.value;
        setTitle(value);
        setTitleError(validateTitle(value));
    };
    const handleDescChange = (e) => {
        const value = e.target.value;
        setDescription(value);
        setDescError(validateDesc(value));
    };


    const handleSubmit = async (e) => {
        e.preventDefault();

        const tError = validateTitle(title);
        const dError = validateDesc(description);
        setTitleError(tError);
        setDescError(dError);
        if (tError || dError) return;

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
                console.log(res)
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
                <h4>Group Title</h4>
                <input
                    type="text"
                    value={title}
                    onChange={e => setTitle(e.target.value)}
                    placeholder='Write your Group Title'
                    required
                    minLength={3}
                    maxLength={100}
                    aria-invalid={!!titleError}
                    aria-describedby="title-error"
                />
                </label>
            </div>
            <div className="form-group">
                <label>
                <h4>Group Description</h4>
                <textarea
                    value={description}
                    onChange={e => setDescription(e.target.value)}
                    placeholder='Your group descrption'
                    required
                    minLength={10}
                    maxLength={1000}
                    aria-invalid={!!descError}
                    aria-describedby="desc-error"
                />
                </label>
            </div>
            <button 
                type="submit" 
                disabled={loading || !!titleError || !!descError}
            >
                {loading ? 'Creating...' : 'Create'}
            </button>
        </form>
    );
}
