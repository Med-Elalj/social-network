"use client";

import { useState } from "react";
import Style from "./chat.module.css";
import Image from "next/image";
import Users from "./[tab]/Users";
import Unread from "./[tab]/Unread";
import Groups from "./[tab]/Groups";

export default function Chat() {
    const [activeTab, setActiveTab] = useState("all");
    const [selectedUser, setSelectedUser] = useState(null); // ✅ FIXED HERE

    const handleTabClick = (selectedTab) => {
        setActiveTab(selectedTab);
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
                                <Image src={`/${selectedUser.avatar}` || "iconMale.png"} width={50} height={50} alt="userProfile" />
                                <div className={Style.userInfo}>
                                    <h5>{selectedUser.name}</h5>
                                    <h6
                                        style={{
                                            color: selectedUser.status === "online" ? "green" : "red",
                                        }}
                                    >
                                        {selectedUser.status || "offline"}
                                    </h6>
                                </div>
                            </div>

                            <div className={Style.body}>
                                {/* {selectedUser.role === 'sender' && ( */}
                                <div className={Style.user1}>
                                    <div>
                                        <p>{selectedUser.name}</p>
                                        <p>00.00 10-07-2002</p>
                                    </div>
                                    <p>Lorem ipsum dolor...</p>
                                </div>
                                {/* )} */}

                                {/* {selectedUser.role === 'receiver' && ( */}
                                <div className={Style.user2}>
                                    <div>
                                        <p>{selectedUser.name}</p>
                                        <p>00.00 10-07-2002</p>
                                    </div>
                                    <p>Lorem ipsum dolor...</p>
                                </div>
                                {/* )} */}

                                <div className={Style.user1}>
                                    <div>
                                        <p>{selectedUser.name}</p>
                                        <p>00.00 10-07-2002</p>
                                    </div>
                                    <p>Lorem ipsum dolor...</p>
                                </div>

                            </div>

                            <div className={Style.bottom}>
                                <Image src="upload.svg" width={30} height={30} alt="upload" />
                                <input type="text" name="message" id="message" />
                                <Image src="send.svg" width={30} height={30} alt="upload" />
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
