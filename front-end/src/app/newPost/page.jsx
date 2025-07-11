"use client";

import { useState, useRef } from "react";
import Styles from "./newPost.module.css";
import { SendData } from "../sendData.js";
import { useRouter } from "next/navigation";
import { HandleUpload } from "../utils.jsx";

export default function NewPost() {
  const [content, setContent] = useState("");
  const [image, setImage] = useState(null);
  const [privacy, setPrivacy] = useState("public");
  const [previewUrl, setPreviewUrl] = useState(null);
  const fileInputRef = useRef(null);
  const router = useRouter();

  const handleImageChange = (e) => {
    const file = e.target.files[0];
    if (file) {
      setImage(file);
      setPreviewUrl(URL.createObjectURL(file));
    }
  };

  const cancelImage = () => {
    setImage(null); // clear the file from state
    setPreviewUrl(null); // remove the preview URL
    if (fileInputRef.current) {
      fileInputRef.current.value = ""; // reset the <input>
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();

    if (!content.trim()) {
      console.log("Content is required.");
      return;
    }

    let imagePath = null;
    if (image) {
      imagePath = await HandleUpload(image);
      console.log("path:", imagePath);
    }

    const formData = {
      content: content,
      privacy: privacy,
      image: imagePath,
    };

    const response = await SendData("/api/v1/set/Post", formData);

    if (response.status !== 200) {
      const errorBody = await response.json();
      console.log(errorBody);
    } else {
      console.log("✅ Post created!");
      router.push("/");
    }
  };

  return (
    <div className={Styles.form}>
      <h2>Create New Post</h2>

      <form onSubmit={handleSubmit}>
        {/* Content Field */}
        <div>
          <label htmlFor="content">Content</label>
          <br />
          <textarea
            id="content"
            rows="4"
            value={content}
            onChange={(e) => setContent(e.target.value)}
            placeholder="What's on your mind?"
          />
        </div>

        {/* Upload Image */}
        <div className={Styles.upload}>
          <label htmlFor="image" style={{ cursor: "pointer" }}>
            <img src="/Image.svg" alt="Upload" width="24" height="24" />
            &nbsp;&nbsp; Upload Image
          </label>
          <input
            type="file"
            name="image"
            id="image"
            style={{ display: "none" }}
            accept="image/*,video/*"
            onChange={handleImageChange}
            ref={fileInputRef}
          />
          {previewUrl && (
            <div className={Styles.previewContainer}>
              <img src={previewUrl} alt="Preview" />
              <button
                type="button"
                className={Styles.cancelButton}
                onClick={cancelImage} // ← wire up cancel
              >
                ✕
              </button>
            </div>
          )}
        </div>

        {/* Dropdown Privacy */}
        <div className={Styles.dropdown}>
          <label htmlFor="privacy">Privacy</label>
          <br />
          <select
            id="privacy"
            className={Styles.input}
            value={privacy}
            onChange={(e) => setPrivacy(e.target.value)}
            style={{ padding: "0.5rem", marginTop: "0.5rem" }}
          >
            <option value="public">Public</option>
            <option value="private">Private</option>
            <option value="folowors">Folowors</option>
          </select>
        </div>

        <button type="submit">Post</button>
      </form>
    </div>
  );
}
