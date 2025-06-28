import Image from "next/image";
import Style from "../profile.module.css";

export default function Followers() {
    return (
        <div className={Style.followList}>
            <h1>Followers</h1>
            <div className={Style.NewUser}>
                <div>
                    <Image src="/db.png" alt="profile" width={40} height={40} style={{ borderRadius: '50%' }} />
                    <h5>username</h5>
                </div>
                {/* <Link href="/join"><Image src="/addUser.svg" alt="profile" width={25} height={25} /></Link> */}
            </div>

            <div className={Style.NewUser}>
                <div>
                    <Image src="/db.png" alt="profile" width={40} height={40} style={{ borderRadius: '50%' }} />
                    <h5>username</h5>
                </div>
                {/* <Link href="/join"><Image src="/addUser.svg" alt="profile" width={25} height={25} /></Link> */}
            </div>

            <div className={Style.NewUser}>
                <div>
                    <Image src="/db.png" alt="profile" width={40} height={40} style={{ borderRadius: '50%' }} />
                    <h5>username</h5>
                </div>
                {/* <Link href="/join"><Image src="/addUser.svg" alt="profile" width={25} height={25} /></Link> */}
            </div>

            <div className={Style.NewUser}>
                <div>
                    <Image src="/db.png" alt="profile" width={40} height={40} style={{ borderRadius: '50%' }} />
                    <h5>username</h5>
                </div>
                {/* <Link href="/join"><Image src="/addUser.svg" alt="profile" width={25} height={25} /></Link> */}
            </div>
        </div>
    )
}