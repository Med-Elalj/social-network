"use client";

import { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import { SendData } from "../../../../../utils/sendData.js";
import Style from "./profile.module.css";
import Image from "next/image.js";
import Settings from "./[tab]/Settings.jsx";
import Posts from "./[tab]/posts.jsx";
import Members from "./[tab]/members.jsx";
import CreatePost from "./[tab]/createPost.jsx";

export default function Profile() {
  const { groupname } = useParams();
  const [activeSection, setActiveSection] = useState("posts");
  const [data, setData] = useState(null);
  const [isPublic, setIsPublic] = useState(false);
  const [isLoading, setIsLoading] = useState(true);
  const [hasError, setHasError] = useState(false);
  const [events, setEvents] = useState([]);

  useEffect(() => {
    async function fetchData() {
      try {
        const res = await SendData(`/api/v1/get/groupData`, { "groupName": groupname });
        if (res.ok) {
          const profileData = await res.json();
          setData(profileData);
          setIsPublic(profileData?.Privacy || false);
        } else {
          setHasError(true);
        }
      } catch (err) {
        setHasError(true);
        console.error("Error fetching profile:", err);
      } finally {
        setIsLoading(false);
      }
    }
    fetchData();
  }, [groupname]);

  // get events
  useEffect(() => {
    async function fetchEvents() {
      if (!data?.ID) return;
      try {
        const res = await SendData(`/api/v1/get/groupEvents`, data?.ID);
        if (res.ok) {
          const eventData = await res.json();
          setEvents(eventData);
        } else {
          setHasError(true);
        }
      } catch (err) {
        setHasError(true);
        console.error("Error fetching events:", err);
      }
    }
    fetchEvents();
  }, [data?.ID]);

  if (isLoading) {
    return <div>Loading...</div>;
  }

  if (hasError) {
    return <div>Error loading group data.</div>;
  }

  return (
    <div className={Style.container}>
      <div className={Style.header}>
        <Image
          src={data?.Avatar?.Valid ? data.Avatar.String : "/groupsBg.png"}
          fill
          alt="cover"
        />
      </div>

      <div className={Style.body}>
        <div className={Style.first}>
          <div className={Style.profileInfo}>
            <div className={Style.avatar}>
              <Image
                src={data?.Avatar?.Valid ? data.Avatar.String : "/groupsBg.png"}
                alt="cover"
                layout="fill"
              />

            </div>
            <h4>{groupname}</h4>

            <div className={Style.privacy}>
              <Image
                src={!isPublic ? `/private.svg` : `/public.svg`}
                alt="privacy"
                width={20}
                height={20}
              />
              <p>{!isPublic ? "Private" : "Public"}</p>
              <p>&nbsp; - &nbsp;</p>
              <p style={{ fontWeight: "bold" }}>{data?.MemberCount} members</p>
            </div>

            <div className={Style.center}>
              {
                data?.Description?.String !== "" ? (
                  <span>
                    <h5>Description:</h5>
                    <p>&nbsp;&nbsp;</p>
                    <h5>{data?.Description?.String}</h5>
                  </span>
                ) : (
                  <span>
                    <h5>Description:</h5>
                    <p>&nbsp;&nbsp;</p>
                    <h5>No description provided.</h5>
                  </span>
                )
              }
            </div>

            <div>

              <button
                className={activeSection === "posts" ? Style.active : ""}
                onClick={() => setActiveSection("posts")}
              >
                Posts
              </button>
              <button
                className={activeSection === "members" ? Style.active : ""}
                onClick={() => setActiveSection("members")}
              >
                Members
              </button>
              <button
                className={activeSection === "Settings" ? Style.active : ""}
                onClick={() => setActiveSection("Settings")}
              >
                Settings
              </button>
            </div>
          </div>
        </div>

        <div className={Style.second}>
          {isPublic && (
            <>
              <div>
              </div>
              <div>
                {activeSection === "posts" && <Posts activeSection={activeSection} setActiveSection={setActiveSection} groupId={data?.ID} />}
                {activeSection === "members" && <Members groupId={data?.ID} />}
                {activeSection === "createPost" && <CreatePost groupId={data?.ID} setActiveSection={setActiveSection} />}
                {activeSection === "Settings" && <Settings />}
              </div>
            </>
          )}
        </div>

        {/* events */}
        <div className={Style.events}>
          <h2>Upcoming Events</h2>
          {events.length > 0 ? (
            events.map((event, index) => (
              <div className={Style.event} key={index}>
                <h3>{event.name}</h3>
                <p>{event.description}</p>
                <div>
                  <button>Going</button>
                  <button>Not going</button>
                </div>
              </div>
            ))
          ) : (
            <p>No upcoming events.</p>
          )}
        </div>
      </div >
    </div >
  );
}

