"use client";
import { useEffect, useState } from "react";
import Image from "next/image";
import { SendData } from "./sendData.js";

// Notifications
let notificationCooldown = false;
function PlaySound(name) {
  if (typeof window === "undefined") return;

  const sounds = {
    alert: new Audio("/sounds/alert.mp3"),
    notification: new Audio("/sounds/notification.mp3"),
  };

  const sound = sounds[name];
  if (sound) {
    sound.currentTime = 0;
    sound.play().catch((e) => console.warn("Playback failed:", e));
  }
}

export const CapitalizeFirstLetter = (str) => {
  if (typeof str !== "string") return "";
  return str.charAt(0).toUpperCase() + str.slice(1);
};

export function showNotification(message, type = "success", sound = true, duration = 3000) {
  if (typeof window === "undefined") return; // SSR guard
  if (notificationCooldown) return;

  notificationCooldown = true;

  const toast = document.createElement("div");
  toast.className = `toast toast--${type}`;

  const toastConfig = {
    success: {
      title: "Success!",
      icon: `<svg class="toast__svg" viewBox="0 0 24 24">
               <path d="M9 16.17L4.83 12l-1.42 1.41L9 19 21 7l-1.41-1.41z"/>
             </svg>`,
    },
    error: {
      title: "Error!",
      icon: `<svg class="toast__svg" viewBox="0 0 24 24">
               <path d="M19 6.41L17.59 5 12 10.59 6.41 5 5 6.41 10.59 12 5 17.59 6.41 19 12 13.41 17.59 19 19 17.59 13.41 12z"/>
             </svg>`,
    },
  };

  const config = toastConfig[type] || toastConfig.success;

  toast.innerHTML = `
    <div class="toast__icon">${config.icon}</div>
    <div class="toast__content">
      <p class="toast__type">${config.title}</p>
      <p class="toast__message">${message}</p>
    </div>
  `;

  document.body.appendChild(toast);

  requestAnimationFrame(() => {
    toast.classList.add("show");
  });

  if (sound) {
    PlaySound(type === "success" ? "notification" : "alert");
  }

  const hideTimeout = setTimeout(() => {
    hideToast(toast);
  }, duration);

  toast._hideTimeout = hideTimeout;

  function hideToast(toastElement) {
    if (toastElement._hideTimeout) {
      clearTimeout(toastElement._hideTimeout);
    }

    toastElement.classList.remove("show");
    toastElement.style.top = "-200px";

    toastElement.addEventListener(
      "transitionend",
      () => {
        if (toastElement.parentNode) toastElement.remove();
      },
      { once: true }
    );

    setTimeout(() => {
      notificationCooldown = false;
    }, 500);
  }
}

const toastConfig = {
  success: {
    title: "Success!",
    icon: (
      <svg className="toast__svg" viewBox="0 0 24 24">
        <path d="M9 16.17L4.83 12l-1.42 1.41L9 19 21 7l-1.41-1.41z" />
      </svg>
    ),
  },
  error: {
    title: "Error!",
    icon: (
      <svg className="toast__svg" viewBox="0 0 24 24">
        <path d="M19 6.41L17.59 5 12 10.59 6.41 5 5 6.41 10.59 12 5 17.59 6.41 19 12 13.41 17.59 19 19 17.59 13.41 12z" />
      </svg>
    ),
  },
  info: {
    title: "Info",
    icon: (
      <svg className="toast__svg" viewBox="0 0 24 24">
        <circle cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="2" fill="none" />
        <line x1="12" y1="8" x2="12" y2="12" stroke="currentColor" strokeWidth="2" />
        <circle cx="12" cy="16" r="1" fill="currentColor" />
      </svg>
    ),
  },
};

function playSound(name) {
  if (typeof window === "undefined") return;
  const sounds = {
    alert: new Audio("/sounds/alert.mp3"),
    notification: new Audio("/sounds/notification.mp3"),
  };
  const sound = sounds[name];
  if (sound) {
    sound.currentTime = 0;
    sound.play().catch((e) => console.warn("Playback failed:", e));
  }
}

export function Notification({
  open,
  message,
  type = "success",
  duration = 3000,
  sound = true,
  onClose,
}) {
  const [show, setShow] = useState(open);

  useEffect(() => {
    setShow(open);
    if (open && sound) {
      playSound(type === "success" ? "notification" : "alert");
    }
    if (open && duration > 0) {
      const timer = setTimeout(() => {
        setShow(false);
        if (onClose) onClose();
      }, duration);
      return () => clearTimeout(timer);
    }
  }, [open, duration, sound, type, onClose]);

  if (!show) return null;

  const config = toastConfig[type] || toastConfig.success;

  return (
    <div className={`toast toast--${type} show`}>
      <div
        className="toast__icon"
        onClick={() => {
          setShow(false);
          if (onClose) onClose();
        }}
        aria-label="Close"
      >
        {config.icon}
      </div>
      <div className="toast__content">
        <p className="toast__type">{config.title}</p>
        <p className="toast__message">{message}</p>
      </div>
    </div>
  );
}

// Password toggle
export function usePasswordToggle() {
  useEffect(() => {
    const handleClick = (event) => {
      const btn = event.target.closest(".togglePwd");
      if (!btn) return;

      const input = btn.previousElementSibling;
      const icon = btn.querySelector(".vis_icon");

      if (input && icon) {
        if (input.type === "password") {
          input.type = "text";
          icon.innerText = "visibility_off";
        } else {
          input.type = "password";
          icon.innerText = "visibility";
        }
      }
    };

    document.addEventListener("click", handleClick);
    return () => document.removeEventListener("click", handleClick);
  }, []);
}

// Like/Dislike button
export default function LikeDeslike({ EntityID, EntityType, isLiked, currentLikeCount }) {
  const [loading, setLoading] = useState(false);
  const [liked, setLiked] = useState(isLiked);
  const [likeCount, setLikeCount] = useState(currentLikeCount);
  if (!EntityID || !EntityType) return null;
  const handleLikeDeslike = async () => {
    if (loading) return;

    setLoading(true);

    const likeInfo = {
      entity_id: EntityID,
      entity_type: EntityType,
      is_liked: liked,
    };

    try {
      setLiked(!liked);
      setLikeCount(liked ? likeCount - 1 : likeCount + 1);

      const response = await SendData("/api/v1/set/like", likeInfo);
      const body = await response.json();

      if (response.status === 200) {
        console.log("Like/Dislike processed successfully!");
      } else {
        console.log(body);
        setLiked(liked);
        setLikeCount(liked ? likeCount - 1 : likeCount + 1);
      }
    } catch (error) {
      console.error("Error while processing like/dislike:", error);
      setLiked(liked);
      setLikeCount(liked ? likeCount - 1 : likeCount + 1);
    }

    setLoading(false);
  };

  return (
    <div
      onClick={handleLikeDeslike}
      style={{ cursor: "pointer", marginRight: "10px", display: "flex", alignItems: "center" }}
    >
      <Image
        src={liked ? "/Like.svg" : "/Like2.svg"}
        alt={liked ? "liked" : "like"}
        width={20}
        height={20}
      />
      <p style={{ marginLeft: "5px" }}>{likeCount}</p>
    </div>
  );
}

// upload
export async function HandleUpload(image) {
  console.log("iamge : ", image);

  if (!image) return null;

  const formData = new FormData();
  formData.append("file", image);

  const response = await SendData("/api/v1/upload", formData);

  if (!response.ok) {
    console.error("Image upload failed");
    return null;
  }

  const { path } = await response.json();
  return path;
}


export function TimeAgo(timestamp) {
  const now = new Date();
  const past = new Date(timestamp.replace(" ", "T")); // Parse properly
  const diffMs = now - past;

  const seconds = Math.floor(diffMs / 1000);
  const minutes = Math.floor(seconds / 60);
  const hours = Math.floor(minutes / 60);
  const days = Math.floor(hours / 24);
  const weeks = Math.floor(days / 7);
  const months = Math.floor(days / 30);
  const years = Math.floor(days / 365);

  if (seconds < 60) return "Just now";
  if (minutes < 60) return `${minutes} minute${minutes !== 1 ? "s" : ""} ago`;
  if (hours < 24) return `${hours} hour${hours !== 1 ? "s" : ""} ago`;
  if (days < 7) return `${days} day${days !== 1 ? "s" : ""} ago`;
  if (weeks < 5) return `${weeks} week${weeks !== 1 ? "s" : ""} ago`;
  if (months < 12) return `${months} month${months !== 1 ? "s" : ""} ago`;
  return `${years} year${years !== 1 ? "s" : ""} ago`;
}
