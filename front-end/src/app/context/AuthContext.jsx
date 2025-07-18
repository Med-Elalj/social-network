"use client";

import { createContext, useContext, useState, useEffect, useRef, useMemo } from "react";

import { BACKEND_URL } from "@/app/sendData.js";
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
  const didFetchRef = useRef(false);

useEffect(() => {
  if (didFetchRef.current) return; // prevent double-fetch
  didFetchRef.current = true;

  const checkAuth = async () => {
    try {
      const response = await fetch(BACKEND_URL+"/api/v1/auth/status", {
        method: "GET",
        credentials: "include",
      });
      const data = await response.json();
      setIsLoggedIn(data.authenticated);
    } catch (err) {
      console.error(err);
      setIsLoggedIn(false);
    } finally {
      setLoading(false);
    }
  };

  checkAuth();
}, []);

console.log("AuthProvider: checking auth status...");

  useEffect(() => {
    const checkAuth = async () => {
      try {
        const response = await fetch(BACKEND_URL+"/api/v1/auth/status", {
          credentials: "include",
          method: "GET",
        });
        const data = await response.json();

        if (response.ok) {
          setIsLoggedIn(data.authenticated);
        } else {
          setIsLoggedIn(false);
        }
      } catch (error) {
        console.error("Auth error:", error);
        setIsLoggedIn(false);
      } finally {
        setLoading(false);
      }
    };

    checkAuth();
  }, []);

  const value = useMemo(
    () => ({ isLoggedIn, loading, setIsLoggedIn }),
    [isLoggedIn, loading]
  );

  return (
    <AuthContext.Provider value={value}>
      {children}
    </AuthContext.Provider>
  );
};
