package db

import (
	"fmt"
	"strconv"
	"time"

	"social-network/server/logs"
	"social-network/sn/structs"
)

func GetPosts(soffset string) ([]structs.PostGet, error) {
	// offset, err := strconv.Atoi(soffset) // TODO ADD
	// if err != nil {
	// 	logs.Errorf("Error converting pid to int: %q", err.Error())
	// 	return nil, err
	// }
	rows, err := DB.Query(`WITH followed_profiles AS (
	SELECT following_id FROM follow WHERE follower_id = ? -- Replace with the actual user ID
),
pivate_post_see AS (
	SELECT post_id FROM pivate_post_visibility WHERE user_id = ? -- Replace with the actual user ID
)
SELECT 
   		p.id, p.user_id, u.display_name,p.title,p.content, p.created_at,p.privacy,
		CASE 
			WHEN g.id IS NULL THEN 0 ELSE g.id
		END AS 'group_id',
        CASE
            WHEN g.id IS NULL THEN '' ELSE g.display_name
        END AS 'name'
	FROM
    	posts p
	JOIN 
    	profile u ON p.user_id = u.id
    LEFT JOIN 
    	profile g ON p.group_id = g.id
	WHERE
			p.privacy = 0
		OR 
			(
				p.privacy = 1
				AND 
				(p.user_id IN followed_profiles OR p.group_id IN followed_profiles)
			)
		OR 
			(
					p.privacy = 2
				AND 
					p.id IN pivate_post_see
			)

		
	ORDER BY
		p.created_at DESC
	LIMIT 10; -- TODO Add pagination
`)
	if err != nil {
		logs.Errorf("Error getting posts: %q", err.Error())
		return nil, err
	}
	defer rows.Close()
	var posts []structs.PostGet
	for rows.Next() {
		var post structs.PostGet
		err := rows.Scan(&post.Pid, &post.AuthorId, &post.Author, &post.Title, &post.Content, &post.CreationTime, &post.Privacy, &post.GroupId, &post.GroupName)
		if err != nil {
			logs.Errorf("Error scanning posts: %q", err.Error())
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func GetComments(pid string) ([]structs.CommentGet, error) {
	if pid == "" {
		return nil, nil
	}
	Pid, err := strconv.Atoi(pid)
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
	`, Pid)
	if err != nil {
		logs.Errorf("Error getting comments: %q", err.Error())
		return nil, err
	}
	defer rows.Close()
	var comments []structs.CommentGet
	for rows.Next() {
		var comment structs.CommentGet
		err := rows.Scan(&comment.Author, &comment.Content, &comment.CreationTime)
		if err != nil {
			logs.Errorf("Error scanning comments: %q", err.Error())
			return nil, err
		}
		comment.Pid = structs.ID(Pid)
		comments = append(comments, comment)
	}
	return comments, nil
}

func GetUserNames(uid int) ([]structs.UsersGet, error) {
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

	var userNames []structs.UsersGet

	for rows.Next() {
		var username structs.UsersGet
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
