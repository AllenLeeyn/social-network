package models

import (
	"database/sql"
	"fmt"
	"log"
	userManagementModels "social-network/pkg/userManagement/models"
	"social-network/pkg/utils"
	"sort"
	"time"
)

// Post struct represents the user data model
type Post struct {
	ID           int           `json:"id"`
	UUID         string        `json:"uuid"`
	UserId       int           `json:"user_id"`
	GroupId      sql.NullInt64 `json:"group_id"`
	Title        string        `json:"title"`
	Content      string        `json:"content"`
	Visibility   string        `json:"visibility"`
	LikeCount    int           `json:"like_count"`
	DisikeCount  int           `json:"dilike_count"`
	CommentCount int           `json:"comment_count"`
	Status       string        `json:"status"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    *time.Time    `json:"updated_at"`
	UpdatedBy    *int          `json:"updated_by"`

	IsLikedByUser    bool                      `json:"liked"`
	IsDislikedByUser bool                      `json:"disliked"`
	User             userManagementModels.User `json:"user"` // Embedded user data

	Categories []Category `json:"categories"` // List of categories related to the post
	PostFiles  []PostFile `json:"post_files"` // List of files related to the post
}

func InsertPost(post *Post, categoryIds []int, uploadedFiles map[string]string) (int, error) {
	// Start a transaction for atomicity
	tx, err := sqlDB.Begin()
	if err != nil {
		return -1, err
	}

	post.UUID, err = utils.GenerateUuid()
	if err != nil {
		tx.Rollback() // Rollback if UUID generation fails
		return -1, err
	}

	insertQuery := `INSERT INTO posts (uuid, title, content, user_id, group_id, visibility) VALUES (?, ?, ?, ?, ?, ?);`
	result, insertErr := tx.Exec(insertQuery, post.UUID, post.Title, post.Content, post.UserId, post.GroupId, post.Visibility)
	if insertErr != nil {
		return -1, insertErr
	}

	// Retrieve the last inserted ID
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		tx.Rollback() // Rollback on error
		log.Fatal(err)
		return -1, err
	}

	insertPostCategoriesErr := InsertPostCategories(int(lastInsertID), categoryIds, post.UserId, tx)
	if insertPostCategoriesErr != nil {
		tx.Rollback() // Rollback on error
		return -1, insertPostCategoriesErr
	}

	insertPostFilesErr := InsertPostFiles(int(lastInsertID), uploadedFiles, post.UserId, tx)
	if insertPostFilesErr != nil {
		tx.Rollback() // Rollback on error
		return -1, insertPostFilesErr
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		tx.Rollback() // Rollback on error
		return -1, err
	}

	return int(lastInsertID), nil
}

func UpdatePost(post *Post, categories []int, uploadedFiles map[string]string, user_id int) error {
	// Start a transaction for atomicity
	tx, err := sqlDB.Begin()
	if err != nil {
		return err
	}

	updateQuery := `UPDATE posts
					SET title = ?,
						content = ?,
						visibility = ?,
						updated_at = CURRENT_TIMESTAMP,
						updated_by = ?
					WHERE id = ?;`
	_, updateErr := tx.Exec(updateQuery, post.Title, post.Content, post.Visibility, user_id, post.ID)
	if updateErr != nil {
		return updateErr
	}

	deletePostCategoriesErr := UpdateStatusPostCategories(post.ID, user_id, "delete", tx)
	if deletePostCategoriesErr != nil {
		tx.Rollback() // Rollback on error
		return deletePostCategoriesErr
	}

	if len(uploadedFiles) != 0 {
		deletePostFilesErr := UpdateStatusPostFiles(post.ID, user_id, "delete", tx)
		if deletePostFilesErr != nil {
			tx.Rollback() // Rollback on error
			return deletePostFilesErr
		}
	}

	insertPostCategoriesErr := InsertPostCategories(post.ID, categories, user_id, tx)
	if insertPostCategoriesErr != nil {
		tx.Rollback() // Rollback on error
		return insertPostCategoriesErr
	}

	insertPostFilesErr := InsertPostFiles(post.ID, uploadedFiles, user_id, tx)
	if insertPostFilesErr != nil {
		tx.Rollback() // Rollback on error
		return insertPostFilesErr
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		tx.Rollback() // Rollback on error
		return err
	}

	return nil
}

func UpdateStatusPost(post_id int, status string, user_id int) error {
	// Start a transaction for atomicity
	tx, err := sqlDB.Begin()
	if err != nil {
		return err
	}

	updateQuery := `UPDATE posts
					SET status = ?,
						updated_at = CURRENT_TIMESTAMP,
						updated_by = ?
					WHERE id = ?;`
	_, updateErr := tx.Exec(updateQuery, status, user_id, post_id)
	if updateErr != nil {
		return updateErr
	}

	updateStatusPostCategories := UpdateStatusPostCategories(post_id, user_id, status, tx)
	if updateStatusPostCategories != nil {
		tx.Rollback() // Rollback on error
		return updateStatusPostCategories
	}

	UpdateStatusPostFiles := UpdateStatusPostFiles(post_id, user_id, status, tx)
	if UpdateStatusPostFiles != nil {
		tx.Rollback() // Rollback on error
		return UpdateStatusPostFiles
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		tx.Rollback() // Rollback on error
		return err
	}

	return nil
}

func ReadAllPosts(checkLikeForUser int) ([]Post, error) {
	// Query the records
	// todo check
	rows, selectError := sqlDB.Query(`
        SELECT p.id as post_id, p.uuid as post_uuid, p.title as post_title, p.content as post_content, p.visibility as post_visibility, p.status as post_status, p.created_at as post_created_at, p.updated_at as post_updated_at, p.updated_by as post_updated_by,
			p.like_count as post_like_count, p.dislike_count as post_dislike_count, p.comment_count as post_comment_count,
			u.id as user_id, u.first_name as user_first_name, u.last_name as user_last_name, u.nick_name as user_nick_name, u.email as user_email, IFNULL(u.profile_image, '') as profile_image,
			c.id as category_id, c.name as category_name,
			IFNULL(pf.id, 0) as post_file_id, pf.file_uploaded_name, pf.file_real_name,
			CASE 
                WHEN EXISTS (SELECT 1 FROM post_likes WHERE post_id = p.id AND status != 'delete' AND type = 'like' AND user_id = ?) THEN 1
                ELSE 0
            END AS is_liked_by_user,
            CASE 
                WHEN EXISTS (SELECT 1 FROM post_likes WHERE post_id = p.id AND status != 'delete' AND type = 'dislike' AND user_id = ?) THEN 1
                ELSE 0
            END AS is_disliked_by_user
		FROM posts p
			INNER JOIN users u
				ON p.user_id = u.id
			LEFT JOIN post_files pf
				ON p.id = pf.post_id
				AND pf.status = 'enable'
			LEFT JOIN post_categories pc
				ON p.id = pc.post_id
				AND pc.status = 'enable'
			LEFT JOIN categories c
				ON pc.category_id = c.id
				AND c.status = 'enable'
		WHERE p.status != 'delete'
			AND u.status != 'delete'
		ORDER BY p.id desc;
    `, checkLikeForUser, checkLikeForUser)
	if selectError != nil {
		return nil, selectError
	}
	defer rows.Close()

	var posts []Post
	// Map to track posts by their ID to avoid duplicates
	postMap := make(map[int]*Post)

	for rows.Next() {
		var post Post
		var user userManagementModels.User
		var category Category
		var postFile PostFile

		// Scan the post and user data
		err := rows.Scan(
			&post.ID, &post.UUID, &post.Title, &post.Content, &post.Visibility, &post.Status,
			&post.CreatedAt, &post.UpdatedAt, &post.UpdatedBy,
			&post.LikeCount, &post.DisikeCount, &post.CommentCount,
			&post.UserId, &user.FirstName, &user.LastName, &user.NickName, &user.Email, &user.ProfileImage,
			&category.ID, &category.Name,
			&postFile.ID, &postFile.FileUploadedName, &postFile.FileRealName,
			&post.IsLikedByUser, &post.IsDislikedByUser,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		// Check if the post already exists in the map
		existingPost, found := postMap[post.ID]
		if !found {
			post.User = user
			post.Categories = []Category{}
			post.PostFiles = []PostFile{}
			postMap[post.ID] = &post
			existingPost = &post
		}

		// Ensure unique categories
		isCategoryAdded := false
		for _, c := range existingPost.Categories {
			if c.ID == category.ID {
				isCategoryAdded = true
				break
			}
		}
		if !isCategoryAdded && category.ID != 0 {
			existingPost.Categories = append(existingPost.Categories, category)
		}

		// Ensure unique post files
		isFileAdded := false
		for _, f := range existingPost.PostFiles {
			if f.ID == postFile.ID {
				isFileAdded = true
				break
			}
		}
		if !isFileAdded && postFile.ID != 0 {
			existingPost.PostFiles = append(existingPost.PostFiles, postFile)
		}
	}

	// Check for any errors during row iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %v", err)
	}

	// Convert the map of posts into a slice
	for _, post := range postMap {
		posts = append(posts, *post)
	}

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].ID > posts[j].ID
	})

	return posts, nil
}

func ReadPostsByCategoryId(category_id int) ([]Post, error) {
	// Query the records
	rows, selectError := sqlDB.Query(`
        SELECT p.id as post_id, p.uuid as post_uuid, p.title as post_title, p.content as post_content, p.status as post_status, p.created_at as post_created_at, p.updated_at as post_updated_at, p.updated_by as post_updated_by,
			p.like_count as post_like_count, p.dislike_count as post_dislike_count, p.comment_count as post_comment_count,
			u.id as user_id, u.first_name as user_first_name, u.last_name as user_last_name, u.nick_name as user_nick_name, u.email as user_email, IFNULL(u.profile_image, '') as profile_image,
			c.id as category_id, c.name as category_name,
			IFNULL(pf.id, 0) as post_file_id, pf.file_uploaded_name, pf.file_real_name
		FROM posts p
			INNER JOIN users u
				ON p.user_id = u.id
			LEFT JOIN post_files pf
				ON p.id = pf.post_id
				AND pf.status = 'enable'
			INNER JOIN post_categories filterd_pc
				ON p.id = filterd_pc.post_id
				AND filterd_pc.status = 'enable'
				AND filterd_pc.category_id = ?
			INNER JOIN post_categories pc
				ON p.id = pc.post_id
				AND pc.status = 'enable'
			INNER JOIN categories c
				ON pc.category_id = c.id
				AND c.status = 'enable'
		WHERE p.status != 'delete'
			AND u.status != 'delete'
		ORDER BY p.id desc;
    `, category_id)
	if selectError != nil {
		return nil, selectError
	}
	defer rows.Close()

	var posts []Post
	// Map to track posts by their ID to avoid duplicates
	postMap := make(map[int]*Post)

	for rows.Next() {
		var post Post
		var user userManagementModels.User
		var category Category
		var postFile PostFile

		// Scan the post and user data
		err := rows.Scan(
			&post.ID, &post.UUID, &post.Title, &post.Content, &post.Status,
			&post.LikeCount, &post.DisikeCount, &post.CommentCount,
			&post.CreatedAt, &post.UpdatedAt, &post.UpdatedBy, &post.UserId,
			&user.FirstName, &user.LastName, &user.NickName, &user.Email, &user.ProfileImage,
			&category.ID, &category.Name,
			&postFile.ID, &postFile.FileUploadedName, &postFile.FileRealName,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		// Check if the post already exists in the map
		existingPost, found := postMap[post.ID]
		if !found {
			post.User = user
			post.Categories = []Category{}
			post.PostFiles = []PostFile{}
			postMap[post.ID] = &post
			existingPost = &post
		}

		// Ensure unique categories
		isCategoryAdded := false
		for _, c := range existingPost.Categories {
			if c.ID == category.ID {
				isCategoryAdded = true
				break
			}
		}
		if !isCategoryAdded && category.ID != 0 {
			existingPost.Categories = append(existingPost.Categories, category)
		}

		// Ensure unique post files
		isFileAdded := false
		for _, f := range existingPost.PostFiles {
			if f.ID == postFile.ID {
				isFileAdded = true
				break
			}
		}
		if !isFileAdded && postFile.ID != 0 {
			existingPost.PostFiles = append(existingPost.PostFiles, postFile)
		}
	}

	// Check for any errors during row iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %v", err)
	}

	// Convert the map of posts into a slice
	for _, post := range postMap {
		posts = append(posts, *post)
	}

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].ID > posts[j].ID
	})

	return posts, nil
}

func FilterPosts(searchTerm string) ([]Post, error) {
	searchPattern := "%" + searchTerm + "%" // Add wildcards for LIKE comparison

	// Query the records
	rows, selectError := sqlDB.Query(`
        SELECT p.id as post_id, p.uuid as post_uuid, p.title as post_title, p.content as post_content, p.status as post_status, p.created_at as post_created_at, p.updated_at as post_updated_at, p.updated_by as post_updated_by,
			p.like_count as post_like_count, p.dislike_count as post_dislike_count, p.comment_count as post_comment_count,
			u.id as user_id, u.first_name as user_first_name, u.last_name as user_last_name, u.nick_name as user_nick_name, u.email as user_email, IFNULL(u.profile_image, '') as profile_image,
			c.id as category_id, c.name as category_name,
			IFNULL(pf.id, 0) as post_file_id, pf.file_uploaded_name, pf.file_real_name
		FROM posts p
			INNER JOIN users u
				ON p.user_id = u.id
			LEFT JOIN post_files pf
				ON p.id = pf.post_id
				AND pf.status = 'enable'
			LEFT JOIN post_categories pc
				ON p.id = pc.post_id
				AND pc.status = 'enable'
			LEFT JOIN categories c
				ON pc.category_id = c.id
				AND c.status = 'enable'
		WHERE p.status != 'delete'
			AND u.status != 'delete'
      		AND (p.title LIKE ? OR p.content LIKE ?)
		ORDER BY p.id desc;
    `, searchPattern, searchPattern)
	if selectError != nil {
		return nil, selectError
	}
	defer rows.Close()

	var posts []Post
	// Map to track posts by their ID to avoid duplicates
	postMap := make(map[int]*Post)

	for rows.Next() {
		var post Post
		var user userManagementModels.User
		var category Category
		var postFile PostFile

		// Scan the post and user data
		err := rows.Scan(
			&post.ID, &post.UUID, &post.Title, &post.Content, &post.Status,
			&post.LikeCount, &post.DisikeCount, &post.CommentCount,
			&post.CreatedAt, &post.UpdatedAt, &post.UpdatedBy, &post.UserId,
			&user.FirstName, &user.LastName, &user.NickName, &user.Email, &user.ProfileImage,
			&category.ID, &category.Name,
			&postFile.ID, &postFile.FileUploadedName, &postFile.FileRealName,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		// Check if the post already exists in the map
		existingPost, found := postMap[post.ID]
		if !found {
			post.User = user
			post.Categories = []Category{}
			post.PostFiles = []PostFile{}
			postMap[post.ID] = &post
			existingPost = &post
		}

		// Ensure unique categories
		isCategoryAdded := false
		for _, c := range existingPost.Categories {
			if c.ID == category.ID {
				isCategoryAdded = true
				break
			}
		}
		if !isCategoryAdded && category.ID != 0 {
			existingPost.Categories = append(existingPost.Categories, category)
		}

		// Ensure unique post files
		isFileAdded := false
		for _, f := range existingPost.PostFiles {
			if f.ID == postFile.ID {
				isFileAdded = true
				break
			}
		}
		if !isFileAdded && postFile.ID != 0 {
			existingPost.PostFiles = append(existingPost.PostFiles, postFile)
		}
	}

	// Check for any errors during row iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %v", err)
	}

	// Convert the map of posts into a slice
	for _, post := range postMap {
		posts = append(posts, *post)
	}

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].ID > posts[j].ID
	})

	return posts, nil
}

func ReadPostsByUserId(userId int) ([]Post, error) {
	// Query the records
	rows, selectError := sqlDB.Query(`
        SELECT p.id as post_id, p.uuid as post_uuid, p.title as post_title, p.content as post_content, p.status as post_status, p.created_at as post_created_at, p.updated_at as post_updated_at, p.updated_by as post_updated_by,
			p.like_count as post_like_count, p.dislike_count as post_dislike_count, p.comment_count as post_comment_count,
			u.id as user_id, u.first_name as user_first_name, u.last_name as user_last_name, u.nick_name as user_nick_name, u.email as user_email, IFNULL(u.profile_image, '') as profile_image,
			c.id as category_id, c.name as category_name,
			IFNULL(pf.id, 0) as post_file_id, pf.file_uploaded_name, pf.file_real_name,
			CASE 
                WHEN EXISTS (SELECT 1 FROM post_likes WHERE post_id = p.id AND status != 'delete' AND type = 'like' AND user_id = ?) THEN 1
                ELSE 0
            END AS is_liked_by_user,
            CASE 
                WHEN EXISTS (SELECT 1 FROM post_likes WHERE post_id = p.id AND status != 'delete' AND type = 'dislike' AND user_id = ?) THEN 1
                ELSE 0
            END AS is_disliked_by_user
		FROM posts p
			INNER JOIN users u
				ON p.user_id = u.id
				AND u.id = ?
			LEFT JOIN post_files pf
				ON p.id = pf.post_id
				AND pf.status = 'enable'
			LEFT JOIN post_categories pc
				ON p.id = pc.post_id
				AND pc.status = 'enable'
			LEFT JOIN categories c
				ON pc.category_id = c.id
				AND c.status = 'enable'
		WHERE p.status != 'delete'
			AND u.status != 'delete'
		ORDER BY p.id desc;
    `, userId, userId, userId)
	if selectError != nil {
		return nil, selectError
	}
	defer rows.Close()

	var posts []Post
	// Map to track posts by their ID to avoid duplicates
	postMap := make(map[int]*Post)

	for rows.Next() {
		var post Post
		var user userManagementModels.User
		var category Category
		var postFile PostFile

		// Scan the post and user data
		err := rows.Scan(
			&post.ID, &post.UUID, &post.Title, &post.Content, &post.Status,
			&post.CreatedAt, &post.UpdatedAt, &post.UpdatedBy,
			&post.LikeCount, &post.DisikeCount, &post.CommentCount,
			&post.UserId, &user.FirstName, &user.LastName, &user.NickName, &user.Email, &user.ProfileImage,
			&category.ID, &category.Name,
			&postFile.ID, &postFile.FileUploadedName, &postFile.FileRealName,
			&post.IsLikedByUser, &post.IsDislikedByUser,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		// Check if the post already exists in the map
		existingPost, found := postMap[post.ID]
		if !found {
			post.User = user
			post.Categories = []Category{}
			post.PostFiles = []PostFile{}
			postMap[post.ID] = &post
			existingPost = &post
		}

		// Ensure unique categories
		isCategoryAdded := false
		for _, c := range existingPost.Categories {
			if c.ID == category.ID {
				isCategoryAdded = true
				break
			}
		}
		if !isCategoryAdded && category.ID != 0 {
			existingPost.Categories = append(existingPost.Categories, category)
		}

		// Ensure unique post files
		isFileAdded := false
		for _, f := range existingPost.PostFiles {
			if f.ID == postFile.ID {
				isFileAdded = true
				break
			}
		}
		if !isFileAdded && postFile.ID != 0 {
			existingPost.PostFiles = append(existingPost.PostFiles, postFile)
		}
	}

	// Check for any errors during row iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %v", err)
	}

	// Convert the map of posts into a slice
	for _, post := range postMap {
		posts = append(posts, *post)
	}

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].ID > posts[j].ID
	})

	return posts, nil
}

func ReadPostsLikedByUserId(userId int) ([]Post, error) {
	// Query the records
	rows, selectError := sqlDB.Query(`
        SELECT p.id as post_id, p.uuid as post_uuid, p.title as post_title, p.content as post_content, p.status as post_status, p.created_at as post_created_at, p.updated_at as post_updated_at, p.updated_by as post_updated_by,
			p.like_count as post_like_count, p.dislike_count as post_dislike_count, p.comment_count as post_comment_count,
			u.id as user_id, u.first_name as user_first_name, u.last_name as user_last_name, u.nick_name as user_nick_name, u.email as user_email, IFNULL(u.profile_image, '') as profile_image,
			c.id as category_id, c.name as category_name,
			IFNULL(pf.id, 0) as post_file_id, pf.file_uploaded_name, pf.file_real_name,
			CASE 
                WHEN EXISTS (SELECT 1 FROM post_likes WHERE post_id = p.id AND status != 'delete' AND type = 'like' AND user_id = ?) THEN 1
                ELSE 0
            END AS is_liked_by_user,
            CASE 
                WHEN EXISTS (SELECT 1 FROM post_likes WHERE post_id = p.id AND status != 'delete' AND type = 'dislike' AND user_id = ?) THEN 1
                ELSE 0
            END AS is_disliked_by_user
		FROM posts p
			INNER JOIN post_likes pl
				ON pl.post_id = p.id
				AND pl.status = 'enable'
			INNER JOIN users u
				ON p.user_id = u.id
			LEFT JOIN post_files pf
				ON p.id = pf.post_id
				AND pf.status = 'enable'
			INNER JOIN users liked_user
				ON pl.user_id = liked_user.id
				AND liked_user.id = ?
			LEFT JOIN post_categories pc
				ON p.id = pc.post_id
				AND pc.status = 'enable'
			LEFT JOIN categories c
				ON pc.category_id = c.id
				AND c.status = 'enable'
		WHERE p.status != 'delete'
			AND u.status != 'delete'
		ORDER BY p.id desc;
    `, userId, userId, userId)
	if selectError != nil {
		return nil, selectError
	}
	defer rows.Close()

	var posts []Post
	// Map to track posts by their ID to avoid duplicates
	postMap := make(map[int]*Post)

	for rows.Next() {
		var post Post
		var user userManagementModels.User
		var category Category
		var postFile PostFile

		// Scan the post and user data
		err := rows.Scan(
			&post.ID, &post.UUID, &post.Title, &post.Content, &post.Status,
			&post.CreatedAt, &post.UpdatedAt, &post.UpdatedBy,
			&post.LikeCount, &post.DisikeCount, &post.CommentCount,
			&post.UserId, &user.FirstName, &user.LastName, &user.NickName, &user.Email, &user.ProfileImage,
			&category.ID, &category.Name,
			&postFile.ID, &postFile.FileUploadedName, &postFile.FileRealName,
			&post.IsLikedByUser, &post.IsDislikedByUser,
		)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %v", err)
		}

		// Check if the post already exists in the map
		existingPost, found := postMap[post.ID]
		if !found {
			post.User = user
			post.Categories = []Category{}
			post.PostFiles = []PostFile{}
			postMap[post.ID] = &post
			existingPost = &post
		}

		// Ensure unique categories
		isCategoryAdded := false
		for _, c := range existingPost.Categories {
			if c.ID == category.ID {
				isCategoryAdded = true
				break
			}
		}
		if !isCategoryAdded && category.ID != 0 {
			existingPost.Categories = append(existingPost.Categories, category)
		}

		// Ensure unique post files
		isFileAdded := false
		for _, f := range existingPost.PostFiles {
			if f.ID == postFile.ID {
				isFileAdded = true
				break
			}
		}
		if !isFileAdded && postFile.ID != 0 {
			existingPost.PostFiles = append(existingPost.PostFiles, postFile)
		}
	}

	// Check for any errors during row iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %v", err)
	}

	// Convert the map of posts into a slice
	for _, post := range postMap {
		posts = append(posts, *post)
	}

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].ID > posts[j].ID
	})

	return posts, nil
}

func ReadPostById(postId int, checkLikeForUser int) (Post, error) {
	// Query the records
	rows, selectError := sqlDB.Query(`
        SELECT p.id as post_id, p.uuid as post_uuid, p.title as post_title, p.content as post_content, p.status as post_status, p.created_at as post_created_at, p.updated_at as post_updated_at, p.updated_by as post_updated_by,
			p.like_count as post_like_count, p.dislike_count as post_dislike_count, p.comment_count as post_comment_count,
			u.id as user_id, u.first_name as user_first_name, u.last_name as user_last_name, u.nick_name as user_nick_name, u.email as user_email, IFNULL(u.profile_image, '') as profile_image,
			c.id as category_id, c.name as category_name,
			IFNULL(pf.id, 0) as post_file_id, pf.file_uploaded_name, pf.file_real_name,
			CASE 
                WHEN EXISTS (SELECT 1 FROM post_likes WHERE post_id = p.id AND status != 'delete' AND type = 'like' AND user_id = ?) THEN 1
                ELSE 0
            END AS is_liked_by_user,
            CASE 
                WHEN EXISTS (SELECT 1 FROM post_likes WHERE post_id = p.id AND status != 'delete' AND type = 'dislike' AND user_id = ?) THEN 1
                ELSE 0
            END AS is_disliked_by_user
		FROM posts p
			INNER JOIN users u
				ON p.user_id = u.id
				AND p.id = ?
			LEFT JOIN post_files pf
				ON p.id = pf.post_id
				AND pf.status = 'enable'
			LEFT JOIN post_categories pc
				ON p.id = pc.post_id
				AND pc.status = 'enable'
			LEFT JOIN categories c
				ON pc.category_id = c.id
				AND c.status = 'enable'
		WHERE p.status != 'delete'
			AND u.status != 'delete';
    `, checkLikeForUser, checkLikeForUser, postId)
	if selectError != nil {
		return Post{}, selectError
	}
	defer rows.Close()

	var post Post
	var user userManagementModels.User

	// Scan the records
	for rows.Next() {
		var category Category
		var postFile PostFile

		err := rows.Scan(
			&post.ID, &post.UUID, &post.Title, &post.Content, &post.Status,
			&post.CreatedAt, &post.UpdatedAt, &post.UpdatedBy,
			&post.LikeCount, &post.DisikeCount, &post.CommentCount,
			&post.UserId, &user.FirstName, &user.LastName, &user.NickName, &user.Email, &user.ProfileImage,
			&category.ID, &category.Name,
			&postFile.ID, &postFile.FileUploadedName, &postFile.FileRealName,
			&post.IsLikedByUser, &post.IsDislikedByUser,
		)
		if err != nil {
			return Post{}, fmt.Errorf("error scanning row: %v", err)
		}

		// Assign user to post
		if post.UserId == 0 { // If this is the first time we're encountering the post
			post.User = user
		}

		// Ensure unique categories
		isCategoryAdded := false
		for _, c := range post.Categories {
			if c.ID == category.ID {
				isCategoryAdded = true
				break
			}
		}
		if !isCategoryAdded && category.ID != 0 {
			post.Categories = append(post.Categories, category)
		}

		// Ensure unique post files
		isFileAdded := false
		for _, f := range post.PostFiles {
			if f.ID == postFile.ID {
				isFileAdded = true
				break
			}
		}
		if !isFileAdded && postFile.ID != 0 {
			post.PostFiles = append(post.PostFiles, postFile)
		}
	}

	// If no rows were returned, the post doesn't exist
	if post.ID == 0 {
		return Post{}, fmt.Errorf("post with ID %d not found", postId)
	}

	// Check for any errors during row iteration
	if err := rows.Err(); err != nil {
		return Post{}, fmt.Errorf("row iteration error: %v", err)
	}

	return post, nil
}

func ReadPostByUUID(postUUID string, checkLikeForUser int) (Post, error) {
	// Query the records
	rows, selectError := sqlDB.Query(`
        SELECT p.id as post_id, p.uuid as post_uuid, p.title as post_title, p.content as post_content, p.status as post_status, p.created_at as post_created_at, p.updated_at as post_updated_at, p.updated_by as post_updated_by,
			p.like_count as post_like_count, p.dislike_count as post_dislike_count, p.comment_count as post_comment_count,
			u.id as user_id, u.first_name as user_first_name, u.last_name as user_last_name, u.nick_name as user_nick_name, u.email as user_email, IFNULL(u.profile_image, '') as profile_image,
			c.id as category_id, c.name as category_name,
			IFNULL(pf.id, 0) as post_file_id, pf.file_uploaded_name, pf.file_real_name,
			CASE 
                WHEN EXISTS (SELECT 1 FROM post_likes WHERE post_id = p.id AND status != 'delete' AND type = 'like' AND user_id = ?) THEN 1
                ELSE 0
            END AS is_liked_by_user,
            CASE 
                WHEN EXISTS (SELECT 1 FROM post_likes WHERE post_id = p.id AND status != 'delete' AND type = 'dislike' AND user_id = ?) THEN 1
                ELSE 0
            END AS is_disliked_by_user
		FROM posts p
			INNER JOIN users u
				ON p.user_id = u.id
				AND p.uuid = ?
			LEFT JOIN post_files pf
				ON p.id = pf.post_id
				AND pf.status = 'enable'
			LEFT JOIN post_categories pc
				ON p.id = pc.post_id
				AND pc.status = 'enable'
			LEFT JOIN categories c
				ON pc.category_id = c.id
				AND c.status = 'enable'
		WHERE p.status != 'delete'
			AND u.status != 'delete';
    `, checkLikeForUser, checkLikeForUser, postUUID)
	if selectError != nil {
		return Post{}, selectError
	}
	defer rows.Close()

	var post Post
	post.Categories = []Category{}
	post.PostFiles = []PostFile{}
	var user userManagementModels.User

	// Scan the records
	for rows.Next() {
		var category Category
		var postFile PostFile

		err := rows.Scan(
			&post.ID, &post.UUID, &post.Title, &post.Content, &post.Status,
			&post.CreatedAt, &post.UpdatedAt, &post.UpdatedBy,
			&post.LikeCount, &post.DisikeCount, &post.CommentCount,
			&post.UserId, &user.FirstName, &user.LastName, &user.NickName, &user.Email, &user.ProfileImage,
			&category.ID, &category.Name,
			&postFile.ID, &postFile.FileUploadedName, &postFile.FileRealName,
			&post.IsLikedByUser, &post.IsDislikedByUser,
		)
		if err != nil {
			return Post{}, fmt.Errorf("error scanning row: %v", err)
		}

		// Ensure unique categories
		isCategoryAdded := false
		for _, c := range post.Categories {
			if c.ID == category.ID {
				isCategoryAdded = true
				break
			}
		}
		if !isCategoryAdded && category.ID != 0 {
			post.Categories = append(post.Categories, category)
		}

		// Ensure unique post files
		isFileAdded := false
		for _, f := range post.PostFiles {
			if f.ID == postFile.ID {
				isFileAdded = true
				break
			}
		}
		if !isFileAdded && postFile.ID != 0 {
			post.PostFiles = append(post.PostFiles, postFile)
		}
	}

	// If no rows were returned, the post doesn't exist
	if post.ID == 0 {
		return Post{}, fmt.Errorf("post with UUID %s not found", postUUID)
	}

	post.User = user

	// Check for any errors during row iteration
	if err := rows.Err(); err != nil {
		return Post{}, fmt.Errorf("row iteration error: %v", err)
	}

	return post, nil
}

func ReadPostByUserID(postId int, userID int) (Post, error) {
	// Updated query to join comments with posts
	rows, selectError := sqlDB.Query(`
        SELECT p.id as post_id, p.uuid as post_uuid, p.title as post_title, p.content as post_content, p.status as post_status, p.created_at as post_created_at, p.updated_at as post_updated_at, p.updated_by as post_updated_by,
			p.like_count as post_like_count, p.dislike_count as post_dislike_count, p.comment_count as post_comment_count,
			p.user_id as post_user_id, u.id as user_id, u.first_name as user_first_name, u.last_name as user_last_name, u.nick_name as user_nick_name, u.email as user_email, IFNULL(u.profile_image, '') as profile_image,
			c.id as category_id, c.name as category_name,
			IFNULL(pf.id, 0) as post_file_id, pf.file_uploaded_name, pf.file_real_name
		FROM posts p
			INNER JOIN users u
				ON p.user_id = u.id
				AND p.id = ?
			LEFT JOIN post_files pf
				ON p.id = pf.post_id
				AND pf.status = 'enable'
			LEFT JOIN post_categories pc
				ON p.id = pc.post_id
				AND pc.status = 'enable'
			LEFT JOIN categories c
				ON pc.category_id = c.id
				AND c.status = 'enable'
		WHERE p.status != 'delete'
			AND u.status != 'delete';
    `, postId)
	if selectError != nil {
		return Post{}, selectError
	}
	defer rows.Close()

	var post Post
	var user userManagementModels.User

	// Scan the records
	for rows.Next() {
		var category Category
		var postFile PostFile
		var Type string
		err := rows.Scan(
			&post.ID, &post.UUID, &post.Title, &post.Content, &post.Status,
			&post.CreatedAt, &post.UpdatedAt, &post.UpdatedBy, &post.UserId,
			&post.LikeCount, &post.DisikeCount, &post.CommentCount,
			&user.ID, &user.FirstName, &user.LastName, &user.NickName, &user.Email, &user.ProfileImage,
			&category.ID, &category.Name,
			&postFile.ID, &postFile.FileUploadedName, &postFile.FileRealName,
			&Type,
		)
		if err != nil {
			return Post{}, fmt.Errorf("error scanning row: %v", err)
		}
		if user.ID == userID {
			if Type == "like" {
				post.IsLikedByUser = true
			} else if Type == "dislike" {
				post.IsDislikedByUser = true
			}
		}

		// Ensure unique categories
		isCategoryAdded := false
		for _, c := range post.Categories {
			if c.ID == category.ID {
				isCategoryAdded = true
				break
			}
		}
		if !isCategoryAdded && category.ID != 0 {
			post.Categories = append(post.Categories, category)
		}

		// Ensure unique post files
		isFileAdded := false
		for _, f := range post.PostFiles {
			if f.ID == postFile.ID {
				isFileAdded = true
				break
			}
		}
		if !isFileAdded && postFile.ID != 0 {
			post.PostFiles = append(post.PostFiles, postFile)
		}
	}

	post.User = user

	// Check for any errors during row iteration
	if err := rows.Err(); err != nil {
		return Post{}, fmt.Errorf("row iteration error: %v", err)
	}

	return post, nil
}
