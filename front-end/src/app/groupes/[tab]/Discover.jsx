import { SendData } from "@/app/sendData.js";
import { useState, useEffect } from "react";
import Image from "next/image";
import Link from "next/link.js";
import Style from "../groups.module.css";
import { useNotification } from "../../context/notificationContext.jsx";

export default function Discover() {
  const [groups, setGroups] = useState([]);
  const [joinedGroupId, setJoinedGroupId] = useState(null);
  const { showNotification } = useNotification();

  useEffect(() => {
    const fetchData = async () => {
      const response = await SendData("/api/v1/get/userSeggestions", { is_user: 0 });
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
      const response = await SendData("/api/v1/set/joinGroup", {
        groupId: joinedGroupId,
      });
      let type = "error";
      const data = await response.json();
      if (response.ok) {
        type = "succes";
        setGroups((prev) =>
          prev.map((group) => {
            return group.ID == joinedGroupId
              ? { ...prev, IsRequested: true }
              : prev;
          })
        );
      }
      showNotification(data.message, type);
    }
    if (joinedGroupId) {
      sentJoinHandler();
      setJoinedGroupId(null);
    }
  }, [joinedGroupId]);

  return (
    <div className={groups ? Style.groupGrid : Style.noPosts}>
      {groups ? (
        groups.map((Group) => (
          <div className={Style.groupCard} key={Group.ID}>
            <Image
              src={Group?.pfp?.Valid ? Group?.pfp?.String : "/iconGroup.png"}
              alt="profile"
              width={50}
              height={50}
              sizes="(max-width: 768px) 100vw, 250px"
              className={Style.groupAvatar}
            />
            <h4>{Group.name}</h4>
            {!Group.IsRequested ? (
              <h3
                onClick={() => setJoinedGroupId(Group.id)}
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
