'use client';

import Image from "next/image";
import Styles from "./style.module.css";
import Time from "./time";
import Input from "./input";
import { useEffect, useRef, useState, useCallback } from "react";

export async function getMessages(person_name, page) {
    const response = await fetch(`/api/v1/get/dmhistory`, {
        method: 'POST',
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

export default function Discussion({ person }) {
    const [messages, setMessages] = useState([]);
    const [loading, setLoading] = useState(false);
    const [hasMore, setHasMore] = useState(true);
    const scrollContainerRef = useRef(null);
    const earliestTimestampRef = useRef(null); // For pagination

    const fetchInitialMessages = useCallback(async () => {
        setLoading(true);
        try {
            const data = await getMessages(person[1]);
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
    }, [person]);

    const fetchMoreMessages = async () => {
        if (loading || !hasMore) return;

        setLoading(true);
        try {
            const olderData = await getMessages(person[1], earliestTimestampRef.current);
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
        if (person) {
            setMessages([]);
            earliestTimestampRef.current = null;
            setHasMore(true);
            fetchInitialMessages();
        }
    }, [person, fetchInitialMessages]);

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
        if (el.scrollHeight <= el.clientHeight + 10) {
            fetchMoreMessages();
        }
        // eslint-disable-next-line
    }, [messages, loading, hasMore]);

    const addMsg = (message) => {
        messages ? setMessages(prev => [...prev, message]) : setMessages([message]);
        if (scrollContainerRef.current) {
            scrollContainerRef.current.scrollTop = scrollContainerRef.current.scrollHeight;
        }
    };

    return person && messages ? (
        <section
            ref={scrollContainerRef}
            onScroll={handleScroll}
            style={{ height: '80vh', overflowY: 'auto' }}
        >
            {messages.map((message) => (
                <div className={Styles.post} key={message.sent_at}>
                    <section className={Styles.userinfo}>
                        <div>
                            <Image src="/iconMale.png" alt="notification" width={25} height={25} />
                            <p>{message.author_name}</p>
                        </div>
                    </section>

                    <section className={Styles.content}>
                        {message.content}
                    </section>

                    <section className={Styles.footer}>
                        <Time rfc3339={message.sent_at} />
                    </section>
                </div>
            ))}
            {loading && <div className={Styles.post}>Loading more...</div>}
            <Input addMessage={addMsg} target={person} />
        </section>
    ) :
        messages ?
            (<div className={Styles.post}>Select a person to start chatting</div>)
            :
            (<div className={Styles.post}>No messages found<br /><Input addMessage={addMsg} target={person} /></div>)
        ;
}