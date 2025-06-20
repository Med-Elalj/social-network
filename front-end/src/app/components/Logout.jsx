import { SendData } from "../../../utils/sendData.js";

export async function LogoutAndRedirect(router) {
    try {
        const response = await SendData('/api/v1/auth/logout', null);
        if (response.status === 200) {
            console.log("Logout successful:", response.body);
            localStorage.removeItem("UserInfo");
            router.push("/login");
        } else {
            console.error("Logout failed with status", response.status);
        }
    } catch (err) {
        console.error("Logout error:", err);
    }
}
