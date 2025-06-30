"use client";

import Styles from "./global.module.css";
import Groups from "./components/Groups";
import Friends from "./components/Friends";
import Image from 'next/image';
import { GetData, SendData } from "../../utils/sendData";
import { useEffect, useState } from "react";

export default function Home() {
  const [posts, setPosts] = useState([]);
  useEffect(() => {
    const fetchData = async () => {
      const formData = {
        start: 0,
      };
      const response = await SendData("/api/v1/get/posts", formData);
      const Body = await response.json();
      if (response.status !== 200) {
        console.log(Body);
      } else {
        setPosts(Body.posts);
        console.log('Posts fetched successfully!');
      }
    };

    fetchData();
  }, []);

  return (
    <div className={Styles.global}>
      {/* Left Sidebar */}
      <div className={Styles.firstSide}>
        <Groups />
      </div>

      {/* Center Content */}
      <div className={Styles.centerContent}>
        {posts && posts.map((Post) => (
          <div key={Post.ID} className={Styles.post} style={{ width: '80%', marginLeft: '10%' }}>
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
                src={`${Post.ImagePath.String}`}
                alt="post"
                width={250}
                height={200}
                sizes="(max-width: 768px) 100vw, 250px"
                style={{ height: 'auto', width: '100%', borderRadius: '10px' }}
              />
            ) : ""
            }

            <section className={Styles.footer}>
              <div className={Styles.action}>
                <Image src="/Like2.svg" alt="like" width={20} height={20} />
                <p>{Post.LikeCount}</p>
              </div>
              <div className={Styles.action}>
                <Image src="/comment.svg" alt="comment" width={20} height={20} />
                <p>{Post.CommentCount}</p>
              </div>
            </section>
          </div>
        ))}
      </div>


      {/* Right Sidebar */}
      <div className={Styles.thirdSide} >
        <Friends />
      </div>
    </div>
  );
}
