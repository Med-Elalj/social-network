import Image from "next/image";
import Styles from "../global.module.css";
import { useState, useEffect } from "react";
import { GetData, SendData } from "@/app/sendData.js";
import { useNotification } from "../context/notificationContext.jsx";

export default function Groups() {
  const [groups, setGroups] = useState([]);
  const [joinedGroupId, setJoinedGroupId] = useState(null);
  const {showNotification} = useNotification();

  useEffect(() => {
    const fetchData = async () => {
      const response = await GetData("/api/v1/get/groupToJoin");
      const body = await response.json();

      if (response.status !== 200) {
        console.log("Faild to get groups");
      } else {
        if (body.groups?.length>0) {
          setGroups(body.groups.filter(group=>!group.IsRequested));
        }
        console.log("Groups fetched successfully!");
      }
    };

    fetchData();
  }, []);

  useEffect( () => {
     async function sentJoinHandler() {
            console.log("group id to join", joinedGroupId)
            const response = await SendData("/api/v1/set/joinGroup", { "groupId": joinedGroupId })
            let type = "error"
            const data = await response.json()
            if (response.ok) {
                type = "succes"
                setGroups(prev => prev.filter(group => group.ID != joinedGroupId))
            }
            showNotification(data.message, type)
        }
        if (joinedGroupId) {
            sentJoinHandler()
            setJoinedGroupId(null)
        }
  }, [joinedGroupId])

  return (
    <div className={Styles.groups}>
      <h1>Groups</h1>
      {groups?.length>0 ? groups.slice(0, 5).map((Group, i) => (
        <div key={Group.ID}>
          <div>
            {/* <Image src={Group?.Avatar?.String || "/db.png"} alt="profile" width={40} height={40} style={{ borderRadius: '50%' }} /> */}
            <h5>{Group.GroupName}</h5>
          </div>
          <div >
            <Image onClick={()=>setJoinedGroupId(Group.ID)} src="/join.svg" alt="join" width={25} height={25} />
          </div>
        </div>
      )) : <h3 style={{ textAlign: "center" }}>No groups to join</h3>}
    </div>
  );
}
