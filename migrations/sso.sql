BEGIN;
SELECT * FROM pg_available_extensions WHERE name LIKE 'uuid-ossp';
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

create table auth_db.public.client
(
    client_uuid varchar(255) default uuid_generate_v4(),
    language varchar(55) default 'ru',
    nickname varchar(255) default '',
    registration_date timestamp default NOW() + '3 hour',
    last_activity timestamp default NOW() + '3 hour',
    last_login timestamp default NOW()+ '3 hour'
);

create table auth_db.public.client_contact
(
    client_uuid varchar(255) unique not null ,
    email varchar(255)  unique not null,
    phone varchar(255) not null default '',
    tg varchar(255) not null default ''
);

create table auth_db.public.ver_level
(
    client_uuid varchar(255) unique not null ,
    kyc boolean default false,
    email boolean default false,
    phone boolean default false,
    tg boolean default false,
    totp boolean default false,
    resolved_ip boolean default false,
    strong_password boolean default false
);

create table auth_db.public.credential
(
    client_uuid varchar(255) unique not null,
    password varchar(255) not null ,
    totp_secret varchar(255) default '',
    tg_id bigint  default 0
);

create table auth_db.public.session_history
(
    id bigserial,
    client_uuid varchar(255) unique not null ,
    ip varchar(255) not null,
    ua varchar(255) not null,
    login_time timestamp default NOW() + '3 hour',
    logout_time timestamp default '0001-01-01 00:00:00',
    is_logout boolean default false
);

create table auth_db.public.history_passwords
(
    id bigserial not null
        constraint history_passwords_pk
            primary key,
    client_uuid varchar(255) unique not null,
    password varchar(255) not null,
    change_time timestamp default NOW() + '3 hour'
);

create table auth_db.public.history_nickname
(
    id bigserial not null
        constraint history_nickname_pk
            primary key,
    client_uuid varchar(255) unique not null,
    old_nickname varchar(255) not null,
    change_time timestamp default NOW() + '3 hour'
);


CREATE TYPE codes AS ENUM ('recovery_by_email', 'email_confirm', 'KYC', 'phone_confirm', 'confirm_withdraw');


create table auth_db.public.history_passwords
(
    client_uuid varchar(255) unique not null,
    type_code codes not null,
    code_need varchar(255) not null,
    create_time timestamp default NOW() + '3 hour',
    destination varchar(255) default ''
);



create table auth_db.public.service
(
    service_id integer unique not null,
    service_name varchar(55) unique not null,
    public varchar(255) not null,
    private varchar(255) not null,
    base_url varchar(255) not null
);

create table auth_db.public.redirect_url
(
    service_id integer unique not null,
    to_redirect_url_id bigint not null,
    create_time timestamp default NOW() + '3 hour'
);

create table auth_db.public.redirect_url
(
    url_id bigserial unique not null,
    to_redirect_url varchar(255) not null
);
COMMIT;