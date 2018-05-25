package bankCallout

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/weAutomateEverything/go2hal/firstCall"
	"github.com/weAutomateEverything/go2hal/telegram"
	"log"
	"os"
	"strings"
	"time"
)

type bankCallout struct {
	store         Store
	telegram      telegram.Service
	telegramStore telegram.Store
}

type Service interface {
	setGroup(ctx context.Context, chat uint32, group string) (name string, number string, err error)
	getGroup(ctx context.Context, chat uint32) (string, error)
}

func NewService(store Store, telegram telegram.Service, telegramStore telegram.Store) firstCall.Service {
	return bankCallout{
		store:         store,
		telegram:      telegram,
		telegramStore: telegramStore,
	}

}

func (s bankCallout) GetFirstCall(ctx context.Context, chat uint32) (name string, number string, err error) {

	group, err := s.store.getCalloutGroup(chat)
	if err != nil {
		return
	}
	return s.getCallout(group)

}

func (s bankCallout) setGroup(ctx context.Context, chat uint32, group string) (name string, number string, err error) {
	name, number, err = s.getCallout(group)
	if err != nil {
		return
	}
	err = s.store.setCallout(chat, group)

	if err != nil {
		g, err := s.telegramStore.GetRoomKey(chat)
		if err != nil {
			s.telegram.SendMessage(ctx, g, fmt.Sprintf("Your callout group has been successfully changed to %v. On firstcall is %v, %v", group, name, number), 0)

		}
	}

	return
}

func (s bankCallout) getGroup(ctx context.Context, chat uint32) (string, error) {
	return s.store.getCalloutGroup(chat)
}

func (s bankCallout) getCallout(group string) (name string, number string, err error) {
	c := fmt.Sprintf("server=%v;user id=%v;password=%v;encrypt=disable;database=%v", getCalloutDbServer(), getCalloutDbUser(), getCalloutDbPassword(), getCalloutDBSchema())
	log.Println(c)
	db, err := sql.Open("mssql", c)
	if err != nil {
		return
	}
	defer db.Close()

	t := time.Now().Format("2006-01-02 15:04:05")
	log.Println(s)

	q := fmt.Sprintf("SELECT FirstName,Primary1st FROM CalloutListingFirstCall where DateFrom < '%v' and DateTo > '%v' and Team = '%v'", t, t, group)
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
