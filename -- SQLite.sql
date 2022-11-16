-- SQLite

create table tasks
(
    id   INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER,
    title        varchar(255),
    description  text,
    end_date   datetime,
    created_at  datetime default CURRENT_TIMESTAMP,
    updated_at  datetime,
    deleted_at  datetime,
    FOREIGN KEY(user_id) REFERENCES users(telegram_id)
);

.quit