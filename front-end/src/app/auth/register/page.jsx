"use client";

import { useState } from 'react';
import styles from './register.module.css';

export default function Register() {
    const [formData, setFormData] = useState({
        email: '',
        password: '',
        firstName: '',
        lastName: '',
        nickName: '',
        dob: '',
        about: '',
        profileImg: null,
    });

    const handleChange = (e) => {
        const { name, value, files } = e.target;
        if (name === 'profileImg') {
            setFormData({ ...formData, profileImg: files[0] });
        } else {
            setFormData({ ...formData, [name]: value });
        }
    };

    const handleSubmit = (e) => {
        e.preventDefault();

        const data = new FormData();
        for (const key in formData) {
            data.append(key, formData[key]);
        }

        console.log('Submitting form...', formData);
    };

    return (
        <div>
            <div className={styles.messageBox}>
                <h2>Join Our Social Network 👋</h2>
            </div>
            <div className={styles.formContainer}>
                <form className={styles.form} onSubmit={handleSubmit}>
                    <label className={styles.label} htmlFor="email">Email</label>
                    <input className={styles.input} type="email" name="email" id="email" onChange={handleChange} />

                    <label className={styles.label} htmlFor="password">Password</label>
                    <input className={styles.input} type="password" name="password" id="password" onChange={handleChange} />

                    <label className={styles.label} htmlFor="firstName">First Name</label>
                    <input className={styles.input} type="text" name="firstName" id="firstName" onChange={handleChange} />

                    <label className={styles.label} htmlFor="lastName">Last Name</label>
                    <input className={styles.input} type="text" name="lastName" id="lastName" onChange={handleChange} />

                    <label className={styles.label} htmlFor="nickName">Nickname</label>
                    <input className={styles.input} type="text" name="nickName" id="nickName" onChange={handleChange} />

                    <label className={styles.label} htmlFor="dob">Date of Birth</label>
                    <input className={styles.input} type="date" name="dob" id="dob" onChange={handleChange} />

                    <label className={styles.label} htmlFor="profileImg">Profile Image</label>
                    <input className={styles.input} type="file" name="profileImg" id="profileImg" accept="image/*" onChange={handleChange} />

                    <label className={styles.label} htmlFor="about">About Me</label>
                    <input className={styles.input} type="text" name="about" id="about" onChange={handleChange} />

                    <button className={styles.submitButton} type="submit">Register</button>
                </form>
            </div>
        </div>
    );
}
