-- 1. register normal user 
drop procedure register_user;
delimiter $$ 
create procedure register_user(
	in n_username varchar(255),
	in n_email varchar(255),
	in n_password text
)
begin
	declare exit handler for sqlexception
	begin
		rollback;
		resignal;
	end;
	
	start transaction;
	
	if (select count(username) from user where user.username = n_username) > 0 then
		signal sqlstate '45000' set message_text = "username already used";
		rollback;
	end if ;
	if (select count(email) from user where user.email = n_email) > 0 then
		signal sqlstate '45000' set message_text = "email already used";
		rollback;
	end if ;
	if (LENGTH(n_email)) < 3 then
		signal sqlstate '45000' set message_text = "invalid email";
		rollback;
	end if ;
	if (LENGTH(n_password)) < 8 then
		signal sqlstate '45000' set message_text = "invalid password minimum 8";
		rollback;
	end if ;
	if (LENGTH(n_username)) < 5 then
		signal sqlstate '45000' set message_text = "invalid username minimum character is 5";
		rollback;
	end if ;
	if (LENGTH(n_username)) > 10 then
		signal sqlstate '45000' set message_text = "invalid username max character is 10";
		rollback;
	end if ;
	if INSTR(n_username,' ') > 0 then
		signal sqlstate '45000' set message_text = "invalid password cant contain whitespace";
		rollback;
	end if;
	if INSTR(n_email,' ') > 0 then
		signal sqlstate '45000' set message_text = "invalid email cant contain whitespace";
		rollback;
	end if;
	if INSTR(n_password,' ') > 0 then
		signal sqlstate '45000' set message_text = "invalid password cant contain whitespace";
		rollback;
	end if;
	
	insert into user(username,email,`password`)values
	(n_username,n_email,SHA2(N_PASSWORD,256));	
	commit;
end$$
delimiter ;


-- 2. login
DELIMITER $$ 
CREATE PROCEDURE login_user(
   IN n_email VARCHAR(100),
   IN n_password TEXT
) 
BEGIN
    DECLARE user_password TEXT;
    DECLARE user_status INT;
   
   	start transaction;
    SELECT `password`
    INTO user_password
    FROM user
    WHERE email = n_email;
	
    IF user_password IS NULL THEN
      	signal sqlstate '45000' set message_text = "Invalid Credential";
      	rollback;
    ELSEIF user_password != SHA2(n_password, 256) THEN
        signal sqlstate '45000' set message_text = "Invalid Credential";
        rollback;
    ELSE
       	SELECT user.id ,user.username, user.email,user.is_admin, user.created_at from user where user.email = n_email; 
       	commit;
    END IF;
END$$
DELIMITER ;


-- 3. search user
delimiter $$
create procedure get_all_user_with_search_by_username(
	in in_username varchar(255)
)
begin
	declare exit handler for sqlexception
	begin 
		rollback;
		resignal;
	end;

	start transaction;

	if in_username = '' then
		select user.id, user.username, user.email , user.created_at  from user;
	else 
		select user.id, user.username, user.email, user.created_at  from user
		where user.username like concat(in_username, '%');
		commit;
	end if; 
end$$
delimiter ;




-- 5 update user 
drop procedure update_user;
delimiter $$
create procedure update_user(
	in user_id int,
	in username_target varchar(255),
	in n_username varchar(255),
	in n_password text,
	in n_email varchar(255)
)
begin 
	declare target_id int;
	declare exit handler for sqlexception
	begin 
		rollback;
		resignal;
	end;
	
	start transaction;

	if (select user.id from user where user.username = username_target) != user_id then 
		signal sqlstate '45000' set message_text = "Unauthorized";
		rollback;
	else select user.id into target_id from user where user.username = username_target;
	end if ;
	if n_username != '' then 
		if (LENGTH(n_username)) < 5 or (LENGTH(n_username)) > 10 then
			signal sqlstate '45000' set message_text = "malformed : username char minimum 8 and maximum 10";
			rollback;
		end if ;
		if (select count(username) from user where user.username = n_username) > 0 then
			signal sqlstate '45000' set message_text = "conflict : username already used";
			rollback;
		end if ;
		update user set user.username = n_username where user.id = target_id ;
	end if;
	if n_password != '' then 
		if (LENGTH(n_password)) < 8 then
			signal sqlstate '45000' set message_text = "malformed : password minimum 8";
			rollback;
		end if ;
		update user set user.password = SHA2(n_password,256) where user.id = target_id;
	end if;
	if n_email != '' then 
		if (LENGTH(n_email)) < 3 then
			signal sqlstate '45000' set message_text = "malformed : invalid email";
			rollback;
		end if ;
		if (select count(email) from user where user.email = n_email) > 0 then
			signal sqlstate '45000' set message_text = "conflict : email already used";
			rollback;
		end if ;
		update user set user.email = n_email where user.id = target_id;
	end if ;
	commit;
end$$
delimiter ; 

-- 6 get user detail by username
drop procedure get_user_detail_by_username;
delimiter $$
create procedure get_user_detail_by_username(
	in username varchar(255)
)
begin 
	start transaction;
	select user.id , user.username, user.email , user.created_at from user where
	user.username = username;
	commit;
end$$ 
delimiter ;

-- 7. get manga that user rate lately
delimiter $$
create procedure get_user_rated_manga(
	in username varchar(255)
)
begin 
	start transaction;
	SELECT 
        m.id AS id, 
        m.title AS title, 
        m.manga_status AS status, 
        m.published_at, 
        m.finished_at, 
		MAX(r2.created_at) AS latest_rating_date,
        COALESCE(AVG(r2.rating), NULL) AS average_rating,
        COUNT(DISTINCT r2.user_id) AS total_user_rated,
        MAX(CASE WHEN r2.user_id = user_id THEN r2.rating END) AS your_rating
        FROM 
        	manga m
        LEFT JOIN 
            rating r2 ON m.id = r2.manga_id	
		left join user on r2.user_id = user.id 
		where user.username = username
		group by m.id, m.title, m.manga_status, m.published_at, m.finished_at
		order by latest_rating_date desc;
	commit;
end$$
delimiter ;


-- 8. get manga that user like
drop procedure get_user_liked_manga;
delimiter $$
create procedure get_user_liked_manga(
	in username varchar(255)
)
begin 
	start transaction;
	SELECT 
        m.id AS id, 
        m.title AS title, 
        m.manga_status AS status, 
        m.published_at, 
        m.finished_at, 
		MAX(lm.created_at) AS latest_liked_date,
    	COUNT(DISTINCT lm.user_id) AS total_likes
        FROM 
        	manga m
        LEFT JOIN 
            liked_manga lm ON m.id = lm.manga_id	
		left join user on lm.user_id = user.id 
		where user.username = username
		group by m.id, m.title, m.manga_status, m.published_at, m.finished_at
		order by latest_liked_date desc;
	commit;
end$$
delimiter ;


-- 9 delete user

drop procedure delete_user;
delimiter $$
create procedure delete_user(
	in username varchar(255),
	in user_id int
)
begin
	
	declare exit handler for sqlexception
	begin 
		rollback;
		resignal;
	end;
	
	if (select is_admin(user_id)) < 1 then
		signal sqlstate '45000' set message_text = "Unauthorized";
		rollback;
	end if; 
	if (select is_admin((select user.id from user where user.username = username))) < 1 then
		signal sqlstate '45000' set message_text = "admin user can't be deleted";
		rollback;
	end if;
	if (select count(id) from user where user.username = username) < 1 then
		signal sqlstate '45000' set message_text = "User not found";
		rollback;
	end if;

	delete from user where user.username = username;
	commit;
end$$
delimiter ;

-- 10. get readlist of a user
drop procedure get_readlist_from_user;
delimiter $$
create procedure get_readlist_from_user(
	in username varchar(255)
)
begin 
	declare user_id int;
	declare exit handler for sqlexception
	begin 
		rollback;
		resignal;
	end;
	start transaction;
	select user.id into user_id where user.username = username;
	if user_id = null then
		signal sqlstate '45000' set message_text = "user not found";
		rollback;
	end if;
	select readlist.id, readlist.name, readlist.description, readlist.created_at, readlist.updated_at
	from readlist where readlist.user_id = user_id ;
	commit;
end$$
delimiter ;
--