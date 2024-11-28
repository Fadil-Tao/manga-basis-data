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


delimiter $$ 
create function author_exist(
    in author_id int
)
returns tinyint DETERMINISTIC
begin 
    declare is_exist tinyint;

    select count(id) into is_exist from author where author.id = author_id;
    return is_exist;
end $$
delimiter ;



-- check if user is exist and active 
delimiter $$ 
create function is_active(user_id int)
returns tinyint DETERMINISTIC
begin 
    declare is_active tinyint;
    
    select count(id) into is_active from user 
    where user.id = user_id and user.is_active = 1;
    return is_active;
end$$
delimiter ; 