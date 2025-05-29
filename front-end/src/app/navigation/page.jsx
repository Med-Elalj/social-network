import Link from "next/link";
import Styles from "./nav.module.css"

export default function Routing() {
    return (
        <div className={Styles.nav}>
            <div>
                <Link href="/">Social Network</Link>
            </div>
            <div>
                <Link href="/">Dashboard</Link>
            </div>
            <div className={Styles.auth}>
                <Link href="/auth/login">Login</Link>
                <Link href="/auth/register">Register</Link>
                <Link href="">Logout</Link>
            </div>
        </div>
    );
}
