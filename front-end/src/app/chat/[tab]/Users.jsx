
import { useEffect, useState } from "react";
import Style from "../chat.module.css";
import Image from "next/image";
// import { useSearchParams } from 'next/navigation';


export default function Users({ users, onUserSelect }) {
    const [activeIndex, setActiveIndex] = useState(null);
    // const searchParams = useSearchParams();
    // const goTo = searchParams.get('goTo');


    // useEffect(() => {
    //     if (goTo && users?.length > 0) {

    //         users.forEach((element, idx) => {
    //             if (element.name == goTo) {
    //                 setActiveIndex(idx);
    //                 onUserSelect(element);
    //             }
    //         });
    //     }
    // }, [users]);

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
                            key={user.id}
                            className={`${Style.user} ${activeIndex === index ? Style.active : ""}`}
                            onClick={() => handleUserClick(user, index)}
                        >
                            <div className={Style.userImageWrapper}>
                                <Image
                                    src={`${user?.pfp?.Valid ? user.pfp.String : "/iconMale.png"}`}
                                    width={50}
                                    height={50}
                                    alt="userProfile"
                                    style={{borderRadius:"50%"}}
                                />
                                {user.online && <span className={Style.activeIndicator} />}
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
                <>
                    <div className={Style.message}>
                        <h4>No messages</h4>
                        <h5>When you have chats, you’ll see them here.</h5>
                    </div>
                </>
            )
            }
        </div >
    );
}
