package main

import (
	"fmt"
	"html/template"
	"strings"
	"time"
)

type User struct {
	Id       uint32
	Name     string
	Headshot uint8 `json:"headshot,string"`
}

type Room struct {
	Id   int64
	Name string
}

type Message struct {
	User    *User
	Message string
	Room    *Room
	Action  string
	Created time.Time
}

type Page struct {
	Room        *Room
	User        *User
	HeadshotImg template.HTML
	Env         string
}

func (u *User) IsGuest() bool {
	guest_len := 5
	username_len := len(u.Name)

	if username_len < guest_len {
		guest_len = username_len
	}

	if "Guest" == u.Name[:guest_len] {
		return true
	}

	return false
}

func (u *User) IsLoggedIn() bool {
	if u.IsGuest() {
		return false
	}

	return true
}

func (u *User) IncrementHeadshot() uint8 {
	if u.Headshot < 255 {
		u.Headshot = u.Headshot + 1
	} else {
		u.Headshot = 1
	}

	sql := "UPDATE users SET headshot = ? WHERE id = ?"
	_, err := DbExec(sql, u.Headshot, u.Id)
	if err != nil {
		fmt.Println(err)
	}

	return u.Headshot
}

func (m *Message) Save() bool {
	var guest_id uint32 = 0
	var user_id uint32 = 0
	if m.User.IsGuest() {
		guest_id = m.User.Id
	} else {
		user_id = m.User.Id
	}

	// joins are handled in hub.go
	if m.Action == "join" {
		return false
	}

	if "" == m.Message {
		return true
	}

	_, err := DbExec("INSERT INTO messages (room_id, user_id, guest_id, message, created) VALUES (?, ?, ?, ?, UTC_TIMESTAMP())", m.Room.GetId(), user_id, guest_id, m.Message)
	defer DbClose()
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true
}

func (m *Message) GetCreated() string {
	return fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02dZ", m.Created.Year(), m.Created.Month(), m.Created.Day(), m.Created.Hour(), m.Created.Minute(), m.Created.Second())
}

func (m *Message) GetJSON() string {
	var json string
	var created string = ""

	if !m.Created.IsZero() {
		created = fmt.Sprintf(",\"created\":\"%s\"", m.GetCreated())
	}

	if "" == m.Action {
		json = fmt.Sprintf("{\"message\":\"%s\",\"user\":{\"name\":\"%s\",\"headshot\":\"%d\"}%s}", strings.Replace(m.Message, "\"", "&quot;", -1), m.User.Name, m.User.Headshot, created)
	} else {
		json = fmt.Sprintf("{\"action\":\"%s\",\"user\":{\"name\":\"%s\",\"headshot\":\"%d\"}%s}", m.Action, m.User.Name, m.User.Headshot, created)
	}

	return json
}

func (r *Room) GetName() string {
	if "" == r.Name {
		return "lobby"
	}

	return strings.ToLower(r.Name)
}

func (r *Room) GetId() int64 {
	if r.Id > 0 {
		return r.Id
	}

	r.Name = r.GetName()
	sql := "SELECT id FROM rooms WHERE name = ?"
	row, err := DbSelectOne(sql, r.Name)
	if err != nil {
		fmt.Println(err)
	}

	if err = row.Scan(&r.Id); err != nil {
		// sql: no rows in result set
		//fmt.Println(err)
	}

	if r.Id > 0 {
		return r.Id
	}

	sql = "INSERT INTO rooms (name) VALUES (?)"
	result, err := DbExec(sql, r.Name)
	if err != nil {
		fmt.Println(err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		fmt.Println(err)
	}

	r.Id = id

	return r.Id
}

func (r *Room) GetLastMessages(limit uint8) []*Message {
	sql := `SELECT
			IF(user_id > 0, u.id, m.guest_id) AS user_id,
			IF(user_id > 0, u.name, CONCAT('Guest', guest_id)) AS name,
			IF(user_id > 0, u.headshot, 0) AS headshot,
			m.message,
			m.created
		FROM messages m
		LEFT JOIN users u ON m.user_id = u.id
		WHERE m.room_id = ?
		ORDER BY m.created DESC
		LIMIT ?`
	rows, err := DbSelect(sql, r.GetId(), limit)
	defer DbClose()
	if err != nil {
		fmt.Println(err)
	}

	var messages = make([]*Message, limit)
	var i uint8 = 0

	for rows.Next() {
		var created string
		user := &User{}
		message := &Message{}
		message.User = user
		rows.Scan(&message.User.Id, &message.User.Name, &message.User.Headshot, &message.Message, &created)

		message.Created, err = time.Parse("2006-01-02 15:04:05", created)
		if err != nil {
			fmt.Println(err)
		}

		messages[i] = message
		i++
	}
	rows.Close()

	return messages
}

func (p *Page) IsProduction() bool {
	if "production" == p.Env {
		return true
	}

	return false
}
