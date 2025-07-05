// "use client";

// import Style from "./profile.module.css";
// import Image from "next/image";
// import Link from "next/link";
// import { useEffect, useState } from "react";
// import Posts from "./[tab]/Posts";
// import Following from "./[tab]/Following";
// import Followers from "./[tab]/Followers";
// import { GetData } from "../../../utils/sendData.js";
// import { CapitalizeFirstLetter } from "../utils.jsx";
// import { useRouter } from "next/navigation";
// import { url } from "inspector";

// export default function PubProfile() {
//   const [activeTab, setActiveTab] = useState("info");
//   const [activeSection, setActiveSection] = useState("posts");
//   const [profileData, setProfileData] = useState(null);

//   const router = useRouter();
//     // print url
//     console.log("🔗 Current URL:", window.location.href);
//   useEffect(() => {

//     console.log("🔄 Fetching profile data for:", router.query.nickname);
//     const fetchAndRedirect = async () => {
//       const res = await GetData("/api/v1/profile/"+router.query.nickname, {
//         credentials: "include",
//       });
//       if (res.ok) {
//         const data = await res.json();
//         router.replace(`/profile/@${data.display_name}`);
//         setProfileData(data);
//       }
//     };
    
//     fetchAndRedirect();
//   }, [router]);

//   return (
//     <div className={Style.container}>
//       <div className={Style.header}>
//         <Image src="/db.png" fill alt="cover" />
//       </div>
      
//       <div className={Style.body}>
//         <div className={Style.first}>
//           <div className={Style.ProfileInfo}>
//             <div className={Style.top}>
//               <div style={{ position: "relative", width: "200px", height: "200px" }}>
//                 <Image src="/db.png" alt="cover" fill />
//               </div>
//               <h4>@{CapitalizeFirstLetter(profileData?.display_name)}</h4>
//             </div>

//             <div className={Style.tabs}>
//               <button onClick={() => setActiveTab("info")}>Info</button>
//               <button onClick={() => setActiveTab("connections")}>Connections</button>
//             </div>

//             {activeTab === "info" && (
//               <>
//                 <div className={Style.center}>
//                   <span>
//                     <h5>About me:</h5>&nbsp;&nbsp;
//                     <h5>
//                       {profileData?.description || "No description provided."}
//                     </h5>
//                   </span>
//                 </div>
//                 <div className={Style.center}>
//                   <span>
//                     <h5>Full Name:</h5>&nbsp;&nbsp;
//                     <h5>
//                       {profileData?.first_name} {profileData?.last_name}
//                     </h5>
//                   </span>
//                   <span>
//                     <h5>Email:</h5>&nbsp;&nbsp;
//                     <h5>{profileData?.email}</h5>
//                   </span>
//                   <span>
//                     <h5>Age:</h5>&nbsp;&nbsp;
//                     <h5>
//                       {profileData?.date_of_birth 
//                         ? (new Date().getFullYear() - new Date(profileData.date_of_birth).getFullYear()).toString()
//                         : "N/A"}
//                     </h5>
//                   </span>
//                   <span>
//                     <h5>Birthdate:</h5>&nbsp;&nbsp;
//                     <h5>{profileData?.date_of_birth || "N/A"}</h5>
//                   </span>
//                 </div>
//               </>
//             )}

//             {activeTab === "connections" && (
//               <div className={Style.numbers}>
//                 <span onClick={() => setActiveSection("posts")}>
//                   <h4>Posts</h4>
//                   <h5>0</h5>
//                 </span>
//                 <span onClick={() => setActiveSection("followers")}>
//                   <h4>Followers</h4>
//                   <h5>0</h5>
//                 </span>
//                 <span onClick={() => setActiveSection("following")}>
//                   <h4>Following</h4>
//                   <h5>0</h5>
//                 </span>
//               </div>
//             )}
//           </div>
//         </div>

//         <div className={Style.second}>
//           {activeSection === "posts" && <Posts />}
//           {activeSection === "followers" && <Followers />}
//           {activeSection === "following" && <Following />}
//         </div>

//         <div className={Style.end}>
//           <div className={Style.requists}>
//             <h3>Suggestions</h3>
//             {[...Array(5)].map((_, i) => (
//               <div key={i}>
//                 <div>
//                   <Image
//                     src="/db.png"
//                     alt="profile"
//                     width={40}
//                     height={40}
//                     style={{ borderRadius: "50%" }}
//                   />
//                   <h5>username</h5>
//                 </div>
//                 <Link href="/addUser">
//                   <Image src="/addUser.svg" alt="profile" width={25} height={25} />
//                 </Link>
//               </div>
//             ))}
//           </div>

//           <div className={Style.requists}>
//             {[...Array(4)].map((_, i) => (
//               <div key={i}>
//                 <div>
//                   <Image src="/iconMale.png" alt="profile" width={40} height={40} />
//                   <h5>Username</h5>
//                 </div>
//                 <div className={Style.Buttons}>
//                   <Link href="/accept">
//                     <Image src="/accept.svg" alt="profile" width={25} height={25} />
//                   </Link>
//                   <Link href="/reject">
//                     <Image src="/reject.svg" alt="profile" width={25} height={25} />
//                   </Link>
//                 </div>
//               </div>
//             ))}
//           </div>
//         </div>
//       </div>
//     </div>
//   );
// }