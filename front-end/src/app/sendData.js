"use client";

export async function refreshAccessToken() {

  console.log("🔄 Attempting to refresh access token...");
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

async function fetchWithAuth(url, options = {}) {

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

export async function SendData(url, data) {
  return await fetchWithAuth(url, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  });
}

export async function GetData(url, params = {}) {
  // Build query string if any params are passed
  const query = new URLSearchParams(params).toString();
  const fullUrl = query ? `${url}?${query}` : url;

  return await fetchWithAuth(fullUrl, {
    method: 'GET',
    headers: { 'Content-Type': 'application/json' },
  });
}

export async function SendAuthData(url, Data) {   
    try {
        const response = await fetch(url, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            credentials: 'include',
            body: JSON.stringify(Data),
        })
        
        return response
    } catch (error) {
        return error
    }
}
