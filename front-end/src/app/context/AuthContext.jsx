"use client"; // Ensures this is a client-side component

import { createContext, useState, useContext, useEffect } from "react";
import { useRouter } from "next/router";

// 1. Create the Authentication Context
const AuthContext = createContext();

// 2. Custom hook to use the auth context
export const useAuth = () => {
  return useContext(AuthContext);
};

// 3. AuthProvider to manage the authentication state
export const AuthProvider = ({ children }) => {
  const [isLoggedIn, setIsLoggedIn] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const checkAuth = async () => {
      try {
        const response = await fetch("/api/v1/auth/status", {
          credentials: "include",
        });
        const data = await response.json();
        console.log("Auth status:", data);
        
        if (data.authenticated) {
          setIsLoggedIn(true);
        } 

      } catch (error) {
        console.error("Error during auth check:", error);
        setIsLoggedIn(false);
      } finally {
        setLoading(false);
      }
    };

    checkAuth();
  }, []);

  return (
    <AuthContext.Provider value={{ isLoggedIn, loading, setIsLoggedIn }}>
      {children}
    </AuthContext.Provider>
  );
};

// 4. AuthCheck Component to conditionally render content based on authentication status
const AuthCheck = ({ children }) => {
  const { isLoggedIn, loading } = useAuth();
  const router = useRouter();
  
  const [isMounted, setIsMounted] = useState(false);

  useEffect(() => {
    setIsMounted(true);
  }, []);

  useEffect(() => {
    if (isMounted && !loading && !isLoggedIn) {
      router.push("/login");
    }
  }, [isLoggedIn, loading, router, isMounted]);

  if (loading || !isMounted) {
    return <div>Loading...</div>; // Show loading state while checking auth and when component is mounting
  }

  if (!isLoggedIn) {
    return <div>Please log in to access this content.</div>;
  }

  return <>{children}</>;
};

export default AuthCheck;
