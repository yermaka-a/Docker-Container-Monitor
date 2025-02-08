package db

import (
	"back/internal/config"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

const (
	createTableSQL = `CREATE TABLE IF NOT EXISTS containers (
        id SERIAL PRIMARY KEY,
        ip TEXT NOT NULL,
		timeMs TEXT NOT NULL,
		pingDate TEXT NOT NULL
    );`
)

func InitDB() {
	conf := config.New()
	var err error
	var db *sql.DB
	DBConfig := conf.DataBase
	psqlInfo := fmt.Sprintf("postgres://%s:%s@%s:%s?sslmode=disable",
		DBConfig.POSTGRES_USER, DBConfig.POSTGRES_PASSWORD, DBConfig.POSTGRES_HOST, DBConfig.POSTGRES_PORT)
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Panic(err)
	}
	query := fmt.Sprintf(`DO $$ BEGIN IF NOT EXISTS (SELECT FROM pg_database WHERE datname = '%s') THEN CREATE DATABASE %s; END IF; END $$;`, DBConfig.POSTGRES_DB, DBConfig.POSTGRES_DB)
	_, err = db.Exec(query)
	if err != nil {
		log.Panic(err)
	}
	db.Close()
	psqlInfo = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		DBConfig.POSTGRES_USER, DBConfig.POSTGRES_PASSWORD, DBConfig.POSTGRES_HOST, DBConfig.POSTGRES_PORT, DBConfig.POSTGRES_DB)
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Panic(err)
	}
	db.Exec(createTableSQL)
	log.Printf("Database %s sucessfully created\n", DBConfig.POSTGRES_DB)
	defer db.Close()
}

func connectDB() *sql.DB {
	conf := config.New()
	var err error
	var db *sql.DB
	DBConfig := conf.DataBase
	psqlInfo := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		DBConfig.POSTGRES_USER, DBConfig.POSTGRES_PASSWORD, DBConfig.POSTGRES_HOST, DBConfig.POSTGRES_PORT, DBConfig.POSTGRES_DB)
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Panic(err)
	}
	err = db.Ping()
	if err != nil {
		log.Panic(err)
	}
	log.Printf("%s connected\n", DBConfig.POSTGRES_DB)
	return db
}

const (
	getSQL    = `SELECT * FROM containers`
	selectSQL = "SELECT ip FROM containers WHERE ip = $1"
	insertSQL = "INSERT INTO containers (ip, timeMs, pingDate) VALUES ($1, $2, $3);"
	updateSQL = `UPDATE containers SET timeMs = $1,
											pingDate = $2
   										WHERE ip = $3`
)

func GetContainers() *[]Container {
	// var db *sql.DB
	db := connectDB()
	rows, err := db.Query(getSQL)
	if err != nil {
		log.Panic(err)
	}
	defer rows.Close()

	// Обработка результатов
	var containers []Container
	for rows.Next() {
		var container Container
		if err := rows.Scan(&container.ID, &container.Ip, &container.TimeMs, &container.PingDate); err != nil {
			log.Panic(err)
		}
		containers = append(containers, container)
	}

	// Проверка на ошибки после итерации
	if err := rows.Err(); err != nil {
		log.Panic(err)
	}

	return &containers
}

func UpdateContainers(containers *[]Container) bool {
	db := connectDB()
	for _, container := range *containers {
		var ip string
		err := db.QueryRow(selectSQL, container.Ip).Scan(&ip)
		defer db.Close()

		// if err := recover(); err != nil {
		// 	fmt.Print(err)
		// }

		if err == sql.ErrNoRows {
			// Запись не найдена, можно добавить новую
			_, err = db.Exec(insertSQL,
				container.Ip,
				container.TimeMs,
				container.PingDate)
			if err != nil {
				log.Panicf("Error inserting record: %s", err)
			}
			fmt.Println("Record added successfully!")
		} else if err != nil {
			// Ошибка при выполнении запроса
			log.Panicf("Error checking record: %s", err)
		} else {
			// Запись найдена, обновляем её
			_, err = db.Exec(
				updateSQL,
				container.TimeMs,
				container.PingDate,
				container.Ip)
			if err != nil {
				log.Panicf("Error updating record: %s", err)
			}
			fmt.Println("Record updated successfully!")
		}
	}
	return true
}
