delimiter $$
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
	select is_active into user_active from user where user.user_id = n_user_id;

	if user_active = 0 then 
		SIGNAL SQLSTATE '45000'
        SET MESSAGE_TEXT = 'User Not Active';
	else
   		insert into review(manga_id, user_id, review_text,tag)
		values (n_manga_id,n_user_id,n_review_text,n_tag);
		commit;
	end if ;
end$$
delimiter ;

