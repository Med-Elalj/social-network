import { SendData } from "@/app/sendData.js";
import { externalNotification } from "../context/NotificationContext";

export async function LogoutAndRedirect({ router, isLoggedIn, setIsLoggedIn }) {
  if (!isLoggedIn) return;

  try {
    const response = await SendData("/api/v1/auth/logout", null);

    if (response.ok) {
      externalNotification("Logout successful!");
      setIsLoggedIn(false);
      router.push("/login");
    } else {
      const body = await response.json();
      externalNotification(body?.message || "Logout failed", "error");
    }
  } catch (err) {
    externalNotification("Logout failed", "error");
  }
}

