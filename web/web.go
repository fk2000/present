package web

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"

	"github.com/hakobe/present/accesslogs"
	"github.com/hakobe/present/entries"
	slackOutgoing "github.com/hakobe/present/slack/outgoing"
)

var bind string = ":" + os.Getenv("PORT")

func Start(db *sql.DB) chan *slackOutgoing.Op {
	op := make(chan *slackOutgoing.Op, 1000)

	http.HandleFunc(
		"/hook",
		func(rw http.ResponseWriter, r *http.Request) {
			slackOutgoing.Handle(op, rw, r)
		},
	)

	http.HandleFunc(
		"/upcommings",
		func(rw http.ResponseWriter, r *http.Request) {
			es, err := entries.Upcommings(db)
			if err != nil {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
				return
			}
			tmpl, err := template.New("upcommings").Parse(`
<html>
<body>
  <table>
	<tr>
		<th>Tag</th><th>Title</th><th>Url</th>
	</tr>
	{{ range . }}
	<tr>
		<td>{{.Tag}}</td><td>{{.Title}}</td><td>{{.Url}}</td>
	</tr>
	{{ end }}
  </table>
</body>
</html>
		`)
			if err != nil {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
				return
			}

			tmpl.Execute(rw, es)
		},
	)

	http.HandleFunc(
		"/entry/",
		func(rw http.ResponseWriter, r *http.Request) {
			matches := regexp.MustCompile(`/entry/(\d+)`).FindStringSubmatch(r.URL.Path)
			if !(matches != nil && matches[1] != "") {
				http.Error(rw, "Invalid URL", http.StatusBadRequest)
				return
			}

			var id int
			var err error
			if id, err = strconv.Atoi(matches[1]); err != nil {
				http.Error(rw, err.Error(), http.StatusNotFound)
				return
			}

			entry, err := entries.Find(db, id)
			if err != nil {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
				return
			}
			accesslogs.Access(db, id)
			tmpl, err := template.New("entry").Parse(`
<html>
<body></body>
<head>
  <meta http-equiv="refresh" content="0;URL=data:text/html,%3Cmeta%20http-equiv%3D%22refresh%22%20content%3D%220%3BURL%3D{{.}}%22%3E"></meta>
</head>
</html>
		`)
			if err != nil {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
				return
			}

			tmpl.Execute(rw, template.HTML(url.QueryEscape(entry.Url())))
		},
	)

	go func() {
		log.Printf("Starting slack webhook on \"%s\"\n", bind)
		err := http.ListenAndServe(bind, nil)
		if err != nil {
			log.Fatalf("ListenAndServe: %v", err)
		}
	}()

	return op
}
