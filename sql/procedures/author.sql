-- 1. create author
delimiter $$
create procedure add_author(
	n_name varchar(255),
	n_birthday date,
	n_biography text,
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
		insert into author(name, birthday, biography)
		values (n_name,n_birthday,n_biography);
		commit;
	else
		signal sqlstate '45000' set message_text = "Unauthorized";
		rollback;
	end if;
end$$
delimiter ;

-- 2. update author
DELIMITER $$
CREATE PROCEDURE update_author(
    IN p_author_id INT,
    IN p_name VARCHAR(255),
    IN p_birthday DATE,
    IN p_biography TEXT,
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
    	UPDATE author
    	SET name = p_name,
        	birthday = p_birthday,
        	biography = p_biography
    	WHERE author.id = p_author_id;
	else 
		signal sqlstate '45000' set message_text = "Unauthorized";
		rollback;
	end if;
END$$
DELIMITER ;

-- 3. delete author
DELIMITER $$
CREATE PROCEDURE delete_author(IN p_author_id INT , in user_id int)
BEGIN
	declare exit handler for sqlexception
		begin 
			rollback;
			resignal;
		end;
	start transaction;
	if (select is_admin(user_id)) = 1 then 	
    	DELETE FROM author WHERE author.id = p_author_id;
	else
		signal sqlstate '45000' set message_text = "Unauthorized";
		rollback;
	end if;
END$$
DELIMITER ;

-- 4. get all author
delimiter $$ 
create procedure get_all_author()
begin
	select a.id , a.name , a.birthday  from author a; 
end$$
delimiter ;


-- 5. get author by id
delimiter $$
create procedure get_author_by_id(
	author_id int
)
begin 
	declare exit handler for sqlexception
	begin 
		rollback;
		resignal;
	end;
	start transaction;
	select a.id, a.name , a.birthday, a.biography  from author a where a.id  = author_id;
	commit;
end$$ 
delimiter ;

-- 6. search author by its name
delimiter $$ 
create procedure get_author_by_name(
	input varchar(255)
)
begin
	declare exit handler for sqlexception
	begin 
		rollback;
		resignal;
	end;
	
	start transaction;
		select a.id, a.name, a.birthday, a.biography from author a where a.name like CONCAT(INPUT,'%') ;
	commit;
end $$
delimiter;

-- 7. get author's manga  
delimiter $$
create procedure get_author_manga(in author_id int)
begin 
	declare exit handler for sqlexception
	begin 
		rollback;
		resignal;
	end;
	
	start transaction;
		select m.id,m.title, m.manga_status,m.synopsys, m.published_at from manga m
		inner join author_manga_pivot amp on m.id = amp.id_manga where
		amp.id_author = author_id order by m.published_at desc;
	commit;
end $$
delimiter ;