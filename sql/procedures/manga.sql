use manga_basis_data;

-- 1. add manga 
delimiter $$
create procedure add_manga(
	in n_title varchar(255),
	in n_synopsys text,
	in n_manga_status enum('finished','in_progress'),
	in n_published_at date,
	in n_finished_at date,
	in user_id int
)
begin 
	declare exit handler for sqlexception 
	begin 
		rollback;
		resignal;
	end;
	
	start transaction;
	if (select is_admin(user_id)) = 1 then 
		insert into manga (title,synopsys,manga_status,published_at,finished_at)
		values (n_title,n_synopsys,n_manga_status,n_published_at,n_finished_at);
		commit;
	else
		signal sqlstate '45000' set message_text = "Unauthorized";
		rollback;
	end if;
end$$
delimiter ;

-- 2. connect author and manga
delimiter $$
create procedure connect_author_manga(
	in n_id_manga int,
	in n_id_author int,
	in user_id int
)
begin
	declare exit handler for sqlexception
	begin 
		rollback;
		resignal;
	end;

	start transaction; 
	if (select is_admin(user_id)) = 1 then 
		insert into author_manga_pivot(id_manga, id_author)
		values (n_id_manga, n_id_author);
		commit;
	else 
		signal sqlstate '45000' set message_text = "Unauthorized";
		rollback;
	end if;
end$$
delimiter ;


-- 3. connect genre and manga
delimiter $$
create procedure connect_genre_manga(
	in n_id_manga int,
	in n_id_genre int,
	in user_id int
)
begin 
	declare exit handler for sqlexception
	begin 
		rollback;
		resignal;
	end;
	start transaction; 
	if (select is_admin(user_id)) = 1 then 
		insert into manga_genre_pivot(id_manga, id_genre)
		values (n_id_manga, n_id_genre); 
		commit;
	else
		signal sqlstate '45000' set message_text = "Unauthorized";
		rollback;
	end if;
end$$
delimiter ;


-- 4. get manga by id
delimiter $$
create procedure get_manga_detail(
	in in_manga_id int
) 
begin
	SELECT m.id,m.title, m.synopsys, m.manga_status, 
    m.published_at, 
    m.finished_at 
	FROM manga m
	WHERE m.id = in_manga_id;
end$$
delimiter ; 
-- 5. get manga's author
delimiter $$
create procedure get_manga_genre(
	in in_manga_id int
) 
begin
	select g.id,g.name from genre g 
	inner join manga_genre_pivot mgp ON g.id = mgp.id_genre  
	inner join manga m on mgp.id_manga = m.id where m.id = in_manga_id;
end$$
delimiter ; 
-- 6. get manga's genre
delimiter $$
create procedure get_manga_author(
	in in_manga_id int
) 
begin
	select a.id,a.name from author a 
	inner join author_manga_pivot amp ON a.id = amp.id_author 
	inner join manga m on amp.id_manga = m.id where m.id = in_manga_id;
end$$
delimiter ; 

-- 7. search by its name
delimiter $$ 
create procedure get_manga_by_title(
	in input varchar(255)
)
begin 
	declare exit handler for sqlexception
		begin 
			rollback;
			resignal;
		end;
	start transaction;
		select m.id, m.title, m.synopsys, m.manga_status, m.published_at , m.finished_at from manga m
		where m.title like CONCAT(input, '%') ;
	commit;
end$$
delimiter ;


-- 8. get sort