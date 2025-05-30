"use client";

import React, { useState, useEffect } from "react";
import "../styles/CreatePost.css";
import { createPost, fetchCategories } from "../lib/apiPosts";
import { toast } from 'react-toastify';
import { handleImage } from "../lib/handleImage"; 

export default function CreatePost({ onClose }) {
  const [title, setTitle] = useState("");
  const [content, setcontent] = useState("");
  const [selectedCategories, setSelectedCategories] = useState([]);
  const [categories, setCategories] = useState([]);
  const [postVisibility, setVisibility] = useState("");
  const [postImages, setImages] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    async function loadCategories() {
      try {
        const data = await fetchCategories();
        setCategories(data.data); // Adjust if your API returns a different shape
      } catch (err) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    }
    loadCategories();
  }, []);

  function handleCategoryChange(e) {
    const value = e.target.value;
    setSelectedCategories((prev) =>
      prev.includes(value)
        ? prev.filter((cat) => cat !== value)
        : [...prev, value]
    );
  }

  async function handleSubmit(e) {
    e.preventDefault();
    const categoryNameToId = {};
    categories.forEach((cat) => {
      categoryNameToId[cat.name] = cat.id;
    });
    const categoryIds = selectedCategories.map(
      (name) => categoryNameToId[name]
    );


    let imageUUIDs = null;
    if (postImages) {
      try {
        imageUUIDs = await handleImage(postImages);
      } catch (err) {
        toast.error("Image upload failed: " + err.message);
        return;
      }
    }

    const postData = { title, content, 
      category_ids: categoryIds, 
      file_attachments: imageUUIDs, 
      visibility: postVisibility};
    try {
      const data = await createPost(postData);
      if (data) {
        window.location.href = `/post?id=${data.data}`;
      } else {
        toast.error(data.message || "Failed to create post");
      }
    } catch (err) {
      toast.error(err.message || "Error creating post");
    }
    if (onClose) onClose();
  }

  if (loading) return <div>Loading categories...</div>;
  if (error) return <div>Error loading categories: {error}</div>;

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
      <select
        value={postVisibility}
        onChange={(e) => setVisibility(e.target.value)}
        required
      >
        <option value="public">Public</option>
        <option value="private">Private</option>
        <option value="selected">Select users</option>
      </select>
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
      <input
        type="file"
        accept="image/*"
        multiple={true}
        onChange={(e) => setImages([...e.target.files])}
      />
      {/* Submit button */}
      <div className="input-group">
        <button className="new-post" type="submit">
          Create Post
        </button>
      </div>
    </form>
  );
}
