package database

import (
	"database/sql"
	"errors"
	"log"
	"strconv"
)

func Database() {
	databaseAll := InitDatabase("forum.db")
	defer databaseAll.Close()
	sqlStmt := `

	CREATE TABLE IF NOT EXISTS users (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
		email TEXT NOT NULL, 
		username TEXT NOT NULL, 
		password TEXT NOT NULL,
		img TEXT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS reponses (
		id INTEGER NOT NULL,
		idPost INTEGER NOT NULL,
		userName TEXT NOT NULL,
		contenu TEXT NOT NULL,
		date TEXT NOT NULL,
		img TEXT NOT NULL,
		PRIMARY KEY (id)
	);

	CREATE TABLE IF NOT EXISTS posts (
		id INTEGER NOT NULL,
		userName Text NOT NULL,
		tag TEXT NOT NULL,
		titre TEXT NOT NULL,
		description TEXT NOT NULL,
		nbrLikes INTEGER,
		nbrDislikes INTEGER,
		date TEXT NOT NULL,
		img TEXT NOT NULL,
		PRIMARY KEY (id)
	);

	CREATE TABLE IF NOT EXISTS sessions (
		id INTEGER NOT NULL,
		email TEXT NOT NULL,
		uuid TEXT NOT NULL,
		PRIMARY KEY (id)
	);
		`

	_, err := databaseAll.Exec(sqlStmt)
	if err != nil {
		log.Fatal(err)
	}
}

func DatabaseAndUsers(values []string) {
	db := InitDatabase("forum.db")
	defer db.Close()
	sqlStmtInsertUsers := `
		INSERT INTO users (email, username, password, img) VALUES (?, ?, ?, ?);
		`
	InsertIntoRow(db, values, sqlStmtInsertUsers)
}

func DatabaseAndReponse(values []string) {
	db := InitDatabase("forum.db")
	defer db.Close()
	sqlStmtInsertReponses := `
		INSERT INTO reponses (idPost, userName,contenu, date, img) VALUES (?, ?, ?, ?, ?);
		`
	InsertIntoRow(db, values, sqlStmtInsertReponses)
}

func DatabaseAndPost(values []string) {
	db := InitDatabase("forum.db")
	defer db.Close()
	sqlStmtInsertPosts := `
		INSERT INTO posts (userName, tag, titre, description, nbrLikes, nbrDislikes, date, img) VALUES (?, ?, ?, ?, ?, ?, ?, ?);
		`
	InsertIntoRow(db, values, sqlStmtInsertPosts)
}

func DatabaseAndSession(values []string) error {
	db := InitDatabase("forum.db")
	defer db.Close()
	if sessionExists(db, values[0]) {
		return errors.New("Votre session est déjà ouverte. Vous ne pouvez pas en ouvrir une autre.")
	}

	sqlStmtInsertPosts := `
        INSERT INTO sessions (email, uuid) VALUES (?, ?);
    `
	_, err := InsertIntoRow(db, values, sqlStmtInsertPosts)
	return err
}

func sessionExists(db *sql.DB, email string) bool {
	query := "SELECT COUNT(*) FROM sessions WHERE email = ?;"
	var count int
	err := db.QueryRow(query, email).Scan(&count)
	if err != nil {
		log.Fatal(err)
	}
	return count > 0
}

func InitDatabase(dataBaseName string) *sql.DB {
	database, err := sql.Open("sqlite3", "./Model/"+dataBaseName)
	if err != nil {
		log.Fatal(err)
	}
	return database
}

func InsertIntoRow(db *sql.DB, values []string, stmt string) (int64, error) {
	args := make([]interface{}, len(values))
	for i, v := range values {
		args[i] = v
	}
	row, err := db.Exec(stmt, args...)
	if err != nil {
		log.Fatal(err)
	}
	return row.LastInsertId()
}

func GetAllPost() []Post {
	db := InitDatabase("forum.db")
	defer db.Close()

	query := "SELECT * FROM posts;"
	result, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	var res []Post
	for result.Next() {
		var post Post
		err := result.Scan(&post.Id, &post.UserName, &post.Tag, &post.Titre, &post.Description, &post.NbrLikes, &post.NbrDislikes, &post.Date, &post.Img)
		if err != nil {
			log.Fatal(err)
		}
		res = append(res, post)
	}

	return res
}

func GetOnePost(id string) Post {
	db := InitDatabase("forum.db")
	defer db.Close()

	query := "SELECT * FROM posts WHERE id= '" + id + "';"
	result, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	var post Post
	for result.Next() {
		err2 := result.Scan(&post.Id, &post.UserName, &post.Tag, &post.Titre, &post.Description, &post.NbrLikes, &post.NbrDislikes, &post.Date, &post.Img)
		if err2 != nil {
			log.Fatal(err)
		}
	}
	return post
}

func GetResponses(id string) []Reponse {
	db := InitDatabase("forum.db")
	defer db.Close()

	query := "SELECT * FROM reponses WHERE idPost = '" + id + "';"
	result, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}

	var res []Reponse
	for result.Next() {
		var reponse Reponse
		err := result.Scan(&reponse.Id, &reponse.IdPost, &reponse.UserName, &reponse.Contenu, &reponse.Date, &reponse.Img)
		if err != nil {
			log.Fatal(err)
		}
		res = append(res, reponse)
	}
	return res
}

func GetUser(email string) User {
	db := InitDatabase("forum.db")
	defer db.Close()

	query := "SELECT * FROM users WHERE email= '" + email + "';"
	result, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	var user User
	for result.Next() {
		err := result.Scan(&user.Id, &user.Email, &user.Username, &user.Password, &user.Img)
		if err != nil {
			log.Fatal(err)
		}
	}
	return user
}

func GetEmail(email string) bool {
	db := InitDatabase("forum.db")
	defer db.Close()

	query := "SELECT * FROM users WHERE email= '" + email + "';"
	result, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	var res []User
	for result.Next() {
		var user User
		err := result.Scan(&user.Id, &user.Email, &user.Username, &user.Password, &user.Img)
		if err != nil {
			log.Fatal(err)
		}
		res = append(res, user)
	}

	if len(res) != 0 {
		return true
	} else {
		return false
	}
}

func GetTagAnime() []Post {
	db := InitDatabase("forum.db")
	defer db.Close()

	query := "SELECT * FROM posts WHERE tag='anime';"
	result, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	var res []Post
	for result.Next() {
		var post Post
		err := result.Scan(&post.Id, &post.UserName, &post.Tag, &post.Titre, &post.Description, &post.NbrLikes, &post.NbrDislikes, &post.Date, &post.Img)
		if err != nil {
			log.Fatal(err)
		}
		res = append(res, post)
	}

	return res
}

func GetTagSerie() []Post {
	db := InitDatabase("forum.db")
	defer db.Close()

	query := "SELECT * FROM posts WHERE tag='serie';"
	result, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	var res []Post
	for result.Next() {
		var post Post
		err := result.Scan(&post.Id, &post.UserName, &post.Tag, &post.Titre, &post.Description, &post.NbrLikes, &post.NbrDislikes, &post.Date, &post.Img)
		if err != nil {
			log.Fatal(err)
		}
		res = append(res, post)
	}

	return res
}

func GetSession(mail string) bool {
	db := InitDatabase("forum.db")
	defer db.Close()

	query := "SELECT * FROM sessions  WHERE email='" + mail + "';"
	result, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}

	var res []Session
	for result.Next() {
		var session Session
		err := result.Scan(&session.Id, &session.Email, &session.Uuid)
		if err != nil {
			log.Fatal(err)
		}
		res = append(res, session)
	}

	if len(res) != 0 {
		return false
	} else {
		return true
	}
}

func DeleteSession(uuid string) {
	db := InitDatabase("forum.db")
	defer db.Close()

	query := "DELETE FROM sessions WHERE uuid = '" + uuid + "';"
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}

func SelectByAscending(filter string) []Post {
	db := InitDatabase("forum.db")
	defer db.Close()
	query := "SELECT * FROM posts ORDER BY " + filter + " ASC;"
	result, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	var res []Post
	for result.Next() {
		var post Post
		err := result.Scan(&post.Id, &post.UserName, &post.Tag, &post.Titre, &post.Description, &post.NbrLikes, &post.NbrDislikes, &post.Date, &post.Img)
		if err != nil {
			log.Fatal(err)
		}
		res = append(res, post)
	}
	return res
}

func SelectByDescending(filter string) []Post {
	db := InitDatabase("forum.db")
	defer db.Close()
	query := "SELECT * FROM posts ORDER BY " + filter + " DESC;"
	result, err := db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	var res []Post
	for result.Next() {
		var post Post
		err := result.Scan(&post.Id, &post.UserName, &post.Tag, &post.Titre, &post.Description, &post.NbrLikes, &post.NbrDislikes, &post.Date, &post.Img)
		if err != nil {
			log.Fatal(err)
		}
		res = append(res, post)
	}
	return res
}

func RecupNbr(data string, id string) int {
	db := InitDatabase("forum.db")
	defer db.Close()

	quer := "SELECT " + data + " FROM posts WHERE id=" + id + ";"
	resul, e := db.Query(quer)
	if e != nil {
		log.Fatal(e)
	}
	var nbrLikes int
	for resul.Next() {
		er := resul.Scan(&nbrLikes)
		if er != nil {
			log.Fatal(er)
		}
	}
	return nbrLikes
}

func UpdateNbr(data string, nbr int, id string) {
	db := InitDatabase("forum.db")
	defer db.Close()
	query := "UPDATE posts SET " + data + " = " + strconv.Itoa(nbr) + " WHERE id= " + id + ";"
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}

func UpdateUserInfo(oldEmail, newUsername, newEmail, newPassword string) error {
	db := InitDatabase("forum.db")
	defer db.Close()

	// Vérifiez d'abord si le nouvel e-mail n'est pas déjà utilisé par un autre utilisateur
	if newEmail != oldEmail && GetEmail(newEmail) {
		return errors.New("Cet e-mail est déjà utilisé par un autre utilisateur")
	}

	// Mettez à jour les informations de l'utilisateur
	query := `
		UPDATE users
		SET username = ?, email = ?, password = ?
		WHERE email = ?;
	`
	_, err := db.Exec(query, newUsername, newEmail, newPassword, oldEmail)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}
