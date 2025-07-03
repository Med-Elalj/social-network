"use client";

import Image from "next/image";
import Styles from "./global.module.css";

export default function Comments({ Post, onClose }) {
    return (
        <div className={Styles.commentPopup}>
            <button className={Styles.closeBtn} onClick={onClose}>
                <Image src="/exit.svg" alt="exit" width={30} height={30} />
            </button>

            <section className={Styles.userinfo}>
                {/* post info */}
                <div className={Styles.user}>
                    {/* Main Avatar */}
                    <Image
                        src={Post.AvatarGroup?.string ? `/${Post.AvatarGroup.String}` : '/iconMale.png'}
                        alt="avatar"
                        width={25}
                        height={25}
                    />

                    {/* Texts block */}
                    <div className={Styles.texts}>
                        {Post.GroupId?.Valid ? (
                            <>
                                <p>{Post.GroupName.String}</p>
                                <div className={Styles.user}>
                                    <Image
                                        src={Post.AvatarUser?.String ? `/${Post.Avatar.String}` : '/iconMale.png'}
                                        alt="avatar"
                                        width={20}
                                        height={20}
                                    />
                                    <p>{Post.UserName}</p>
                                </div>
                            </>
                        ) : (
                            <>
                                <p>{Post.UserName}</p>
                                <div className={Styles.user}></div>
                            </>
                        )}
                    </div>
                </div>
                {/* Add your comments list and form here */}
            </section>
        </div>
    );
}
