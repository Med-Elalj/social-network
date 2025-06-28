import { useState } from "react";
import Style from "../chat.module.css";
import Image from "next/image";

export default function Users({ onUserSelect }) {
    const [activeIndex, setActiveIndex] = useState(null);

    const users = [
        { name: "User 1", message: "Lorem ipsum, dolor sit amet consectetu", status:"online" },
        { name: "User 2", message: "Lorem ipsum, dolor sit amet consectetu", status:"oflline" },
        { name: "User 3", message: "Lorem ipsum, dolor sit amet consectetu", status:"oflline" },
    ];

    const handleUserClick = (user, index) => {
        setActiveIndex(index);
        onUserSelect(user);
    };

    return (
        <div className={Style.users}>
            {users && users.length > 0 ? (
                <>
                    {users.map((user, index) => (
                        <div
                            key={index}
                            className={`${Style.user} ${activeIndex === index ? Style.active : ""}`}
                            onClick={() => handleUserClick(user,index)}
                        >
                            <div className={Style.userImageWrapper}>
                                <Image src={`/${user.avatar ?? "iconMale.png"}`} width={50} height={50} alt="userProfile" />
                                {activeIndex === index && <span className={Style.activeIndicator} />}
                            </div>
                            <div>
                                <h4>{user.name}</h4>
                                <p>{user.message}</p>
                            </div>
                            {activeIndex === index && (
                                <Image className={Style.last} src={`/${user.avatar ?? "iconMale.png"}`} width={20} height={20} alt="userProfile" />
                            )}
                        </div>
                    ))}
                </>
            ) : (
                <h5>No messages</h5>
            )}
        </div>
    );
}
