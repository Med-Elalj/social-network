// "use client";
// import { createContext, useContext, useEffect, useState } from "react";
// import { GetData } from "../../../utils/sendData.js";

// const AuthContext = createContext();

// export const AuthProvider = ({ children }) => {
//   const [isLoggedIn, setIsLoggedIn] = useState(null);

//   useEffect(() => {
//     const checkAuth = async () => {
//       try {
//         const response = await GetData("/api/v1/auth/status", {
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
