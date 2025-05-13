'use client';

import React, { useState } from 'react';

export default function CreatePost({categories, onClose }) {
    // state for title, content, selected categories
    const [ title, setTitle ] = useState('');
    const [ content, setContent ] = useState('');
    const [ selectedCategories, setSelectedCategories ] = useState([]);

    // handleChange for inputs and checkboxes
    function handleCategoryChange(e) {
        const value = e.target.value;

        setSelectedCategories( prev => 
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
                    {categories.map((cat,index) => (
                        <div className='checkbox-item' key={cat.id || index}>
                            <input
                                type='checkbox'
                                id={`category${index}`}
                                name='categories'
                                value={cat.id || cat.name}
                                checked={selectedCategories.includes(cat.id || cat.name)}
                                onChange={handleCategoryChange}
                            />
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





// export const templateNewPost = (categories) => `
//   <form id="newPostForm">
//     <div class="input-group">
//       <input type="text" name="title" placeholder="Title" required />
//     </div>
//     <div class="input-group">
//       <textarea name="content" placeholder="Write your content here..." rows="10" required ></textarea>
//     </div>
//     <div class="input-group">
//       <h4>Click to select categories:</h4>
//       <div class="checkbox-group">
//         ${categories.map((cat, index) => `
//           <div class="checkbox-item">
//             <input type="checkbox" id="category${index}" name="categories" value=${index}>
//             <label for="category${index}">${cat}</label>
//           </div>
//           `).join('')}
//       </div>
//     </div>
//     <div class="input-group">
//       <button class="new-post" id="newPostSubmit" type="submit">Create Post</button>
//     </div>
//   </form>`;
