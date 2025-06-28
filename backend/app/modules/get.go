package modules

import (
	"fmt"
	"time"

	"social-network/app/structs"
	"social-network/server/logs"
)

func GetPosts(start, uid, groupId int) ([]structs.Post, error) {
	query := `
	WITH
	    user_groups AS (
	        SELECT group_id
	        FROM groupmember
	        WHERE user_id = ?          -- <-- Current user ID
	          AND active = 1            -- <-- Must be an active member
	    ),
	    followed_profiles AS (
	        SELECT following_id
	        FROM follow
	        WHERE follower_id = ?       -- <-- Current user ID again
	          AND status = 1            -- <-- Follow relationship is accepted
	    )
	SELECT
	    p.id,
	    p.group_id,
	    p.user_id,
	    author.display_name AS UserName,
	    group_profile.display_name AS GroupName,
		author.avatar AS AvatarUser,
    	group_profile.avatar AS AvatarGroup,
	    p.content,
	    p.image_path,
	    p.created_at,
	    (
	        SELECT COUNT(*) FROM comments c WHERE c.post_id = p.id
	    ) AS comment_count,
	    (
	        SELECT COUNT(*) FROM likes l WHERE l.post_id = p.id
	    ) AS like_count
	FROM posts p
	JOIN profile author ON author.id = p.user_id
	LEFT JOIN profile group_profile ON group_profile.id = p.group_id
	LEFT JOIN user_groups ug ON p.group_id = ug.group_id
	WHERE
	    (? = 0 OR p.id < ?)  -- pagination
	    AND (? = 0 OR p.group_id = ?)  -- group filter
	    AND p.privacy != 'private'
	    AND (
	        p.privacy = 'public'
	        OR p.user_id = ?  -- current user is the author
	        OR (
	            p.privacy = 'friends'
	            AND p.user_id IN (SELECT following_id FROM followed_profiles)
	        )
	        OR (
	            p.group_id IS NOT NULL
	            AND ug.group_id IS NOT NULL  -- user is group member
	        )
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
		err := rows.Scan(
			&post.ID,
			&post.GroupId,
			&post.UserId,
			&post.UserName,
			&post.GroupName,
			&post.AvatarUser,
			&post.AvatarGroup,
			&post.Content,
			&post.ImagePath,
			&post.CreatedAt,
			&post.CommentCount,
			&post.LikeCount,
		)
		if err != nil {
			logs.Errorf("Scan error: %q", err.Error())
			return nil, err
		}
		// TODO:if post of group get name of group
		posts = append(posts, post)
	}

	return posts, nil
}

func GetMembers(groupid int) ([]structs.Gusers, error) {
	var adminid int

	rows, err := DB.Query(`    
    SELECT p.id p.display_name, p.avatar
    FROM profile p
    JOIN follow ON p.id = follow.follower_id
    WHERE follow.following_id = ?;`, groupid)
	if err != nil {
		logs.Errorf("GetMembers query error: %q", err.Error())
		return nil, err
	}
	defer rows.Close()
	err = DB.QueryRow(`SELECT g.creator_id FROM group g WHERE g.id = ?;`, groupid).Scan(adminid)
	if err != nil {
		return []structs.Gusers{}, fmt.Errorf("error fetching user: %v", err)
	}
	var admin structs.Gusers
	err = DB.QueryRow(`select p.id p.display_name, p.avatar from profile p where p.id = ?`, adminid).Scan(admin.Uid, admin.Name, admin.Avatar)
	if err != nil {
		//
	}
	var members []structs.Gusers
	members = append(members, admin)
	for rows.Next() {
		var member structs.Gusers
		if err := rows.Scan(member.Uid, member.Name, member.Avatar); err != nil {
			logs.Errorf("Error scanning message: %q", err.Error())
			return nil, err
		}
		members = append(members, member)
	}
	return members, nil
}

func GetGroupFeed(uid int) ([]structs.Post, error) {
	rows, err := DB.Query(`SELECT
	    sub.id,
	    sub.group_id,
	    sub.user_id,
	    sub.content,
	    author.display_name AS UserName,
	    group_profile.display_name AS GroupName,
	    author.avatar AS AvatarUser,
		group_profile.avatar AS AavatarGroup,
	    sub.image_path,
		sub.created_at,
	    sub.like_count,
	    sub.comment_count
	FROM (
	    SELECT
	        p.id,
			p.group_id,
	        p.user_id,
	        p.content,
	        p.image_path,
			p.created_at,
	        p.group_id,
	        (
	            SELECT COUNT(*) FROM likes l WHERE l.post_id = p.id
	        ) AS like_count,
	        (
	            SELECT COUNT(*) FROM comments c WHERE c.post_id = p.id
	        ) AS comment_count,
	        ROW_NUMBER() OVER (
	            PARTITION BY p.group_id
	            ORDER BY p.created_at DESC
	        ) AS rn
	    FROM posts p
	    JOIN groupmember gm ON p.group_id = gm.group_id
	    WHERE gm.user_id = ?
	) AS sub
	JOIN "group" g ON sub.group_id = g.id
	JOIN profile group_profile ON group_profile.id = g.id         -- group profile
	JOIN profile author ON author.id = sub.user_id               -- post author
	WHERE sub.rn <= 2;
`, uid)
	if err != nil {
		logs.Errorf("GetgroupFeed query error: %q", err.Error())
		return nil, err
	}

	var posts []structs.Post
	for rows.Next() {
		var pt structs.Post
		if err := rows.Scan(&pt.ID, &pt.GroupId, &pt.UserId, &pt.Content, &pt.UserName, &pt.GroupName, &pt.AvatarUser, &pt.AvatarGroup, &pt.ImagePath, &pt.CreatedAt, &pt.LikeCount, &pt.CommentCount); err != nil {
			logs.Errorf("Error scanning groups %q", err.Error())
			return nil, err
		}
		posts = append(posts, pt)
	}
	return posts, nil
}

func GetGroupToJoin(uid int) ([]structs.GroupGet, error) {
	rows, err := DB.Query(`SELECT 
	  p.id,
	  p.display_name,
	  p.avatar,
	  p.description
	FROM profile p
	JOIN "group" g ON p.id = g.id
	LEFT JOIN groupmember gm 
	  ON g.id = gm.group_id AND gm.user_id = ?
	WHERE gm.user_id IS NULL
	ORDER BY RANDOM()
	LIMIT 10;
	`, uid)
	if err != nil {
		logs.Errorf("GetGroupToJoin query error: %q", err.Error())
		return nil, err
	}

	var grs []structs.GroupGet

	for rows.Next() {
		var gr structs.GroupGet
		if err := rows.Scan(&gr.ID, &gr.GroupName, &gr.Avatar, &gr.Description); err != nil {
			logs.Errorf("Error scanning groups %q", err.Error())
			return nil, err
		}
		grs = append(grs, gr)
	}
	return grs, nil
}

func GetGroupImIn(uid int) ([]structs.GroupGet, error) {
	rows, err := DB.Query(`
	SELECT
	    p.id,
	    p.display_name,
	    p.avatar,
	    p.description
	FROM
	    profile p
	    JOIN "group" g ON p.id = g.id
	    JOIN groupmember gm ON g.id = gm.group_id
	WHERE
	    gm.user_id = ?
	ORDER BY
	    RANDOM()
	LIMIT
	    10;
	`, uid)
	if err != nil {
		logs.Errorf("GetGroupImIn query error: %q", err.Error())
		return nil, err
	}

	var grs []structs.GroupGet

	for rows.Next() {
		var gr structs.GroupGet
		if err := rows.Scan(&gr.ID, &gr.GroupName, &gr.Avatar, &gr.Description); err != nil {
			logs.Errorf("Error scanning groups %q", err.Error())
			return nil, err
		}
		grs = append(grs, gr)
	}
	return grs, nil
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

func GetUserNames(uid int) ([]structs.UsersGet, error) {
	rows, err := DB.Query(`
	SELECT
        p.id,
		p.display_name,
        NOT p.is_person AS is_group
	FROM
		user u
    JOIN
        profile p
	LEFT JOIN
		message m
	ON
		(u.id = m.sender_id OR u.id = m.receiver_id)
	AND
		(m.sender_id = ? OR m.receiver_id = ? )
	WHERE
		p.id != ?
	GROUP BY
		p.id, p.display_name
	ORDER BY
		CASE WHEN MAX(m.created_at) IS NOT NULL THEN 1
	    ELSE 2
	END,
		MAX(m.created_at) DESC,
	p.display_name ASC;`, uid, uid, uid)
	if err != nil {
		return nil, fmt.Errorf("could not execute query: %w", err)
	}
	defer rows.Close()

	var userS []structs.UsersGet

	for rows.Next() {
		var user structs.UsersGet
		if err := rows.Scan(&user.ID, &user.Username, &user.Is_Group); err != nil {
			return userS, fmt.Errorf("could not scan row: %w", err)
		}
		// TODO IMPLEMENT ONLINE STATUS
		// if _, e := helpers.Sockets[username.Username]; e {
		// 	username.Online = true
		// }
		userS = append(userS, user)
	}

	if err := rows.Err(); err != nil {
		return userS, fmt.Errorf("row iteration error: %w", err)
	}

	return userS, nil
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
                sender.id,sender.display_name, d.content, d.created_at
            FROM
                message d
            JOIN
                profile sender ON d.sender_id = sender.id
            JOIN
                profile recipient ON d.receiver_id = recipient.id
            WHERE
                d.created_at < ?
                AND (
                    (sender.display_name = ? AND recipient.display_name = ?)
                    OR
                    (sender.display_name = ? AND recipient.display_name = ?)
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
		if err := rows.Scan(&message.Sender, &message.SenderName, &message.Content, &message.Time); err != nil {
			logs.Errorf("Error scanning message: %q", err.Error())
			return nil, err
		}
		messages = append(messages, message)
	}
	return filter(messages, d), nil
}
