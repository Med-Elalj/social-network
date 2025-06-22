import Image from "next/image";
import Styles from "../global.module.css";
import Link from "next/link";

export default function Friends() {
    return (
        <>
            <div className={Styles.Requiests}>
                <h1>Friend requests</h1>
                {[1, 2, 3].map((_, i) => (
                    <div key={i}>
                        <div>
                            <Image src="/iconMale.png" alt="profile" width={40} height={40} />
                            <h5>Username</h5>
                        </div>
                        <div className={Styles.Buttons}>
                            <Link href="/accept">Accept</Link>
                            <Link href="/reject">Reject</Link>
                        </div>
                    </div>
                ))}
            </div>

            <div className={Styles.friends}>
                <h1>Contacts</h1>
                {[1, 2, 3, 4].map((_, i) => (
                    <div key={i}>
                        <div>
                            <Image src="/iconMale.png" alt="profile" width={40} height={40} />
                            <h5>User Name</h5>
                        </div>
                        <p>online</p>
                    </div>
                ))}
                <div>
                    <div>
                        <Image src="/iconMale.png" alt="profile" width={40} height={40} />
                        <h5>User Name</h5>
                    </div>
                    <p style={{ color: 'red' }}>offline</p>
                </div>
            </div>
        </>

    )
}