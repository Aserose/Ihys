create table users (
    tg_user_id bigint,
    tg_chat_id bigint,
    PRIMARY KEY(tg_user_id, tg_chat_id)
);

create table spotify (
    log_id integer UNIQUE default(1),
    encrypted_key varchar,
    update_date timestamp,
     Constraint singlerow CHECK (log_id = 1)
) inherits (users);

create table lastFm (
    log_id integer UNIQUE default(1),
    encrypted_key varchar,
    update_date timestamp,
     Constraint singlerow CHECK (log_id = 1)
) inherits (users);

create table vk (
    log_id integer UNIQUE default(1),
    encrypted_key varchar,
    update_date timestamp,
     Constraint singlerow CHECK (log_id = 1)
) inherits (users);