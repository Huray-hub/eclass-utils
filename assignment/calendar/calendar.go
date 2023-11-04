package calendar

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"time"

	as "github.com/Huray-hub/eclass-utils/assignment"
	ics "github.com/arran4/golang-ical"
)

func Export(a []as.Assignment, baseDomain, targetDirectory string) (filePath string, err error) {
	var buffer *bytes.Buffer
	buffer, err = createCalendar(a, baseDomain)
	if err != nil {
		return
	}

	filePath = filepath.Join(
		targetDirectory,
		fmt.Sprintf("assignments_%s.ics", time.Now().Format("02-01-2006")),
	)

	err = os.WriteFile(filePath, buffer.Bytes(), 0644)
	if err != nil {
		return
	}

	return
}

func createCalendar(
	a []as.Assignment,
	baseDomain string,
) (*bytes.Buffer, error) {
	cal := ics.NewCalendar()
	cal.SetProductId("eclass-utils")
	cal.SetCalscale("GREGORIAN")
	cal.SetName("Προθεσμίες")
	cal.SetDescription("Calendar for eclass' assignments")
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
	if a.Deadline == nil {
		return nil
	}
	event := cal.AddEvent(fmt.Sprintf("%v-%v-%v", "eclass-utils", a.Course.ID, a.ID))
	event.SetCreatedTime(time.Now())
	event.SetDtStampTime(time.Now())
	event.SetModifiedAt(time.Now())
	event.SetStartAt(*a.Deadline)
	event.SetEndAt(*a.Deadline)
	event.SetSummary(fmt.Sprintf("%v: %v", a.Course.Name, a.Title))

	assignmentURL, err := a.PrepareURL(baseDomain)
	if err != nil {
		return err
	}
	description := "https://" + assignmentURL

	if a.IsSent {
		description = description + "\n" + "Έχει σταλεί"
	}
	event.SetDescription(description)
	event.SetURL(assignmentURL)

	return nil
}
