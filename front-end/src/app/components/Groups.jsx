import Image from "next/image";
import Styles from "../global.module.css";
import { useState, useEffect } from "react";
import { GetData, SendData } from "@/app/sendData.js";
import { useNotification } from "../context/notificationContext.jsx";

export default function Groups() {
  const [groups, setGroups] = useState([]);
  const [joinedGroupId, setJoinedGroupId] = useState(null);
  const [requests, setRequests] = useState([]);
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
          setGroups([]);
        }
        console.log("Groups fetched successfully!");
      }
    };

    fetchData();
  }, []);

  useEffect(() => {
    const fetchGroupRequests = async () => {
      try {
        const response = await SendData("/api/v1/get/requests", { type: 1 });

        if (!response.ok) {
          console.error("Failed to fetch group requests");
          return;
        }

        const data = await response.json();
        setRequests(data);
        console.log("requests groups:", data);
      } catch (err) {
        console.error("Error fetching group requests:", err);
      }
    };

    fetchGroupRequests();
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
      } else {
        console.error("failed to send join request to this group");
      }
      showNotification(data.message, type);
    }
    console.log("group id to join", joinedGroupId);
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
            <div key={Group.id}>
              <div>
                <Image
                  src={Group.pfp?.String || "/iconGroup.png"}
                  alt="profile"
                  width={40}
                  height={40}
                  style={{ borderRadius: "50%" }}
                />
                <h5>{Group.name}</h5>
              </div>
              <div>
                <Image
                  onClick={() => setJoinedGroupId(Group.id)}
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
        {requests?.length > 0 ? requests.map((Group,_) => (
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
