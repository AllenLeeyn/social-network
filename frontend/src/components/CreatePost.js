'use client';

import React, { useState } from 'react';
import "../styles/CreatePost.css"

export default function CreatePost({categories, onClose }) {
    // state for title, content, selected categories (store names)
    const [ title, setTitle ] = useState('');
    const [ content, setContent ] = useState('');
    const [ selectedCategories, setSelectedCategories ] = useState([]);

    // handleChange for inputs and checkboxes
    function handleCategoryChange(e) {
        const value = e.target.value; // use category name (string)
        setSelectedCategories(prev =>
            prev.includes(value)
                ? prev.filter(cat => cat !== value)
                : [...prev, value]
        );
    }

    // handlesubmit for the form
    function handleSubmit(e) {
        e.preventDefault();
        // Form Data
        const postData = { title, content, categories: selectedCategories };
        // do something with postData // send to api or update or something
        if (onClose) onClose();
    }

    return (
        <form onSubmit ={handleSubmit}>
            {/* Title input */}
            <label htmlFor="title">Title</label>
            <div className='input-group'>
                <input
                    type='text'
                    name='title'
                    placeholder='Title'
                    value={title}
                    onChange={e => setTitle(e.target.value)}
                    required
                />
            </div>
            {/* Content textarea */}
            <label htmlFor="content">Content</label>
            <div className='input-group'>
                <textarea
                    name='content'
                    placeholder='Write your post here...'
                    rows={10}
                    value={content}
                    onChange={e => setContent(e.target.value)}
                    required
                />
            </div>
            {/* Categories checkboxes */}
            <div className='input-group'>
                <h4>Click to select categories</h4>
                <div className='checkbox-group'>
                    {categories.map((cat, index) => (
                        <div className='checkbox-item' key={cat.id || index}>
                            <input
                                type='checkbox'
                                id={`category${cat.id || index}`}
                                name='categories'
                                value={cat.name}
                                checked={selectedCategories.includes(cat.name)}
                                onChange={handleCategoryChange}
                            />
                            <label htmlFor={`category${cat.id || index}`}>{cat.name}</label>
                        </div>
                    ))}
                </div>
            </div>
            {/* Submit button */}
            <div className='input-group'>
                <button className='new-post' type='submit'>
                    Create Post
                </button>
            </div>
        </form>
    )
}



