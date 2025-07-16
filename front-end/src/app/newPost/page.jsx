"use client";

import { useState, useRef, useEffect, use } from "react";
import Styles from "./newPost.module.css";
import { GetData, SendData } from "../sendData.js";
import { useRouter } from "next/navigation";
import { HandleUpload } from "../utils.jsx";

export default function NewPost() {
  const [content, setContent] = useState("");
  const [image, setImage] = useState(null);
  const [privacy, setPrivacy] = useState("public");
  const [previewUrl, setPreviewUrl] = useState(null);
  const [selectedFriends, setSelectedFriends] = useState([]);
  const [showDropdown, setShowDropdown] = useState(false);
  const [users, setUsers] = useState(null);
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
      content,
      privacy,
      image: imagePath,
      groupId: null,
      privetids: privacy === "private" ? selectedFriends : [],
    };

    const response = await SendData("/api/v1/set/Post", formData);

    if (response.status !== 200) {
      const errorBody = await response.json();
      console.log(errorBody);
    } else {
      router.push("/");
    }
  };


  useEffect(() => {
    const fetchUsers = async () => {
      try {
        const res = await GetData(`/api/v1/get/myFollowers`);
        if (res.ok) {
          const data = await res.json();
          setUsers(data);
        } else {
          console.error("Failed to fetch followers");
        }
      } catch (err) {
        console.error("Error fetching followers:", err);
      }
    };

    fetchUsers();
  }, []);

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
        <div className={Styles.privacy}>
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
              <option value="almost_private">almost private</option>
              <option value="private">private</option>
            </select>
          </div>

          {privacy === "private" && (
            <div className={Styles.friendsSection}>
              <label
                className={Styles.dropdownToggle}
                onClick={() => setShowDropdown((prev) => !prev)}
              >
                Select Friends ▾
              </label>


              {showDropdown && (
                <div className={Styles.friendList}>
                  {users ? users?.map((friend) => {
                    const isSelected = selectedFriends.includes(friend);

                    return (
                      <div
                        key={friend.id}
                        className={`${Styles.friendItem} ${isSelected ? Styles.selected : ""}`}
                        onClick={() =>
                          setSelectedFriends((prev) => {
                            return isSelected
                              ? prev.filter((f) => f !== friend)
                              : [...prev, friend]
                          }
                          )
                        }
                      >
                        {friend.name}
                      </div>
                    );
                  }) : "No Friends"}
                </div>

              )}
            </div>
          )}
        </div>

        {/* Show Selected Friends */}
        {privacy === "private" && selectedFriends.length > 0 && (
          <div className={Styles.selectedFriends}>
            <p>Selected Friends:</p>
            <ul>
              {selectedFriends.map((friend) => (
                <li key={friend.id}>{friend.name}</li>
              ))}
            </ul>
          </div>
        )}

        <button type="submit">Post</button>
      </form>
    </div>
  );
}
