create table transactions(
   id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
   amount int,
   account_id int UNIQUE,
   category INT,
   created_at timestamp default CURRENT_TIMESTAMP, 
   updated_at timestamp default CURRENT_TIMESTAMP
);