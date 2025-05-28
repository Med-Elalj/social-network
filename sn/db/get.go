package db

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"social-network/server/logs"
	"social-network/sn/structs"
)

func GetPosts(soffset string) ([]structs.Post1, error) {
	offset, err := strconv.Atoi(soffset)
	if err != nil {
		logs.Errorf("Error converting pid to int: %q", err.Error())
		return nil, err
	}
	rows, err := DB.Query(`
	SELECT 
   		p.post_id AS pid, 
    	p.title, 
    	p.content, 
    	p.categories, 
    	p.created_at, 
    	u.username AS author
	FROM
    	posts p
	JOIN 
    	users u ON p.uid = u.id
	ORDER BY
		p.created_at DESC
	LIMIT 3 OFFSET ?
	`, offset)
	if err != nil {
		logs.Errorf("Error getting posts: %q", err.Error())
		return nil, err
	}
	defer rows.Close()
	var categories string
	var posts []structs.Post1
	for rows.Next() {
		var post structs.Post1
		err := rows.Scan(&post.Pid, &post.Title, &post.Content, &categories, &post.CreationTime, &post.Author)
		if err != nil {
			logs.Errorf("Error scanning posts: %q", err.Error())
			return nil, err
		}
		post.Categories = strings.Split(categories, ", ")
		posts = append(posts, post)
	}
	return posts, nil
}

func GetComments(pid string) ([]structs.Comment, error) {
	if pid == "" {
		return nil, nil
	}
	iPid, err := strconv.Atoi(pid)
	if err != nil {
		logs.Errorf("Error converting pid to int: %q", err.Error())
		return nil, err
	}
	rows, err := DB.Query(`
	SELECT 
    	u.username AS author,
    	c.content,
    	c.created_at
	FROM 
    	comments c
	JOIN 
    	users u ON c.uid = u.id
	WHERE
		c.post_id = ?
	ORDER BY
		c.created_at DESC
	`, iPid)
	if err != nil {
		logs.Errorf("Error getting comments: %q", err.Error())
		return nil, err
	}
	defer rows.Close()
	var comments []structs.Comment
	for rows.Next() {
		var comment structs.Comment
		err := rows.Scan(&comment.Author, &comment.Content, &comment.CreationTime)
		if err != nil {
			logs.Errorf("Error scanning comments: %q", err.Error())
			return nil, err
		}
		comment.Pid = iPid
		comments = append(comments, comment)
	}
	return comments, nil
}

func GetUserNames(uid int) ([]structs.User1, error) {
	rows, err := DB.Query(`
	SELECT
		u.username
	FROM
		users u
	LEFT JOIN
		dms m
	ON
		(u.id = m.sender_id OR u.id = m.recipient_id)
	AND
		(m.sender_id = ? OR m.recipient_id = ? )
	WHERE
		u.id != ?
	GROUP BY
		u.id, u.username
	ORDER BY
		CASE WHEN MAX(m.created_at) IS NOT NULL THEN 1
	    ELSE 2
	END,
		MAX(m.created_at) DESC,
	u.username ASC;`, uid, uid, uid)
	if err != nil {
		return nil, fmt.Errorf("could not execute query: %w", err)
	}
	defer rows.Close()

	var userNames []structs.User1

	for rows.Next() {
		var username structs.User1
		if err := rows.Scan(&username.Username); err != nil {
			return userNames, fmt.Errorf("could not scan row: %w", err)
		}
		// TODO
		// if _, e := helpers.Sockets[username.Username]; e {
		// 	username.Online = true
		// }
		// userNames = append(userNames, username)
	}

	if err := rows.Err(); err != nil {
		return userNames, fmt.Errorf("row iteration error: %w", err)
	}

	return userNames, nil
}

func GetdmHistory(uname1, uname2, date string) ([]structs.Message, error) {
	var d time.Time
	if date == "" {
		d = time.Now()
	} else {
		var err error
		d, err = time.Parse(time.RFC3339, date)
		if err != nil {
			return nil, fmt.Errorf("invalid date format")
		}
	}
	rows, err := DB.Query(`
        SELECT * 
        FROM (
            SELECT
                sender.username, d.message, d.created_at
            FROM
                dms d
            JOIN
                users sender ON d.sender_id = sender.id
            JOIN
                users recipient ON d.recipient_id = recipient.id
            WHERE
                d.created_at < ?
                AND (
                    (sender.username = ? AND recipient.username = ?)
                    OR
                    (sender.username = ? AND recipient.username = ?)
                )
            ORDER BY
                d.created_at DESC
            LIMIT 10
        ) AS sub 
        ORDER BY created_at ASC;
    `, d, uname1, uname2, uname2, uname1)
	if err != nil {
		logs.Errorf("Error getting messages: %q", err.Error())
		return nil, err
	}
	defer rows.Close()

	var messages []structs.Message
	for rows.Next() {
		var message structs.Message
		if err := rows.Scan(&message.Sender, &message.Content, &message.Time); err != nil {
			logs.Errorf("Error scanning message: %q", err.Error())
			return nil, err
		}
		messages = append(messages, message)
	}
	return filter(messages, d), nil
}
