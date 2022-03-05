CREATE TABLE transactions(
   id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
   transaction_id VARCHAR UNIQUE,
   amount int,
   account_id VARCHAR,
   category INT,
   created_at timestamp default CURRENT_TIMESTAMP, 
   updated_at timestamp default CURRENT_TIMESTAMP
);