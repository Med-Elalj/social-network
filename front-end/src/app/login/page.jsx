"use client";

import { usePasswordToggle } from "../utils";

import { useState } from "react";
import Styles from "./login.module.css";
import { SendAuthData } from "../sendData.js";
import { useRouter } from "next/navigation";
import { useWebSocket } from "../context/WebSocketContext.jsx";
import { useNotification } from "../context/NotificationContext.jsx";
import { useAuth } from "../context/AuthContext.jsx";

export default function Login() {
  usePasswordToggle();
  const { connectWebSocket, isConnected } = useWebSocket();
  const { showNotification } = useNotification();
  const router = useRouter();
  const { isLoggedIn, setIsLoggedIn } = useAuth();
  const [formData, setFormData] = useState({
    login: "",
    pwd: "",
  });

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData({ ...formData, [name]: value });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    if (isLoggedIn) {
      router.push("/");
    }
    const response = await SendAuthData("/api/v1/auth/login", formData);

    if (response.ok) {
      const res = await response.json();
      if (!isConnected) connectWebSocket();
      setIsLoggedIn(true);
      router.push("/?login=success");

      if (res.avatar) {
        localStorage.setItem("avatar", res.avatar);
      } else if (!localStorage.getItem("avatar") && res.avatar) {
        localStorage.removeItem("avatar");
      }

      // router.push("/");
      showNotification("Login successful!", "success", 5000);
    } else {
      let errorBody;
      try {
        errorBody = await response.json();
      } catch {
        errorBody = { error: "Unknown error" };
      }
      showNotification(
        errorBody.error || "Login failed. Please try again.",
        "error",
        5000
      );
    }
  };

  return (
    <>
      <div className={Styles.container}>
        <div className={Styles.messageBox}>
          <h2>Welcome back to Our Social Network 👋</h2>
        </div>

        <form className={Styles.form} onSubmit={handleSubmit}>
          <label className={Styles.label} htmlFor="login">
            Email or Nickname
          </label>
          <input
            className={Styles.input}
            type="text"
            name="login"
            id="login"
            onChange={handleChange}
            required
          />

          <label className={Styles.label} htmlFor="pwd">
            Password
          </label>
          <div className={Styles.inputWrapper}>
            <input
              className={Styles.input}
              type="password"
              name="pwd"
              id="password"
              onChange={handleChange}
              required
            />
            <i className="togglePwd">
              <span className="icon vis_icon material-symbols-outlined">
                visibility
              </span>
            </i>
          </div>

          <button className={Styles.button} type="submit">
            Login
          </button>
        </form>
      </div>
    </>
  );
}
