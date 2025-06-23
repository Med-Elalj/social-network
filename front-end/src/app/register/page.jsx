"use client";

import { useState } from 'react';
import Styles from './register.module.css';
import { SendData } from '../../../utils/sendData.js';
import { showNotification, usePasswordToggle } from '../utils.jsx';

export default function Register() {
    const [formData, setFormData] = useState({
        username: '',
        email: '',
        password: '',
        gender: '',
        fname: '',
        lname: '',
        birthdate: '',
        avatar: null,
        about: null,
    });
    const [previewUrl, setPreviewUrl] = useState(null);
    usePasswordToggle();
    const handleChange = (e) => {
        const { name, value, files } = e.target;
        if (name === 'avatar') {
            setFormData({ ...formData, avatar: files[0].name });
            setPreviewUrl(URL.createObjectURL(files[0]))
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

        const response = await SendData('/api/v1/auth/register', formData);

        if (response.status !== 200) {
            const errorBody = await response.json();
            // console.log(errorBody);
            showNotification(errorBody.error || "Registration failed. Please try again.", "error", true, 5000);
        } else {
            // console.log('Form submitted successfully!');
            showNotification("Registration successful! Welcome to our social network!", "success", true, 5000);
            // load home page 
            // todo: login refers to login
        }
    };


    return (
        <div className={Styles.container}>
            <div className={Styles.messageBox}>
                <h2>Join Our Social Network 👋</h2>
            </div>
            <div className={Styles.formContainer}>
                <form className={Styles.form} onSubmit={handleSubmit}>


                    <label className={Styles.label} htmlFor="fname">First Name</label>
                    <input className={Styles.input} type="text" name="fname" id="firstName" onChange={handleChange} />

                    <label className={Styles.label} htmlFor="lname">Last Name</label>
                    <input className={Styles.input} type="text" name="lname" id="lastName" onChange={handleChange} />
                    
                    <label className={Styles.label} htmlFor="username">Nickname</label>
                    <input className={Styles.input} type="text" name="username" id="nickName" onChange={handleChange} />

                    <label className={Styles.label} htmlFor="email">Email</label>
                    <input className={Styles.input} type="email" name="email" id="email" onChange={handleChange} />
                    
                    <label className={Styles.label} htmlFor="password">Password</label>
                    {/* <input className={styles.input} type="password" name="password" id="password" onChange={handleChange} /> */}
                    <div className={Styles.inputWrapper}>
                        <input
                            className={Styles.input}
                            type="password"
                            name="password"
                            id="password"
                            onChange={handleChange}
                            required
                        />
                        <i className="togglePwd">
                            <span className="icon vis_icon material-symbols-outlined">visibility</span>
                        </i>
                    </div>

                    <label className={Styles.label} htmlFor="birthdate">Date of Birth</label>
                    <input className={Styles.input} type="date" name="birthdate" id="dob" onChange={handleChange} />

                    <label className={Styles.label} htmlFor="gender">Gender</label>
                    <div>
                        <label htmlFor="male">Male</label>
                        <input type="radio" name="gender" id="male" value="male" checked={formData.gender === "male"} onChange={handleChange} />

                        <label htmlFor="female">Female</label>
                        <input type="radio" name="gender" id="female" value="female" checked={formData.gender === "female"} onChange={handleChange} />
                    </div>

                    <label htmlFor="image" style={{ cursor: "pointer" }} className={Styles.label}>
                        <img src="/Image.svg" alt="Upload" width="24" height="24" />&nbsp;&nbsp;
                        Profile Image
                    </label>
                    <input className={Styles.input} type="file" name="avatar" id="profileImg" style={{ display: "none" }} accept="image/*" onChange={handleChange} />
                    {previewUrl && (
                        <img src={previewUrl} alt="Preview" />
                    )}

                    <label className={Styles.label} htmlFor="about">About Me</label>
                    <input className={Styles.input} type="text" name="about" id="about" onChange={handleChange} />

                    <button className={Styles.submitButton} type="submit">Register</button>
                </form>
            </div>
        </div>
    );
}
