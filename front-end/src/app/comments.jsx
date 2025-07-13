import { useEffect, useState, useRef } from "react";
import Image from "next/image";
import Styles from "./global.module.css";
import LikeDeslike from "./utils.jsx";
import { SendData } from "./sendData.js";
import { HandleUpload } from "./utils.jsx";

export default function Comments({ Post, onClose }) {
  const [content, setContent] = useState("");
  const [Comments, setComments] = useState([]);
  const [LastCommentID, setLastCommentID] = useState(0);
  const [hasMore, setHasMore] = useState(true);
  const [loading, setLoading] = useState(false);
  const scrollRef = useRef(null); //ref to scroll container
  const [commentImage, setCommentImage] = useState(null);
  const commentFileRef = useRef(null);

  const handleCommentImageChange = (e) => {
    const file = e.target.files[0];
    if (file) setCommentImage(file);
  };

  const fetchData = async (reset = false) => {
    if (loading || (!hasMore && !reset)) return;
    setLoading(true);

    const startID = reset ? 0 : LastCommentID;
    const formData = { post_id: Post.ID, start: startID };

    try {
      const response = await SendData("/api/v1/get/comments", formData);
      const Body = await response.json();

      const newComments = Body.comments || [];

      setComments((prev) => {
        const combined = reset ? newComments : [...prev, ...newComments];
        const unique = Array.from(new Map(combined.map((c) => [c.ID, c])).values());
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

  useEffect(() => {
    fetchData(true);
  }, []);

  //Scroll pagination inside comments box
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
    if (!content.trim() && !commentImage) return;

    let image_path = null;
    if (commentImage) {
      image_path = await HandleUpload(commentImage);
    }
    const formData = { post_id: Post.ID, content, image_path: image_path };

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

  Comments.map(com => {
    console.log("Comments:", com);
    console.log("img path:", com.image_path)
  })

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
                src={Post.AvatarGroup?.String ? `/${Post.AvatarGroup.String}` : "/iconMale.png"}
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
                        src={Post.AvatarUser?.String ? `/${Post.Avatar.String}` : "/iconMale.png"}
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
          </div>

          <div className={Styles.postContent}>
            <p>{Post.Content}</p>
            {Post.ImagePath.valid ? (
              <Image src={`/${Post.ImagePath}`} alt="comment" width={30} height={30} />
            ) : null}
          </div>
        </section>

        {/* Comments list */}
        <div className={Styles.comments}>
          {Comments.map((Comment) => (
            <div key={Comment.ID} className={Styles.comment}>
              <div className={Styles.first}>
                <div>
                  <Image
                    src={
                      Comment.AvatarUser?.String ? `${Comment.AvatarUser.String}` : "/iconMale.png"
                    }
                    alt="avatar"
                    width={25}
                    height={25}
                  />
                  <p>{Comment.Author}</p>
                </div>
                <div>
                  <p>{Comment.CreatedAt.replace("T", " ").slice(0, -1)}</p>
                </div>
              </div>
              <div className={Styles.texts}>
                <p>{Comment.Content}</p>

                {Comment.image_path?.Valid && (
                  <Image
                    src={
                      Comment.image_path.String.startsWith("/")
                        ? Comment.image_path.String
                        : `/${Comment.image_path.String}`
                    }
                    alt="comment attachment"
                    width={200}
                    height={150}
                    unoptimized
                    style={{ width: "100%", height: "auto", marginTop: "0.5rem" }}
                  />
                )}
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

        <button type="button" onClick={() => commentFileRef.current.click()}>
          Add Image
        </button>
        <input
          ref={commentFileRef}
          id="commentImage"
          type="file"
          accept="image/*"
          style={{ display: "none" }}
          onChange={handleCommentImageChange}
        />

        {commentImage && (
          <div className={Styles.previewContainer}>
            <span className={Styles.fileName}>{commentImage.name}</span>
            <button
              type="button"
              onClick={() => {
                setCommentImage(null);
                commentFileRef.current.value = "";
              }}
            >
              ✕
            </button>
          </div>
        )}

        <button type="submit">Comment</button>
      </form>
    </div>
  );
}
