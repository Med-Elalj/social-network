import React, { useState, useRef } from "react";
import Styles from "./style.module.css";
const WEBSOCKET_URL = process.env.NEXT_PUBLIC_BACKEND_URL + '/api/v1/ws';

export default function ChatInput({ userId, recipientId }) {
    const [message, setMessage] = useState("");
    const ws = useRef(null);

    // Connect WebSocket on mount
    React.useEffect(() => {
        ws.current = new WebSocket(WEBSOCKET_URL);

        ws.current.onopen = () => {
            console.log("WebSocket connected");
        };

        ws.current.onclose = () => {
            console.log("WebSocket disconnected");
        };

        return () => {
            ws.current.close();
        };
    }, []);

    const sendMessage = (e) => {
        e.preventDefault();
        if (!message.trim() || ws.current.readyState !== 1) return;

        const dm = {
            type: "dm",
            from: userId,
            to: recipientId,
            content: message,
        };

        ws.current.send(JSON.stringify(dm));
        setMessage("");
    };

    return (
        <form onSubmit={sendMessage} className={Styles.chat_input_form} >
            <input
                type="text"
                value={message}
                placeholder="Type a message..."
                onChange={(e) => setMessage(e.target.value)}
            />
            <button type="submit" disabled={!message.trim()}>
                Send
            </button>
        </form>
    );
}