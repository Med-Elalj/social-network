"use client";

import { useEffect, useState } from "react";
import { GetData, SendData } from "../../../utils/sendData.js";
import { useNotification } from "../context/notificationContext.jsx";
import Image from "next/image";

export function SuggestionList() {
  const [users, setUsers] = useState([]);
  const [sentUser, setSentUser] = useState(null);
  const { showNotification } = useNotification();

  useEffect(() => {
    const fetchSugguestion = async () => {
      const response = await GetData("/api/v1/get/userSeggestions", {
        is_user: 1,
      });
      if (response.ok) {
        const data = await response.json();
        Array.isArray(data) ? setUsers(data) : setUsers([]);
      } else {
        console.error("error getting user suggesions ");
      }
    };

    fetchSugguestion();
  }, []);

  useEffect(() => {
    const responseFetch = async () => {
      const response = await SendData("/api/v1/set/follow", {
        target: sentUser.id,
        status: sentUser.status,
      });
      if (response.ok) {
        showNotification(`${sentUser.status} sent succeffully`);
        setUsers((prev) => {
          return prev.map((user) => {
            {
              console.log(user);
              return user.id == sentUser.id
                ? { ...user, isFollowSent: true }
                : user;
            }
          });
        });
      }
    };

    if (sentUser != null) {
      responseFetch();
    }
  }, [sentUser]);

  useEffect(() => console.log("users suggestions", users), [users]);

  return users?.length > 0 ? (
    users.map((user) => {
      return (
        <div key={user.id}>
          <div>
            <Image
              src={user?.pfp?.String ? user.pfp.String : "/iconeMale.png"}
              alt={user.name}
              width={40}
              height={40}
              style={{ borderRadius: "50%" }}
            />
            <h5>{user.name}</h5>
          </div>
          {!user.isFollowSent ? (
            <div
              onClick={() => setSentUser({ id: user.id, status: user.status })}
            >
              <Image src="/addUser.svg" alt="profile" width={25} height={25} />
            </div>
          ) : (
            <></>
          )}
        </div>
      );
    })
  ) : (
    <h3
      style={{
        color: "var(--secondary-color)",
        fontWeight: 400,
        fontSize: "0.9em",
      }}
    >
      No user suggestions for now
    </h3>
  );
}
