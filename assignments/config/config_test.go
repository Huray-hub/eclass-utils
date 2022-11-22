package config_test

import (
	"testing"

	"github.com/Huray-hub/eclass-utils/assignments/config"
)

func TestImport(t *testing.T) {
	t.Skip("currently I use this only as a shorcut to my workflow")
	//Arrange

	//Act
	opts, _, err := config.Import()
	if err != nil {
		t.Fatalf("failed")
	}

	if len(opts.ExcludedCourses) == 0 {
		t.Fatalf("falied to import")
	}

	if len(opts.ExcludedCourses) == 5 {
		t.Fatalf("falied to import")
	}

	//Assert
}
