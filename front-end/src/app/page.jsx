"use client";

import Styles from "./global.module.css";
import Groups from "./components/Groups";
import Friends from "./components/Friends";
import Image from "next/image";
import { SendData } from "./sendData.js";
import { useEffect, useState } from "react";
import LikeDeslike, { TimeAgo } from "./utils.jsx";
import Comments from "./comments.jsx";
import { useAuth } from "./context/AuthContext";

export default function Home() {
  const [openComments, setOpenComments] = useState(null);
  const [posts, setPosts] = useState([]);
  const [lastPostID, setLastPostID] = useState(0);
  const [hasMore, setHasMore] = useState(true);
  const [loading, setLoading] = useState(false );

  const { isloading, isLoggedIn } = useAuth();

  /* ---------------- data fetching (unchanged) ---------------- */
  const fetchData = async (reset = false) => {
    if (loading || (!hasMore && !reset)) return;
    if ( isloading || !isLoggedIn) return;
    setLoading(true);

    try {
      const startID = reset ? 0 : lastPostID;
      const response = await SendData("/api/v1/get/posts", {
        start: startID,
        fetch: "home",
      });
      const Body = await response.json();

      if (response.status !== 200) throw Body;

      const newPosts = Body.posts ?? [];

      setPosts((prev) => {
        const combined = reset ? newPosts : [...prev, ...newPosts];
        return Array.from(new Map(combined.map((p) => [p.ID, p])).values());
      });

      setLastPostID(newPosts.at(-1)?.ID ?? lastPostID);
      setHasMore(newPosts.length > 0);
    } catch (err) {
      console.error("Fetch error:", err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    if (!isLoggedIn) return;
    fetchData(true);
  }, [isLoggedIn]);

  useEffect(() => {
    const onScroll = () => {
      if (
        window.innerHeight + window.scrollY >=
          document.body.offsetHeight - 100 &&
        hasMore &&
        !loading
      )
        setTimeout(fetchData, 1_000);
    };
    window.addEventListener("scroll", onScroll);
    return () => window.removeEventListener("scroll", onScroll);
  }, [hasMore, loading]);

  /* ---------------- render ---------------- */
  return (
    <div className={Styles.global}>
      <div className={Styles.firstSide}>
        <Groups />
      </div>

      <div className={Styles.centerContent}>
        {posts.map((Post) => {
          /* ► ONE canonical avatar for the author ◄ */
          const authorAvatar = Post?.AvatarUser?.String
            ? `${Post.AvatarUser.String}` // or full URL if stored externally
            : "/iconMale.png";

          return (
            <div key={Post.ID} className={Styles.post}>
              {/* ---------- header ---------- */}
              <section className={Styles.userinfo}>
                <div className={Styles.user}>
                  {/* left-most avatar or group badge */}
                  {Post.GroupId?.Valid ? (
                    <Image
                      src={
                        Post.AvatarGroup?.String
                          ? `${Post.AvatarGroup.String}`
                          : "/iconGroup.png"
                      }
                      alt="group avatar"
                      width={25}
                      height={25}
                    />
                  ) : (
                    <Image
                      src={authorAvatar}
                      alt="author avatar"
                      width={25}
                      height={25}
                    />
                  )}

                  {/* texts block */}
                  <div className={Styles.texts}>
                    {Post.GroupId?.Valid ? (
                      <>
                        {/* group name */}
                        <p>{Post.GroupName.String}</p>

                        {/* author info (small) */}
                        <div className={Styles.user}>
                          <Image
                            src={
                              Post.AvatarUser.Valid
                                ? `${Post.AvatarUser.String}`
                                : "/iconMale.png"
                            }
                            alt="avatar"
                            width={20}
                            height={20}
                          />
                          <p>{Post.UserName}</p>
                        </div>
                      </>
                    ) : (
                      <p>{Post.UserName}</p>
                    )}
                  </div>
                </div>

                {/* Timestamp */}
                <div>
                  <p>{TimeAgo(Post.CreatedAt)}</p>
                </div>
              </section>

              {/* ---------- body ---------- */}
              <section className={Styles.content}>{Post.Content}</section>

              {Post.ImagePath?.String && (
                <Image
                  src={Post.ImagePath.String}
                  alt="post illustration"
                  width={250}
                  height={200}
                  sizes="(max-width:768px) 100vw, 250px"
                  style={{
                    width: "100%",
                    height: "auto",
                    borderRadius: "10px",
                  }}
                  unoptimized
                />
              )}

              {/* ---------- footer ---------- */}
              <section className={Styles.footer}>
                <LikeDeslike
                  EntityID={Post.ID}
                  EntityType="post"
                  isLiked={Post.IsLiked}
                  currentLikeCount={Post.LikeCount}
                />

                <div
                  className={Styles.action}
                  onClick={() =>
                    setOpenComments((open) =>
                      open === Post.ID ? null : Post.ID
                    )
                  }
                >
                  <Image
                    src="/comment.svg"
                    alt="comment"
                    width={20}
                    height={20}
                  />
                  <p>{Post.CommentCount}</p>
                </div>
              </section>

              {/* ---------- comments popup ---------- */}
              {openComments === Post.ID && (
                <div
                  className={Styles.commentPopup}
                  onClick={() => setOpenComments(null)}
                >
                  <div onClick={(e) => e.stopPropagation()}>
                    <Comments
                      Post={Post}
                      onClose={() => setOpenComments(null)}
                    />
                  </div>
                </div>
              )}
            </div>
          );
        })}
      </div>

      <div className={Styles.thirdSide}>
        <Friends />
      </div>
    </div>
  );
}
