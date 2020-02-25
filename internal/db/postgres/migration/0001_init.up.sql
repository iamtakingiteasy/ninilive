create table users (
    user_id       bigserial primary key,
    user_name     text,
    user_login    text unique,
    user_password text,
    user_mod      bool
);

create table channels (
    channel_id   bigserial primary key,
    channel_name text unique,
    channel_order int
);

create table messages (
    message_id         bigserial primary key,
    message_channel_id bigint references channels(channel_id),
    message_body       text,
    message_time       timestamp,
    message_edit       timestamp,
    message_trip       text,
    message_origin     text,
    message_remote     text,
    message_file_name  text,
    message_file_path  text,
    message_user       bigint references users(user_id)
);

insert into users (
    user_id,
    user_name,
    user_login,
    user_password,
    user_mod
) values
(0, 'Anonymous', '', '', false);