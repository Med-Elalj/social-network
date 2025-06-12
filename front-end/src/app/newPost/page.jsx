"use client";

import { useState } from "react";
import Styles from "./newPost.module.css";

export default function NewPost() {
  const [content, setContent] = useState("");
  const [image, setImage] = useState(null);
  const [previewUrl, setPreviewUrl] = useState(null);

  const handleImageChange = (e) => {
    const file = e.target.files[0];
    if (file) {
      setImage(file);
      setPreviewUrl(URL.createObjectURL(file));
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();

    if (!content.trim()) {
      setError("Content is required.");
      return;
    }
    console.log(image);
    console.log(content);
  };

  return (
    <div className={Styles.form}>
      <h2>Create New Post</h2>

      <form onSubmit={handleSubmit}>
        <div>
          <label htmlFor="content">Content</label><br />
          <textarea
            id="content"
            rows="4"
            value={content}
            onChange={(e) => setContent(e.target.value)}
            placeholder="What's on your mind?"
          />
        </div>

        <div>
          <label htmlFor="image">Upload Image</label><br />
          <input type="file" id="image" accept="image/*" onChange={handleImageChange} />
          {previewUrl && (
            <div >
              <img src={previewUrl} alt="Preview" />
            </div>
          )}
        </div>

        <button type="submit">Post</button>
      </form>
    </div>
  );
}
