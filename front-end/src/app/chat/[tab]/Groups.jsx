import { useState } from "react";
import Style from "../chat.module.css";
import Image from "next/image";

export default function Groups({ onUserSelect }) {
    const [activeIndex, setActiveIndex] = useState(null);

    const groups = [
        { name: "Group 1" },
        { name: "Group 2" },
        { name: "Group 3" },
    ];

    const handleUserClick = (group, index) => {
        setActiveIndex(index);
        onUserSelect(group);
    };

    return (
        <div className={Style.users}>
            {groups && groups.length > 0 ? (
                <>
                    {groups.map((group, index) => (
                        <div
                            key={index}
                            className={`${Style.user} ${activeIndex === index ? Style.active : ""}`}
                            onClick={() => handleUserClick(group,index)}
                        >
                            <div className={Style.userImageWrapper}>
                                <Image src="/iconMale.png" width={50} height={50} alt="userProfile" />
                            </div>
                            <h4>{group.name}</h4>
                        </div>
                    ))}
                </>
            ) : (
                <>
                    <div className={Style.message}>
                        <h4>No messages</h4>
                        <h5>When you have groups chats, you’ll see them here.</h5>
                    </div>
                </>
            )}
        </div>
    );
}
