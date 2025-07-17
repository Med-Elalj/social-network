"use client";

import Image from "next/image";
import Styles from "../global.module.css";
import { useState, useEffect } from "react";
import { useNotification } from "../context/notificationContext.jsx";
import { useWebSocket } from "../context/WebSocketContext.jsx";
import Link from "next/link";
import { SendData } from "../sendData.js";

export default function Friends() {
  const [followRequests, setFollowRequests] = useState([]);
  const [contacts, setContacts] = useState([]);
  const { showNotification } = useNotification();
  const { updateOnlineUser, newFollowRequest } = useWebSocket();

  useEffect(() => {
    async function fetchRequest(url, body) {
      try {
        const response = await SendData(url, body);

        if (!response.ok) {
          console.log(`HTTP error! Status: ${response.status}`);
          return;
        }

        const data = await response.json();
        return data;
      } catch (error) {
        console.error("Error fetching requests:", error);
      }
    }

    async function handleRequests() {
      const usersData = await fetchRequest("/api/v1/get/users");
      usersData ? setContacts(usersData) : setContacts([]);

      const followRequestData = await fetchRequest("/api/v1/get/requests", {
        type: 0,
      });
      followRequestData
        ? setFollowRequests(followRequestData)
        : setFollowRequests([]);
    }

    handleRequests();
  }, []);

  useEffect(() => {
    if (newFollowRequest) {
      console.log("Received newFollowRequest:", newFollowRequest);

      setFollowRequests((prev) => {


        const exists = prev.some((req) => req.sender_id === newFollowRequest.sender_id);
        if (!exists) {
          return [
            ...prev,
            {
              sender_id: newFollowRequest.sender.id,
              username: newFollowRequest.receiver.display_name,
              status: null,
            },
          ];
        }
        return prev;
      });
    }
  }, [newFollowRequest]);

  useEffect(() => {
    if (updateOnlineUser?.uid) {
      setContacts((prev) => {
        if (prev?.length > 0) {
          return prev.map((user) =>
            user.id === updateOnlineUser.uid
              ? { ...user, online: updateOnlineUser.value }
              : user
          );
        }
      });
    }
  }, [updateOnlineUser]);

  async function responseHandle(id, status) {
    const response = await SendData(`/api/v1/set/acceptFollow`, {
      status: status,
      sender: id,
      type: 0,
    });
    const data = await response.json();
    if (response.ok) {
      showNotification(data.message, response.ok ? "succes" : "error");
      setFollowRequests((prev) =>
        prev.map((item) =>
          item.sender_id === id ? { ...item, status: status } : item
        )
      );
    } else {
      console.error(error);
      showNotification(`can't ${status} request, try again`, "error");
    }
  }

  return (
    <>
      <div className={Styles.Requiests}>
        <h1>Follow requests</h1>
        {followRequests?.length > 0 ? (
          followRequests.map((user) => (
            <div
              key={user.id}
              style={{
                display: "flex",
                alignItems: "center",
                alignItems: "center",
              }}
            >
              <Link
                href={`profile/${user.username}`}
                style={{ display: "flex", alignItems: "center" }}
              >
                <Image
                  src={
                    user.group_avatar?.String
                      ? user.group_avatar.String
                      : "/iconMale.png"
                  }
                  style={{ borderRadius: "50%" }}
                  alt="profile"
                  width={40}
                  height={40}
                />
                <h5>{user.username}</h5>
              </Link>
              {user.status ? (
                <h2>{`Follow ${user.status}ed`}</h2>
              ) : (
                <div className={Styles.Buttons}>
                  <div onClick={() => responseHandle(user.sender_id, "accept")}>
                    <Image
                      src="/accept.svg"
                      alt="accept"
                      width={30}
                      height={30}
                    />
                  </div>
                  <div onClick={() => responseHandle(user.sender_id, "reject")}>
                    <Image
                      src="/reject.svg"
                      alt="reject"
                      width={30}
                      height={30}
                    />
                  </div>
                </div>
              )}
            </div>
          ))
        ) : (
          <h3 style={{ textAlign: "center" }}>No Requests</h3>
        )}
      </div>

      <div className={Styles.friends}>
        <h1>Contacts</h1>
        {contacts?.length > 0 ? (
          contacts.map((user) => (
            <Link href={`/chat?goTo=${user.name}`} key={user.id}>
              <div>
                <Image
                  src={user.pfp?.String ? user.pfp.String : "/iconMale.png"}
                  alt="profile"
                  width={40}
                  height={40}
                  style={{ borderRadius: "50%" }}
                />
                <h5>{user.name}</h5>
              </div>
              <p className={user.online ? Styles.online : Styles.offline}>
                {user.online ? "online" : "offline"}
              </p>
            </Link>
          ))
        ) : (
          <h3 style={{ textAlign: "center" }}>Go Get followers</h3>
        )}
      </div>
    </>
  );
}
