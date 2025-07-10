"use client";

import { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import { SendData } from "../../../../../utils/sendData.js";
import Style from "./profile.module.css";
import Image from "next/image.js";

export default function Profile() {
  const { groupname } = useParams();
  const [activeSection, setActiveSection] = useState("posts");

  console.log("groupname:", groupname);

  const [data, setData] = useState(null);

  useEffect(() => {
    async function fetchData() {
      try {
        const res = await SendData(`/api/v1/get/groupData`, { "groupName": groupname });
        if (res.ok) {
          const profileData = await res.json();
          setData(profileData);
        }
      } catch (err) {
        console.error("Error fetching profile:", err);
      }
    }
    fetchData();
  }, [groupname]);

  return (
    console.log("data:", data),

    <div className={Style.container}>
      <div className={Style.header}>
        <Image
          src={data?.group?.Avatar?.Valid ? data.group.Avatar.String : "/groupsBg.png"}
          fill
          alt="cover"
        />
      </div>

      <div className={Style.body}>
        <div className={Style.first}>
          <div className={Style.profileInfo}>
            <div
              style={{
                position: "relative",
                width: "200px",
                height: "200px",
              }}
            >
              <Image
                src={data?.group?.Avatar?.Valid ? data.group.Avatar.String : "/groupsBg.png"}
                alt="cover"
                fill
              />
            </div>
            <h4>{groupname}</h4>

            <div className={Style.privacy}>
              <Image
                src={`/${data?.group?.Privacy}.svg`}
                alt="privacy"
                width={20}
                height={20}
              />
              <p>&nbsp;</p>
              <p>{data?.group?.Privacy}</p>
              <p>&nbsp; - &nbsp;</p>
              <p style={{ fontWeight: "bold" }}>{data?.group?.Members?.length}</p>
            </div>

            <div className={Style.center}>
              <span>
                <h5>Description:</h5>
                <h5>{data?.group?.Description}</h5>
              </span>
            </div>

            <div className={Style.center}>
              <span>
                <h5>Members:</h5>
                <h5>{data?.group?.Members?.length}</h5>
              </span>
            </div>
          </div>
        </div>

        <div className={Style.second}>
          {activeSection === "Settings"} {/*&& <Settings />*/}

          {activeSection === "posts"} {/*&& <Posts />*/}

          {activeSection === "members"} {/*&& <Members />*/}
        </div>

        {/* events */}
        <div className={Style.events}>
          <div className={Style.event}>
            <h3>Event 1</h3>
            <p>Event description</p>
          </div>
          <div className={Style.event}>
            <h3>Event 2</h3>
            <p>Event description</p>
          </div>
          <div className={Style.event}>
            <h3>Event 3</h3>
            <p>Event description</p>
          </div>
        </div>
      </div>

      <h1>Profile of group: {groupname}</h1>
      {/* <p>{data.description}</p> */}
      {/* Render more profile info here */}
    </div>
  );
}

