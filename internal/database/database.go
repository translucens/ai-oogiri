package database

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/go-sql-driver/mysql"
	_ "github.com/go-sql-driver/mysql"
)

type Riddle struct {
	ID         int
	CreatedAt  time.Time
	Theme      string
	HotAnswer  string
	ColdAnswer string
}

type Client struct {
	db *sql.DB

	addStmt   *sql.Stmt
	queryStmt *sql.Stmt
}

const timezone = "Asia/Tokyo"

func NewClient(ctx context.Context, username, password, host string, port int, socket, database string) (*Client, error) {

	var net, addr string
	if socket != "" {
		net = "unix"
		addr = socket
	} else if host != "" && port > 0 {
		net = "tcp"
		addr = host + ":" + strconv.Itoa(port)
	} else {
		return nil, fmt.Errorf("invalid database configuration")
	}

	tz, err := time.LoadLocation(timezone)
	if err != nil {
		return nil, err
	}

	cfg := mysql.Config{
		User:                 username,
		Passwd:               password,
		Net:                  net,
		Addr:                 addr,
		DBName:               database,
		ParseTime:            true,
		AllowNativePasswords: true,
		Loc:                  tz,
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, err
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	db.Ping()

	addStmt, err := db.PrepareContext(ctx, "INSERT INTO riddles (theme, hot_answer, cold_answer) VALUES (?, ?, ?)")
	if err != nil {
		return nil, fmt.Errorf("PrepareContext(add): %v", err)
	}

	queryStmt, err := db.PrepareContext(ctx, "SELECT id, theme, hot_answer, cold_answer, created_at FROM riddles ORDER BY created_at DESC")
	if err != nil {
		return nil, fmt.Errorf("PrepareContext(query): %v", err)
	}

	return &Client{db: db, addStmt: addStmt, queryStmt: queryStmt}, nil
}

func (c *Client) Close() error {
	c.addStmt.Close()
	c.queryStmt.Close()

	return c.db.Close()
}

func (c *Client) AddRiddle(ctx context.Context, theme string, hotAns string, coldAns string) error {

	_, err := c.addStmt.ExecContext(ctx, theme, hotAns, coldAns)
	if err != nil {
		return fmt.Errorf("ExecContext: %v", err)
	}

	return nil
}

func (c *Client) GetHistory(ctx context.Context) ([]Riddle, error) {

	rows, err := c.queryStmt.QueryContext(ctx)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var id int
	var theme, hotAnswer, ColdAnswer string
	var timestamp time.Time

	var conversations []Riddle

	for rows.Next() {
		if err := rows.Scan(&id, &theme, &hotAnswer, &ColdAnswer, &timestamp); err != nil {
			return nil, err
		}

		conversations = append(conversations, Riddle{ID: id, CreatedAt: timestamp, Theme: theme, HotAnswer: hotAnswer, ColdAnswer: ColdAnswer})
	}

	return conversations, nil
}
