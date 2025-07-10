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
            {notifications.length > 0 && (
                notifications.map((notification) => (
                    <Link key={notification.id} href={`/`} onClick={() => setIsOpen(false)}>
                        {notification.message}
                    </Link>
                ))
            )}

        </div>
    )
}