import { GetData } from "@/app/sendData.js";
import { useState, useEffect } from "react";
import Image from "next/image";
import Link from "next/link.js";
import Style from '../groupes/groups.module.css';

export default function Discover() {
    const [groups, setGroups] = useState([]);

    useEffect(() => {
        const fetchData = async () => {
            const formData = { userId: 1 };
            const response = await GetData("/api/v1/get/groupToJoin", formData);
            const body = await response.json();

            if (response.status !== 200) {
                console.error(body);
            } else {
                setGroups(body.groups);
                console.log("Groups fetched successfully!");
            }
        };

        fetchData();
    }, []);

    return (
        <div className={Style.groupGrid}>
            {groups && groups.map((Group, i) => (
                <div className={Style.groupCard} key={Group.ID}>
                    <Image
                        src={Group.Avatar?.String || "/bannerBG.png"}
                        alt="profile"
                        width={50}
                        height={50}
                        sizes="(max-width: 768px) 100vw, 250px"
                        className={Style.groupAvatar}
                    />
                    <h4>{Group.GroupName}</h4>
                    <p>{Group.Description}</p>
                    <Link href="/join" className={Style.acceptBtn}>Join Group</Link>
                </div>
            ))}
        </div>

    )
}