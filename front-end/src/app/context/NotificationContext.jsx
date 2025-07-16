"use client";
import { createContext, useContext, useState, useCallback, useEffect } from "react";
import Notification from "../components/notification.jsx";

const NotificationContext = createContext();

// Global reference for use outside React
let externalNotify = null;

export function externalNotification(message, type = "success", duration = 3000) {
  if (typeof window === "undefined") return;

  // 🔒 Ensure message is string (handle object input)
  if (typeof message === "object") {
    try {
      message = message.message || JSON.stringify(message);
    } catch (e) {
      message = "Unknown error";
    }
  }

  if (externalNotify) {
    externalNotify(message, type, duration);
  }
}


export function NotificationProvider({ children }) {
  const [notification, setNotification] = useState({
    open: false,
    message: "",
    type: "success",
    duration: 3000,
  });

const showNotification = useCallback((message, type = "success", duration = 3000) => {
  if (Array.isArray(message)) {
    message = message.join(", ");
  } else if (typeof message === "object" && message !== null) {
    try {
      message = message.message || JSON.stringify(message);
    } catch {
      message = "Invalid message";
    }
  }

  setNotification({ open: true, message, type, duration });
}, []);


  const closeNotification = useCallback(() => {
    setNotification((n) => ({ ...n, open: false }));
  }, []);

  useEffect(() => {
    externalNotify = showNotification;
  }, [showNotification]);

  return (
    <NotificationContext.Provider value={{ showNotification }}>
      {children}
      <Notification
        open={notification.open}
        message={notification.message}
        type={notification.type}
        duration={notification.duration}
        onClose={closeNotification}
      />
    </NotificationContext.Provider>
  );
}

export function useNotification() {
  return useContext(NotificationContext);
}
