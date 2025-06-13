import { SendData } from "../utils/sendData.js";

export async function Logout() {
    try {
        const { status, body } = await SendData('/api/v1/auth/logout', null);

        if (status === 200) {
            console.log("Logout successful:", body);
            localStorage.removeItem("UserInfo");
            window.location.href = "/auth/login";
        } else {
            console.error("Logout failed with status", status);
        }
    } catch (err) {
        console.error("Logout error:", err);
    }
}
