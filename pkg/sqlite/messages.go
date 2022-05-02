package sqlite

import (
	"forumAPIs/pkg/models"
	"log"
	"math"
	"strconv"
)

func GetAllChats(userId int) ([]*models.Message, error) {
	rows, err := DB.Query("SELECT content, user_id, creation_date from messages_" + strconv.Itoa(userId) + " ORDER BY creation_date desc")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []*models.Message
	var tempUserId float64 // tempUserId = 0
	checkedUserIds := []float64{}
	for rows.Next() {
		message := &models.Message{}
		err = rows.Scan(&message.Content, &message.ToFromUser, &message.CreationDate)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		if math.Abs(float64(message.ToFromUser)) != math.Abs(float64(tempUserId)) {
			if !contains(checkedUserIds, math.Abs(float64(message.ToFromUser))) {
				checkedUserIds = append(checkedUserIds, math.Abs(float64(message.ToFromUser)))
				messages = append(messages, message)
			}
			tempUserId = math.Abs(float64(message.ToFromUser))
		}
	}
	return messages, nil
}

func contains(checked []float64, userId float64) bool {
	for _, value := range checked {
		if value == userId {
			return true
		}
	}
	return false
}

func GetChatWithUser(userId, personToChatId int) ([]*models.Message, error) {
	var messages []*models.Message
	rows, err := DB.Query("select content, user_id, creation_date from messages_"+strconv.Itoa(userId)+" where (user_id = ? or user_id = - ?) order by id desc", personToChatId, personToChatId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		message := &models.Message{}
		err = rows.Scan(&message.Content, &message.ToFromUser, &message.CreationDate)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}
	return messages, nil
}

func SendMessage(content string, senderId, receiverId int) error {
	// TODO: handle possible error from quering for existing tables tableExists()
	var err error
	//sender part
	if tableExists(senderId) {
		err = insertMessage(senderId, receiverId, content)
		if err != nil {
			return models.InternalServerError
		}
	} else {
		err = createMessageTable(senderId)
		if err != nil {
			return models.InternalServerError
		}
		err = insertMessage(senderId, receiverId, content)
		if err != nil {
			return models.InternalServerError
		}
	}

	// recevier part
	if tableExists(receiverId) {
		insertMessage(receiverId, -(senderId), content)
	} else {
		err = createMessageTable(receiverId)
		if err != nil {
			return models.InternalServerError
		}
		err = insertMessage(receiverId, -(senderId), content)
		if err != nil {
			return models.InternalServerError
		}
	}
	return nil
}

func insertMessage(messageTableHolder, toFromUserId int, content string) error {
	_, err := DB.Exec("insert into messages_"+strconv.Itoa(messageTableHolder)+" (user_id, content, creation_date) values (?,?, strftime('%s','now'))",
		toFromUserId, content)
	if err != nil {
		log.Println("error with inserting message into DB")
		return err
	}
	return nil
}

func tableExists(userId int) bool {
	_, table_check := DB.Query("select * from " + "messages_" + strconv.Itoa(userId) + ";")
	if table_check == nil { // if no error, table exists
		return true
	} else {
		return false
	}
}

func createMessageTable(userId int) error {
	var create = `create table messages_` + strconv.Itoa(userId) + `(
		id INTEGER not null primary key autoincrement,
		user_id INTEGER not null,
		content TEXT not null,
		creation_date INTEGER not null);
		`
	_, err := DB.Exec(create)
	if err != nil {
		log.Println("error with creating message table into DB")
		return err
	}
	if err = DB.Ping(); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
