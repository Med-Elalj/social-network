import Styles from "./nav.module.css";
import Link from "next/link";
import { useWebSocket } from "@/app/context/WebSocketContext.jsx";
import { useEffect } from "react";

export default function NotificationList({ notifications, setIsOpen }) {
    const { newNotification } = useWebSocket();

    useEffect(() => {
        if (newNotification) {
            notifications = [...notifications, newNotification]
        }
    }, [newNotification])

    return (
        <div className={`${Styles.dropdownMenu} ${Styles.notification}`}>
            {notifications?.length > 0 ?  (
                notifications.map((notification) => (
                    <Link key={notification.id} href={notification.type==0? `/profile/${notification.username}` : notification.type== 1 ? `/groupes/profile/${notification.group_name}` : ""} onClick={() => setIsOpen(false)}>
                        <span style={{color:"var(--third-color)"}}>{notification.type==0? "Follow request: ": notification.type==1? "Join group: " : "New event: "}</span>{notification.message}
                    </Link>
                ))
            ) : <h3>you don't have any notification</h3>}

        </div>
    )
}