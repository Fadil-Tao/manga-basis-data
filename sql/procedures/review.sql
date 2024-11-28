-- 1. Buat sebuat review
drop procedure add_review;
delimiter $$
create procedure add_review(
	in n_manga_id int,
	in n_user_id int,
	in n_review_text text,
	in n_tag enum('Reccomended' , 'Mixed Feelings' ,'Not Reccomended')
)
begin
	declare user_active tinyint;
	declare exit handler for sqlexception
	begin
		rollback;
		resignal;
	end;	
		    
	start transaction;
	select is_active into user_active from user where user.id = n_user_id;

	if user_active = 0 then 
		SIGNAL SQLSTATE '45000'
        SET MESSAGE_TEXT = 'Unauthorized';
	else
		if (select count(id) from manga where manga.id = n_manga_id) > 0 then
			if (select count(id) from review where review.user_id = n_user_id and review.manga_id = n_manga_id) > 0 then 
				signal sqlstate '45000' 
				set message_text = 'conflict : Youre already submitted a review for this manga';
				rollback;
			else
   				insert into review(manga_id, user_id, review_text,tag)
				values (n_manga_id,n_user_id,n_review_text,n_tag);
				commit;
			end if;
		else 
			signal sqlstate '45000'
			set message_text = 'manga not found';
		end if;
	end if ;
end$$
delimiter ;


-- 2. Get list review of a manga
DELIMITER $$
CREATE PROCEDURE get_review_from_manga(
    IN manga_id INT
)
BEGIN
	declare exit handler for sqlexception
	begin
		rollback;
		resignal;
	end;

	start transaction;
	if (select count(id) from manga where manga.id = manga_id) > 0 then 	
		SELECT 
			review.id, 
			review.user_id, 
			user.username, 
			review.review_text,
			review.tag, 
			review.created_at,
			COUNT(liked_review.id) AS total_likes
		FROM 
			review
		JOIN 
			user ON review.user_id = user.id
		LEFT JOIN 
			liked_review ON review.id = liked_review.review_id 
		WHERE 
			review.manga_id = manga_id
		GROUP BY 
			review.id, review.user_id, user.username, review.review_text, 
			review.tag, review.created_at
		ORDER BY 
			review.created_at DESC;
		commit;
	else
		signal sqlstate '45000'
		set message_text = 'manga not found';
		rollback;
	end if;
END$$
DELIMITER ;



-- 3. delete a review
delimiter $$
create procedure delete_a_review(
	in n_user_id int,
	in n_review_id int
)
begin 
	declare an_admin tinyint;
	declare is_owner tinyint;
	declare exit handler for sqlexception
	begin
		rollback;
		resignal;
	end;
	
	select is_admin into an_admin from user
	where user.id = n_user_id; 

	select count(*) into is_owner from review 
	where review.id = n_user_id and review.id = n_review_id ;
		if an_admin = 0  and  is_owner = 0 then
			SIGNAL SQLSTATE '45000'
        	SET MESSAGE_TEXT = 'Unauthorized';
			rollback;
  		else 
		 	start transaction;
		 	delete from review where review.id = n_review_id;
			commit;
		 end if;
end$$
delimiter ;

-- 4 update a review
delimiter $$
create procedure update_a_review(
	in user_id int,
	in review_id int,
	in review_text text,
	in n_tag enum('Reccomended' , 'Mixed Feelings' ,'Not Reccomended')
)
begin
	declare exit handler for sqlexception
		begin 
			rollback;
			resignal;
		end;
	start transaction;
	if (select count(*) from review where review.user_id = user_id 
	and review.id  = review_id ) > 0 then
		update review set review.review_text = review_text, review.tag = n_tag
		where review.user_id = user_id and review.id = review_id;
		commit;
	else
		signal sqlstate '45000' set message_text = 'Unauthorized';
		rollback;
	end if;
end$$
delimiter ;


-- 5 Get one review by its id 
delimiter $$
create procedure get_review_by_id(
	in review_id int 
)
begin
	declare exit handler for sqlexception
		begin 
			rollback;
			resignal;
		end;
	start transaction;
	if (select count(id) from review where review.id = review_id) > 0 then 
		SELECT 
			review.id, 
			review.manga_id,
			review.user_id, 
			user.username, 
			review.review_text,
			review.tag, 
			review.created_at,
			COUNT(liked_review.id) AS total_likes
		FROM 
			review
		JOIN 
			user ON review.user_id = user.id
		LEFT JOIN 
			liked_review ON review.id = liked_review.review_id 
		WHERE 
			review.id = review_id
		 GROUP BY 
            review.id, 
            review.user_id, 
            user.username, 
            review.review_text,
            review.tag, 
            review.created_at;
		commit; 
	else
		signal sqlstate '45000'
		set message_text = 'manga not found';
		rollback;
	end if;
end$$ 
delimiter ;

-- 6. Toggle Like review
delimiter $$
create procedure like_unlike_a_review(
    in n_user_id int,
    in n_review_id int
)
begin
	declare user_active tinyint;
    declare exit handler for sqlexception
    begin
        rollback;
        resignal;
    end;
    start transaction;

   	select is_active into user_active from user where user.id = n_user_id;
  	if user_active = 0 then 
		SIGNAL SQLSTATE '45000'
        SET MESSAGE_TEXT = 'Unauthorized';
	else
        if (select count(user_id) from liked_review  
            where user_id = n_user_id and review_id = n_review_id) > 0 then
            delete from liked_review where user_id = n_user_id 
            and review_id = n_review_id;
        else 
            insert into liked_review (user_id, review_id)
            values (n_user_id, n_review_id);
        end if;
     end if;
    commit;
end$$
delimiter ;