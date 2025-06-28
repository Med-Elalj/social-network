import Style from "./profile.module.css";
import Image from "next/image";

export default function Profile() {
    return (
        <div className={Style.container}>
            <div className={Style.header}>
                <Image src="/db.png" fill alt="cover" />
            </div>
            <div className={Style.body}>
                <div className={Style.first}>
                    <div className={Style.ProfileInfo}>
                        <div className={Style.top}>
                            <Image src="/db.png" alt="cover" width={0} height={0} />
                            <h4>User name</h4>
                        </div>

                        <div className={Style.center}>
                            <span>
                                <h5>User name :</h5>&nbsp;&nbsp;<h5>Lorem ipsum dolor sit amet.</h5>
                            </span>
                            <span>
                                <h5>User name :</h5>&nbsp;&nbsp;<h5>Lorem ipsum dolor sit amet.</h5>
                            </span>
                            <span>
                                <h5>User name :</h5>&nbsp;&nbsp;<h5>Lorem ipsum dolor sit amet.</h5>
                            </span>
                        </div>

                        <button type="button">
                            Update profile
                        </button>

                        <button type="button">
                            Delete profile
                        </button>
                    </div>
                </div>

                <div className={Style.second}></div>
                <div className={Style.end}></div>
            </div>
        </div>
    )
}