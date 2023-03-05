package assignment

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/Huray-hub/eclass-utils/assignments/config"
	auth "github.com/Huray-hub/eclass-utils/authentication"
	"github.com/Huray-hub/eclass-utils/course"
	"github.com/PuerkitoBio/goquery"
)

type Service struct {
	opts     *config.Options
	location *time.Location
	client   *http.Client
}

func NewService(
	ctx context.Context,
	opts *config.Options,
	creds auth.Credentials,
	client *http.Client,
) (*Service, error) {
	location, err := time.LoadLocation("Europe/Athens")
	if err != nil {
		return nil, err
	}

	client, err = auth.Login(ctx, "https://"+opts.BaseDomain, creds, client)
	if err != nil {
		return nil, err
	}

	return &Service{
		opts:     opts,
		location: location,
		client:   client,
	}, nil
}


// FetchAssignments method will retrieve all assignments of your assignments and filter them out 
// based on Service's options.
//
// 1. retrieves enrolled courses, 
//
// 2. concurrently fans-out each course's http request + assignment scrapping through channels, 
//
// 3. fans-in channels' result into one, 
//
// 4. sorts them by deadline
func (svc *Service) FetchAssignments(ctx context.Context) ([]Assignment, error) {
	courses, err := course.GetEnrolled(ctx, svc.opts.Options, svc.client)
	if err != nil {
		return nil, err
	}

	assignmentCh, err := svc.getAssignments(ctx, courses)
	if err != nil {
		return nil, err
	}

	// The usual peek number of non-expired assignments is 10
	assignments := make([]Assignment, 0, 10)
	for asg := range assignmentCh {
		assignments = append(assignments, asg)
	}

	SortByDeadline(assignments)

	return assignments, nil
}

func (svc *Service) getAssignments(ctx context.Context, courses []course.Course,
) (<-chan Assignment, error) {
	// The max number of courses per semester is 7. Till half the semester students expect
	// courses' grades from previous semester's exam, so they are enrolled to maximum 14
	// courses in total
	assignmentChans := make([]<-chan Assignment, 0, 14)

	// fanning out http call + html scrapping per course
	for _, crs := range courses {
		assignmentCh, err := svc.getAssignmentsPerCourse(ctx, crs)
		if err != nil {
			// TODO: log
			continue
		}

		assignmentChans = append(assignmentChans, assignmentCh)
	}

	multiplexedAssignmentCh := fanInAssignments(ctx, assignmentChans)

	return multiplexedAssignmentCh, nil
}

func fanInAssignments(ctx context.Context, assignmentChans []<-chan Assignment) <-chan Assignment {
	var wg sync.WaitGroup

	multiplexedCh := make(chan Assignment)

	multiplex := func(c <-chan Assignment) {
		defer wg.Done()
		for i := range c {
			select {
			case <-ctx.Done():
				return
			case multiplexedCh <- i:
			}

		}
	}

	wg.Add(len(assignmentChans))
	for _, c := range assignmentChans {
		go multiplex(c)
	}

	go func() {
		wg.Wait()
		close(multiplexedCh)
	}()

	return multiplexedCh
}

func (svc *Service) getAssignmentsPerCourse(ctx context.Context, course course.Course,
) (<-chan Assignment, error) {
	assignmentCh := make(chan Assignment)
	go func() {
		defer close(assignmentCh)

		assignmentsURL, err := course.PrepareAssignmentsURL(svc.opts.BaseDomain)
		if err != nil {
			return
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://"+assignmentsURL, nil)
		if err != nil {
			return
		}

		resp, err := svc.client.Do(req)
		if err != nil {
			return
		}
		defer func() {
			err = resp.Body.Close()
			if err != nil {
				// TODO: avoid panic. Propagate error instead
				panic("could not close response body")
			}
		}()

		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			return
		}

		// instead of query first rows and then query again each row to get table cells, 
		// knowing the each row has 4, I immediately get query tds for performance improvement
		tds := make([]*goquery.Selection, 4)
		doc.Find("table#assignment_table tbody tr td").
			Each(func(i int, td *goquery.Selection) {
				tds[i%4] = td

				// on each 4th element process with assignment creation
				if i%4 != 3 {
					return
				}

				assignment, err := newAssignment(tds, &course, svc.location)
				if err != nil {
					return
				}

				if assignment.IsExcluded(svc.opts, course.ID, svc.location) {
					return
				}

				assignmentCh <- assignment
			})
	}()

	return assignmentCh, nil
}
