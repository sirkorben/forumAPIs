package sqlite

import "forumAPIs/pkg/models"

func GetAllCategories() ([]*models.Category, error) {
	rows, err := DB.Query("select id, name from categories order by id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*models.Category
	for rows.Next() {
		c := &models.Category{}

		err = rows.Scan(&c.Id, &c.Name)

		if err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return categories, nil
}

func GetCategory(catId int) (*models.Category, error) {
	c := &models.Category{Id: catId}
	row := DB.QueryRow("select name from categories where id = ?", catId)
	err := row.Scan(&c.Name)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func GetCategoriesByPost(postID int) ([]*models.Category, error) {
	rows, err := DB.Query("select category_id from posts_categories where post_id = ?", postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []*models.Category
	for rows.Next() {
		c := &models.Category{}
		err := rows.Scan(&c.Id)
		if err != nil {
			return nil, err
		}

		c, err = GetCategory(c.Id)
		if err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return categories, nil
}
