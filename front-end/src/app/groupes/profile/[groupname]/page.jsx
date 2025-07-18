"use client";

import { useEffect, useState, useCallback } from "react";
import { useParams } from "next/navigation";
import { SendData } from "../../../sendData.js";
import Style from "./profile.module.css";
import Image from "next/image.js";
import CreateEvent from "./[tab]/CreateEvenet.jsx";
import Posts from "./[tab]/posts.jsx";
import Members from "./[tab]/members.jsx";
import CreatePost from "./[tab]/createPost.jsx";
import { useNotification } from "../../../context/NotificationContext.jsx";
import { SearchInput } from "../../../components/navigation/search.jsx";
import { Countdown } from "../../../utils.jsx";

export default function Profile() {
  const [activeSection, setActiveSection] = useState("posts");
  const [data, setData] = useState(null);
  const [isPublic, setIsPublic] = useState(false);
  const [isLoading, setIsLoading] = useState(true);
  const [hasError, setHasError] = useState(false);
  const [requests, setRequests] = useState([]);
  const [events, setEvents] = useState([]);
  const [respondUserRequest, setRespondUserRequest] = useState(null);
  const [showSearch, setShowSearch] = useState(false);
  const [joinedGroupId, setJoinedGroupId] = useState(null);
  const [reactionEventRequest, setReactionEventRequest] = useState(null);
  const [fetchedReactionEvents, setFetchedReactionEvents] = useState([]);


  const { showNotification } = useNotification();
  const { groupname } = useParams();

  // Memoize the search close handler
  const handleSearchClose = useCallback(() => {
    setShowSearch(false);
    // Reset active section when search is closed
    setActiveSection("posts");
  }, []);

  // Fetch group data
  useEffect(() => {
    async function fetchData() {
      try {
        const res = await SendData(`/api/v1/get/groupData`, {
          groupName: groupname,
        });
        if (res.ok) {
          const profileData = await res.json();
          setData(profileData);
          setIsPublic(profileData?.isPublic || false);
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

    if (groupname) {
      fetchData();
    }
  }, [groupname]);

  // Fetch events - only when data.ID changes
  useEffect(() => {
    async function fetchEvents() {
      if (!data?.ID) return;

      try {
        const res = await SendData(`/api/v1/get/groupEvents`, data.ID);
        if (res.ok) {
          const eventData = await res.json();
          setEvents(eventData.events);
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

  // Fetch requests - only once on component mount
  useEffect(() => {
    const fetchRequests = async () => {
      try {
        const response = await SendData("/api/v1/get/requests", { type: 1 });
        const body = await response.json();
        if (response.ok) {
          setRequests(body);
          console.log("requests fetched successfully!");
        } else {
          console.log("Error fetching requests:", body);
        }
      } catch (err) {
        console.error("Error fetching requests:", err);
      }
    };

    fetchRequests();
  }, []); // Empty dependency array - only run once

  // Handle user request response
  useEffect(() => {
    const fetchResponse = async () => {
      if (!respondUserRequest || !data?.ID) return;

      try {
        const response = await SendData("/api/v1/set/acceptFollow", {
          sender: respondUserRequest.id,
          target: data.ID,
          status: respondUserRequest.status,
          type: 1,
        });

        if (response.ok) {
          const responseData = await response.json();
          showNotification(responseData.message);

          // Remove the processed request from the requests array
          setRequests(prev => prev.filter(req => req.sender_id !== respondUserRequest.id));
        } else {
          showNotification("failed to accept or refuse request", "error");
        }
      } catch (err) {
        console.error("Error responding to request:", err);
        showNotification("failed to accept or refuse request", "error");
      } finally {
        // Clear the respondUserRequest after processing
        setRespondUserRequest(null);
      }
    };

    fetchResponse();
  }, [respondUserRequest, data?.ID, showNotification]);

  // Handle "add members" section
  useEffect(() => {
    if (activeSection === "add members") {
      setShowSearch(true);
    } else {
      setShowSearch(false);
    }
  }, [activeSection]);

  // Send join request
  useEffect(() => {
    async function sentJoinHandler() {
      const response = await SendData("/api/v1/set/joinGroup", {
        groupId: joinedGroupId,
      });
      let type = "error";
      const data = await response.json();
      if (response.ok) {
        type = "succes";
      }
      showNotification(data.message, type);
    }
    if (joinedGroupId) {
      sentJoinHandler();
      setJoinedGroupId(null);
    }
  }, [joinedGroupId]);

  // send reaction events
  useEffect(() => {
    async function fetchReactionEvents() {
      if (!reactionEventRequest) return;

      try {
        const response = await SendData("/api/v1/set/reactionEvents", reactionEventRequest);
        if (!response.ok) return;

        const reactionEventData = await response.json();

        // ✅ Update the 'respond' field in the current events state
        setEvents((prevEvents) =>
          prevEvents.map((event) =>
            event.event_id === reactionEventRequest.event_id
              ? { ...event, respond: reactionEventRequest.response }
              : event
          )
        );
      } catch (err) {
        setHasError(true);
        console.error("Error reacting to event:", err);
      } finally {
        setReactionEventRequest(null);
      }
    }

    fetchReactionEvents();
  }, [reactionEventRequest]);


  useEffect(() => console.log("aaaaaaaaaaaaaaaaaaaaaaaaaaa events", events), [events])
  useEffect(() => console.log("aaaaaaaaaaaaaaaaaaaaaaaaaaa respondUserRequest", respondUserRequest), [respondUserRequest])

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
          src={data?.Avatar?.String ? data.Avatar.String : "/groupsBg.png"}
          fill
          alt="cover"
        />
      </div>

      <div className={Style.body}>
        <div className={Style.first}>
          <div className={Style.profileInfo}>
            <div className={Style.avatar}>
              <Image
                src={data?.Avatar?.Valid ? data.Avatar.String : "/iconGroup.png"}
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
              <span>
                {data?.Description?.String !== "" && (
                  <>
                    <h5>Description:</h5>
                    <p>&nbsp;&nbsp;</p>
                    <h5>{data?.Description?.String}</h5>
                  </>
                )}
              </span>
            </div>

            <div className={Style.buttons}>
              {(data?.IsAdmin || data?.IsMember) ? (
                <>
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
                  <button
                    className={Style.button}
                    onClick={() => setActiveSection("add members")}
                  >
                    Add members
                  </button>
                </>
              ) : (
                <>
                  <button
                    className={Style.button}
                    onClick={() => setJoinedGroupId(data.ID)}
                  >
                    Join Group
                  </button>
                </>
              )}
            </div>
          </div>

          {data?.IsAdmin && (
            <div className={Style.requests}>
              <h2>Requests</h2>
              <div className={Style.requestsContainer}>
                {requests && requests.length > 0 ? (
                  requests.map((request) => (
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
                        <p>
                          {respondUserRequest?.id !== request.sender_id
                            ? request.message
                            : `request ${respondUserRequest.status}ed`}
                        </p>
                      </div>
                      <div>
                        {respondUserRequest?.id !== request.sender_id ? (
                          <>
                            <button
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
                            </button>
                          </>
                        ) : null}
                      </div>
                    </div>
                  ))
                ) : (
                  <h3>no requests</h3>
                )}
              </div>
            </div>
          )}
        </div>

        <div className={Style.second}>

          {(data.IsAdmin || data.IsMember) && (
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
              {activeSection === "CreateEvent" && <CreateEvent groupId={data?.ID} setActiveSection={setActiveSection} />}
            </>
          )}
        </div>

        {/* events */}
        {(data.IsAdmin || data.IsMember) && (
          <div className={Style.events}>
            <h2>Upcoming Events</h2>
            {events?.length > 0 ? (

              events?.map((event) => (
                console.log(event),
                <div key={event?.event_id} className={Style.event}>
                  <h3>{event.title}</h3>
                  <p>{event.description}</p>
                  <Countdown
                    targetTimeISO={event.time}
                    onComplete={() => console.log("event started")}
                  />
                  {event.respond ?
                    <>
                      <p>You are going</p>
                    </> :
                    <>
                      <p>You are not going</p>
                    </>}
                  <div>
                    <button
                      className={Style.button}
                      onClick={() => {
                        setReactionEventRequest({
                          event_id: event.event_id,
                          response: true,
                          is_reacted: true,
                        })
                      }

                      }
                    >
                      Going
                    </button>

                    <button
                      className={Style.button}
                      onClick={() =>
                        setReactionEventRequest({
                          event_id: event.event_id,
                          response: false,
                          is_reacted: true,
                        })
                      }
                    >
                      Not Going
                    </button>
                  </div>
                </div>
              ))
            ) : (
              <p>No upcoming events.</p>
            )}
          </div>
        )}

        {showSearch && <SearchInput onClose={handleSearchClose} groupId={data?.ID} />}
      </div>
    </div >
  );
}