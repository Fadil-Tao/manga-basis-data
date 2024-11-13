-- +goose Up
-- +goose StatementBegin
create table user(
	id int not null auto_increment primary key,
	username varchar(255) not null,
	email varchar(255) not null unique,
	password text not null,
	created_at TIMESTAMP DEFAULT current_timestamp ,
    updated_at TIMESTAMP DEFAULT current_timestamp on update current_timestamp,
   	is_deleted tinyint default 0
);

create table ban_history(
	id int auto_increment primary key not null,
	banned_user_id int not null,
	reason text not null,
	ban_start_date timestamp default current_timestamp not null,
	ban_end_date timestamp not null, 
	admin_id int not null,
	foreign key (admin_id) references user(id), 
	foreign key (banned_user_id) references user(id)
);

create table user_report(
	id int not null auto_increment primary key,
	reason enum('harrasment','hate','spam') not null,
	user_reported_id int not null,
	report_by_user_id int not null,
	foreign key (user_reported_id) references user(id),
	foreign key (report_by_user_id) references user(id),
	created_at TIMESTAMP DEFAULT current_timestamp 
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists ban_history;
drop table if exists user_report;
drop table if exists user;
-- +goose StatementEnd
