package bankCallout

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/pkg/errors"
	"github.com/weAutomateEverything/go2hal/firstCall"
	"github.com/weAutomateEverything/go2hal/telegram"
	"log"
	"os"
	"strings"
	"time"
)

type Service interface {
	setGroup(ctx context.Context, chat uint32, group string) (name string, number string, err error)
	getGroup(ctx context.Context, chat uint32) (string, error)
}

func NewService(store Store, telegram telegram.Service, telegramStore telegram.Store) firstCall.CalloutFunction {
	return bankCallout{
		store:         store,
		telegram:      telegram,
		telegramStore: telegramStore,
	}

}

type bankCallout struct {
	store         Store
	telegram      telegram.Service
	telegramStore telegram.Store
}

func (s bankCallout) Escalate(ctx context.Context, count int, chat uint32) (name string, number string, err error) {
	if count < 3 {
		return s.GetFirstCallDetails(ctx, chat)
	}
	if count < 6 {
		return s.GetSecondCallDetails(ctx, chat)
	}
	if count < 9 {
		return s.GetManagementDetails(ctx, chat)
	}
	err = errors.New("Unable to escalate any further. Giving up. ")
	return

}

func (s bankCallout) Configured(chat uint32) bool {
	_, err := s.store.getCalloutGroup(chat)
	return err != nil
}
func (s bankCallout) GetSecondCallDetails(ctx context.Context, chat uint32) (name string, number string, err error) {

	group, err := s.store.getCalloutGroup(chat)
	if err != nil {
		return
	}
	return s.getCallout(s.getCalloutString(group, "SecondName", "Primary2nd", "CalloutListingSecondCall"))

}

func (s bankCallout) GetFirstCallDetails(ctx context.Context, chat uint32) (name string, number string, err error) {

	group, err := s.store.getCalloutGroup(chat)
	if err != nil {
		return
	}
	return s.getCallout(s.getCalloutString(group, "FirstName", "Primary1st", "CalloutListingFirstCall"))

}

func (s bankCallout) GetManagementDetails(ctx context.Context, chat uint32) (name string, number string, err error) {

	group, err := s.store.getCalloutGroup(chat)
	if err != nil {
		return
	}
	return s.getCallout(s.getManagementCalloutString(group))

}

func (s bankCallout) setGroup(ctx context.Context, chat uint32, group string) (name string, number string, err error) {
	name, number, err = s.getCallout(s.getCalloutString(group, "FirstName", "Primary1st", "CalloutListingFirstCall"))
	if err != nil {
		return
	}
	err = s.store.setCallout(chat, group)

	if err != nil {
		return
	}
	g, err := s.telegramStore.GetRoomKey(chat)
	if err != nil {
		return
	}
	s.telegram.SendMessage(ctx, g, fmt.Sprintf("Your callout group has been successfully changed to %v. On firstcall is %v, %v", group, name, number), 0)

	return
}

func (s bankCallout) getGroup(ctx context.Context, chat uint32) (string, error) {
	return s.store.getCalloutGroup(chat)
}

func (s bankCallout) getCallout(query string) (name string, number string, err error) {
	name, number = "", ""
	c := fmt.Sprintf("server=%v;user id=%v;password=%v;encrypt=disable;database=%v", getCalloutDbServer(), getCalloutDbUser(), getCalloutDbPassword(), getCalloutDBSchema())
	log.Println(c)
	db, err := sql.Open("mssql", c)
	if err != nil {
		return
	}
	defer db.Close()

	log.Println(s)

	stmt, err := db.Query(query)

	if err != nil {
		return

	}
	defer stmt.Close()

	stmt.Next()
	err = stmt.Scan(&name, &number)
	if err != nil {
		return
	}

	number = strings.Replace(number, " ", "", -1)
	number = strings.Replace(number, "-", "", -1)
	if strings.HasPrefix(number, "0") {
		number = strings.Replace(number, "0", "+27", 1)

	}

	return
}

func (s bankCallout) getManagementCalloutString(group string) string {
	t := time.Now().Format("2006-01-02 15:04:05")
	return fmt.Sprintf("SELECT ManagementFirstName,Primary1stManagement FROM vCalloutManagement where ((DateFrom < '%v' and DateTo > '%v') or Always = true) and Team = '%v'", t, t, group)
}

func (s bankCallout) getCalloutString(group string, nameFields, NumberField, table string) string {
	t := time.Now().Format("2006-01-02 15:04:05")
	return fmt.Sprintf("SELECT %v,%v FROM %v where DateFrom < '%v' and DateTo > '%v' and Team = '%v'", nameFields, NumberField, table, t, t, group)
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
