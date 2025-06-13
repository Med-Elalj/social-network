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
        const response = await SendData('/api/v1/auth', null);
        if (response.status === 200) {
          setIsLoggedIn(true);
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
        <div className={Styles.leftSection}>
          <Link className={Styles.loginTitle} href="/">Social Network</Link>
        </div>

        {isLoggedIn && (
          <div className={Styles.centerSection}>
            <NavLink href="/" icon="home" pathname={pathname} />
            <NavLink href="/newPost" icon="posts" pathname={pathname} />
            <NavLink href="/groupes" icon="groupe" pathname={pathname} />
            <NavLink href="/chat" icon="messages" pathname={pathname} />
          </div>
        )}

        <div className={Styles.rightSection}>
          {isLoggedIn ? (
            <>
              <div className={Styles.dropdownWrapper}>
                <div
                  className={Styles.notif}
                  onClick={() => setIsOpen(true)}
                  onMouseLeave={() => setIsOpen(false)}
                >
                  <span>
                    <Image src="/notification.svg" alt="notification" width={25} height={25} />
                  </span>
                  {isOpen && (
                    <div className={`${Styles.dropdownMenu} ${Styles.notification}`}>
                      <Link href={`/`} onClick={() => setIsOpen(false)}>Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum.</Link>
                      <Link href={`/`} onClick={() => setIsOpen(false)}>test2</Link>
                    </div>
                  )}
                </div>
              </div>
              <div className={Styles.dropdownWrapper}>
                <div
                  className={Styles.profile}
                  onClick={() => setIsOpen(true)}
                  onMouseLeave={() => setIsOpen(false)}
                >
                  <span className={Styles.iconUser}>
                    <Image src="/iconMale.png" alt="profile" width={40} height={40} />
                  </span>
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
            </>
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
            <NavLink href="/newPost" icon="posts" pathname={pathname} />
            <NavLink href="/groups" icon="groupe" pathname={pathname} />
            <NavLink href="/chat" icon="messages" pathname={pathname} />
          </>
        )}
      </div>
    </div >
  );
}

function NavLink({ href, icon, pathname }) {
  return (
    <Link className={`${Styles.linkWithIcon} ${pathname === href ? Styles.active : ""}`} href={href}>
      <span className={Styles.iconWrapper}>
        <Image src={`/${icon}2.svg`} alt={`${icon}-hover`} width={25} height={25} className={Styles.iconHover} />
        <Image src={`/${icon}.svg`} alt={icon} width={25} height={25} className={Styles.iconDefault} />
      </span>
    </Link>
  );
}
