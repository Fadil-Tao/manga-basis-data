use manga_basis_data;

-- 1. add manga 
drop procedure add_manga;
delimiter $$
create procedure add_manga(
	in n_title varchar(255),
	in n_synopsys text,
	in n_manga_status varchar(20),
	in n_published_at date,
	in n_finished_at date,
	in user_id int
)
begin 
	declare finished_date date;
	declare exit handler for sqlexception 
	begin 
		rollback;
		resignal;
	end;
	
	start transaction;
	set finished_date = n_finished_at; 
	if n_manga_status not in ('finished','in_progress') then 
	 	signal sqlstate '45000' set message_text = "invalid manga status";
		rollback;
	end if;
   	if n_manga_status = 'in_progress' and n_finished_at is not null then
        set finished_date = null; 
    elseif n_manga_status = 'finished' then
        if n_finished_at is null then
            signal sqlstate '45000' set message_text = "invalid finished manga must have a finished date";
            rollback;
        else
            set finished_date = n_finished_at;
        end if;
    end if;
	if (select is_admin(user_id)) = 1 then 
		insert into manga (title,synopsys,manga_status,published_at,finished_at)
		values (n_title,n_synopsys,n_manga_status,n_published_at,finished_date);
		commit;
	else
		signal sqlstate '45000' set message_text = "Unauthorized";
		rollback;
	end if;
end$$
delimiter ;

-- 2. connect author and manga
drop procedure connect_author_manga;
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
			if (select count(id_manga) from author_manga_pivot where author_manga_pivot.id_manga = n_id_manga and author_manga_pivot.id_author  = n_id_author) > 0 then
				signal sqlstate '45000' set message_text = "already associated";
				rollback ;
			end if;
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
drop procedure connect_genre_manga;
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
			if (select count(id_manga) from manga_genre_pivot where manga_genre_pivot.id_manga = n_id_manga and manga_genre_pivot.id_genre = n_id_genre) > 0 then
				signal sqlstate '45000' set message_text = "already associated";
				rollback ;
			end if;
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
drop procedure get_manga_detail;
delimiter $$
create procedure get_manga_detail(
	in in_manga_id int
) 
begin
	SELECT m.id,m.title, m.synopsys, m.manga_status, 
    m.published_at, 
    m.finished_at,
	COALESCE(AVG(r2.rating), 0) AS average_rating,
    COUNT(DISTINCT r.user_id) AS total_reviews,
    COUNT(lm.user_id) AS total_likes,
	COUNT(r2.user_id) as total_user_give_rating  
	FROM manga m 
	LEFT JOIN review r ON m.id = r.manga_id
	LEFT join liked_manga lm  on lm.manga_id = m.id
	left join rating r2 on m.id = r2.manga_id 
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
drop procedure update_manga;
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

	if (select  is_admin(user_id)) < 1 then
		signal sqlstate '45000' set message_text = "Unauthorized";
		rollback;	
	end if;
	if (select count(id) from manga where manga.id = id_manga) < 1 then
		signal sqlstate '45000' set message_text = "manga not found";
		rollback;	
	end if;
	if n_title != '' then 
		update manga set manga.title = n_title where manga.id = id_manga;
	end if ;
	if n_synopsys is not null then 
        update manga set manga.synopsys = n_synopsys where manga.id = id_manga;
    end if;
	if n_manga_status is not null and n_manga_status in ('finished', 'in_progress') then
        if n_manga_status = 'in_progress' then
            update manga set manga.manga_status = n_manga_status, manga.finished_at = null where manga.id = id_manga;
        elseif n_manga_status = 'finished' then
            if n_finished_at is null then
                signal sqlstate '45000' set message_text = "invalid finished manga must have finished date";
                rollback;
            else
                update manga set manga.manga_status = n_manga_status, manga.finished_at = n_finished_at where manga.id = id_manga;
            end if;
        end if;
    end if;
	if n_published_at is not null then 
        update manga set manga.published_at = n_published_at where manga.id = id_manga;
    end if;

    if n_finished_at is not null and n_manga_status != 'in_progress' then
        update manga set manga.finished_at = n_finished_at where manga.id = id_manga;
    end if;
    commit;
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
    COALESCE(AVG(r2.rating), 0) AS average_rating,
    COUNT(DISTINCT r.user_id) AS total_reviews,
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
    COALESCE(AVG(r2.rating), 0) AS average_rating,
    COUNT(DISTINCT r.user_id) AS total_reviews,
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
    COALESCE(AVG(r2.rating), 0) AS average_rating,
    COUNT(DISTINCT r.user_id) AS total_reviews,
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
drop procedure toggle_like_manga;
delimiter $$
create procedure toggle_like_manga(
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

  	if (select count(id) from user where user.id = n_user_id) < 1 then 
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

  	if  (select count(id) from user where user.id = n_user_id) < 1 then 
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

    IF orderby NOT IN ('published_at', 'total_likes', 'average_rating', 'title', 'total_user_rated', 'total_reviews') THEN
        SET orderby = 'title';
    END IF;

    IF sort NOT IN ('asc', 'desc') THEN
        SET sort = 'desc';
    END IF;

    IF `n_limit` < 0 THEN
        SET `n_limit` = 0;
    END IF;

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

    PREPARE stmt FROM @query;
    EXECUTE stmt;
    DEALLOCATE PREPARE stmt;

    COMMIT;
END$$
DELIMITER ;

-- 17. delete association between manga and author
drop procedure delete_association_manga_author;
delimiter $$
create procedure delete_association_manga_author(
	in user_id int,
	in manga_id int,
	in author_id int
)
begin
declare exit handler for sqlexception
	begin
	rollback;
	resignal;
	end;
	
	start transaction;
	if (select is_admin(user_id)) < 1 then 
		signal sqlstate '45000'
		set message_text= "Unauthorized";
	end if;
	if (select count(id) from manga where manga.id = manga_id) < 1 then
		signal sqlstate '45000'
		set message_text = "Manga not found";
	end if ;
	if (select count(id)from author where author.id = author_id) < 1 then
		signal sqlstate '45000'
		set message_text = "Author not found";
	end if;

	delete from author_manga_pivot amp 
	where amp.id_manga = manga_id and amp.id_author = author_id;
	commit;
end$$
delimiter ;



-- 18. delete assoctiaon manga and genre
drop procedure  delete_association_manga_genre;
delimiter $$
create procedure delete_association_manga_genre(
	in user_id int,
	in manga_id int,
	in genre_id int
)
begin
declare exit handler for sqlexception
	begin
	rollback;
	resignal;
	end;
	
	start transaction;
	if (select is_admin(user_id)) < 1 then 
		signal sqlstate '45000'
		set message_text= "Unauthorized";
	end if;
	if (select count(id) from manga where manga.id = manga_id) < 1 then
		signal sqlstate '45000'
		set message_text = "Manga not found";
	end if ;
	if (select count(id)from genre where genre.id = genre_id) then
		signal sqlstate '45000'
		set message_text = "Genre not found";
	end if;

	delete from manga_genre_pivot mgp 
	where mgp.id_manga = manga_id and mgp.id_genre = genre_id;
	commit;
end$$
delimiter ;


-- 19 . get all manga with search its name
drop procedure get_all_manga_with_search;
delimiter $$
create procedure get_all_manga_with_search(
	in input varchar(255)
)
begin
declare exit handler for sqlexception
	begin
		rollback;
		resignal;
	end;

	start transaction;
	if input != '' then
	SELECT 
        m.id AS id, 
        m.title AS title, 
        m.manga_status AS status, 
        m.published_at, 
        m.finished_at, 
        COALESCE(AVG(r2.rating), 0) AS average_rating,
        COUNT(DISTINCT r.user_id) AS total_reviews,
    	COUNT(DISTINCT lm.user_id) AS total_likes,
        COUNT(DISTINCT r2.user_id) AS total_user_rated
        FROM 
            manga m
        LEFT JOIN 
            liked_manga lm ON lm.manga_id = m.id
        LEFT JOIN 
            rating r2 ON m.id = r2.manga_id
        left join 
        	review r on r.manga_id = m.id
		where m.title like concat(input, '%')
        GROUP BY 
            m.id
        ORDER BY m.title asc;
	commit;
	else
		SELECT 
			m.id AS id, 
			m.title AS title, 
			m.manga_status AS status, 
			m.published_at, 
			m.finished_at, 
			COALESCE(AVG(r2.rating), 0) AS average_rating,
			COUNT(DISTINCT r.user_id) AS total_reviews,
			COUNT(DISTINCT lm.user_id) AS total_likes,
			COUNT(DISTINCT r2.user_id) AS total_user_rated
			FROM 
				manga m
			LEFT JOIN 
				liked_manga lm ON lm.manga_id = m.id
			LEFT JOIN 
				rating r2 ON m.id = r2.manga_id
			left join 
        		review r on r.manga_id = m.id
			GROUP BY 
				m.id
			ORDER BY m.title asc;
		commit;
	end if;
end$$
delimiter ;
	
-- 20 .  update set manga status to finish

