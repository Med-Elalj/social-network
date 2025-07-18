"use client";

import { useState, useRef } from "react";
import Styles from "./register.module.css";
import { SendAuthData } from "../sendData.js";
import { useNotification } from "../context/NotificationContext.jsx";
import { usePasswordToggle } from "../utils.jsx";
import { useRouter } from "next/navigation";
import { useWebSocket } from "../context/WebSocketContext.jsx";
import { useAuth } from "../context/AuthContext.jsx";
import { HandleUpload } from "@/app/components/upload.jsx";

export default function Register() {
  const { connectWebSocket } = useWebSocket();
  const Router = useRouter();
  const fileInputRef = useRef(null);
  const { isLoggedIn, setIsLoggedIn } = useAuth();
  const [formData, setFormData] = useState({
    username: "",
    email: "",
    password: "",
    gender: "",
    fname: "",
    lname: "",
    birthdate: "",
    avatar: "",
    about: "",
  });
  const { showNotification } = useNotification();
  const [avatarName, setAvatarName] = useState("");
  usePasswordToggle();

  const handleChange = (e) => {
    const { name, value, files } = e.target;
    if (name === "avatar") {
      const file = files[0];
      setAvatarName(file.name);
      if (file) {
        setAvatarName(file.name);
        //rm the upload and get back a public path
        HandleUpload(file)
          .then((path) => {
            setFormData((f) => ({ ...f, avatar: path }));
          })
          .catch((err) => {
            console.error("Avatar upload failed", err);
          });
      }
    } else {
      setFormData((prev) => ({
        ...prev,
        [name]: value,
      }));
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (isLoggedIn) {
      return;
    }
    const data = new FormData();
    for (const key in formData) {
      data.append(key, formData[key]);
    }

    const response = await SendAuthData("/api/v1/auth/register", formData);

    if (response.status !== 200) {
      const errorBody = await response.json();
      if (formData.avatar) {
        localStorage.setItem("avatar", formData.avatar);
      }
      showNotification(errorBody.error || "Registration failed. Please try again.", "error", 5000);
    } else {
      setIsLoggedIn(true);
      showNotification("Registration successful! Welcome to our social network!", "success", 5000);
      Router.push("/");
      connectWebSocket();
    }
  };

  return (
    <div className={Styles.container}>
      <div className={Styles.messageBox}>
        <h2>Join Our Social Network 👋</h2>
      </div>
      <div className={Styles.formContainer}>
        <form className={Styles.form} onSubmit={handleSubmit}>
          <label className={Styles.label} htmlFor="fname">
            First Name
          </label>
          <input
            className={Styles.input}
            type="text"
            name="fname"
            id="firstName"
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
            id="lastName"
            required
            onChange={handleChange}
          />

          <label className={Styles.label} htmlFor="username">
            Nickname
          </label>
          <input
            className={Styles.input}
            type="text"
            name="username"
            id="nickName"
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

          {/* <input className={styles.input} type="password" name="password" id="password" onChange={handleChange} /> */}
          <label className={Styles.label} htmlFor="password">
            Password
          </label>
          <div className={Styles.inputWrapper}>
            <input
              className={Styles.input}
              type="password"
              name="password"
              id="password"
              onChange={handleChange}
              required
            />
            <i className="togglePwd">
              <span className="icon vis_icon material-symbols-outlined">visibility</span>
            </i>
          </div>

          <label className={Styles.label} htmlFor="birthdate">
            Date of Birth
          </label>
          <input
            className={Styles.input}
            type="date"
            name="birthdate"
            id="dob"
            required
            onChange={handleChange}
          />

          <label className={Styles.label} htmlFor="gender" required>
            Gender
          </label>
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

          <button
            type="button"
            className={Styles.label}
            onClick={() => fileInputRef.current.click()}
            style={{ cursor: "pointer", display: "flex", alignItems: "center",justifyContent: "center" }}
          >
            <img src="/Image.svg" alt="Upload" width="24" height="24" style={{ marginTop: "8px" }}/> Upload Avatar
          </button>
          <input
            ref={fileInputRef}
            id="profileImg"
            type="file"
            name="avatar"
            style={{ display: "none" }}
            accept="image/*"
            onChange={handleChange}
          />
          {avatarName && <div className={Styles.fileName}>{avatarName}</div>}

          <label className={Styles.label} htmlFor="about">
            About Me
          </label>
          <input
            className={Styles.input}
            type="text"
            name="about"
            id="about"
            onChange={handleChange}
          />

          <button className={Styles.submitButton} type="submit">
            Register
          </button>
        </form>
      </div>
    </div>
  );
}
