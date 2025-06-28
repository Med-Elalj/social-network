"use client"

import Style from "./profile.module.css";
import Image from "next/image";
import Link from "next/link";
import { useState } from "react";
import Posts from "./[tab]/Posts";
import Following from "./[tab]/Following";
import Followers from "./[tab]/Followers";

export default function Profile() {
    const [activeSection, setActiveSection] = useState("posts"); // default section

    return (
        <div className={Style.container}>
            <div className={Style.header}>
                <Image src="/db.png" fill alt="cover" />
            </div>
            <div className={Style.body}>
                <div className={Style.first}>
                    <div className={Style.ProfileInfo}>
                        <div className={Style.top}>
                            <div style={{ position: "relative", width: "200px", height: "200px" }}>
                                <Image src="/db.png" alt="cover" fill />
                            </div>
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

                        <button type="button" style={{ backgroundColor: "red", border: "1px solid red" }}>
                            Delete profile
                        </button>

                        <div className={Style.numbers}>
                            <span onClick={() => setActiveSection("posts")}>
                                <h4>Posts</h4>
                                <h5>0</h5>
                            </span>

                            <span onClick={() => setActiveSection("followers")}>
                                <h4>Followers</h4>
                                <h5>0</h5>
                            </span>

                            <span
                                onClick={() => setActiveSection("following")}>
                                <h4>Following</h4>
                                <h5>0</h5>
                            </span>

                        </div>
                    </div>
                </div>

                <div className={Style.second}>
                    {activeSection === "posts" && (
                        <Posts />
                    )}

                    {activeSection === "followers" && (
                        <Followers />
                    )}

                    {activeSection === "following" && (
                        <Following />
                    )}
                </div>


                <div className={Style.end}>
                    <div className={Style.requists}>
                        <h3>Suggestion</h3>
                        <div>
                            <div>
                                <Image src="/db.png" alt="profile" width={40} height={40} style={{ borderRadius: '50%' }} />
                                <h5>username</h5>
                            </div>
                            <Link href="/addUser"><Image src="/addUser.svg" alt="profile" width={25} height={25} /></Link>
                        </div>
                        <div>
                            <div>
                                <Image src="/db.png" alt="profile" width={40} height={40} style={{ borderRadius: '50%' }} />
                                <h5>username</h5>
                            </div>
                            <Link href="/addUser"><Image src="/addUser.svg" alt="profile" width={25} height={25} /></Link>
                        </div>
                        <div>
                            <div>
                                <Image src="/db.png" alt="profile" width={40} height={40} style={{ borderRadius: '50%' }} />
                                <h5>username</h5>
                            </div>
                            <Link href="/addUser"><Image src="/addUser.svg" alt="profile" width={25} height={25} /></Link>
                        </div>
                        <div>
                            <div>
                                <Image src="/db.png" alt="profile" width={40} height={40} style={{ borderRadius: '50%' }} />
                                <h5>username</h5>
                            </div>
                            <Link href="/addUser"><Image src="/addUser.svg" alt="profile" width={25} height={25} /></Link>
                        </div>
                        <div>
                            <div>
                                <Image src="/db.png" alt="profile" width={40} height={40} style={{ borderRadius: '50%' }} />
                                <h5>username</h5>
                            </div>
                            <Link href="/addUser"><Image src="/addUser.svg" alt="profile" width={25} height={25} /></Link>
                        </div>
                    </div>

                    <div className={Style.requists}>
                        <div>
                            <div>
                                <Image src="/iconMale.png" alt="profile" width={40} height={40} />
                                <h5>Username</h5>
                            </div>
                            <div className={Style.Buttons}>
                                <Link href="/accept"><Image src="/accept.svg" alt="profile" width={25} height={25} /></Link>
                                <Link href="/reject"><Image src="/reject.svg" alt="profile" width={25} height={25} /></Link>
                            </div>
                        </div>
                        <div>
                            <div>
                                <Image src="/iconMale.png" alt="profile" width={40} height={40} />
                                <h5>Username</h5>
                            </div>
                            <div className={Style.Buttons}>
                                <Link href="/accept"><Image src="/accept.svg" alt="profile" width={25} height={25} /></Link>
                                <Link href="/reject"><Image src="/reject.svg" alt="profile" width={25} height={25} /></Link>
                            </div>
                        </div>
                        <div>
                            <div>
                                <Image src="/iconMale.png" alt="profile" width={40} height={40} />
                                <h5>Username</h5>
                            </div>
                            <div className={Style.Buttons}>
                                <Link href="/accept"><Image src="/accept.svg" alt="profile" width={25} height={25} /></Link>
                                <Link href="/reject"><Image src="/reject.svg" alt="profile" width={25} height={25} /></Link>
                            </div>
                        </div>
                        <div>
                            <div>
                                <Image src="/iconMale.png" alt="profile" width={40} height={40} />
                                <h5>Username</h5>
                            </div>
                            <div className={Style.Buttons}>
                                <Link href="/accept"><Image src="/accept.svg" alt="profile" width={25} height={25} /></Link>
                                <Link href="/reject"><Image src="/reject.svg" alt="profile" width={25} height={25} /></Link>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    )
}