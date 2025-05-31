"use client";

import { useState } from 'react';
import Link from "next/link";
import Image from 'next/image';
import Styles from "./nav.module.css";

export default function Routing() {
    const [isOpen, setIsOpen] = useState(false);
    const isLoggedIn = false;
    const nickname = "nickname";

    return (
        <div className={Styles.nav}>
            {/* Left - Logo */}
            <div className={Styles.leftSection}>
                <Link className={Styles.loginTitle} href="/">Social Network</Link>
            </div>

            {/* Center - Navigation Links */}
            <div className={`${Styles.centerSection} ${isOpen ? Styles.open : ""}`}>
                {isLoggedIn && (
                    <>
                        <Link className={Styles.linkWithIcon} href="/">
                            <span className={Styles.iconWrapper}>
                                <Image src="/home.svg" alt="home" width={20} height={20} className={Styles.iconDefault} />
                                <Image src="/home2.svg" alt="home-hover" width={20} height={20} className={Styles.iconHover} />
                            </span>
                            <span>Dashboard</span>
                        </Link>
                        <Link className={Styles.linkWithIcon} href="/posts">
                            <span className={Styles.iconWrapper}>
                                <Image src="/posts.svg" alt="posts" width={20} height={20} className={Styles.iconDefault} />
                                <Image src="/posts2.svg" alt="posts-hover" width={20} height={20} className={Styles.iconHover} />
                            </span>
                            <span>Posts</span>
                        </Link>
                        <Link className={Styles.linkWithIcon} href="/groups">
                            <span className={Styles.iconWrapper}>
                                <Image src="/groupe.svg" alt="groups" width={20} height={20} className={Styles.iconDefault} />
                                <Image src="/groupe2.svg" alt="groups-hover" width={20} height={20} className={Styles.iconHover} />
                            </span>
                            <span>Groups</span>
                        </Link>
                        <Link className={Styles.linkWithIcon} href="/chat">
                            <span className={Styles.iconWrapper}>
                                <Image src="/messages.svg" alt="messages" width={20} height={20} className={Styles.iconDefault} />
                                <Image src="/messages2.svg" alt="messages-hover" width={20} height={20} className={Styles.iconHover} />
                            </span>
                            <span>Messages</span>
                        </Link>
                    </>
                )}
            </div>

            {/* Right - Auth/Profile */}
            <div className={`${Styles.rightSection} ${isOpen ? Styles.open : ""}`}>
                {isLoggedIn ? (
                    <div className={Styles.dropdownWrapper}>
                        <div
                            className={Styles.profile}
                            onMouseEnter={() => setIsOpen(true)}
                            onMouseLeave={() => setIsOpen(false)}
                        >
                            <span className={Styles.iconWrapper}>
                                <Image src="/login.svg" alt="profile" width={20} height={20} className={Styles.iconDefault} />
                                <Image src="/login2.svg" alt="profile-hover" width={20} height={20} className={Styles.iconHover} />
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
                        <Link className={Styles.linkWithIcon} href="/auth/login" onClick={() => setIsOpen(false)}>
                            <span className={Styles.iconWrapper}>
                                <Image src="/login.svg" alt="login" width={20} height={20} className={Styles.iconDefault} />
                                <Image src="/login2.svg" alt="login-hover" width={20} height={20} className={Styles.iconHover} />
                            </span>
                            <span>Login</span>
                        </Link>
                        <Link className={Styles.linkWithIcon} href="/auth/register" onClick={() => setIsOpen(false)}>
                            <span className={Styles.iconWrapper}>
                                <Image src="/register.svg" alt="register" width={20} height={20} className={Styles.iconDefault} />
                                <Image src="/register2.svg" alt="register-hover" width={20} height={20} className={Styles.iconHover} />
                            </span>
                            <span>Register</span>
                        </Link>
                    </>
                )}
            </div>

            {/* Hamburger Button */}
            <button
                className={Styles.menuToggle}
                onClick={() => setIsOpen(!isOpen)}
                aria-label="Toggle menu"
            >
                {isOpen ? (
                    <Image src="/close.svg" alt="Close menu" width={24} height={24} />
                ) : (
                    <Image src="/menu.svg" alt="Open menu" width={24} height={24} />
                )}
            </button>
        </div>
    );
}
