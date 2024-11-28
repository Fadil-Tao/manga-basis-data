-- 1. add rating
delimiter $$ 
create procedure rate_manga(
	in n_manga_id int,
	in n_user_id int,
	in n_rating smallint
)
begin
	declare user_active tinyint;
	declare exit handler for sqlexception
	begin
		rollback;
		resignal;
	end;
    
    start transaction;
    if n_rating > 10 or n_rating < 0 then 
        signal sqlstate '45000' set message_text = "malformed request : rating can't be lower than 0 or higher than 10 ";
		rollback;
    end if;
    if (select count(id) from manga where manga.id = n_manga_id) = 0 then
    	signal sqlstate '45000' set message_text = "Manga not found";
		rollback;
    end if;
	if (is_active(n_user_id)) = 0 then
        signal sqlstate  '45000' set message_text = "Unauthorized";
        rollback;
    end if;
    if (select count(id) from rating where rating.manga_id = n_manga_id and rating.user_id = n_user_id) > 0 then
        update rating set rating.rating = n_rating where rating.user_id = n_user_id and rating.manga_id = n_manga_id;
	else
        insert into rating(manga_id, user_id,rating)
        values (n_manga_id,n_user_id, n_rating);
        commit;        
	end if;
end$$
delimiter ;

-- get rating of a manga

