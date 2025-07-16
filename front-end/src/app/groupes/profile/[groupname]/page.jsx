"use client";

import { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import { SendData } from "../../../sendData.js";
import Style from "./profile.module.css";
import Image from "next/image.js";
import CreateEvent from "./[tab]/CreateEvenet.jsx";
import Posts from "./[tab]/posts.jsx";
import Members from "./[tab]/members.jsx";
import CreatePost from "./[tab]/createPost.jsx";
import { useNotification } from "../../../context/notificationContext.jsx";

export default function Profile() {
  const { groupname } = useParams();
  const [activeSection, setActiveSection] = useState("posts");
  const [data, setData] = useState(null);
  const [isPublic, setIsPublic] = useState(false);
  const [isLoading, setIsLoading] = useState(true);
  const [hasError, setHasError] = useState(false);
  const [requests, setRequests] = useState([]);
  const [events, setEvents] = useState([]);
  const [respondUserRequest, setRespondUserRequest] = useState(null);
  const {showNotification} = useNotification();

  useEffect(() => {
    async function fetchData() {
      try {
        const res = await SendData(`/api/v1/get/groupData`, {
          groupName: groupname,
        });
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

  // get requests
  useEffect(() => {
    const fetchData = async () => {
      const response = await SendData("/api/v1/get/requests", { type: 1 });
      const Body = await response.json();
      if (!response.ok) {
        console.log(Body);
      } else {
        setRequests(Body);
        console.log("requests fetched successfully!");
      }
    };

    fetchData();
  }, []);

  useEffect( ()=> {
    const fetchResponse = async () => {
      const response = await SendData("/api/v1/set/acceptFollow",{sender: respondUserRequest.id, target: data.ID, status:respondUserRequest.status, type:1})
      if (response.ok) {
        const responseData=await response.json();
        showNotification(responseData.message);
      } else {
        showNotification("failed to accept or refuse request", "error")
      }
    }

    if (respondUserRequest) {
      fetchResponse()
    }
  }, [respondUserRequest])

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
              {data?.Description?.String !== "" ? (
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
              )}
            </div>

            <div className={Style.buttons}>
              <button
                className={Style.button}
                onClick={() => setActiveSection("posts")}
              >
                Posts
              </button>
              <button
                className={Style.button}
                onClick={() => setActiveSection("members")}
              >
                Members
              </button>
              <button
                className={Style.button}
                onClick={() => setActiveSection("CreateEvent")}
              >
                Create Event
              </button>
            </div>
          </div>
          {data?.IsAdmin ? (
            <div className={Style.requests}>
              {/* reauests to join group */}
              <h2>Requests</h2>
              <div className={Style.requestsContainer}>
                {requests ? (requests.map((request) => (
                  <div key={request.sender_id} className={Style.request}>
                    <div className={Style.avatar}>
                      <Image
                        src={
                          request.avatar?.Valid
                            ? request.avatar.String
                            : "/iconMale.png"
                        }
                        width={25}
                        height={25}
                        alt="avatar"
                      />
                      <h4>{request.username ?? "User"}</h4>
                    </div>
                    <div>
                      <p>{respondUserRequest?.id != request.sender_id ? request.message : `request ${respondUserRequest.status}ed`}</p>
                    </div>
                    <div>
                      {respondUserRequest?.id != request.sender_id ?  
                      (<><button
                        className={Style.button}
                        onClick={() => {
                          setRespondUserRequest({
                            id: request.sender_id,
                            status: "accept",
                          });
                        }}
                      >
                        Accept
                      </button>
                      <button
                        className={Style.button}
                        onClick={() => {
                          setRespondUserRequest({
                            id: request.sender_id,
                            status: "refuse",
                          });
                        }}
                      >
                        Decline
                      </button></> ): <></>}
                    </div>
                  </div>
                ))):<h3>no requests</h3>}
              </div>
            </div>
          ) : (
            ""
          )}
        </div>

        <div className={Style.second}>
          {isPublic && (
            <>
              <>
                {activeSection === "posts" && (
                  <Posts
                    activeSection={activeSection}
                    setActiveSection={setActiveSection}
                    groupId={data?.ID}
                  />
                )}
                {activeSection === "members" && <Members groupId={data?.ID} />}
                {activeSection === "createPost" && (
                  <CreatePost
                    groupId={data?.ID}
                    setActiveSection={setActiveSection}
                  />
                )}
                {activeSection === "CreateEvent" && <CreateEvent />}
              </>
            </>
          )}
        </div>

        {/* events */}
        <div className={Style.events}>
          <h2>Upcoming Events</h2>
          {events.length > 0 ? (
            events.map((event, index) => (
              <div key={index} className={Style.event}>
                <h3>{event.name}</h3>
                <p>{event.description}</p>
                <p>{event.date}</p>
                <div>
                  <button className={Style.button}>Going</button>
                  <button className={Style.button}>Not going</button>
                </div>
              </div>
            ))
          ) : (
            <p>No upcoming events.</p>
          )}
        </div>
      </div>
    </div>
  );
}
