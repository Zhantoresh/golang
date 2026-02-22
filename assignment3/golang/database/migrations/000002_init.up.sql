create table if not exists users ( 
  id serial primary key, 
  name varchar(255) not null, 
  email varchar(255) not null, 
  city varchar(255) not null, 
  age int not null 
); 
 
insert into users (name, email, city, age) values ('John Doe', 'john@example.com', 'New York', 30);