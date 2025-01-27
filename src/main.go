package main

import (
	"os"
	"io"
	"github.com/Philio/GoMySQL"
	"github.com/garyburd/twister/server"
	"github.com/garyburd/twister/web"
	// "strings"
	"log"
	// "time"
)

var (
	db *mysql.Client
	db_err os.Error
)

func displayIndex(req *web.Request) {
	w := req.Respond(web.StatusOK, web.HeaderContentType, "text/html; charset=\"utf-8\"")
	io.WriteString(w, RenderFile("templates/index.html", map[string]string{"c":"world"}))
}

func createEntry(req *web.Request) {
	message := req.Param.Get("message")
	extra := req.Param.Get("extra")
	tags := splitTags(extra)
	id := NewUUID()
	storeEntry(id, message, tags)
	w := req.Respond(web.StatusOK, web.HeaderContentType, "text/html; charset=\"utf-8\"")
	io.WriteString(w, RenderFile("templates/index.html", map[string]string{"c":"world"}))
}

func displayArchive(req *web.Request) {
	entries := getEntries()
	entryGroups := flattenEntryGroups(groupEntries(entries))
	w := req.Respond(web.StatusOK, web.HeaderContentType, "text/html; charset=\"utf-8\"")
	params := make(map[string]interface{})
	params["entry_groups"] = entryGroups
	io.WriteString(w, RenderFile("templates/archive.html", params))
}

func main() {
	/* log.Println(splitTags("Nick Carolyn Vanessa Hannah"))
	log.Println(splitTags("\"Hello World\""))
	log.Println(splitTags("\"@Carolyn Gerakines\" #dinner #date"))
	log.Println(splitTags("#meeting \"@Steve McGarrity\" #port #battle.net    \"\"")) */
	db, db_err = mysql.DialTCP("localhost", "root", "asd123", "gobook")
	if db_err != nil {
		log.Println(db_err)
	    os.Exit(1)
	}

	port := ":8080"
	if envPort := os.Getenv("GOBOOK_PORT"); envPort != "" {
		port = envPort
	}

	h := web.FormHandler(10000, false,
		web.NewRouter().
			Register("/", "GET", displayIndex, "POST", createEntry).
			Register("/archive", "GET", displayArchive).
			Register("/static/<path:.*>", "GET", web.DirectoryHandler("./static/", new(web.ServeFileOptions))))
	server.Run(port, h)
}

