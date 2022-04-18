package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"forumAPIs/pkg/models"
	"sort"
)

func GetPosts() ([]*models.Post, error) {
	rows, err := DB.Query("select id, title, contents, create_date, user_id from Posts order by id desc")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*models.Post

	for rows.Next() {
		post := &models.Post{}
		post.Categories = []*models.Category{}
		post.User = &models.User{}
		err = rows.Scan(&post.Id, &post.Title, &post.Content, &post.CreateDate, &post.User.Id)
		if err != nil {
			return nil, err
		}
		post.User, err = GetUserForPostInfo(post.User.Id)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		post.Categories, err = GetCategoriesByPost(post.Id)
		if err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}
	if err != nil {
		return nil, err
	}
	return posts, nil
}

func GetPostById(id int) (*models.Post, error) {
	post := &models.Post{}
	post.Categories = []*models.Category{}
	post.User = &models.User{}

	row := DB.QueryRow("select id, title, contents, create_date, user_id from Posts where id = ?", id)
	err := row.Scan(&post.Id, &post.Title, &post.Content, &post.CreateDate, &post.User.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		}
		return nil, err
	}

	post.Categories, err = GetCategoriesByPost(post.Id)
	if err != nil {
		return nil, err
	}

	post.User, err = GetUserForPostInfo(post.User.Id)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	likes, err := GetPostLikeUsers(post.Id)
	if err != nil {
		return nil, err
	}
	dislikes, err := GetPostDislikeUsers(post.Id)
	if err != nil {
		return nil, err
	}
	for _, userId := range likes {
		post.Likes = append(post.Likes, &models.User{Id: userId})
	}
	for _, userId := range dislikes {
		post.Dislikes = append(post.Dislikes, &models.User{Id: userId})
	}
	post.Comments, err = GetCommentsByPostId(post.Id)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func GetPostsByCategory(catID int) ([]*models.Post, error) {
	rows, err := DB.Query("select post_id from posts_categories where category_id =  ? ", catID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*models.Post

	for rows.Next() {
		post := &models.Post{}

		err := rows.Scan(&post.Id)
		if err != nil {
			return nil, err
		}

		post, err = GetPostById(post.Id)
		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Id > posts[j].Id
	})

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func InsertPost(title, contents string, categories []string, userId int) (int, error) {
	result, err := DB.Exec("insert into posts (title, contents, create_date, user_id) values (?, ?, DATETIME('now', 'localtime'), ?)", title, contents, userId)
	if err != nil {
		return -1, err
	}
	postId, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}

	for _, catName := range categories {
		row := DB.QueryRow("select id from categories where name = ?", catName)
		var catId int
		err := row.Scan(&catId)

		_, err = DB.Exec("insert into posts_categories (post_id, category_id) values (?,?);", postId, catId)
		if err != nil {
			return -1, err
		}

		if err != nil {
			return -1, err
		}
	}

	return int(postId), nil
}
