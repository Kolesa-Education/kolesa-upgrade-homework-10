create table users
(
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    name        varchar(255),
    telegram_id INT,
    first_name  varchar(255),
    last_name   varchar(255),
    chat_id     INT,
    created_at  datetime default CURRENT_TIMESTAMP,
    updated_at  datetime,
    deleted_at  datetime
);

create table tasks
(
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    title        varchar(255),
    description  text,
    end_date     varchar(10),
    user_id      INTEGER,
    created_at  datetime default CURRENT_TIMESTAMP,
    updated_at  datetime,
    deleted_at  datetime,
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);