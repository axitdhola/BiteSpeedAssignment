package customerApi

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Axit88/BiteSpeedTask/src/adapters"
	"github.com/Axit88/BiteSpeedTask/src/constants"
	"github.com/Axit88/BiteSpeedTask/src/model"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func GetDbConnection() *sql.DB {
	connection := fmt.Sprintf("%v:%v@tcp(%v)/", constants.DatabaseUsername, constants.DatabasePassword, constants.DatabaseHostURL)
	db, err := sql.Open("mysql", connection)
	if err != nil {
		log.Fatalf("db connection failure: %v", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf("db ping failure: %v", err)
	}

	_, err = db.Exec(`CREATE DATABASE IF NOT EXISTS CustomerDB`)
	if err != nil {
		log.Fatalf("db init failure: %v", err)
	}

	connection = fmt.Sprintf("%v:%v@tcp(%v)/%v", constants.DatabaseUsername, constants.DatabasePassword, constants.DatabaseHostURL, constants.DatabaseName)
	db, err = sql.Open("mysql", connection)
	if err != nil {
		log.Fatalf("db connection failure: %v", err)
	}

	_, err = db.Exec(`
			CREATE TABLE IF NOT EXISTS contact (
			id INT PRIMARY KEY,
			phoneNumber VARCHAR(255),
			email VARCHAR(255),
			linkedId INT,
			linkPrecedence ENUM('secondary', 'primary'),
			createdAt DATETIME,
			updatedAt DATETIME,
			deletedAt DATETIME
		)
	`)
	if err != nil {
		log.Fatalf("db init failure: %v", err)
	}

	return db
}

func getContactsByPhoneNumberOrEmail(db *sql.DB, phoneNumber string, email string) ([]model.Customer, error) {
	query := `
		SELECT id, phoneNumber, email, linkedId, linkPrecedence, createdAt, updatedAt, deletedAt
		FROM contact
		WHERE phoneNumber = ? OR email = ?
		ORDER BY createdAt
	`

	rows, err := db.Query(query, phoneNumber, email)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contacts []model.Customer
	for rows.Next() {
		var contact model.Customer
		_ = rows.Scan(
			&contact.ID,
			&contact.PhoneNumber,
			&contact.Email,
			&contact.LinkedID,
			&contact.LinkPrecedence,
			&contact.CreatedAt,
			&contact.UpdatedAt,
			&contact.DeletedAt,
		)

		contacts = append(contacts, contact)
	}

	return contacts, nil
}

func AddNewContact(db *sql.DB, phone string, email string, linkPrecedence string, linkedId int) error {
	query := `
		INSERT INTO contact (id, phoneNumber, email, linkedId, linkPrecedence, createdAt, updatedAt, deletedAt)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	_, err := db.Exec(query, GetTotalID(db)+1, phone, email, linkedId, "primary", time.Now(), time.Now(), nil)

	return err
}

func GetTotalID(db *sql.DB) int {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM contact").Scan(&count)
	if err != nil {
		fmt.Println(err)
	}
	return count
}

func ResolveMoreThanOnePrimaryContact(db *sql.DB, phone string, email string) error {
	cnt := 0
	primaryId := 0
	res, _ := getContactsByPhoneNumberOrEmail(db, phone, email)
	for _, c := range res {
		if cnt > 1 && c.LinkPrecedence == "primary" {
			_, err := db.Exec("UPDATE contact SET linkPrecedence = ?, updatedAt = ?, linkedId = ? WHERE id = ?", "secondary", time.Now(), primaryId, c.ID)
			if err != nil {
				return err
			}
		}

		if cnt == 1 {
			primaryId = c.ID
		}
		cnt++
	}
	return nil
}

func isInsert(db *sql.DB, phone string, email string) int {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM contact WHERE (phoneNumber=? OR email=?) AND linkPrecedence=?", phone, email, "primary").Scan(&count)
	if err != nil {
		fmt.Println(err)
	}
	return count
}

func PostRequest(context *gin.Context) {
	db := GetDbConnection()
	var newContact model.Payload

	err := context.ShouldBindJSON(&newContact)
	if err != nil {
		context.JSON(400, gin.H{"error": err.Error()})
		return
	}

	email := newContact.Email
	phoneNumber := newContact.PhoneNumber

	if len(email) == 0 && len(phoneNumber) == 0 {
		context.JSON(http.StatusBadRequest, "Insert At least One Not Null Value")
		return
	}

	res, err := getContactsByPhoneNumberOrEmail(db, phoneNumber, email)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"Failed To Get Data From Database": err})
		return
	}

	apiResponse := adapters.GetReults()

	if len(res) == 0 {
		err := AddNewContact(db, phoneNumber, email, "primary", -1)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"Failed Database Insert": err})
			return
		}

		apiResponse.Emails = append(apiResponse.Emails, email)
		apiResponse.PhoneNumbers = append(apiResponse.PhoneNumbers, phoneNumber)
		apiResponse.PrimaryContactID = GetTotalID(db)

		context.JSON(http.StatusOK, apiResponse)
		return
	}

	if res[0].Email == email && res[0].PhoneNumber == phoneNumber {
		context.IndentedJSON(http.StatusCreated, "Contact Already Exist")
		return
	}

	totalPrimaryId := isInsert(db, phoneNumber, email)
	if totalPrimaryId < 2 {
		err = AddNewContact(db, phoneNumber, email, "secondary", res[0].ID)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"Failed Database Insert": err})
			return
		}
	}

	distinctEmails := make(map[string]int)
	distinctPhone := make(map[string]int)
	apiResponse.PrimaryContactID = res[0].ID

	for i, c := range res {
		if len(c.Email) > 0 {
			distinctEmails[c.Email] = 1
		}
		if len(c.PhoneNumber) > 0 {
			distinctPhone[c.PhoneNumber] = 1
		}
		if i > 0 {
			apiResponse.SecondaryContactIDs = append(apiResponse.SecondaryContactIDs, c.ID)
		}
	}

	if totalPrimaryId < 2 {
		if len(email) > 0 {
			distinctEmails[email] = 1
		}

		if len(phoneNumber) > 0 {
			distinctPhone[phoneNumber] = 1
		}
	}

	for key, _ := range distinctEmails {
		apiResponse.Emails = append(apiResponse.Emails, key)
	}

	for key, _ := range distinctPhone {
		apiResponse.PhoneNumbers = append(apiResponse.PhoneNumbers, key)
	}

	err = ResolveMoreThanOnePrimaryContact(db, phoneNumber, email)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"Failed ResolveMoreThanOnePrimaryContact": err})
	}

	if totalPrimaryId < 2 {
		apiResponse.SecondaryContactIDs = append(apiResponse.SecondaryContactIDs, GetTotalID(db))
	}
	context.IndentedJSON(http.StatusOK, apiResponse)
}
