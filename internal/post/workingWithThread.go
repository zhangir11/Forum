package post

import (
	"forum/pkg/sserver"
	"log"
	"regexp"
	"strings"
)

var rxAlphaNumeric = regexp.MustCompile(`^[a-zA-Z0-9]*$`)

// InsertThreadInfo ..
func InsertThreadInfo(threadName string, postID int) error {
	stmnt, err := sserver.Db.Db.Prepare("INSERT INTO Threads (Name) VALUES (?)")
	if err != nil {
		log.Println("Failed on CreatePost" + "insert Threads error")
		return err
	}
	defer stmnt.Close()
	var threadID int64

	res, err := stmnt.Exec(threadName)
	if err != nil {
		log.Println("Failed on CreatePost"+err.Error(), "---> exec Threads error")

		// if thread exists in DB, we will get his threadID, and insert to PostMapping Table
		if err.Error() == "UNIQUE constraint failed: Threads.Name" {
			threadID = int64(GetThreadID(threadName))
		} else {
			return err
		}
	} else {
		threadID, _ = res.LastInsertId()
	}

	InsertPostMapInfo(postID, int(threadID))
	return nil
}

// GetThreadID ...
func GetThreadID(name string) int {
	var id int
	if err := sserver.Db.Db.QueryRow("SELECT ID FROM Threads WHERE Name = ?", name).Scan(&id); err != nil {
		log.Println("Failed on CreatePost " + name + " threadID retrieve error" + err.Error())
	}

	return id
}

// GetPostIDsByTID ...
func GetPostIDsByTID(tid int) []int {
	postIDs := []int{}

	query, err := sserver.Db.Db.Query("SELECT postID FROM PostMapping WHERE threadID=? ORDER BY postID DESC", tid)
	if err != nil {
		log.Println("Failed on CreatePost"+"GetPostIDsByTID", err.Error())
		return nil
	}
	defer query.Close()

	for query.Next() {
		var id int
		if err := query.Scan(&id); err != nil {
			log.Println("Failed on CreatePost"+"GetPostIDsByTID", err.Error())
			return nil
		}
		postIDs = append(postIDs, id)
	}
	return postIDs
}

// CheckNumberOfThreads ... (when user create post)
func CheckNumberOfThreads(input string) []string {
	input = strings.ToLower(input)
	done := []string{}
	Threads := strings.Split(input, ",")
	for _, thread := range Threads {
		thread = strings.ReplaceAll(thread, " ", "")
		if thread == "" {
			continue
		}

		if rxAlphaNumeric.Match([]byte(thread)) {
			done = append(done, thread)
		}
	}
	return done
}

// InsertPostMapInfo ...
func InsertPostMapInfo(postID, threadID int) {
	stmnt, err := sserver.Db.Db.Prepare("INSERT INTO PostMapping (postID, threadID) VALUES (?, ?)")
	if err != nil {
		log.Println("Failed on CreatePost" + "postmap insert error" + err.Error())
		return
	}
	stmnt.Exec(postID, threadID)
	stmnt.Close()
}

// GetThreadOfPost ...
func GetThreadOfPost(postID int) ([]*Thread, error) {
	threadIDs := []int{}  // for getting all threadIDs from postID
	var threads []*Thread // for getting names of threads
	res, err := sserver.Db.Db.Query("SELECT threadID FROM PostMapping WHERE postID = ?", postID)
	if err != nil { // if err == sql.ErrNoRows ---> if no category in the post
		return nil, err
	}
	defer res.Close()
	/* Here we retrieve all threads relating with one single post*/
	for res.Next() {
		ID := 0
		if err := res.Scan(&ID); err != nil {
			log.Println("error func\"GetThreadOfPost()\"")
			return nil, err
		}
		threadIDs = append(threadIDs, ID)
	}
	// After getting IDs of threads, by these IDs we will get names of the threads
	for _, threadID := range threadIDs {
		thread, err := GetThreadByID(threadID) // exactly here
		if err != nil {
			return nil, err
		}
		threads = append(threads, thread)
	}

	return threads, nil
}
