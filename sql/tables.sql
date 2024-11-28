show tables;

-- 1. table user
create table user(
	id int not null auto_increment primary key,
	username varchar(255) not null,
	email varchar(255) not null unique,
	password text not null,
	created_at TIMESTAMP DEFAULT current_timestamp ,
    updated_at TIMESTAMP DEFAULT current_timestamp on update current_timestamp,
    is_active tinyint default 1,
   	is_deleted tinyint default 0,
   	is_admin tinyint default 0
);
drop table if exists user;
select * from user;

-- 2. table author
create table author(
	id int not null auto_increment primary key,
	name varchar(255) not null,
	birthday date not null,
	biography text
);	
drop table if exists author;
-- 3. table genre
create table genre(
	id int not null auto_increment primary key,
	name varchar(255) not null,
	description text not null
);
drop table if exists genre;
-- 4. table manga
create table manga(
	id int primary key auto_increment,
    title varchar(255) not null,
    synopsys text not null,
    manga_status ENUM('finished' , 'in_progress'),
    published_at date not null,
    finished_at date
);
drop table if exists manga;
-- 5. manga genre pivot
CREATE TABLE manga_genre_pivot(
    id_manga int not null,
    id_genre int not null,
    primary key(id_manga, id_genre),
    foreign key( id_manga) references manga(id),
    foreign key (id_genre) references genre(id)
);
drop table if exists manga_genre_pivot;
-- 6. author manga pivot
create table manga_author_pivot(
    id_manga integer not null, 
    id_author integer not null,
    primary key(id_manga,id_author),
    foreign key(id_manga) references manga(id),
    foreign key(id_author) references author(id)
);
drop table if exists manga_author_pivot;
-- 7. table review
create table review(
    id int primary key auto_increment,
    manga_id int not null,
    user_id integer not null,
    unique(manga_id, user_id),
    review_text text not null,
    tag enum('Reccomended' , 'Mixed Feelings' ,'Not Reccomended'),
    foreign key(user_id) references user(id),
    foreign key(manga_id) references manga(id),
    created_at TIMESTAMP DEFAULT current_timestamp
);
drop table if exists review;
-- 8. table liked manga
create table liked_manga(
    user_id int not null,
    manga_id int not null,
    primary key (user_id, manga_id),
    created_at timestamp default current_timestamp ,
    foreign key (user_id) references user(id),
    foreign key (manga_id) references manga(id)
);
drop table if exists liked_manga;
-- 9. table liked review
create table liked_review(
	user_id int not null ,
    review_id int not null,
    primary key (user_id, review_id),
    foreign key (user_id) references user(id),
    foreign key (review_id) references review(id),
    created_at timestamp default current_timestamp
);
drop table if exists liked_review;
-- 10. table readlist
create table readlist(
 	id int not null auto_increment primary key,
    user_id int not null,
    nama varchar(255) not null,
    description text,
    foreign key (user_id) references user(id),
    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp on update current_timestamp
);
-- 11. table readlist item
create table readlist_item(
	id int not null primary key auto_increment,
	readlist_id int not null,
    read_status enum ('done', 'in_progress', 'later'),
    manga_id integer not null,
    unique (readlist_id, manga_id),
    foreign key(manga_id) references manga(id),
    foreign key(readlist_id) references readlist(id),
    added_at timestamp default current_timestamp
);
-- 12. table rating
create table rating(
    manga_id integer not null,
    user_id integer not null,
    PRIMARY KEY (manga_id, user_id),
    foreign key(manga_id) references manga(id),
    foreign key(user_id) references user(id),
    rating SMALLINT NOT NULL CHECK (rating between 1 and 10),
    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp on update current_timestamp
);
