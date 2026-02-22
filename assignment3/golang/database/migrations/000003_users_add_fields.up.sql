alter table users
add column if not exists email varchar(255) not null default 'john@example.com',
add column if not exists age int not null default 20,
add column if not exists city varchar(255) not null default 'Almaty';