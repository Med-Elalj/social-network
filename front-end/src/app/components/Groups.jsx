import Image from "next/image";
import Styles from "../global.module.css";
import { useState, useEffect } from "react";
import { GetData, SendData } from "@/app/sendData.js";
import { useNotification } from "../context/NotificationContext.jsx";
import  {useAuth} from "../context/AuthContext.jsx";
import { useWebSocket } from "../context/WebSocketContext.jsx";

export default function Groups() {
  const [groups, setGroups] = useState([]);
  const [requests, setRequests] = useState([]);
  const [respond, setRespond] = useState(null);
  const { showNotification } = useNotification();
  const { newGroupRequest } = useWebSocket(); // Get from WebSocket context
  const {isloading,isLoggedIn} = useAuth();
  if (isloading) return null;

  useEffect(() => {
    if (!isLoggedIn) return;
    const fetchData = async () => {
      try {
        const response = await GetData("/api/v1/get/userSeggestions", {
          is_user: 0,
        });
        const body = await response.json();

        if (response.status !== 200) {
          console.log("Failed to get groups");
          showNotification("Failed to load groups", "error");
        } else {
          if (body?.length > 0) {
            setGroups(body);
          } else {
            setGroups([]);
          }
          console.log("Groups fetched successfully!");
        }
      } catch (error) {
        console.error("Error fetching groups:", error);
        showNotification("Error loading groups", "error");
      }
    };

    fetchData();
  }, [isLoggedIn]);

  useEffect(() => {
    console.log("New group request received:", newGroupRequest);

    if (newGroupRequest) {
      setRequests((prev) => {
        const exists = prev.some(
          (req) =>
            req.sender_id === newGroupRequest.sender.id &&
            req.group_id === newGroupRequest.target.id
        );

        if (!exists) {
          const newRequest = {
            sender_id: newGroupRequest.sender.id,
            group_id: newGroupRequest.target.id,
            username: newGroupRequest.receiver.display_name,
            group_name: newGroupRequest.target.display_name,
            message: newGroupRequest.message,
            isRespond: false,
          };

          console.log("Adding new group request:", newRequest);
          return [...prev, newRequest];
        }

        return prev;
      });
    }
  }, [newGroupRequest]);

  useEffect(() => {
    if (!isLoggedIn) return;
    const fetchGroupRequests = async () => {
      try {
        const response = await SendData("/api/v1/get/requests", { type: 1, is_special:true });

        if (!response.ok) {
          console.error("Failed to fetch group requests");
          showNotification("Failed to load group requests", "error");
          return;
        }

        const data = await response.json();
        setRequests(data || []);
        console.log("Existing group requests:", data);
      } catch (err) {
        console.error("Error fetching group requests:", err);
        showNotification("Error loading group requests", "error");
      }
    };

    fetchGroupRequests();
  }, [isLoggedIn]);

  useEffect(() => {
    const handleRequestResponse = async () => {
      try {
        console.log("Processing request:", respond);
        const response = await SendData("/api/v1/set/acceptFollow", respond);

        if (response.ok) {
          showNotification(
            `Request ${respond.status}ed successfully`,
            "success"
          );

          // Update the processed request in the requests list
          setRequests((prevRequests) =>
            prevRequests.map((req) =>
              req.sender_id === respond.sender &&
              req.group_id === respond.target
                ? { ...req, isRespond: true, processedStatus: respond.status }
                : req
            )
          );
        } else {
          console.error("Failed to process request");
          showNotification(
            "Failed to process request. Please try again.",
            "error"
          );
        }
      } catch (error) {
        console.error("Error processing request:", error);
        showNotification(
          "Error processing request. Please try again.",
          "error"
        );
      }
    };

    if (respond) {
      handleRequestResponse();
    }
  }, [respond]);

  const handleJoinGroup = async (groupId) => {
    try {
      const response = await SendData("/api/v1/set/sendRequest", {
        target: groupId,
        type: 1,
      });

      if (response.ok) {
        showNotification("Join request sent successfully", "success");

        // Update the groups list to reflect the request was sent
        setGroups((prevGroups) =>
          prevGroups.map((group) => {
            return group.id === groupId
              ? { ...group, IsRequested: true }
              : group;
          })
        );
      } else {
        console.error("Failed to send join request");
        showNotification(
          "Failed to send join request. Please try again.",
          "error"
        );
      }
    } catch (error) {
      console.error("Error sending join request:", error);
      showNotification(
        "Error sending join request. Please try again.",
        "error"
      );
    }
  };

  return (
    <>
      <div className={Styles.groups}>
        <h1>Groups</h1>
        {groups?.length > 0 ? (
          groups.slice(0, 5).map((Group) => (
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
              {!Group?.IsRequested ? (
                <div>
                  <Image
                    onClick={() => handleJoinGroup(Group.id)}
                    src="/join.svg"
                    alt="join"
                    width={25}
                    height={25}
                    style={{
                      cursor: "pointer",
                      opacity: 1,
                    }}
                  />
                </div>
              ) : (
                <div className={Styles.requested}>
                  <span>Requested</span>
                </div>
              )}
            </div>
          ))
        ) : (
          <h3 style={{ textAlign: "center" }}>No groups to join</h3>
        )}
      </div>

      <div className={Styles.groups}>
        <h1>Group Requests</h1>
        {requests?.length > 0 ? (
          requests.map((request) => {
            return (
              <div
                key={`${request.sender_id}_${request.group_id}`}
                className={Styles.grouprequest}
              >
                <div>
                  <Image
                    src={request.avatar?.String || "/bannerBG.png"}
                    alt="profile"
                    width={40}
                    height={40}
                    style={{ borderRadius: "50%" }}
                  />
                  <div>
                    <h5>
                      {!request?.isRespond
                        ? request.message ||
                          `${request.username} wants to join ${request.group_name}`
                        : `You ${request.processedStatus}ed the request`}
                    </h5>
                    {request.username && request.group_name && (
                      <p style={{ fontSize: "12px", color: "#666" }}>
                        {request.username} → {request.group_name}
                      </p>
                    )}
                  </div>
                </div>
                {!request?.isRespond ? (
                  <div className={Styles.Buttons}>
                    <Image
                      onClick={() =>
                        setRespond({
                          sender: request.sender_id,
                          target: request.group_id,
                          type: 1,
                          status: "accept",
                          isSpecial: true,
                        })
                      }
                      src="/accept2.svg"
                      alt="accept"
                      width={25}
                      height={25}
                      style={{
                        marginRight: "10px",
                        cursor: "pointer",
                        opacity: 1,
                      }}
                    />
                    <Image
                      onClick={() =>
                        setRespond({
                          sender: request.sender_id,
                          target: request.group_id,
                          type: 1,
                          status: "refuse",
                        })
                      }
                      src="/decline.svg"
                      alt="decline"
                      width={25}
                      height={25}
                      style={{
                        cursor: "pointer",
                        opacity: 1,
                      }}
                    />
                  </div>
                ) : (
                  <div className={Styles.processedStatus}>
                    <span>Request {request.processedStatus}ed</span>
                  </div>
                )}
              </div>
            );
          })
        ) : (
          <h3 style={{ textAlign: "center" }}>No group requests</h3>
        )}
      </div>
    </>
  );
}
