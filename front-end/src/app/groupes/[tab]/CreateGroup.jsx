import Image from "next/image";
import Style from "../groups.module.css";

export default function CreateGroup({ groupName, privacy, about, imagePreview }) {
    return (
        <div className={Style.container}>
            <h2>Group Preview</h2>
            <div className={Style.header}>
                <Image
                    src={imagePreview || "/groupsBg.png"}
                    alt="Group Avatar"
                    fill
                    style={{ objectFit: 'inherit' }}
                />
            </div>


            <div className={Style.preview}>
                <div className={Style.name}>
                    <h2>{groupName || "Group name"}</h2>
                    <div className={Style.privacy}>
                        <Image
                            src={`/${privacy}.svg`}
                            alt="privacy"
                            width={20}
                            height={20}
                        />
                        <p>&nbsp;</p>
                        <p>{privacy}</p>
                        <p>&nbsp; - &nbsp;</p>
                        <p style={{ fontWeight: 'bold' }}>1 member</p>
                    </div>
                </div>

                <div className={Style.navPreview}>
                    <h5>About</h5>
                    <h5>Posts</h5>
                    <h5>Members</h5>
                    <h5>Events</h5>
                </div>
            </div>

            <div className={Style.example}>
                <h4>{about}</h4>
            </div>
        </div>
    )
}