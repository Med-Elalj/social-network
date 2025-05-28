# API v1 Documentation

**Base Path:** `/api/v1/`

---

## 📌 Endpoints Overview

### ✅ General

- **GET** `/api/v1/ws`  
  _WebSocket connection endpoint_  
  **TODO**:

  - Being logged in WILL be required,

  - other requirments to be decided at later date.

---

### 🔐 Auth (`/api/v1/auth`)

- **GET** `/auth`  
  _Check login status_  
  **TODO**: simple get request will return will return `{"status": true/false}`

- **POST** `/auth/register`  
  _User registration_  
  **Requirements**:

  - _not logged in_

  - {"username": "User_Name","email": "example@web.site","birthdate": "2001-11-09","fname": "Firstname", "lname": "LastName", "password": "password", "gender": "DFK"}.

    NOTE: check "sn/structs/input.go" for detail on what is valid

- **POST** `/auth/login`  
  _User login_  
  **TODO**:

  - _not logged in_
  - `{"login":"username/email","password":"eazerazer"}`

- **POST** `/auth/logout`  
  _User logout_  
  **TODO**: -Just post to it

---

### 📥 GET Endpoints (`/api/v1/get/`)

- **GET** `/get/profile`  
  _Retrieve user profile_  
  **TODO**: Add requirements

- **GET** `/get/posts`  
  _Fetch posts_  
  **TODO**: Add requirements

- **GET** `/get/categories`  
  _List all categories_  
  **TODO**: Add requirements

- **GET** `/get/conversations`  
  _List user conversations_  
  **TODO**: Add requirements

- **GET** `/get/messages`  
  _Get messages from a conversation_  
  **TODO**: Add requirements

---

### 📤 SET Endpoints (`/api/v1/set/`)

- **POST** `/set/profile`  
  _Update user profile_  
  **TODO**: Add requirements

- **POST** `/set/posts`  
  _Create or update posts_  
  **TODO**: Add requirements

- **POST** `/set/comments`  
  _Add a comment_  
  **TODO**: Add requirements

- **POST** `/set/follow`  
  _Follow or unfollow a user_  
  **TODO**: Add requirements

---
