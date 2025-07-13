"use client";
import { usePathname } from "next/navigation";
import { useState, useEffect } from "react";
import { GetData, SendData } from "@/app/sendData.js";
import { LogoutAndRedirect } from "../Logout.jsx";
import { useRouter } from "next/navigation";
import Link from "next/link";
import Image from "next/image";
import Styles from "./nav.module.css";
import NotificationList from "./notificationList.jsx";
import { refreshAccessToken } from "@/app/sendData.js";
import { useWebSocket } from "@/app/context/WebSocketContext.jsx";
import { SearchIcon, SearchInput } from "./search.jsx"; // Import SearchInput too
import { useAuth } from "@/app/context/AuthContext.jsx";
const RefreshFrequency = 14 * (60 * 1000); // 14 mins since jwt expiry is 15mins

export default function Routing() {
  const [isOpen, setIsOpen] = useState(false);
  const { isLoggedIn, setIsLoggedIn } = useAuth();
  const pathname = usePathname();
  const router = useRouter();
  const { closeWebSocket, isConnected } = useWebSocket();
  const [notifications, setNotifications] = useState([]);
  const [showSearch, setShowSearch] = useState(false);

  useEffect(() => {
    const fetchNotifications = async () => {
      try {
        const response = await SendData("/api/v1/get/requests", 3);

        if (!response.ok) {
          console.error("Failed to fetch notifications");
          return;
        }

        const data = await response.json();
        setNotifications(data);
        console.log("Notifications:", data);
      } catch (err) {
        console.error("Error fetching notifications:", err);
      }
    };

    if (isOpen) {
      fetchNotifications();
      fetchNotifications();
    }
  }, [isOpen]);

  // // 🛰️ Route Protection
  const publicRoutes = [ "/login", "/register"];
  const isPublic = publicRoutes.some((route) => pathname.startsWith(route));

  // ⛔ Route Redirect
  useEffect(() => {
    if (!isLoggedIn && !isPublic) {
      console.log("Redirecting to /login", isLoggedIn);
      router.push("/login");
    }else if (isLoggedIn && isPublic) {
      router.push("/");
    }
  }, [isLoggedIn, pathname]);

  // Function to handle search close
  const handleSearchClose = () => {
    setShowSearch(false);
  };

  return (
    <div>
      <div className={Styles.nav}>
        <div className={Styles.leftSection}>
          
          <Link className={Styles.loginTitle} href={isLoggedIn ? "/" : "/login"}>
            Social Network
          </Link>
        </div>

        {isLoggedIn && (
          <div className={Styles.centerSection}>
            <NavLink href="/" icon="home" pathname={pathname} />
            <NavLink href="/newPost" icon="posts" pathname={pathname} />
            <NavLink href="/groupes/feed" icon="groupe" pathname={pathname} />
            <NavLink href="/chat" icon="messages" pathname={pathname} />
            <SearchIcon
              onClick={() => setShowSearch(true)}
              showSearch={showSearch}
            />
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
                    <Image
                      src="/notification.svg"
                      alt="notification"
                      width={25}
                      height={25}
                    />
                  </span>
                  {isOpen && (
                    <NotificationList
                      notifications={notifications}
                      setIsOpen={setIsOpen}
                    />
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
                    <Image
                      src="/iconMale.png"
                      alt="profile"
                      width={40}
                      height={40}
                    />
                  </span>
                  {isOpen && (
                    <div className={Styles.dropdownMenu}>
                      <Link
                        href={`/profile/me`}
                        onClick={() => setIsOpen(false)}
                      >
                        Profile
                      </Link>
                      <button
                        onClick={async () => {
                          await LogoutAndRedirect(router);
                          if (isConnected) closeWebSocket();
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
              <Link
                className={`${Styles.linkWithIcon} ${
                  pathname === "/login" ? Styles.active : ""
                }`}
                href="/login"
                onClick={() => setIsOpen(false)}
              >
                Login
              </Link>
              <Link
                className={`${Styles.linkWithIcon} ${
                  pathname === "/register" ? Styles.active : ""
                }`}
                href="/register"
                onClick={() => setIsOpen(false)}
              >
                Register
              </Link>
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
            <NavLink href="/groupes/feed" icon="groupe" pathname={pathname} />
            <NavLink href="/chat" icon="messages" pathname={pathname} />
            <SearchIcon
              onClick={() => setShowSearch(true)}
              showSearch={showSearch}
            />
          </>
        )}
      </div>

      {/* Search popup - This was missing! */}
      {showSearch && <SearchInput onClose={handleSearchClose} />}
    </div>
  );
}

function NavLink({ href, icon, pathname }) {
  return (
    <Link
      className={`${Styles.linkWithIcon} ${
        pathname === href ? Styles.active : ""
      }`}
      href={href}
    >
      <span className={Styles.iconWrapper}>
        <Image
          src={`/${icon}2.svg`}
          alt={`${icon}-hover`}
          width={25}
          height={25}
          className={Styles.iconHover}
        />
        <Image
          src={`/${icon}.svg`}
          alt={icon}
          width={25}
          height={25}
          className={Styles.iconDefault}
        />
      </span>
    </Link>
  );
}