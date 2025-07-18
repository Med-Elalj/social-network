package modules

import (
	"database/sql"
	"fmt"

	"social-network/app/logs"
	"social-network/app/structs"
)

// function to get group page posts
func GetGroupPosts(start, uid, groupId int) ([]structs.Post, error) {
	query := `
	SELECT
	    p.id,
	    p.group_id,
	    p.user_id,
	    creator.display_name,
	    pg.display_name,
	    creator.avatar,
	    pg.avatar,
	    p.content,
	    p.privacy,
	    p.image_path,
	    p.created_at,
	    (SELECT COUNT(*) FROM comments c WHERE c.post_id = p.id) AS CommentCount,
	    (SELECT COUNT(*) FROM likes l WHERE l.post_id = p.id) AS LikeCount,
	    CASE 
	        WHEN EXISTS (
	            SELECT 1 FROM likes l 
	            WHERE l.post_id = p.id AND l.user_id = ?
	        ) THEN 1 ELSE 0 
	    END AS IsLiked
	FROM posts p
	JOIN profile creator ON p.user_id = creator.id
	LEFT JOIN profile pg ON p.group_id = pg.id
	WHERE 
	    p.group_id = ?
	    AND (? = 0 OR p.id < ?)
	    AND (
	        EXISTS (
	            SELECT 1 FROM follow f
	            WHERE f.follower_id = ? AND f.following_id = p.group_id
	        )
	        OR EXISTS (
	            SELECT 1 FROM "group" g
	            WHERE g.id = p.group_id AND g.creator_id = ?
	        )
	    )
	ORDER BY p.created_at DESC
	LIMIT 10;`
	rows, err := DB.Query(query, uid, groupId, start, start, uid, uid)
	if err != nil {
		logs.ErrorLog.Printf("GetGroupPosts query error: %q", err.Error())
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
			&post.Privacy,
			&post.ImagePath,
			&post.CreatedAt,
			&post.CommentCount,
			&post.LikeCount,
			&post.IsLiked,
		)
		if err != nil {
			logs.ErrorLog.Printf("Error scanning post: %q", err.Error())
			return nil, err
		}
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		logs.ErrorLog.Printf("Error iterating rows: %q", err.Error())
		return nil, err
	}
	return posts, nil
}

// function to get home page posts
func GetHomePosts(start, uid int) ([]structs.Post, error) {
	query := `
	SELECT
		p.id AS ID,
		p.group_id AS GroupId,
		p.user_id AS UserId,
		creator.display_name AS UserName,
		pg.display_name AS GroupName,
		creator.avatar AS AvatarUser,
		pg.avatar AS AvatarGroup,
		p.content AS Content,
		p.Privacy AS Privacy,
		p.image_path AS ImagePath,
		p.created_at AS CreatedAt,
		(
			SELECT
				COUNT(*)
			FROM
				comments c
			WHERE
				c.post_id = p.id
		) AS CommentCount,
		(
			SELECT
				COUNT(*)
			FROM
				likes l
			WHERE
				l.post_id = p.id
		) AS LikeCount,
		CASE
			WHEN EXISTS (
				SELECT
					1
				FROM
					likes l
				WHERE
					l.post_id = p.id
					AND l.user_id = :me
			) THEN 1
			ELSE 0
		END AS IsLiked
	FROM
		posts p
		JOIN profile creator ON p.user_id = creator.id
		LEFT JOIN profile pg ON p.group_id = pg.id
	WHERE
		(
			:last_post_id = 0
			OR p.id < :last_post_id
		)
		AND (
			-- Public posts (not in group)
			(
				p.group_id IS NULL
				AND p.privacy = 'public'
			)
			-- Almost private posts (only if I'm following the user)
			OR (
				p.group_id IS NULL
				AND p.privacy = 'almost_private'
				AND EXISTS (
					SELECT
						1
					FROM
						follow f
					WHERE
						f.following_id = p.user_id
						AND f.follower_id = :me
				)
			)
			-- Private posts: only if mine OR shared with me via postrack
			OR (
				p.group_id IS NULL
				AND p.privacy = 'private'
				AND (
					p.user_id = :me
					OR EXISTS (
						SELECT
							1
						FROM
							postrack pt
						WHERE
							pt.post_id = p.id
							AND pt.follower_id = :me
					)
				)
			)
			-- Group posts: only if I'm in the group (follow the group)
			OR (
				p.group_id IS NOT NULL
				AND (
					EXISTS (
						SELECT
							1
						FROM
							follow f
						WHERE
							f.following_id = p.group_id
							AND f.follower_id = :me
					)
					OR EXISTS (
						SELECT
							1
						FROM
							"group" g
						WHERE
							g.id = p.group_id
							AND g.creator_id = :me
					)
				)
			)
		)
	ORDER BY
		p.created_at DESC
	LIMIT
		10;`
	rows, err := DB.Query(query, sql.Named("me", uid), sql.Named("last_post_id", start))
	if err != nil {
		logs.ErrorLog.Printf("GetHomePosts query error: %q", err.Error())
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
			&post.Privacy,
			&post.ImagePath,
			&post.CreatedAt,
			&post.CommentCount,
			&post.LikeCount,
			&post.IsLiked,
		)
		if err != nil {
			logs.ErrorLog.Printf("Error scanning post: %q", err.Error())
			return nil, err
		}
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		logs.ErrorLog.Printf("Error iterating rows: %q", err.Error())
		return nil, err
	}
	return posts, nil
}

// function to get profile page posts
func GetProfilePosts(start int, uid int, userId int) ([]structs.Post, error) {
	query := `
	SELECT 
	    p.id,
	    p.group_id,
	    p.user_id,
	    author.display_name,
	    pg.display_name,
	    author.avatar,
	    pg.avatar,
	    p.content,
	    p.privacy,
	    p.image_path,
	    p.created_at,
	    (SELECT COUNT(*) FROM comments c WHERE c.post_id = p.id) AS CommentCount,
	    (SELECT COUNT(*) FROM likes l WHERE l.post_id = p.id) AS LikeCount,
	    CASE WHEN EXISTS (
	        SELECT 1 FROM likes l 
	        WHERE l.post_id = p.id AND l.user_id = ?
	    ) THEN 1 ELSE 0 END AS IsLiked
	FROM posts p
	JOIN profile author ON p.user_id = author.id
	LEFT JOIN profile pg ON p.group_id = pg.id
	WHERE 
	    p.user_id = ?
	    AND (? = 0 OR p.id < ?)
	    AND (
	        (author.is_public = 1 AND p.privacy = 'public')
	        OR (
	            EXISTS (
	                SELECT 1 FROM follow 
	                WHERE follower_id = ? AND following_id = ?
	            ) AND p.privacy IN ('public', 'almost_private')
	        )
	        OR (
	            p.privacy = 'private'
	            AND EXISTS (
	                SELECT 1 FROM postrack 
	                WHERE post_id = p.id AND follower_id = ?
	            )
	        )
	    )
	ORDER BY p.created_at DESC
	LIMIT 10;`
	rows, err := DB.Query(query, uid, userId, start, start, uid, userId, uid)
	if err != nil {
		logs.ErrorLog.Printf("GetProfilePosts query error: %q", err.Error())
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
			&post.Privacy,
			&post.ImagePath,
			&post.CreatedAt,
			&post.CommentCount,
			&post.LikeCount,
			&post.IsLiked,
		)
		if err != nil {
			logs.ErrorLog.Printf("Error scanning post: %q", err.Error())
			return nil, err
		}
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		logs.ErrorLog.Printf("Error iterating rows: %q", err.Error())
		return nil, err
	}
	return posts, nil
}

// function to get own profile page posts
func GetOwnProfilePosts(start int, uid int) ([]structs.Post, error) {
	query := `
	SELECT
	    p.id,
	    p.group_id,
	    p.user_id,
	    creator.display_name,
	    pg.display_name,
	    creator.avatar,
	    pg.avatar,
	    p.content,
	    p.image_path,
	    p.created_at,
	    (
	        SELECT
	            COUNT(*)
	        FROM
	            comments c
	        WHERE
	            c.post_id = p.id
	    ) AS CommentCount,
	    (
	        SELECT
	            COUNT(*)
	        FROM
	            likes l
	        WHERE
	            l.post_id = p.id
	    ) AS LikeCount,
	    CASE
	        WHEN EXISTS (
	            SELECT
	                1
	            FROM
	                likes l
	            WHERE
	                l.post_id = p.id
	                AND l.user_id = ?
	        ) THEN 1
	        ELSE 0
	    END AS IsLiked
	FROM
	    posts p
	    JOIN profile creator ON p.user_id = creator.id
	    LEFT JOIN profile pg ON p.group_id = pg.id
	WHERE
	    p.user_id = ?
	    AND (
	        ? = 0
	        or p.id < ?
	    )
	ORDER BY
	    p.created_at DESC
	LIMIT
	    10;`
	rows, err := DB.Query(query, uid, uid, start, start)
	if err != nil {
		logs.ErrorLog.Printf("GetOwnProfilePosts query error: %q", err.Error())
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
			logs.ErrorLog.Printf("Error scanning post: %q", err.Error())
			return nil, err
		}
		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		logs.ErrorLog.Printf("Error iterating rows: %q", err.Error())
		return nil, err
	}

	return posts, nil
}

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

// function to get group events
func GetEvents(group_id int, uid int) ([]structs.GroupEvent, error) {
	rows, err := DB.Query(`    
	SELECT
	    e.id,
	    e.user_id,
	    e.description,
	    e.title,
	    e.timeof,
	    COALESCE(eu.respond, 0) as respond
	FROM
	    events e
	LEFT JOIN userevents eu ON e.id = eu.event_id AND eu.user_id = ?
	WHERE e.group_id = ?
	ORDER BY e.timeof ASC;`, uid, group_id)
	if err != nil {
		logs.ErrorLog.Printf("Getevent query error: %q", err.Error())
		return nil, err
	}
	var events []structs.GroupEvent
	for rows.Next() {
		var event structs.GroupEvent
		if err := rows.Scan(&event.ID, &event.Userid, &event.Description, &event.Title, &event.Timeof, &event.Respond); err != nil {
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
        OR p.id IN ( -- or you follow this group
            SELECT
                f.following_id
            FROM
                follow f
            WHERE
                f.follower_id = ?
                AND p.is_user = 0
        )
    )
LIMIT
    10;`, uid, uid)
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
	LEFT JOIN
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

func GetdmHistory(uid1 int, uname1, uname2 string, page int) (structs.Chat, error) {
	// var d time.Time
	var chat structs.Chat

	pageSize := 10

	offset := (page - 1) * pageSize

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
    		WHERE(
                (
                    recipient.is_user
                    AND
                    (? IN (sender.display_name ,recipient.display_name))
                    AND
                    (? IN (sender.display_name ,recipient.display_name))
                ) OR 
    		    (
                    recipient.display_name = ? 
                    AND (NOT recipient.is_user)
                    AND (
    		        EXISTS (
    		            SELECT 1 FROM follow f WHERE f.follower_id = ? 
    		            AND f.following_id = recipient.id
    		        )
                )
    		    )
    		)
    		ORDER BY

    		    d.created_at DESC
    		LIMIT 11 OFFSET ?
        ) AS sub
        ORDER BY created_at ASC;
    `, uname1, uname2, uname2, uid1, offset)
	if err != nil {
		logs.ErrorLog.Printf("Error getting messages: %q", err.Error())
		return chat, err
	}
	defer rows.Close()

	// var messages []structs.Message
	for rows.Next() {

		var message structs.Message
		if err := rows.Scan(&message.Sender, &message.SenderName, &message.Content, &message.Time); err != nil {
			logs.ErrorLog.Printf("Error scanning message: %q", err.Error())
			return chat, err
		}
		chat.Messages = append(chat.Messages, message)
	}

	if len(chat.Messages) == 11 {
		chat.HasMore = true
		chat.Messages = chat.Messages[1:]
	}
	return chat, nil
}

func GetSearchprofile(query string, page, groupId, uid int) (structs.SearchProfile, error) {
	offset := (page - 1) * 10
	var Query string
	var rows *sql.Rows
	var err error
	if groupId != 0 {
		Query = `
	SELECT
	    p.id,
	    p.display_name,
	    p.avatar,
	    p.is_user
	FROM
	    profile p
	JOIN 
	    follow f ON p.id = f.follower_id
	WHERE
	    f.following_id = ?  -- Your user ID (the person being followed)
	    AND p.display_name LIKE ?
	    AND p.is_user = 1   -- Only users (not groups)
	ORDER BY p.display_name ASC
	LIMIT 11 OFFSET ?;`
		rows, err = DB.Query(Query, uid, "%"+query+"%", offset)

	} else {
		Query = `
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
	LIMIT 11 OFFSET ?;`
		rows, err = DB.Query(Query, "%"+query+"%", offset)
	}
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

func GetSuggestions(uid int, is_user int) ([]structs.UsersGet, error) {
	var users []structs.UsersGet

	query := `
    	SELECT p.id, p.avatar, p.display_name, p.is_user,p.description
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
    	AND p.id NOT IN (
    	    -- Exclude groups where uid is creator
    	    SELECT g.id 
    	    FROM "group" g 
    	    WHERE g.creator_id = ?
    	)
    	ORDER BY p.created_at DESC
    	LIMIT 20`

	rows, err := DB.Query(query, uid, is_user, uid, uid, uid, uid)
	if err != nil {
		return nil, fmt.Errorf("failed to query suggestions: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var user structs.UsersGet
		var isUser bool

		err := rows.Scan(&user.ID, &user.Avatar, &user.Username, &isUser, &user.Description)
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

// UserInfoForNotification retrieves id, display_name, and a message for a given user id from profile table
func UserInfoForNotification(senderId, receiverId, targetId int) (structs.UserNotification, error) {
	var info structs.UserNotification
	row := DB.QueryRow(`SELECT id, display_name, is_user FROM profile WHERE id = ?`, senderId)
	err := row.Scan(&info.Sender.ID, &info.Sender.DisplayName, &info.Sender.IsUser)
	if err != nil {
		return info, err
	}

	row = DB.QueryRow(`SELECT id, display_name, is_user FROM profile WHERE id = ?`, receiverId)
	err = row.Scan(&info.Receiver.ID, &info.Receiver.DisplayName, &info.Receiver.IsUser)
	if err != nil {
		return info, err
	}

	row = DB.QueryRow(`SELECT id, display_name, is_user FROM profile WHERE id = ?`, targetId)
	err = row.Scan(&info.Target.ID, &info.Target.DisplayName, &info.Target.IsUser)
	if err != nil {
		return info, err
	}

	return info, nil
}
