"use client";
import { createContext, useContext, useState, useCallback } from "react";
import Notification from "../utils.jsx"; // Adjust path

const NotificationContext = createContext();

export function NotificationProvider({ children }) {
  const [notification, setNotification] = useState({
    open: false,
    message: "",
    type: "success",
    duration: 3000,
  });

  const showNotification = useCallback((message, type = "success", duration = 3000) => {
    setNotification({ open: true, message, type, duration });
  }, []);

  const closeNotification = useCallback(() => {
    setNotification((n) => ({ ...n, open: false }));
  }, []);

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