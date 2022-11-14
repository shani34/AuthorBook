DROP DATABASE IF EXISTS AuthorBook;
CREATE DATABASE AuthorBook;
USE  AuthorBook;

 CREATE TABLE  books(
     id int not null AUTO_INCREMENT,
     author_id int,
     title varchar(50),
     publication varchar(50),
     published_date varchar(50),
     PRIMARY KEY(id),
     FOREIGN KEY(author_id)REFERENCES author1(author_id)
     );

CREATE TABLE author1(
    author_id int not null AUTO_INCREMENT,
    first_name varchar(50),
    last_name varchar(50),
    dob varchar(10),
    pen_name varchar(50),
    PRIMARY KEY(author_id)
);


