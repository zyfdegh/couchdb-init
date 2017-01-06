package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/cesanta/docker_auth/auth_server/authz"
	"github.com/cesanta/docker_auth/auth_server/utils"
	"github.com/segmentio/pointer"
	"github.com/zemirco/couchdb"
)

type DocACLEntry struct {
	couchdb.Document
	Seq     *int
	Match   authz.MatchConditions
	Actions *[]string
	Comment *string
}

type DocAccount struct {
	couchdb.Document
	Email    *string `yaml:"email,omitempty" json:"email,omitempty"`
	Password *string `yaml:"password,omitempty" json:"password,omitempty"`
	Username *string `yaml:"username,omitempty" json:"username,omitempty"`
}

// just some helper function
func check(err error) {
	if err != nil {
		log.Printf("error: %v\n", err)
	}
}

// StringInSlice return true if list contains a
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func main() {
	// create a new client
	client, err := couchdb.NewClient("http://127.0.0.1:5984/")
	check(err)
	// get some information about your CouchDB
	info, err := client.Info()
	check(err)
	log.Println(info)

	// create a database
	log.Println(">>> Creating db acl")
	allDBs, err := client.All()
	check(err)
	dbName := "acl"
	if !utils.StringInSlice(dbName, allDBs) {
		_, err = client.Create(dbName)
		check(err)
	}

	// use new database
	db := client.Use(dbName)

	// add acl design doc
	log.Println(">>> Adding design doc _design/acl")
	aclDesignDoc := &couchdb.DesignDocument{
		Document: couchdb.Document{
			ID: "_design/acl",
		},
		Language: "javascript",
		Views: map[string]couchdb.DesignDocumentView{
			"getBySeq": couchdb.DesignDocumentView{
				Map: `function(doc) { emit(doc.seq, doc); }`,
			},
		},
	}
	_, err = db.Post(aclDesignDoc)
	check(err)

	// add documents
	aclArr := [][]byte{
		[]byte(`{"seq": 10, "match" : {"account" : "admin"}, "actions" : ["*"], "comment" : "Admin has full access to everything."}`),
		[]byte(`{"seq": 20, "match" : {"account" : "test", "name" : "test-*"}, "actions" : ["*"], "comment" : "User \"test\" has full access to test-* images but nothing else. (1)"}`),
		[]byte(`{"seq": 30, "match" : {"account" : "test"}, "actions" : [], "comment" : "User \"test\" has full access to test-* images but nothing else. (2)"}`),
		[]byte(`{"seq": 40, "match" : {"account" : "/.+/"}, "actions" : ["pull"], "comment" : "All logged in users can pull all images."}`),
		[]byte(`{"seq": 50, "match" : {"account" : "/.+/", "name" : "${account}/*"}, "actions" : ["*"], "comment" : "All logged in users can push all images that are in a namespace beginning with their name"}`),
		[]byte(`{"seq": 60, "match" : {"account" : "", "name" : "hello-world"}, "actions" : ["pull"], "comment" : "Anonymous users can pull \"hello-world\"."}`),
	}

	for _, a := range aclArr {
		aclDoc := &DocACLEntry{}
		err = json.Unmarshal(a, aclDoc)
		check(err)

		result, errPost := db.Post(aclDoc)
		check(errPost)
		log.Println(result)
	}

	// query
	log.Println(">>> Query db acl")
	db = client.Use(dbName)
	view := db.View(dbName)
	queryParams := couchdb.QueryParameters{
		Key: pointer.String(fmt.Sprintf("%q", "10")),
	}
	res, err := view.Get("getBySeq", queryParams)
	check(err)
	if res != nil {
		log.Println(res)
		for _, r := range res.Rows {
			log.Println(r.Value)
		}
	}

	// create a database
	log.Println(">>> Creating db user")
	allDBs, err = client.All()
	check(err)
	dbName = "user"
	if !utils.StringInSlice(dbName, allDBs) {
		_, err = client.Create(dbName)
		check(err)
	}

	// use new database
	db = client.Use(dbName)

	// add user design doc
	log.Println(">>> Adding design doc _design/user")
	userDesignDoc := &couchdb.DesignDocument{
		Document: couchdb.Document{
			ID: "_design/user",
		},
		Language: "javascript",
		Views: map[string]couchdb.DesignDocumentView{
			"getByEmail": couchdb.DesignDocumentView{
				Map: `function(doc) { emit(doc.email, doc); }`,
			},
		},
	}
	_, err = db.Post(userDesignDoc)
	check(err)

	// add users
	log.Println(">>> Adding users")
	uArr := [][]byte{
		// admin badmin
		[]byte(`{"email" : "admin@email.com","password" : "secret123","username":"admin"}`),
		// zhang password
		[]byte(`{"email" : "zhang@email.com","password" : "password","username":"zhang"}`),
		[]byte(`{"email" : "zhang@email.com","password2" : "password","username":"zhang2"}`),
	}
	for _, u := range uArr {
		accoutDoc := &DocAccount{}
		err = json.Unmarshal(u, accoutDoc)
		check(err)

		result, errPost := db.Post(accoutDoc)
		check(errPost)
		log.Println(result)
	}

	// query
	log.Println(">>> Query db user")
	db = client.Use(dbName)
	view = db.View(dbName)
	queryParams = couchdb.QueryParameters{
		Key: pointer.String(fmt.Sprintf("%q", "admin@email.com")),
	}
	res, err = view.Get("getByEmail", queryParams)
	check(err)
	if res != nil {
		log.Println(res)
		for _, r := range res.Rows {
			// {_id: "2695ab910092c67a1eb6ffb0af00412e", _rev: "1-4e9b44eafd0fbc6bd4a04c0f73868627", email: "admin@email.com", password: "secret123", username: "admin"}
			log.Println(r.Value)
			// couchdbUser := r.Value.(DocAccount)
			// log.Println(couchdbUser.Password)
		}
	}
}
