"use client";

import Style from "../profile.module.css";
import Image from "next/image";
import Link from "next/link";
import { useEffect, useState } from "react";
import { useParams } from "next/navigation";
import { GetData, SendData } from "../../../../utils/sendData.js";
import { CapitalizeFirstLetter, showNotification } from "../../utils.jsx";
import Posts from "@/app/profile/[nickname]/[tab]/Posts.jsx";
import Following from "@/app/profile/[nickname]/[tab]/Following";
import Followers from "@/app/profile/[nickname]/[tab]/Followers";
import Settings from "@/app/profile/[nickname]/[tab]/Settings";

function FollowButton({ targetUsername, isPublic, following }) {
  const [requested, setRequested] = useState(false);
  const [isFollowing, setIsFollowing] = useState(false);
  const [isFollowBack, setIsFollowBack] = useState(false);

  useEffect(() => {
    if (following === "following") {
      setIsFollowing(true);
      setRequested(true);
      setIsFollowBack(false);
    } else if (following === "requested") {
      setIsFollowing(false);
      setRequested(true);
      setIsFollowBack(false);
    } else if (following === "followback") {
      setIsFollowing(false);
      setRequested(false);
      setIsFollowBack(true);
    } else {
      setIsFollowing(false);
      setRequested(false);
      setIsFollowBack(false);
    }
  }, [following]);

  const handleFollow = async () => {
    const FollowData = {
      targetUsername,
      isFollowing,
    };

    const res = await SendData(`/api/v1/set/follows`, FollowData);

    if (res.status === 200) {
      const result = await res.json();

      if (result.status === "following") {
        setIsFollowing(true);
        setRequested(true);
        setIsFollowBack(false);
      } else if (result.status === "requested") {
        setIsFollowing(false);
        setRequested(true);
        setIsFollowBack(false);
      } else {
        setIsFollowing(false);
        setRequested(false);
        setIsFollowBack(false);
      }
    } else {
      showNotification("Failed to follow user.", "error");
    }
  };


  return (
    <div style={{ marginTop: "1rem", alignSelf: "center" }}>
      {!requested && !isFollowBack && (
        <button
          onClick={handleFollow}
          className={`${Style.followBtn} ${Style.follow}`}
        >
          {isPublic ? "Follow" : "Request to Follow"}
        </button>
      )}

      {isFollowBack && (
        <button
          onClick={handleFollow}
          className={`${Style.followBtn} ${Style.follow}`}
        >
          Follow Back
        </button>
      )}

      {requested && isFollowing && (
        <button
          onClick={handleFollow}
          className={`${Style.followBtn} ${Style.unfollow}`}
        >
          Unfollow
        </button>
      )}

      {requested && !isFollowing && (
        <button
          onClick={handleFollow}
          className={`${Style.followBtn} ${Style.unfollow}`}>
          Requested
        </button>
      )}
    </div>
  );
}


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

export default function Profile() {
  const [profileData, setProfileData] = useState(null);
  const [isPublic, setIsPublic] = useState(false);
  const [activeTab, setActiveTab] = useState("info");
  const [activeSection, setActiveSection] = useState("posts");
  const [notFound, setNotFound] = useState(false);

  const { nickname } = useParams() || {};

  // const avatarUrl = profileData.avatar?.String
  //   ? `/${profileData.avatar.String}`
  //   : "/default-avatar.png";

  useEffect(() => {
    const fetchProfile = async () => {
      try {
        const res = await GetData(`/api/v1/profile/${nickname}`);
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
        <p style={{ color: "#e0e0e0" }}>Sorry, this profile doesn't exist.</p>
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
                <Image
                  src={profileData.avatar.String || "/default-avatar.png"}
                  alt="user avatar"
                  fill
                  style={{ borderRadius: "50%" }}
                />
              </div>
              <h4>@{CapitalizeFirstLetter(profileData.display_name)}</h4>
            </div>
            {!profileData.isSelf && (
              <FollowButton
                targetUsername={profileData.display_name}
                isPublic={profileData.isPublic}
                following={profileData.isFollowed}
              />
            )}
            {profileData.isSelf && <PrivacyToggle isPublic={isPublic} setIsPublic={setIsPublic} />}

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
                      <h5>
                        {profileData.first_name} {profileData.last_name}
                      </h5>
                    </span>
                     <span>
                      <h5>Gender:</h5>&nbsp;&nbsp;
                      <h5>
                        {profileData.gender}
                      </h5>
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
              <div className={Style.center}>{activeSection === "Settings" && <Settings />}</div>
            )}

            {activeTab === "connections" && (
              <div className={Style.numbers}>
                <span onClick={() => setActiveSection("posts")}>
                  <h4>Posts</h4>
                  <h5>{profileData.post_count}</h5>
                </span>
                <span onClick={() => setActiveSection("followers")}>
                  <h4>Followers</h4>
                  <h5>{profileData.follower_count}</h5>
                </span>
                <span onClick={() => setActiveSection("following")}>
                  <h4>Following</h4>
                  <h5>{profileData.following_count}</h5>
                </span>
              </div>
            )}
          </div>
        </div>

        <div className={Style.second}>
          {activeSection === "posts" && <Posts userId={profileData.id} />}
          {activeSection === "followers" && <Followers userId={profileData.id}/>}
          {activeSection === "following" && <Following userId={profileData.id}/>}
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
                <Image src="/addUser.svg" alt="profile" width={25} height={25} />
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
                <Image src="/addUser.svg" alt="profile" width={25} height={25} />
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
                <Image src="/addUser.svg" alt="profile" width={25} height={25} />
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
                <Image src="/addUser.svg" alt="profile" width={25} height={25} />
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
                <Image src="/addUser.svg" alt="profile" width={25} height={25} />
              </Link>
            </div>
          </div>

          <div className={Style.requists}>
            <h3>Requests</h3>
            <div>
              <div>
                <Image src="/iconMale.png" alt="profile" width={40} height={40} />
                <h5>Username</h5>
              </div>
              <div className={Style.Buttons}>
                <Link href="/accept">
                  <Image src="/accept.svg" alt="profile" width={25} height={25} />
                </Link>
                <Link href="/reject">
                  <Image src="/reject.svg" alt="profile" width={25} height={25} />
                </Link>
              </div>
            </div>
            <div>
              <div>
                <Image src="/iconMale.png" alt="profile" width={40} height={40} />
                <h5>Username</h5>
              </div>
              <div className={Style.Buttons}>
                <Link href="/accept">
                  <Image src="/accept.svg" alt="profile" width={25} height={25} />
                </Link>
                <Link href="/reject">
                  <Image src="/reject.svg" alt="profile" width={25} height={25} />
                </Link>
              </div>
            </div>
            <div>
              <div>
                <Image src="/iconMale.png" alt="profile" width={40} height={40} />
                <h5>Username</h5>
              </div>
              <div className={Style.Buttons}>
                <Link href="/accept">
                  <Image src="/accept.svg" alt="profile" width={25} height={25} />
                </Link>
                <Link href="/reject">
                  <Image src="/reject.svg" alt="profile" width={25} height={25} />
                </Link>
              </div>
            </div>
            <div>
              <div>
                <Image src="/iconMale.png" alt="profile" width={40} height={40} />
                <h5>Username</h5>
              </div>
              <div className={Style.Buttons}>
                <Link href="/accept">
                  <Image src="/accept.svg" alt="profile" width={25} height={25} />
                </Link>
                <Link href="/reject">
                  <Image src="/reject.svg" alt="profile" width={25} height={25} />
                </Link>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
