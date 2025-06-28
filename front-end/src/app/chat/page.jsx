"use client";

import { useState } from "react";
import Style from "./chat.module.css";
import Image from "next/image";
import { useEffect } from "react";
import Users from "./[tab]/Users";
import Unread from "./[tab]/Unread";
import Groups from "./[tab]/Groups";

export default function Chat() {
    const [activeTab, setActiveTab] = useState("all");
    const [selectedUser, setSelectedUser] = useState(null);
    const [image, setImage] = useState(null);
    // const [previewUrl, setPreviewUrl] = useState(null);
    const [content, setContent] = useState("");

    useEffect(() => {
        setContent("");
        setImage(null);
        // setPreviewUrl(null);
    }, [selectedUser]);


    const handleTabClick = (selectedTab) => {
        setActiveTab(selectedTab);
    };

    const handleMediaChange = (e) => {
        const file = e.target.files[0];
        if (file) {
            setImage(file);
            // setPreviewUrl(URL.createObjectURL(file));
        }
    };

    const handleSend = async (e) => {

        console.log(image);
        console.log(content);

    };

    return (
        <div className={Style.container}>
            <div className={Style.first}>
                <div className={Style.header}>
                    <h1>Chats</h1>
                    <Image src="/newMessage.svg" width={35} height={35} alt="newMessage" />
                </div>

                <div className={Style.select}>
                    <Tab name="all" icon="messages" activeTab={activeTab} onClick={handleTabClick} />
                    <Tab name="unread" icon="unread" activeTab={activeTab} onClick={handleTabClick} />
                    <Tab name="groups" icon="groupe" activeTab={activeTab} onClick={handleTabClick} />
                </div>

                {{
                    all: <Users onUserSelect={setSelectedUser} />,
                    unread: <Unread onUserSelect={setSelectedUser} />,
                    groups: <Groups onUserSelect={setSelectedUser} />
                }[activeTab] || <p>Invalid tab</p>}
            </div>

            <div className={Style.second}>
                <div className={Style.chat}>
                    {selectedUser ? (
                        <>
                            <div className={Style.top}>
                                <Image src={`/${selectedUser.avatar ?? "iconMale.png"}`} width={50} height={50} alt="userProfile" />
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
                                {/* {selectedUser.role === 'sender' && ( */}
                                <div className={Style.user1}>
                                    <p>{selectedUser.name}</p>
                                    <p>Lorem ipsum dolor...</p>
                                    <p>00.00 10-07-2002</p>
                                </div>
                                {/* )} */}

                                {/* {selectedUser.role === 'receiver' && ( */}
                                <div className={Style.user2}>
                                    <p>{selectedUser.name}</p>
                                    <p>Lorem ipsum dolor...</p>
                                    <p>00.00 10-07-2002</p>
                                </div>
                                {/* )} */}

                                <div className={Style.user1}>
                                    <p>{selectedUser.name}</p>
                                    <p>Lorem ipsum dolor...</p>
                                    <p>00.00 10-07-2002</p>
                                </div>

                            </div>

                            <div className={Style.bottom}>
                                <div>
                                    <label htmlFor="media" style={{ cursor: "pointer" }}>
                                        <Image src="upload.svg" width={30} height={30} alt="upload" />
                                    </label>
                                    <input
                                        type="file"
                                        name="media"
                                        style={{ display: "none" }}
                                        id="media"
                                        accept="image/*,video/*"
                                        onChange={handleMediaChange}
                                    />
                                </div>

                                <input type="text" name="message" id="message" value={content} onChange={(e) => setContent(e.target.value)} />

                                <Image
                                    src="send.svg"
                                    width={30}
                                    height={30}
                                    alt="send"
                                    onClick={handleSend}
                                    style={{ cursor: "pointer", marginRight: "3%" }}
                                />
                            </div>

                        </>
                    ) : (
                        <div className={Style.emptyChat}>
                            {activeTab == "groups"
                                ?
                                <h1>Select a group to start chat</h1>
                                :
                                <h1>Select a user to start chat</h1>
                            }
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
            <Image src={`/${icon}.svg`} alt={name} width={30} height={30} />
            <p>{name}</p>
        </div>
    );
}
