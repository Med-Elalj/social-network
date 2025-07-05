"use client";
import { createContext, useContext, useEffect, useRef, useState } from "react";

const WebSocketContext = createContext(null);
const WEBSOCKET_URL = "/api/v1/ws";

export const WebSocketProvider = ({ children }) => {
  const ws = useRef(null);

  const reconnectTimeout = useRef(null);
  const reconnectAttempts = useRef(0);
  const MAX_RECONNECTS = 5;

  const [newMessage, setNewMessage] = useState(null);
  const [isConnected, setIsConnected] = useState(false);
  const [target, setTarget] = useState(null);

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
      console.log("data received via websocket: ", data);
      if ([data.sender, data.receiver].includes(target)) {
        data.sent_at = new Date().toISOString().replace(/\.\d{3}Z$/, "Z");
        setNewMessage(data);
      } else {
        console.warn("Message not for this chat:", data); // TODO: notify user
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
    };

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
    };
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
      }}
    >
      {children}
      
    </WebSocketContext.Provider>
  );
};

export const useWebSocket = () => useContext(WebSocketContext);
