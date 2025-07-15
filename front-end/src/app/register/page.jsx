"use client";

import { useState, useRef } from "react";
import Styles from "./register.module.css";
import { SendAuthData } from "../sendData.js";
import { useNotification } from "../context/notificationContext.jsx";
import { HandleUpload } from "../utils.jsx";
import { useRouter } from "next/navigation";
import { useWebSocket } from "../context/WebSocketContext.jsx";
import { useAuth } from "../context/AuthContext.jsx";

export default function Register() {
  const { connectWebSocket } = useWebSocket();
  const router = useRouter();
  const fileInputRef = useRef(null);
  const { isLoggedIn, setIsLoggedIn } = useAuth();
  const { showNotification } = useNotification();

  // keep track of the raw file so we can upload on submit
  const [selectedFile, setSelectedFile] = useState(null);
  const [avatarName, setAvatarName] = useState("");

  // all your other form fields
  const [formData, setFormData] = useState({
    fname: "",
    lname: "",
    username: "",
    email: "",
    password: "",
    birthdate: "",
    gender: "",
    about: "",
    // avatarPath will hold the *returned* URL after upload
    avatarPath: "",
  });

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData((fd) => ({ ...fd, [name]: value }));
  };

  const handleFileChange = (e) => {
    const file = e.target.files?.[0] ?? null;
    setSelectedFile(file);
    setAvatarName(file?.name || "");
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (isLoggedIn) return;

    // 1) upload avatar if one was selected
    let avatarPath = "";
    if (selectedFile) {
      avatarPath = await HandleUpload(selectedFile);
      console.log("k", avatarPath);
      if (!avatarPath) {
        return showNotification("Avatar upload failed", "error");
      }
    }

    // 2) build your payload
    const payload = {
      ...formData,
      avatar: avatarPath,
    };

    // 3) send it as JSON
    const response = await SendAuthData("/api/v1/auth/register", payload);

    if (response.status !== 200) {
      const errorBody = await response.json();
      showNotification(errorBody.error || "Registration failed", "error", 5000);
      return;
    }

    // 4) on success...
    setIsLoggedIn(true);
    showNotification("Registration successful!", "success", 5000);
    connectWebSocket();
    router.push("/");
  };

  return (
    <div className={Styles.container}>
      <div className={Styles.messageBox}>
        <h2>Join Our Social Network 👋</h2>
      </div>
      <div className={Styles.formContainer}>
        <form className={Styles.form} onSubmit={handleSubmit}>
          {/* First & Last Name */}
          <label className={Styles.label} htmlFor="fname">
            First Name
          </label>
          <input
            className={Styles.input}
            type="text"
            name="fname"
            id="fname"
            required
            onChange={handleChange}
          />

          <label className={Styles.label} htmlFor="lname">
            Last Name
          </label>
          <input
            className={Styles.input}
            type="text"
            name="lname"
            id="lname"
            required
            onChange={handleChange}
          />

          {/* Nickname, Email, Password, DOB, Gender, About */}
          <label className={Styles.label} htmlFor="username">
            Nickname
          </label>
          <input
            className={Styles.input}
            type="text"
            name="username"
            id="username"
            required
            onChange={handleChange}
          />

          <label className={Styles.label} htmlFor="email">
            Email
          </label>
          <input
            className={Styles.input}
            type="email"
            name="email"
            id="email"
            required
            onChange={handleChange}
          />

          <label className={Styles.label} htmlFor="password">
            Password
          </label>
          <input
            className={Styles.input}
            type="password"
            name="password"
            id="password"
            required
            onChange={handleChange}
          />

          <label className={Styles.label} htmlFor="birthdate">
            Date of Birth
          </label>
          <input
            className={Styles.input}
            type="date"
            name="birthdate"
            id="birthdate"
            required
            onChange={handleChange}
          />

          <label className={Styles.label}>Gender</label>
          <div>
            <label htmlFor="male">Male</label>
            <input
              type="radio"
              name="gender"
              id="male"
              value="male"
              required
              checked={formData.gender === "male"}
              onChange={handleChange}
            />

            <label htmlFor="female">Female</label>
            <input
              type="radio"
              name="gender"
              id="female"
              value="female"
              required
              checked={formData.gender === "female"}
              onChange={handleChange}
            />
          </div>

          <label className={Styles.label} htmlFor="about">
            About Me
          </label>
          <textarea
            className={Styles.input}
            name="about"
            id="about"
            rows={3}
            onChange={handleChange}
          />

          {/* Avatar Upload */}
          <button
            type="button"
            className={Styles.label}
            onClick={() => fileInputRef.current.click()}
          >
            Upload Avatar
          </button>
          <input
            ref={fileInputRef}
            type="file"
            name="avatar"
            accept="image/*"
            style={{ display: "none" }}
            onChange={handleFileChange}
          />
          {avatarName && <div className={Styles.fileName}>{avatarName}</div>}

          <button type="submit" className={Styles.submitButton}>
            Register
          </button>
        </form>
      </div>
    </div>
  );
}
