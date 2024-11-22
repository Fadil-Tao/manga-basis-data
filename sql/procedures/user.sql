-- 1. register normal user 
delimiter $$ 
create or replace procedure register_user(
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
	insert into user(username,email,`password`)values
	(n_username,n_email,SHA2(N_PASSWORD,256));	
	commit;
end$$

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
    SELECT `password`, is_active
    INTO user_password, user_status
    FROM user
    WHERE email = n_email;
	
    IF user_password IS NULL THEN
      	signal sqlstate '45000' set message_text = "Email Did not exist";
      	rollback;
    ELSEIF user_password != SHA2(n_password, 256) THEN
        signal sqlstate '45000' set message_text = "Wrong Credential";
        rollback;
    ELSEIF user_status < 1 THEN 
    	signal sqlstate '45000' set message_text ="User Inactive";
    	rollback;
    ELSE
       	SELECT user.id ,user.username, user.email,user.is_admin, user.created_at from user where user.email = n_email; 
       	commit;
    END IF;
END$$
DELIMITER ;
drop procedure login_user;


--  change user into an admin
delimiter $$
create procedure