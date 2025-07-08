"use client";

import Style from "../profile.module.css";
import Image from "next/image";
import Link from "next/link";
import { useEffect, useState } from "react";
import Posts from "../[tab]/Posts.jsx";
import Following from "../[tab]/Following.jsx";
import Followers from "../[tab]/Followers.jsx";
import { GetData } from "../../../../utils/sendData.js";
import { CapitalizeFirstLetter } from "../../utils.jsx";
import { useParams } from "next/navigation";

export default function PubProfile() {
  const [activeTab, setActiveTab] = useState("info");
  const [activeSection, setActiveSection] = useState("posts");
  const [profileData, setProfileData] = useState(null);
  const params = useParams();
  const [notFound, setNotFound] = useState(false);
  const nickname = params.nickname;

  useEffect(() => {
    const fetchProfile = async () => {
      try {
        const res = await GetData(`/api/v1/profile/${nickname}`);
        if (res.ok) {
          const data = await res.json();
          setProfileData(data);
        } else {
          setNotFound(true);
          console.log("Profile not found", nickname);
          
        }
      } catch (err) {
        console.error("Error fetching profile:", err);
        setNotFound(true);
      }
    };

    if (nickname) {
      fetchProfile();
    }
  }, [nickname]);
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

            <div className={Style.tabs}>
              <button onClick={() => setActiveTab("info")}>Info</button>
              <button onClick={() => setActiveTab("connections")}>
                Connections
              </button>
            </div>

            {activeTab === "info" && (
              <>
                <div className={Style.center}>
                  <span>
                    <h5>About me:</h5>&nbsp;&nbsp;
                    <h5>
                      {profileData?.description || "No description provided."}
                    </h5>
                  </span>
                </div>

                {profileData?.isPublic === true && (
                  <div className={Style.center}>
                    <span>
                      <h5>Full Name:</h5>&nbsp;&nbsp;
                      <h5>
                        {profileData?.first_name} {profileData?.last_name}
                      </h5>
                    </span>
                    <span>
                      <h5>Email:</h5>&nbsp;&nbsp;
                      <h5>{profileData?.email}</h5>
                    </span>
                    <span>
                      <h5>Age:</h5>&nbsp;&nbsp;
                      <h5>
                        {profileData?.date_of_birth
                          ? (
                              new Date().getFullYear() -
                              new Date(profileData.date_of_birth).getFullYear()
                            ).toString()
                          : "N/A"}
                      </h5>
                    </span>
                    <span>
                      <h5>Birthdate:</h5>&nbsp;&nbsp;
                      <h5>{profileData?.date_of_birth || "N/A"}</h5>
                    </span>
                  </div>
                )}
              </>
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
          <div className={Style.requists}>
            <h3>Suggestions</h3>
            {[...Array(1)].map((_, i) => (
              <div key={i}>
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
            ))}
          </div>

          <div className={Style.requists}>
            {[...Array(1)].map((_, i) => (
              <div key={i}>
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
            ))}
          </div>
        </div>
      </div>
    </div>
  );
}
