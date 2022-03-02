create DATABASE if not exists mydb;

create USER if not exists myuser password 'secret';
grant all privileges on database sales to myuser;

create table if not exists sites
(
    id bigserial primary key,
    created_at timestamp,
    updated_at timestamp,
    deleted_at timestamp,
    name varchar(255),
    loading_time interval
);