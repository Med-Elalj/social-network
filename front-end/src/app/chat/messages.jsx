'use client';

// import Image from "next/image";
import Style from "./chat.module.css";
import ChatInput from "./input";
import Time from "./time";
import { useEffect, useRef, useState, useCallback } from "react";
const me = 0

export async function getMessages(person_name, page) {
    console.log("persone name", person_name)
    const response = await fetch(`/api/v1/get/dmhistory`, {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
            'target': person_name,
            'page': page || new Date().toISOString(), // Use ISO string as page
        },
    });

    if (!response.ok) {
        throw new Error('Failed to fetch messages');
    }

    return response.json();
}

export default function Messages({ user }) {
    const [messages, setMessages] = useState([]);
    const [loading, setLoading] = useState(false);
    const [hasMore, setHasMore] = useState(true);
    const scrollContainerRef = useRef(null);
    const earliestTimestampRef = useRef(null); // For pagination

    const fetchInitialMessages = useCallback(async () => {
        setLoading(true);
        try {
            const data = await getMessages(user.name);
            setMessages(data);
            if (data != null && data.length > 0) {
                earliestTimestampRef.current = data[0].sent_at;
            } else {
                setHasMore(false);
            }
        } catch (error) {
            console.error("Initial fetch error:", error);
        } finally {
            setLoading(false);
        }
    }, [user]);

    const fetchMoreMessages = async () => {
        if (loading || !hasMore) return;

        setLoading(true);
        try {
            const olderData = await getMessages(user[1].name, earliestTimestampRef.current);
            if (olderData != null && olderData.length > 0) {
                setMessages(prev => [...olderData, ...prev]);
                earliestTimestampRef.current = olderData[0].sent_at;
            } else {
                setHasMore(false);
            }
        } catch (error) {
            console.error("Fetch more error:", error);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => {
        if (user) {
            console.log("1")
            setMessages([]);
            earliestTimestampRef.current = null;
            setHasMore(true);
            fetchInitialMessages();
        }
    }, [user, fetchInitialMessages]);

    // Scroll handler
    const handleScroll = () => {
        const el = scrollContainerRef.current;
        if (!el) return;
        if (el.scrollTop < 50) {
            fetchMoreMessages();
        }
    };

    // If messages are too few to allow scrolling, fetch more automatically
    useEffect(() => {
        const el = scrollContainerRef.current;
        if (!el || loading || !hasMore) return;
        // if (el.scrollHeight <= el.clientHeight + 10) {
        //     fetchMoreMessages();
        // }
        // eslint-disable-next-line
    }, [messages, loading, hasMore]);

    const addMsg = (message) => {
        messages ? setMessages(prev => [...prev, message]) : setMessages([message]);
        if (scrollContainerRef.current) {
            scrollContainerRef.current.scrollTop = scrollContainerRef.current.scrollHeight;
        }
    };
    if (user && messages) {
        return <section
            ref={scrollContainerRef}
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
            (<h1>No messages found<br /><ChatInput addMessage={addMsg} target={user} /></h1>)
            ;
    }
}