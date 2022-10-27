package internal

import (
	"bytes"
	"fmt"
	"os"
	"time"

	ics "github.com/arran4/golang-ical"
)

func ExportICS(a []Assignment) (string, error) {
	buffer, err := createCalendar(a)
	if err != nil {
		return "", err
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	t := time.Now()
	path := fmt.Sprintf(
		"%v/Documents/assignments_%v%v%v.ics",
		home, t.Day(), int(t.Month()), t.Year())

	err = os.WriteFile(path, buffer.Bytes(), 0644)
	if err != nil {
		return "", err
	}

	return path, nil
}

func createCalendar(a []Assignment) (*bytes.Buffer, error) {
	// cal, err := ParseCalendar(strings.NewReader(input))
	cal := ics.NewCalendar()
	cal.SetMethod(ics.MethodRequest)
	// prop:=cal.CalendarProperties
    panic("not implemented")
	// cal.SetName("Προθεσμίες", ))
	cal.SetColor("red", nil)

	addCalEvent(a[len(a)-1], cal)
	// for _, v := range a {
	// 	addCalEvent(v, cal)
	// }

	b := bytes.NewBufferString("")
	err := cal.SerializeTo(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func addCalEvent(a Assignment, cal *ics.Calendar) {
	// event := cal.AddEvent(fmt.Sprintf("id@domain", p.SessionKey.IntID()))
	event := cal.AddEvent("")
	event.SetCreatedTime(time.Now())
	event.SetDtStampTime(time.Now())
	event.SetModifiedAt(time.Now())
	event.SetStartAt(time.Now())
	event.SetEndAt(a.Deadline)
	event.SetSummary(a.Title)
	// event.SetLocation("Address")
	// event.SetDescription("Description")
	// event.SetURL("https://URL/")
	// event.AddRrule(fmt.Sprintf("FREQ=YEARLY;BYMONTH=%d;BYMONTHDAY=%d", time.Now().Month(), time.Now().Day()))
	// event.SetOrganizer("sender@domain", ics.WithCN("This Machine"))
}
