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
	declare exit handler for sqlexception
	begin
		rollback;
		resignal;
	end;	
		    
	start transaction;
	if (select count(id) from user where user.id = n_user_id) < 1 then 
		SIGNAL SQLSTATE '45000'
        SET MESSAGE_TEXT = 'Unauthorized';
	else
		if (select count(id) from manga where manga.id = n_manga_id) > 0 then
			if (select count(user_id) from review where review.user_id = n_user_id and review.manga_id = n_manga_id) > 0 then 
				signal sqlstate '45000' 
				set message_text = 'Youre already submitted a review for this manga';
				rollback;
			else
				if n_tag not in ('Reccomended' , 'Mixed Feelings' ,'Not Reccomended') then
					signal sqlstate '45000'
					set message_text = 'invalid tag';
					rollback;
				end if;
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
drop procedure get_review_from_manga;
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
			review.user_id, 
			user.username, 
			review.review_text,
			review.tag, 
			review.created_at,
			COUNT(liked_review.user_id) AS total_likes
		FROM 
			review
		JOIN 
			user ON review.user_id = user.id
		LEFT JOIN 
			liked_review ON review.user_id = liked_review.review_user_id and review.manga_id = liked_review.manga_id 
		WHERE 
			review.manga_id = manga_id
		GROUP BY 
		review.user_id, user.username, review.review_text, 
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
	in doer_id int,
	in n_manga_id int,
	in n_user_id int
)
begin 
	declare exit handler for sqlexception
	begin
		rollback;
		resignal;
	end;

	start transaction;
	if (select count(manga_id) from review where review.manga_id  = n_manga_id and review.user_id = n_user_id) < 1 then
		SIGNAL SQLSTATE '45000'
        SET MESSAGE_TEXT = 'review not found';
		rollback;	
	end if;
	if (select user_id from review where review.manga_id  = n_manga_id and review.user_id = n_user_id) != doer_id then  
		SIGNAL SQLSTATE '45000'
        SET MESSAGE_TEXT = 'Unauthorized';
		rollback;	
	end if;
	if (select is_admin(n_user_id)) < 1 then
		SIGNAL SQLSTATE '45000'
        SET MESSAGE_TEXT = 'Unauthorized';
		rollback;	
	end if;

	delete from review where review.user_id  = n_user_id and review.manga_id  = n_manga_id;
	commit;
end$$
delimiter ;

-- 4 update a review
drop procedure update_a_review;
delimiter $$
create procedure update_a_review(
	in doer_id int,
	in n_user_id int,
	in n_manga_id int,
	in review_text text,
	in n_tag varchar(255)
)
begin
	declare exit handler for sqlexception
		begin 
			rollback;
			resignal;
		end;
	start transaction;
	if (select count(manga_id) from review where review.manga_id  = n_manga_id and review.user_id = n_user_id) < 1 then
		SIGNAL SQLSTATE '45000'
        SET MESSAGE_TEXT = 'review not found';
		rollback;	
	end if;
	if (select user_id from review where review.manga_id  = n_manga_id and review.user_id = n_user_id) != doer_id then  
		SIGNAL SQLSTATE '45000'
        SET MESSAGE_TEXT = 'Unauthorized';
		rollback;	
	end if;

  IF n_tag IS NOT NULL AND n_tag NOT IN ('Reccomended', 'Mixed Feelings', 'Not Reccomended') THEN
        SIGNAL SQLSTATE '45000'
        SET MESSAGE_TEXT = 'Invalid tag';
    END IF;
    IF review_text IS NOT NULL THEN
        UPDATE review
        SET review_text = review_text
        WHERE user_id = n_user_id AND manga_id = n_manga_id;
    END IF;
    IF n_tag IS NOT NULL THEN
        UPDATE review
        SET tag = n_tag
        WHERE user_id = n_user_id AND manga_id = n_manga_id;
    END IF;

    COMMIT;
end$$
delimiter ; 	


-- 5 Get one review by its id 
drop procedure get_review_by_id;
delimiter $$
create procedure get_review_by_id(
	in manga_id int,
	in user_id int 
)
begin
	declare exit handler for sqlexception
		begin 
			rollback;
			resignal;
		end;
	start transaction;
	if (select count(user_id) from review where review.user_id = user_id and review.manga_id = manga_id) > 0 then 
		SELECT 
			review.manga_id,
			review.user_id, 
			user.username, 
			review.review_text,
			review.tag, 
			review.created_at,
			COUNT(liked_review.user_id) AS total_likes
		FROM 
			review
		JOIN 
			user ON review.user_id = user.id
		LEFT JOIN 
			liked_review ON review.user_id = liked_review.review_user_id and review.manga_id = liked_review.manga_id 
		WHERE 
			review.user_id = user_id and review.manga_id = manga_id
		 GROUP BY 
		 	review.manga_id,
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
drop procedure like_unlike_a_review;
drop procedure like_unlike_a_review;
delimiter $$
create procedure like_unlike_a_review(
    in doer_id int,
    in n_manga_id int,
    in n_user_id int
)
begin
    declare exit handler for sqlexception
    begin
        rollback;
        resignal;
    end;
    start transaction;

  	if (select count(id) from user where user.id = doer_id) < 1 then 
		SIGNAL SQLSTATE '45000'
        SET MESSAGE_TEXT = 'Unauthorized';
       	rollback;
    end  if ;
   if (select count(user_id) from review where review.manga_id = n_manga_id and review.user_id = n_user_id) < 1 then
   		SIGNAL SQLSTATE '45000'
        SET MESSAGE_TEXT = 'review not found';
       	rollback;
    end  if ;
    if (select count(user_id) from liked_review  where user_id = doer_id 
    and (liked_review.manga_id = n_manga_id and liked_review.review_user_id = n_user_id)) > 0 then
        delete from liked_review where user_id = doer_id 
        and liked_review.manga_id = n_manga_id and liked_review.review_user_id = n_user_id;
    else 
        insert into liked_review (user_id, manga_id,review_user_id)
        values (doer_id, n_manga_id,n_user_id);
        end if;
    commit;
end$$
delimiter ;