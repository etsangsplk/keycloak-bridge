package mobileapi

import (
	"context"
	"testing"
	"time"

	"github.com/cloudtrust/keycloak-bridge/internal/dto"

	"github.com/cloudtrust/common-service/log"
	"github.com/stretchr/testify/assert"
)

func TestSetAccreditations(t *testing.T) {
	var userInfo UserInformationRepresentation
	var ctx = context.TODO()
	var logger = log.NewNopLogger()

	t.Run("Nil accreditations", func(t *testing.T) {
		userInfo.SetAccreditations(ctx, nil, logger)
		assert.Nil(t, userInfo.Accreditations)
	})

	t.Run("Empty accreditations", func(t *testing.T) {
		userInfo.SetAccreditations(ctx, []string{}, logger)
		assert.Nil(t, userInfo.Accreditations)
	})

	t.Run("3 accreditations with one non-unmarshallable", func(t *testing.T) {
		userInfo.SetAccreditations(ctx, []string{"{}", "{}", "{"}, logger)
		assert.NotNil(t, userInfo.Accreditations)
		assert.Len(t, *userInfo.Accreditations, 2)
	})
}

func TestSetChecks(t *testing.T) {
	var userInfo UserInformationRepresentation

	t.Run("Nil checks", func(t *testing.T) {
		userInfo.SetChecks(nil)
		assert.Nil(t, userInfo.Checks)
	})

	t.Run("Empty checks", func(t *testing.T) {
		userInfo.SetChecks([]dto.DBCheck{})
		assert.Nil(t, userInfo.Checks)
	})

	t.Run("With checks", func(t *testing.T) {
		var oneDate = time.Now()
		var checks = []dto.DBCheck{dto.DBCheck{DateTime: nil}, dto.DBCheck{DateTime: &oneDate}, dto.DBCheck{DateTime: &oneDate}}
		userInfo.SetChecks(checks)
		assert.Len(t, *userInfo.Checks, len(checks))
	})
}

func TestSetActions(t *testing.T) {
	var userInfo UserInformationRepresentation

	t.Run("Nil actions", func(t *testing.T) {
		userInfo.SetActions(nil)
		assert.Nil(t, userInfo.Actions)
	})

	t.Run("Empty actions", func(t *testing.T) {
		userInfo.SetActions(map[string]bool{})
		assert.Nil(t, userInfo.Actions)
	})

	t.Run("Wich actions", func(t *testing.T) {
		userInfo.SetActions(map[string]bool{"one": false, "two": true, "three": false})
		assert.Len(t, *userInfo.Actions, 1)
	})
}
