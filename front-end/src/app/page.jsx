"use client";

import Styles from "./global.module.css";
import Groups from "./components/Groups";
import Friends from "./components/Friends";
import Image from "next/image";
import { SendData } from "../../utils/sendData";
import { useEffect, useState } from "react";
import LikeDeslike, { TimeAgo } from "./utils.jsx";
import Comments from "./comments.jsx";

export default function Home() {
  const [openComments, setOpenComments] = useState(null);
  const [posts, setPosts] = useState([]);
  const [lastPostID, setLastPostID] = useState(0);
  const [hasMore, setHasMore] = useState(true);
  const [loading, setLoading] = useState(false);

  const fetchData = async (reset = false) => {
    if (loading || (!hasMore && !reset)) return;
    setLoading(true);

    let startID = lastPostID;
    if (reset) {
      startID = 0;
    }

    const formData = { start: startID };

    try {
      const response = await SendData("/api/v1/get/posts", formData);
      const Body = await response.json();

      if (response.status !== 200) {
        console.log("Fetch error:", Body);
        setLoading(false);
        return;
      }

      const newPosts = Body.posts || [];

      if (reset) {
        setPosts(newPosts);
      } else {
        setPosts((prev) => {
          const combined = [...prev, ...newPosts];
          const unique = Array.from(new Map(combined.map((p) => [p.ID, p])).values());
          return unique;
        });
      }

      if (newPosts.length === 0) {
        setHasMore(false);
      } else {
        setLastPostID(newPosts[newPosts.length - 1].ID);
        setHasMore(true);
      }
    } catch (err) {
      console.error("Fetch exception:", err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchData(true);
  }, []);

  useEffect(() => {
    const onScroll = () => {
      const nearBottom = window.innerHeight + window.scrollY >= document.body.offsetHeight - 100;
      if (nearBottom && hasMore && !loading) {
        setTimeout(() => {
          fetchData();
        }, 1000);
      }
    };

    window.addEventListener("scroll", onScroll);
    return () => window.removeEventListener("scroll", onScroll);
  }, [hasMore, loading]);

  return (
    <div className={Styles.global}>
      {/* Left Sidebar */}
      <div className={Styles.firstSide}>
        <Groups />
      </div>

      {/* Center Content */}
      <div className={Styles.centerContent}>
        {posts &&
          posts.map((Post) => (
            <div key={Post.ID} className={Styles.post}>
              <section className={Styles.userinfo}>
                <div className={Styles.user}>
                  {/* Main Avatar */}
                  <Image
                    src={Post.AvatarGroup?.String ? `/${Post.AvatarGroup.String}` : "/iconMale.png"}
                    alt="avatar"
                    width={25}
                    height={25}
                  />

                  {/* Texts block */}
                  <div className={Styles.texts}>
                    {Post.GroupId?.Valid ? (
                      <>
                        <p>{Post.GroupName.String}</p>
                        <div className={Styles.user}>
                          <Image
                            src={
                              Post.AvatarUser.Valid ? `${Post.AvatarUser.String}` : "/iconMale.png"
                            }
                            alt="avatar"
                            width={20}
                            height={20}
                          />
                          <p>{Post.UserName}</p>
                        </div>
                      </>
                    ) : (
                      <>
                        <p>{Post.UserName}</p>
                        <div className={Styles.user}></div>
                      </>
                    )}
                  </div>
                </div>

                {/* Timestamp */}
                <div>
                  <p>{TimeAgo(Post.CreatedAt)}</p>
                </div>
              </section>

              <section className={Styles.content}>{Post.Content}</section>

              {/* Post Image (optional) */}
              {Post.ImagePath?.String && (
                <Image
                  src={Post.ImagePath.String}
                  alt="post"
                  width={250}
                  height={200}
                  sizes="(max-width: 768px) 100vw, 250px"
                  style={{ height: "auto", width: "100%", borderRadius: "10px" }}
                  unoptimized
                />
              )}

              <section className={Styles.footer}>
                {/* TODO:add to websocket to be updated for all users */}
                <LikeDeslike
                  EntityID={Post.ID}
                  EntityType={"post"}
                  isLiked={Post.IsLiked}
                  currentLikeCount={Post.LikeCount}
                />

                <div className={Styles.action} onClick={() => setOpenComments(Post.ID)}>
                  <Image src="/comment.svg" alt="comment" width={20} height={20} />
                  <p>{Post.CommentCount}</p>
                </div>

                {openComments === Post.ID && (
                  <div className={Styles.commentPopup} onClick={() => setOpenComments(null)}>
                    <div onClick={(e) => e.stopPropagation()}>
                      <Comments Post={Post} onClose={() => setOpenComments(null)} />
                    </div>
                  </div>
                )}
              </section>
            </div>
          ))}
      </div>

      {/* Right Sidebar */}
      <div className={Styles.thirdSide}>
        <Friends />
      </div>
    </div>
  );
}
