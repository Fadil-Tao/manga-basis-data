show tables;

-- table user
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

-- table author
create table author(
	id int not null auto_increment primary key,
	name varchar(255) not null,
	birthday date not null,
	biography text
);	
drop table if exists author;
-- table genre
create table genre(
	id int not null auto_increment primary key,
	name varchar(255) not null,
	description text not null
);
drop table if exists genre;
create table manga(
	id int primary key auto_increment,
    title varchar(255) not null,
    synopsys text not null,
    manga_status ENUM('finished' , 'in_progress'),
    published_at date not null,
    finished_at date
);
drop table if exists manga;
-- manga genre pivot
CREATE TABLE manga_genre_pivot(
    id_manga int not null,
    id_genre int not null,
    foreign key( id_manga) references manga(id),
    foreign key (id_genre) references genre(id)
);
drop table if exists manga_genre_pivot;
-- author manga pivot
create table manga_author_pivot(
    id_manga integer not null, 
    id_author integer not null,
    foreign key(id_manga) references manga(id),
    foreign key(id_author) references author(id)
);
drop table if exists manga_author_pivot;
-- table review
create table review(
    id int primary key auto_increment,
    manga_id int not null,
    user_id integer not null,
    review_text text not null,
    tag enum('Reccomended' , 'Mixed Feelings' ,'Not Reccomended'),
    foreign key(user_id) references user(id),
    foreign key(manga_id) references manga(id),
    created_at TIMESTAMP DEFAULT current_timestamp
);
drop table if exists review;
-- table liked manga
create table liked_manga(
    id int auto_increment primary key ,
    user_id int not null,
    manga_id int not null, 
    created_at timestamp default current_timestamp ,
    foreign key (user_id) references user(id),
    foreign key (manga_id) references manga(id)
);
drop table if exists liked_manga;
-- table liked review
create table liked_review(
	id int auto_increment primary key not null,
	user_id int not null ,
    review_id int not null,
    foreign key (user_id) references user(id),
    foreign key (review_id) references review(id),
    created_at timestamp default current_timestamp
);
-- table readlist
create table readlist(
 	id int not null auto_increment primary key,
    user_id int not null,
    nama varchar(255) not null,
    description text,
    foreign key (user_id) references user(id),
    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp on update current_timestamp
);
-- table readlist item
create table readlist_item(
	id int not null primary key auto_increment,
	readlist_id int not null,
    read_status enum ('done', 'in_progress', 'later'),
    manga_id integer not null,
    foreign key(manga_id) references manga(id),
    foreign key(readlist_id) references readlist(id),
    added_at timestamp default current_timestamp
);
-- table rating
create table rating(
	id int primary key auto_increment not null,
    manga_id integer not null,
    user_id integer not null,
    foreign key(manga_id) references manga(id),
    foreign key(user_id) references user(id),
    rating SMALLINT NOT NULL CHECK (rating between 1 and 10),
    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp on update current_timestamp
);
