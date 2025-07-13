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
            {notifications?.length > 0 && (
                notifications.map((notification) => (
                    <Link key={notification.Id} href={`/profile/${notification.Type == 0 ? notification.Username : `${notification.GroupName}${notification.Type == 1 ? "" : `/${notification.ID}`}`}`} onClick={() => setIsOpen(false)}>
                        {notification.Message}
                    </Link>
                ))
            )}
        </div>
    )
}