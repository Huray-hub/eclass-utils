package calendar

import (
	"bytes"
	"fmt"
	"os"
	"time"

	as "github.com/Huray-hub/eclass-utils/deadlines/assignments"
	ics "github.com/arran4/golang-ical"
)

func Export(a []as.Assignment, baseDomain string) (string, error) {
	buffer, err := createCalendar(a, baseDomain)
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

func createCalendar(
	a []as.Assignment,
	baseDomain string,
) (*bytes.Buffer, error) {
	cal := ics.NewCalendar()
	cal.SetProductId("eclass-utils")
	cal.SetCalscale("GREGORIAN")
	// cal.SetMethod(ics.MethodRefresh)
	cal.SetName("Προθεσμίες")
	cal.SetDescription("Calendar for eclass assignments' deadlines")
	cal.SetColor("red")

	for _, v := range a {
		err := addEvent(v, cal, baseDomain)
		if err != nil {
			return nil, err
		}
	}

	b := bytes.NewBufferString("")
	err := cal.SerializeTo(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func addEvent(a as.Assignment, cal *ics.Calendar, baseDomain string) error {
	event := cal.AddEvent(fmt.Sprintf("%v-%v-%v", "eclass-utils", a.Course.ID, a.ID))
	event.SetCreatedTime(time.Now())
	event.SetDtStampTime(time.Now())
	event.SetModifiedAt(time.Now())
	event.SetStartAt(a.Deadline)
	event.SetEndAt(a.Deadline)
	event.SetSummary(fmt.Sprintf("%v: %v", a.Course.Name, a.Title))

	assignmentURL, err := a.PrepareURL(baseDomain)
	if err != nil {
		return err
	}
	assignmentURL = "https://" + assignmentURL
	event.SetDescription(assignmentURL)
	event.SetURL(assignmentURL)

	return nil
}
