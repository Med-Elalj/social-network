package main

import (
	"database/sql"
	"log"
)

var db *sql.DB


func init() {
	var err error
	db, err = sql.Open("sqlite3", "./test1.db?_foreign_keys=1")
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("Database connection failed:", err)
	}
	createFollowersTable := ` 
	CREATE TABLE IF NOT EXISTS followers (
	follower TEXT NOT NULL,
	followed TEXT NOT NULL,
	FOREIGN KEY(follower) REFERENCES users(username),
	FOREIGN KEY(followed) REFERENCES users(username)
	)
	`
	_, err = db.Exec(createFollowersTable)
	if err != nil {
		log.Fatal("Failed to create users table:", err)
	}
	createUsersTable := `
	CREATE TABLE IF NOT EXISTS users (
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
);`
_, err = db.Exec(createUsersTable)
if err != nil {
	log.Fatal("Failed to create users table:", err)
}

	createnotificationTable := `
	CREATE TABLE IF NOT EXISTS notifications (
	notif TEXT NOT NULL,
	user TEXT NOT NULL,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY(user) REFERENCES users(username)
	)`
	_,err = db.Exec(createnotificationTable)
	if err != nil {
		log.Fatal("Failed to create notif table")
	}
	
	createGroupTable := `
	CREATE TABLE IF NOT EXISTS groupTable (
	 	id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		description TEXT NOT NULL,
		creator_id INT NOT NULL ,
		FOREIGN KEY(creator_id) REFERENCES users(id)
	)	
	`
	_, err = db.Exec(createGroupTable)
	if err != nil {
		log.Fatal("Failed to create groupe table:", err)
	}
	createGroupmember := `
	CREATE TABLE IF NOT EXISTS groupmembers (
	 	group_id INT not null,
		user_id INT NOT NULL,
		goined_time DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(group_id) REFERENCES groupTable(id),
		FOREIGN KEY(user_id) REFERENCES users(id)
	)	
	`
	_, err = db.Exec(createGroupmember)
	if err != nil {
		log.Fatal("Failed to create groupemember table:", err)
	}

	createCommentgroups := `
	CREATE TABLE IF NOT EXISTS groupcomment (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    content TEXT NOT NULL,
	post_id INTEGER NOT NULL,
	user_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(post_id) REFERENCES grouposts(id),
	FOREIGN KEY(user_id) REFERENCES users(id)
);`

_, err = db.Exec(createCommentgroups)
if err != nil {
	log.Fatal("Failed to create groupemember table:", err)
}
	createGrouPost := `
	CREATE TABLE IF NOT EXISTS grouposts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
	 	group_id INT not null,
		user_id INT NOT NULL,
		title TEXT NOT NULL,
		post_content TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		attachement TEXT DEFAULT NULL,
		FOREIGN KEY(group_id) REFERENCES groupTable(id),
		FOREIGN KEY(user_id) REFERENCES users(id)
	)	
	`
	_, err = db.Exec(createGrouPost)
	if err != nil {
		log.Fatal("Failed to create groupepost table:", err)
	}
	
	createjoinrequest := `
	CREATE TABLE IF NOT EXISTS joinrequest (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
	 	group_id INT not null,
		sender_id INT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		stqtus TEXT DEFAULT NULL,
		FOREIGN KEY(group_id) REFERENCES groupTable(id),
		FOREIGN KEY(sender_id) REFERENCES users(id)
	)	
	`
	_, err = db.Exec(createjoinrequest)
	if err != nil {
		log.Fatal("Failed to create joinrequest table:", err)
	}
	createPostsTable := `
	CREATE TABLE IF NOT EXISTS posts (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	id_users INTEGER NOT NULL,
	title TEXT NOT NULL,
	content TEXT NOT NULL,
	PostPlace TEXT DEFAULT NULL,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	attachement TEXT DEFAULT NULL,
	status TEXT NOT NULL,
	FOREIGN KEY(id_users) REFERENCES users(id)
);`
_, err = db.Exec(createPostsTable)
if err != nil {
	log.Fatal("Failed to create posts table:", err)
}

	createCategoriesTable := ` 
		CREATE TABLE IF NOT EXISTS category(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	 );
	`

	_, err = db.Exec(createCategoriesTable)
	if err != nil {
		log.Fatal("Failed to create posts table:", err)
	}

	createPostCategories := ` 
	CREATE TABLE IF NOT EXISTS post_category(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	catego_id INTEGER NOT NULL,
	post_id INTEGER NOT NULL,
  	FOREIGN KEY(catego_id) REFERENCES category(id),
   	FOREIGN KEY(post_id) REFERENCES posts(id)
);
`
	_, err = db.Exec(createPostCategories)
	if err != nil {
		log.Fatal("Failed to create posts table:", err)
	}

	createCommentsTable := `
	CREATE TABLE IF NOT EXISTS comments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    content TEXT NOT NULL,
	post_id INTEGER NOT NULL,
	user_id INTEGER NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(post_id) REFERENCES posts(id),
	FOREIGN KEY(user_id) REFERENCES users(id)
);`
	_, err = db.Exec(createCommentsTable)
	if err != nil {
		log.Fatal("Failed to create comments table:", err)
	}

	groupMessageTable := `
	CREATE TABLE IF NOT EXISTS groupMessage (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		sender_id INTEGER NOT NULL,
		receiver_id INTEGER NOT NULL,
		content TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(sender_id) REFERENCES users(id),
		FOREIGN KEY(receiver_id) REFERENCES groupTable(id)
	);`

	_, err = db.Exec(groupMessageTable)
	if err != nil {
		log.Fatal("Failed to create groupemsg table:", err)
	}
	createMessagesTable := `
	CREATE TABLE IF NOT EXISTS messages (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		sender_id INTEGER NOT NULL,
		receiver_id INTEGER NOT NULL,
		content TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		read BOOLEAN DEFAULT 0,  -- New column (0=unread, 1=read)
		FOREIGN KEY(sender_id) REFERENCES users(id),
		FOREIGN KEY(receiver_id) REFERENCES users(id)
	);`

	_, err = db.Exec(createMessagesTable)
	if err != nil {
		log.Fatal("Failed to create messages table:", err)
	}
}

// func init() {
// 	var err error
// 	db, err = sql.Open("sqlite3", "./test8.db?_foreign_keys=1")
// 	if err != nil {
// 		log.Fatal("Failed to open database:", err)
// 	}

// 	if err := db.Ping(); err != nil {
// 		log.Fatal("Database connection failed:", err)
// 	}

// 	createUsersTable := `
// 		CREATE TABLE IF NOT EXISTS users (
// 		id INTEGER PRIMARY KEY AUTOINCREMENT,
// 		uuid TEXT UNIQUE NOT NULL,
// 		username TEXT UNIQUE NOT NULL,
// 		email TEXT UNIQUE NOT NULL,
// 		age TEXT,
// 		gender TEXT NOT NULL,
// 		fname TEXT NOT NULL,
// 		lname TEXT NOT NULL,
// 		password TEXT NOT NULL
// 	);`
// 	_, err = db.Exec(createUsersTable)
// 	if err != nil {
// 		log.Fatal("Failed to create users table:", err)
// 	}

// 	createPostsTable := `
// 		CREATE TABLE IF NOT EXISTS posts (
// 		id INTEGER PRIMARY KEY AUTOINCREMENT,
// 		id_users INTEGER NOT NULL,
// 		title TEXT NOT NULL,
// 		content TEXT NOT NULL,
// 		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
// 		FOREIGN KEY(id_users) REFERENCES users(id)

// 	);`
// 	_, err = db.Exec(createPostsTable)
// 	if err != nil {
// 		log.Fatal("Failed to create posts table:", err)
// 	}

// 	createCategoriesTable := ` 
// 		CREATE TABLE IF NOT EXISTS category(
// 		id INTEGER PRIMARY KEY AUTOINCREMENT,
// 		name TEXT NOT NULL,
// 		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
// 	 );
// 	`

// 	_, err = db.Exec(createCategoriesTable)
// 	if err != nil {
// 		log.Fatal("Failed to create posts table:", err)
// 	}

// 	createPostCategories := ` 
// 	CREATE TABLE IF NOT EXISTS post_category(
// 	id INTEGER PRIMARY KEY AUTOINCREMENT,
// 	catego_id INTEGER NOT NULL,
// 	post_id INTEGER NOT NULL,
//   	FOREIGN KEY(catego_id) REFERENCES category(id),
//    	FOREIGN KEY(post_id) REFERENCES posts(id)
// );
// `
// 	_, err = db.Exec(createPostCategories)
// 	if err != nil {
// 		log.Fatal("Failed to create posts table:", err)
// 	}

// 	createCommentsTable := `
// 	CREATE TABLE IF NOT EXISTS comments (
//     id INTEGER PRIMARY KEY AUTOINCREMENT,
//     content TEXT NOT NULL,
// 	post_id INTEGER NOT NULL,
// 	user_id INTEGER NOT NULL,
//     created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
//     FOREIGN KEY(post_id) REFERENCES posts(id),
// 	FOREIGN KEY(user_id) REFERENCES users(id)
// );`
// 	_, err = db.Exec(createCommentsTable)
// 	if err != nil {
// 		log.Fatal("Failed to create comments table:", err)
// 	}

// 	createMessagesTable := `
// 	CREATE TABLE IF NOT EXISTS messages (
// 		id INTEGER PRIMARY KEY AUTOINCREMENT,
// 		sender_id INTEGER NOT NULL,
// 		receiver_id INTEGER NOT NULL,
// 		content TEXT NOT NULL,
// 		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
// 		read BOOLEAN DEFAULT 0,  -- New column (0=unread, 1=read)
// 		FOREIGN KEY(sender_id) REFERENCES users(id),
// 		FOREIGN KEY(receiver_id) REFERENCES users(id)
// 	);`

// 	_, err = db.Exec(createMessagesTable)
// 	if err != nil {
// 		log.Fatal("Failed to create messages table:", err)
// 	}
// }

