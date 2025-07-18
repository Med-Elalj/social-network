"use client";
import { usePathname } from "next/navigation";
import { useState, useEffect, use } from "react";
import { SendData, refreshAccessToken } from "@/app/sendData.js";
import { LogoutAndRedirect } from "../Logout.jsx";
import { useRouter } from "next/navigation";
import Link from "next/link";
import Image from "next/image";
import Styles from "./nav.module.css";
import NotificationList from "./notificationList.jsx";
import { useWebSocket } from "@/app/context/WebSocketContext.jsx";
import { SearchIcon, SearchInput } from "./search.jsx";
import {  useAuth } from "@/app/context/AuthContext.jsx";
import { externalNotification } from "@/app/context/NotificationContext.jsx";
import { UserAvatar } from "../upload.jsx";


const RefreshFrequency = 10 * (60 * 1000); // 14 mins since JWT expiry is 15mins

function ProtectedLink({ isLoggedIn, href, className, children }) {
  const handleClick = (e) => {
    if (!isLoggedIn) {
      e.preventDefault();
    }
  };

  return (
    <Link
      href={href}
      className={className}
      onClick={handleClick}
      style={{
        opacity: isLoggedIn ? 1 : 0.5,
        cursor: isLoggedIn ? "pointer" : "not-allowed",
      }}
    >
      {children}
    </Link>
  );
}


const Routing = () => {
  const [isOpen, setIsOpen] = useState(false);
  const { isLoggedIn, loading, setIsLoggedIn } = useAuth();
  const pathname = usePathname();
  const router = useRouter();
  const { closeWebSocket, isConnected } = useWebSocket();
  const [notifications, setNotifications] = useState([]);
  const [showSearch, setShowSearch] = useState(false);

  const publicRoutes = ["/login", "/register"];
  const isPublic = publicRoutes.some((route) => pathname.startsWith(route));
  // 🧑‍💻 Redirect Logic and Route Protection
  useEffect(() => {
    if (loading) return;

    if (!isLoggedIn && !isPublic) {
      router.replace("/login");
    } else if (isLoggedIn && isPublic && pathname !== "/") {
      router.replace("/");
    }
  }, [isLoggedIn, pathname, loading, router]);

  useEffect(() => {
    if (!isLoggedIn) return;

    const interval = setInterval(() => {
      refreshAccessToken();
    }, RefreshFrequency);

    return () => clearInterval(interval);
  }, [isLoggedIn]);
  console.log(isLoggedIn, "isLoggedIn",loading,"loading");
  
  // Fetch notifications when dropdown opens
  useEffect(() => {
    const fetchNotifications = async () => {
      try {
        const response = await SendData("/api/v1/get/requests", { type: 3 });
        if (!response.ok) throw new Error("Failed to fetch notifications");

        const data = await response.json();
        setNotifications(data);
      } catch (err) {
        console.error("Error fetching notifications:", err);
      }
    };

    if (isOpen) fetchNotifications();
  }, [isOpen]);

  // Close search modal
  const handleSearchClose = () => setShowSearch(false);

  return (
    <div>
      <div className={Styles.nav}>
        <div className={Styles.leftSection}>
          <ProtectedLink
            className={Styles.loginTitle}
            href={"/"}
            isLoggedIn={isLoggedIn}
          >
            Social Network
          </ProtectedLink>
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
                  <Image
                    src="/notification.svg"
                    alt="notification"
                    width={25}
                    height={25}
                  />
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
                  <UserAvatar className={Styles.profileIcon} />
                  {/* <span className={Styles.iconUser}>
                    <Image
                      src="/iconMale.png"
                      alt="profile"
                      width={40}
                      height={40}
                    />
                  </span> */}
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
                          await LogoutAndRedirect({
                            router,
                            isLoggedIn,
                            setIsLoggedIn,
                          });
                          if (isConnected) closeWebSocket();
                          setIsOpen(false);
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
              >
                Login
              </Link>
              <Link
                className={`${Styles.linkWithIcon} ${
                  pathname === "/register" ? Styles.active : ""
                }`}
                href="/register"
              >
                Register
              </Link>
            </>
          )}
        </div>
      </div>

      {/* Bottom nav */}
      <div className={`${isLoggedIn ? Styles.bottomNav : Styles.logged}`}>
          <>
            <NavLink
              href="/"
              icon="home"
              pathname={pathname}
            />
            <NavLink
              href="/newPost"
              icon="posts"
              pathname={pathname}
            />
            <NavLink
              href="/groupes/feed"
              icon="groupe"
              pathname={pathname}
            />
            <NavLink
              href="/chat"
              icon="messages"
              pathname={pathname}
            />
            <SearchIcon
              onClick={() => setShowSearch(true)}
              showSearch={showSearch}
            />
          </>
      </div>

      {showSearch && <SearchInput onClose={handleSearchClose} />}
    </div>
  );
};

// Helper component for NavLink
function NavLink({ href, icon, pathname }) {
  const { isLoggedIn, loading } = useAuth();

  if (loading || isLoggedIn === null) return null;

  return (
    <ProtectedLink
      className={`${Styles.linkWithIcon} ${
        pathname === href ? Styles.active : ""
      }`}
      href={href}
      isLoggedIn={isLoggedIn}
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
    </ProtectedLink>
  );
}

export default Routing;