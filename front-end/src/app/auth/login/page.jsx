"use client";

import { useState } from 'react';
import Styles from "./login.module.css";
import { SendData } from '../../../../utils/sendData.js';

export default function Login() {
    const [formData, setFormData] = useState({
        login: '',
        pwd: '',
    });

    const handleChange = (e) => {
        const { name, value } = e.target;
        setFormData({ ...formData, [name]: value });
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        const data = new FormData();
        for (const key in formData) {
            data.append(key, formData[key]);
        }

        let status
        let res
        await status, res = SendData('/api/v1/auth/login', formData)
    };

    return (
        <div className={Styles.container}>
            <div className={Styles.messageBox}>
                <h2>Welcome back to Our Social Network 👋</h2>
            </div>

            <form className={Styles.form} onSubmit={handleSubmit}>
                <label className={Styles.label} htmlFor="login">Email or Nickname</label>
                <input
                    className={Styles.input}
                    type="text"
                    name="login"
                    id="login"
                    onChange={handleChange}
                    required
                />

                <label className={Styles.label} htmlFor="pwd">Password</label>
                <input
                    className={Styles.input}
                    type="password"
                    name="pwd"
                    id="password"
                    onChange={handleChange}
                    required
                />

                <button className={Styles.button} type="submit">Login</button>
            </form>
        </div>
    );
}
