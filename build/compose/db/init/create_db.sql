create DATABASE mydb;

create USER myuser password 'secret';
grant all privileges on database mydb to myuser;

drop table if exists sites;
create table sites
(
    id bigserial primary key,
    created_at timestamp,
    updated_at timestamp,
    deleted_at timestamp,
    name varchar(255),
    loading_time interval
);