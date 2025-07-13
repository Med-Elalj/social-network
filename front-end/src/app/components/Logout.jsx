import { SendData } from "@/app/sendData.js";
// import { useAuth } from "../context/AuthContext.jsx";

export async function LogoutAndRedirect(router) {
    // const { setIsLoggedIn } = useAuth();
    try {
        const response = await SendData('/api/v1/auth/logout', null);
        if (response.status === 200) {
            console.log("Logout successful:", response.body);
            localStorage.removeItem("UserInfo");
            // setIsLoggedIn(false);
            router.push("/login");
        } else {
            console.error("Logout failed with status", response.status);
        }
    } catch (err) {
        console.error("Logout error:", err);
    }
}
