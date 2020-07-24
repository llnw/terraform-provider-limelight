package limelight

import (
	"os"
)

func getShortname() string {
	return os.Getenv("LLNW_TEST_SHORTNAME")
}
