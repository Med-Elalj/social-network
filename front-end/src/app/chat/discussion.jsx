'use client';

import Image from "next/image";
import Styles from "./style.module.css"
import Time from "./time";
import { useEffect, useState } from "react";
// get messages from backend POST /api/v1/get/dmhistory {person_id:123}

// export async function getMessages(person_id) {
//     const response = await fetch(`/api/v1/get/dmhistory?person_id=${person_id}`, {
//         method: 'POST',
//         headers: {
//             'Content-Type': 'application/json',
//         },
//     });

//     if (!response.ok) {
//         throw new Error('Failed to fetch messages');
//     }

//     return response.json();
// }

// Sample messages data
// This is a mockup of the messages data structure
const messages = [
    { id: 1, author_id: 1, author_name: "Name 1", sent_at: "2025-06-18T18:37:52Z", content: "Lorem ipsum dolor sit amet, consectetur adipiscing elit." },
    { id: 2, author_id: 2, author_name: "Name 2", sent_at: "2025-06-18T18:37:50Z", content: "Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua." },
    { id: 3, author_id: 3, author_name: "Name 3", sent_at: "2025-06-18T18:37:48Z", content: "Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat." },
    { id: 4, author_id: 4, author_name: "Name 4", sent_at: "2025-06-18T18:37:46Z", content: "Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur." },
    { id: 5, author_id: 5, author_name: "Name 5", sent_at: "2025-06-18T18:37:44Z", content: "Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum." },
    { id: 6, author_id: 1, author_name: "Name 1", sent_at: "2025-06-18T18:37:42Z", content: "Lorem ipsum dolor sit amet, consectetur adipiscing elit." },
    { id: 7, author_id: 2, author_name: "Name 2", sent_at: "2025-06-18T18:37:40Z", content: "Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua." },
    { id: 8, author_id: 1, author_name: "Name 1", sent_at: "2025-06-18T17:00Z", content: "Azer" },
];

// export default function Groups() {
//     return (
//         <section>

//             {messages.map((message) => (
//                 <div className={Styles.post} key={message.id}>
//                     <section className={Styles.userinfo}>

//                         <div key={message.author_id}>
//                             <Image src="/iconMale.png" alt="notification" width={25} height={25} />
//                             <p>{message.author_name}</p>
//                         </div>
//                     </section>

//                     <section className={Styles.content}>
//                         {message.content}
//                     </section>

//                     <section className={Styles.footer}>
//                         <p>time/time/time</p>
//                     </section>
//                 </div>
//             ))}
//         </section>
//     )
// }

export default function Discussion({ person_id }) {
    // const [messages, setMessages] = useState([]);

    // useEffect(() => {
    //     async function fetchMessages() {
    //         try {
    //             const data = await getMessages(person_id);
    //             setMessages(data);
    //         } catch (error) {
    //             console.error("Error fetching messages:", error);
    //         }
    //     }

    //     fetchMessages();
    // }, [person_id]);

    return (
        <section>
            {messages.map((message) => (
                <div className={Styles.post} key={message.id}>
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
        </section>
    );
}