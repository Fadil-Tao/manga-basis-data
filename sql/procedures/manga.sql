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
		if (select count(id) from manga where manga.id = n_id_manga)> 0 and 
			(select count(id) from author where author.id = n_id_author) > 0 then
				insert into author_manga_pivot(id_manga, id_author)
				values (n_id_manga, n_id_author);
				commit;
		else
			signal sqlstate '45000' set message_text = "resource not found";
			rollback;
		end if;
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
		if 	(select count(id) from manga where manga.id = n_id_manga) > 0 and
			(select count(id) from genre where genre.id = n_id_genre) > 0 then
			insert into manga_genre_pivot(id_manga, id_genre)
			values (n_id_manga, n_id_genre); 
			commit;
		else
			signal sqlstate '45000' set message_text = "resource not found";
			rollback ;
		end if;
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


-- 8 Delete manga procedure
delimiter $$
create procedure delete_manga(
	in id_manga int,	
	in id_user int
)
begin 
	declare exit handler for sqlexception
	begin
		rollback;
		resignal;
	end;
	start transaction;
	if(select is_admin(id_user)) = 1 then 
		if (select count(id) from manga where manga.id = id_manga) > 0 then
			delete from manga where manga.id = id_manga;
			commit;
		else
			signal sqlstate '45000' set message_text = "Manga not found";
			rollback;
		end if;
	else
		signal sqlstate '45000' set message_text = "Unauthorized";
		rollback;
	end if;
end$$
delimiter ;
		

-- 9. update manga
delimiter $$ 
create procedure update_manga(
	in id_manga int,
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
	if(select is_admin(user_id)) = 1 then
		if(select count(id) from manga where manga.id = id_manga) > 0 then
			update manga set title = n_title,
				synopsys = n_synopsys,
				manga_status = n_manga_status,
				published_at = n_published_at,
				finished_at = n_finished_at
			where manga.id = id_manga;
			commit ;
		else 
			signal sqlstate '45000' set message_text = "Manga not found";
			rollback;
		end if;
	else 
		signal sqlstate '45000' set message_text = "Unauthorized";
		rollback;
	end if;
end$$
delimiter ;
		

-- 10 . Get ranked manga based on parameter
DELIMITER $$

CREATE PROCEDURE get_manga_ranking(
    IN period ENUM('today', 'month', 'all')
)
BEGIN
    DECLARE exit handler FOR SQLEXCEPTION
    BEGIN
        ROLLBACK;
        RESIGNAL;
    END;

    START TRANSACTION;

    IF period = 'all' THEN
        SELECT * FROM ranking_manga_century;
    ELSEIF period = 'month' THEN
        SELECT * FROM ranking_manga_this_month;
    ELSEIF period = 'today' THEN
        SELECT * FROM ranking_manga_today;
    ELSE
        SIGNAL SQLSTATE '45000'
        SET MESSAGE_TEXT = 'Invalid period parameter';
    END IF;

    COMMIT;
END $$

DELIMITER ;


-- 11 . rangking manga this century
CREATE OR REPLACE VIEW ranking_manga_century AS
SELECT
    m.id,
    m.title,
    m.synopsys,
    m.published_at,
    COALESCE(AVG(r2.rating), NULL) AS average_rating,
    COUNT(DISTINCT r.id) AS total_reviews,
    COUNT(lm.user_id) AS total_likes,
	COUNT(r2.user_id) as total_user_give_rating   
FROM
    manga m
LEFT JOIN review r ON m.id = r.manga_id
LEFT join liked_manga lm  on lm.manga_id = m.id
left join rating r2 on m.id = r2.manga_id 
WHERE
    r.created_at >= '2000-01-01' or
    lm.created_at >= '2000-01-01' or
    r2.created_at >= '2000-01-01'
GROUP BY
    m.id
ORDER BY
    (total_reviews + total_likes + total_user_give_rating ) DESC
LIMIT 10;

--  12 . rangking manga this month
CREATE OR REPLACE VIEW ranking_manga_this_month AS
SELECT
    m.id,
    m.title,
    m.synopsys,
    m.published_at,
    COALESCE(AVG(r2.rating), NULL) AS average_rating,
    COUNT(DISTINCT r.id) AS total_reviews,
    COUNT(DISTINCT lm.user_id) AS total_likes,
    COUNT(DISTINCT r2.user_id) AS total_user_give_rating
FROM
    manga m
LEFT JOIN review r ON m.id = r.manga_id
LEFT JOIN liked_manga lm ON lm.manga_id = m.id
LEFT JOIN rating r2 ON m.id = r2.manga_id
WHERE
    (r.created_at >= DATE_FORMAT(CURRENT_DATE, '%Y-%m-01') OR
    lm.created_at >= DATE_FORMAT(CURRENT_DATE, '%Y-%m-01') OR
    r2.created_at >= DATE_FORMAT(CURRENT_DATE, '%Y-%m-01'))
GROUP BY
    m.id
ORDER BY
    (total_reviews + total_likes + total_user_give_rating) DESC
LIMIT 10;

-- 13. rangking manga today
CREATE OR REPLACE VIEW ranking_manga_today AS
SELECT
    m.id,
    m.title,
    m.synopsys,
    m.published_at,
    COALESCE(AVG(r2.rating), NULL) AS average_rating,
    COUNT(DISTINCT r.id) AS total_reviews,
    COUNT(DISTINCT lm.user_id) AS total_likes,
    COUNT(DISTINCT r2.user_id) AS total_user_give_rating
FROM
    manga m
LEFT JOIN review r ON m.id = r.manga_id
LEFT JOIN liked_manga lm ON lm.manga_id = m.id
LEFT JOIN rating r2 ON m.id = r2.manga_id
WHERE
    (DATE(r.created_at) = CURRENT_DATE OR
    DATE(lm.created_at) = CURRENT_DATE OR
    DATE(r2.created_at) = CURRENT_DATE)
GROUP BY
    m.id
ORDER BY
    (total_reviews + total_likes + total_user_give_rating) DESC
LIMIT 10;

-- 14. toggle like manga
delimiter $$
create procedure toggle_like_manga(
	in n_user_id int,
	in n_manga_id int
)
begin 
	declare user_active tinyint;
    declare exit handler for sqlexception
    begin
        rollback;
        resignal;
    end;
    start transaction;

   	select is_active into user_active from user where user.user_id = n_user_id;
  	if user_active = 0 then 
		SIGNAL SQLSTATE '45000'
        SET MESSAGE_TEXT = 'User Not Active';
	else
		if(select count(user_id) from liked_manga 
		where user_id = n_user_id and manga_id = n_manga_id) > 0 then 
			delete from liked_manga where user_id = n_user_id 
			and manga_id = n_manga_id;
		else 
			insert into liked_manga(user_id,manga_id)
			values(n_user_id,n_manga_id);
		end if;
	end if;
	commit;
end$$


-- 15. like a manga
delimiter $$
create procedure like_unlike_a_manga(
	in n_user_id int,
	in n_manga_id int
)
begin 
    declare exit handler for sqlexception
    begin
        rollback;
        resignal;
    end;
    start transaction;

  	if (is_active(n_user_id)) = 0 then 
		SIGNAL SQLSTATE '45000'
        SET MESSAGE_TEXT = 'Unauthorized';
	else
		if(select count(user_id) from liked_manga 
		where user_id = n_user_id and manga_id = n_manga_id) > 0 then 
			delete from liked_manga where user_id = n_user_id 
			and manga_id = n_manga_id;
		else 
			insert into liked_manga(user_id,manga_id)
			values(n_user_id,n_manga_id);
		end if;
	end if;
	commit;
end$$
delimiter ; 


-- 16 get all manga

DELIMITER $$

CREATE PROCEDURE get_all_manga(
    IN `n_limit` INT,
    IN orderby VARCHAR(30),
    IN sort ENUM('asc', 'desc'),
    IN manga_title VARCHAR(255)
)
BEGIN
    DECLARE query TEXT;
    DECLARE exit handler FOR sqlexception 
    BEGIN 
        ROLLBACK;
        RESIGNAL;
    END;

    START TRANSACTION;

    -- Ensure safe defaults for inputs
    IF orderby NOT IN ('published_at', 'total_likes', 'average_rating', 'title', 'total_user_rated', 'total_reviews') THEN
        SET orderby = 'title';
    END IF;

    IF sort NOT IN ('asc', 'desc') THEN
        SET sort = 'desc';
    END IF;

    IF `n_limit` < 0 THEN
        SET `n_limit` = 0;
    END IF;

    -- Build query based on whether manga_title is empty
    IF manga_title = '' THEN
        SET @query = CONCAT(
            'SELECT 
                m.id AS id, 
                m.title AS title, 
                m.manga_status AS status, 
                m.published_at, 
                m.finished_at, 
                COALESCE(AVG(r2.rating), NULL) AS average_rating,
                COUNT(DISTINCT r.id) AS total_reviews,
                COUNT(DISTINCT lm.user_id) AS total_likes,
                COUNT(DISTINCT r2.user_id) AS total_user_rated
            FROM 
                manga m
            LEFT JOIN 
                liked_manga lm ON lm.manga_id = m.id
            LEFT JOIN 
                rating r2 ON m.id = r2.manga_id
            GROUP BY 
                m.id
            ORDER BY ', orderby, ' ', sort, '
            LIMIT ', n_limit);
    ELSE
        SET @query = CONCAT(
            'SELECT 
                m.id AS id, 
                m.title AS title, 
                m.manga_status AS status, 
                m.published_at, 
                m.finished_at, 
                COALESCE(AVG(r2.rating), NULL) AS average_rating,
                COUNT(DISTINCT r.id) AS total_reviews,
                COUNT(DISTINCT lm.user_id) AS total_likes,
                COUNT(DISTINCT r2.user_id) AS total_user_rated
            FROM 
                manga m
            LEFT JOIN 
                liked_manga lm ON lm.manga_id = m.id
            LEFT JOIN 
                rating r2 ON m.id = r2.manga_id
            WHERE 
                m.title LIKE CONCAT(', QUOTE(manga_title), ', "%")
            GROUP BY 
                m.id
            ORDER BY ', orderby, ' ', sort, '
            LIMIT ', n_limit);
    END IF;

    -- Prepare and execute the query
    PREPARE stmt FROM @query;
    EXECUTE stmt;
    DEALLOCATE PREPARE stmt;

    COMMIT;
END$$

DELIMITER ;
