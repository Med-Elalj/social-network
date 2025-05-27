PRAGMA foreign_keys=OFF;
BEGIN TRANSACTION;
CREATE TABLE followers (
        follower TEXT NOT NULL,
        followed TEXT NOT NULL,
        FOREIGN KEY (follower) REFERENCES users (username),
        FOREIGN KEY (followed) REFERENCES users (username)
    );
CREATE TABLE users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT UNIQUE NOT NULL,
        email TEXT UNIQUE NOT NULL,
        birthdate DATE DEFAULT (DATE('now')),
        gender INT NOT NULL,
        fname TEXT NOT NULL,
        lname TEXT NOT NULL,
        password TEXT NOT NULL,
        avatar TEXT DEFAULT NULL,
        aboutme TEXT DEFAULT NULL,
        status TEXT NOT NULL CHECK (status IN ('public', 'private')) DEFAULT 'public'
    );
CREATE TABLE notifications (
        notif TEXT NOT NULL,
        user TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (user) REFERENCES users (username)
    );
CREATE TABLE groupTable (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        description TEXT NOT NULL,
        creator_id INT NOT NULL,
        FOREIGN KEY (creator_id) REFERENCES users (id)
    );
CREATE TABLE groupmembers (
        group_id INT not null,
        user_id INT NOT NULL,
        goined_time DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (group_id) REFERENCES groupTable (id),
        FOREIGN KEY (user_id) REFERENCES users (id)
    );
CREATE TABLE posts (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        id_users INTEGER NOT NULL,
        title TEXT NOT NULL,
        content TEXT NOT NULL,
        PostPlace TEXT DEFAULT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        attachement TEXT DEFAULT NULL,
        status TEXT NOT NULL,
        group_id INT DEFAULT null,
        FOREIGN KEY (id_users) REFERENCES users (id),
        FOREIGN KEY (group_id) REFERENCES groupTable (id)
    );
CREATE TABLE category (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );
CREATE TABLE post_category (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        catego_id INTEGER NOT NULL,
        post_id INTEGER NOT NULL,
        FOREIGN KEY (catego_id) REFERENCES category (id),
        FOREIGN KEY (post_id) REFERENCES posts (id)
    );
CREATE TABLE comments (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        content TEXT NOT NULL,
        post_id INTEGER NOT NULL,
        user_id INTEGER NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (post_id) REFERENCES posts (id),
        FOREIGN KEY (user_id) REFERENCES users (id)
    );
CREATE TABLE groupMessage (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        sender_id INTEGER NOT NULL,
        receiver_id INTEGER NOT NULL,
        content TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (sender_id) REFERENCES users (id),
        FOREIGN KEY (receiver_id) REFERENCES groupTable (id)
    );
CREATE TABLE messages (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        sender_id INTEGER NOT NULL,
        receiver_id INTEGER NOT NULL,
        content TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        read BOOLEAN DEFAULT 0, -- New column (0=unread, 1=read)
        FOREIGN KEY (sender_id) REFERENCES users (id),
        FOREIGN KEY (receiver_id) REFERENCES users (id)
    );
CREATE TABLE postreaction (
        post_id INTEGER NOT NULL,
        user_id INTEGER NOT NULL,
        action BOOLEAN NOT NULL,
        FOREIGN KEY (post_id) REFERENCES posts (id),
        FOREIGN KEY (user_id) REFERENCES users (id),
        PRIMARY KEY (user_id, post_id)
    );
CREATE TABLE commentreaction (
        comment_id INTEGER NOT NULL,
        user_id INTEGER NOT NULL,
        action BOOLEAN NOT null,
        FOREIGN KEY (user_id) REFERENCES users (id),
        FOREIGN KEY (comment_id) REFERENCES comments (id),
        PRIMARY KEY (user_id, comment_id)
    );

INSERT INTO users(username, email, birthdate, password, gender, fname, lname) VALUES
    ("Uname_1","email_1@web.site","2001-11-09","pwd hash","0","FirstName","LastName"),
    ("Uname_2","email_2@web.site","2001-11-09","pwd hash","0","FirstName","LastName"),
    ("Uname_3","email_3@web.site","2001-11-09","pwd hash","0","FirstName","LastName"),
    ("Uname_4","email_4@web.site","2001-11-09","pwd hash","0","FirstName","LastName"),
    ("Uname_5","email_5@web.site","2001-11-09","pwd hash","0","FirstName","LastName"),
    ("Uname_6","email_6@web.site","2001-11-09","pwd hash","0","FirstName","LastName"),
    ("Uname_7","email_7@web.site","2001-11-09","pwd hash","0","FirstName","LastName");
DELETE FROM sqlite_sequence;
COMMIT;
