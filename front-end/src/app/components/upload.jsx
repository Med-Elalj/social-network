"use client";
import { useState } from "react";
import { fetchWithAuth } from "../sendData";

// Reusable Upload Function
export async function HandleUpload(image) {
  if (!image) return null;

  const formData = new FormData();
  formData.append("file", image);

  try {
    const response = await fetchWithAuth("/api/v1/upload", {
      method: "POST",
      body: formData,
    });

    if (!response.ok) {
      console.error("Image upload failed with status:", response.status);
      return null;
    }

    const data = await response.json();
    return data?.path || null;

  } catch (error) {
    console.error("Upload error:", error);
    return null;
  }
}

// Upload Form Component
export default function UploadForm() {
  const [file, setFile] = useState(null);
  const [message, setMessage] = useState("");
  const [imageUrl, setImageUrl] = useState("");

  const handleFileChange = (e) => {
    setFile(e.target.files[0]);
    setImageUrl("");
    setMessage("");
  };

  const handleUpload = async () => {
    if (!file) {
      setMessage("Please select a file first.");
      return;
    }

    const uploadedPath = await HandleUpload(file);

    if (uploadedPath) {
      setMessage("Upload successful!");
      setImageUrl(uploadedPath); // Show preview
    } else {
      setMessage("Upload failed.");
    }
  };

  return (
    <div>
      <input type="file" accept="image/*" onChange={handleFileChange} />
      <button onClick={handleUpload}>Upload</button>
      <p>{message}</p>
      {imageUrl && (
        <img
          src={imageUrl}
          alt="Uploaded preview"
          style={{ maxWidth: "200px", marginTop: "10px" }}
        />
      )}
    </div>
  );
}
