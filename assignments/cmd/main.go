package main

import (
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

		fmt.Println("stored in " + path)
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

// TODO: Fix this ugly code
func printAssignmentsPretty(a []assignment.Assignment) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetRowLine(true)
	table.SetHeader([]string{"ΜΑΘΗΜΑ", "ΕΡΓΑΣΙΑ", "ΠΡΟΘΕΣΜΙΑ", "ΕΧΕΙ ΥΠΟΒΛΗΘΕΙ"})
	for _, v := range a {
		calcRemainingTime(v)
		var isSent string
		if v.IsSent {
			isSent = "yes"
		} else {
			isSent = "no"
		}
		table.Append([]string{
			v.Course.Name,
			v.Title,
			v.Deadline.Format("02/01/2006 15:04") + " " + calcRemainingTime(v),
			isSent,
		})
	}
	table.Render()

	return nil
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
