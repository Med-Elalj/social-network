import Image from "next/image";
import Style from "../../profile.module.css";
import { GetData } from "@/app/sendData.js";
import { useEffect, useState } from "react";

export default function Followers({ userId }) {
    const [users, setUsers] = useState(null);

    useEffect(() => {
        const fetchProfile = async () => {
            try {
                const res = await GetData(`/api/v1/followers?userId=${userId}`);
                if (res.ok) {
                    const data = await res.json();
                    setUsers(data);
                } else {
                    console.error("Failed to fetch followers");
                }
            } catch (err) {
                console.error("Error fetching followers:", err);
            }
        };

        if (userId) {
            fetchProfile(); // ✅ only run if userId exists
        }
    }, [userId]); // ✅ runs only when userId changes

    return (
        <div className={Style.followList}>
            <h1>Followers</h1>
            {users && users?.map((user, i) => (
                <div key={i} className={Style.NewUser}>
                    <div>
                        <Image
                            src={user.pfp?.Valid ? user.pfp.String : "/iconMale.png"}
                            alt="profile"
                            width={40}
                            height={40}
                            style={{ borderRadius: "50%" }}
                        />
                        <h5>{user.name}</h5>
                    </div>
                </div>
            ))}
        </div>
    );
}
