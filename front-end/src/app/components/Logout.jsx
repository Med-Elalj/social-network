import { SendData } from "@/app/sendData.js";
import { useAuth } from "@/app/context/AuthContext.jsx";
import { showNotification } from "../utils";

export async function LogoutAndRedirect(router) {
  const { setIsLoggedIn } = useAuth();

  try {
    const response = await SendData("/api/v1/auth/logout", null);

    if (response.status === 200) {
      setIsLoggedIn(false);
      router.push("/login");
    } else {
      const body = await response.json();
      showNotification(body?.message || "Logout failed", "error");
    }
  } catch (err) {
    showNotification("Logout failed", "error");
  }
}
