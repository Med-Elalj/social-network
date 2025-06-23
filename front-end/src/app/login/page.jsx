"use client";
import { showNotification, usePasswordToggle } from "../utils";

import { useState } from "react";
import Styles from "./login.module.css";
import { SendData } from "../../../utils/sendData.js";
import { useRouter } from "next/navigation";

export default function Login() {
    usePasswordToggle();
    const router = useRouter();
    const [formData, setFormData] = useState({
        login: "",
        pwd: "",
    });

    const handleChange = (e) => {
        const { name, value } = e.target;
        setFormData({ ...formData, [name]: value });
    };

    const handleSubmit = async (e) => {
        e.preventDefault();

        const response = await SendData("/api/v1/auth/login", formData);
        const res = await response.json();
        if (response.status == 200) {
            console.log("res", res);
            router.push("/");
        } else {
            showNotification(
                res.error || "Login failed. Please try again.",
                "error",
                true,
                5000
            );
            // console.log(res.error || "Login failed. Please try again.");
        }
    };

    return (
        <div className={Styles.container}>
            <div className={Styles.messageBox}>
                <h2>Welcome back to Our Social Network 👋</h2>
            </div>

            <form className={Styles.form} onSubmit={handleSubmit}>
                <label className={Styles.label} htmlFor="login">
                    Email or Nickname
                </label>
                <input
                    className={Styles.input}
                    type="text"
                    name="login"
                    id="login"
                    onChange={handleChange}
                    required
                />

                <label className={Styles.label} htmlFor="pwd">
                    Password
                </label>
                <div className={Styles.inputWrapper}>
                        <input
                            className={Styles.input}
                            type="password"
                            name="pwd"
                            id="password"
                            onChange={handleChange}
                            required
                        />
                        <i className="togglePwd">
                            <span className="icon vis_icon material-symbols-outlined">visibility</span>
                        </i>
                </div>
        
                <button className={Styles.button} type="submit">
                    Login
                </button>
            </form>
        </div>
    );
}
