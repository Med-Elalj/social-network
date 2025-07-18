"use client";
import { createContext, useContext, useEffect, useRef, useState } from "react";
import { useNotification } from "./NotificationContext.jsx";

const WebSocketContext = createContext(null);
const WEBSOCKET_URL = "http://localhost:8080/api/v1/ws";

export const WebSocketProvider = ({ children }) => {
  const ws = useRef(null);

  const reconnectTimeout = useRef(null);
  const reconnectAttempts = useRef(0);
  const MAX_RECONNECTS = 1;

  const [newMessage, setNewMessage] = useState(null);
  const [isConnected, setIsConnected] = useState(false);
  const [updateOnlineUser, setUpdateOnlineUser] = useState(null);
  const [newNotification, setNewNotification] = useState(null);
  const [newFollowRequest, setNewFollowRequest] = useState(null);
  const [newGroupRequest, setNewGroupRequest] = useState(null);
  const { showNotification } = useNotification();
  const target = useRef(null);

  const setTarget = (newTarget) => {
    target.current = newTarget;
  };

  const connectWebSocket = () => {
    if (ws.current) ws.current.close();

    ws.current = new WebSocket(WEBSOCKET_URL);

    ws.current.onopen = () => {
      console.log("✅ WebSocket connected");
      setIsConnected(true);
      reconnectAttempts.current = 0;
    };

    ws.current.onmessage = (event) => {
      const data = JSON.parse(event.data);

      if (data.content) {
        if ([data.sender, data.receiver].includes(target.current)) {
          data.sent_at = new Date().toISOString().replace(/\.\d{3}Z$/, "Z");
          setNewMessage(data);
        } else {
          if (data?.author_name != "system") {
            showNotification(
              `New messsage from ${data.author_name}`,
              "success"
            );
            console.warn("Message not for this chat:", data);
          } else {
            showNotification(`${data.content}`, "error");
          }
        }
      } else if ((data.sender === "<system>", data.command)) {
        if (data.command == "online") {
          setUpdateOnlineUser(data);
        } else {
          showNotification(data.value.message, "info");
          // setNewNotification({ ...data.value, type: data.command });
          if (data?.command === "group_request") {
            console.log("new group request received: ", data.value);
            setNewGroupRequest(data.value);
          } else if (data?.command === "follow_request") {
            console.log("Received newFollowRequest:", data.value);
            setNewFollowRequest(data.value);
          }
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
        console.log("❌ Max reconnect attempts reached");
      }
    };

    ws.current.onerror = (err) => {
      console.log("⚠️ WebSocket error:", err);
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
    }; // if _, e := helpers.Sockets[username.Username]; e {
    // 	username.Online = true
    // }
  }, []);

  const sendMessage = (msg) => {
    if (ws.current?.readyState === WebSocket.OPEN) {
      console.log("to send", msg);
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
        newFollowRequest,
        newGroupRequest,
      }}
    >
      {children}
    </WebSocketContext.Provider>
  );
};

export const useWebSocket = () => useContext(WebSocketContext);
