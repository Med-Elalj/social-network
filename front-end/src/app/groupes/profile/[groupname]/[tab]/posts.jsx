import { useEffect, useState } from "react";
import { SendData } from "@/app/sendData.js";
import Styles from "../profile.module.css";
import Image from "next/image";
import LikeDeslike, { TimeAgo } from "@/app/utils.jsx";
import Comments from "@/app/comments.jsx";

export default function Posts({ activeSection, setActiveSection, groupId }) {
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

        const formData = { start: startID, groupId: groupId, fetch: "group" };

        try {
            const response = await SendData("/api/v1/get/posts", formData);
            const Body = await response.json();

            if (response.status !== 200) {
                console.error("Fetch error:", Body);
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
        <div>
            <button
                className={Styles.CreatePostBtn}
                onClick={() => setActiveSection("createPost")}
            >
                Create Post
            </button>
            {posts ?
                posts.map((Post) => {
                    /* ► ONE canonical avatar for the author ◄ */
                    const authorAvatar = Post?.AvatarUser?.String
                        ? `${Post.AvatarUser.String}`  // or full URL if stored externally
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
                                            style={{ borderRadius: "50%" }}
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
                                    style={{ width: "100%", height: "auto", borderRadius: "10px" }}
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
                                        setOpenComments((open) => (open === Post.ID ? null : Post.ID))
                                    }
                                >
                                    <Image src="/comment.svg" alt="comment" width={20} height={20} />
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
                                        <Comments Post={Post} onClose={() => setOpenComments(null)} />
                                    </div>
                                </div>
                            )}
                        </div>
                    );
                })
                : <p>Go Join Groupes.</p>}
        </div>
    )
}