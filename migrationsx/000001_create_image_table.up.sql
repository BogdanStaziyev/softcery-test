create table images
(
    id         serial primary key,
    image_path varchar(255) unique not null
);