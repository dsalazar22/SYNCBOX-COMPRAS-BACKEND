package libs

import (
	"database/sql"
	"fmt"

	"github.com/go-redis/redis"
	_ "github.com/lib/pq"
)

func PGconnect() *sql.DB {
	//dsn := "postgres://postgres:syncboxdev@127.0.0.1:5432/inside?sslmode=disable"
	dsn := StrConn.Pg
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}
	return db
}

func RedisConnectVars(dbNew int) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     StrConn.Vars, //"127.0.0.1:6379",
		Password: StrConn.Pass,
		DB:       dbNew,
	})
	// fmt.Println(client)
	return client
}

func RedisConnectDest(dbNew int) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     StrConn.Dest, //"127.0.0.1:6379",
		Password: StrConn.Pass,
		DB:       dbNew,
	})
	// fmt.Println(client)
	return client
}

func RedisConnect(dbNew int) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     StrConn.Addr, //"127.0.0.1:6379",
		Password: StrConn.Pass,
		DB:       dbNew,
	})
	// fmt.Println(client)
	return client
}

func SendDB(query string) (string, error) {
	var result string
	conn := PGconnect()
	rows, err := conn.Query(query)
	conn.Close()

	if err != nil {
		fmt.Println(err)
	} else {
		for rows.Next() {
			rows.Scan(&result)
		}
	}

	return result, err
}

func SendOpenDB(query string) *sql.DB {
	// var result string
	conn := PGconnect()
	// conn.Exec(query) //Query(query)
	// conn.Close()
	return conn
}

func CloseOpenDB(conn *sql.DB) {
	conn.Close()
}

func ContentDB(query string) (*sql.Rows, error) {
	conn := PGconnect()
	rows, err := conn.Query(query)
	conn.Close()

	if err == nil {
		fmt.Println(err)
	}

	return rows, err
}
