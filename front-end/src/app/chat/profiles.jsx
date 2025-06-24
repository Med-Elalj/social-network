import { useEffect, useState } from "react";
import Image from "next/image";
import Styles from "../global.module.css";
const BACKEND_URL = process.env.NEXT_PUBLIC_BACKEND_URL;

export default function Profile() {
  const [personalDiscussions, setPersonalDiscussions] = useState([]);
  const [groupDiscussions, setGroupDiscussions] = useState([]);

  useEffect(() => {
    console.log("Resolved backend URL:", BACKEND_URL);
    const fetchConversations = async () => {
      try {
        const response = await fetch(BACKEND_URL+"/api/v1/get/users", {
          method: "GET",
          credentials: "include", // Send cookies (auth)
          headers: {
            "Content-Type": "application/json"
          }
        });

        if (!response.ok) {
          throw new Error(`HTTP error! Status: ${response.status}`);
        }

        const data = await response.json();
        console.log("Conversations:", data);
        // Filter data by profile.isgroup
        const personal = (data || []).filter(profile => profile.is_group === false);
        const groups = (data || []).filter(profile => profile.is_group === true);
        setPersonalDiscussions(personal);
        setGroupDiscussions(groups);
      } catch (error) {
        console.error("Error fetching conversations:", error);
      }
    };

    fetchConversations();
  }, []);

  return (
    <span>
      <div className={Styles.groups}>
        {personalDiscussions.map((discussion) => (
          <div key={discussion.profile_id}>
            <div>
              <Image src="/iconMale.png" alt="profile" width={40} height={40} />
              <h5>{discussion.profile_name}</h5>
            </div>
          </div>
        ))}
      </div>
      <div className={Styles.groups}>
        {groupDiscussions.map((discussion) => (
          <div key={discussion.profile_id}>
            <div>
              <Image src="/iconMale.png" alt="profile" width={40} height={40} />
              <h5>{discussion.profile_name}</h5>
            </div>
          </div>
        ))}
      </div>
    </span>
  );
}