package main

import (
	"database/sql"
	"fmt"
	"os"
	"regexp"
	"strings"

	_ "github.com/lib/pq"
)

type Lexer struct {
	Token string
}

type LexersInterface interface {
	Handler(string) interface{}
}

func NewLexer(token string) *Lexer {

	return &Lexer{
		Token: token,
	}
}

var lexers = map[string]func(string) any{
	"$$dsn": Connect,
	"sql":   Handler2,
}

var db *sql.DB

func Connect(line string) interface{} {

	split := strings.Split(line, ":=")

	pattern := regexp.MustCompile(`(?P<conn>\w+)>(?P<dsn>.+)`)

	dbType := pattern.FindStringSubmatch(split[1])

	db, _ = sql.Open(dbType[1], dbType[2])

	return nil
}

func Handler2(line string) interface{} {

	pattern := regexp.MustCompile(`sql\s+(.*)`)

	dbType := pattern.FindStringSubmatch(line)

	if db == nil {
		return nil
	}

	row := db.QueryRow(dbType[1])

	if row == nil {

		return nil
	}
	var col string

	row.Scan(&col)

	patternSelector := regexp.MustCompile(`\w+:`)

	s := patternSelector.FindStringSubmatch(line)

	selector := fmt.Sprintf("%s%s;", s[0], col)

	return selector
}

func main() {

	file, err := os.ReadFile("main.csql")

	if err != nil {

		panic(err)
	}

	str := string(file)

	lines := strings.Split(str, "\n")

	var varCssFile = ``
line:
	for _, line := range lines {

		lineTokens := strings.Split(line, " ")

		for _, token := range lineTokens {

			if lexers[token] != nil {
				res := lexers[token](line)
				if res == nil {
					continue line
				}
				if res != nil {
					varCssFile += res.(string)
					continue line
				}
			}
		}
		varCssFile += line

	}

	f, _ := os.Create("main.css")

	f.WriteString(varCssFile)

}
