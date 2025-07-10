"use client";

import { useState } from "react";
import Style from "./chat.module.css";
import Image from "next/image";
import { useEffect } from "react";
import Users from "./[tab]/Users";
import Groups from "./[tab]/Groups";
import Messages from "./messages.jsx";
import ChatInput from "./input.jsx";
import { useWebSocket } from "../context/WebSocketContext.jsx";
const BACKEND_URL =
  process.env.NEXT_PUBLIC_BACKEND_URL || "https://localhost:8080";

export default function Chat() {
  const [activeTab, setActiveTab] = useState("all");
  const [selectedUser, setSelectedUser] = useState(null);
  const [image, setImage] = useState(null);
  // const [previewUrl, setPreviewUrl] = useState(null);
  const [content, setContent] = useState("");
  const [newMessage, setNewMessage] = useState(null)
  const [personalDiscussions, setPersonalDiscussions] = useState([]);
  const [groupDiscussions, setGroupDiscussions] = useState([]);
  const { setTarget, updateOnlineUser } = useWebSocket()



  useEffect(() => {
    const fetchConversations = async () => {
      try {
        console.log("Resolved backend URL:", BACKEND_URL);
        const response = await fetch(BACKEND_URL + "/api/v1/get/users", {
          method: "POST",
          credentials: "include",
          headers: {
            "Content-Type": "application/json",
          },
        });

        if (!response.ok) {
          throw new Error(`HTTP error! Status: ${response.status}`);
        }

        const data = await response.json();
        // Filter data by profile.isgroup
        const personal = (data || []).filter(
          (profile) => profile.is_group === false
        );
        const groups = (data || []).filter(
          (profile) => profile.is_group === true
        );
        setPersonalDiscussions(personal);
        setGroupDiscussions(groups);
      } catch (error) {
        console.error("Error fetching conversations:", error);
      }
    };

    fetchConversations();
  }, []);

  useEffect(() => {
    if (updateOnlineUser && updateOnlineUser.uid) {
      setPersonalDiscussions((prev) =>
        prev.map((user) =>
          user.id === updateOnlineUser.uid
            ? { ...user, online: updateOnlineUser.value }
            : user
        )
      );
    }
  }, [updateOnlineUser]);


  useEffect(() => {
    if (selectedUser) setTarget(selectedUser.id);

    setContent("");
    setImage(null);
    // setPreviewUrl(null);
  }, [selectedUser]);

  const handleTabClick = (selectedTab) => {
    console.log("selected tab: ", selectedTab);
    setActiveTab(selectedTab);
  };

  const handleMediaChange = (e) => {
    const file = e.target.files[0];
    if (file) {
      setImage(file);
      // setPreviewUrl(URL.createObjectURL(file));
    }
  };

  return (
    <div className={Style.container}>
      <div className={Style.first}>
        <div className={Style.header}>
          <h1>Chats</h1>
          <Image
            src="/newMessage.svg"
            width={25}
            height={25}
            alt="newMessage"
          />
        </div>

        <div className={Style.select}>
          <Tab
            name="all"
            icon="messages"
            activeTab={activeTab}
            onClick={handleTabClick}
          />
          <Tab
            name="groups"
            icon="groupe"
            activeTab={activeTab}
            onClick={handleTabClick}
          />
        </div>

        {{
          all: (
            <Users users={personalDiscussions} onUserSelect={setSelectedUser} />
          ),
          groups: (
            <Groups groups={groupDiscussions} onUserSelect={setSelectedUser} />
          ),
        }[activeTab] || <p>Invalid tab</p>}
      </div>

      <div className={Style.second}>
        <div className={Style.chat}>
          {selectedUser ? (
            <>
              <div className={Style.top}>
                <Image
                  src={`/${selectedUser.pfp.String ?? "iconMale.png"}`}
                  width={50}
                  height={50}
                  alt="userProfile"
                />
                <div className={Style.userInfo}>
                  <h5>{selectedUser.name}</h5>
                  <h6
                    style={{
                      color: selectedUser.status === "online" ? "green" : "red",
                    }}
                  >
                    {selectedUser.status}
                  </h6>
                </div>
              </div>

              <div className={Style.body}>
                <Messages user={selectedUser} />
              </div>

              <div className={Style.bottom}>


                {/* <input type="text" name="message" id="message" value={content} onChange={(e) => setContent(e.target.value)} />

                                <Image
                                    src="send.svg"
                                    width={25}
                                    height={25}
                                    alt="send"
                                    onClick={handleSend}
                                    style={{ cursor: "pointer", marginRight: "6%" }}
                                /> */}
                <ChatInput target={selectedUser.id} />
              </div>
            </>
          ) : (
            <div className={Style.emptyChat}>
              {activeTab == "groups" ? (
                <h1>Select a group to start chat</h1>
              ) : (
                <h1>Select a user to start chat</h1>
              )}
            </div>
          )}
        </div>
      </div>
    </div>
  );
}

function Tab({ name, icon, activeTab, onClick }) {
  return (
    <div
      className={`${Style.tab} ${activeTab === name ? Style.active : ""}`}
      onClick={() => onClick(name)}
    >
      <Image src={`/${icon}.svg`} alt={name} width={25} height={25} />
      <p>{name}</p>
    </div>
  );
}


