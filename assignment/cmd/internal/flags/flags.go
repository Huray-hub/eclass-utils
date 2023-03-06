package flags

import (
	"flag"
	"strings"

	"github.com/Huray-hub/eclass-utils/assignment/config"
	"github.com/Huray-hub/eclass-utils/auth"
)

func Read(cfg *config.Config) {
	flag.BoolVar(
		&cfg.Options.PlainText,
		"p",
		cfg.Options.PlainText,
		"Print results in plain csv format",
	)
	flag.BoolVar(
		&cfg.Options.IncludeExpired,
		"i",
		cfg.Options.IncludeExpired,
		"Include expired assignments",
	)
	flag.BoolVar(&cfg.Options.ExportICS, "c", cfg.Options.ExportICS, "Export calendar file")
	baseDomain := flag.String(
		"d",
		"",
		"Specify base e-class domain (ex. -d=eclass.uniwa.gr)",
	)
	excludedCourses := flag.String(
		"e",
		"",
		"Exclude courses by ID (ex. -e=ICE262,CS152)",
	)
	excludedAssignments := flag.String(
		"a",
		"",
		`Exclude assignments by pattern. 
Use course ID and a part of the assignment's title to ignore it from results
(ex. -a=ICE262:"τμήματα Tετάρτης,τμήματα Παρασκευής"_CS152:...)`)

	username := flag.String("username", "", "Your e-class username")
	password := flag.String("password", "", "Your e-class password")

	flag.Parse()

	flagsToOptions(*baseDomain, *excludedCourses, *excludedAssignments,&cfg.Options)
	flagsToCredentials(*username, *password, &cfg.Credentials)
}

func flagsToOptions(
	baseDomain string,
	excludedCourses string,
	excludedAssignments string,
	opts *config.Options,
) {
	if baseDomain != "" {
		opts.BaseDomain = baseDomain
	}

	if excludedCourses != "" {
		opts.ExcludedCourses = parseExcludedCourses(excludedCourses)
	}

	if excludedAssignments != "" {
		opts.ExcludedAssignments = parseExcludedAssignments(excludedAssignments)
	}
}

func parseExcludedCourses(raw string) map[string]struct{} {
	excludedCourses := strings.Split(raw, ",")

	res := make(map[string]struct{}, len(excludedCourses))

	for _, v := range excludedCourses {
		res[strings.TrimSpace(v)] = struct{}{}
	}
	return res
}

func parseExcludedAssignments(raw string) map[string][]string {
	kvPairs := strings.Split(raw, "_")
	res := make(map[string][]string, len(kvPairs))

	if len(kvPairs) == 0 {
		return res
	}

	for _, kv := range kvPairs {
		key, valuesCSV, found := strings.Cut(kv, ":")
		if !found {
			continue
		}
		values := strings.Split(valuesCSV, ",")
		res[key] = values
	}
	return res
}

func flagsToCredentials(username string, password string, creds *auth.Credentials) {
	if username != "" {
		creds.Username = username
	}

	if password != "" {
		creds.Password = password
	}
}
