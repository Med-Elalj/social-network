.
|____output.md
|____backend
| |____init.py
| |____readme.md
| |____app
| | |____ws
| | | |____ws.go
| | |____modules
| | | |____commets.go
| | | |____get.go
| | | |____azer.go
| | | |____files.go
| | | |____upload.go
| | | |____db.go
| | | |____like.go
| | | |____follow.go
| | | |____dms.go
| | | |____set.go
| | |____Auth
| | | |____Auth.go
| | | |____jwt
| | | | |____generate.go
| | | | |____utils.go
| | | |____ValidationInsertion.go
| | | |____session.go
| | | |____middleware
| | | | |____middleware.go
| | | |____authHelpers.go
| | |____index_handler.go
| | |____upload
| | | |____upload.go
| | | |____svg_check.go
| | |____structs
| | | |____output.go
| | | |____structs.go
| | |____routes.go
| | |____handlers
| | | |____GetSet.go
| | | |____Profile
| | | | |____Profile.go
| | | | |____ProfileSettings.go
| | | |____Auth
| | | | |____register.go
| | | | |____login.go
| | | | |____refreshToken.go
| | | |____Like.go
| | | |____Group.go
| | | |____upload.go
| | | |____Comments.go
| | | |____refreshToken.go
| | | |____follow.go
| | | |____Post.go
| | |____requests
| | | |____commentCreation.go
| | |____docs.md
| |____main.go
| |____server
| | |____sql
| | | |____migrations
| | | | |____013_create_eventuser_table.up.sql
| | | | |____001_create_user_table.up.sql
| | | | |____006_create_comments_table.down.sql
| | | | |____007_create_message_table.down.sql
| | | | |____005_create_posts_table.up.sql
| | | | |____008_create_request_table.up.sql
| | | | |____009_create_sessions_table.up.sql
| | | | |____012_create_event_table.up.sql
| | | | |____011_create_dm_trigger.down.sql
| | | | |____003_create_group_table.up.sql
| | | | |____008_create_request_table.down.sql
| | | | |____004_create_follow_table.up.sql
| | | | |____010_create_likes_table.down.sql
| | | | |____011_create_dm_trigger.up.sql
| | | | |____012_create_event_table.down.sql
| | | | |____010_create_likes_table.up.sql
| | | | |____002_create_profile_table.up.sql
| | | | |____007_create_message_table.up.sql
| | | | |____001_create_user_table.down.sql
| | | | |____006_create_comments_table.up.sql
| | | | |____004_create_follow_table.down.sql
| | | | |____013_create_eventuser_table.down.sql
| | | | |____009_create_sessions_table.down.sql
| | | | |____005_create_posts_table.down.sql
| | | | |____002_create_profile_table.down.sql
| | | | |____003_create_group_table.down.sql
| | |____logs
| | | |____fatal.log
| | | |____app.log
| | | |____error.log
| | | |____logs.go
| | |____db
| | | |____main.db-shm
| | | |____uploads
| | | | |____Untitled.png
| | | | |____123.png
| | | | |____Screenshot 2024-11-25 185542.png
| | | | |____Untitled.jpeg
| | | | |____1233.png
| | | | |____Screenshot 2024-11-25 185645.png
| | | |____main.db-wal
| | | |____main.db
| |____go.sum
| |____go.mod
| |____bannerBG.png
|____Makefile
|____package-lock.json
|____front-end
| |____public
| | |____groupe2.svg
| | |____404.svg
| | |____home2.svg
| | |____groupsBg.png
| | |____reject.svg
| | |____sounds
| | | |____notification.mp3
| | | |____alert.mp3
| | |____iconFemale.png
| | |____upload.svg
| | |____Image.svg
| | |____send.svg
| | |____comment.svg
| | |____unread.svg
| | |____newMessage.svg
| | |____messages.svg
| | |____Like2.svg
| | |____postGroups.svg
| | |____accept.svg
| | |____messages2.svg
| | |____exit.svg
| | |____addUser.svg
| | |____posts2.svg
| | |____posts.svg
| | |____uploads
| | |____discover.svg
| | |____home.svg
| | |____public.svg
| | |____groupe.svg
| | |____create.svg
| | |____Like.svg
| | |____bannerBG.png
| | |____iconMale.png
| | |____notification.svg
| | |____private.svg
| |____utils
| | |____sendData.js
| |____eslint.config.mjs
| |____next-env.d.ts
| |____next.config.ts
| |____tsconfig.json
| |____package-lock.json
| |____src
| | |____app
| | | |____page.jsx
| | | |____newPost
| | | | |____page.jsx
| | | | |____newPost.module.css
| | | |____register
| | | | |____page.jsx
| | | | |____register.module.css
| | | |____chat
| | | | |____[tab]
| | | | | |____Unread.jsx
| | | | | |____Groups.jsx
| | | | | |____Users.jsx
| | | | |____page.jsx
| | | | |____input.jsx
| | | | |____time.jsx
| | | | |____messages.jsx
| | | | |____chat.module.css
| | | |____utils.jsx
| | | |____comments.jsx
| | | |____groupes
| | | | |____[tab]
| | | | | |____CreateGroup.jsx
| | | | | |____page.jsx
| | | | | |____GroupPosts.jsx
| | | | | |____Discover.jsx
| | | | | |____YourGroups.jsx
| | | | |____groups.module.css
| | | |____layout.jsx
| | | |____components
| | | | |____footer
| | | | | |____page.jsx
| | | | |____Friends.jsx
| | | | |____Groups.jsx
| | | | |____navigation
| | | | | |____page.jsx
| | | | | |____nav.module.css
| | | | |____Logout.jsx
| | | | |____Discover.jsx
| | | |____global.module.css
| | | |____login
| | | | |____page.jsx
| | | | |____login.module.css
| | | |____profile
| | | | |____page.jsx
| | | | |____[nickname]
| | | | | |____[tab]
| | | | | | |____Posts.jsx
| | | | | | |____Following.jsx
| | | | | | |____Followers.jsx
| | | | | | |____Settings.jsx
| | | | | |____page.jsx
| | | | |____profile.module.css
| | | |____context
| | | | |____notificationContext.jsx
| | | | |____WebSocketContext.jsx
| | | | |____AuthContext.jsx
| | | |____global.css
| | | |____not-found.jsx
| |____package.json
|____private
| |____cert.pem
| |____key.pem
|____Dockerfile
|____deprecated
| |____profilePage.jsx.old
|____logs
