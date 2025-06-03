package data

import (
	"database/sql"

	"log"
	"ytdlp-bot/internal/environment"

	_ "github.com/mattn/go-sqlite3"
)

type Message struct {
	ID   int
	Date string
}

func GetAllMessagesId(clearKey bool) []int {
	db, err := sql.Open("sqlite3", environment.Environment.DBFilePath)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	if err != nil {
		log.Panic(err)
	}
	tx, err := db.Begin()
	defer tx.Rollback()

	rows, err := db.Query("SELECT * FROM MESSAGES")
	if err != nil {
		log.Print("История сообщений пуста")
	}
	messageIDS := []int{}
	for rows.Next() {
		message := Message{}
		err = rows.Scan(&message.ID, &message.Date)
		if err == nil {
			messageIDS = append(messageIDS, message.ID)
		}
	}
	if clearKey {
		db.Exec("DELETE FROM MESSAGES")
	}
	tx.Commit()
	return messageIDS
}

func PutMessageToDB(message Message) {

	db, err := sql.Open("sqlite3", environment.Environment.DBFilePath)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()
	_, err = db.Exec("insert into MESSAGES (ID, DATE_TIME_OF_CREATION) values ($1, $2)", message.ID, message.Date)
	if err != nil {
		log.Printf("Попытка поторно добавить сообщение с id=%d в базу данных", message.ID)
	}
}

func GetAllMessagesIdExcludeInput(excludeMessageId int) []Message {
	db, err := sql.Open("sqlite3", environment.Environment.DBFilePath)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	if err != nil {
		log.Panic(err)
	}

	rows, err := db.Query("SELECT * FROM MESSAGES WHERE ID!=$1", excludeMessageId)
	if err != nil {
		log.Print("История сообщений пуста")
	}
	messages := []Message{}
	for rows.Next() {
		message := Message{}
		err = rows.Scan(&message.ID, &message.Date)
		if err == nil {
			messages = append(messages, message)
		}
	}

	return messages
}

func DeleteMessageWithId(messageId int) {
	db, err := sql.Open("sqlite3", environment.Environment.DBFilePath)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()
	db.Exec("DELETE FROM MESSAGES WHERE ID=$1", messageId)
}
