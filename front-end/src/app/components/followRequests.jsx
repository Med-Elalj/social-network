"use client";

import { useEffect, useState } from "react";
import { useNotification } from "../context/NotificationContext";
import { SendData } from "@/app/sendData.js";
import Style from "../profile/profile.module.css";
import Image from "next/image";
import Link from "next/link";

export function FollowRequestsList() {
  const [followRequests, setFollowRequests] = useState([]);
  const { showNotification } = useNotification();

  useEffect(() => {
    async function handleRequests() {
      const response = await SendData("/api/v1/get/requests", { type: 0 });
      if (response.ok) {
        const followRequestData = await response.json();

        followRequestData?.length > 0
          ? setFollowRequests(followRequestData)
          : setFollowRequests([]);
      } else {
        console.log("error here")
      }
    }

    handleRequests();
  }, []);

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
      {followRequests?.length > 0 ? (
        followRequests.map((user) => (
          <div key={user.sender_id}>
            <Link href={`profile/${user.username}`}>
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
              <h2 className={Style.newStatus}>{`Follow ${user.status}ed`}</h2>
            ) : (
              <div className={Style.Buttons}>
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
        <h3
          style={{
            textAlign: "center",
            color: "var(--secondary-color",
            fontWeight: 500,
            fontSize: ".87em",
          }}
        >
          No Requests
        </h3>
      )}
    </>
  );
}
