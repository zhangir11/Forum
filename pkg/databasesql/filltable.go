package databasesql

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3" // SQLite driver
)

// Initiate ...
func (d *DatabaseSQL) Initiate(DbAdress string) error {
	fmt.Println("there :", DbAdress)
	d.DatabaseAddress = DbAdress
	store, err := sql.Open("sqlite3", DbAdress)
	if err != nil {
		log.Println("Failed Creating db: ", err.Error())
		return err
	}
	// Complete checking of connection with DB
	if err := store.Ping(); err != nil {
		log.Println("Failed accesing db: ", err.Error())
		return err
	}
	d.Db = store // fill "db" field with completely configured DB

	users, err := d.Db.Prepare(`CREATE TABLE IF NOT EXISTS Users (
		id INTEGER PRIMARY KEY NOT NULL, 
		firstname TEXT NOT NULL, 
		lastname TEXT NOT NULL, 
		username TEXT NOT NULL UNIQUE, 
		email TEXT NOT NULL UNIQUE, 
		password TEXT NOT NULL
	)`)
	if err != nil {
		return err
	}

	_, err = users.Exec()
	if err != nil {
		fmt.Println(err.Error())
	}

	err = users.Close()
	if err != nil {
		return err
	}

	sessions, err := d.Db.Prepare(`CREATE TABLE IF NOT EXISTS Sessions (
		userID INTEGER NOT NULL,
		cookieName TEXT NOT NULL,
		cookieValue TEXT NOT NULL UNIQUE,
		FOREIGN KEY(userID) REFERENCES Users(id)
	)`)
	if err != nil {
		return err
	}
	_, err = sessions.Exec()
	if err != nil {
		fmt.Println(err.Error())
	}

	err = sessions.Close()
	if err != nil {
		return err
	}
	posts, err := d.Db.Prepare(`CREATE TABLE IF NOT EXISTS Posts (
		postid INTEGER PRIMARY KEY NOT NULL,
		userid INTEGER NOT NULL,
		author TEXT NOT NULL, 
		title TEXT NOT NULL, 
		content TEXT NOT NULL,
		creationDate TEXT NOT NULL,
		FOREIGN KEY(userid) REFERENCES Users(id) 
	)`)
	if err != nil {
		return err
	}
	_, err = posts.Exec()
	if err != nil {
		fmt.Println(err.Error())
	}

	err = posts.Close()
	if err != nil {
		return err
	}

	threads, err := d.Db.Prepare(`CREATE TABLE IF NOT EXISTS Threads (
		ID INTEGER PRIMARY KEY NOT NULL, 
		Name TEXT NOT NULL UNIQUE
	)`)
	if err != nil {
		return err
	}
	_, err = threads.Exec()
	if err != nil {
		fmt.Println(err.Error())
	}

	err = threads.Close()
	if err != nil {
		return err
	}

	postMap, err := d.Db.Prepare(`CREATE TABLE IF NOT EXISTS PostMapping (
		postID INTEGER NOT NULL, 
		threadID INTEGER NOT NULL,
		FOREIGN KEY(postID) REFERENCES Posts(post_id),
		FOREIGN KEY(threadID) REFERENCES Threads(ID)
	)`)
	if err != nil {
		return err
	}
	_, err = postMap.Exec()
	if err != nil {
		fmt.Println(err.Error())
	}

	err = postMap.Close()
	if err != nil {
		return err
	}

	comments, err := d.Db.Prepare(`CREATE TABLE IF NOT EXISTS Comments (
		ID INTEGER PRIMARY KEY NOT NULL,
		postID INTEGER NOT NULL,
		authorID INTEGER NOT NULL,
		content TEXT NOT NULL, 
		creationDate TEXT NOT NULL,
		FOREIGN KEY(postID) REFERENCES Posts(post_id)
		)`)
	if err != nil {
		return err
	}

	_, err = comments.Exec()
	if err != nil {
		fmt.Println(err.Error())
	}

	err = comments.Close()
	if err != nil {
		return err
	}

	postRating, err := d.Db.Prepare(`CREATE TABLE IF NOT EXISTS PostRating (
		postID INTEGER NOT NULL UNIQUE,
		likeCount INTEGER NOT NULL,
		dislikeCount INTEGER NOT NULL,
		FOREIGN KEY(postID) REFERENCES Posts(post_id)
	)`)
	if err != nil {
		return err
	}
	_, err = postRating.Exec()
	if err != nil {
		fmt.Println(err.Error())
	}

	err = postRating.Close()
	if err != nil {
		return err
	}

	commRating, err := d.Db.Prepare(`CREATE TABLE IF NOT EXISTS CommentRating (
		commentID INTEGER NOT NULL UNIQUE,
		likeCount INTEGER NOT NULL,
		dislikeCount INTEGER NOT NULL,
		FOREIGN KEY(commentID) REFERENCES Comments(ID)
		)`)

	if err != nil {
		return errors.New(err.Error() + " comentrating")
	}

	_, err = commRating.Exec()
	if err != nil {
		fmt.Println(err.Error() + " comentrating")
	}

	err = commRating.Close()
	if err != nil {
		return err
	}

	rateUP, err := d.Db.Prepare(`CREATE TABLE IF NOT EXISTS RateUserPost (
		userID INTEGER NOT NULL,
		postID INTEGER NOT NULL,
		kind INTEGER NOT NULL,
		FOREIGN KEY(userID) REFERENCES Users(ID),
		FOREIGN KEY(postID) REFERENCES Posts(post_id)
	)`)
	if err != nil {
		return err
	}
	_, err = rateUP.Exec()
	if err != nil {
		fmt.Println(err.Error())
	}

	err = rateUP.Close()
	if err != nil {
		return err
	}

	rateUC, err := d.Db.Prepare(`CREATE TABLE IF NOT EXISTS RateUserComment (
		commentID INTEGER NOT NULL,
		userID INTEGER NOT NULL,
		kind INTEGER NOT NULL,
		FOREIGN KEY(commentID) REFERENCES Comments(ID),
		FOREIGN KEY(userID) REFERENCES Users(ID)
	)`)
	if err != nil {
		return err
	}
	_, err = rateUC.Exec()
	if err != nil {
		fmt.Println(err.Error())
	}

	err = rateUC.Close()
	if err != nil {
		return err
	}
	return nil
}
