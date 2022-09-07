package search

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildStackOverflowUrl(t *testing.T) {
	e1 := "https://api.stackexchange.com/a?order=desc&site=stackoverflow&sort=votes"
	assert.Equal(t, e1, buildStackOverflowUrl("a", false), "Invalid stackoverflow URL")
	e2 := "https://api.stackexchange.com/a?filter=%216VvPDzQ%29xXOrL&order=desc&site=stackoverflow&sort=votes"
	assert.Equal(t, e2, buildStackOverflowUrl("a", true), "Invalid stackoverflow URL")
}
