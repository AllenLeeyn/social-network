"use client";

import React, { useState } from "react";
import "../styles/CreatePost.css";
import { createPost } from "../lib/apiPosts";

export default function CreatePost({ categories, onClose }) {
  // state for title, content, selected categories (store names)
  const [title, setTitle] = useState("");
  const [content, setcontent] = useState("");
  const [selectedCategories, setSelectedCategories] = useState([]);

  // handleChange for inputs and checkboxes
  function handleCategoryChange(e) {
    const value = e.target.value; // use category name (string)
    setSelectedCategories((prev) =>
      prev.includes(value)
        ? prev.filter((cat) => cat !== value)
        : [...prev, value]
    );
  }

  // handlesubmit for the form
  async function handleSubmit(e) {
    e.preventDefault();
    // Map selected category names to IDs
    const categoryNameToId = {};
    categories.forEach((cat) => {
      categoryNameToId[cat.name] = cat.id;
    });
    const categoryIds = selectedCategories.map(
      (name) => categoryNameToId[name]
    );

    const postData = { title, content, category_ids: categoryIds};
    try {
      const data = await createPost(postData);
      if (data) {
        window.location.href = `/post?id=${data.data}`;
      } else {
        alert(data.message || "Failed to create post");
      }
    } catch (err) {
      alert(err.message || "Error creating post");
    }
    if (onClose) onClose();
  }

  return (
    <form onSubmit={handleSubmit}>
      {/* Title input */}
      <label htmlFor="title">Title</label>
      <div className="input-group">
        <input
          type="text"
          name="title"
          placeholder="Title"
          value={title}
          onChange={(e) => setTitle(e.target.value)}
          required
        />
      </div>
      {/* content textarea */}
      <label htmlFor="content">content</label>
      <div className="input-group">
        <textarea
          name="content"
          placeholder="Write your post here..."
          rows={10}
          value={content}
          onChange={(e) => setcontent(e.target.value)}
          required
        />
      </div>
      {/* Categories checkboxes */}
      <div className="input-group">
        <h4>Click to select categories</h4>
        <div className="checkbox-group">
          {categories.map((cat, index) => (
            <div className="checkbox-item" key={cat.id || index}>
              <input
                type="checkbox"
                id={`category${cat.id || index}`}
                name="categories"
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
      <div className="input-group">
        <button className="new-post" type="submit">
          Create Post
        </button>
      </div>
    </form>
  );
}
