"use client";

import { usePathname } from 'next/navigation';
import { useState, useEffect } from 'react';
import Link from "next/link";
import Image from 'next/image';
import Styles from "./nav.module.css";
import { SendData } from '../../../utils/sendData.js';
import { Logout } from "../../../EndPoints/Logout.js"

export default function Routing() {
  const [isOpen, setIsOpen] = useState(false);
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const pathname = usePathname();

  //TODO: if logged forbedin to use login/logout and the oposet true
  useEffect(() => {
    const fetchAuthStatus = async () => {
      try {
        const { status, body } = await SendData('/api/v1/auth', null);
        if (status === 200) {
          setIsLoggedIn(true);
          const { status: st, body: body2 } = await SendData('/api/v1/profile', { id: body.id });

          if (st === 200) {
            localStorage.setItem("UserInfo", JSON.stringify(body2.Userinfo));
          } else {
            console.log("err", st, body2);
          }
        }
      } catch (err) {
        console.error("Auth check failed", err);
      }
    };

    fetchAuthStatus();
  }, []);

  //TODO:fucntion to get data of user from local storage

  return (
    <div>
      <div className={Styles.nav}>
        {/* Left - Logo */}
        <div className={Styles.leftSection}>
          <Link className={Styles.loginTitle} href="/">Social Network</Link>
        </div>

        {/* Center - Navigation Links */}
        <div className={Styles.centerSection}>
          {isLoggedIn && (
            <>
              <NavLink href="/" icon="home" pathname={pathname} />
              <NavLink href="/posts" icon="posts" pathname={pathname} />
              <NavLink href="/groups" icon="groupe" pathname={pathname} />
              <NavLink href="/chat" icon="messages" pathname={pathname} />
            </>
          )}
        </div>

        {/* Right - Auth/Profile */}
        <div className={Styles.rightSection}>
          {isLoggedIn ? (
            <div className={Styles.dropdownWrapper}>
              <div
                className={Styles.profile}
                onMouseEnter={() => setIsOpen(true)}
                onMouseLeave={() => setIsOpen(false)}
              >
                <span className={Styles.iconUser}>
                  <Image src="/iconMale.png" alt="profile" width={40} height={40} />
                </span>
                <span>nickname</span>
                {isOpen && (
                  <div className={Styles.dropdownMenu}>
                    <Link href={`/profile/nickname`} onClick={() => setIsOpen(false)}>Profile</Link>
                    <button
                      onClick={async () => {
                        await Logout();
                        setIsOpen(false);
                        setIsLoggedIn(false);
                      }}
                      className={Styles.dropdownItem}
                    >
                      Logout
                    </button>
                  </div>
                )}
              </div>
            </div>
          ) : (
            <>
              <Link className={`${Styles.linkWithIcon} ${pathname === "/auth/login" ? Styles.active : ""}`} href="/auth/login" onClick={() => setIsOpen(false)}>Login</Link>
              <Link className={`${Styles.linkWithIcon} ${pathname === "/auth/register" ? Styles.active : ""}`} href="/auth/register" onClick={() => setIsOpen(false)}>Register</Link>
            </>
          )}
        </div>
      </div>

      {/* Bottom nav */}
      <div className={`${isLoggedIn ? Styles.bottomNav : Styles.logged}`}>
        {isLoggedIn && (
          <>
            <NavLink href="/" icon="home" pathname={pathname} />
            <NavLink href="/posts" icon="posts" pathname={pathname} />
            <NavLink href="/groups" icon="groupe" pathname={pathname} />
            <NavLink href="/chat" icon="messages" pathname={pathname} />
          </>
        )}
      </div>
    </div >
  );
}

// 🔧 Helper for nav links
function NavLink({ href, icon, pathname }) {
  return (
    <Link className={`${Styles.linkWithIcon} ${pathname === href ? Styles.active : ""}`} href={href}>
      <span className={Styles.iconWrapper}>
        <Image src={`/${icon}2.svg`} alt={icon} width={25} height={25} className={Styles.iconDefault} />
        <Image src={`/${icon}.svg`} alt={`${icon}-hover`} width={25} height={25} className={Styles.iconHover} />
      </span>
    </Link>
  );
}
