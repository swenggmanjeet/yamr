package main

import (
	"code.google.com/p/go.net/websocket"
	"code.google.com/p/gorilla/sessions"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	store.Options = &sessions.Options{
		Path:   "/",
		MaxAge: 31536000, // 1 year
	}
	session, _ := store.Get(r, "yamr")

	if session.Values["id"] == nil {
		res, err := DbExec("INSERT INTO guests (agent, ip, created) VALUES (?, ?, UTC_TIMESTAMP())", r.Header.Get("User-agent"), r.Header.Get("X-Real-IP"))
		if err != nil {
			fmt.Println(err)
			t, _ := template.ParseFiles("views/error.html")
			page := &Page{}
			t.Execute(w, page)
			return
		}
		defer DbClose()
		id, err := res.LastInsertId()

		session.Values["id"] = uint32(id)
		session.Values["user"] = "Guest" + strconv.FormatInt(id, 10)
		session.Values["headshot"] = uint8(0)
		session.Save(r, w)
	}

	user := &User{
		Id:       session.Values["id"].(uint32),
		Name:     fmt.Sprintf("%s", session.Values["user"]),
		Headshot: session.Values["headshot"].(uint8),
	}

	room := &Room{
		Name: r.URL.Path[1:],
	}

	page := &Page{
		Room: room,
		User: user,
	}

	if user.Headshot > 0 {
		page.HeadshotImg = template.HTML(fmt.Sprintf("<img src=\"/headshots/%s.jpg?%d\" width=\"35\" height=\"35\" class=\"left\"/>", user.Name, user.Headshot))
	} else {
		page.HeadshotImg = template.HTML("<img src=\"/images/no_photo.gif\" width=\"35\" height=\"35\" class=\"left\"/>")
	}

	// t, _ := template.ParseFiles("views/index.html")
	// t.Execute(w, page)

	stripped_html := index_html
	stripped_html = strings.Replace(stripped_html, "\t", "", -1)
	stripped_html = strings.Replace(stripped_html, "\n", "", -1)

	t, err := template.New("index").Parse(stripped_html)
	if err != nil {
		fmt.Println(err)
	}
	t.Execute(w, page)
}

func signupHandler(w http.ResponseWriter, r *http.Request) {
	store.Options = &sessions.Options{
		Path:   "/",
		MaxAge: 31536000, // 1 year
	}
	session, _ := store.Get(r, "yamr")

	reg, err := regexp.Compile("[^a-z0-9]+")
	if err != nil {
		fmt.Println(err)
		return
	}

	username := strings.ToLower(strings.Trim(r.FormValue("username"), " "))
	username = reg.ReplaceAllString(username, "")

	password := strings.Trim(r.FormValue("password"), " ")
	guest_len := 5
	username_len := len(username)

	if username_len < guest_len {
		guest_len = username_len
	}

	if "" == username {
		fmt.Fprintf(w, "Please enter a username")
		return
	} else if "" == password {
		fmt.Fprintf(w, "Please enter a password")
		return
	} else if "Guest" == username[:guest_len] {
		fmt.Fprintf(w, "\"Guest\" usernames not allowed")
		return
	}

	sql := "SELECT name, headshot FROM users WHERE name = ?"
	row, err := DbSelectOne(sql, username)
	defer DbClose()
	if err != nil {
		log.Fatal(err)
	}

	user := new(User)
	err = row.Scan(&user.Name, &user.Headshot)

	if "" == user.Name {
		sql := "INSERT INTO users (name, password, created, ip, agent, last_login) VALUES (?, ?, UTC_TIMESTAMP(), ?, ?, CURDATE())"
		res, err := DbExec(sql, username, Sha1(password), r.Header.Get("X-Real-IP"), r.Header.Get("User-agent"))
		defer DbClose()
		id, err := res.LastInsertId()

		if err != nil {
			fmt.Println(err)
			return
		}

		session.Values["id"] = uint32(id)
		session.Values["user"] = username
		session.Values["headshot"] = uint8(0)
		session.Save(r, w)

		fmt.Fprintf(w, "ok")
	} else {
		fmt.Fprintf(w, "Username already exists")
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	store.Options = &sessions.Options{
		Path:   "/",
		MaxAge: 31536000, // 1 year
	}
	session, _ := store.Get(r, "yamr")

	reg, err := regexp.Compile("[^a-z0-9]+")
	if err != nil {
		fmt.Println(err)
		return
	}

	username := strings.ToLower(strings.Trim(r.FormValue("username"), " "))
	username = reg.ReplaceAllString(username, "")

	password := strings.Trim(r.FormValue("password"), " ")

	if "" == username || "" == password {
		fmt.Fprintf(w, "Invalid login")
		return
	}

	sql := "SELECT id, headshot FROM users WHERE name = ? AND password = ?"
	row, err := DbSelectOne(sql, username, Sha1(password))
	defer DbClose()

	if err != nil {
		log.Fatal(err)
	}

	user := new(User)
	err = row.Scan(&user.Id, &user.Headshot)

	if 0 == user.Id {
		fmt.Fprintf(w, "Invalid login")
		return
	}

	session.Values["id"] = user.Id
	session.Values["user"] = username
	session.Values["headshot"] = user.Headshot
	session.Save(r, w)

	sql = "UPDATE users SET ip = ?, agent = ?, last_login = CURDATE() WHERE id = ?"
	DbExec(sql, r.Header.Get("X-Real-IP"), r.Header.Get("User-agent"), user.Id)

	fmt.Fprintf(w, "ok")
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	store.Options = &sessions.Options{
		Path:   "/",
		MaxAge: 31536000, // 1 year
	}
	session, _ := store.Get(r, "yamr")

	delete(session.Values, "id")
	delete(session.Values, "user")
	delete(session.Values, "headshot")
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusFound)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	store.Options = &sessions.Options{
		Path:   "/",
		MaxAge: 31536000, // 1 year
	}
	session, _ := store.Get(r, "yamr")

	if session.Values["id"] == nil {
		fmt.Fprintf(w, "you must be logged in")
		return
	}

	user := &User{
		Id:       session.Values["id"].(uint32),
		Name:     fmt.Sprintf("%s", session.Values["user"]),
		Headshot: session.Values["headshot"].(uint8),
	}

	file, header, err := r.FormFile("headshot")
	if err != nil {
		fmt.Println(err)
	}

	content_type := header.Header.Get("Content-Type")
	response := "error adding headshot"

	if "image/jpeg" == content_type ||
		"image/jpg" == content_type ||
		"image/png" == content_type {

		error := false

		data, err := ioutil.ReadAll(file)
		if err != nil {
			fmt.Println(err)
			error = true
		}

		var temp_file string = ""
		var new_file string = ""

		if false == error {
			temp_file = fmt.Sprintf("/tmp/%s", header.Filename)
			// path, err := syscall.Getwd()

			path, err := filepath.Abs(filepath.Dir(os.Args[0]))
			if err != nil {
				log.Fatal(err)
			}

			new_file = fmt.Sprintf("%s/public/headshots/%s.jpg", path, user.Name)
			err = ioutil.WriteFile(temp_file, data, 0777)
			if err != nil {
				fmt.Println(err)
				error = true
			}
		}

		if false == error {
			cmd := exec.Command("convert", temp_file, new_file)
			if err = cmd.Run(); err != nil {
				fmt.Println("convert ", temp_file, new_file)
				fmt.Println(err)
				error = true
			}
		}

		if false == error {
			cmd := exec.Command("convert", "-scale", "40x40", new_file, new_file)
			if err = cmd.Run(); err != nil {
				fmt.Println("convert -scale 40x40 %s %s", new_file, new_file)
				fmt.Println(err)
				error = true
			}
		}

		if false == error {
			session.Values["headshot"] = user.IncrementHeadshot()
			session.Save(r, w)
			response = fmt.Sprintf("/headshots/%s.jpg?%d", user.Name, user.Headshot)
		}
	} else {
		response = "please use a JPG or PNG"
	}

	fmt.Fprintf(w, response)
}

func wsHandler(ws *websocket.Conn) {
	r := ws.Request()

	store.Options = &sessions.Options{
		Path:   "/",
		MaxAge: 31536000, // 1 year
	}
	session, _ := store.Get(r, "yamr")

	room := &Room{Name: r.URL.Query().Get("r")}

	user := &User{
		Id:       session.Values["id"].(uint32),
		Name:     session.Values["user"].(string),
		Headshot: session.Values["headshot"].(uint8),
	}
	c := &connection{ws: ws, send: make(chan string, 256), user: user, room: room}
	h.register <- c
	defer func() { h.unregister <- c }()
	go c.writer()
	c.reader()
}

func findHandler(w http.ResponseWriter, r *http.Request) {
	handled := false

	if "POST" == r.Method {
		if "/login" == r.URL.Path {
			handled = true
			loginHandler(w, r)
		} else if "/signup" == r.URL.Path {
			handled = true
			signupHandler(w, r)
		} else if "/logout" == r.URL.Path {
			handled = true
			logoutHandler(w, r)
		} else if "/upload" == r.URL.Path {
			handled = true
			uploadHandler(w, r)
		}
	}

	if !handled {
		indexHandler(w, r)
	}
}
