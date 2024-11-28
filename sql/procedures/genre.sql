-- 1. add genre
delimiter $$
create procedure add_genre(
	n_name varchar(255),
	n_description text,
	user_id int
)
begin 
	declare exit handler for sqlexception
		begin 
			rollback;
			resignal;
		end;
	start transaction;
	if (select is_admin(user_id)) = 1 then 
		insert into genre(name, description)
		values (n_name,n_description);
		commit;
	else
		signal sqlstate '45000' set message_text = "Unauthorized";
		rollback;
	end if ;
end$$
delimiter ;


-- 2. update genre
drop procedure update_genre;
DELIMITER $$
CREATE PROCEDURE update_genre(
    IN p_genre_id INT,
    IN p_name VARCHAR(255),
    IN p_description TEXT,
    IN user_id int
)
BEGIN
	declare exit handler for sqlexception
	begin 
		rollback;
		resignal;
	end;
	start transaction;
	if (select is_admin(user_id)) = 1 then 
		if (select count(id) from genre where genre.id = p_genre_id)>0 then
    		UPDATE genre
    			SET name = p_name,
        		description = p_description
    		WHERE genre.id = p_genre_id;
    		commit;
    	else 
    		signal sqlstate '45000' set message_text = "Genre not found";
    		rollback;
    	end if;
	else
		signal sqlstate '45000' set message_text = "Unauthorized";
		rollback;
	end if ;
END$$
DELIMITER ;


-- 3. delete genre
DELIMITER $$
CREATE PROCEDURE delete_genre(IN p_genre_id INT, user_id int)
BEGIN
	declare exit handler for sqlexception
	begin 
		rollback;
		resignal;
	end;
	start transaction;
	if (select is_admin(user_id)) = 1 then 
		if (select count(id) from genre where genre.id = p_genre_id) >0 then
			DELETE FROM genre WHERE genre.id = p_genre_id;
			commit;
		else
			signal sqlstate '45000' set message_text = "Genre not found";
			rollback;
		end if; 
	else
		signal sqlstate '45000' set message_text = "Unauthorized";
		rollback;
	end if ;
END$$
DELIMITER ;

-- 4. get all genre
delimiter $$
create procedure get_all_genre() 
begin
	select g.id , g.name , g.description  from genre g;
end$$
delimiter ; 


-- 5. search genre by its name
delimiter $$
create procedure get_genre_by_name(in input varchar(255))
begin
    declare exit handler for sqlexception
    begin 
        rollback;
        resignal;
    end;
    start transaction;
        select g.id, g.name, g.description from genre g where g.name like concat(input,'%');
    commit;
end$$
delimiter;