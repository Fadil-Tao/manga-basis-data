use manga_basis_data;

-- isinnya trigger doang


-- trigger user delete
delimiter $$
create trigger before_user_delete
before delete on user
for each row
BEGIN 
	delete from readlist where readlist.user_id = old.id;
	delete from liked_manga where liked_manga.user_id = old.id;
	delete from liked_review where liked_review.user_id = old.id;
	delete from review where review.user_id = old.id;
	delete from rating where rating.user_id = old.id;
END$$	
delimiter ;



-- trigger for deleting readlist item if readlist is deleted
delimiter $$
create trigger before_readlist_delete
before delete on readlist 
for each row 
begin 
	delete from readlist_item where readlist_item.readlist_id = old.readlist_id;
end$$
delimiter ;


-- trigger before manga delete
delimiter $$ 
create trigger before_manga_delete
before delete on manga
for each row 
begin 
	delete from readlist_item where readlist_item.manga_id = old.id;
	delete from review where review.manga_id = old.id;
	delete from liked_manga where liked_manga.manga_id = old.id;
	delete from author_manga_pivot where author_manga_pivot.id_manga = old.id;
	delete from manga_genre_pivot where manga_genre_pivot.id_manga = old.id;
	delete from rating where rating.manga_id = old.id;
end$$
delimiter ;


-- trigger before genre delete
delimiter $$
create trigger before_genre_delete
before delete on genre
for each row
begin
	delete from manga_genre_pivot where manga_genre_pivot.id_genre = old.id;
end$$
delimiter ;
-- trigger before author delete
delimiter $$
create trigger before_author_delete
before delete on genre
for each row
begin
	delete from manga_author_pivot where manga_author_pivot.id_author = old.id;	
end$$
delimiter ;


-- trigger before review
delimiter $$ 
create trigger before_review_delete
before delete on review
for each row 
begin
	delete from liked_review where liked_review.review_id = old.id;
end$$
delimiter ;



