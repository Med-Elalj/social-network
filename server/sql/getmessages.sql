
-- uID: user ID  , rID: receiver ID , offsetTime: time offset
WITH follow_exists AS (
    SELECT 1 as following FROM follow f
    WHERE
        f.follower_id = 2 -- uID
        AND f.following_id = 3 -- rID
        AND f.status = 1
    LIMIT 1
), group_exists AS (
    SELECT 1 as exits FROM "group" g WHERE g.id = 3 -- rID
    LIMIT 1
)
SELECT
    m.sender_id,
    p.display_name AS sender_display_name,
    m.receiver_id,
    m.content,
    m.created_at
FROM
    message m
    JOIN profile p ON p.id = m.sender_id
    LEFT JOIN group_exists g ON 1=1
    LEFT JOIN follow_exists f ON 1=1
WHERE
    m.created_at < '2023-10-01 12:25:01' -- offsetTime
    AND (
        (
            g.exits = 1
            AND (m.receiver_id = 3) -- rID
            AND f.following = 1
        )
        OR (
            g.exits IS NULL
            AND 2 IN (m.sender_id, m.receiver_id) -- uID
            AND 3 IN (m.sender_id, m.receiver_id) -- rID
        )
    )
ORDER BY
    m.created_at ASC;