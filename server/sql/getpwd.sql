-- database: ../db/test.db
-- Use the ▷ button in the top right corner to run the entire file.
SELECT
    password
FROM
    users u
WHERE
    ? IN (username, email);
