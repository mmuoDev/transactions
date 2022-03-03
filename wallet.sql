create table wallet(
   id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
   previous_balance int,
   current_balance int,
   account_id int UNIQUE,
   created_at timestamp default CURRENT_TIMESTAMP, 
   updated_at timestamp default CURRENT_TIMESTAMP
);