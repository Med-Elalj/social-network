"use client";

import Image from "next/image";
import Styles from "./global.module.css";
import LikeDeslike from "./utils.jsx";
import { useEffect, useState } from "react";
import { SendData } from "../../utils/sendData.js";

export default function Comments({ Post, onClose }) {
    const [content, setContent] = useState("");
    const [Comments, setComments] = useState([]);

    const fetchData = async () => {
        const formData = {
            post_id: Post.ID,
            start: 0
        };
        const response = await SendData("/api/v1/get/comments", formData);

        if (response.status !== 200) {
            const errorBody = await response.text();
            console.error("Error fetching comments:", errorBody);
        } else {
            const Body = await response.json();
            setComments(Body.comments);
            console.log(Body);
            console.log('Comments fetched successfully!');
        }
    };

    useEffect(() => { fetchData(); }, []);

    const handleSubmit = async (e) => {
        e.preventDefault();

        if (!content.trim()) {
            console.log("Content is required.");
            return;
        }

        const formData = {
            post_id: Post.ID,
            content: content,
        };

        const response = await SendData("/api/v1/set/comment", formData);
        const Body = await response.json();
        if (response.status !== 200) {
            const errorBody = await response.json();
            console.log(errorBody);
        } else {
            setContent("");
            fetchData();
            console.log(Body.message);
        }
    };

    return (
        <div className={Styles.commentPopup}>
            <button className={Styles.closeBtn} onClick={onClose}>
                <Image src="/exit.svg" alt="exit" width={30} height={30} />
            </button>
            <div className={Styles.commentPopupContent}>
                <section className={Styles.postSmall}>
                    <div className={Styles.userInfo}>
                        <div className={Styles.first}>
                            <Image
                                src={Post.AvatarGroup?.String ? `/${Post.AvatarGroup.String}` : '/iconMale.png'}
                                alt="avatar"
                                width={25}
                                height={25}
                            />

                            <div>
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
                        <div>
                            <p>{Post.CreatedAt.replace('T', ' ').slice(0, -1)}</p>
                        </div>
                    </div>

                    <div className={Styles.postContent}>
                        <p>{Post.Content}</p>
                        <Image src="/comment.svg" alt="comment" width={30} height={30} />
                    </div>
                </section>


                <div className={Styles.comments}>
                    {Comments && Comments.map((Comment, i) => (
                        <div key={i} className={Styles.comment}>
                            <div className={Styles.first}>
                                <div>
                                    <Image
                                        src={Comment.AvatarUser.String ? `/${Comment.AvatarUser.String}` : "/iconMale.png"}
                                        alt="avatar"
                                        width={25}
                                        height={25}
                                    />
                                    <p>{Comment.Author}</p>
                                </div>
                                <div>
                                    <p>{Post.CreatedAt.replace('T', ' ').slice(0, -1)}</p>
                                </div>
                            </div>
                            <div className={Styles.texts}>
                                <p>{Comment.Content}</p>
                            </div>
                            <div className={Styles.like}>
                                <LikeDeslike
                                    EntityID={Comment.ID}
                                    EntityType={"comment"}
                                    isLiked={Comment.IsLiked}
                                    currentLikeCount={Comment.LikeCount}
                                />
                            </div>
                        </div>
                    ))}
                </div>

            </div>
            <form className={Styles.commentCreation} onSubmit={handleSubmit}>
                <input type="text"
                    placeholder="Comment"
                    value={content}
                    onChange={(e) => setContent(e.target.value)}
                />
                <button type="submit">Send</button>
            </form>
        </div>
    );
}
