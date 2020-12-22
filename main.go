package main

import (
	"database/sql"
	"fmt"
	"math"
	"os"
	"os/user"
	"path"
	"strings"

	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"

	_ "github.com/mattn/go-sqlite3"
)

var q1 = `
	SELECT C.ROWID, C.chat_identifier, M.ROWID, M.text, datetime(M.date/1000000000 + strftime("%s", "2001-01-01") ,"unixepoch","localtime")  as date_utc, M.is_from_me
	FROM chat C
	INNER JOIN chat_handle_join H
		ON H.chat_id = C.ROWID
	INNER JOIN message M
		ON M.handle_id = H.handle_id
	WHERE C.chat_identifier=?
	ORDER by M.date;
`

// Row is a result from sqlite3
type Row struct {
	ID             int
	RowID          int
	chatIdentifier string
	text           sql.NullString
	date           sql.NullString
	fromMe         string
}

func main() {
	fmt.Println("Starting... 🚀")
	m := pdf.NewMaroto(consts.Portrait, consts.Letter)
	m.SetBorder(true)
	m.SetPageMargins(1, 1, 1)

	usr, _ := user.Current()
	dir := path.Join(usr.HomeDir, "/Library/Messages/chat.db")
	fmt.Printf("Using DB %s \n", dir)

	if len(os.Args) == 1 {
		fmt.Println("Please provide the recipient phone number or email: '+14151231234'")
		os.Exit(1)
	}
	recipient := os.Args[1]

	db, err := sql.Open("sqlite3", dir)
	defer db.Close()
	if err != nil {
		panic(err)
	}

	stmt, err := db.Prepare(q1)
	defer stmt.Close()

	rows, err := stmt.Query(recipient)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	for rows.Next() {
		r := Row{}
		err = rows.Scan(&r.ID, &r.chatIdentifier, &r.RowID, &r.text, &r.date, &r.fromMe)
		if err != nil {
			fmt.Println(err)
		}

		t := strings.TrimSpace(r.text.String)
		align := consts.Left
		if r.fromMe == "1" {
			align = consts.Right
		}

		h := math.RoundToEven(float64(len(t)/60) + 1)
		h = h * 4
		m.Row(float64(h), func() {
			m.Col(2, func() {
				m.Text(r.date.String)
			})
			m.Col(10, func() {
				m.Text(t, props.Text{
					Align: align,
				})
			})
		})

	}
	err = m.OutputFileAndClose("output.pdf")
	if err != nil {
		fmt.Println("Could not save PDF:", err)
	}
	fmt.Println("All done! 🎉")
}
