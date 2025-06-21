package modules

import (
	"fmt"

	"social-network/app/structs"
	"social-network/server/logs"
)

func GetPosts(start, uid, groupId int) ([]structs.Post, error) {
	query := `
		WITH
			user_groups AS (
				SELECT group_id
				FROM groupmember
				WHERE person_id = ? AND active = 1
			),
			followed_profiles AS (
				SELECT following_id
				FROM follow
				WHERE follower_id = ? AND status = 1
			)
		SELECT 
			p.id,
			p.content,
			p.image_path,
			p.created_at,
			pr.display_name,
			p.privacy,
			(SELECT COUNT(*) FROM comments c WHERE c.post_id = p.id) AS comment_count,
			(SELECT COUNT(*) FROM likes l WHERE l.post_id = p.id) AS like_count
		FROM
			posts p
		JOIN profile pr ON pr.id = p.user_id
		LEFT JOIN user_groups ug ON p.group_id = ug.group_id
		WHERE
			(? = 0 OR p.id < ?) AND
			(? = 0 OR p.group_id = ?) AND
			p.privacy != 'private'
 			AND (
				p.privacy = 'public'
				OR p.user_id = ?
				OR (p.privacy = 'friends' AND p.user_id IN (SELECT following_id FROM followed_profiles))
				OR (p.group_id IS NOT NULL AND ug.group_id IS NOT NULL)
			)
		ORDER BY p.id DESC
		LIMIT 10;
	`

	rows, err := DB.Query(query,
		uid,          // for user_groups
		uid,          // for followed_profiles
		start, start, // pagination
		groupId, groupId, // group filter
		uid, // post owner visibility
	)
	if err != nil {
		logs.Errorf("GetPosts query error: %q", err.Error())
		return nil, err
	}
	defer rows.Close()

	var posts []structs.Post
	for rows.Next() {
		var post structs.Post
		err := rows.Scan(&post.ID, &post.Content, &post.ImagePath, &post.CreatedAt, &post.Username, &post.CommentCount, &post.LikeCount)
		if err != nil {
			logs.Errorf("Scan error: %q", err.Error())
			return nil, err
		}
		posts = append(posts, post)
	}
	return posts, nil
}

func GetMembers(groupid int) ([]structs.Gusers, error) {
	var adminid int

	rows, err := DB.Query(`
	SELECT p.id p.display_name, p.avatar
FROM profile p
JOIN groupmember ON p.id = groupmember.person_id
WHERE groupmember.group_id = ?;`, groupid)
	if err != nil {
		// anas
	}
	defer rows.Close()
	err = DB.QueryRow(`SELECT g.creator_id FROM group g WHERE g.id = ?;`, groupid).Scan(adminid)
	if err != nil {
		return []structs.Gusers{}, fmt.Errorf("error fetching user: %v", err)
	}
	var members []structs.Gusers
	for rows.Next() {
		var member structs.Gusers
		if err := rows.Scan(member.Uid, member.Name, member.Avatar); err != nil {
			logs.Errorf("Error scanning message: %q", err.Error())
			return nil, err
		}
		if member.Uid == adminid {
			member.Adm = true
		} else {
			member.Adm = false
		}
		members = append(members, member)
	}
	return members, nil
}

// func GetComments(pid string) ([]structs.CommentGet, error) {
// 	if pid == "" {
// 		return nil, nil
// 	}
// 	Pid, err := strconv.Atoi(pid)
// 	if err != nil {
// 		logs.Errorf("Error converting pid to int: %q", err.Error())
// 		return nil, err
// 	}
// 	rows, err := DB.Query(`
// 	SELECT
//     	u.username AS author,
//     	c.content,
//     	c.created_at
// 	FROM
//     	comments c
// 	JOIN
//     	users u ON c.uid = u.id
// 	WHERE
// 		c.post_id = ?
// 	ORDER BY
// 		c.created_at DESC
// 	`, Pid)
// 	if err != nil {
// 		logs.Errorf("Error getting comments: %q", err.Error())
// 		return nil, err
// 	}
// 	defer rows.Close()
// 	var comments []structs.CommentGet
// 	for rows.Next() {
// 		var comment structs.CommentGet
// 		err := rows.Scan(&comment.Author, &comment.Content, &comment.CreationTime)
// 		if err != nil {
// 			logs.Errorf("Error scanning comments: %q", err.Error())
// 			return nil, err
// 		}
// 		comment.Pid = structs.ID(Pid)
// 		comments = append(comments, comment)
// 	}
// 	return comments, nil
// }

// func GetUserNames(uid int) ([]structs.UsersGet, error) {
// 	rows, err := DB.Query(`
// 	SELECT
// 		u.username
// 	FROM
// 		users u
// 	LEFT JOIN
// 		dms m
// 	ON
// 		(u.id = m.sender_id OR u.id = m.recipient_id)
// 	AND
// 		(m.sender_id = ? OR m.recipient_id = ? )
// 	WHERE
// 		u.id != ?
// 	GROUP BY
// 		u.id, u.username
// 	ORDER BY
// 		CASE WHEN MAX(m.created_at) IS NOT NULL THEN 1
// 	    ELSE 2
// 	END,
// 		MAX(m.created_at) DESC,
// 	u.username ASC;`, uid, uid, uid)
// 	if err != nil {
// 		return nil, fmt.Errorf("could not execute query: %w", err)
// 	}
// 	defer rows.Close()
// 	var userNames []structs.UsersGet
// 	for rows.Next() {
// 		var username structs.UsersGet
// 		if err := rows.Scan(&username.Username); err != nil {
// 			return userNames, fmt.Errorf("could not scan row: %w", err)
// 		}
// 		// TODO
// 		// if _, e := helpers.Sockets[username.Username]; e {
// 		// 	username.Online = true
// 		// }
// 		// userNames = append(userNames, username)
// 	}
// 	if err := rows.Err(); err != nil {
// 		return userNames, fmt.Errorf("row iteration error: %w", err)
// 	}
// 	return userNames, nil
// }

// func GetdmHistory(uname1, uname2, date string) ([]structs.Message, error) {
// 	var d time.Time
// 	if date == "" {
// 		d = time.Now()
// 	} else {
// 		var err error
// 		d, err = time.Parse(time.RFC3339, date)
// 		if err != nil {
// 			return nil, fmt.Errorf("invalid date format")
// 		}
// 	}
// 	rows, err := DB.Query(`
//         SELECT *
//         FROM (
//             SELECT
//                 sender.username, d.message, d.created_at
//             FROM
//                 dms d
//             JOIN
//                 users sender ON d.sender_id = sender.id
//             JOIN
//                 users recipient ON d.recipient_id = recipient.id
//             WHERE
//                 d.created_at < ?
//                 AND (
//                     (sender.username = ? AND recipient.username = ?)
//                     OR
//                     (sender.username = ? AND recipient.username = ?)
//                 )
//             ORDER BY
//                 d.created_at DESC
//             LIMIT 10
//         ) AS sub
//         ORDER BY created_at ASC;
//     `, d, uname1, uname2, uname2, uname1)
// 	if err != nil {
// 		logs.Errorf("Error getting messages: %q", err.Error())
// 		return nil, err
// 	}
// 	defer rows.Close()
// 	var messages []structs.Message
// 	for rows.Next() {
// 		var message structs.Message
// 		if err := rows.Scan(&message.Sender, &message.Content, &message.Time); err != nil {
// 			logs.Errorf("Error scanning message: %q", err.Error())
// 			return nil, err
// 		}
// 		messages = append(messages, message)
// 	}
// 	return filter(messages, d), nil
// }
