"use client";

import Image from "next/image";
import Styles from "../global.module.css";
import { useState, useEffect } from "react";
import { useNotification } from "../context/notificationContext.jsx";
import { useWebSocket } from "../context/WebSocketContext.jsx";
import Link from "next/link";

export default function Friends() {
    const [followRequests, setFollowRequests] = useState([]);
    const [contacts, setContacts] = useState([]);
    const { showNotification } = useNotification();
    const { updateOnlineUser, newNotification } = useWebSocket();

    useEffect(() => {
        async function fetchRequest(url, body) {
            try {
                const response = await fetch(url, {
                    method: 'POST',
                    credentials: 'include',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: body,
                });

                if (!response.ok) {
                    console.log(`HTTP error! Status: ${response.status}`);
                    return
                }

                const data = await response.json();
                return data


            } catch (error) {
                console.error("Error fetching requests:", error)
            }
        }

        async function handleRequests() {
            const usersData = await fetchRequest("/api/v1/get/users")
            usersData ? setContacts(usersData) : setContacts([])

            const followRequestData = await fetchRequest("/api/v1/get/requests", JSON.stringify({ type: 0 }))
            followRequestData ? setFollowRequests(followRequestData) : setFollowRequests([])
        }

        handleRequests()


    }, [])

    useEffect(() => {
        if (updateOnlineUser?.uid) {
            setContacts(prev => {
                if (prev?.length > 0) {
                    return prev.map((user) =>
                        user.id === updateOnlineUser.uid
                            ? { ...user, online: updateOnlineUser.value }
                            : user
                    );
                }
            });
        }
    }, [updateOnlineUser])

    useEffect(() => {
        if (newNotification?.command == "followRequest" && newNotification?.value == "request" && followRequests) {
            setRequests(prev => [...prev, { Uid: newNotification.uid, Name: newNotification.sender }])
        }
    }, [newNotification])

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
        } catch (error) {
            console.error(error)
            showNotification(`can't ${status} request, try again`, "error")
        }
    }

    return (
        <>
            <div className={Styles.Requiests}>
                <h1>Follow requests</h1>
                {followRequests?.length > 0 ? (followRequests.map((user) => (
                    <Link href={`/profle/${user.Name}`} key={user.Uid}>
                        <div>
                            <Image src={user.pfp?.String ? user.pfp.String : "/iconMale.png"} style={{ borderRadius: "50%" }} alt="profile" width={40} height={40} />
                            <h5>{user.Name}</h5>
                        </div>
                        {user.status ? <h2>{`Follow ${user.status}ed`}</h2> : (<div className={Styles.Buttons}>
                            <div onClick={() => responseHandle(user.uid, "accept")}>
                                <Image src="/accept.svg" alt="accept" width={30} height={30} />
                            </div>
                            <div onClick={() => responseHandle(user.uid, "reject")}>
                                <Image src="/reject.svg" alt="reject" width={30} height={30} />
                            </div>
                        </div>)}

                    </Link>
                ))) : <h3 style={{ textAlign: "center" }}>No Requests</h3>}
            </div>

            <div className={Styles.friends}>
                <h1>Contacts</h1>
                {contacts?.length > 0 ? (contacts.map((user) => (
                    <Link href={`/chat?goTo=${user.name}`} key={user.id}>
                        <div>
                            <Image src={user.pfp?.String ? user.pfp.String : "/iconMale.png"} alt="profile" width={40} height={40} style={{ borderRadius: "50%" }} />
                            <h5>{user.name}</h5>
                        </div>
                        <p className={user.online ? Styles.online : Styles.offline}>{user.online ? "online" : "offline"}</p>
                    </Link>
                ))) : <h3 style={{textAlign:"center"}}>Go Get followers</h3>}

            </div>
        </>

    )
}