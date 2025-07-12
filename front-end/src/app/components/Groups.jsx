import Image from "next/image";
import Link from "next/link";
import Styles from "../global.module.css";
import { useState, useEffect } from "react";
import { GetData } from "../../../utils/sendData.js";

export default function Groups() {
  const [groups, setGroups] = useState([]);

  useEffect(() => {
    const fetchData = async () => {
      const response = await GetData("/api/v1/get/groupToJoin");
      const body = await response.json();

      if (response.status !== 200) {
        console.log("Faild to get groups");
      } else {
        setGroups(body.groups);
        console.log("Groups fetched successfully!");
      }
    };

    fetchData();
  }, []);

  return (
    <div className={Styles.groups}>
      <h1>Groups</h1>
      {groups ? groups.slice(0, 5).map((Group, i) => (
        <div key={i}>
          <div>
            <Image src={Group.Avatar?.String || "/db.png"} alt="profile" width={40} height={40} style={{ borderRadius: '50%' }} />
            <h5>{Group.GroupName}</h5>
          </div>
          <div >
            <Image src="/join.svg" alt="join" width={25} height={25} />
          </div>
        </div>
      )) : <h3 style={{ textAlign: "center" }}>No groups to join</h3>}
    </div>
  );
}
