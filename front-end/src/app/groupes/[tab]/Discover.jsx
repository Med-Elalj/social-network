import { GetData } from "../../../../utils/sendData.js";
import { useState, useEffect } from "react";
import Image from "next/image";
import Link from "next/link.js";
import Style from '../groups.module.css';

export default function Discover() {
    const [groups, setGroups] = useState([]);

    useEffect(() => {
        const fetchData = async () => {
            const response = await GetData(process.env.NEXT_PUBLIC_API_URL + "/get/groupToJoin");
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
        <div className={groups ? Style.groupGrid : Style.noPosts}>
            {groups ? (
                groups.map((Group) => (
                    <div className={Style.groupCard} key={Group.ID}>
                        <Image
                            src={Group.Avatar?.String || "/db.png"}
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
                ))
            ) : (
                <>
                    <h3>Join groups to see feeds</h3>
                    <Link href="/groupes/create">Create a group</Link>
                </>
            )}
        </div>
    );
}
