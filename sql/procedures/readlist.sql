-- 1. add new readlist
delimiter $$
create procedure add_readlist(
    in n_owned_by int,
    in n_name varchar(255),
    in n_description text
)
begin
    declare exit handler for sqlexception
    begin
        rollback;
        resignal;
    end;
    start transaction;

    if (is_active(n_owned_by)) = 0 then
        signal sqlstate '45000' set message_text = 'Unauthorized';
        rollback;	
    else 
        insert into readlist(user_id,name,description)
        values(n_owned_by,n_name,n_description);
        commit;
    end if;
end$$
delimiter ;


-- 2. delete readlist
delimiter $$ 
create procedure delete_readlist(
	in n_readlist_id int,
    in user_id int
)
begin
	declare exit handler for sqlexception 
	begin
		rollback;
		resignal;
	end;

	start transaction;
    if (is_active(user_id)) = 0 then 
        signal sqlstate '45000' set message_text = 'Unauthorized';
        rollback;
    end if;

    if (select count(id) from readlist where readlist.id = n_readlist_id and readlist.user_id = user_id) < 1 then
        signal sqlstate '45000' set message_text = 'Unauthorized';
        rollback;
    end if;
	delete from readlist where readlist.id = n_readlist_id;
	commit;
end$$
delimiter ;


-- 3. update readlist
delimiter $$
create procedure update_readlist(
	in n_name varchar(255),
	in n_description text,
	in n_readlist_id int,
	in user_id int
)
begin
	declare exit handler for sqlexception
	begin
		rollback;
		resignal;
	end;
	
	start transaction;
        if (is_active(user_id)) = 0 then 
        signal sqlstate '45000' set message_text = 'Unauthorized';
        rollback;
    end if;

    if (select count(id) from readlist where readlist.id = n_readlist_id and readlist.user_id = user_id) < 1 then
        signal sqlstate '45000' set message_text = 'Unauthorized';
        rollback;
    end if;
	update readlist set name = n_name, description = n_description
	where readlist.id = n_readlist_id;
	commit;
end$$
delimiter ;

-- 4. search readlist
delimiter $$
create procedure search_readlist(
    in input varchar (255)
)
begin 
    declare exit handler for sqlexception
    begin 
        rollback;
        resignal;
    end;

    start transaction;
    if input = '' then 
    	 select readlist.id, user.username as 'Owner', readlist.name, readlist.description, readlist.created_at, readlist.updated_at 
    	from readlist join user on readlist.user_id = user.id ;
    	commit;
    else
    	select readlist.id, user.username as 'Owner', readlist.name, readlist.description, readlist.created_at, readlist.updated_at 
    	from readlist join user on readlist.user_id = user.id  where readlist.name like concat(input, '%');
    	commit;
   	end if;
end$$
delimiter ;


-- 5. Get user's readlist
create procedure get_readlists_from_user(
    in user_id int
)
begin 
    declare exit handler for sqlexception
    begin 
        rollback;
        resignal;
    end;

    start transaction;
    select r.id, user.username  as 'Owner', r.name, r.description, r.created_at, r.updated_at 
    from readlist r join user on r.user_id = user.id where r.user_id = user_id order by r.updated_at desc;
    commit;
end$$
delimiter ;


-- 6. add manga into readlist
delimiter $$
create procedure add_to_readlist(
	in n_user_id int,
	in n_manga_id int,
	in n_read_status  enum ('done', 'in_progress', 'later'),
	in n_readlist_id int
)
begin 
	declare exit handler for sqlexception
	begin 
		rollback;
		resignal;
	end;
	start transaction ; 
	if (is_active(n_user_id)) = 0 then 
		signal sqlstate '45000'
		set MESSAGE_TEXT = "Unauthorized";
        rollback;
    end if;
	if (select count(*) from readlist where 
	readlist.user_id = n_user_id and readlist.id = n_readlist_id) = 0 then
		signal sqlstate '45000'
		set MESSAGE_TEXT = "Unauthorized";
		rollback;
	end if;
	if (select count(readlist_item.manga_id) from readlist_item where readlist_item.manga_id = n_manga_id) > 0 then
		signal sqlstate '45000'
		set MESSAGE_TEXT = "conflict : manga already in readlist";
		rollback;
	end if;
		insert into readlist_item ( readlist_id ,read_status, manga_id  ) 
		values (n_readlist_id,n_read_status,n_manga_id);
		commit;
end$$
delimiter ;  

-- 8. delete item from readlist
delimiter $$
create procedure delete_readlist_item(
	in n_readlist_item_id int,
    in n_user_id int
)
begin
	declare exit handler for sqlexception
	begin
		rollback;
		resignal;
	end;
	
	start transaction;
    if (is_active(n_user_id)) = 0 then 
		signal sqlstate '45000'
		set MESSAGE_TEXT = "Unauthorized";
        rollback;
    end if;
	    delete from readlist_item 
	    where readlist_item.id = n_readlist_item_id;
	commit;
end$$
delimiter ;

-- 9. get readlist manga list 
delimiter $$
create procedure get_manga_list_from_readlist(
    in readlist_id int    
)
begin 
    declare exit handler for sqlexception
    	begin
		rollback;
		resignal;
	end;
	
	start transaction;
    select ri.id, ri.manga_id, m.title, ri.read_status, ri.added_at from readlist_item ri 
    join manga m on ri.manga_id = m.id where ri.readlist_id = readlist_id;
    commit ;
end $$
delimiter ;

-- 10. update readlist item read status
delimiter $$ 
create procedure update_readlist_item_status(
    in readlist_item_id int,
    in n_read_status enum('done', 'in_progress', 'later'),
    in user_id int
)
begin 
    declare exit handler for sqlexception
    begin
		rollback;
		resignal;
	end;
	
	start transaction;
    if (is_active(user_id)) = 0 then 
        signal sqlstate '45000'
		set MESSAGE_TEXT = "Unauthorized";
        rollback;
    end if;
    if (select count(id) from readlist_item where id = readlist_item_id) < 1 then
		signal sqlstate '45000'
		set MESSAGE_TEXT = "Unauthorized";
        rollback;
    end if;
    if (n_read_status not in ('done', 'in_progress', 'later')) then 
        signal sqlstate '45000'
		set MESSAGE_TEXT = "malformed : incorrect status";
        rollback;
    end if;
        update readlist_item set read_status = n_read_status 
        where readlist_item.id = readlist_item_id;
    commit;
end$$
delimiter ; 