"use client";

import Style from "../groups.module.css";
import Image from "next/image";
import Link from "next/link";
import { useState, useEffect } from "react";
import { useParams, useRouter } from "next/navigation";
import GroupPosts from "./GroupPosts.jsx";
import Discover from "../../components/Discover.jsx";
import YourGroups from "./YourGroups.jsx";

export default function Groupes() {
    const router = useRouter();
    const { tab } = useParams(); // e.g., "feed", "discover", etc.
    const [activeTab, setActiveTab] = useState(tab || "feed");

    useEffect(() => {
        setActiveTab(tab || "feed");
    }, [tab]);

    const handleTabClick = (selectedTab) => {
        if (selectedTab !== activeTab) {
            router.push(`/groupes/${selectedTab}`);
        }
    };

    return (
        <div className={Style.content}>
            <div className={Style.first}>
                <div className={Style.select}>
                    <Tab name="feed" icon="postGroups" activeTab={activeTab} onClick={handleTabClick} />
                    <Tab name="discover" icon="discover" activeTab={activeTab} onClick={handleTabClick} />
                    <Tab name="groups" icon="groupe" activeTab={activeTab} onClick={handleTabClick} />
                    <Tab name="create" icon="create" activeTab={activeTab} onClick={handleTabClick} />
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
                {{
                    feed: <GroupPosts />,
                    discover: <Discover />,
                    groups: <YourGroups />,
                    create: <p>Create Group</p>,
                }[activeTab] || <p>Invalid tab</p>}
            </div>
        </div>
    );
}

function Tab({ name, icon, activeTab, onClick }) {
    return (
        <div
            className={activeTab === name ? Style.active : ""}
            onClick={() => onClick(name)}
        >
            <Image src={`/${icon}.svg`} alt={name} width={30} height={30} />
            <p>{name.charAt(0).toUpperCase() + name.slice(1)}</p>
        </div>
    );
}
