package output

import (
	"fmt"
	"math"
	"os"
	"time"

	"github.com/Huray-hub/eclass-utils/assignment"
	"github.com/olekukonko/tablewriter"
)

func PrintAssignments(assignments []assignment.Assignment, plainText bool) error {
	if plainText {
		return printAssignmentsPlain(assignments)
	}

	return printAssignmentsPretty(assignments)
}

func printAssignmentsPlain(assignments []assignment.Assignment) error {
	for _, a := range assignments {
		_, err := fmt.Println(a.String())
		if err != nil {
			return err
		}
	}
	return nil
}

func printAssignmentsPretty(assignments []assignment.Assignment) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetRowLine(true)
	table.SetHeader([]string{"ΜΑΘΗΜΑ", "ΕΡΓΑΣΙΑ", "ΠΡΟΘΕΣΜΙΑ", "ΥΠΟΒΛΗΘΗΚΕ"})
	table.SetColumnAlignment(
		[]int{
			tablewriter.ALIGN_DEFAULT,
			tablewriter.ALIGN_DEFAULT,
			tablewriter.ALIGN_DEFAULT,
			tablewriter.ALIGN_CENTER,
		},
	)
	appendToTable(assignments, table)
	table.Render()

	return nil
}

func appendToTable(assignments []assignment.Assignment, table *tablewriter.Table) {
	for _, asg := range assignments {
		var isSent string
		if asg.IsSent {
			isSent = "✓"
		} else {
			isSent = "✗"
		}

		var deadlineMsg string
		if asg.Deadline == nil {
			deadlineMsg = assignment.NoDeadline
		} else {
			deadlineMsg = asg.Deadline.Format("02/01/2006 15:04") + " " + remainingTime(asg)
		}

		table.Append([]string{
			asg.Course.Name,
			asg.Title,
			deadlineMsg,
			isSent,
		})
	}
}

func remainingTime(asg assignment.Assignment) string {
	t := time.Until(*asg.Deadline)

	switch {
	case t < 0:
		return "(Έληξε)"
	case t.Hours()/24 >= 1:
		return "(" + fmt.Sprint(math.Floor(t.Hours()/24)) + " μέρες)"
	case t.Minutes()/60 >= 1:
		return "(" + fmt.Sprint(math.Floor(t.Hours())) + " ώρες)"
	default:
		return "(" + fmt.Sprint(math.Floor(t.Minutes())) + " λεπτά)"
	}
}
