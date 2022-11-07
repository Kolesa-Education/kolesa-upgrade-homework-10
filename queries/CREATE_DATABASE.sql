upgrade.db 

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


CREATE TABLE tasks( 
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title   varchar(255),
    description     varchar(255),
    end_date  datetime,
    user_id INTEGER NOT NULL,
    FOREIGN KEY (user_id)
    REFERENCES users (id)
);
