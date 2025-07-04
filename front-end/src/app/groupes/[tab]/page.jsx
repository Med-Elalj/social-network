"use client";

import Style from "../groups.module.css";
import Image from "next/image";
import Link from "next/link";
import { useState, useEffect } from "react";
import { useParams, usePathname, useRouter } from "next/navigation";
import GroupPosts from "./GroupPosts.jsx";
import Discover from "./Discover.jsx";
import YourGroups from "./YourGroups.jsx";
import CreateGroup from "./CreateGroup.jsx";
import { SendData } from "../../../../utils/sendData.js";

export default function Groupes() {
    const router = useRouter();
    const pathname = usePathname();
    const { tab } = useParams(); // e.g., "feed", "discover", etc.
    const [activeTab, setActiveTab] = useState(tab || "feed");

    const [GroupName, setGroupName] = useState("");
    const [image, setImage] = useState(null);
    const [privacy, setPrivacy] = useState("public");
    const [about, setAbout] = useState("");
    const [previewUrl, setPreviewUrl] = useState(null);

    const handleImageChange = (e) => {
        const file = e.target.files[0];
        if (file) {
            setImage(file);
            setPreviewUrl(URL.createObjectURL(file));
        }
    };

    const handleExit = (e) => {
        router.push('/groupes/feed')
    };

    const handleSubmit = async (e) => {
        e.preventDefault();

        const fetchData = async () => {
            const formData = {
                "groupName": GroupName,
                "privacy": privacy,
                "about": about,
                "avatar": image
            };

            const response = await SendData("/api/v1/set/GroupCreation", formData);
            const Body = await response.json();
            if (response.status !== 200) {
                console.log(Body);
            } else {
                router.push('/groupes/feed')
                console.log('Posts fetched successfully!');
            }
        };

        fetchData();
    }

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
                {pathname == "/groupes/create"
                    ?
                    <>
                        <div className={Style.CreateForm}>
                            <Image
                                src="/exit.svg"
                                width={40}
                                height={40}
                                alt="Exit"
                                style={{ cursor: "pointer" }}
                                onClick={handleExit}
                            />

                            <form className={Style.form} onSubmit={handleSubmit}>
                                <div style={{ marginBottom: "15px" }}>
                                    <label htmlFor="groupName">
                                        <span>Group Name</span>
                                        <input
                                            className={Style.input}
                                            type="text"
                                            name="groupName"
                                            id="groupName"
                                            required
                                            value={GroupName}
                                            onChange={(e) => setGroupName(e.target.value)}
                                        />
                                    </label>
                                </div>

                                <div style={{ marginBottom: "15px" }}>
                                    <label htmlFor="privacy">
                                        <span>Privacy</span>
                                        <select
                                            className={Style.input}
                                            name="privacy"
                                            id="privacy"
                                            required
                                            value={privacy}
                                            onChange={(e) => setPrivacy(e.target.value)
                                            }
                                        >
                                            <option className={Style.option} value="public">Public Group</option>
                                            <option value="private">Private Group</option>
                                        </select>
                                    </label>
                                </div>

                                <div style={{ marginBottom: "20px" }}>
                                    <label htmlFor="about">
                                        <span>About (optional)</span>
                                        <input
                                            className={Style.input}
                                            type="text"
                                            name="about"
                                            id="about"
                                            value={about}
                                            onChange={(e) => setAbout(e.target.value)} />
                                    </label>
                                </div>

                                <div style={{ marginBottom: "20px" }}>
                                    <label htmlFor="image" style={{ cursor: "pointer", color: "var(--third-color)" }}>
                                        <span>Avatar (optional)</span>
                                        <div style={{ marginTop: "8px", padding: "10px", border: "1px dashed var(--border-color)", borderRadius: "6px", display: "flex" }}>
                                            <Image src="/Image.svg" alt="Upload" width="24" height="24" />&nbsp;&nbsp;
                                            Click to choose image
                                        </div>
                                    </label>
                                    <input
                                        type="file"
                                        name="image"
                                        id="image"
                                        accept="image/*"
                                        style={{ display: "none" }}
                                        onChange={handleImageChange}
                                    />
                                </div>

                                <div className={Style.buttonWrapper}>
                                    <button type="submit" className={Style.submit}>Create</button>
                                </div>
                            </form>
                        </div>

                    </>
                    : <>
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
                                        <Image src="/iconMale.png" alt="profile" width={25} height={25} />
                                        <h5>Username</h5>
                                    </div>
                                    <div className={Style.Buttons}>
                                        <Link href="/accept">
                                            <Image src="/accept.svg" alt="accept" width={25} height={25} />
                                        </Link>
                                        <Link href="/reject">
                                            <Image src="/reject.svg" alt="reject" width={25} height={25} />
                                        </Link>
                                    </div>
                                </div>
                            ))}
                        </div>
                    </>
                }

            </div >

            <div className={Style.second}>
                {{
                    feed: <GroupPosts />,
                    discover: <Discover />,
                    groups: <YourGroups />,
                    create: (
                        <CreateGroup
                            groupName={GroupName}
                            privacy={privacy}
                            about={about}
                            imagePreview={previewUrl}
                        />
                    ),
                }[activeTab] || <p>Invalid tab</p>}
            </div>
        </div >
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
