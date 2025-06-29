"use client";

export async function refreshAccessToken() {
  try {
    const res = await fetch('/api/v1/auth/refresh', {
      method: 'POST',
      credentials: 'include',
    });

    if (res.ok) {
      console.log("✅ Token refreshed");
      return true;
    } else {
      console.warn("⚠️ Refresh failed, redirecting to login...");
      return false;
    }
  } catch (err) {
    console.error("❌ Refresh error:", err);
    return false;
  }
}

export async function fetchWithAuth(url, options = {}) {
  let res = await fetch(url, {
    ...options,
    credentials: 'include',
  });

  if (res.status === 401) {
    const refreshed = await refreshAccessToken();
    if (refreshed) {
      res = await fetch(url, {
        ...options,
        credentials: 'include',
      });
    }
  }

  return res;
}
useEffect(() => {
  const interval = setInterval(() => {
    refreshAccessToken();
  }, 14 * 60 * 1000); 

  return () => clearInterval(interval);
}, []);
