"use client";

import { usePathname } from 'next/navigation';
import { useState } from 'react';
import Link from "next/link";
import Image from 'next/image';
import Styles from "./nav.module.css";

export default function Routing() {
    const [isOpen, setIsOpen] = useState(false);
    const isLoggedIn = true;
    const nickname = "nickname";
    const pathname = usePathname();

    return (
        <div>
            <div className={Styles.nav}>
                {/* Left - Logo */}
                <div className={Styles.leftSection}>
                    <Link className={Styles.loginTitle} href="/">Social Network</Link>
                </div>

                {/* Center - Navigation Links */}
                <div className={`${Styles.centerSection} `}>
                    {isLoggedIn && (
                        <>
                            <Link className={`${Styles.linkWithIcon} ${pathname === "/" ? Styles.active : ""}`} href="/">
                                <span className={Styles.iconWrapper}>
                                    <Image src="/home2.svg" alt="home" width={25} height={25} className={Styles.iconDefault} />
                                    <Image src="/home.svg" alt="home-hover" width={25} height={25} className={Styles.iconHover} />
                                </span>
                            </Link>
                            <Link className={`${Styles.linkWithIcon} ${pathname === "/posts" ? Styles.active : ""}`} href="/posts">
                                <span className={Styles.iconWrapper}>
                                    <Image src="/posts2.svg" alt="posts" width={25} height={25} className={Styles.iconDefault} />
                                    <Image src="/posts.svg" alt="posts-hover" width={25} height={25} className={Styles.iconHover} />
                                </span>
                            </Link>
                            <Link className={`${Styles.linkWithIcon} ${pathname === "/groups" ? Styles.active : ""}`} href="/groups">
                                <span className={Styles.iconWrapper}>
                                    <Image src="/groupe2.svg" alt="groups" width={25} height={25} className={Styles.iconDefault} />
                                    <Image src="/groupe.svg" alt="groups-hover" width={25} height={25} className={Styles.iconHover} />
                                </span>
                            </Link>
                            <Link className={`${Styles.linkWithIcon} ${pathname === "/chat" ? Styles.active : ""}`} href="/chat">
                                <span className={Styles.iconWrapper}>
                                    <Image src="/messages2.svg" alt="messages" width={25} height={25} className={Styles.iconDefault} />
                                    <Image src="/messages.svg" alt="messages-hover" width={25} height={25} className={Styles.iconHover} />
                                </span>
                            </Link>
                        </>
                    )}
                </div>

                {/* Right - Auth/Profile */}
                <div className={`${Styles.rightSection}`}>
                    {isLoggedIn ? (
                        <div className={Styles.dropdownWrapper}>
                            <div
                                className={Styles.profile}
                                onMouseEnter={() => setIsOpen(true)}
                                onMouseLeave={() => setIsOpen(false)}
                            >
                                <span className={Styles.iconUser}>
                                    <Image src="/iconMale.png" alt="profile" width={40} height={40} className={Styles.iconDefault} />
                                </span>
                                <span>{nickname}</span>

                                {isOpen && (
                                    <div className={Styles.dropdownMenu}>
                                        <Link href={`/profile/${nickname}`} onClick={() => setIsOpen(false)}>Profile</Link>
                                        <Link href="/logout" onClick={() => setIsOpen(false)}>Logout</Link>
                                    </div>
                                )}
                            </div>
                        </div>
                    ) : (
                        <>
                            <Link className={`${Styles.linkWithIcon} ${pathname === "/auth/login" ? Styles.active : ""}`} href="/auth/login" onClick={() => setIsOpen(false)}>
                                <span>Login</span>
                            </Link>
                            <Link className={`${Styles.linkWithIcon} ${pathname === "/auth/register" ? Styles.active : ""}`} href="/auth/register" onClick={() => setIsOpen(false)}>
                                <span>Register</span>
                            </Link>
                        </>
                    )}
                </div>
            </div>
            
            {/* navigation on small screens */}
            <div className={`${isLoggedIn ? Styles.bottomNav : Styles.logged}`}>
                {isLoggedIn && (
                    <>
                        <Link className={`${Styles.linkWithIcon} ${pathname === "/" ? Styles.active : ""}`} href="/">
                            <span className={Styles.iconWrapper}>
                                <Image src="/home2.svg" alt="home" width={25} height={25} className={Styles.iconDefault} />
                                <Image src="/home.svg" alt="home-hover" width={25} height={25} className={Styles.iconHover} />
                            </span>
                        </Link>
                        <Link className={`${Styles.linkWithIcon} ${pathname === "/posts" ? Styles.active : ""}`} href="/posts">
                            <span className={Styles.iconWrapper}>
                                <Image src="/posts2.svg" alt="posts" width={25} height={25} className={Styles.iconDefault} />
                                <Image src="/posts.svg" alt="posts-hover" width={25} height={25} className={Styles.iconHover} />
                            </span>
                        </Link>
                        <Link className={`${Styles.linkWithIcon} ${pathname === "/groups" ? Styles.active : ""}`} href="/groups">
                            <span className={Styles.iconWrapper}>
                                <Image src="/groupe2.svg" alt="groups" width={25} height={25} className={Styles.iconDefault} />
                                <Image src="/groupe.svg" alt="groups-hover" width={25} height={25} className={Styles.iconHover} />
                            </span>
                        </Link>
                        <Link className={`${Styles.linkWithIcon} ${pathname === "/chat" ? Styles.active : ""}`} href="/chat">
                            <span className={Styles.iconWrapper}>
                                <Image src="/messages2.svg" alt="messages" width={25} height={25} className={Styles.iconDefault} />
                                <Image src="/messages.svg" alt="messages-hover" width={25} height={25} className={Styles.iconHover} />
                            </span>
                        </Link>
                    </>
                )
                }
            </div >
        </div >



    );
}
