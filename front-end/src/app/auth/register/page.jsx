"use client";

import { useState } from 'react';
import styles from './register.module.css';
import { SendData } from '../../../../utils/sendData.js';

export default function Register() {
    const [formData, setFormData] = useState({
        username: '',
        email: '',
        password: '',
        gender: 'male',
        fname: '',
        lname: '',
        birthdate: '',
        avatar: null,
        about: null,
    });

    const handleChange = (e) => {
        const { name, value, files } = e.target;
        if (name === 'avatar') {
            console.log(files);

            setFormData({ ...formData, avatar: files[0].name });
        } else {
            setFormData({ ...formData, [name]: value });
        }
    };

    const handleSubmit = async (e) => {
        e.preventDefault();

        const data = new FormData();
        for (const key in formData) {
            data.append(key, formData[key]);
        }

        console.log('Submitting form...', formData);
        let status
        let res
        await status, res = SendData('/api/v1/auth/register', formData)
        if (status != 200) {
            console.log(res);
        }

    };

    return (
        <div className={styles.container}>
            <div className={styles.messageBox}>
                <h2>Join Our Social Network 👋</h2>
            </div>
            <div className={styles.formContainer}>
                <form className={styles.form} onSubmit={handleSubmit}>
                    <label className={styles.label} htmlFor="email">Email</label>
                    <input className={styles.input} type="email" name="email" id="email" onChange={handleChange} />

                    <label className={styles.label} htmlFor="password">Password</label>
                    <input className={styles.input} type="password" name="password" id="password" onChange={handleChange} />

                    <label className={styles.label} htmlFor="fname">First Name</label>
                    <input className={styles.input} type="text" name="fname" id="firstName" onChange={handleChange} />

                    <label className={styles.label} htmlFor="lname">Last Name</label>
                    <input className={styles.input} type="text" name="lname" id="lastName" onChange={handleChange} />

                    <label className={styles.label} htmlFor="username">Nickname</label>
                    <input className={styles.input} type="text" name="username" id="nickName" onChange={handleChange} />

                    <label className={styles.label} htmlFor="birthdate">Date of Birth</label>
                    <input className={styles.input} type="date" name="birthdate" id="dob" onChange={handleChange} />

                    <label className={styles.label} htmlFor="avatar">Profile Image</label>
                    <input className={styles.input} type="file" name="avatar" id="profileImg" accept="image/*" onChange={handleChange} />

                    <label className={styles.label} htmlFor="about">About Me</label>
                    <input className={styles.input} type="text" name="about" id="about" onChange={handleChange} />

                    <button className={styles.submitButton} type="submit">Register</button>
                </form>
            </div>
        </div>
    );
}
