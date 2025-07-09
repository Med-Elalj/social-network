"use client"; // REQUIRED if you're using App Router (in `/app` folder)

import { useState } from "react";
import Style from "../profile.module.css"; // Adjust the path as necessary
import { useNotification } from "../../../context/notificationContext.jsx";
import { SendData } from "../../../../../utils/sendData.js";
import { Router } from "next/dist/client/router.js";

export default function Settings() {
  const [activeForm, setActiveForm] = useState(null); // 'nickname', 'password', or 'delete'
  const [formData, setFormData] = useState({
    nickname: "",
    email: "",
    currentPassword: "",
    newPassword: "",
    confirmPassword: "",
    deletePassword: "",
  });
  const [confirmationWord, setConfirmationWord] = useState("");
  const [userTypedWord, setUserTypedWord] = useState("");
  const {showNotification} = useNotification();

  const handleChange = (e) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value,
    });
  };

  const handleNicknameUpdate = () => {
    setActiveForm(activeForm === "nickname" ? null : "nickname");
  };

  const handleNicknameSubmit = () => {
    if (!formData.nickname.trim()) {
      showNotification("Please enter a nickname.", "error");
      return;
    }

    SendData("/api/v1/settings/updateUsername", {
      nickname: formData.nickname,
    });
    setActiveForm(null);
    setFormData({ ...formData, nickname: "" }); // Clear nickname field
  };

  const handleChangePassword = () => {
    setActiveForm(activeForm === "password" ? null : "password");
  };

  const handlePasswordSubmit = () => {
    if (
      !formData.currentPassword ||
      !formData.newPassword ||
      !formData.confirmPassword
    ) {
      showNotification("Please fill in all password fields.", "error");
      return;
    }

    if (formData.newPassword !== formData.confirmPassword) {
      showNotification("New password and confirmation do not match.", "error");
      return;
    }

    if (formData.newPassword.length < 8) {
      showNotification(
        "New password must be at least 8 characters long.",
        "error"
      );
      return;
    }

    SendData("/api/v1/settings/updatePassword", {
      currentPassword: formData.currentPassword,
      newPassword: formData.newPassword,
    });

    // Clear password fields after submission
    setFormData({
      ...formData,
      currentPassword: "",
      newPassword: "",
      confirmPassword: "",
    });
    setActiveForm(null);
  };

  const handleDeleteProfile = () => {
    if (activeForm === "delete") {
      setActiveForm(null);
      return;
    }

    // Generate random string user must type (e.g. 6-8 chars)
    const words = [
      "removeMe",
      "delete123",
      "finalExit",
      "lastStep",
      "goodbyeNow",
    ];
    const random = words[Math.floor(Math.random() * words.length)];
    setConfirmationWord(random);
    setUserTypedWord("");
    setActiveForm("delete");
  };

  const handleDeleteConfirm = async () => {
    if (!formData.deletePassword.trim()) {
      showNotification(
        "Please enter your password to confirm deletion.",
        "error"
      );
      return;
    }
    // Check if user typed the confirmation string correctly
    if (userTypedWord !== confirmationWord) {
      showNotification("Confirmation text does not match.", "error");
      return;
    }

    const response = await SendData("/api/v1/settings/delete", {
      confirmDelete: true,
      deletePassword: formData.deletePassword,
    });

    if (response.ok) {
      fetch("/api/v1/auth/logout", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
      }).then(() => {
        // Redirect to homepage or login page after logout
        Router.push("/");
      });
      showNotification(
        "Your account has been successfully deleted.",
        "success"
      );

      // Optionally redirect after short delay (e.g., logout or homepage)
      setTimeout(() => {
        window.location.href = "/goodbye"; // Or login/homepage
      }, 2000);
    } else {
      showNotification(
        "Failed to delete account. Please check your password.",
        "error"
      );
    }
  };

  const handleCancel = () => {
    setActiveForm(null);
    // Clear all form data when canceling
    setFormData({
      nickname: "",
      currentPassword: "",
      newPassword: "",
      confirmPassword: "",
    });
  };

  return (
    <div className={Style.post} style={{ width: "80%", margin: "20px auto" }}>
      <h2 className={Style.title}>Profile Settings</h2>
      <section
        className={Style.footer}
        style={{ display: "flex", gap: "10px", flexWrap: "wrap" }}
      >
        <button
          type="button"
          className={Style.button}
          onClick={handleNicknameUpdate}
          style={
            activeForm === "nickname"
              ? { backgroundColor: "#007bff", color: "white" }
              : {}
          }
        >
          Update Nickname
        </button>
        <button
          type="button"
          className={Style.button}
          onClick={handleChangePassword}
          style={
            activeForm === "password"
              ? { backgroundColor: "#007bff", color: "white" }
              : {}
          }
        >
          Change Password
        </button>
        <button
          type="button"
          className={Style.button}
          style={{
            backgroundColor: activeForm === "delete" ? "#dc3545" : "red",
            color: "white",
            border: "1px solid black",
          }}
          onClick={handleDeleteProfile}
        >
          Delete Profile
        </button>
      </section>

      {activeForm === "nickname" && (
        <div style={{ marginTop: "20px" }}>
          <h4>Update Nickname</h4>

          <div className={Style.inputWrapper}>
            <input
              className={Style.input}
              type="text"
              name="nickname"
              placeholder="Enter new nickname"
              value={formData.nickname}
              onChange={handleChange}
              required
            />
          </div>

          <div style={{ marginTop: "10px", display: "flex", gap: "10px" }}>
            <button
              type="button"
              className={Style.button}
              onClick={handleNicknameSubmit}
            >
              Update Nickname
            </button>
            <button
              type="button"
              className={Style.button}
              style={{ backgroundColor: "gray" }}
              onClick={handleCancel}
            >
              Cancel
            </button>
          </div>
        </div>
      )}

      {activeForm === "password" && (
        <div style={{ marginTop: "20px" }}>
          <h4>Change Password</h4>

          <div className={Style.inputWrapper}>
            <input
              className={Style.input}
              type="password"
              name="currentPassword"
              placeholder="Current Password"
              value={formData.currentPassword}
              onChange={handleChange}
              required
            />
          </div>

          <div className={Style.inputWrapper}>
            <input
              className={Style.input}
              type="password"
              name="newPassword"
              placeholder="New Password"
              value={formData.newPassword}
              onChange={handleChange}
              required
            />
          </div>

          <div className={Style.inputWrapper}>
            <input
              className={Style.input}
              type="password"
              name="confirmPassword"
              placeholder="Confirm New Password"
              value={formData.confirmPassword}
              onChange={handleChange}
              required
            />
          </div>

          <div style={{ marginTop: "10px", display: "flex", gap: "10px" }}>
            <button
              type="button"
              className={Style.button}
              onClick={handlePasswordSubmit}
            >
              Update Password
            </button>
            <button
              type="button"
              className={Style.button}
              style={{ backgroundColor: "gray" }}
              onClick={handleCancel}
            >
              Cancel
            </button>
          </div>
        </div>
      )}

      {activeForm === "delete" && (
        <div style={{ marginTop: "20px" }}>
          <h4 style={{ color: "red" }}>Delete Profile</h4>
          <div>
            <div className={Style.inputWrapper}>
              <input
                className={Style.input}
                type="password"
                name="deletePassword" // ✅ fixed casing
                placeholder="Enter password to confirm"
                value={formData.deletePassword} // ✅ match state
                onChange={handleChange}
                required
              />
            </div>
            <p style={{ fontSize: "14px", color: "red" }}>
              Type the following to confirm:
            </p>
            <p
              style={{
                fontFamily: "monospace",
                userSelect: "none",
                pointerEvents: "none",
                background: "#f5f5f5",
                padding: "4px 8px",
                display: "inline-block",
              }}
            >
              {confirmationWord}
            </p>
            <div className={Style.inputWrapper}>
              <input
                className={Style.input}
                type="text"
                name="userTypedWord"
                placeholder="Re-type confirmation text"
                value={userTypedWord}
                onChange={(e) => setUserTypedWord(e.target.value)}
                required
              />
            </div>

            <p style={{ color: "#666", marginBottom: "15px" }}>
              ⚠️ This action cannot be undone. Your profile and all associated
              data will be permanently deleted.
            </p>

            <button
              type="button"
              className={Style.button}
              style={{
                backgroundColor: "#dc3545",
                color: "white",
                border: "1px solid #dc3545",
              }}
              onClick={handleDeleteConfirm}
            >
              Yes, Delete My Profile
            </button>
            <button
              type="button"
              className={Style.button}
              style={{ backgroundColor: "gray" }}
              onClick={handleCancel}
            >
              Cancel
            </button>
          </div>
        </div>
      )}
    </div>
  );
}
