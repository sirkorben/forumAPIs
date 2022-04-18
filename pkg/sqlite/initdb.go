package sqlite

import (
	"database/sql"
	"log"
)

var DB = &sql.DB{}

func DataBase() {
	db, err := sql.Open("sqlite3", "forum.db")
	if err != nil {
		log.Fatal(err)
		return
	}

	DB = db

	// PS! changed User struct in models.go
	const CREATE = `
	create table users(	
		id INTEGER not null primary key autoincrement,
		firstname TEXT not null,
		lastname TEXT not null,
		age INTEGER not null,
		gender TEXT not null,
		username TEXT not null unique,
		email TEXT not null unique, 
		password BLOB not null, 
		created_date TEXT not null );

	create table posts(
		id INTEGER not null primary key autoincrement,
		title TEXT not null,
		contents TEXT not null, 
		create_date TEXT not null, 
		user_id INTEGER not null );

	create table posts_categories(
		post_id int not null, 
		category_id int not null );

	create table categories(
		id INTEGER primary key autoincrement,
		name TEXT not null );

	create table likes(
		post_id INTEGER not null, 
		user_id INTEGER not null);

	create table dislikes(
		post_id INTEGER not null, 
		user_id TEXT not null);

	create table comments(
		id INTEGER not null primary key autoincrement, 
		post_id INTEGER not null, 
		content TEXT not null, 
		user_id TEXT not null);

	create table comment_likes(
		id integer not null constraint comment_likes_pk primary key autoincrement,
		comment_id integer not null, 
		user_id integer not null);

	create table comment_dislikes(
		id INTEGER not null constraint comment_dislikes_pk primary key autoincrement, 
		comment_id INTEGER not null, 
		user_id INTEGER not null);

	create table sessions(
		id TEXT not null primary key, 
		user_id INTEGER not null unique, 
		created_date TEXT not null);
	`

	const INSERT = `
	INSERT INTO posts (id, title, contents, create_date, user_id) VALUES (1, 'Cats are cool!', 'Dogs are same', '2022-02-07 22:11:45', 1);
	INSERT INTO posts (id, title, contents, create_date, user_id) VALUES (2, 'Porsche', 'Fast, rich car originally from Italy', '2022-02-07 22:09:45', 1);
	INSERT INTO posts (id, title, contents, create_date, user_id) VALUES (3, 'API', 'Work in progress, trying hard', '2022-02-07 22:13:45', 1);

	INSERT INTO categories (id, name) VALUES (1, 'Cars');
	INSERT INTO categories (id, name) VALUES (2, 'Animals');
	INSERT INTO categories (id, name) VALUES (3, 'Art');
	INSERT INTO categories (id, name) VALUES (4, 'Games');
	INSERT INTO categories (id, name) VALUES (5, 'Movies');
	INSERT INTO categories (id, name) VALUES (6, 'Misc');

	INSERT INTO posts_categories (post_id, category_id) VALUES (1, 2);
	INSERT INTO posts_categories (post_id, category_id) VALUES (2, 1);
	INSERT INTO posts_categories (post_id, category_id) VALUES (3, 6);

	`
	_, err = db.Exec(CREATE)
	if err != nil {
		return
	}
	_, err = db.Exec(INSERT)
	if err != nil {
		return
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
}
