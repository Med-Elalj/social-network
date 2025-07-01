'use client';

import Style from "./chat.module.css";
import ChatInput from "./input";
import Time from "./time";
import { useEffect, useRef, useState } from "react";

const me = 0;

export async function getMessages(person_name, page) {
    const response = await fetch(`/api/v1/get/dmhistory`, {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
            'target': person_name,
            'page': page,
        },
    });

    if (!response.ok) {
        throw new Error('Failed to fetch messages');
    }

    return await response.json();
}

export default function Messages({ user }) {
    const [messages, setMessages] = useState([]);
    const [hasMore, setHasMore] = useState(false);
    const [loading, setLoading] = useState(false);
    const [page, setPage] = useState(0);

    const containerRef = useRef(null);
    const prevScrollHeight = useRef(0);

    useEffect(() => {
        if (user) {
            setMessages([]);
            setPage(1);
            setLoading(true);
        }
    }, [user]);

    useEffect(() => {
        const fetchMessages = async () => {
            try {
                const data = await getMessages(user.name, page);
                if (data) {
                    const newMessages = data.messages;
                    newMessages == null ? setMessages([]) : setMessages((prev) => [...newMessages, ...prev]);
                    setHasMore(data.has_more);

                    setTimeout(() => {
                        if (containerRef.current && prevScrollHeight.current) {
                            const newHeight = containerRef.current.scrollHeight;
                            containerRef.current.scrollTop = newHeight - prevScrollHeight.current;
                        }
                    }, 0);
                }
            } catch (error) {
                console.error("Fetch error:", error);
            }
        };

        if (loading && page > 0) {
            fetchMessages();
        }
    }, [loading, page]);

    const handleScroll = (e) => {
        const container = e.target;
        if (container.scrollTop < 50 && hasMore && !loading) {
            prevScrollHeight.current = container.scrollHeight;
            setLoading(true);
            setPage((prev) => prev + 1);
        }
    };


    if (user && messages) {
        return <section
            onScroll={handleScroll}
            style={{ height: '80vh', overflowY: 'auto' }}
        >
            {messages.map((message) => (
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
                <div className={message.author_name == me ? Style.user1 : Style.user1}>
                    <p>{message.author_name}</p>
                    <p>{message.content}</p>
                    <p><Time rfc3339={message.sent_at} /></p>
                </div>
            ))}
            {/* {loading && <div className={Styles.post}>Loading more...</div>} */}
        </section>
    } else {
        messages ?
            (<h1>Select a person to start chatting</h1>)
            :
            (<h1>No messages found<br /><ChatInput addMessage={setMessages} target={user} /></h1>)
            ;
    }
}