import { SendData } from "../utils/sendData.js";

export async function Logout() {
    try {
        const { status, body } = await SendData('/api/v1/auth/logout', null);

        if (status === 200) {
            console.log("Logout successful:", body);

            // Clear localStorage/sessionStorage or cookies
            localStorage.removeItem("UserInfo");

            // Redirect to home or login
            window.location.href = "/auth/login"; // or "/" if you prefer
        } else {
            console.error("Logout failed with status", status);
        }
    } catch (err) {
        console.error("Logout error:", err);
    }
}
