package sqlite

import (
	"forumAPIs/pkg/models"
	"log"
)

func InsertComment(postID int, content string, userId int) error {
	_, err := DB.Exec("insert into comments (post_id, content, user_id) values ((select id from posts where id = ?), ?, ?)",
		postID, content, userId)
	return err
}

func GetCommentsByPostId(postID int) ([]*models.Comment, error) {
	rows, err := DB.Query("select id,content,user_id from comments where post_id = ?", postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var comments []*models.Comment
	for rows.Next() {
		c := &models.Comment{User: &models.User{}}
		err = rows.Scan(&c.Id, &c.Content, &c.User.Id)
		if err != nil {
			return nil, err
		}
		c.User, err = GetUserForPostInfo(c.User.Id)
		if err != nil {
			return nil, err
		}
		c.LikeUsers, err = GetCommentLikes(c.Id)
		if err != nil {
			return nil, err
		}
		c.LikeCount = len(c.LikeUsers)
		c.DislikeUsers, err = GetCommentDislikes(c.Id)
		if err != nil {
			return nil, err
		}
		c.DislikeCount = len(c.DislikeUsers)
		comments = append(comments, c)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return comments, nil
}

func GetCommentLikes(commentId int) ([]int, error) {
	rows, err := DB.Query("select user_id from comment_likes where comment_id = ?", commentId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var likes []int
	for rows.Next() {
		userId := -1
		err = rows.Scan(&userId)
		if err != nil {
			log.Println(err)
		}
		likes = append(likes, userId)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return likes, nil
}

func GetCommentDislikes(commentId int) ([]int, error) {
	rows, err := DB.Query("select user_id from comment_dislikes where comment_id = ?", commentId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var dislikes []int
	for rows.Next() {
		userId := -1
		err = rows.Scan(&userId)
		if err != nil {
			log.Println(err)
		}
		dislikes = append(dislikes, userId)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return dislikes, nil
}
