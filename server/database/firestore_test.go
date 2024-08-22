package database_test

import (
	"os"
	"testing"
	"time"
	"url-trigger/database"
)

func TestUploadSimpleData(t *testing.T) {
	db := database.NewDataBase()
	string1 := "This is a test data OTP 1234 time is " + time.Now().String()
	err := db.InsertOtpToTable("otps", string1)

	if err != nil {
		t.Error("Error while inserting data to table")
	}
	string2 := "This is a test data OTP 123456 " + time.Now().String()
	db.InsertOtpToTable("otps", string2)
	val, err := db.RetrieveLatest("otps")
	if err != nil {
		t.Error("Error while retrieving data from table")
	}
	if val[0] != string2 {
		t.Error("Data mismatch")
	}
	if val[1] != string1 {
		t.Error("Data mismatch")
	}
}

func TestInsertAllowedUserCOllection(t *testing.T) {
	t.Skip("Skipping this test")
	db := database.NewDataBase()
	db.DeleteCollection("allowed_users")
	allowedUsers := []string{"dummyemail@email.com", "dummyemail1@email.com"}
	for _, user := range allowedUsers {

		db.InsDataToTable("allowed_users", map[string]interface{}{
			"username": user,
		})
	}
	userNames, err := db.AllowedUsers()
	if err != nil {
		t.Error("Error while retrieving data from table")
	}
	if len(userNames) != 4 {
		t.Error("Data mismatch")
	}
}

func TestRetrieveData(t *testing.T) {
	db := database.NewDataBase()
	val, err := db.RetrieveLatest("otps")

	if err != nil && err.Error() != "no data found" {
		t.Error("Error while retrieving data from table")
	}
	if len(val) == 0 {
		t.Error("No data found")
	}
}

func TestVerifyAccessToken(t *testing.T) {
	db := database.NewDataBase()
	os.Setenv("testing", "1")
	_, _, err := db.AuthenticateBearer("")
	if err != nil {
		t.Errorf("Error while verifying ID token: %v", err)
	}
}
