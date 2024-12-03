create table users
(
    id   char(26)    not null
        primary key,
    name varchar(50) not null,
    email  varchar(50) not null,
    deleted_at timestamp default null,
    created_at timestamp default current_timestamp not null,
    updated_at timestamp default current_timestamp not null
);
