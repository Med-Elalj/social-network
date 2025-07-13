-- Profile Data (Users and Groups)
INSERT INTO profile (id, email, first_name, last_name, display_name, date_of_birth, gender, avatar, description, is_public, is_user, created_at) VALUES
(1, 'john.doe@example.com', 'John', 'Doe', 'johndoe', '1990-05-15', 'male', 'avatars/john.jpg', 'Software engineer and photography enthusiast', 1, 1, '2023-01-10 09:15:22'),
(2, 'jane.smith@example.com', 'Jane', 'Smith', 'janesmith', '1988-11-23', 'female', 'avatars/jane.png', 'Digital marketer and travel blogger', 1, 1, '2023-01-12 14:30:45'),
(3, 'mike.johnson@example.com', 'Mike', 'Johnson', 'mikej', '1992-07-30', 'male', NULL, 'Fitness trainer and nutrition expert', 0, 1, '2023-02-05 11:20:33'),
(4, 'sarah.williams@example.com', 'Sarah', 'Williams', 'sarahw', '1995-03-18', 'female', 'avatars/sarah.jpg', 'Graphic designer and illustrator', 1, 1, '2023-02-15 16:45:12'),
(5, 'david.brown@example.com', 'David', 'Brown', 'davidb', '1985-09-12', 'male', NULL, 'Financial analyst and stock market enthusiast', 1, 1, '2023-03-01 08:10:29'),
(6, NULL, NULL, NULL, 'Tech Enthusiasts', NULL, NULL, 'groups/tech.jpg', 'Group for technology lovers', 1, 0, '2023-01-20 10:00:00'),
(7, NULL, NULL, NULL, 'Photography Club', NULL, NULL, 'groups/photo.jpg', 'Group for photography enthusiasts', 1, 0, '2023-02-10 15:30:00'),
(8, NULL, NULL, NULL, 'Fitness Community', NULL, NULL, 'groups/fitness.jpg', 'Group for fitness and health', 1, 0, '2023-02-28 09:45:00'),
(9, NULL, NULL, NULL, 'Art & Design', NULL, NULL, 'groups/art.jpg', 'Group for artists and designers', 1, 0, '2023-03-05 14:15:00'),
(10, NULL, NULL, NULL, 'Finance Discussion', NULL, NULL, 'groups/finance.jpg', 'Group for financial topics', 1, 0, '2023-03-10 11:20:00');

-- User Authentication Data
INSERT INTO user (id, password_hash) VALUES
(1, '$2a$10$xJwL5v5z5U6Uo6b5Y5X5Xu5X5Xu5X5Xu5X5Xu5X5Xu5X5Xu5X5Xu'), -- johndoe
(2, '$2a$10$yKvL6w6z6V7Vp7c6Z6Y6Yv6Y6Yv6Y6Yv6Y6Yv6Y6Yv6Y6Yv6Y6Yv'), -- janesmith
(3, '$2a$10$zLwM7x7A8B9C0d1E2F3G4H5I6J7K8L9M0N1O2P3Q4R5S6T7U8V9'), -- mikej
(4, '$2a$10$aNbOcPdQeRfTgUhVjWkXlYmZnXlYmZnXlYmZnXlYmZnXlYmZnXlYm'), -- sarahw
(5, '$2a$10$bMcNdPeQfRgShTiUjVkWlYmZnXlYmZnXlYmZnXlYmZnXlYmZnXlYm'); -- davidb

-- Groups Data
INSERT INTO "group" (id, creator_id) VALUES
(6, 1), -- Tech Enthusiasts created by John
(7, 2), -- Photography Club created by Jane
(8, 3), -- Fitness Community created by Mike
(9, 4), -- Art & Design created by Sarah
(10, 5); -- Finance Discussion created by David

-- Follow Relationships
INSERT INTO follow (follower_id, following_id, status) VALUES
(1, 2, 1), -- John follows Jane
(2, 1, 1), -- Jane follows John
(1, 3, 1), -- John follows Mike
(3, 1, 0), -- Mike has a pending follow request to John
(2, 4, 1), -- Jane follows Sarah
(4, 2, 1), -- Sarah follows Jane
(5, 1, 1), -- David follows John
(1, 6, 1), -- John follows Tech Enthusiasts group
(2, 7, 1), -- Jane follows Photography Club
(3, 8, 1), -- Mike follows Fitness Community
(4, 9, 1), -- Sarah follows Art & Design
(5, 10, 1); -- David follows Finance Discussion

-- Posts Data
INSERT INTO posts (id, user_id, group_id, content, image_path, privacy, created_at) VALUES
(1, 1, NULL, 'Just finished a new project using React and Node.js! So excited to share it with everyone.', 'projects/react-node.jpg', 'public', '2023-01-15 12:30:00'),
(2, 2, NULL, 'Beautiful sunset at the beach today. Nature is truly amazing!', 'photos/sunset.jpg', 'public', '2023-01-18 18:45:00'),
(3, 1, 6, 'Anyone interested in a meetup to discuss the latest in web development?', NULL, 'public', '2023-01-22 10:15:00'),
(4, 3, NULL, 'New workout routine is showing great results! Consistency is key.', 'fitness/workout.jpg', 'private', '2023-02-10 07:30:00'),
(5, 4, 9, 'Check out my latest digital painting. What do you think?', 'art/painting.png', 'public', '2023-02-20 14:20:00'),
(6, 5, NULL, 'Market trends are looking interesting this quarter. Time to review portfolios.', NULL, 'public', '2023-03-05 09:00:00'),
(7, 2, 7, 'Upcoming photography workshop next weekend. DM me for details!', 'events/workshop.jpg', 'public', '2023-03-12 16:30:00'),
(8, 3, 8, 'Join our fitness challenge starting next month!', 'fitness/challenge.jpg', 'public', '2023-03-15 11:00:00'),
(9, 4, NULL, 'Working on a new series of illustrations inspired by nature.', 'art/series.jpg', 'public', '2023-03-18 15:45:00'),
(10, 5, 10, 'Interesting article about cryptocurrency trends.', NULL, 'public', '2023-03-20 10:30:00');

-- Comments Data
INSERT INTO comments (id, post_id, user_id, content, image_path, created_at) VALUES
(1, 1, 2, 'Looks amazing John! What tech stack did you use?', NULL, '2023-01-15 13:45:00'),
(2, 1, 3, 'Great work! Would love to see the code if it''s open source.', NULL, '2023-01-15 14:20:00'),
(3, 2, 1, 'Stunning photo Jane! Where was this taken?', NULL, '2023-01-18 19:30:00'),
(4, 3, 5, 'I''d be interested in the web dev meetup. When are you thinking?', NULL, '2023-01-22 11:45:00'),
(5, 5, 2, 'Love the colors in this piece Sarah!', NULL, '2023-02-20 15:10:00'),
(6, 7, 4, 'I might be interested in the workshop. What''s the focus?', NULL, '2023-03-12 17:20:00'),
(7, 8, 1, 'Count me in for the fitness challenge!', NULL, '2023-03-15 12:30:00'),
(8, 9, 2, 'These illustrations are gorgeous! Do you take commissions?', NULL, '2023-03-18 16:20:00'),
(9, 10, 1, 'Very insightful article. Thanks for sharing!', NULL, '2023-03-20 11:45:00'),
(10, 4, 2, 'Looking fit Mike! What''s your routine?', NULL, '2023-02-10 09:15:00');

-- Messages Data (only between users who follow each other)
INSERT INTO message (sender_id, receiver_id, isread, content, created_at) VALUES
(1, 2, 1, 'Hey Jane, how are you doing?', '2023-01-14 08:30:00'),
(2, 1, 1, 'I''m good John! Just working on some new photos.', '2023-01-14 08:35:00'),
(1, 2, 1, 'That sunset pic was amazing!', '2023-01-18 20:00:00'),
(2, 4, 1, 'Sarah, would you be interested in collaborating on a project?', '2023-03-15 10:20:00'),
(4, 2, 1, 'Sounds interesting! Tell me more.', '2023-03-15 10:25:00'),
(1, 5, 0, 'David, I have a question about investments.', '2023-03-22 14:15:00'),
(4, 2, 1, 'I was thinking about a combined art and photography exhibit.', '2023-03-16 11:30:00'),
(2, 4, 1, 'That sounds like a great idea! Let''s discuss details.', '2023-03-16 11:35:00');

-- Requests Data
INSERT INTO request (sender_id, receiver_id, target_id, type, created_at) VALUES
(3, 1, NULL, 0, '2023-02-01 09:15:00'), -- Mike sent follow request to John
(5, 6, 6, 1, '2023-01-25 11:30:00'), -- David requested to join Tech Enthusiasts
(4, 7, 7, 1, '2023-02-15 14:45:00'), -- Sarah requested to join Photography Club
(2, 8, 8, 1, '2023-03-01 16:20:00'), -- Jane requested to join Fitness Community
(1, 9, 9, 1, '2023-03-10 10:30:00'); -- John requested to join Art & Design

-- Sessions Data
INSERT INTO sessions (user_id, session_id, refresh_token, expires_at, ip_address, user_agent) VALUES
(1, 'a1b2c3d4-e5f6-g7h8-i9j0-k1l2m3n4o5p6', 'q1w2e3r4t5y6u7i8o9p0', '2023-04-10 08:30:00', '192.168.1.100', 'Mozilla/5.0 (Windows NT 10.0)'),
(2, 'b2c3d4e5-f6g7-h8i9-j0k1-l2m3n4o5p6q7', 'r3t4y5u6i7o8p9a0s1d2', '2023-04-12 12:45:00', '192.168.1.101', 'Mozilla/5.0 (Macintosh; Intel Mac OS X)'),
(3, 'c3d4e5f6-g7h8-i9j0-k1l2-m3n4o5p6q7r8', 's4d5f6g7h8j9k0l1z2x3', '2023-04-15 18:20:00', '192.168.1.102', 'Mozilla/5.0 (iPhone; CPU iPhone OS)'),
(4, 'd4e5f6g7-h8i9-j0k1-l2m3-n4o5p6q7r8s9', 't5u6v7w8x9y0z1a2b3c4', '2023-04-18 14:10:00', '192.168.1.103', 'Mozilla/5.0 (Android 12)'),
(5, 'e5f6g7h8-i9j0-k1l2-m3n4-o5p6q7r8s9t0', 'u6v7w8x9y0z1a2b3c4d5', '2023-04-20 09:45:00', '192.168.1.104', 'Mozilla/5.0 (Linux x86_64)');

-- Likes Data
INSERT INTO likes (user_id, post_id, comment_id, created_at) VALUES
(2, 1, NULL, '2023-01-15 13:40:00'), -- Jane liked John's post
(3, 1, NULL, '2023-01-15 14:15:00'), -- Mike liked John's post
(1, 2, NULL, '2023-01-18 19:00:00'), -- John liked Jane's post
(4, 2, NULL, '2023-01-18 19:30:00'), -- Sarah liked Jane's post
(5, 3, NULL, '2023-01-22 11:40:00'), -- David liked John's group post
(1, NULL, 5, '2023-02-20 15:30:00'), -- John liked Sarah's comment
(2, 5, NULL, '2023-02-20 15:00:00'), -- Jane liked Sarah's art post
(3, 8, NULL, '2023-03-15 12:00:00'), -- Mike liked fitness challenge post
(4, NULL, 7, '2023-03-15 13:20:00'), -- Sarah liked comment on fitness post
(5, 10, NULL, '2023-03-20 11:00:00'); -- David liked finance post

-- Events Data
INSERT INTO events (id, user_id, group_id, description, title, timeof, created_at) VALUES
(1, 1, NULL, 'Web development workshop covering React and Node.js', 'Web Dev Workshop', '2023-04-15 14:00:00', '2023-03-20 10:00:00'),
(2, 2, 7, 'Photography field trip to the national park', 'Photo Field Trip', '2023-04-22 09:00:00', '2023-03-18 15:30:00'),
(3, 3, 8, 'Group fitness session at the central park', 'Fitness Bootcamp', '2023-04-08 07:00:00', '2023-03-15 08:45:00'),
(4, 4, NULL, 'Digital art exhibition at the downtown gallery', 'Art Exhibition', '2023-05-10 18:00:00', '2023-03-25 14:20:00'),
(5, 5, 10, 'Financial planning seminar with guest speakers', 'Finance Seminar', '2023-05-05 13:00:00', '2023-03-22 11:15:00');

-- User Events (RSVPs)
INSERT INTO userevents (event_id, user_id, respond, created_at) VALUES
(1, 2, 1, '2023-03-21 11:20:00'), -- Jane attending Web Dev Workshop
(1, 5, 1, '2023-03-22 09:15:00'), -- David attending Web Dev Workshop
(2, 1, 1, '2023-03-19 16:30:00'), -- John attending Photo Field Trip
(2, 4, 1, '2023-03-20 14:00:00'), -- Sarah attending Photo Field Trip
(3, 1, 1, '2023-03-16 07:30:00'), -- John attending Fitness Bootcamp
(3, 2, 0, '2023-03-17 12:45:00'), -- Jane declined Fitness Bootcamp
(4, 2, 1, '2023-03-26 10:15:00'), -- Jane attending Art Exhibition
(4, 3, 1, '2023-03-27 09:30:00'), -- Mike attending Art Exhibition
(5, 1, 1, '2023-03-23 15:20:00'), -- John attending Finance Seminar
(5, 4, 0, '2023-03-24 11:40:00'); -- Sarah declined Finance Seminar