package main

import (
	"fmt"
	"log"
	"math"
	"os"
	"time"

	in "github.com/Huray-hub/eclass-utils/deadlines/internal"
	dl "github.com/Huray-hub/eclass-utils/deadlines/pkg"
	"github.com/olekukonko/tablewriter"
)

func main() {
	opts, err := in.GetOptions()
	if err != nil {
		log.Fatal(err.Error())
	}

	credentials, err := in.GetCredentials()
	if err != nil {
		log.Fatal(err.Error())
	}

	assignments, err := dl.Deadlines(opts, credentials)
	if err != nil {
		log.Fatal(err.Error())
	}

	if opts.IgnoreExpired {
		assignments, err = in.FilterExpiredDeadlines(assignments)
		if err != nil {
			log.Fatal(err.Error())
		}
	}

	err = printAssignments(assignments, opts.PlainText)
	if err != nil {
		log.Fatal(err.Error())
	}

	if opts.ExportICS {
		path, err := in.ExportICS(assignments, opts.BaseDomain)
		if err != nil {
			log.Fatal(err.Error())
		}

		fmt.Println("stored in " + path)
	}
}

func printAssignments(a []in.Assignment, plainText bool) error {
	if plainText {
		return printAssignmentsPlain(a)
	}

	return printAssignmentsPretty(a)
}

func printAssignmentsPlain(a []in.Assignment) error {
	for _, v := range a {
		_, err := fmt.Println(v.String())
		if err != nil {
			return err
		}
	}
	return nil
}

// TODO: Fix this ugly code
func printAssignmentsPretty(a []in.Assignment) error {
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

func calcRemainingTime(a in.Assignment) string {
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
