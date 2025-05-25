PRAGMA foreign_keys=ON;
BEGIN TRANSACTION;
CREATE TABLE followers (
	follower TEXT NOT NULL,
	followed TEXT NOT NULL,
	FOREIGN KEY(follower) REFERENCES users(username),
	FOREIGN KEY(followed) REFERENCES users(username)
	);
CREATE TABLE users (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
uuid TEXT UNIQUE NOT NULL,
username TEXT UNIQUE NOT NULL,
email TEXT UNIQUE NOT NULL,
birthdate DATE  DEFAULT NULL,
age TEXT DEFAULT NULL,
gender TEXT NOT NULL,
fname TEXT NOT NULL,
lname TEXT NOT NULL,
password TEXT NOT NULL,
avatar TEXT DEFAULT NULL,
aboutme TEXT DEFAULT NULL,
status TEXT NOT NULL CHECK (status IN ('public', 'private')) DEFAULT 'public'
);
INSERT INTO users VALUES(1,'6352337196a2449cb772b524818bea36','bbb','bbb@gmail.com','1999-05-10','19','male','bbb','bbbbb','bbb','Screenshot 2024-11-25 185542.png','fsjklsdklfjsdjgksdjklghkldsdjfkdjhfsdhgjlsdjglksdjgkjksdjgkldjgkjfklgjlfkljgkfjklgjklgjfkjgfjgklfjkgfdkl','public');
INSERT INTO users VALUES(2,'ee9774d09b884691b54f1df39fa4e03f','aaa','aaa@gmaol.clp','1992-12-31','20','male','aaa','aaa','aaa','','fdgdfgfdg','public');
CREATE TABLE notifications (
	notif TEXT NOT NULL,
	user TEXT NOT NULL,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY(user) REFERENCES users(username)
	);
INSERT INTO notifications VALUES('aaa start to follow you','bbb','2025-05-07 16:33:07');
CREATE TABLE groupTable (
	 	id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		description TEXT NOT NULL,
		userid INT NOT NULL ,
		FOREIGN KEY(userid) REFERENCES users(id)
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
	FOREIGN KEY(id_users) REFERENCES users(id)
);
INSERT INTO posts VALUES(1,1,'bbbbbbb','dkfhklsdjklfajlkadsjfkjsdjfsdkjf',NULL,'2025-05-06 15:58:10','/uploads\Screenshot 2024-11-25 185645.png','public');
INSERT INTO posts VALUES(2,1,'fsdfklsdjjkfsd','dskfkjjksdlf',NULL,'2025-05-06 16:04:45','/uploads\Screenshot 2024-11-25 185542.png','public');
INSERT INTO posts VALUES(3,1,'fsdfsdf','fsdfsdf',NULL,'2025-05-06 16:07:49','/uploads\Screenshot 2024-11-16 171359.png','public');
INSERT INTO posts VALUES(4,1,'dfsfsdf','fsdfsd',NULL,'2025-05-06 16:11:03','/uploads\Screenshot 2024-11-16 171359.png','public');
INSERT INTO posts VALUES(5,1,'asdfsdafsdaf','dfsafasfsda',NULL,'2025-05-06 16:19:00','/uploads\Screenshot 2024-11-16 171359.png','public');
INSERT INTO posts VALUES(6,1,'dsjfsdf','fsdfsdf',NULL,'2025-05-07 10:29:20','/uploads\Screenshot 2024-11-16 171359.png','public');
INSERT INTO posts VALUES(7,1,'dsfsdffsdfsdfsd','fsdfsdf',NULL,'2025-05-07 10:34:00','/uploads\Screenshot 2024-11-16 171359.png','public');
INSERT INTO posts VALUES(8,1,'fsdfsdfs','fsdfsdfsd',NULL,'2025-05-07 10:36:32','/uploads/Screenshot 2024-11-16 171359.png','public');
INSERT INTO posts VALUES(9,1,'kljkljkljkl','lkjkljkjklj',NULL,'2025-05-07 10:41:05','','public');
INSERT INTO posts VALUES(10,1,'hghghjgj','jhghghjghjgh',NULL,'2025-05-07 10:45:07','','public');
INSERT INTO posts VALUES(11,1,'fdsjfjsdfjkhsdjfjk','jkdfhjkfhsdjfhsdjk',NULL,'2025-05-07 10:47:15','','public');
INSERT INTO posts VALUES(12,1,'dfhjkfhdsjk','jhfdjshjk',NULL,'2025-05-07 10:50:09','\uploads\Screenshot 2024-11-16 171359.png','public');
INSERT INTO posts VALUES(13,1,'dsfjsdhf','hfdjhsdjkfh',NULL,'2025-05-07 10:56:23','/uploads/Screenshot 2024-11-16 171359.png','public');
INSERT INTO posts VALUES(14,1,'filsdjflj','lfsksjfkslj',NULL,'2025-05-07 10:58:30','/uploads/Screenshot 2024-11-16 171359.png','public');
INSERT INTO posts VALUES(15,1,'fsdlkjf','klsdklfjskl',NULL,'2025-05-07 11:09:33','/uploads/Screenshot 2024-11-16 171359.png','public');
INSERT INTO posts VALUES(16,1,'aaa','fdfdfd',NULL,'2025-05-07 11:13:43','/uploads/1233.png','public');
INSERT INTO posts VALUES(17,1,'hjkjhkj','khjkjhkhj',NULL,'2025-05-07 11:37:39','/uploads/Untitled.png','public');
INSERT INTO posts VALUES(18,1,'dsjkfdh','kjkdjlkgjd',NULL,'2025-05-07 11:39:23','/uploads/Untitled.jpeg','public');
CREATE TABLE category(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	 );
INSERT INTO category VALUES(1,'Technology','2025-05-06 15:58:10');
INSERT INTO category VALUES(2,'Travel','2025-05-06 15:58:10');
INSERT INTO category VALUES(3,'Gaming','2025-05-06 15:58:10');
INSERT INTO category VALUES(4,'Health','2025-05-06 16:07:49');
INSERT INTO category VALUES(5,'Science','2025-05-07 10:47:15');
CREATE TABLE post_category(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	catego_id INTEGER NOT NULL,
	post_id INTEGER NOT NULL,
  	FOREIGN KEY(catego_id) REFERENCES category(id),
   	FOREIGN KEY(post_id) REFERENCES posts(id)
);
INSERT INTO post_category VALUES(1,1,1);
INSERT INTO post_category VALUES(2,2,1);
INSERT INTO post_category VALUES(3,3,1);
INSERT INTO post_category VALUES(4,1,2);
INSERT INTO post_category VALUES(5,4,3);
INSERT INTO post_category VALUES(6,2,4);
INSERT INTO post_category VALUES(7,1,5);
INSERT INTO post_category VALUES(8,1,6);
INSERT INTO post_category VALUES(9,2,7);
INSERT INTO post_category VALUES(10,2,8);
INSERT INTO post_category VALUES(11,2,9);
INSERT INTO post_category VALUES(12,1,10);
INSERT INTO post_category VALUES(13,5,11);
INSERT INTO post_category VALUES(14,5,12);
INSERT INTO post_category VALUES(15,5,13);
INSERT INTO post_category VALUES(16,3,14);
INSERT INTO post_category VALUES(17,3,15);
INSERT INTO post_category VALUES(18,2,16);
INSERT INTO post_category VALUES(19,4,17);
INSERT INTO post_category VALUES(20,4,18);
CREATE TABLE comments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    content TEXT NOT NULL,
	post_id INTEGER NOT NULL,
	user_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(post_id) REFERENCES posts(id),
	FOREIGN KEY(user_id) REFERENCES users(id)
);
CREATE TABLE groupMessage (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		sender_id INTEGER NOT NULL,
		receiver_id INTEGER NOT NULL,
		content TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(sender_id) REFERENCES users(id),
		FOREIGN KEY(receiver_id) REFERENCES groupTable(id)
	);
CREATE TABLE messages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		sender_id INTEGER NOT NULL,
		receiver_id INTEGER NOT NULL,
		content TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		read BOOLEAN DEFAULT 0,  -- New column (0=unread, 1=read)
		FOREIGN KEY(sender_id) REFERENCES users(id),
		FOREIGN KEY(receiver_id) REFERENCES users(id)
	);
INSERT INTO messages VALUES(1,1,2,'hiu','2025-05-07 11:43:20',1);
INSERT INTO messages VALUES(2,2,1,'hi','2025-05-07 11:43:26',1);
CREATE TABLE groupmembers (
	 	group_id INT not null,
		user_id INT NOT NULL,
		goined_time DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(group_id) REFERENCES groupTable(id),
		FOREIGN KEY(user_id) REFERENCES users(id)
	);
CREATE TABLE grouposts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
	 	group_id INT not null,
		user_id INT NOT NULL,
		title TEXT NOT NULL,
		post_content TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		attachement TEXT DEFAULT NULL,
		FOREIGN KEY(group_id) REFERENCES groupTable(id),
		FOREIGN KEY(user_id) REFERENCES users(id)
	);
CREATE TABLE joinrequest (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
	 	group_id INT not null,
		sender_id INT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		stqtus TEXT DEFAULT NULL,
		FOREIGN KEY(group_id) REFERENCES groupTable(id),
		FOREIGN KEY(sender_id) REFERENCES users(id)
	);
DELETE FROM sqlite_sequence;
INSERT INTO sqlite_sequence VALUES('users',2);
INSERT INTO sqlite_sequence VALUES('posts',18);
INSERT INTO sqlite_sequence VALUES('category',5);
INSERT INTO sqlite_sequence VALUES('post_category',20);
INSERT INTO sqlite_sequence VALUES('messages',2);
COMMIT;
