package main

import (
	"testing"
)

func TestUserIsGuest(t *testing.T) {
	user := &User{Name: "Guest1234"}
	is_guest := user.IsGuest()

	if false == is_guest {
		t.Errorf("Expecting: true but got: %s", is_guest)
	}

	user.Name = "brandon"
	is_guest = user.IsGuest()

	if true == is_guest {
		t.Errorf("Expecting: false but got: %s", is_guest)
	}
}

func TestUserIsLoggedIn(t *testing.T) {
	user := &User{Name: "Guest1234"}
	is_logged_in := user.IsLoggedIn()

	if true == is_logged_in {
		t.Errorf("Expecting: false but got: %s", is_logged_in)
	}

	user.Name = "brandon"
	is_logged_in = user.IsLoggedIn()

	if false == is_logged_in {
		t.Errorf("Expecting: true but got: %s", is_logged_in)
	}
}

func TestUserIncrementHeadshot(t *testing.T) {
	user := &User{Id: 1, Name: "brandon", Headshot: 1}
	user.IncrementHeadshot()

	var headshot uint8
	row, _ := DbSelectOne("SELECT headshot FROM users WHERE id = 1")
	row.Scan(&headshot)

	if 2 != headshot {
		t.Errorf("Expecting: 2 but got: %s", headshot)
	}

	// after 255, should go back to 1
	user.Headshot = 255
	user.IncrementHeadshot()

	row, _ = DbSelectOne("SELECT headshot FROM users WHERE id = 1")
	row.Scan(&headshot)

	if 1 != headshot {
		t.Errorf("Expecting: 1 but got: %s", headshot)
	}
}

func TestMessageSaveWithGuest(t *testing.T) {
	user := &User{Id: 1, Name: "Guest1"}
	room := &Room{Id: 3, Name: "lobby"}
	message := &Message{User: user, Message: "testing", Room: room}

	saved := message.Save()
	if true != saved {
		t.Errorf("Expecting: !0 but got: %s", saved)
	}

	var id int
	var room_id int
	var user_id int
	var guest_id int
	var message_text string
	row, _ := DbSelectOne("SELECT id, room_id, user_id, guest_id, message FROM messages ORDER BY id DESC LIMIT 1")
	row.Scan(&id, &room_id, &user_id, &guest_id, &message_text)

	if 3 != room_id {
		t.Errorf("Expecting: 3 but got: %s", room_id)
	}

	if 0 != user_id {
		t.Errorf("Expecting: 0 but got: %s", user_id)
	}

	if 1 != guest_id {
		t.Errorf("Expecting: 1 but got: %s", guest_id)
	}

	if "testing" != message_text {
		t.Errorf("Expecting: testing but got: %s", message_text)
	}

	DbExec("DELETE FROM messages WHERE id = ?", id)
}

func TestMessageSaveWithUser(t *testing.T) {
	user := &User{Id: 1, Name: "brandon"}
	room := &Room{Id: 3, Name: "lobby"}
	message := &Message{User: user, Message: "testing", Room: room}

	saved := message.Save()
	if true != saved {
		t.Errorf("Expecting: !0 but got: %s", saved)
	}

	var id int
	var room_id int
	var user_id int
	var guest_id int
	var message_text string
	row, _ := DbSelectOne("SELECT id, room_id, user_id, guest_id, message FROM messages ORDER BY id DESC LIMIT 1")
	row.Scan(&id, &room_id, &user_id, &guest_id, &message_text)

	if 3 != room_id {
		t.Errorf("Expecting: 3 but got: %s", room_id)
	}

	if 1 != user_id {
		t.Errorf("Expecting: 1 but got: %s", user_id)
	}

	if 0 != guest_id {
		t.Errorf("Expecting: 0 but got: %s", guest_id)
	}

	if "testing" != message_text {
		t.Errorf("Expecting: testing but got: %s", message_text)
	}

	DbExec("DELETE FROM messages WHERE id = ?", id)
}

// func TestMessageGetCreated(t *testing.T) {
// 	today := time.Now()
// 	format := "2006-01-02 15:04:05"
// 	created, _ := time.Parse(format, "2010-01-02 01:02:00")
// 	expects := "1/2/2010 @ 1:02am UTC"
// 	message := &Message{ Created: created }

// 	if expects != message.GetCreated() {
// 		t.Errorf("Expecting: %s but got: %s", expects, message.GetCreated())
// 	}

// 	message.Created, _ = time.Parse(format, "2010-01-02 16:02:00")
// 	expects = "1/2/2010 @ 4:02pm UTC"

// 	if expects != message.GetCreated() {
// 		t.Errorf("Expecting: %s but got: %s", expects, message.GetCreated())
// 	}

// 	message.Created, _ = time.Parse(format, fmt.Sprintf("%d-01-02 16:02:00", today.Year()))
// 	expects = "1/2 @ 4:02pm UTC"

// 	if expects != message.GetCreated() {
// 		t.Errorf("Expecting: %s but got: %s", expects, message.GetCreated())
// 	}

// 	message.Created, _ = time.Parse(format, fmt.Sprintf("%d-%02d-%02d 16:02:00", today.Year(), today.Month(), today.Day()))
// 	expects = "4:02pm UTC"

// 	if expects != message.GetCreated() {
// 		t.Errorf("Expecting: %s but got: %s", expects, message.GetCreated())
// 	}

// 	fmt.Println(message.Created)
// 	fmt.Println(message.GetCreated())
// }

func TestMessageGetJSON(t *testing.T) {
	user := &User{Id: 1, Name: "brandon", Headshot: 1}
	room := &Room{Id: 3, Name: "lobby"}
	message := &Message{User: user, Message: "testing", Room: room}

	resp := "{\"message\":\"testing\",\"user\":{\"name\":\"brandon\",\"headshot\":\"1\"}}"
	json := message.GetJSON()
	if resp != json {
		t.Errorf("Expecting: %s but got: %s", resp, json)
	}

	message.Action = "join"
	resp = "{\"action\":\"join\",\"user\":{\"name\":\"brandon\",\"headshot\":\"1\"}}"
	json = message.GetJSON()
	if resp != json {
		t.Errorf("Expecting: %s but got: %s", resp, json)
	}
}

func TestRoomGetName(t *testing.T) {
	room := &Room{}

	if "lobby" != room.GetName() {
		t.Errorf("Expecting: %s but got: %s", "test", "lobby", room.GetName())
	}

	room.Name = "TEST"
	if "test" != room.GetName() {
		t.Errorf("Expecting: %s but got: %s", "test", room.GetName())
	}
}

func TestRoomGetId(t *testing.T) {
	room := &Room{}

	if 3 != room.GetId() {
		t.Errorf("Expecting: %s but got: %s", 1, room.GetId())
	}

	room.Id = 0
	room.Name = "new-room-name"
	if 0 == room.GetId() {
		t.Errorf("Expecting: !0 but got: %s", room.GetId())
	}
}

func TestRoomGetLastMessages(t *testing.T) {
	test_message := "|||hiya|||"
	user := &User{Id: 1, Name: "brandon"}
	room := &Room{Id: 3, Name: "lobby"}
	message := &Message{User: user, Message: test_message, Room: room}
	message.Save()

	messages := room.GetLastMessages(3)
	if test_message != messages[0].Message {
		t.Errorf("Expecting: %s but got: %s", test_message, messages[0].Message)
	}

	DbExec("DELETE FROM messages WHERE user_id = 1 AND room_id = 3 AND message = ?", test_message)
}

func TestPageIsProduction(t *testing.T) {
	page := &Page{Env: "development"}

	if page.IsProduction() {
		t.Errorf("Expecting: false but got: %s", page.IsProduction())
	}

	page.Env = "production"

	if !page.IsProduction() {
		t.Errorf("Expecting: true but got: %s", page.IsProduction())
	}
}
