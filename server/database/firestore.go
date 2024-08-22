package database

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"golang.org/x/exp/slices"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

func ConnectToFirestore() (*firestore.Client, context.Context, error) {
	// Replace with your service account key file path
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return nil, nil, err
	}
	credentialsFile := "serviceAccount.json"
	if strings.Contains(pwd, "database") {
		credentialsFile = "../serviceAccount.json"
	}
	opt := option.WithCredentialsFile(credentialsFile)
	ctx := context.Background()

	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, nil, fmt.Errorf("error initializing app: %v", err)
	}
	client, err := app.Firestore(ctx)
	if err != nil {
		return nil, nil, err
	}
	return client, ctx, nil
}

type DataBase struct {
	dbType string `json:"type"`
}

func NewDataBase() *DataBase {
	return &DataBase{dbType: "firestore"}
}
func (db *DataBase) AuthenticateBearer(authHeader string) (string, string, error) {
	t := os.Getenv("testing")
	if t != "" {
		return "test", "test", nil
	}
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", "", fmt.Errorf("invalid authorization header")
	}
	authString := authHeader[7:]

	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		return "", "", err
	}
	credentialsFile := "serviceAccount.json"
	if strings.Contains(pwd, "database") {
		credentialsFile = "../serviceAccount.json"
	}
	opt := option.WithCredentialsFile(credentialsFile)
	ctx := context.Background()
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return "", "", fmt.Errorf("error initializing app: %v", err)
	}

	client, err := app.Auth(ctx)
	if err != nil {
		return "", "", err
	}

	authToken, err := client.VerifyIDToken(ctx, authString)
	if err != nil {
		fmt.Printf("error verifying ID token: %v\n", err)
		return "", "", err
	}
	userRecord, err := client.GetUser(ctx, authToken.UID)
	if err != nil {
		return "", "", err
	}
	//we wanna check if its in allowed users
	allowedUsers, err := db.AllowedUsers()
	if err != nil {
		return "", "", err
	}
	if !slices.Contains(allowedUsers, userRecord.Email) {
		return "", "", errors.New("user not allowed")
	}
	return userRecord.Email, userRecord.UID, nil

}
func (db *DataBase) InsDataToTable(table string, data map[string]interface{}) error {
	client, ctx, err := ConnectToFirestore()
	if err != nil {
		return err
	}
	defer client.Close()
	insertData := data
	insertData["created_at"] = time.Now().UTC()
	_, _, err = client.Collection(table).Add(ctx, insertData)
	if err != nil {
		log.Fatalf("Failed adding data: %v", err)
	}
	return nil
}
func (db *DataBase) InsertOtpToTable(table string, data string) error {
	err := db.InsDataToTable(table, map[string]interface{}{
		"otp": data,
	})
	return err
}

func (db *DataBase) RetrieveLatest(table string) ([]string, error) {
	limit := 10
	client, ctx, err := ConnectToFirestore()
	if err != nil {
		return nil, err
	}
	defer client.Close()
	otps := client.Collection(table)
	var lastOtps []string
	iter := otps.OrderBy("created_at", firestore.Desc).Where("created_at", ">", time.Now().Add(-time.Hour).UTC()).Limit(100).Documents(ctx)
	counter := 0
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate: %v", err)
		}
		lastOtps = append(lastOtps, doc.Data()["otp"].(string))
		if counter > int(limit) {
			doc.Ref.Delete(ctx)
		}
		counter++
	}

	db.deleteOldOtps(ctx, client, time.Hour)
	if len(lastOtps) > 0 {
		return lastOtps, nil
	}
	return nil, fmt.Errorf("no data found")
}

func (db *DataBase) AllowedUsers() ([]string, error) {
	userNames := []string{}
	client, ctx, err := ConnectToFirestore()
	if err != nil {
		return userNames, err
	}
	defer client.Close()
	tableData := client.Collection("allowed_users")
	iter := tableData.Limit(100).Documents(ctx)

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		userNames = append(userNames, doc.Data()["username"].(string))
	}
	return userNames, nil
}

func (db *DataBase) DeleteCollection(table string) error {
	client, ctx, err := ConnectToFirestore()
	if err != nil {
		return err
	}
	defer client.Close()
	tableData := client.Collection(table)
	iter := tableData.Where("created_at", "<", time.Now()).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		doc.Ref.Delete(ctx)
	}
	if err != nil {
		return err
	}
	return nil
}

func (db *DataBase) deleteOldOtps(ctx context.Context, client *firestore.Client, olderThan time.Duration) {
	otps := client.Collection("otps")
	iter := otps.Where("created_at", "<", time.Now().Add(-olderThan).UTC()).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		doc.Ref.Delete(ctx)
	}
}
