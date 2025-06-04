import React, { useState } from "react";
import { toast } from 'react-toastify';
import { handleImage } from "../../../../lib/handleImage";

export default function CreatePostForm({ groupID, onSubmit, onClose }) {
    const [title, setTitle] = useState("");
    const [content, setContent] = useState("");

    const [postImages, setImages] = useState(null);

    const handleSubmit = async (e) => {
        e.preventDefault();
        // You can add validation here if needed
        
        let imageUUIDs = null;
        if (postImages) {
            try {
            imageUUIDs = await handleImage(postImages);
            } catch (err) {
            toast.error("Image upload failed: " + err.message);
            return;
            }
        }
        const postData = {
            title,
            content,
            category_ids: [0],
            group_id: groupID,
            file_attachments: imageUUIDs,
            visibility: "public",
            type:"group"};
        onSubmit && onSubmit(postData);
        onClose && onClose();
    };

    return (
        <form className="create-post-form" onSubmit={handleSubmit}>
        <h3>Create Post</h3>
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
            onChange={(e) => setContent(e.target.value)}
            required
            />
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
