CREATE TABLE
    IF NOT EXISTS followers (
        follower TEXT NOT NULL,
        followed TEXT NOT NULL,
        FOREIGN KEY (follower) REFERENCES users (username),
        FOREIGN KEY (followed) REFERENCES users (username)
    );

CREATE TABLE
    IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        uuid TEXT UNIQUE NOT NULL,
        username TEXT UNIQUE NOT NULL,
        email TEXT UNIQUE NOT NULL,
        birthdate DATE DEFAULT NULL,
        age TEXT DEFAULT NULL,
        gender TEXT NOT NULL,
        fname TEXT NOT NULL,
        lname TEXT NOT NULL,
        password TEXT NOT NULL,
        avatar TEXT DEFAULT NULL,
        aboutme TEXT DEFAULT NULL,
        status TEXT NOT NULL CHECK (status IN ('public', 'private')) DEFAULT 'public'
    );

CREATE TABLE
    IF NOT EXISTS notifications (
        notif TEXT NOT NULL,
        user TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (user) REFERENCES users (username)
    );

CREATE TABLE
    IF NOT EXISTS groupTable (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        description TEXT NOT NULL,
        creator_id INT NOT NULL,
        FOREIGN KEY (creator_id) REFERENCES users (id)
    );

CREATE TABLE
    IF NOT EXISTS groupmembers (
        group_id INT not null,
        user_id INT NOT NULL,
        goined_time DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (group_id) REFERENCES groupTable (id),
        FOREIGN KEY (user_id) REFERENCES users (id)
    );

CREATE TABLE
    IF NOT EXISTS groupcomment (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        content TEXT NOT NULL,
        post_id INTEGER NOT NULL,
        user_id INTEGER NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (post_id) REFERENCES grouposts (id),
        FOREIGN KEY (user_id) REFERENCES users (id)
    );

CREATE TABLE
    IF NOT EXISTS grouposts (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        group_id INT not null,
        user_id INT NOT NULL,
        title TEXT NOT NULL,
        post_content TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        attachement TEXT DEFAULT NULL,
        FOREIGN KEY (group_id) REFERENCES groupTable (id),
        FOREIGN KEY (user_id) REFERENCES users (id)
    );

CREATE TABLE
    IF NOT EXISTS joinrequest (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        group_id INT not null,
        sender_id INT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        stqtus TEXT DEFAULT NULL,
        FOREIGN KEY (group_id) REFERENCES groupTable (id),
        FOREIGN KEY (sender_id) REFERENCES users (id)
    );

CREATE TABLE
    IF NOT EXISTS posts (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        id_users INTEGER NOT NULL,
        title TEXT NOT NULL,
        content TEXT NOT NULL,
        PostPlace TEXT DEFAULT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        attachement TEXT DEFAULT NULL,
        status TEXT NOT NULL,
        FOREIGN KEY (id_users) REFERENCES users (id)
    );

CREATE TABLE
    IF NOT EXISTS category (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );

CREATE TABLE
    IF NOT EXISTS post_category (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        catego_id INTEGER NOT NULL,
        post_id INTEGER NOT NULL,
        FOREIGN KEY (catego_id) REFERENCES category (id),
        FOREIGN KEY (post_id) REFERENCES posts (id)
    );

CREATE TABLE
    IF NOT EXISTS comments (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        content TEXT NOT NULL,
        post_id INTEGER NOT NULL,
        user_id INTEGER NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (post_id) REFERENCES posts (id),
        FOREIGN KEY (user_id) REFERENCES users (id)
    );

CREATE TABLE
    IF NOT EXISTS groupMessage (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        sender_id INTEGER NOT NULL,
        receiver_id INTEGER NOT NULL,
        content TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (sender_id) REFERENCES users (id),
        FOREIGN KEY (receiver_id) REFERENCES groupTable (id)
    );

CREATE TABLE
    IF NOT EXISTS messages (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        sender_id INTEGER NOT NULL,
        receiver_id INTEGER NOT NULL,
        content TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        read BOOLEAN DEFAULT 0, -- New column (0=unread, 1=read)
        FOREIGN KEY (sender_id) REFERENCES users (id),
        FOREIGN KEY (receiver_id) REFERENCES users (id)
    );