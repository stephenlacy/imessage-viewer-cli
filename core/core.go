package core

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"os/user"
	"path"
	"strings"
	"time"

	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"

	_ "github.com/mattn/go-sqlite3"
)

var SMSQuery = `
	SELECT C.ROWID, C.chat_identifier, M.ROWID, M.text, datetime(M.date/1000000000 + strftime("%s", "2001-01-01") ,"unixepoch","localtime")  as date_utc, M.is_from_me
	FROM chat C
	INNER JOIN chat_handle_join H
		ON H.chat_id = C.ROWID
	INNER JOIN message M
		ON M.handle_id = H.handle_id
	WHERE C.chat_identifier=?
	ORDER by M.date;
`

var BackupQuery = `
SELECT fileID, domain, relativePath from Files;
`

// Row is a result from sqlite3
type Row struct {
	ID             int
	RowID          int
	ChatIdentifier string
	Text           sql.NullString
	Date           sql.NullString
	FromMe         string
}

func Process(dir string, recipient string) {

	m := pdf.NewMaroto(consts.Portrait, consts.Letter)
	m.SetBorder(true)
	m.SetPageMargins(1, 1, 1)

	db, err := sql.Open("sqlite3", dir)
	defer db.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	stmt, err := db.Prepare(SMSQuery)
	defer stmt.Close()

	rows, err := stmt.Query(recipient)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()

	for rows.Next() {
		r := Row{}
		err = rows.Scan(&r.ID, &r.ChatIdentifier, &r.RowID, &r.Text, &r.Date, &r.FromMe)
		if err != nil {
			fmt.Println(err)
		}
		WritePDFRow(r, m)

	}
	saved := fmt.Sprintf("%s.pdf", recipient)
	err = m.OutputFileAndClose(saved)
	if err != nil {
		fmt.Println("Could not save PDF:", err)
	}
	fmt.Printf("All done! 🎉 \nSaved to: %s\n", saved)

}

func HandleiOSBackups(_ string, recipient string) {
	var rootDir = "/Library/Application Support/MobileSync/Backup/"
	var smsDB = "Library/SMS/sms.db"
	usr, _ := user.Current()

	dir := path.Join(usr.HomeDir, rootDir)
	files, err := ioutil.ReadDir(dir)

	// Get latest backup db
	latestBackupDir := ""
	var modTime time.Time
	for _, f := range files {
		if f.IsDir() {
			if f.ModTime().After(modTime) {
				modTime = f.ModTime()
				latestBackupDir = f.Name()
			}
		}
	}

	dir = path.Join(dir, latestBackupDir, "Manifest.db")
	fmt.Printf("Using Manifest.db %s \n", dir)

	if len(os.Args) == 1 {
		fmt.Println("Please provide the recipient phone number or email: '+14151231234'")
		os.Exit(1)
	}

	db, err := sql.Open("sqlite3", dir)
	defer db.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	rows, err := db.Query(BackupQuery)
	if err != nil {
		fmt.Println(err)
		return
	}

	var dbHash string

	var fileID string
	var domain string
	var filePath string

	for rows.Next() {
		rows.Scan(&fileID, &domain, &filePath)
		if filePath == smsDB {
			dbHash = fileID
		}
	}

	pth := path.Join(usr.HomeDir, rootDir, latestBackupDir, dbHash[0:2], dbHash)

	Process(pth, recipient)
}

func WritePDFRow(r Row, m pdf.Maroto) {
	t := strings.TrimSpace(r.Text.String)
	align := consts.Left
	if r.FromMe == "1" {
		align = consts.Right
	}

	h := math.RoundToEven(float64(len(t)/60) + 1)
	h = h * 4

	m.Row(float64(h), func() {
		m.Col(2, func() {
			m.Text(r.Date.String)
		})
		m.Col(10, func() {
			m.Text(t, props.Text{
				Align: align,
			})
		})
	})
}
