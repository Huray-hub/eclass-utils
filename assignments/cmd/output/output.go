package output

import (
	"fmt"
	"math"
	"os"
	"time"

	"github.com/Huray-hub/eclass-utils/assignments/assignment"
	"github.com/olekukonko/tablewriter"
)

func PrintAssignments(a []assignment.Assignment, plainText bool) error {
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
		remainingTime(v)
		var isSent string
		if v.IsSent {
			isSent = "Ναι"
		} else {
			isSent = "Όχι"
		}
		table.Append([]string{
			v.Course.Name,
			v.Title,
			v.Deadline.Format("02/01/2006 15:04") + " " + remainingTime(v),
			isSent,
		})
	}
}

func remainingTime(a assignment.Assignment) string {
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
