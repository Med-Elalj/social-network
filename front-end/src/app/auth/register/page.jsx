"use client"

import { useState } from 'react';

export default function Login() {
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

        // TODO: send `data` to backend via fetch
        console.log('Submitting form...', formData);
    };

    return (
        <form onSubmit={handleSubmit}>
            <label htmlFor="email">Email</label>
            <input type="email" name="email" id="email" onChange={handleChange} />

            <label htmlFor="password">Password</label>
            <input type="password" name="password" id="password" onChange={handleChange} />


            <label htmlFor="firstName">First Name</label>
            <input type="text" name="firstName" id="firstName" onChange={handleChange} />


            <label htmlFor="lastName">Last Name</label>
            <input type="text" name="lastName" id="lastName" onChange={handleChange} />


            <label htmlFor="nickName">Nickname</label>
            <input type="text" name="nickName" id="nickName" onChange={handleChange} />


            <label htmlFor="dob">Date of Birth</label>
            <input type="date" name="dob" id="dob" onChange={handleChange} />


            <label htmlFor="profileImg">Profile Image</label>
            <input type="file" name="profileImg" id="profileImg" accept="image/*" onChange={handleChange} />

            <label htmlFor="about">About Me</label>
            <input type="text" name="about" id="about" onChange={handleChange} />


            <button type="submit">Register</button>
        </form>
    );
}
