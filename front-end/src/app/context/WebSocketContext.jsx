"use client";
import { createContext, useContext, useEffect, useRef, useState } from "react";
import { useNotification } from "./notificationContext";

// import { Notification } from "../components/notification/notification.jsx";

const WebSocketContext = createContext(null);

export const WebSocketProvider = ({ children }) => {
  const ws = useRef(null);

  const reconnectTimeout = useRef(null);
  const reconnectAttempts = useRef(0);
  const MAX_RECONNECTS = 1;

  const [newMessage, setNewMessage] = useState(null);
  const [isConnected, setIsConnected] = useState(false);
  const [updateOnlineUser, setUpdateOnlineUser] = useState(null);
  const [newNotification, setNewNotification] = useState(null);
  const { showNotification } = useNotification();
  const target = useRef(null);

  const setTarget = (newTarget) => {
    target.current = newTarget;
  };

  const connectWebSocket = () => {
    if (ws.current) ws.current.close();

    ws.current = new WebSocket(process.env.NEXT_PUBLIC_API_URL + "/ws");

    ws.current.onopen = () => {
      console.log("✅ WebSocket connected");
      setIsConnected(true);
      reconnectAttempts.current = 0;
    };

    ws.current.onmessage = (event) => {
      const data = JSON.parse(event.data);

      console.log("the target is: ", target.current);
      if (data.content) {
        if ([data.sender, data.receiver].includes(target.current)) {
          data.sent_at = new Date().toISOString().replace(/\.\d{3}Z$/, "Z");
          setNewMessage(data);
        } else {
          if (data?.author_name != "system") {
            showNotification(`New messsage from ${data.author_name}`, "success");
            // showNotification(`New message from ${data.author_name}`, "success", true, 5000);
            console.warn("Message not for this chat:", data);
          } else {
            showNotification(`${data.content}`, "error")
          }
        }
      } else if (data.sender === "<system>", data.command) {
        if (data.command == "online") {
          console.log("data received via websocket: ", data);
          setUpdateOnlineUser(data)
        } else {
          setNewNotification(data)
          showNotification(data.content, "info")
        }
      }

    };

    ws.current.onclose = () => {
      console.warn("❌ WebSocket closed");
      setIsConnected(false);

      if (reconnectAttempts.current < MAX_RECONNECTS) {
        reconnectTimeout.current = setTimeout(() => {
          reconnectAttempts.current += 1;
          console.log(reconnectAttempts.current);
          connectWebSocket();
        }, 2000);
      } else {
        console.error("❌ Max reconnect attempts reached");
      }
    }


    ws.current.onerror = (err) => {
      console.error("⚠️ WebSocket error:", err);
      ws.current?.close();
    };
  };

  useEffect(() => {
    connectWebSocket();

    return () => {
      ws.current?.close();
      if (reconnectTimeout.current) {
        clearTimeout(reconnectTimeout.current);
      }
    };		// if _, e := helpers.Sockets[username.Username]; e {
    // 	username.Online = true
    // }
  }, []);

  const sendMessage = (msg) => {
    if (ws.current?.readyState === WebSocket.OPEN) {
      console.log("this message sent via websocket: ", msg);
      ws.current.send(JSON.stringify(msg));
    } else {
      console.warn("Can't send message: WebSocket not open");
    }
  };

  const closeWebSocket = () => {
    if (ws.current) {
      console.log("close websocket ❌");
      ws.current.close();
      if (reconnectTimeout.current) {
        clearTimeout(reconnectTimeout.current);
      }
    }
  };

  return (
    <WebSocketContext.Provider
      value={{
        sendMessage,
        closeWebSocket,
        connectWebSocket,
        newMessage,
        setTarget,
        isConnected,
        updateOnlineUser,
        newNotification,
      }}
    >
      {children}

    </WebSocketContext.Provider>
  );
};

export const useWebSocket = () => useContext(WebSocketContext);
