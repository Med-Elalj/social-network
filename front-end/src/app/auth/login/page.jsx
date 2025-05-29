"use client"

import { useState } from 'react';
import Styles from './login.module.css';

export default function Login() {
    const [formData, setFormData] = useState({
        login: '',
        password: "",
    });

    const handleChange = (e) => {
        const { name, value } = e.target;
        setFormData({ ...formData, [name]: value });
    };

    const handleSubmit = (e) => {
        e.preventDefault();

        const data = new FormData();
        for (const key in formData) {
            data.append(key, formData[key]);
        }

        // TODO: send `data` to backend via fetch
        console.log('Submitting form...', formData);
    };

    return (
        <form className={Styles.form} onSubmit={handleSubmit}>
            <label htmlFor="login">Email Or Nickname</label>
            <input type="text" name="login" id="login" onChange={handleChange} />

            <label htmlFor="password">Password</label>
            <input type="password" name="password" id="password" onChange={handleChange} />

            <button type="submit">Login</button>
        </form>
    );
}
