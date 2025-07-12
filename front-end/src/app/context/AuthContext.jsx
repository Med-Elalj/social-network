// "use client";
// import { createContext, useContext, useEffect, useState } from "react";
// import { GetData } from "../../../utils/sendData.js";
// const { NEXT_PUBLIC_API_URL } = process.env.NEXT_PUBLIC_API_URL;


// const AuthContext = createContext();

// export const AuthProvider = ({ children }) => {
//   const [isLoggedIn, setIsLoggedIn] = useState(null);

//   useEffect(() => {
//     const checkAuth = async () => {
//       try {
//         const response = await GetData(process.env.NEXT_PUBLIC_API_URL + "/auth/status", {
//           credentials: "include",
//         });
//         const data = await response.json();
//         setIsLoggedIn(data.authenticated === true);
//       } catch (err) {
//         setIsLoggedIn(false);
//       }
//     };

//     checkAuth();
//   }, []);

//   return (
//     <AuthContext.Provider value={{ isLoggedIn, setIsLoggedIn }}>
//       {children}
//     </AuthContext.Provider>
//   );
// };

// export const useAuth = () => useContext(AuthContext);
