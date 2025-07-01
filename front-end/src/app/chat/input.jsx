'use client';

import React, { useState, useRef } from "react";
import Styles from "./chat.module.css";
const WEBSOCKET_URL = '/api/v1/ws';

export default function ChatInput({ addMessage, target }) {
    const [message, setMessage] = useState("");
    const ws = useRef(null);
    const reconnectTimeout = useRef(null);

    React.useEffect(() => {
        let isMounted = true;

        function connect() {
            ws.current = new WebSocket(WEBSOCKET_URL);

            ws.current.onopen = () => {
                console.log("WebSocket connected");
            };

            ws.current.onclose = () => {
                console.log("WebSocket disconnected, retrying in 2s...");
                if (isMounted) {
                    reconnectTimeout.current = setTimeout(connect, 2000);
                }
            };

            ws.current.onmessage = (event) => {
                const data = JSON.parse(event.data);
                console.log("Received:", data);
                if ([data.sender, data.receiver].includes(target[0]) || data.author_name === 'system') {
                    data.sent_at = data.author_name == 'system' ? new Date().toISOString() : new Date().toISOString().replace(/\.\d{3}Z$/, 'Z');
                    addMessage(data);
                } else {
                    console.warn("Message not for this user:", data); // TODO NOTIFY
                }
            };
        }

        connect();

        return () => {
            isMounted = false;
            if (ws.current) ws.current.close();
            if (reconnectTimeout.current) clearTimeout(reconnectTimeout.current);
        };
    }, [target]);

    const sendMessage = (e) => {
        e.preventDefault();
        console.log("here")
        if (!message.trim() || !ws.current || ws.current.readyState !== 1) return;
        const dm = {
            receiver: target[0],
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