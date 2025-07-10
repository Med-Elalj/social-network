"use client";

import Image from "next/image";
import Styles from "../global.module.css";
import { useState, useEffect } from "react";
import { useNotification } from "../context/notificationContext.jsx";
import { useWebSocket } from "../context/WebSocketContext.jsx";
import Link from "next/link";

export default function Friends() {
    const [requests, setRequests] = useState([]);
    const [contacts, setContacts] = useState([]);
    const { showNotification } = useNotification();
    const { updateOnlineUser } = useWebSocket();

    useEffect(() => {
        async function fetchFollowRequest(request) {
            try {
                const response = await fetch(`/api/v1/get/${request}`, {
                    method: 'GET',
                    credentials: 'include',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                });

                if (!response.ok) {
                    console.error(`HTTP error! Status: ${response.status}`);
                }

                const data = await response.json();
                console.log(request, data)
                if (request == "follow") {
                    setRequests(data.users || [])
                } else if (request == "users") {
                    data ? setContacts(data) : setContacts([])
                }
            } catch (error) {
                console.error("Error fetching follow requests:", error)
            }
        }

        fetchFollowRequest("follow")
        fetchFollowRequest("users")

    }, [])

    useEffect(() => {
        if (updateOnlineUser?.uid) {
            setContacts(prev => {
                if (prev?.length > 0) {
                    // Return the updated array from map
                    return prev.map(user =>
                        user.id === updateOnlineUser.uid
                            ? { ...user, online: updateOnlineUser.value }
                            : user
                    );
                }
                return prev; // Return the previous state if no contacts
            });
        }
    }, [updateOnlineUser])

    async function responseHandle(id, status) {
        try {
            const response = await fetch(`/api/v1/follow/${status}`, {
                method: 'POST',
                credentials: 'include',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    uid: id,
                    status: status
                })
            });
            const data = await response.json(); // Added 'const' declaration
            showNotification(data.message, response.ok ? "succes" : "error")
            setRequests(prev => prev.filter((user) => user.id != id))
        } catch (error) {
            console.error(error)
            showNotification(`can't ${status} request, try again`, "error")
        }
    }

    return (
        <>
            <div className={Styles.Requiests}>
                <h1>Friend requests</h1>
                {requests?.length > 0 ? (requests.map((user) => (
                    <Link href={`/profle/${user.Name}`} key={user.uid}>
                        <div>
                            <Image src={user.Avatar} alt="profile" width={40} height={40} />
                            <h5>{user.Name}</h5> {/* Fixed: was "user.Name" as string */}
                        </div>
                        <div className={Styles.Buttons}>
                            <div onClick={() => responseHandle(user.uid, "accept")}> {/* Fixed: wrapped in arrow function */}
                                <Image src="/accept.svg" alt="accept" width={30} height={30} />
                            </div>
                            <div onClick={() => responseHandle(user.uid, "reject")}> {/* Fixed: changed to "reject" */}
                                <Image src="/reject.svg" alt="reject" width={30} height={30} />
                            </div>
                        </div>
                    </Link>
                ))) : <h3 style={{ textAlign: "center" }}>No Requests</h3>}
            </div>

            <div className={Styles.friends}>
                <h1>Contacts</h1>
                {contacts?.length > 0 ? (contacts.map((user) => (
                    <Link href={`/chat?goTo=${user.name}`} key={user.id}>
                        <div>
                            <Image src="/iconMale.png" alt="profile" width={40} height={40} />
                            <h5>{user.name}</h5>
                        </div>
                        <p className={user.online ? Styles.online : Styles.offline}>{user.online ? "online" : "offline"}</p>
                    </Link>
                ))) : <h2>Go Get followers</h2>}

            </div>
        </>

    )
}