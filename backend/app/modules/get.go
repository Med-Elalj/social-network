package modules

import (
	"database/sql"
	"fmt"

	"social-network/app/logs"
	// "social-network/app/modules"
	"social-network/app/structs"
)

func GetPosts(start, uid, groupId, userId int) ([]structs.Post, error) {
	query := `
	WITH
	    user_groups AS (
	        SELECT g.id FROM "group" g WHERE g.creator_id = ?
	        UNION
	        SELECT g.id FROM "group" g JOIN follow f ON f.following_id = g.id
	        WHERE f.follower_id = ? AND f.status = 1
	    ),
	    followed_profiles AS (
	        SELECT following_id FROM follow
	        WHERE follower_id = ? AND status = 1
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
	    (SELECT COUNT(*) FROM comments c WHERE c.post_id = p.id) AS comment_count,
	    (SELECT COUNT(*) FROM likes l WHERE l.post_id = p.id) AS like_count,
		CASE 
	    WHEN EXISTS (
	        SELECT 1 FROM likes l 
	        WHERE l.user_id = ? AND l.post_id = p.id AND l.comment_id IS NULL
	    ) THEN 1
	    ELSE 0
	    END AS is_liked
	FROM posts p
	JOIN profile author ON author.id = p.user_id
	LEFT JOIN profile group_profile ON group_profile.id = p.group_id
	LEFT JOIN follow f ON p.group_id = f.follower_id
	WHERE
	    (? = 0 OR p.id < ?)
	    AND (? = 0 OR p.group_id = ?)
	    AND (? = 0 OR p.user_id = ?) -- condition to filter by user
	    AND p.privacy != 'private'
	    AND (
	        p.privacy = 'public'
	        OR p.user_id = ?
	        OR (p.privacy = 'friends' AND p.user_id IN (SELECT following_id FROM followed_profiles))
	        OR (p.group_id IS NOT NULL AND f.follower_id IS NOT NULL)
	    )
	ORDER BY p.id DESC
	LIMIT 10;
`

	rows, err := DB.Query(query,
		uid, uid, // user_groups
		uid,          // followed_profiles
		uid,          // is_liked check
		start, start, // pagination
		groupId, groupId, // group filter
		userId, userId, // user filter
		uid, // privacy condition
	)
	if err != nil {
		logs.ErrorLog.Printf("GetPosts query error: %q", err.Error())
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
			&post.IsLiked,
		)
		if err != nil {
			logs.ErrorLog.Printf("Scan error: %q", err.Error())
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, nil
}

// to do offset
func GetRequests(uid, tpdefind int) ([]structs.RequestsGet, error) {
	rows, err := DB.Query(`
	SELECT
	    r.sender_id,
	    r.target_id,
	    r.type,
	    COALESCE(pg.display_name, pe_group.display_name, ''),
	    COALESCE(pg.avatar, pe_group.avatar, ''),
	    r.created_at,
	    ps.display_name,
	    ps.avatar
	FROM request r
	JOIN profile ps ON r.sender_id = ps.id
	LEFT JOIN profile pg ON pg.id = r.target_id AND r.type = 1
	LEFT JOIN events e ON e.id = r.target_id AND r.type = 2
	LEFT JOIN profile pe_group ON pe_group.id = e.group_id
	WHERE
	    r.receiver_id = ? AND
	    (? = 3 OR r.type = ?)
	ORDER BY r.created_at DESC;
	`, uid, tpdefind, tpdefind)
	if err != nil {
		logs.ErrorLog.Printf("GetRequests query error: %q", err.Error())
		return nil, err
	}
	defer rows.Close()

	var requests []structs.RequestsGet
	for rows.Next() {
		var request structs.RequestsGet
		if err := rows.Scan(&request.SenderId, &request.GroupId, &request.Type, &request.GroupName, &request.GroupAvatar, &request.Time, &request.Username, &request.Avatar); err != nil {
			logs.ErrorLog.Printf("Error scanning requests: %q", err.Error())
			return nil, err
		}

		switch request.Type {
		case 0:
			request.Message = fmt.Sprintf("%s sent you a follow request", request.Username)
		case 1:
			request.Message = fmt.Sprintf("%s wants to join %s group", request.Username, request.GroupName)
		case 2:
			request.Message = fmt.Sprintf("%s create a new event on %s group", request.Username, request.GroupName)
		}

		requests = append(requests, request)
	}
	return requests, nil
}

// anas
func GetEvents(group_id int, uid int) ([]structs.GroupEvent, error) {
	rows, err := DB.Query(`    
	SELECT
		e.id,
	    e.user_id,
	    e.description,
		e.title,
		e.timeof,
		e.created_at,
		eu.respond
	FROM
	    "events" e
	    JOIN "group" g ON e.group_id = g.id
		JOIN userevents eu ON e.id = eu.event_id
	WHERE
	    g.id = ? AND eu.user_id = ?;`, group_id, uid, "event")
	if err != nil {
		logs.ErrorLog.Printf("Getevent query error: %q", err.Error())
		return nil, err
	}
	var events []structs.GroupEvent
	for rows.Next() {
		var event structs.GroupEvent
		if err := rows.Scan(event.ID, event.Userid, event.Description, event.Title, event.Timeof, event.CreationTime, event.Respond); err != nil {
			logs.ErrorLog.Printf("Error scanning events: %q", err.Error())
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}

func GetMembers(groupid int) ([]structs.Gusers, error) {
	rows, err := DB.Query(`    
	SELECT
	    p.id,
	    p.display_name,
	    p.avatar
	FROM
	    profile p
	    JOIN follow ON p.id = follow.follower_id
	WHERE
	    follow.following_id = ?;`, groupid)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logs.ErrorLog.Printf("GetMembers query error: %q", err.Error())
		return nil, err
	}
	defer rows.Close()

	var admin structs.Gusers
	err = DB.QueryRow(`
		select
		    p.id,
		    p.display_name,
		    p.avatar
		from
		    profile p
		    JOIN "group" g ON g.creator_id = p.id
		where
		    g.id = ?;`, groupid).Scan(&admin.Uid, &admin.Name, &admin.Avatar)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		logs.ErrorLog.Printf("GetMembers query error: %q", err.Error())
		return nil, err
	}

	var members []structs.Gusers
	members = append(members, admin)
	for rows.Next() {
		var member structs.Gusers
		if err := rows.Scan(&member.Uid, &member.Name, &member.Avatar); err != nil {
			logs.ErrorLog.Printf("Error scanning message: %q", err.Error())
			return nil, err
		}
		members = append(members, member)
	}
	return members, nil
}

func GetGroupFeed(uid int) ([]structs.Post, error) {
	rows, err := DB.Query(`
		WITH
		    user_groups AS (
		        SELECT
		            id
		        FROM
		            "group"
		        WHERE
		            creator_id = ?
		        UNION
		        SELECT
		            following_id
		        FROM
		            follow
		        WHERE
		            follower_id = ?
		            AND status = 1
		    ),
		    posts_with_rn AS (
		        SELECT
		            p.*,
		            ROW_NUMBER() OVER (
		                PARTITION BY
		                    p.group_id
		                ORDER BY
		                    p.created_at DESC
		            ) AS rn
		        FROM
		            posts p
		            JOIN user_groups ug ON p.group_id = ug.id
		    )
		SELECT
		    pwr.id,
		    pwr.group_id,
		    pwr.user_id,
		    pwr.content,
		    author.display_name AS UserName,
		    group_profile.display_name AS GroupName,
		    author.avatar AS AvatarUser,
		    group_profile.avatar AS AvatarGroup,
		    pwr.image_path,
		    pwr.created_at,
		    (
		        SELECT
		            COUNT(*)
		        FROM
		            likes l
		        WHERE
		            l.post_id = pwr.id
		    ) AS like_count,
		    (
		        SELECT
		            COUNT(*)
		        FROM
		            comments c
		        WHERE
		            c.post_id = pwr.id
		    ) AS comment_count,
		    CASE
		        WHEN EXISTS (
		            SELECT
		                1
		            FROM
		                likes l
		            WHERE
		                l.user_id = ?
		                AND l.post_id = pwr.id
		                AND l.comment_id IS NULL
		        ) THEN 1
		        ELSE 0
		    END AS is_liked
		FROM
		    posts_with_rn pwr
		    JOIN profile author ON author.id = pwr.user_id
		    JOIN profile group_profile ON group_profile.id = pwr.group_id
		WHERE
		    pwr.rn <= 2
		ORDER BY
		    pwr.group_id,
		    pwr.created_at DESC;`, uid, uid, uid)
	if err != nil {
		logs.ErrorLog.Printf("GetgroupFeed query error: %q", err.Error())
		return nil, err
	}

	var posts []structs.Post
	for rows.Next() {
		var pt structs.Post
		if err := rows.Scan(&pt.ID, &pt.GroupId, &pt.UserId, &pt.Content, &pt.UserName, &pt.GroupName, &pt.AvatarUser, &pt.AvatarGroup, &pt.ImagePath, &pt.CreatedAt, &pt.LikeCount, &pt.CommentCount, &pt.IsLiked); err != nil {
			logs.ErrorLog.Printf("Error scanning groups %q", err.Error())
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
    p.description,
    CASE
        WHEN EXISTS (
            SELECT 1
            FROM request r
            WHERE r.sender_id = ?
              AND r.target_id = p.id
              AND r.type = 1
        ) THEN 1
        ELSE 0
    END AS is_requested
FROM
    profile p
    JOIN "group" g ON p.id = g.id
WHERE
    p.is_user = 0
    AND p.id NOT IN (
        SELECT g2.id FROM "group" g2 WHERE g2.creator_id = ?
        UNION
        SELECT f.following_id FROM follow f
        WHERE f.follower_id = ? AND f.status = 1
    )
LIMIT 10;`, uid, uid, uid)
	if err != nil {
		logs.ErrorLog.Printf("GetGroupToJoin query error: %q", err.Error())
		return nil, err
	}

	var grs []structs.GroupGet

	for rows.Next() {
		var gr structs.GroupGet
		if err := rows.Scan(&gr.ID, &gr.GroupName, &gr.Avatar, &gr.Description, &gr.IsRequested); err != nil {
			logs.ErrorLog.Printf("Error scanning groups %q", err.Error())
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
		WHERE
		    p.is_user = 0
		    AND (
		        g.creator_id = ? -- you are the creator
		        OR p.id IN (    -- or you follow this group
		            SELECT f.following_id
		            FROM follow f
		            WHERE f.follower_id = ? AND f.status = 1
		        )
		    )
		LIMIT 10;`, uid, uid)
	if err != nil {
		logs.ErrorLog.Printf("GetGroupImIn query error: %q", err.Error())
		return nil, err
	}

	var grs []structs.GroupGet

	for rows.Next() {
		var gr structs.GroupGet
		if err := rows.Scan(&gr.ID, &gr.GroupName, &gr.Avatar, &gr.Description); err != nil {
			logs.ErrorLog.Printf("Error scanning groups %q", err.Error())
			return nil, err
		}
		grs = append(grs, gr)
	}
	return grs, nil
}

func GetUserNames(uid int) ([]structs.UsersGet, error) {
	rows, err := DB.Query(`
	SELECT
		p.id,
		p.display_name,
		p.avatar,
		NOT p.is_user AS is_group
	FROM
		profile p
	JOIN
		user u ON u.id = p.id
	INNER JOIN 
		follow f ON (f.follower_id = ? OR f.following_id = ?)
		AND (f.follower_id = p.id OR f.following_id = p.id)
	LEFT JOIN
		message m ON (
			(u.id = m.sender_id OR u.id = m.receiver_id)
			AND (m.sender_id = ? OR m.receiver_id = ?)
		)
	WHERE
		p.id != ?
		AND p.is_user = 1
	GROUP BY
		p.id, p.display_name
	ORDER BY
		CASE WHEN MAX(m.created_at) IS NOT NULL THEN 1 ELSE 2 END,
		MAX(m.created_at) DESC,
		p.display_name ASC;
`, uid, uid, uid, uid, uid)
	if err != nil {
		return nil, fmt.Errorf("could not execute query: %w", err)
	}
	defer rows.Close()

	var userS []structs.UsersGet

	for rows.Next() {
		var user structs.UsersGet
		if err := rows.Scan(&user.ID, &user.Username, &user.Avatar, &user.Is_Group); err != nil {
			return userS, fmt.Errorf("could not scan row: %w", err)
		}
		// TODO IMPLEMENT ONLINE STATUS
		_, user.Online = structs.Sockets[int(user.ID)]
		userS = append(userS, user)
	}

	if err := rows.Err(); err != nil {
		return userS, fmt.Errorf("row iteration error: %w", err)
	}

	return userS, nil
}

func GetdmHistory(uname1, uname2 string, page int) (structs.Chat, error) {
	// var d time.Time
	var chat structs.Chat

	pageSize := 10

	offset := (page - 1) * pageSize

	fmt.Println("uname1: ", uname1)
	fmt.Println("uname2: ", uname2)
	fmt.Println("page: ", page)

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
                (sender.display_name = ? AND recipient.display_name = ?)
                OR
                (sender.display_name = ? AND recipient.display_name = ?)
            
            ORDER BY
                d.created_at DESC
            LIMIT 11 OFFSET ?
        ) AS sub
        ORDER BY created_at ASC;
    `, uname1, uname2, uname2, uname1, offset)
	if err != nil {
		logs.ErrorLog.Printf("Error getting messages: %q", err.Error())
		return chat, err
	}
	defer rows.Close()

	// var messages []structs.Message
	var count int
	for rows.Next() {
		if count == 10 {
			chat.HasMore = true
			break
		}

		var message structs.Message
		if err := rows.Scan(&message.Sender, &message.SenderName, &message.Content, &message.Time); err != nil {
			logs.ErrorLog.Printf("Error scanning message: %q", err.Error())
			return chat, err
		}
		chat.Messages = append(chat.Messages, message)
		count++
	}
	fmt.Println("before", chat.Messages)

	return chat, nil
}

func GetSearchprofile(query string, page int) (structs.SearchProfile, error) {
	offset := (page - 1) * 10
	rows, err := DB.Query(`
	SELECT
		p.id,
		p.display_name,
		p.avatar,
		p.is_user
	FROM
		profile p
	WHERE
		p.display_name LIKE ?
	ORDER BY p.display_name ASC
	LIMIT 11 OFFSET ?;`, "%"+query+"%", offset)
	if err != nil {
		logs.ErrorLog.Printf("GetSearchProfile query error: %q", err.Error())
		return structs.SearchProfile{}, err
	}
	defer rows.Close()

	var profiles []structs.UsersGet
	var rtn structs.SearchProfile
	rtn.HasMore = false
	for i := 0; rows.Next(); i++ {
		if i == 10 {
			rtn.HasMore = true
			break
		}
		var profile structs.UsersGet
		if err := rows.Scan(&profile.ID, &profile.Username, &profile.Avatar, &profile.Is_Group); err != nil {
			logs.ErrorLog.Printf("Error scanning profile: %q", err.Error())
			return structs.SearchProfile{}, err
		}
		profile.Is_Group = !profile.Is_Group
		profiles = append(profiles, profile)

	}
	rtn.Profiles = profiles
	return rtn, nil
}

// func GetSuggestions(uid int, is_user int) ([]structs.UsersGet, error) {
// 	var users []structs.UsersGet

// 	// Query to get all profiles excluding:
// 	// 1. The user themselves
// 	// 2. Users they already follow
// 	// 3. Users they have sent/received requests to/from
// 	query := `
//         SELECT p.id, p.avatar, p.display_name, p.is_user
//         FROM profile p
//         WHERE p.id != ?
//         AND p.is_user = ?
//         AND p.id NOT IN (
//             -- Exclude users already being followed
//             SELECT f.following_id
//             FROM follow f
//             WHERE f.follower_id = ?
//         )
//         AND p.id NOT IN (
//             -- Exclude users with pending requests (as sender)
//             SELECT r.target_id
//             FROM request r
//             WHERE r.sender_id = ?
//         )
//         AND p.id NOT IN (
//             -- Exclude users with pending requests (as receiver)
//             SELECT r.sender_id
//             FROM request r
//             WHERE r.target_id = ?
//         )
//         ORDER BY p.created_at DESC
//         LIMIT 20`

// 	rows, err := DB.Query(query, uid, is_user, uid, uid, uid)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to query suggestions: %v", err)
// 	}
// 	defer rows.Close()

// 	for rows.Next() {
// 		var user structs.UsersGet
// 		var isUser bool

// 		err := rows.Scan(&user.ID, &user.Avatar, &user.Username, &isUser)
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to scan user: %v", err)
// 		}

// 		// Set Is_Group based on is_user field (inverted)
// 		user.Is_Group = !isUser

// 		// Set online status (you might want to implement this based on your logic)
// 		user.Online = false // Default to false, implement your online logic here

// 		// Get relationship status
// user.FollowStatus, err = GetRelationship(uid, int(user.ID))
// if err != nil {
// 	return nil, fmt.Errorf("failed to get relation ship: %v", err)
// }

// 		users = append(users, user)
// 	}

// 	if err = rows.Err(); err != nil {
// 		return nil, fmt.Errorf("error iterating rows: %v", err)
// 	}

// 	return users, nil
// }

func GetSuggestions(uid int, Type int) ([]structs.UsersGet, error) {
	var users []structs.UsersGet

	query := `
        SELECT p.id, p.avatar, p.display_name, p.is_user
        FROM profile p
        WHERE p.id != ?
        AND p.is_user = ?
        AND p.id NOT IN (
            -- Exclude users where uid is follower (following them)
            SELECT f.following_id 
            FROM follow f 
            WHERE f.follower_id = ?
        )
        AND p.id NOT IN (
            -- Exclude users where uid is sender in request
            SELECT r.target_id 
            FROM request r 
            WHERE r.sender_id = ?
        )
        AND p.id NOT IN (
            -- Exclude users where uid is receiver/target in request
            SELECT r.sender_id 
            FROM request r 
            WHERE r.target_id = ?
        )
        ORDER BY p.created_at DESC
        LIMIT 20`

	rows, err := DB.Query(query, uid, Type, uid, uid, uid)
	if err != nil {
		return nil, fmt.Errorf("failed to query suggestions: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var user structs.UsersGet
		var isUser bool

		err := rows.Scan(&user.ID, &user.Avatar, &user.Username, &isUser)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %v", err)
		}

		// Set Is_Group based on is_user field (inverted)
		user.Is_Group = !isUser

		// Set online status
		user.Online = false

		// Get relationship status
		user.FollowStatus, err = GetRelationship(uid, int(user.ID))
		if err != nil {
			return nil, fmt.Errorf("failed to get relation ship: %v", err)
		}

		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %v", err)
	}

	return users, nil
}
