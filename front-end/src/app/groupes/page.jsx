"use client";

import Style from "./groups.module.css";
import Image from "next/image";
import Link from "next/link";
import { useState } from "react";
import GroupPosts from "../components/GroupPosts.jsx";
import Discover from "../components/Discover.jsx";
import YourGroups from "../components/YourGroups.jsx";

export default function Groupes() {
    const [activeTab, setActiveTab] = useState("feed");

    return (
        <div className={Style.content}>
            <div className={Style.first}>
                <div className={Style.select}>
                    <div
                        className={activeTab === "feed" ? Style.active : ""}
                        onClick={() => setActiveTab("feed")}
                    >
                        <Image src="/postGroups.svg" alt="profile" width={30} height={30} />
                        <p>Your feed</p>
                    </div>
                    <div
                        className={activeTab === "discover" ? Style.active : ""}
                        onClick={() => setActiveTab("discover")}
                    >
                        <Image src="/discover.svg" alt="profile" width={30} height={30} />
                        <p>Discover</p>
                    </div>
                    <div
                        className={activeTab === "groups" ? Style.active : ""}
                        onClick={() => setActiveTab("groups")}
                    >
                        <Image src="/groupe.svg" alt="profile" width={30} height={30} />
                        <p>Your groups</p>
                    </div>
                </div>
                <div className={Style.Requiests}>
                    <h1>Groups requests</h1>
                    {[1, 2, 3].map((_, i) => (
                        <div key={i} className={Style.RequestItem}>
                            <div>
                                <Image src="/iconMale.png" alt="profile" width={40} height={40} />
                                <h5>Username</h5>
                            </div>
                            <div className={Style.Buttons}>
                                <Link href="/accept">Accept</Link>
                                <Link href="/reject">Reject</Link>
                            </div>
                        </div>
                    ))}
                </div>
            </div>
            <div className={Style.second}>
                {(() => {
                    switch (activeTab) {
                        case "feed":
                            return <GroupPosts />;
                        case "discover":
                            return <Discover />;
                        case "groups":
                            return <YourGroups />;
                        default:
                            return <p>Profile Group</p>;
                    }
                })()}
            </div>
        </div>
    )
}
