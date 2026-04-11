package pgxredis_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestPgxredis(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "pgxredis Suite")
}
