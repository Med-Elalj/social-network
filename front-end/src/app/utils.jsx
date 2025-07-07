"use client";
import { useEffect, useState } from "react";
import Image from "next/image";
import { SendData } from "../../utils/sendData.js";

// Notifications
let notificationCooldown = false;
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
//todo: reduce body opacity a little when using notification

export const CapitalizeFirstLetter = (str) => {
  if (typeof str !== "string") return "";
  return str.charAt(0).toUpperCase() + str.slice(1);
}

export function showNotification(
  message,
  type = "success",
  sound = true,
  duration = 3000
) {
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
    playSound(type === "success" ? "notification" : "alert");
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
export default function LikeDeslike({ EntityID, EntityType,isLiked, currentLikeCount }) {
  const [loading, setLoading] = useState(false);
  const [liked, setLiked] = useState(isLiked);
  const [likeCount, setLikeCount] = useState(currentLikeCount);

  const handleLikeDeslike = async () => {
    if (loading) return;

    setLoading(true);

    const likeInfo = {
      entity_id: EntityID,
      entity_type: EntityType,
      is_liked: liked
    };

    try {
      setLiked(!liked);
      setLikeCount(liked ? likeCount - 1 : likeCount + 1);

      const response = await SendData("/api/v1/set/like", likeInfo);
      const body = await response.json();

      if (response.status === 200) {
        console.log('Like/Dislike processed successfully!');
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
    <div onClick={handleLikeDeslike} style={{ cursor: "pointer",marginRight:"10px",display:"flex",alignItems:"center" }}>
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