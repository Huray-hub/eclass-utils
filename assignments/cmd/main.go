package main

import (
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"time"

	"github.com/Huray-hub/eclass-utils/assignments/assignment"
	"github.com/Huray-hub/eclass-utils/assignments/calendar"
	"github.com/Huray-hub/eclass-utils/assignments/config"
	"github.com/olekukonko/tablewriter"
)

func init() {
	homeCache, err := os.UserCacheDir()
	if err != nil {
        log.Fatal(err.Error())
	}

    path:= homeCache+"/eclass-utils"
    if _, err2 := os.Stat(path); errors.Is(err2, os.ErrNotExist) {
		err3 := os.Mkdir(path, os.ModePerm)
		if err3 != nil {
			log.Fatal(err3)
		}
	}

	file, err := os.OpenFile(
		path + "/assignments.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(file)
}

func main() {
	opts, creds, err := config.Import()
	if err != nil {
		log.Fatal(err.Error())
	}

	a, err := assignment.Get(opts, creds)
	if err != nil {
		log.Fatal(err.Error())
	}

	err = printAssignments(a, opts.PlainText)
	if err != nil {
		log.Fatal(err.Error())
	}

	if opts.ExportICS {
		path, err := calendar.Export(a, opts.BaseDomain)
		if err != nil {
			log.Fatal(err.Error())
		}

        fmt.Printf("stored in\n%v", path)
	}
}

func printAssignments(a []assignment.Assignment, plainText bool) error {
	if plainText {
		return printAssignmentsPlain(a)
	}

	return printAssignmentsPretty(a)
}

func printAssignmentsPlain(a []assignment.Assignment) error {
	for _, v := range a {
		_, err := fmt.Println(v.String())
		if err != nil {
			return err
		}
	}
	return nil
}

func printAssignmentsPretty(a []assignment.Assignment) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetRowLine(true)
	table.SetHeader([]string{"ΜΑΘΗΜΑ", "ΕΡΓΑΣΙΑ", "ΠΡΟΘΕΣΜΙΑ", "ΥΠΟΒΛΗΘΗΚΕ"})
	appendToTable(a, table)
	table.Render()

	return nil
}

func appendToTable(a []assignment.Assignment, table *tablewriter.Table) {
	for _, v := range a {
		calcRemainingTime(v)
		var isSent string
		if v.IsSent {
			isSent = "Ναι"
		} else {
			isSent = "Όχι"
		}
		table.Append([]string{
			v.Course.Name,
			v.Title,
			v.Deadline.Format("02/01/2006 15:04") + " " + calcRemainingTime(v),
			isSent,
		})
	}
}

func calcRemainingTime(a assignment.Assignment) string {
	t := time.Until(a.Deadline)

	switch {
	case t < 0:
		return "(Έληξε)"
	case t.Hours()/24 > 0:
		return "(" + fmt.Sprint(math.Floor(t.Hours()/24)) + " μέρες)"
	case t.Minutes()/60 > 0:
		return "(" + fmt.Sprint(math.Floor(t.Hours())) + " ώρες)"
	default:
		return "(" + fmt.Sprint(math.Floor(t.Minutes())) + " λεπτά)"
	}
}
