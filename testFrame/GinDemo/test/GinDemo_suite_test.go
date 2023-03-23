package test_test

import (
	"testing"

	//_ "GinDemoTest/controllers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestGinDemo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GinDemo Suite")
}
