"use client";

import Style from "./profile.module.css";
import Image from "next/image";
import Link from "next/link";
import { useEffect, useState } from "react";
import Posts from "./[tab]/Posts";
import Following from "./[tab]/Following";
import Followers from "./[tab]/Followers";
import Settings from "./[tab]/Settings";
import { GetData, SendData } from "../../../utils/sendData.js";
import { CapitalizeFirstLetter, showNotification } from "../utils.jsx";

// ✅ Named sub-component
function PrivacyToggle({ isPublic, setIsPublic }) {
  const [loading, setLoading] = useState(false);

  const handlePrivacyChange = async () => {
    setLoading(true);

    console.log("Sending:", { privacy: !isPublic });
    const res = await SendData("/api/v1/settings/changePrivacy", {
      privacy: !isPublic,
    });

    const result = await res.json();
    setLoading(false);

    if (res.ok) {
      setIsPublic((prev) => !prev);
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

// ✅ Main default component
export default function Profile() {
  const [activeTab, setActiveTab] = useState("info");
  const [activeSection, setActiveSection] = useState("posts");
  const [profileData, setProfileData] = useState(null);
  const [isPublic, setIsPublic] = useState(false);

  useEffect(() => {
    const fetchProfile = async () => {
      try {
        const res = await GetData("/api/v1/profile");
        if (res.status === 200) {
          const data = await res.json();
          setProfileData(data);
          if (typeof data.is_public === "boolean") {
            setIsPublic(data.is_public);
          }
        } else {
          console.error("Failed to fetch profile data");
        }
      } catch (error) {
        console.error("Error fetching profile data:", error);
      }
    };
    fetchProfile();
  }, []);

  useEffect(() => {
    if (activeTab === "settings") {
      setActiveSection("Settings");
    }
  }, [activeTab]);

  return (
    <div className={Style.container}>
      <div className={Style.header}>
        <Image src="/db.png" fill alt="cover" />
      </div>
      <div className={Style.body}>
        <div className={Style.first}>
          <div className={Style.ProfileInfo}>
            <div className={Style.top}>
              <div
                style={{
                  position: "relative",
                  width: "200px",
                  height: "200px",
                }}
              >
                <Image src="/db.png" alt="cover" fill />
              </div>
              <h4>@{CapitalizeFirstLetter(profileData?.display_name)}</h4>
            </div>
            <PrivacyToggle isPublic={isPublic} setIsPublic={setIsPublic} />
            <div className={Style.tabs}>
              <button onClick={() => setActiveTab("info")}>Info</button>
              <button onClick={() => setActiveTab("connections")}>
                Connections
              </button>
              <button onClick={() => setActiveTab("settings")}>Settings</button>
            </div>
            {activeTab === "info" && (
              <>
                <div className={Style.center}>
                  <span>
                    <h5>About me:</h5>&nbsp;&nbsp;
                    <h5>
                      {profileData?.description
                        ? profileData?.description
                        : "No description provided."}
                    </h5>
                  </span>
                </div>
                <div className={Style.center}>
                  <span>
                    <h5>Full Name :</h5>&nbsp;&nbsp;
                    <h5>
                      {profileData?.first_name} {profileData?.last_name}
                    </h5>
                  </span>
                  <span>
                    <h5>Email :</h5>&nbsp;&nbsp;
                    <h5>{profileData?.email}</h5>
                  </span>
                  <span>
                    <h5>Age :</h5>&nbsp;&nbsp;
                    <h5>
                      {(
                        new Date().getFullYear() -
                        new Date(profileData?.date_of_birth).getFullYear()
                      ).toString()}
                    </h5>
                  </span>
                  <span>
                    <h5>Birthdate :</h5>&nbsp;&nbsp;
                    <h5>{profileData?.date_of_birth}</h5>
                  </span>
                </div>
              </>
            )}

            {activeTab === "settings" && (
              <>
                <span onClick={() => setActiveSection("Settings")}></span>
              </>
            )}

            {activeTab === "connections" && (
              <>
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
              </>
            )}
          </div>
        </div>

        <div className={Style.second}>
          {activeSection === "Settings" && <Settings />}

          {activeSection === "posts" && <Posts />}

          {activeSection === "followers" && <Followers />}

          {activeSection === "following" && <Following />}
        </div>

        <div className={Style.end}>
          <div className={Style.requists}>
            <h3>Suggestion</h3>
            <div>
              <div>
                <Image
                  src="/db.png"
                  alt="profile"
                  width={40}
                  height={40}
                  style={{ borderRadius: "50%" }}
                />
                <h5>username</h5>
              </div>
              <Link href="/addUser">
                <Image
                  src="/addUser.svg"
                  alt="profile"
                  width={25}
                  height={25}
                />
              </Link>
            </div>
            <div>
              <div>
                <Image
                  src="/db.png"
                  alt="profile"
                  width={40}
                  height={40}
                  style={{ borderRadius: "50%" }}
                />
                <h5>username</h5>
              </div>
              <Link href="/addUser">
                <Image
                  src="/addUser.svg"
                  alt="profile"
                  width={25}
                  height={25}
                />
              </Link>
            </div>
            <div>
              <div>
                <Image
                  src="/db.png"
                  alt="profile"
                  width={40}
                  height={40}
                  style={{ borderRadius: "50%" }}
                />
                <h5>username</h5>
              </div>
              <Link href="/addUser">
                <Image
                  src="/addUser.svg"
                  alt="profile"
                  width={25}
                  height={25}
                />
              </Link>
            </div>
            <div>
              <div>
                <Image
                  src="/db.png"
                  alt="profile"
                  width={40}
                  height={40}
                  style={{ borderRadius: "50%" }}
                />
                <h5>username</h5>
              </div>
              <Link href="/addUser">
                <Image
                  src="/addUser.svg"
                  alt="profile"
                  width={25}
                  height={25}
                />
              </Link>
            </div>
            <div>
              <div>
                <Image
                  src="/db.png"
                  alt="profile"
                  width={40}
                  height={40}
                  style={{ borderRadius: "50%" }}
                />
                <h5>username</h5>
              </div>
              <Link href="/addUser">
                <Image
                  src="/addUser.svg"
                  alt="profile"
                  width={25}
                  height={25}
                />
              </Link>
            </div>
          </div>

          <div className={Style.requists}>
            <div>
              <div>
                <Image
                  src="/iconMale.png"
                  alt="profile"
                  width={40}
                  height={40}
                />
                <h5>Username</h5>
              </div>
              <div className={Style.Buttons}>
                <Link href="/accept">
                  <Image
                    src="/accept.svg"
                    alt="profile"
                    width={25}
                    height={25}
                  />
                </Link>
                <Link href="/reject">
                  <Image
                    src="/reject.svg"
                    alt="profile"
                    width={25}
                    height={25}
                  />
                </Link>
              </div>
            </div>
            <div>
              <div>
                <Image
                  src="/iconMale.png"
                  alt="profile"
                  width={40}
                  height={40}
                />
                <h5>Username</h5>
              </div>
              <div className={Style.Buttons}>
                <Link href="/accept">
                  <Image
                    src="/accept.svg"
                    alt="profile"
                    width={25}
                    height={25}
                  />
                </Link>
                <Link href="/reject">
                  <Image
                    src="/reject.svg"
                    alt="profile"
                    width={25}
                    height={25}
                  />
                </Link>
              </div>
            </div>
            <div>
              <div>
                <Image
                  src="/iconMale.png"
                  alt="profile"
                  width={40}
                  height={40}
                />
                <h5>Username</h5>
              </div>
              <div className={Style.Buttons}>
                <Link href="/accept">
                  <Image
                    src="/accept.svg"
                    alt="profile"
                    width={25}
                    height={25}
                  />
                </Link>
                <Link href="/reject">
                  <Image
                    src="/reject.svg"
                    alt="profile"
                    width={25}
                    height={25}
                  />
                </Link>
              </div>
            </div>
            <div>
              <div>
                <Image
                  src="/iconMale.png"
                  alt="profile"
                  width={40}
                  height={40}
                />
                <h5>Username</h5>
              </div>
              <div className={Style.Buttons}>
                <Link href="/accept">
                  <Image
                    src="/accept.svg"
                    alt="profile"
                    width={25}
                    height={25}
                  />
                </Link>
                <Link href="/reject">
                  <Image
                    src="/reject.svg"
                    alt="profile"
                    width={25}
                    height={25}
                  />
                </Link>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
