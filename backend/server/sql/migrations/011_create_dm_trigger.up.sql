CREATE TRIGGER IF NOT EXISTS insert_dm BEFORE INSERT ON message
FOR EACH ROW
BEGIN
    SELECT RAISE(
        ABORT,
        'You can not send a message to a user you are not following'
    )
    WHERE
        NEW.sender_id <> NEW.receiver_id
        AND NOT EXISTS (
            SELECT
                1
            FROM
                follow
            WHERE
                (follow.follower_id = NEW.sender_id AND follow.following_id = NEW.receiver_id)
                OR
                (follow.follower_id = NEW.receiver_id AND follow.following_id = NEW.sender_id)
            LIMIT 1
        );
END;