package mail

import (
	"testing"

	"github.com/jurgisjaska/binbogami/internal"
	"github.com/stretchr/testify/assert"
	"gopkg.in/gomail.v2"
)

func TestCreateInvitation(t *testing.T) {
	var dialer *gomail.Dialer
	var config *internal.Config

	result := CreateInvitation(dialer, config)

	assert.Equal(t, dialer, result.d)
	assert.Equal(t, config, result.c)
}
