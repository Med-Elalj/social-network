"use client";

import { useWebSocket } from "../context/WebSocketContext";
import Style from "./chat.module.css";
import Time from "./time";
import { useEffect, useRef, useState } from "react";

const me = 0;

export async function getMessages(person_name, page) {
  const response = await fetch(`${NEXT_PUBLIC_API_URL}/get/dmhistory`, {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
      target: person_name,
      page: page,
    },
  });

  if (!response.ok) {
    throw new Error("Failed to fetch messages");
  }

  return await response.json();
}

export default function Messages({ user }) {
  const [messages, setMessages] = useState([]);
  const [hasMore, setHasMore] = useState(true);
  const [loading, setLoading] = useState(false);
  const { newMessage } = useWebSocket();
  const [page, setPage] = useState(0);
  const containerRef = useRef(null);
  const prevScrollHeight = useRef(0);

  useEffect(() => {
    setMessages([...messages, newMessage]);
  }, [newMessage]);

  useEffect(() => {
    if (
      containerRef.current &&
      containerRef.current.scrollHeight > prevScrollHeight.current
    ) {
      const newHeight = containerRef.current.scrollHeight;
      console.log("new height is: ", newHeight);
      containerRef.current.scrollTop = newHeight - prevScrollHeight.current;
    }
  }, [messages]);

  useEffect(() => {
    if (user) {
      setMessages([]);
      setPage(1);
      setHasMore(true);
      setLoading(true);
      containerRef.current.scrollTop = 0;
      prevScrollHeight.current = containerRef.current.scrollHeight;
    }
  }, [user]);

  useEffect(() => {
    const fetchMessages = async () => {
      try {
        const data = await getMessages(user.name, page);
        console.log("fetch messages data",data)
        if (data) {
          const newMessages = data.messages;
          newMessages == null
            ? setMessages([])
            : setMessages((prev) => [...newMessages, ...prev]);
          setHasMore(data.has_more);
        }
      } catch (error) {
        console.error("Fetch error:", error);
      }
    };

    if (loading && page > 0 && hasMore) {
      fetchMessages();
      setLoading(false);
    }
  }, [loading, page]);

  const handleScroll = (e) => {
    const container = e.target;
    if (container.scrollTop < 10 && hasMore && !loading) {
      prevScrollHeight.current = container.scrollHeight;
      setPage((prev) => prev + 1);
      setLoading(true);
    }
  };

  if (user && messages) {
    return (
      <section
        ref={containerRef}
        onScroll={handleScroll}
        style={{ height: "100%", overflowY: "auto" }}
      >
        {messages.length > 0 &&
          messages.map((message, idx) => (
            // <div className={Styles.post} key={message.sent_at}>
            //     <section className={Styles.userinfo}>
            //         <div>
            //             <Image src="/iconMale.png" alt="notification" width={25} height={25} />
            //             <p>{message.author_name}</p>
            //         </div>
            //     </section>

            //     <section className={Styles.content}>
            //         {message.content}
            //     </section>

            //     <section className={Styles.footer}>
            //         <Time rfc3339={message.sent_at} />
            //     </section>
            // </div>
            <div
              key={idx}
              className={message.author_name == me ? Style.user1 : Style.user1}
              style={{ marginBottom: "10px" }}
            >
              <p>{message.author_name}</p>
              <p>{message.content}</p>
              <p>
                <Time rfc3339={message.sent_at} />
              </p>
            </div>
          ))}
        {/* {loading && <div className={Styles.post}>Loading more...</div>} */}
      </section>
    );
  } else {
    messages ? (
      <h1>Select a person to start chatting</h1>
    ) : (
      <h1>
        No messages found
        <br />
      </h1>
    );
  }
}
