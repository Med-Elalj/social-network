// "use client";

// import Image from "next/image";
// import Styles from "./global.module.css";
// import LikeDeslike from "./utils.jsx";
// import { useEffect, useState } from "react";
// import { SendData } from "../../utils/sendData.js";

// export default function Comments({ Post, onClose }) {
//     const [content, setContent] = useState("");
//     const [Comments, setComments] = useState([]);
//     const [LastCommentID, setLastCommentID] = useState(0);
//     const [hasMore, setHasMore] = useState(true);
//     const [loading, setLoading] = useState(false);

//     const fetchData = async (reset = false) => {
//         if (loading || (!hasMore && !reset)) return;
//         setLoading(true);

//         let startID = LastCommentID;
//         if (reset) {
//             startID = 0;
//         }

//         const formData = {
//             post_id: Post.ID,
//             start: startID
//         };
//         try {
//             const response = await SendData("/api/v1/get/comments", formData);

//             if (response.status !== 200) {
//                 const errorBody = await response.text();
//                 console.error("Error fetching comments:", errorBody);
//                 setLoading(false);
//             }

//             const Body = await response.json();
//             const newComments = Body.comments || [];

//             if (reset) {
//                 setComments(newComments);
//                 console.log(Body.message);
//                 console.log('Comments fetched successfully!');
//             } else {
//                 setComments((prev) => {
//                     const combined = [...prev, ...newComments];
//                     const unique = Array.from(new Map(combined.map(p => [p.ID, p])).values());
//                     return unique;
//                 });
//             }

//             if (newComments.length === 0) {
//                 setHasMore(false);
//             } else {
//                 setLastCommentID(newComments[newComments.length - 1].ID);
//                 setHasMore(true);
//             }
//         } catch (err) {
//             console.error("Fetch exception:", err);
//         } finally {
//             setLoading(false);
//         }
//     };

//     useEffect(() => { fetchData(true); }, []);

//     useEffect(() => {
//         const onScroll = () => {
//             const nearBottom = commentPopupContent.innerHeight + commentPopupContent.scrollY >= document.body.offsetHeight - 100;
//             if (nearBottom && hasMore && !loading) {
//                 setTimeout(() => {
//                     fetchData();
//                 }, 1000);
//             }
//         };

//         commentPopupContent.addEventListener("scroll", onScroll);
//         return () => commentPopupContent.removeEventListener("scroll", onScroll);
//     }, [hasMore, loading]);

//     const handleSubmit = async (e) => {
//         e.preventDefault();

//         if (!content.trim()) {
//             console.log("Content is required.");
//             return;
//         }

//         const formData = {
//             post_id: Post.ID,
//             content: content,
//         };

//         const response = await SendData("/api/v1/set/comment", formData);
//         if (response.status !== 200) {
//             const errorBody = await response.json();
//             console.log(errorBody);
//         } else {
//             setContent("");
//             fetchData();
//             const Body = await response.json();
//             console.log(Body.message);
//         }
//     };

//     return (
//         <div className={Styles.commentPopup}>
//             <button className={Styles.closeBtn} onClick={onClose}>
//                 <Image src="/exit.svg" alt="exit" width={30} height={30} />
//             </button>
//             <div className={Styles.commentPopupContent}>
//                 <section className={Styles.postSmall}>
//                     <div className={Styles.userInfo}>
//                         <div className={Styles.first}>
//                             <Image
//                                 src={Post.AvatarGroup?.String ? `/${Post.AvatarGroup.String}` : '/iconMale.png'}
//                                 alt="avatar"
//                                 width={25}
//                                 height={25}
//                             />

//                             <div>
//                                 {Post.GroupId?.Valid ? (
//                                     <>
//                                         <p>{Post.GroupName.String}</p>
//                                         <div className={Styles.user}>
//                                             <Image
//                                                 src={Post.AvatarUser?.String ? `/${Post.Avatar.String}` : '/iconMale.png'}
//                                                 alt="avatar"
//                                                 width={20}
//                                                 height={20}
//                                             />
//                                             <p>{Post.UserName}</p>
//                                         </div>
//                                     </>
//                                 ) : (
//                                     <>
//                                         <p>{Post.UserName}</p>
//                                         <div className={Styles.user}></div>
//                                     </>
//                                 )}
//                             </div>
//                         </div>
//                         <div>
//                             <p>{Post.CreatedAt.replace('T', ' ').slice(0, -1)}</p>
//                         </div>
//                     </div>

//                     <div className={Styles.postContent}>
//                         <p>{Post.Content}</p>
//                         <Image src="/comment.svg" alt="comment" width={30} height={30} />
//                     </div>
//                 </section>


//                 <div className={Styles.comments}>
//                     {Comments && Comments.map((Comment, i) => (
//                         <div key={i} className={Styles.comment}>
//                             <div className={Styles.first}>
//                                 <div>
//                                     <Image
//                                         src={Comment.AvatarUser.String ? `/${Comment.AvatarUser.String}` : "/iconMale.png"}
//                                         alt="avatar"
//                                         width={25}
//                                         height={25}
//                                     />
//                                     <p>{Comment.Author}</p>
//                                 </div>
//                                 <div>
//                                     <p>{Post.CreatedAt.replace('T', ' ').slice(0, -1)}</p>
//                                 </div>
//                             </div>
//                             <div className={Styles.texts}>
//                                 <p>{Comment.Content}</p>
//                             </div>
//                             <div className={Styles.like}>
//                                 <LikeDeslike
//                                     EntityID={Comment.ID}
//                                     EntityType={"comment"}
//                                     isLiked={Comment.IsLiked}
//                                     currentLikeCount={Comment.LikeCount}
//                                 />
//                             </div>
//                         </div>
//                     ))}
//                 </div>

//             </div>
//             <form className={Styles.commentCreation} onSubmit={handleSubmit}>
//                 <input type="text"
//                     placeholder="Comment"
//                     value={content}
//                     onChange={(e) => setContent(e.target.value)}
//                 />
//                 <button type="submit">Send</button>
//             </form>
//         </div>
//     );
// }
import { useEffect, useState, useRef } from "react";
import Image from "next/image";
import Styles from "./global.module.css";
import LikeDeslike from "./utils.jsx";
import { SendData } from "../../utils/sendData.js";

export default function Comments({ Post, onClose }) {
    const [content, setContent] = useState("");
    const [Comments, setComments] = useState([]);
    const [LastCommentID, setLastCommentID] = useState(0);
    const [hasMore, setHasMore] = useState(true);
    const [loading, setLoading] = useState(false);
    const scrollRef = useRef(null); // ✅ ref to scroll container

    const fetchData = async (reset = false) => {
        if (loading || (!hasMore && !reset)) return;
        setLoading(true);

        const startID = reset ? 0 : LastCommentID;
        const formData = { post_id: Post.ID, start: startID };

        try {
            const response = await SendData("/api/v1/get/comments", formData);
            const Body = await response.json();

            const newComments = Body.comments || [];

            setComments(prev => {
                const combined = reset ? newComments : [...prev, ...newComments];
                const unique = Array.from(new Map(combined.map(c => [c.ID, c])).values());
                return unique;
            });

            if (newComments.length === 0) {
                setHasMore(false);
            } else {
                setLastCommentID(newComments[newComments.length - 1].ID);
                setHasMore(true);
            }
        } catch (err) {
            console.error("Error fetching comments:", err);
        } finally {
            setLoading(false);
        }
    };

    useEffect(() => { fetchData(true); }, []);

    // ✅ Scroll pagination inside comments box
    useEffect(() => {
        const el = scrollRef.current;
        if (!el) return;

        const onScroll = () => {
            const nearBottom = el.scrollTop + el.clientHeight >= el.scrollHeight - 100;
            if (nearBottom && hasMore && !loading) {
                setTimeout(() => {
                    fetchData();
                }, 1000);
            }
        };

        el.addEventListener("scroll", onScroll);
        return () => el.removeEventListener("scroll", onScroll);
    }, [hasMore, loading]);

    const handleSubmit = async (e) => {
        e.preventDefault();
        if (!content.trim()) return;

        const formData = { post_id: Post.ID, content };

        try {
            const response = await SendData("/api/v1/set/comment", formData);
            const Body = await response.json();

            if (response.status !== 200) {
                console.log("Comment error:", Body);
            } else {
                setContent("");
                fetchData(true); // reset + reload comments
                console.log(Body.message);
            }
        } catch (err) {
            console.error("Submit error:", err);
        }
    };

    return (
        <div className={Styles.commentPopup}>
            <button className={Styles.closeBtn} onClick={onClose}>
                <Image src="/exit.svg" alt="exit" width={30} height={30} />
            </button>

            <div className={Styles.commentPopupContent} ref={scrollRef}>
                {/* Post preview */}
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

                {/* Comments list */}
                <div className={Styles.comments}>
                    {Comments.map((Comment) => (
                        <div key={Comment.ID} className={Styles.comment}>
                            <div className={Styles.first}>
                                <div>
                                    <Image
                                        src={Comment.AvatarUser?.String ? `/${Comment.AvatarUser.String}` : "/iconMale.png"}
                                        alt="avatar"
                                        width={25}
                                        height={25}
                                    />
                                    <p>{Comment.Author}</p>
                                </div>
                                <div>
                                    <p>{Comment.CreatedAt.replace('T', ' ').slice(0, -1)}</p>
                                </div>
                            </div>
                            <div className={Styles.texts}>
                                <p>{Comment.Content}</p>
                            </div>
                            <div className={Styles.like}>
                                <LikeDeslike
                                    EntityID={Comment.ID}
                                    EntityType="comment"
                                    isLiked={Comment.IsLiked}
                                    currentLikeCount={Comment.LikeCount}
                                />
                            </div>
                        </div>
                    ))}
                </div>
            </div>

            {/* Add new comment */}
            <form className={Styles.commentCreation} onSubmit={handleSubmit}>
                <input
                    type="text"
                    placeholder="Comment"
                    value={content}
                    onChange={(e) => setContent(e.target.value)}
                />
                <button type="submit">Send</button>
            </form>
        </div>
    );
}
