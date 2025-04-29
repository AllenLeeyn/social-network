package dbTools

import "fmt"

func (db *DBContainer) SelectUserByField(fieldName string, fieldValue interface{}) (*User, error) {
	if fieldName != "id" && fieldName != "nick_name" && fieldName != "email" {
		return nil, fmt.Errorf("invalid field")
	}
	qry := `SELECT * FROM users WHERE ` + fieldName + ` = ?`
	var u User
	err := db.conn.QueryRow(qry, fieldValue).Scan(
		&u.ID,
		&u.TypeID,
		&u.FirstName,
		&u.LastName,
		&u.NickName,
		&u.Gender,
		&u.Age,
		&u.Email,
		&u.PwHash,
		&u.RegDate,
		&u.LastLogin)
	if err != nil {
		return nil, checkErrNoRows(err)
	}
	return &u, nil
}

// db.InserUser() insert a User into the database
func (db *DBContainer) InsertUser(u *User) (int, error) {
	qry := `INSERT INTO users 
			(type_id, first_name, last_name, nick_name, gender, age, email, pw_hash) 
			VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	res, err := db.conn.Exec(qry,
		u.TypeID,
		u.FirstName,
		u.LastName,
		u.NickName,
		u.Gender,
		u.Age,
		u.Email,
		u.PwHash)
	if err != nil {
		return -1, err
	}
	userID, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}
	return int(userID), nil
}

// db.UpdateUser() info like name, pwHash and lastLogin
func (db *DBContainer) UpdateUser(u *User) error {
	qry := `UPDATE users
			SET first_name = ?, last_name = ?, gender = ?, age = ?, pw_hash = ?, last_login = ?
			WHERE id = ?`
	_, err := db.conn.Exec(qry,
		u.FirstName,
		u.LastName,
		u.Gender,
		u.Age,
		u.PwHash,
		u.LastLogin,
		u.ID)
	return err
}
