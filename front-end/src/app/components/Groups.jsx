import Image from "next/image";
import Styles from "../global.module.css";
import { useState, useEffect } from "react";
import { GetData, SendData } from "@/app/sendData.js";
import { useNotification } from "../context/notificationContext.jsx";

export default function Groups() {
  const [groups, setGroups] = useState([]);
  const [joinedGroupId, setJoinedGroupId] = useState(null);
  const [requests, setRequests] = useState([]);
  const [respond, setRespond] = useState(null);
  const [loading, setLoading] = useState(false);
  const [processedRequests, setProcessedRequests] = useState(new Set());
  const { showNotification } = useNotification();

  useEffect(() => {
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
  }, [showNotification]);

  useEffect(() => {
    const fetchGroupRequests = async () => {
      try {
        const response = await SendData("/api/v1/get/requests", { type: 1 });

        if (!response.ok) {
          console.error("Failed to fetch group requests");
          showNotification("Failed to load group requests", "error");
          return;
        }

        const data = await response.json();
        setRequests(data);
        console.log("requests groups:", data);
      } catch (err) {
        console.error("Error fetching group requests:", err);
        showNotification("Error loading group requests", "error");
      }
    };

    fetchGroupRequests();
  }, [showNotification]);

  useEffect(() => {
    if (!respond) return;

    const handleRequestResponse = async () => {
      setLoading(true);
      try {
        console.log("Processing request:", respond);
        const response = await SendData("/api/v1/set/acceptFollow", respond);
        
        if (response.ok) {
          showNotification(`Request ${respond.status}ed successfully`, "success");
          
          // Add to processed requests
          const requestKey = `${respond.sender}_${respond.target}`;
          setProcessedRequests(prev => new Set([...prev, requestKey]));
          
          // Remove the processed request from the requests list
          setRequests(prevRequests => 
            prevRequests.filter(req => 
              !(req.sender_id === respond.sender && req.group_id === respond.target)
            )
          );
        } else {
          console.error("Failed to process request");
          showNotification("Failed to process request. Please try again.", "error");
        }
      } catch (error) {
        console.error("Error processing request:", error);
        showNotification("Error processing request. Please try again.", "error");
      } finally {
        setLoading(false);
        setRespond(null);
      }
    };

    handleRequestResponse();
  }, [respond, showNotification]);

  const handleJoinGroup = async (groupId) => {
    if (loading) return;
    
    setLoading(true);
    try {
      const response = await SendData("/api/v1/set/sendRequest", { 
        target: groupId,
        type:1 
      });
      
      if (response.ok) {
        showNotification("Join request sent successfully", "success");
        setJoinedGroupId(groupId);
        
        // Update the groups list to reflect the request was sent
        setGroups(prevGroups => 
          prevGroups.map(group => 
            group.id === groupId ? { ...group, IsRequested: true } : group
          )
        );
      } else {
        console.error("Failed to send join request");
        showNotification("Failed to send join request. Please try again.", "error");
      }
    } catch (error) {
      console.error("Error sending join request:", error);
      showNotification("Error sending join request. Please try again.", "error");
    } finally {
      setLoading(false);
      setJoinedGroupId(null);
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
              <div>
                <Image
                  onClick={() => handleJoinGroup(Group.id)}
                  src="/join.svg"
                  alt="join"
                  width={25}
                  height={25}
                  style={{ 
                    cursor: loading ? "not-allowed" : "pointer",
                    opacity: loading ? 0.5 : 1
                  }}
                />
              </div>
            </div>
          ))
        ) : (
          <h3 style={{ textAlign: "center" }}>No groups to join</h3>
        )}
      </div>

      <div className={Styles.groups}>
        <h1>Groups Request</h1>
        {requests?.length > 0 ? (
          requests.map((Group) => {
            const requestKey = `${Group.sender_id}_${Group.group_id}`;
            const isProcessed = processedRequests.has(requestKey);
            
            return (
              <div key={`${Group.sender_id}_${Group.group_id}`} className={Styles.grouprequest}>
                <div>
                  <Image
                    src={Group.Avatar?.String || "/db.png"}
                    alt="profile"
                    width={40}
                    height={40}
                    style={{ borderRadius: "50%" }}
                  />
                  <h5>
                    {isProcessed 
                      ? `Request has been processed`
                      : `${Group?.username} sent you a join request to group ${Group?.group_name}`
                    }
                  </h5>
                </div>
                
                {!isProcessed && (
                  <div className={Styles.Buttons}>
                    <Image
                      onClick={() =>
                        setRespond({
                          sender: Group.sender_id,
                          target: Group.group_id,
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
                        cursor: loading ? "not-allowed" : "pointer",
                        opacity: loading ? 0.5 : 1
                      }}
                    />
                    <Image
                      onClick={() =>
                        setRespond({
                          sender: Group.sender_id,
                          target: Group.group_id,
                          type: 1,
                          status: "refuse",
                        })
                      }
                      src="/decline.svg"
                      alt="decline"
                      width={25}
                      height={25}
                      style={{ 
                        cursor: loading ? "not-allowed" : "pointer",
                        opacity: loading ? 0.5 : 1
                      }}
                    />
                  </div>
                )}
              </div>
            );
          })
        ) : (
          <h3 style={{ textAlign: "center" }}>No groups request</h3>
        )}
      </div>
    </>
  );
}