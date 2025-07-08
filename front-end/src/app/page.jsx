"use client";

import Styles from "./global.module.css";
import Groups from "./components/Groups";
import Friends from "./components/Friends";
import Image from 'next/image';
import { SendData } from "../../utils/sendData";
import { useEffect, useState } from "react";
import LikeDeslike from "./utils.jsx";
import Comments from "./comments.jsx";

export default function Home() {
  const [posts, setPosts] = useState([]);
  const [openComments, setOpenComments] = useState(null);
  const [lastPostID, setLastPostID] = useState(0);
  {/*const [isLoading, setIsLoading] = useState(false);*/}
  {/*const [hasMore, setHasMore] = useState(true);*/}

  const fetchData = async () => {
    {/*if (isLoading || !hasMore) return;*/}
    {/*setIsLoading(true);*/}

    const formData = { start: lastPostID };
    const response = await SendData("/api/v1/get/posts", formData);
    const Body = await response.json();

    if (response.status !== 200) {
      console.log(Body);
    } else {
      {/*const newPosts = Body.posts;
      if (!newPosts || newPosts.length === 0) {
        setHasMore(false);
      } else {
        setPosts((prev) => [...prev, ...newPosts]);
        const newLastID = newPosts[newPosts.length - 1].ID;
        setLastPostID(newLastID);
      }*/}
      setPosts(Body.posts);
    }

    {/*setIsLoading(false);*/}
  };

  useEffect(() => { fetchData(); }, [lastPostID]);

  {/*useEffect(() => {
    const handleScroll = () => {
      const nearBottom = window.innerHeight + window.scrollY >= document.body.offsetHeight - 50;
      console.log(nearBottom && !isLoading && hasMore);

      if (nearBottom && !isLoading && hasMore) {
        fetchData();
      }
    };

    window.addEventListener("scroll", handleScroll);
    return () => window.removeEventListener("scroll", handleScroll);
  }, [isLoading, hasMore]);*/}

  return (
    <div className={Styles.global}>
      {/* Left Sidebar */}
      <div className={Styles.firstSide}>
        <Groups />
      </div>

      {/* Center Content */}
      <div className={Styles.centerContent}>
        {posts && posts.map((Post) => (
          <div key={Post.ID} className={Styles.post}>
            <section className={Styles.userinfo}>
              <div className={Styles.user}>
                {/* Main Avatar */}
                <Image
                  src={Post.AvatarGroup?.String ? `/${Post.AvatarGroup.String}` : '/iconMale.png'}
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
                          src={Post.AvatarUser?.String ? `/${Post.Avatar.String}` : '/iconMale.png'}
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
                <p>{Post.CreatedAt.replace('T', ' ').slice(0, -1)}</p>
              </div>
            </section>

            <section className={Styles.content}>
              {Post.Content}
            </section>

            {/* Post Image (optional) */}
            {Post.ImagePath?.String ? (
              <Image
                src={`/${Post.ImagePath.String}`}
                alt="post"
                width={250}
                height={200}
                sizes="(max-width: 768px) 100vw, 250px"
                style={{ height: 'auto', width: '100%', borderRadius: '10px' }}
              />
            ) : ""
            }

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
                  <div onClick={e => e.stopPropagation()}>
                    <Comments Post={Post} onClose={() => setOpenComments(null)} />
                  </div>
                </div>
              )}
            </section>
          </div>
        ))}
        {/*{isLoading && (
          <div style={{ textAlign: "center", margin: "20px" }}>
            <p>Loading more posts...</p>
          </div>
        )}*/}
      </div>


      {/* Right Sidebar */}
      <div className={Styles.thirdSide} >
        <Friends />
      </div>
    </div>
  );
}
