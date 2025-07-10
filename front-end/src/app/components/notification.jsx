"use client";
import { useEffect, useState } from "react";

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

export default function Notification({
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
            <div className="toast__icon" onClick={() => {
                setShow(false);
                if (onClose) onClose();
            }}
                aria-label="Close">{config.icon}</div>
            <div className="toast__content">
                <p className="toast__type">{config.title}</p>
                <p className="toast__message">{message}</p>
            </div>
        </div>
    );
}