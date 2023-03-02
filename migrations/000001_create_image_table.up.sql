create table images
(
    id           serial primary key,
    image_path   varchar(255) unique not null,
    content_type varchar(50)         not null
);