import Style from "../../profile.module.css";
import Image from "next/image";
import { useEffect, useState } from "react";
import { GetData } from "../../../../../utils/sendData";

export default function Following(userId) {
  const [users, setUsers] = useState(null);

  useEffect(() => {
    const fetchProfile = async () => {
      try {
        const res = await GetData(`/api/v1/following?userId=${userId.userId}`);
        if (res.ok) {
          const data = await res.json();
          setUsers(data);
        }
      } catch (err) {
        console.error("Error fetching profile:", err);
      }
    };

    if (!users) {
      fetchProfile();
    }
  }, [userId]);
  return (
    <div className={Style.followList}>
      <h1>Following</h1>
      {users && users.map((user, i) => (
        <div key={i} className={Style.NewUser}>
          <div>
            <Image
              src={user.pfp.Valid ? user.pfp.String : "/iconMale.png"}
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
