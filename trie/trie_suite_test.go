package trie_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestTrie(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Trie Suite")
}
