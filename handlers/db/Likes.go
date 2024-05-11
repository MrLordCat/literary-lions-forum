package db

import "database/sql"

func AddOrUpdateLike(db *sql.DB, postID, userID, likeType int) error {
	// Проверяем, существует ли уже лайк или дизлайк от этого пользователя к этому посту
	var existingType int
	err := db.QueryRow("SELECT like_type FROM post_likes WHERE post_id = ? AND user_id = ?", postID, userID).Scan(&existingType)
	if err != nil && err != sql.ErrNoRows {
		// Произошла ошибка при выполнении запроса
		return err
	}

	if err == sql.ErrNoRows {
		// Если лайк/дизлайк не существует, вставляем новый
		_, err = db.Exec("INSERT INTO post_likes (post_id, user_id, like_type) VALUES (?, ?, ?)", postID, userID, likeType)
	} else if existingType == likeType {
		// Если существующий лайк/дизлайк совпадает с новым, удаляем его
		_, err = db.Exec("DELETE FROM post_likes WHERE post_id = ? AND user_id = ?", postID, userID)
	} else {
		// Если существующий лайк/дизлайк отличается, обновляем его
		_, err = db.Exec("UPDATE post_likes SET like_type = ? WHERE post_id = ? AND user_id = ?", likeType, postID, userID)
	}
	return err
}

func LikeComment(db *sql.DB, commentID, userID, likeType int) error {
	// Проверяем существует ли лайк
	var existingType int
	err := db.QueryRow("SELECT like_type FROM comment_likes WHERE comment_id = ? AND user_id = ?", commentID, userID).Scan(&existingType)
	if err == sql.ErrNoRows {
		// Если лайк не существует, добавляем новый
		_, err = db.Exec("INSERT INTO comment_likes (comment_id, user_id, like_type) VALUES (?, ?, ?)", commentID, userID, likeType)
		return err
	} else if err != nil {
		// Произошла ошибка при запросе
		return err
	}

	// Если лайк существует и тип совпадает, удаляем его
	if existingType == likeType {
		_, err = db.Exec("DELETE FROM comment_likes WHERE comment_id = ? AND user_id = ?", commentID, userID)
		return err
	}

	// Если тип не совпадает, обновляем существующий лайк
	_, err = db.Exec("UPDATE comment_likes SET like_type = ? WHERE comment_id = ? AND user_id = ?", likeType, commentID, userID)
	return err
}
