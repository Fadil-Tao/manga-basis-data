use manga_basis_data;

-- isinnya trigger doang


-- trigger user delete
delimiter $$
create trigger before_user_delete 
before delete on user
for each row
begin
	delete from readlist where readlist.owned_by = old.user_id;
end$$
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

-- trigger for deleting readlist item if readlist is deleted
delimiter $$
create trigger before_readlist_delete
before delete on readlist 
for each row 
begin 
	delete from readlist_item where readlist_item.readlist_id = old.readlist_id;
end$$
delimiter ; 
delimiter $$

-- trigger before manga delete

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


