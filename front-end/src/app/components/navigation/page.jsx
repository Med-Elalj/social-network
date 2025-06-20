"use client";

import { usePathname } from 'next/navigation';
import { SendData } from '../../../../utils/sendData.js';
import { useState, useEffect } from 'react';
import { LogoutAndRedirect } from "../Logout.jsx";
import { useRouter } from 'next/navigation';
import Link from "next/link";
import Image from 'next/image';
import Styles from "./nav.module.css";

export default function Routing() {
  const [isOpen, setIsOpen] = useState(false);
  const [isLoggedIn, setIsLoggedIn] = useState(null);
  const pathname = usePathname();
  const router = useRouter();

  useEffect(() => {
    const fetchAuthStatus = async () => {
      try {
        const response = await SendData('/api/v1/auth', null);
        setIsLoggedIn(response.status === 200);
      } catch (err) {
        setIsLoggedIn(false);
      }
    };
    fetchAuthStatus();
  });

  const validPaths = ["/", "/login", "/register", "/newPost", "/groupes", "/chat", "/profile/nickname"];

  useEffect(() => {
    if (isLoggedIn === null) return;

    if (!validPaths.includes(pathname)) return;

    if (isLoggedIn && (pathname === "/login" || pathname === "/register")) {
      router.push('/' || pathname);
    } else if (!isLoggedIn && pathname !== "/login" && pathname !== "/register") {
      router.push("/login");
    }
  }, [isLoggedIn, pathname]);

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
                          await LogoutAndRedirect(router);
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
              <Link className={`${Styles.linkWithIcon} ${pathname === "/login" ? Styles.active : ""}`} href="/login" onClick={() => setIsOpen(false)}>Login</Link>
              <Link className={`${Styles.linkWithIcon} ${pathname === "/register" ? Styles.active : ""}`} href="/register" onClick={() => setIsOpen(false)}>Register</Link>
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
            <NavLink href="/groupes" icon="groupe" pathname={pathname} />
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
