"use client";

import React, { useState, useEffect } from "react";
import "../styles/CreatePost.css";
import { createPost, fetchCategories } from "../lib/apiPosts";
import { toast } from 'react-toastify';
import { handleImage } from "../lib/handleImage"; 
import { fetchFollowees } from "../lib/apiAuth";

export default function CreatePost({ onClose }) {
  const [title, setTitle] = useState("");
  const [content, setcontent] = useState("");
  const [selectedCategories, setSelectedCategories] = useState([]);
  const [categories, setCategories] = useState([]);
  const [postVisibility, setVisibility] = useState("");
  const [postImages, setImages] = useState(null);
  const [followers, setFollowers] = useState([]);
  const [selectedUsers, setSelectedUsers] = useState([]);

  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    async function loadData() {
      try {
        const categoriesData = await fetchCategories();
        setCategories(categoriesData.data);

        let nickname = null
        if (typeof window !== 'undefined') {
          nickname = localStorage.getItem('user-nick_name');
        }

        const followerData = await fetchFollowees();
        const filteredFollowers = (followerData || []).filter(
          (follower) => follower.follower_name !== nickname
        );
        setFollowers(filteredFollowers || []);

      } catch (err) {
        setError(err)
        toast.error(err.message);
      } finally {
        setLoading(false);
      }
    }
    loadData();
  }, []);

  function handleCategoryChange(e) {
    const value = e.target.value;
    setSelectedCategories((prev) =>
      prev.includes(value)
        ? prev.filter((cat) => cat !== value)
        : [...prev, value]
    );
  }

  function handleUserSelect(e) {
    const value = e.target.value;
    setSelectedUsers((prev) =>
      prev.includes(value)
        ? prev.filter((id) => id !== value)
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
      visibility: postVisibility,
      selected_audience_user_uuids: selectedUsers};
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
        <option value="">Select visibility</option>
        <option value="public">Public</option>
        <option value="private">Private</option>
        <option value="selected">Select users</option>
      </select>
      {postVisibility === "selected" && (
        <div className="input-group">
          {followers.length === 0 ? (
            <p>No followers found.</p>
          ) : (
            <ul className="follower-group">
              {followers.map((f) => (
                <li key={f.follower_uuid} className="follower-item">
                  <input
                    type="checkbox"
                    id={`user-${f.follower_uuid}`}
                    value={f.follower_uuid}
                    checked={selectedUsers.includes(f.follower_uuid)}
                    onChange={handleUserSelect}
                  />
                  <label htmlFor={`user-${f.follower_uuid}`}>{f.follower_name}</label>
                </li>
              ))}
            </ul>
          )}
        </div>
      )}

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
