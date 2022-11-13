CREATE TABLE
    users (
              id INT PRIMARY KEY AUTOINCREMENT,
              name VARCHAR (255),
              telegram_id INT,
              first_name VARCHAR(255),
              last_name VARCHAR(255),
              chat_id INT
);

CREATE TABLE
    tasks (
              id INT PRIMARY KEY AUTOINCREMENT,
              title VARCHAR(255),
              description VARCHAR(255),
              end_date DATE,
              user_id INT NOT NULL,
              FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE ON UPDATE CASCADE
);