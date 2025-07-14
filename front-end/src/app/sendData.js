"use client";

const BACKEND_URL = "http://localhost:8080";

export async function refreshAccessToken() {

  console.log("🔄 Attempting to refresh access token...");
  try {
    const res = await fetch(BACKEND_URL + '/api/v1/auth/refresh', {
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

async function fetchWithAuth(path, options = {}) {
  let url = BACKEND_URL + path;
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

export async function SendData(path, data) {
  return await fetchWithAuth( path, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(data),
  });
}
// export async function  SendData (path, data) {
//   SendAuthData(path, data)
// } 

export async function GetData(path, params = {}) {
  // Build query string if any params are passed
  const query = new URLSearchParams(params).toString();
  const fullUrl = query ? `${path}?${query}` : path;

  return await fetchWithAuth(fullUrl, {
    method: 'GET',
    headers: { 'Content-Type': 'application/json' },
  });
}

export async function SendAuthData(path, Data) {   
    try {
        const response = await fetch(BACKEND_URL+path, {
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
