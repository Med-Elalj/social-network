"use client";

import { useState } from "react";

export default function NewPost() {
  const [content, setContent] = useState("");
  const [image, setImage] = useState(null);
  const [previewUrl, setPreviewUrl] = useState(null);
  const [error, setError] = useState("");
  const [success, setSuccess] = useState("");

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

    const formData = new FormData();
    formData.append("content", content);
    if (image) formData.append("image", image);

    try {
      const res = await fetch("/api/v1/posts", {
        method: "POST",
        body: formData,
      });

      const data = await res.json();

      if (res.ok) {
        setSuccess("Post created successfully!");
        setContent("");
        setImage(null);
        setPreviewUrl(null);
        setError("");
      } else {
        setError(data.message || "Failed to create post.");
      }
    } catch (err) {
      setError("Something went wrong.");
    }
  };

  return (
    <div style={{ maxWidth: 600, margin: "60px auto", padding: 20 }}>
      <h2>Create New Post</h2>

      {error && <p style={{ color: "red" }}>{error}</p>}
      {success && <p style={{ color: "green" }}>{success}</p>}

      <form onSubmit={handleSubmit}>
        <div style={{ marginBottom: 15 }}>
          <label htmlFor="content">Content</label><br />
          <textarea
            id="content"
            rows="4"
            value={content}
            onChange={(e) => setContent(e.target.value)}
            style={{ width: "100%", padding: 10 }}
            placeholder="What's on your mind?"
          />
        </div>

        <div style={{ marginBottom: 15 }}>
          <label htmlFor="image">Upload Image</label><br />
          <input type="file" id="image" accept="image/*" onChange={handleImageChange} />
          {previewUrl && (
            <div style={{ marginTop: 10 }}>
              <img src={previewUrl} alt="Preview" style={{ maxWidth: "100%", height: "auto" }} />
            </div>
          )}
        </div>

        <button type="submit" style={{ padding: "10px 20px" }}>Post</button>
      </form>
    </div>
  );
}
