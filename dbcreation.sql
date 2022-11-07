create table users
(
    id          INTEGER PRIMARY KEY AUTO_INCREMENT,
    name        varchar(255),
    telegram_id INTEGER,
    first_name  varchar(255),
    last_name   varchar(255),
    chat_id     INTEGER,
    created_at  datetime default CURRENT_TIMESTAMP,
    updated_at  datetime,
    deleted_at  datetime
);

create table tasks
(
    id          INTEGER PRIMARY KEY AUTO_INCREMENT,
    title        varchar(255),
    description text,
    end_date  datetime,
    user_id INTEGER,
    created_at  datetime default CURRENT_TIMESTAMP,
    updated_at  datetime,
    deleted_at  datetime,
    FOREIGN KEY (user_id) REFERENCES users(id)
);