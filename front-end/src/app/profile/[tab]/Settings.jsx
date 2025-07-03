'use client'; // REQUIRED if you're using App Router (in `/app` folder)

import { useState } from 'react';
import Style from '../profile.module.css'; // Adjust the path as necessary
import { showNotification } from '../../utils.jsx'; 

export default function Settings() {
  const [activeForm, setActiveForm] = useState(null); // 'nickname', 'password', or 'delete'
  const [formData, setFormData] = useState({
    nickname: '',
    currentPassword: '',
    newPassword: '',
    confirmPassword: '',
  });

  const handleChange = (e) => {
    setFormData({
      ...formData,
      [e.target.name]: e.target.value,
    });
  };

  const sendData = async (endpoint, data) => {
    try {
      const res = await fetch(endpoint, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(data),
      });
      const json = await res.json();
      console.log('Response:', json);
      
      // Use your showNotification function instead of alert
      if (res.ok) {
        showNotification(json.message || 'Request completed successfully.', 'success');
      } else {
        showNotification(json.message || 'Request failed.', 'error');
      }
    } catch (err) {
      console.error('Error:', err);
      showNotification('An error occurred.', 'error');
    }
  };

  const handleNicknameUpdate = () => {
    setActiveForm(activeForm === 'nickname' ? null : 'nickname');
  };

  const handleNicknameSubmit = () => {
    if (!formData.nickname.trim()) {
      showNotification('Please enter a nickname.', 'error');
      return;
    }

    sendData('/api/v1/updateProfile', { nickname: formData.nickname });
    setActiveForm(null);
    setFormData({ ...formData, nickname: '' }); // Clear nickname field
  };

  const handleChangePassword = () => {
    setActiveForm(activeForm === 'password' ? null : 'password');
  };

  const handlePasswordSubmit = () => {
    if (!formData.currentPassword || !formData.newPassword || !formData.confirmPassword) {
      showNotification('Please fill in all password fields.', 'error');
      return;
    }

    if (formData.newPassword !== formData.confirmPassword) {
      showNotification('New password and confirmation do not match.', 'error');
      return;
    }

    if (formData.newPassword.length < 8) {
      showNotification('New password must be at least 8 characters long.', 'error');
      return;
    }

    sendData('/api/v1/updatePassword', {
      currentPassword: formData.currentPassword,
      newPassword: formData.newPassword,
    });

    // Clear password fields after submission
    setFormData({
      ...formData,
      currentPassword: '',
      newPassword: '',
      confirmPassword: '',
    });
    setActiveForm(null);
  };

  const handleDeleteProfile = () => {
    setActiveForm(activeForm === 'delete' ? null : 'delete');
  };

  const handleDeleteConfirm = () => {
    sendData('/api/v1/delete', { confirmDelete: true });
  };

  const handleCancel = () => {
    setActiveForm(null);
    // Clear all form data when canceling
    setFormData({
      nickname: '',
      currentPassword: '',
      newPassword: '',
      confirmPassword: '',
    });
  };

  return (
    <div className={Style.post} style={{ width: '80%', margin: '20px auto' }}>
        <h2 className={Style.title}>Profile Settings</h2>
      <section
        className={Style.footer}
        style={{ display: 'flex', gap: '10px', flexWrap: 'wrap' }}
      >
        <button 
          type="button" 
          className={Style.button} 
          onClick={handleNicknameUpdate}
          style={activeForm === 'nickname' ? { backgroundColor: '#007bff', color: 'white' } : {}}
        >
          Update Nickname
        </button>
        <button 
          type="button" 
          className={Style.button} 
          onClick={handleChangePassword}
          style={activeForm === 'password' ? { backgroundColor: '#007bff', color: 'white' } : {}}
        >
          Change Password
        </button>
        <button
          type="button"
          className={Style.button}
          style={{ 
            backgroundColor: activeForm === 'delete' ? '#dc3545' : 'red', 
            color: 'white', 
            border: '1px solid black' 
          }}
          onClick={handleDeleteProfile}
        >
          Delete Profile
        </button>
      </section>

      {activeForm === 'nickname' && (
        <div style={{ marginTop: '20px' }}>
          <h4>Update Nickname</h4>

          <div className={Style.inputWrapper}>
            <input
              className={Style.input}
              type="text"
              name="nickname"
              placeholder="Enter new nickname"
              value={formData.nickname}
              onChange={handleChange}
              required
            />
          </div>

          <div style={{ marginTop: '10px', display: 'flex', gap: '10px' }}>
            <button
              type="button"
              className={Style.button}
              onClick={handleNicknameSubmit}
            >
              Update Nickname
            </button>
            <button
              type="button"
              className={Style.button}
              style={{ backgroundColor: 'gray' }}
              onClick={handleCancel}
            >
              Cancel
            </button>
          </div>
        </div>
      )}

      {activeForm === 'password' && (
        <div style={{ marginTop: '20px' }}>
          <h4>Change Password</h4>

          <div className={Style.inputWrapper}>
            <input
              className={Style.input}
              type="password"
              name="currentPassword"
              placeholder="Current Password"
              value={formData.currentPassword}
              onChange={handleChange}
              required
            />
          </div>

          <div className={Style.inputWrapper}>
            <input
              className={Style.input}
              type="password"
              name="newPassword"
              placeholder="New Password"
              value={formData.newPassword}
              onChange={handleChange}
              required
            />
          </div>

          <div className={Style.inputWrapper}>
            <input
              className={Style.input}
              type="password"
              name="confirmPassword"
              placeholder="Confirm New Password"
              value={formData.confirmPassword}
              onChange={handleChange}
              required
            />
          </div>

          <div style={{ marginTop: '10px', display: 'flex', gap: '10px' }}>
            <button
              type="button"
              className={Style.button}
              onClick={handlePasswordSubmit}
            >
              Update Password
            </button>
            <button
              type="button"
              className={Style.button}
              style={{ backgroundColor: 'gray' }}
              onClick={handleCancel}
            >
              Cancel
            </button>
          </div>
        </div>
      )}

      {activeForm === 'delete' && (
        <div style={{ marginTop: '20px' }}>
          <h4 style={{ color: 'red' }}>Delete Profile</h4>
          <p style={{ color: '#666', marginBottom: '15px' }}>
            ⚠️ This action cannot be undone. Your profile and all associated data will be permanently deleted.
          </p>
          
          <div style={{ marginTop: '10px', display: 'flex', gap: '10px' }}>
            <button
              type="button"
              className={Style.button}
              style={{ backgroundColor: '#dc3545', color: 'white', border: '1px solid #dc3545' }}
              onClick={handleDeleteConfirm}
            >
              Yes, Delete My Profile
            </button>
            <button
              type="button"
              className={Style.button}
              style={{ backgroundColor: 'gray' }}
              onClick={handleCancel}
            >
              Cancel
            </button>
          </div>
        </div>
      )}
    </div>
  );
}