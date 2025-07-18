"use client";

import Image from "next/image";
import Link from "next/link";
import Styles from "../groups.module.css";
import { useState, useEffect } from "react";
import { GetData } from "../../sendData.js";
import Comments from "@/app/comments.jsx";
import LikeDeslike from "@/app/utils.jsx";

export default function GroupPosts() {
    const [posts, setPosts] = useState([]);
    const [openComments, setOpenComments] = useState(null);


  useEffect(() => {
    const fetchData = async () => {
      const response = await GetData("/api/v1/get/groupFeeds");
      const body = await response.json();

      if (response.status !== 200) {
        console.error(body);
      } else {
        setPosts(body.posts);
        console.log("Posts fetched successfully!");
      }
    };

    fetchData();
  }, []);

  return (
    <div>
      {posts ? (
        posts.map((Post, i) => (
          <div key={i} className={Styles.post}>
            <section className={Styles.userinfo}>
              <div className={Styles.user}>
                <Image
                  src={Post.AvatarGroup?.String ? `${Post.AvatarGroup.String}` : "/iconMale.png"}
                  alt="avatar"
                  width={25}
                  height={25}
                />
                <div className={Styles.texts}>
                  {Post.GroupId?.Valid ? (
                    <>
                      <p>{Post.GroupName.String}</p>
                      <div className={Styles.user}>
                        <Image
                          src={
                            Post.AvatarUser?.String ? `${Post.AvatarUser.String}` : "/iconMale.png"
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
              <div>
                <p>{Post.CreatedAt.replace("T", " ").slice(0, -1)}</p>
              </div>
            </section>

            <section className={Styles.content}>{Post.Content || ""}</section>

            {Post.ImagePath?.String && (
              <Image
                src={`${Post.ImagePath.String}`}
                alt="post image"
                width={250}
                height={200}
                sizes="(max-width: 768px) 100vw, 250px"
                style={{ height: "auto", width: "100%", borderRadius: "10px" }}
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

              <div className={Styles.action}>
                <Image src="/comment.svg" alt="comment" width={20} height={20} />
                <p>{Post.CommentCount ?? 0}</p>
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
        ))
      ) : (
        <div className={Styles.noPosts}>
          <h3>Join groups to see feeds</h3>
          <Link href="/groupes/discover">Descover groups</Link>
        </div>
      )}
    </div>
  );
}
