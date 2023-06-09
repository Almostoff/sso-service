CREATE TABLE auth_db.public.auth_admin
(id bigserial not null
     constraint admin_auth_pk
     primary key,
 admin_id bigint not null unique,
 nickname varchar(255),
 password varchar(255),
 second_password varchar(255),
 registration_date timestamp default '0001-01-01 00:00:00'::timestamp without time zone
 );

alter table auth_admin owner to auth_postgres_user;

create unique index users_id_uindex
    on auth_admin (id);

create table auth_db.public.auth_clients
(
    id bigserial not null
        constraint clients_auth_pk
            primary key,
    client_id bigint unique not null ,
    nickname varchar(255),
    email varchar(255),
    phone varchar(255),
    tg varchar(255),
    password varchar(255),
    registration_date timestamp default '0001-01-01 00:00:00'::timestamp without time zone,
    auth_level_id bigint default 0,
    is_dnd bool default false
);

alter table auth_clients owner to auth_postgres_user;

create unique index users_id_uindex
    on auth_clients (id);

CREATE table auth_db.public.auth_clients_level
(
    id bigserial not null
        constraint auth_clients_level_pk
            primary key,
    email_confirm bool,
    phone_confirm bool,
    kyc_confirm bool
);

CREATE table auth_db.public.user_agent
(
    id bigserial not null
        constraint auth_user_agents_pk
            primary key,
    client_id bigint not null,
    ua varchar(255),
    ip varchar(255),
    sign_in_date timestamp default '0001-01-01 00:00:00'::timestamp without time zone,
    logout_date timestamp default '0001-01-01 00:00:00'::timestamp without time zone,
    logout bool
);

Create table auth_codes_confirms
(
    id bigserial not null
        constraint clients_auth_pk
            primary key,
    type varchar(255),
    code_need varchar(255),
    date timestamp default '0001-01-01 00:00:00'::timestamp without time zone,
    destination varchar(255),
    client_id bigint
);

create table inner_connection (
      id bigserial not null
          constraint inner_connection_pm_pk
              primary key,
      base_url text not null,
      public text not null,
      private text not null,
      name varchar(25) not null
);

alter table inner_connection owner to auth_postgres_user;

create unique index inner_connection_id_uindex
    on inner_connection (id);