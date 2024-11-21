DELIMITER $$
CREATE FUNCTION is_admin(user_id int)
RETURNS tinyint DETERMINISTIC
BEGIN
    DECLARE is_admin tinyint;

    SELECT user.is_admin into is_admin from user where
    user.id = user_id;
    RETURN is_admin;
end $$
DELIMITER ;