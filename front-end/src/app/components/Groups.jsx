import Image from "next/image";
import Styles from "../global.module.css";
import { useState, useEffect } from "react";
import { GetData, SendData } from "@/app/sendData.js";
import { useNotification } from "../context/notificationContext.jsx";

export default function Groups() {
  const [groups, setGroups] = useState([]);
  const [joinedGroupId, setJoinedGroupId] = useState(null);
  const { showNotification } = useNotification();

  useEffect(() => {
    const fetchData = async () => {
      const response = await GetData("/api/v1/get/userSeggestions", {
        is_user: 0,
      });
      const body = await response.json();

      if (response.status !== 200) {
        console.log("Faild to get groups");
      } else {
        if (body?.length > 0) {
          // setGroups(body.groups.filter((group) => !group.IsRequested));
          setGroups(body);
        } else {
          setGroups([])
        }
        console.log("Groups fetched successfully!");
      }
    };

    fetchData();
  }, []);

  useEffect(() => {
    async function sentJoinHandler() {
      console.log("group id to join", joinedGroupId);
      const response = await SendData("/api/v1/set/sendRequest", {
        target: joinedGroupId,
        type: 1,
      });
      let type = "error";
      const data = await response.json();
      if (response.ok) {
        type = "succes";
        setGroups((prev) => prev.filter((group) => group.id != joinedGroupId));
      }
      showNotification(data.message, type);
    }
    if (joinedGroupId) {
      sentJoinHandler();
      setJoinedGroupId(null);
    }
  }, [joinedGroupId]);

  return (
    <>
      <div className={Styles.groups}>
        <h1>Groups</h1>
        {groups?.length > 0 ? (
          groups.slice(0, 5).map((Group, i) => (
            <div key={Group.ID}>
              <div>
                <Image
                  src={Group.Avatar?.String || "/db.png"}
                  alt="profile"
                  width={40}
                  height={40}
                  style={{ borderRadius: "50%" }}
                />
                <h5>{Group.GroupName}</h5>
              </div>
              <div>
                <Image
                  onClick={() => setJoinedGroupId(Group.ID)}
                  src="/join.svg"
                  alt="join"
                  width={25}
                  height={25}
                />
              </div>
            </div>
          ))
        ) : (
          <h3 style={{ textAlign: "center" }}>No groups to join</h3>
        )}
      </div>
      <div className={Styles.groups}>
        {/* groups request */}
        <h1>Groups Request</h1>
        {[].length > 0 ? [1, 2].map((Group, i) => (
          <div key={Group.ID} className={Styles.grouprequest}>
            <div >
              <Image
                src={Group.Avatar?.String || "/db.png"}
                alt="profile"
                width={40}
                height={40}
                style={{ borderRadius: "50%" }}
              />
              <h5>{Group.GroupName || "Group name"}</h5>
            </div>
            <div className={Styles.Buttons}>
              <Image
                onClick={() => setJoinedGroupId(Group.ID)}
                src="/accept2.svg"
                alt="join"
                width={25}
                height={25}
                style={{ marginRight: "10px" }}
              />
              <Image
                onClick={() => setJoinedGroupId(Group.ID)}
                src="/decline.svg"
                alt="join"
                width={25}
                height={25}
              />
            </div>
          </div>
        )) : (
          <h3 style={{ textAlign: "center" }}>No groups request</h3>
        )}
      </div>
    </>
  );
}
