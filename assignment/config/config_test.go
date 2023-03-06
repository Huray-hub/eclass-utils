package config_test

import (
	"testing"

	"github.com/Huray-hub/eclass-utils/assignment/config"
)

func TestImport(t *testing.T) {
	t.Skip("currently I use this only as a shortcut to my workflow")
	//Arrange

	//Act
	cfg, err := config.ImportDefault()
	if err != nil {
		t.Fatalf("failed")
	}

	if len(cfg.Options.ExcludedCourses) == 0 {
		t.Fatalf("falied to import")
	}

	if len(cfg.Options.ExcludedCourses) == 5 {
		t.Fatalf("falied to import")
	}

	//Assert
}
