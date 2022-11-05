CREATE TABLE tasks
(
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    title       varchar(255),
    description TEXT,
    end_date    TEXT,
    telegram_id INTEGER,
    created_at  datetime default CURRENT_TIMESTAMP,
    updated_at  datetime,
    deleted_at  datetime,
    FOREIGN KEY(telegram_id) REFERENCES users(telegram_id)
);

INSERT INTO tasks (title, description, end_date, telegram_id)
VALUES("Homework10", "Do telegram bot", "Jan11", 806808831);