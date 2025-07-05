"use client";

import { useState } from "react";
import Styles from "./chat.module.css";
import Image from "next/image"; 
import { useWebSocket } from "../context/WebSocketContext";


export default function ChatInput({ target }) {
  const [message, setMessage] = useState("");
  const {sendMessage}=useWebSocket();

  const sendMessageHandler = (e) => {
    e.preventDefault(); 
    if (!message.trim()) return;

    const dm = {
      receiver: target,
      content: message,
    };

    sendMessage(dm)
    setMessage("");
  };

  return (
    <form onSubmit={sendMessageHandler} className={Styles.chat_input_form}>
      <div className={Styles.uploadContainer}>
        <label htmlFor="media" style={{ cursor: "pointer" }}>
          <Image src="/upload.svg" width={25} height={25} alt="upload" />
        </label>
        <input
          type="file"
          name="media"
          id="media"
          style={{ display: "none" }}
          accept="image/*,video/*"
          // onChange={handleMediaChange} // Optional
        />
      </div>

      <input
        type="text"
        id="message"
        value={message}
        placeholder="Type a message..."
        onChange={(e) => setMessage(e.target.value)}
        className={Styles.messageInput}
      />

      <button type="submit" disabled={!message.trim()} className={Styles.sendBtn}>
        <Image
          src="/send.svg"
          width={25}
          height={25}
          alt="send"
          style={{ cursor: "pointer", marginRight: "6%" }}
        />
      </button>
    </form>
  );
}
