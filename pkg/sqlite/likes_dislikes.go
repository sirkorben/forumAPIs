package sqlite

import (
	"database/sql"
	"errors"
)

func GetPostLikeUsers(postID int) ([]int, error) {
	rows, err := DB.Query("select user_id from likes where post_id = ?", postID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var likes []int
	for rows.Next() {
		userId := -1
		err = rows.Scan(&userId)
		if err != nil {
			return nil, err
		}
		likes = append(likes, userId)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return likes, nil
}

func GetPostDislikeUsers(postID int) ([]int, error) {
	rows, err := DB.Query("select user_id from dislikes where post_id = ?", postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var dislikes []int

	for rows.Next() {
		userId := -1
		err = rows.Scan(&userId)
		if err != nil {
			return nil, err
		}
		dislikes = append(dislikes, userId)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return dislikes, nil
}

func ChangePostLike(postId int, userId int) error {
	row := DB.QueryRow("select user_id from likes where (post_id, user_id) = (?,?)", postId, userId)

	id := -1
	err := row.Scan(&id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			_, err := DB.Exec("delete from dislikes where (post_id,user_id) = (?,?)", postId, userId)
			if err != nil && err != sql.ErrNoRows {
				return err
			}

			_, err = DB.Exec("insert into likes (post_id, user_id) values (?,?)", postId, userId)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		_, err := DB.Exec("delete from likes where (post_id,user_id) = (?,?)", postId, userId)
		if err != nil {
			return err
		}
	}
	return nil
}

func ChangePostDislike(postId int, userId int) error {
	row := DB.QueryRow("select user_id from dislikes where (post_id, user_id) = (?,?)", postId, userId)

	id := -1
	err := row.Scan(&id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			_, err := DB.Exec("delete from likes where (post_id,user_id) = (?,?)", postId, userId)
			if err != nil && err != sql.ErrNoRows {
				return err
			}

			_, err = DB.Exec("insert into dislikes (post_id, user_id) values (?,?)", postId, userId)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		_, err := DB.Exec("delete from dislikes where (post_id,user_id) = (?,?)", postId, userId)
		if err != nil {
			return err
		}
	}
	return nil
}
