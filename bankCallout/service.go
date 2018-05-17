package bankCallout

import (
	"database/sql"
	_ "github.com/denisenkom/go-mssqldb"
	"context"
	"strings"
	"time"
	"log"
	"os"
	"fmt"
	"github.com/weAutomateEverything/go2hal/firstCall"
)

type bankCallout struct {
}

func NewService() firstCall.Service {
	return bankCallout{}

}

func (bankCallout) GetFirstCall(ctx context.Context) (name string, number string, err error) {
	c := fmt.Sprintf("server=%v;user id=%v;password=%v;encrypt=disable;database=%v", getCalloutDbServer(), getCalloutDbUser(), getCalloutDbPassword(), getCalloutDBSchema())
	log.Println(c)
	db, err := sql.Open("mssql", c)
	if err != nil {
		return
	}
	defer db.Close()

	s := time.Now().Format("2006-01-02 15:04:05")
	log.Println(s)

	q := fmt.Sprintf("SELECT FirstName,Primary1st FROM CalloutListingFirstCall where DateFrom < '%v' and DateTo > '%v' and Team = '%v'",s,s,getCalloutGroup())
	stmt, err := db.Query(q)
	defer stmt.Close()

	if err != nil {
		return

	}
	stmt.Next()
	err = stmt.Scan(&name, &number)
	if err != nil {
		return
	}

	number = strings.Replace(number, " ", "", -1)
	number = strings.Replace(number, "-", "", -1)
	number = strings.Replace(number, "0", "+27", 1)

	return
}

func getCalloutDbServer() string {
	return os.Getenv("CALLOUT_DB_SERVER")
}

func getCalloutDbUser() string {
	return os.Getenv("CALLOUT_DB_USER")
}

func getCalloutDbPassword() string {
	return os.Getenv("CALLOUT_DB_PASSWORD")
}

func getCalloutDBSchema() string {
	return os.Getenv("CALLOUT_DB_SCHEMA")
}

func getCalloutGroup() string {
	return os.Getenv("CALLOUT_GROUP")
}
