import { SendData } from "@/app/sendData.js";
import { useState, useEffect } from "react";
import Image from "next/image";
import Link from "next/link.js";
import Style from "../groups.module.css";
import { useNotification } from "../../context/NotificationContext.jsx";
import { useAuth } from "../../context/AuthContext.jsx";

export default function Discover() {
  const [groups, setGroups] = useState([]);
  const [joinedGroup, setJoinedGroup] = useState(null);
  const { showNotification } = useNotification();
  const { isloading, isLoggedIn } = useAuth();


  useEffect(() => {
    const fetchData = async () => {
      const response = await SendData("/api/v1/get/userSeggestions", {
        is_user: 0,
      });
      const body = await response.json();

      if (response.status !== 200) {
        console.error(body);
      } else {
        setGroups(body);
        console.log("Groups fetched successfully!");
      }
    };

    fetchData();
  }, []);

  useEffect(() => {
    async function sentJoinHandler() {
      console.log("group id to join", joinedGroup);
      const response = await SendData("/api/v1/set/sendRequest", joinedGroup);
      let type = "error";
      const data = await response.json();
      if (response.ok) {
        type = "success";
        setGroups((prev) =>
          prev.map((group) => {
            return group.id == joinedGroup.target
              ? { ...prev, IsRequested: true }
              : prev;
          })
        );

      }
      showNotification(data.message, type);
    }
    if (joinedGroup) {
      sentJoinHandler();
      setJoinedGroup(null);
    }
  }, [joinedGroup]);

  useEffect(() => {
    console.log(groups);
  }, [groups]);

  if (isloading || !isLoggedIn) return null;

  return (
    <div className={groups ? Style.groupGrid : Style.noPosts}>
      {groups ? (
        groups.map((Group) => (
          <div className={Style.groupCard} key={Group?.id}>
            <Image
              src={Group?.pfp?.String ? Group?.pfp?.String : "/iconGroup.png"}
              alt="profile"
              width={50}
              height={50}
              sizes="(max-width: 768px) 100vw, 250px"
              className={Style.groupAvatar}
            />
            <h4>{Group.name}</h4>
            <p>
              {Group.Description
                ? Group.Description
                : "No description"}
            </p>
            {!Group.IsRequested ? (
              <h3
                onClick={() => setJoinedGroup({ target: Group.id, type: 1 })}
                className={Style.acceptBtn}
              >
                Join Group
              </h3>
            ) : (
              <h3 style={{ cursor: "not-allowed" }}>waiting ...</h3>
            )}
          </div>
        ))
      ) : (
        <>
          <h3>Join groups to see feeds</h3>
          <Link href="/groupes/create">Create a group</Link>
        </>
      )}
    </div>
  );
}
