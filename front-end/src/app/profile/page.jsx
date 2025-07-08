"use client";

import Style from "./profile.module.css";
import Image from "next/image";
import Link from "next/link";
import { useEffect, useState } from "react";
import { useParams, usePathname } from "next/navigation";
import { GetData, SendData } from "../../../utils/sendData.js";
import { CapitalizeFirstLetter, showNotification } from "../utils.jsx";

function PrivacyToggle({ isPublic, setIsPublic }) {
  const [loading, setLoading] = useState(false);

  const handlePrivacyChange = async () => {
    setLoading(true);
    const res = await SendData("/api/v1/settings/changePrivacy", {
      privacy: !isPublic,
    });

    const result = await res.json();
    setLoading(false);

    if (res.ok) {
      setIsPublic(!isPublic);
      showNotification("Privacy setting updated successfully!", "success");
    } else {
      showNotification(result.error, "error");
    }
  };

  return (
    <div className={Style.privacyToggle}>
      <p>Your profile is {isPublic ? "Public" : "Private"}</p>
      <label className={Style.switch}>
        <input
          type="checkbox"
          checked={isPublic}
          onChange={handlePrivacyChange}
          disabled={loading}
        />
        <span className={Style.slider}></span>
      </label>
    </div>
  );
}

export default function UnifiedProfile() {
  const [profileData, setProfileData] = useState(null);
  const [isPublic, setIsPublic] = useState(false);
  const [activeTab, setActiveTab] = useState("info");
  const [activeSection, setActiveSection] = useState("posts");
  const [notFound, setNotFound] = useState(false);

  const { nickname } = useParams() || {};
  const pathname = usePathname();

  useEffect(() => {
    const fetchProfile = async () => {
      try {
        const res = await GetData(nickname ? `/api/v1/profile/${nickname}` : `/api/v1/profile`);
        if (res.ok) {
          const data = await res.json();
          setProfileData(data);
          if (typeof data.isPublic === "boolean") {
            setIsPublic(data.isPublic);
          }
        } else {
          setNotFound(true);
        }
      } catch (err) {
        console.error("Error fetching profile:", err);
        setNotFound(true);
      }
    };

    fetchProfile();
  }, [nickname]);

  useEffect(() => {
    if (activeTab === "settings") {
      setActiveSection("Settings");
    }
  }, [activeTab]);

  if (notFound) {
    return (
      <div style={{ padding: "4rem", textAlign: "center", marginTop: "60px" }}>
        <Image src="/404.svg" alt="404" width={500} height={500} />
        <h1 style={{ color: "#8D6B0D" }}>404 - Profile Not Found</h1>
        <p style={{ color: "#e0e0e0" }}>
          Sorry, this profile doesn't exist.
        </p>
        <Link href="/" style={{ color: "#8D6B0D" }}>
          <p>Go back home</p>
        </Link>
      </div>
    );
  }

  if (!profileData) {
    return <p style={{ textAlign: "center", marginTop: "80px" }}>Loading...</p>;
  }

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
              <h4>@{CapitalizeFirstLetter(profileData.display_name)}</h4>
            </div>

            {profileData.isSelf && (
              <PrivacyToggle isPublic={isPublic} setIsPublic={setIsPublic} />
            )}

            <div className={Style.tabs}>
              <button onClick={() => setActiveTab("info")}>Info</button>
              <button onClick={() => setActiveTab("connections")}>Connections</button>
              {profileData.isSelf && (
                <button onClick={() => setActiveTab("settings")}>Settings</button>
              )}
            </div>

            {activeTab === "info" && (
              <>
                <div className={Style.center}>
                  <span>
                    <h5>About me:</h5>&nbsp;&nbsp;
                    <h5>{profileData.description || "No description provided."}</h5>
                  </span>
                </div>

                {(profileData.isSelf || profileData.isPublic) && (
                  <div className={Style.center}>
                    <span>
                      <h5>Full Name:</h5>&nbsp;&nbsp;
                      <h5>{profileData.first_name} {profileData.last_name}</h5>
                    </span>
                    <span>
                      <h5>Email:</h5>&nbsp;&nbsp;
                      <h5>{profileData.email}</h5>
                    </span>
                    <span>
                      <h5>Age:</h5>&nbsp;&nbsp;
                      <h5>
                        {profileData.date_of_birth
                          ? (
                              new Date().getFullYear() -
                              new Date(profileData.date_of_birth).getFullYear()
                            ).toString()
                          : "N/A"}
                      </h5>
                    </span>
                    <span>
                      <h5>Birthdate:</h5>&nbsp;&nbsp;
                      <h5>{profileData.date_of_birth || "N/A"}</h5>
                    </span>
                  </div>
                )}
              </>
            )}

            {activeTab === "settings" && profileData.isSelf && (
              <div className={Style.center}>
                {activeSection === "Settings" && <Settings />}
              </div>
            )}

            {activeTab === "connections" && (
              <div className={Style.numbers}>
                <span onClick={() => setActiveSection("posts")}>
                  <h4>Posts</h4>
                  <h5>0</h5>
                </span>
                <span onClick={() => setActiveSection("followers")}>
                  <h4>Followers</h4>
                  <h5>0</h5>
                </span>
                <span onClick={() => setActiveSection("following")}>
                  <h4>Following</h4>
                  <h5>0</h5>
                </span>
              </div>
            )}
          </div>
        </div>

        <div className={Style.second}>
          {activeSection === "posts" && <Posts />}
          {activeSection === "followers" && <Followers />}
          {activeSection === "following" && <Following />}
        </div>

        <div className={Style.end}>
          {/* Suggestions/Requests UI can be left unchanged */}
        </div>
      </div>
    </div>
  );
}
