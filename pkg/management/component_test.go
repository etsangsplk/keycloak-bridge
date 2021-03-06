package management

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"testing"
	"time"

	cs "github.com/cloudtrust/common-service"
	"github.com/cloudtrust/common-service/configuration"
	"github.com/cloudtrust/common-service/database"
	commonhttp "github.com/cloudtrust/common-service/errors"
	errorhandler "github.com/cloudtrust/common-service/errors"
	"github.com/cloudtrust/common-service/log"
	api "github.com/cloudtrust/keycloak-bridge/api/management"
	"github.com/cloudtrust/keycloak-bridge/internal/constants"
	"github.com/cloudtrust/keycloak-bridge/internal/dto"
	"github.com/cloudtrust/keycloak-client"

	"github.com/cloudtrust/keycloak-bridge/pkg/management/mock"
	kc "github.com/cloudtrust/keycloak-client"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestGetActions(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()
	var mockKeycloakClient = mock.NewKeycloakClient(mockCtrl)
	var mockUsersDetailsDBModule = mock.NewUsersDetailsDBModule(mockCtrl)
	var mockEventDBModule = mock.NewEventDBModule(mockCtrl)
	var mockConfigurationDBModule = mock.NewConfigurationDBModule(mockCtrl)
	var mockLogger = log.NewNopLogger()
	var allowedTrustIDGroups = []string{"grp1", "grp2"}

	var managementComponent = NewComponent(mockKeycloakClient, mockUsersDetailsDBModule, mockEventDBModule, mockConfigurationDBModule, allowedTrustIDGroups, mockLogger)

	var accessToken = "TOKEN=="
	var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)

	res, err := managementComponent.GetActions(ctx)

	assert.Nil(t, err)
	assert.Equal(t, len(actions), len(res))
}

func TestGetRealms(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()
	var mockKeycloakClient = mock.NewKeycloakClient(mockCtrl)
	var mockUsersDetailsDBModule = mock.NewUsersDetailsDBModule(mockCtrl)
	var mockEventDBModule = mock.NewEventDBModule(mockCtrl)
	var mockConfigurationDBModule = mock.NewConfigurationDBModule(mockCtrl)
	var mockLogger = log.NewNopLogger()
	var allowedTrustIDGroups = []string{"grp1", "grp2"}

	var managementComponent = NewComponent(mockKeycloakClient, mockUsersDetailsDBModule, mockEventDBModule, mockConfigurationDBModule, allowedTrustIDGroups, mockLogger)

	var accessToken = "TOKEN=="

	// Get realms with succces
	{
		var id = "1245"
		var keycloakVersion = "4.8.3"
		var realm = "master"
		var displayName = "Master"
		var enabled = true

		var kcRealmRep = kc.RealmRepresentation{
			ID:              &id,
			KeycloakVersion: &keycloakVersion,
			Realm:           &realm,
			DisplayName:     &displayName,
			Enabled:         &enabled,
		}

		var kcRealmsRep []kc.RealmRepresentation
		kcRealmsRep = append(kcRealmsRep, kcRealmRep)

		mockKeycloakClient.EXPECT().GetRealms(accessToken).Return(kcRealmsRep, nil).Times(1)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)

		apiRealmsRep, err := managementComponent.GetRealms(ctx)

		var expectedAPIRealmRep = api.RealmRepresentation{
			ID:              &id,
			KeycloakVersion: &keycloakVersion,
			Realm:           &realm,
			DisplayName:     &displayName,
			Enabled:         &enabled,
		}

		var expectedAPIRealmsRep []api.RealmRepresentation
		expectedAPIRealmsRep = append(expectedAPIRealmsRep, expectedAPIRealmRep)

		assert.Nil(t, err)
		assert.Equal(t, expectedAPIRealmsRep, apiRealmsRep)
	}

	//Error
	{
		mockKeycloakClient.EXPECT().GetRealms(accessToken).Return([]kc.RealmRepresentation{}, fmt.Errorf("Unexpected error")).Times(1)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)

		_, err := managementComponent.GetRealms(ctx)

		assert.NotNil(t, err)
	}
}

func TestGetRealm(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()
	var mockKeycloakClient = mock.NewKeycloakClient(mockCtrl)
	var mockUsersDetailsDBModule = mock.NewUsersDetailsDBModule(mockCtrl)
	var mockEventDBModule = mock.NewEventDBModule(mockCtrl)
	var mockConfigurationDBModule = mock.NewConfigurationDBModule(mockCtrl)
	var mockLogger = log.NewNopLogger()
	var allowedTrustIDGroups = []string{"grp1", "grp2"}

	var managementComponent = NewComponent(mockKeycloakClient, mockUsersDetailsDBModule, mockEventDBModule, mockConfigurationDBModule, allowedTrustIDGroups, mockLogger)

	var accessToken = "TOKEN=="
	var realmName = "master"
	var username = "username"

	// Get realm with succces
	{
		var id = "1245"
		var keycloakVersion = "4.8.3"
		var realm = "master"
		var displayName = "Master"
		var enabled = true

		var kcRealmRep = kc.RealmRepresentation{
			ID:              &id,
			KeycloakVersion: &keycloakVersion,
			Realm:           &realm,
			DisplayName:     &displayName,
			Enabled:         &enabled,
		}

		mockKeycloakClient.EXPECT().GetRealm(accessToken, realmName).Return(kcRealmRep, nil).Times(1)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		ctx = context.WithValue(ctx, cs.CtContextRealm, realmName)
		ctx = context.WithValue(ctx, cs.CtContextUsername, username)

		mockEventDBModule.EXPECT().Store(ctx, gomock.Any()).Return(nil).AnyTimes()

		apiRealmRep, err := managementComponent.GetRealm(ctx, "master")

		var expectedAPIRealmRep = api.RealmRepresentation{
			ID:              &id,
			KeycloakVersion: &keycloakVersion,
			Realm:           &realm,
			DisplayName:     &displayName,
			Enabled:         &enabled,
		}

		assert.Nil(t, err)
		assert.Equal(t, expectedAPIRealmRep, apiRealmRep)
	}

	//Error
	{
		mockKeycloakClient.EXPECT().GetRealm(accessToken, realmName).Return(kc.RealmRepresentation{}, fmt.Errorf("Unexpected error")).Times(1)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		//mockEventDBModule.EXPECT().Store(ctx, gomock.Any()).Return(nil).Times(1)

		_, err := managementComponent.GetRealm(ctx, "master")

		assert.NotNil(t, err)
	}
}

func TestGetClient(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()
	var mockKeycloakClient = mock.NewKeycloakClient(mockCtrl)
	var mockUsersDetailsDBModule = mock.NewUsersDetailsDBModule(mockCtrl)
	var mockEventDBModule = mock.NewEventDBModule(mockCtrl)
	var mockConfigurationDBModule = mock.NewConfigurationDBModule(mockCtrl)
	var mockLogger = log.NewNopLogger()
	var allowedTrustIDGroups = []string{"grp1", "grp2"}

	var managementComponent = NewComponent(mockKeycloakClient, mockUsersDetailsDBModule, mockEventDBModule, mockConfigurationDBModule, allowedTrustIDGroups, mockLogger)

	var accessToken = "TOKEN=="
	var realmName = "master"

	// Get client with succces
	{
		var id = "1245-1245-4578"
		var name = "clientName"
		var baseURL = "http://toto.com"
		var clientID = "client-id"
		var protocol = "saml"
		var enabled = true
		var username = "username"

		var kcClientRep = kc.ClientRepresentation{
			ID:       &id,
			Name:     &name,
			BaseURL:  &baseURL,
			ClientID: &clientID,
			Protocol: &protocol,
			Enabled:  &enabled,
		}

		mockKeycloakClient.EXPECT().GetClient(accessToken, realmName, id).Return(kcClientRep, nil).Times(1)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		ctx = context.WithValue(ctx, cs.CtContextRealm, realmName)
		ctx = context.WithValue(ctx, cs.CtContextUsername, username)

		mockEventDBModule.EXPECT().Store(ctx, gomock.Any()).Return(nil).AnyTimes()

		apiClientRep, err := managementComponent.GetClient(ctx, "master", id)

		var expectedAPIClientRep = api.ClientRepresentation{
			ID:       &id,
			Name:     &name,
			BaseURL:  &baseURL,
			ClientID: &clientID,
			Protocol: &protocol,
			Enabled:  &enabled,
		}

		assert.Nil(t, err)
		assert.Equal(t, expectedAPIClientRep, apiClientRep)
	}

	//Error
	{
		var id = "1234-79894-7594"
		mockKeycloakClient.EXPECT().GetClient(accessToken, realmName, id).Return(kc.ClientRepresentation{}, fmt.Errorf("Unexpected error")).Times(1)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)

		_, err := managementComponent.GetClient(ctx, "master", id)

		assert.NotNil(t, err)
	}
}

func TestGetClients(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()
	var mockKeycloakClient = mock.NewKeycloakClient(mockCtrl)
	var mockUsersDetailsDBModule = mock.NewUsersDetailsDBModule(mockCtrl)
	var mockEventDBModule = mock.NewEventDBModule(mockCtrl)
	var mockConfigurationDBModule = mock.NewConfigurationDBModule(mockCtrl)
	var mockLogger = log.NewNopLogger()
	var allowedTrustIDGroups = []string{"grp1", "grp2"}

	var managementComponent = NewComponent(mockKeycloakClient, mockUsersDetailsDBModule, mockEventDBModule, mockConfigurationDBModule, allowedTrustIDGroups, mockLogger)

	var accessToken = "TOKEN=="
	var realmName = "master"

	// Get clients with succces
	{
		var id = "1234-7894-58"
		var name = "clientName"
		var baseURL = "http://toto.com"
		var clientID = "client-id"
		var protocol = "saml"
		var enabled = true

		var kcClientRep = kc.ClientRepresentation{
			ID:       &id,
			Name:     &name,
			BaseURL:  &baseURL,
			ClientID: &clientID,
			Protocol: &protocol,
			Enabled:  &enabled,
		}

		var kcClientsRep []kc.ClientRepresentation
		kcClientsRep = append(kcClientsRep, kcClientRep)

		mockKeycloakClient.EXPECT().GetClients(accessToken, realmName).Return(kcClientsRep, nil).Times(1)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)

		apiClientsRep, err := managementComponent.GetClients(ctx, "master")

		var expectedAPIClientRep = api.ClientRepresentation{
			ID:       &id,
			Name:     &name,
			BaseURL:  &baseURL,
			ClientID: &clientID,
			Protocol: &protocol,
			Enabled:  &enabled,
		}

		var expectedAPIClientsRep []api.ClientRepresentation
		expectedAPIClientsRep = append(expectedAPIClientsRep, expectedAPIClientRep)

		assert.Nil(t, err)
		assert.Equal(t, expectedAPIClientsRep, apiClientsRep)
	}

	//Error
	{
		mockKeycloakClient.EXPECT().GetClients(accessToken, realmName).Return([]kc.ClientRepresentation{}, fmt.Errorf("Unexpected error")).Times(1)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)

		_, err := managementComponent.GetClients(ctx, "master")

		assert.NotNil(t, err)
	}
}

func TestGetRequiredActions(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()
	var mockKeycloakClient = mock.NewKeycloakClient(mockCtrl)
	var mockUsersDetailsDBModule = mock.NewUsersDetailsDBModule(mockCtrl)
	var mockEventDBModule = mock.NewEventDBModule(mockCtrl)
	var mockConfigurationDBModule = mock.NewConfigurationDBModule(mockCtrl)
	var mockLogger = log.NewNopLogger()
	var allowedTrustIDGroups = []string{"grp1", "grp2"}

	var managementComponent = NewComponent(mockKeycloakClient, mockUsersDetailsDBModule, mockEventDBModule, mockConfigurationDBModule, allowedTrustIDGroups, mockLogger)

	var accessToken = "TOKEN=="
	var realmName = "master"

	// Get required actions with succces
	{
		var alias = "ALIAS"
		var name = "name"
		var boolTrue = true
		var boolFalse = false

		var kcRa = kc.RequiredActionProviderRepresentation{
			Alias:         &alias,
			Name:          &name,
			Enabled:       &boolTrue,
			DefaultAction: &boolTrue,
		}

		var kcDisabledRa = kc.RequiredActionProviderRepresentation{
			Alias:         &alias,
			Name:          &name,
			Enabled:       &boolFalse,
			DefaultAction: &boolFalse,
		}

		var kcRasRep []kc.RequiredActionProviderRepresentation
		kcRasRep = append(kcRasRep, kcRa, kcDisabledRa)

		mockKeycloakClient.EXPECT().GetRequiredActions(accessToken, realmName).Return(kcRasRep, nil).Times(1)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)

		apiRasRep, err := managementComponent.GetRequiredActions(ctx, "master")

		var expectedAPIRaRep = api.RequiredActionRepresentation{
			Alias:         &alias,
			Name:          &name,
			DefaultAction: &boolTrue,
		}

		var expectedAPIRasRep []api.RequiredActionRepresentation
		expectedAPIRasRep = append(expectedAPIRasRep, expectedAPIRaRep)

		assert.Nil(t, err)
		assert.Equal(t, expectedAPIRasRep, apiRasRep)
	}

	//Error
	{
		mockKeycloakClient.EXPECT().GetRequiredActions(accessToken, realmName).Return([]kc.RequiredActionProviderRepresentation{}, fmt.Errorf("Unexpected error")).Times(1)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)

		_, err := managementComponent.GetRequiredActions(ctx, "master")

		assert.NotNil(t, err)
	}
}

func TestCreateUser(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()
	var mockKeycloakClient = mock.NewKeycloakClient(mockCtrl)
	var mockUsersDetailsDBModule = mock.NewUsersDetailsDBModule(mockCtrl)
	var mockEventDBModule = mock.NewEventDBModule(mockCtrl)
	var mockConfigurationDBModule = mock.NewConfigurationDBModule(mockCtrl)
	var mockLogger = mock.NewLogger(mockCtrl)
	var allowedTrustIDGroups = []string{"grp1", "grp2"}

	var managementComponent = NewComponent(mockKeycloakClient, mockUsersDetailsDBModule, mockEventDBModule, mockConfigurationDBModule, allowedTrustIDGroups, mockLogger)

	var accessToken = "TOKEN=="
	var username = "test"
	var realmName = "master"
	var targetRealmName = "DEP"
	var userID = "41dbf4a8-32a9-4000-8c17-edc854c31231"
	var locationURL = "http://toto.com/realms/" + userID

	t.Run("Create with minimum properties", func(t *testing.T) {
		var kcUserRep = kc.UserRepresentation{
			Username: &username,
		}

		mockKeycloakClient.EXPECT().CreateUser(accessToken, realmName, targetRealmName, kcUserRep).Return(locationURL, nil).Times(1)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		ctx = context.WithValue(ctx, cs.CtContextRealm, realmName)
		ctx = context.WithValue(ctx, cs.CtContextUsername, username)

		mockEventDBModule.EXPECT().ReportEvent(ctx, "API_ACCOUNT_CREATION", "back-office", gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)

		var userRep = api.UserRepresentation{
			Username: &username,
		}

		location, err := managementComponent.CreateUser(ctx, targetRealmName, userRep)

		assert.Nil(t, err)
		assert.Equal(t, locationURL, location)
	})

	t.Run("Create with minimum properties and having error when storing the event", func(t *testing.T) {
		var kcUserRep = kc.UserRepresentation{
			Username: &username,
		}

		mockKeycloakClient.EXPECT().CreateUser(accessToken, realmName, realmName, kcUserRep).Return(locationURL, nil).Times(1)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		ctx = context.WithValue(ctx, cs.CtContextRealm, realmName)
		ctx = context.WithValue(ctx, cs.CtContextUsername, username)

		mockEventDBModule.EXPECT().ReportEvent(ctx, "API_ACCOUNT_CREATION", "back-office", database.CtEventRealmName, realmName, database.CtEventUserID, userID, database.CtEventUsername, username).Return(errors.New("error")).Times(1)
		m := map[string]interface{}{"event_name": "API_ACCOUNT_CREATION", database.CtEventRealmName: realmName, database.CtEventUserID: userID, database.CtEventUsername: username}
		eventJSON, _ := json.Marshal(m)
		mockLogger.EXPECT().Error(ctx, "err", "error", "event", string(eventJSON))

		var userRep = api.UserRepresentation{
			Username: &username,
		}

		location, err := managementComponent.CreateUser(ctx, realmName, userRep)

		assert.Nil(t, err)
		assert.Equal(t, locationURL, location)
	})

	t.Run("Create with all properties allowed by Bridge API", func(t *testing.T) {
		var email = "toto@elca.ch"
		var enabled = true
		var emailVerified = true
		var firstName = "Titi"
		var lastName = "Tutu"
		var phoneNumber = "+41789456"
		var phoneNumberVerified = true
		var label = "Label"
		var gender = "M"
		var birthDate = "01/01/1988"
		var locale = "de"

		var groups = []string{"145-784-545251"}
		var trustIDGroups = []string{"l1_support_agent"}
		var roles = []string{"445-4545-751515"}

		var birthLocation = "Rolle"
		var idDocumentType = "Card ID"
		var idDocumentNumber = "1234-4567-VD-3"
		var idDocumentExpiration = "23.12.2019"

		mockKeycloakClient.EXPECT().CreateUser(accessToken, realmName, targetRealmName, gomock.Any()).DoAndReturn(
			func(accessToken, realmName, targetRealmName string, kcUserRep kc.UserRepresentation) (string, error) {
				assert.Equal(t, username, *kcUserRep.Username)
				assert.Equal(t, email, *kcUserRep.Email)
				assert.Equal(t, enabled, *kcUserRep.Enabled)
				assert.Equal(t, emailVerified, *kcUserRep.EmailVerified)
				assert.Equal(t, firstName, *kcUserRep.FirstName)
				assert.Equal(t, lastName, *kcUserRep.LastName)
				assert.Equal(t, phoneNumber, *kcUserRep.GetAttributeString(constants.AttrbPhoneNumber))
				verified, _ := kcUserRep.GetAttributeBool(constants.AttrbPhoneNumberVerified)
				assert.Equal(t, phoneNumberVerified, *verified)
				assert.Equal(t, label, *kcUserRep.GetAttributeString(constants.AttrbLabel))
				assert.Equal(t, gender, *kcUserRep.GetAttributeString(constants.AttrbGender))
				assert.Equal(t, birthDate, *kcUserRep.GetAttributeString(constants.AttrbBirthDate))
				assert.Equal(t, locale, *kcUserRep.GetAttributeString(constants.AttrbLocale))
				return locationURL, nil
			}).Times(1)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		ctx = context.WithValue(ctx, cs.CtContextRealm, realmName)
		ctx = context.WithValue(ctx, cs.CtContextUsername, username)

		mockUsersDetailsDBModule.EXPECT().StoreOrUpdateUserDetails(ctx, targetRealmName, gomock.Any()).DoAndReturn(
			func(ctx context.Context, targetRealmName string, user dto.DBUser) {
				assert.Equal(t, userID, *user.UserID)
				assert.Equal(t, birthLocation, *user.BirthLocation)
				assert.Equal(t, idDocumentType, *user.IDDocumentType)
				assert.Equal(t, idDocumentNumber, *user.IDDocumentNumber)
				assert.Equal(t, idDocumentExpiration, *user.IDDocumentExpiration)
			}).Return(nil).Times(1)

		mockEventDBModule.EXPECT().Store(ctx, gomock.Any()).Return(nil).AnyTimes()

		var userRep = api.UserRepresentation{
			ID:                   &userID,
			Username:             &username,
			Email:                &email,
			Enabled:              &enabled,
			EmailVerified:        &emailVerified,
			FirstName:            &firstName,
			LastName:             &lastName,
			PhoneNumber:          &phoneNumber,
			PhoneNumberVerified:  &phoneNumberVerified,
			Label:                &label,
			Gender:               &gender,
			BirthDate:            &birthDate,
			Locale:               &locale,
			Groups:               &groups,
			TrustIDGroups:        &trustIDGroups,
			Roles:                &roles,
			BirthLocation:        &birthLocation,
			IDDocumentType:       &idDocumentType,
			IDDocumentNumber:     &idDocumentNumber,
			IDDocumentExpiration: &idDocumentExpiration,
		}
		mockEventDBModule.EXPECT().ReportEvent(ctx, "API_ACCOUNT_CREATION", "back-office", gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)
		location, err := managementComponent.CreateUser(ctx, targetRealmName, userRep)

		assert.Nil(t, err)
		assert.Equal(t, locationURL, location)
	})

	t.Run("Error from KC client", func(t *testing.T) {
		var kcUserRep = kc.UserRepresentation{}

		mockKeycloakClient.EXPECT().CreateUser(accessToken, realmName, targetRealmName, kcUserRep).Return("", fmt.Errorf("Invalid input")).Times(1)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		ctx = context.WithValue(ctx, cs.CtContextRealm, realmName)

		var userRep = api.UserRepresentation{}
		mockLogger.EXPECT().Warn(ctx, "err", "Invalid input")

		location, err := managementComponent.CreateUser(ctx, targetRealmName, userRep)

		assert.NotNil(t, err)
		assert.Equal(t, "", location)
	})

	t.Run("Error from DB users", func(t *testing.T) {
		mockKeycloakClient.EXPECT().CreateUser(accessToken, realmName, targetRealmName, gomock.Any()).DoAndReturn(
			func(accessToken, realmName, targetRealmName string, kcUserRep kc.UserRepresentation) (string, error) {
				return locationURL, nil
			}).Times(1)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		ctx = context.WithValue(ctx, cs.CtContextRealm, realmName)

		mockUsersDetailsDBModule.EXPECT().StoreOrUpdateUserDetails(ctx, targetRealmName, gomock.Any()).Return(fmt.Errorf("SQL error")).Times(1)

		var birthLocation = "Rolle"
		var userRep = api.UserRepresentation{
			ID:            &userID,
			Username:      &username,
			BirthLocation: &birthLocation,
		}
		mockLogger.EXPECT().Warn(ctx, "msg", "Can't store user details in database", "err", "SQL error")

		location, err := managementComponent.CreateUser(ctx, targetRealmName, userRep)

		assert.NotNil(t, err)
		assert.Equal(t, "", location)
	})
}

func TestDeleteUser(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()
	var mockKeycloakClient = mock.NewKeycloakClient(mockCtrl)
	var mockUsersDetailsDBModule = mock.NewUsersDetailsDBModule(mockCtrl)
	var mockEventDBModule = mock.NewEventDBModule(mockCtrl)
	var mockConfigurationDBModule = mock.NewConfigurationDBModule(mockCtrl)
	var mockLogger = mock.NewLogger(mockCtrl)
	var allowedTrustIDGroups = []string{"grp1", "grp2"}

	var managementComponent = NewComponent(mockKeycloakClient, mockUsersDetailsDBModule, mockEventDBModule, mockConfigurationDBModule, allowedTrustIDGroups, mockLogger)

	var accessToken = "TOKEN=="
	var userID = "41dbf4a8-32a9-4000-8c17-edc854c31231"
	var realmName = "master"
	var username = "username"

	t.Run("Delete user with success", func(t *testing.T) {
		mockKeycloakClient.EXPECT().DeleteUser(accessToken, realmName, userID).Return(nil).Times(1)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		ctx = context.WithValue(ctx, cs.CtContextRealm, realmName)
		ctx = context.WithValue(ctx, cs.CtContextUsername, username)

		mockUsersDetailsDBModule.EXPECT().DeleteUserDetails(ctx, realmName, userID).Return(nil).Times(1)

		mockEventDBModule.EXPECT().ReportEvent(ctx, "API_ACCOUNT_DELETION", "back-office", gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)

		err := managementComponent.DeleteUser(ctx, "master", userID)

		assert.Nil(t, err)
	})

	t.Run("Delete user with success but the having an error when storing the event in the DB", func(t *testing.T) {
		mockKeycloakClient.EXPECT().DeleteUser(accessToken, realmName, userID).Return(nil).Times(1)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		ctx = context.WithValue(ctx, cs.CtContextRealm, realmName)
		ctx = context.WithValue(ctx, cs.CtContextUsername, username)

		mockUsersDetailsDBModule.EXPECT().DeleteUserDetails(ctx, realmName, userID).Return(nil).Times(1)

		mockEventDBModule.EXPECT().ReportEvent(ctx, "API_ACCOUNT_DELETION", "back-office", database.CtEventRealmName, realmName, database.CtEventUserID, userID).Return(errors.New("error")).Times(1)
		m := map[string]interface{}{"event_name": "API_ACCOUNT_DELETION", database.CtEventRealmName: realmName, database.CtEventUserID: userID}
		eventJSON, _ := json.Marshal(m)
		mockLogger.EXPECT().Error(ctx, "err", "error", "event", string(eventJSON))
		err := managementComponent.DeleteUser(ctx, "master", userID)

		assert.Nil(t, err)
	})

	t.Run("Error from KC client", func(t *testing.T) {
		mockKeycloakClient.EXPECT().DeleteUser(accessToken, realmName, userID).Return(fmt.Errorf("Invalid input")).Times(1)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		mockLogger.EXPECT().Warn(ctx, "err", "Invalid input")

		err := managementComponent.DeleteUser(ctx, "master", userID)

		assert.NotNil(t, err)
	})

	t.Run("Error from DB users", func(t *testing.T) {
		mockKeycloakClient.EXPECT().DeleteUser(accessToken, realmName, userID).Return(nil).Times(1)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		mockUsersDetailsDBModule.EXPECT().DeleteUserDetails(ctx, realmName, userID).Return(fmt.Errorf("SQL Error")).Times(1)

		mockLogger.EXPECT().Warn(ctx, "err", "SQL Error")

		err := managementComponent.DeleteUser(ctx, "master", userID)

		assert.NotNil(t, err)
	})
}

func TestGetUser(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()
	var mockKeycloakClient = mock.NewKeycloakClient(mockCtrl)
	var mockUsersDetailsDBModule = mock.NewUsersDetailsDBModule(mockCtrl)
	var mockEventDBModule = mock.NewEventDBModule(mockCtrl)
	var mockConfigurationDBModule = mock.NewConfigurationDBModule(mockCtrl)
	var mockLogger = mock.NewLogger(mockCtrl)
	var allowedTrustIDGroups = []string{"grp1", "grp2"}

	var managementComponent = NewComponent(mockKeycloakClient, mockUsersDetailsDBModule, mockEventDBModule, mockConfigurationDBModule, allowedTrustIDGroups, mockLogger)

	var accessToken = "TOKEN=="
	var realmName = "master"
	var id = "41dbf4a8-32a9-4000-8c17-edc854c31231"
	var username = "username"

	t.Run("Get user with succces", func(t *testing.T) {
		var email = "toto@elca.ch"
		var enabled = true
		var emailVerified = true
		var firstName = "Titi"
		var lastName = "Tutu"
		var phoneNumber = "+41789456"
		var phoneNumberVerified = true
		var label = "Label"
		var gender = "M"
		var birthDate = "01/01/1988"
		var createdTimestamp = time.Now().UTC().Unix()
		var locale = "it"
		var trustIDGroups = []string{"grp1", "grp2"}
		var birthLocation = "Rolle"
		var idDocumentType = "Card ID"
		var idDocumentNumber = "1234-4567-VD-3"
		var idDocumentExpiration = "23.12.2019"

		var attributes = make(kc.Attributes)
		attributes.SetString(constants.AttrbPhoneNumber, phoneNumber)
		attributes.SetString(constants.AttrbLabel, label)
		attributes.SetString(constants.AttrbGender, gender)
		attributes.SetString(constants.AttrbBirthDate, birthDate)
		attributes.SetBool(constants.AttrbPhoneNumberVerified, phoneNumberVerified)
		attributes.SetString(constants.AttrbLocale, locale)
		attributes.Set(constants.AttrbTrustIDGroups, trustIDGroups)

		var kcUserRep = kc.UserRepresentation{
			ID:               &id,
			Username:         &username,
			Email:            &email,
			Enabled:          &enabled,
			EmailVerified:    &emailVerified,
			FirstName:        &firstName,
			LastName:         &lastName,
			Attributes:       &attributes,
			CreatedTimestamp: &createdTimestamp,
		}

		mockKeycloakClient.EXPECT().GetUser(accessToken, realmName, id).Return(kcUserRep, nil).Times(1)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		ctx = context.WithValue(ctx, cs.CtContextRealm, realmName)
		ctx = context.WithValue(ctx, cs.CtContextUsername, username)

		mockUsersDetailsDBModule.EXPECT().GetUserDetails(ctx, realmName, id).Return(dto.DBUser{
			UserID:               &id,
			BirthLocation:        &birthLocation,
			IDDocumentExpiration: &idDocumentExpiration,
			IDDocumentNumber:     &idDocumentNumber,
			IDDocumentType:       &idDocumentType,
		}, nil).Times(1)

		mockEventDBModule.EXPECT().ReportEvent(ctx, "GET_DETAILS", "back-office", gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)

		apiUserRep, err := managementComponent.GetUser(ctx, "master", id)

		assert.Nil(t, err)
		assert.Equal(t, username, *apiUserRep.Username)
		assert.Equal(t, email, *apiUserRep.Email)
		assert.Equal(t, enabled, *apiUserRep.Enabled)
		assert.Equal(t, emailVerified, *apiUserRep.EmailVerified)
		assert.Equal(t, firstName, *apiUserRep.FirstName)
		assert.Equal(t, lastName, *apiUserRep.LastName)
		assert.Equal(t, phoneNumber, *apiUserRep.PhoneNumber)
		assert.Equal(t, phoneNumberVerified, *apiUserRep.PhoneNumberVerified)
		assert.Equal(t, label, *apiUserRep.Label)
		assert.Equal(t, gender, *apiUserRep.Gender)
		assert.Equal(t, birthDate, *apiUserRep.BirthDate)
		assert.Equal(t, createdTimestamp, *apiUserRep.CreatedTimestamp)
		assert.Equal(t, locale, *apiUserRep.Locale)
		assert.Equal(t, trustIDGroups, *apiUserRep.TrustIDGroups)
		assert.Equal(t, birthLocation, *apiUserRep.BirthLocation)
		assert.Equal(t, idDocumentExpiration, *apiUserRep.IDDocumentExpiration)
		assert.Equal(t, idDocumentNumber, *apiUserRep.IDDocumentNumber)
		assert.Equal(t, idDocumentType, *apiUserRep.IDDocumentType)
	})

	t.Run("Get user with succces with empty user info", func(t *testing.T) {
		var email = "toto@elca.ch"
		var enabled = true
		var emailVerified = true
		var firstName = "Titi"
		var lastName = "Tutu"
		var phoneNumber = "+41789456"
		var phoneNumberVerified = true
		var label = "Label"
		var gender = "M"
		var birthDate = "01/01/1988"
		var createdTimestamp = time.Now().UTC().Unix()
		var locale = "it"
		var trustIDGroups = []string{"grp1", "grp2"}

		var attributes = make(kc.Attributes)
		attributes.SetString(constants.AttrbPhoneNumber, phoneNumber)
		attributes.SetString(constants.AttrbLabel, label)
		attributes.SetString(constants.AttrbGender, gender)
		attributes.SetString(constants.AttrbBirthDate, birthDate)
		attributes.SetBool(constants.AttrbPhoneNumberVerified, phoneNumberVerified)
		attributes.SetString(constants.AttrbLocale, locale)
		attributes.Set(constants.AttrbTrustIDGroups, trustIDGroups)

		var kcUserRep = kc.UserRepresentation{
			ID:               &id,
			Username:         &username,
			Email:            &email,
			Enabled:          &enabled,
			EmailVerified:    &emailVerified,
			FirstName:        &firstName,
			LastName:         &lastName,
			Attributes:       &attributes,
			CreatedTimestamp: &createdTimestamp,
		}

		mockKeycloakClient.EXPECT().GetUser(accessToken, realmName, id).Return(kcUserRep, nil).Times(1)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		ctx = context.WithValue(ctx, cs.CtContextRealm, realmName)
		ctx = context.WithValue(ctx, cs.CtContextUsername, username)

		mockUsersDetailsDBModule.EXPECT().GetUserDetails(ctx, realmName, id).Return(dto.DBUser{
			UserID: &id,
		}, nil).Times(1)

		mockEventDBModule.EXPECT().ReportEvent(ctx, "GET_DETAILS", "back-office", gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)

		apiUserRep, err := managementComponent.GetUser(ctx, "master", id)

		assert.Nil(t, err)
		assert.Equal(t, username, *apiUserRep.Username)
		assert.Equal(t, email, *apiUserRep.Email)
		assert.Equal(t, enabled, *apiUserRep.Enabled)
		assert.Equal(t, emailVerified, *apiUserRep.EmailVerified)
		assert.Equal(t, firstName, *apiUserRep.FirstName)
		assert.Equal(t, lastName, *apiUserRep.LastName)
		assert.Equal(t, phoneNumber, *apiUserRep.PhoneNumber)
		assert.Equal(t, phoneNumberVerified, *apiUserRep.PhoneNumberVerified)
		assert.Equal(t, label, *apiUserRep.Label)
		assert.Equal(t, gender, *apiUserRep.Gender)
		assert.Equal(t, birthDate, *apiUserRep.BirthDate)
		assert.Equal(t, createdTimestamp, *apiUserRep.CreatedTimestamp)
		assert.Equal(t, locale, *apiUserRep.Locale)
		assert.Equal(t, trustIDGroups, *apiUserRep.TrustIDGroups)
		assert.Nil(t, apiUserRep.BirthLocation)
		assert.Nil(t, apiUserRep.IDDocumentExpiration)
		assert.Nil(t, apiUserRep.IDDocumentNumber)
		assert.Nil(t, apiUserRep.IDDocumentType)
	})

	t.Run("Get user with succces but with error when storing the event in the DB", func(t *testing.T) {
		var birthLocation = "Rolle"
		var idDocumentType = "Card ID"
		var idDocumentNumber = "1234-4567-VD-3"
		var idDocumentExpiration = "23.12.2019"
		var kcUserRep = kc.UserRepresentation{
			ID:       &id,
			Username: &username,
		}
		mockKeycloakClient.EXPECT().GetUser(accessToken, realmName, id).Return(kcUserRep, nil).Times(1)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		ctx = context.WithValue(ctx, cs.CtContextRealm, realmName)
		ctx = context.WithValue(ctx, cs.CtContextUsername, username)

		mockUsersDetailsDBModule.EXPECT().GetUserDetails(ctx, realmName, id).Return(dto.DBUser{
			UserID:               &id,
			BirthLocation:        &birthLocation,
			IDDocumentExpiration: &idDocumentExpiration,
			IDDocumentNumber:     &idDocumentNumber,
			IDDocumentType:       &idDocumentType,
		}, nil).Times(1)

		mockEventDBModule.EXPECT().ReportEvent(ctx, "GET_DETAILS", "back-office", database.CtEventRealmName, realmName, database.CtEventUserID, id, database.CtEventUsername, username).Return(errors.New("error")).Times(1)
		m := map[string]interface{}{"event_name": "GET_DETAILS", database.CtEventRealmName: realmName, database.CtEventUserID: id, database.CtEventUsername: username}
		eventJSON, _ := json.Marshal(m)
		mockLogger.EXPECT().Error(ctx, "err", "error", "event", string(eventJSON))

		apiUserRep, err := managementComponent.GetUser(ctx, "master", id)
		assert.Nil(t, err)
		assert.Equal(t, username, *apiUserRep.Username)
	})

	t.Run("Error with Users DB", func(t *testing.T) {
		var kcUserRep = kc.UserRepresentation{
			ID:       &id,
			Username: &username,
		}
		mockKeycloakClient.EXPECT().GetUser(accessToken, realmName, id).Return(kcUserRep, nil).Times(1)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)

		mockUsersDetailsDBModule.EXPECT().GetUserDetails(ctx, realmName, id).Return(dto.DBUser{}, fmt.Errorf("SQL Error")).Times(1)
		mockLogger.EXPECT().Warn(ctx, "err", "SQL Error")

		_, err := managementComponent.GetUser(ctx, "master", id)

		assert.NotNil(t, err)
	})

	t.Run("Error with KC", func(t *testing.T) {
		var id = "1234-79894-7594"
		mockKeycloakClient.EXPECT().GetUser(accessToken, realmName, id).Return(kc.UserRepresentation{}, fmt.Errorf("Unexpected error")).Times(1)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		mockLogger.EXPECT().Warn(ctx, "err", "Unexpected error")

		_, err := managementComponent.GetUser(ctx, "master", id)

		assert.NotNil(t, err)
	})
}

func TestUpdateUser(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()
	var mockKeycloakClient = mock.NewKeycloakClient(mockCtrl)
	var mockUsersDetailsDBModule = mock.NewUsersDetailsDBModule(mockCtrl)
	var mockEventDBModule = mock.NewEventDBModule(mockCtrl)
	var mockConfigurationDBModule = mock.NewConfigurationDBModule(mockCtrl)
	var mockLogger = mock.NewLogger(mockCtrl)
	var allowedTrustIDGroups = []string{"grp1", "grp2"}

	var managementComponent = NewComponent(mockKeycloakClient, mockUsersDetailsDBModule, mockEventDBModule, mockConfigurationDBModule, allowedTrustIDGroups, mockLogger)

	var accessToken = "TOKEN=="
	var realmName = "master"
	var id = "41dbf4a8-32a9-4000-8c17-edc854c31231"
	var username = "username"
	var enabled = true

	var email = "toto@elca.ch"
	var emailVerified = true
	var firstName = "Titi"
	var lastName = "Tutu"
	var phoneNumber = "+41789456"
	var phoneNumberVerified = true
	var label = "Label"
	var gender = "M"
	var birthDate = "01/01/1988"
	var locale = "de"
	var birthLocation = "Rolle"
	var idDocumentType = "Card ID"
	var idDocumentNumber = "1234-4567-VD-3"
	var idDocumentExpiration = "23.12.2019"
	var createdTimestamp = time.Now().UTC().Unix()

	var attributes = make(kc.Attributes)
	attributes.SetString(constants.AttrbPhoneNumber, phoneNumber)
	attributes.SetString(constants.AttrbLabel, label)
	attributes.SetString(constants.AttrbGender, gender)
	attributes.SetString(constants.AttrbBirthDate, birthDate)
	attributes.SetBool(constants.AttrbPhoneNumberVerified, phoneNumberVerified)
	attributes.SetString(constants.AttrbLocale, locale)

	var kcUserRep = kc.UserRepresentation{
		ID:               &id,
		Username:         &username,
		Email:            &email,
		Enabled:          &enabled,
		EmailVerified:    &emailVerified,
		FirstName:        &firstName,
		LastName:         &lastName,
		Attributes:       &attributes,
		CreatedTimestamp: &createdTimestamp,
	}

	var dbUserRep = dto.DBUser{
		UserID:               &id,
		BirthLocation:        &birthLocation,
		IDDocumentType:       &idDocumentType,
		IDDocumentNumber:     &idDocumentNumber,
		IDDocumentExpiration: &idDocumentExpiration,
	}

	var userRep = api.UserRepresentation{
		Username:            &username,
		Email:               &email,
		EmailVerified:       &emailVerified,
		FirstName:           &firstName,
		LastName:            &lastName,
		PhoneNumber:         &phoneNumber,
		PhoneNumberVerified: &phoneNumberVerified,
		Label:               &label,
		Gender:              &gender,
		BirthDate:           &birthDate,
		Locale:              &locale,
	}

	var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
	ctx = context.WithValue(ctx, cs.CtContextRealm, realmName)
	ctx = context.WithValue(ctx, cs.CtContextUsername, username)

	t.Run("Update user with succces (without user info update)", func(t *testing.T) {
		mockKeycloakClient.EXPECT().GetUser(accessToken, realmName, id).Return(kcUserRep, nil).Times(1)
		mockUsersDetailsDBModule.EXPECT().GetUserDetails(ctx, realmName, id).Return(dbUserRep, nil).Times(1)

		mockKeycloakClient.EXPECT().UpdateUser(accessToken, realmName, id, gomock.Any()).DoAndReturn(
			func(accessToken, realmName, id string, kcUserRep kc.UserRepresentation) error {
				assert.Equal(t, username, *kcUserRep.Username)
				assert.Equal(t, email, *kcUserRep.Email)
				assert.Equal(t, emailVerified, *kcUserRep.EmailVerified)
				assert.Equal(t, firstName, *kcUserRep.FirstName)
				assert.Equal(t, lastName, *kcUserRep.LastName)
				assert.Equal(t, phoneNumber, *kcUserRep.GetAttributeString(constants.AttrbPhoneNumber))
				verified, _ := kcUserRep.GetAttributeBool(constants.AttrbPhoneNumberVerified)
				assert.Equal(t, phoneNumberVerified, *verified)
				assert.Equal(t, label, *kcUserRep.GetAttributeString(constants.AttrbLabel))
				assert.Equal(t, gender, *kcUserRep.GetAttributeString(constants.AttrbGender))
				assert.Equal(t, birthDate, *kcUserRep.GetAttributeString(constants.AttrbBirthDate))
				assert.Equal(t, locale, *kcUserRep.GetAttributeString(constants.AttrbLocale))
				return nil
			}).Times(1)

		err := managementComponent.UpdateUser(ctx, "master", id, userRep)

		assert.Nil(t, err)
	})

	t.Run("Update user with succces (with user info update)", func(t *testing.T) {
		mockKeycloakClient.EXPECT().GetUser(accessToken, realmName, id).Return(kcUserRep, nil).Times(2)
		mockUsersDetailsDBModule.EXPECT().GetUserDetails(ctx, realmName, id).Return(dbUserRep, nil).Times(2)

		mockKeycloakClient.EXPECT().UpdateUser(accessToken, realmName, id, gomock.Any()).DoAndReturn(
			func(accessToken, realmName, id string, kcUserRep kc.UserRepresentation) error {
				assert.Equal(t, username, *kcUserRep.Username)
				assert.Equal(t, email, *kcUserRep.Email)
				assert.Equal(t, emailVerified, *kcUserRep.EmailVerified)
				assert.Equal(t, firstName, *kcUserRep.FirstName)
				assert.Equal(t, lastName, *kcUserRep.LastName)
				assert.Equal(t, phoneNumber, *kcUserRep.GetAttributeString(constants.AttrbPhoneNumber))
				verified, _ := kcUserRep.GetAttributeBool(constants.AttrbPhoneNumberVerified)
				assert.Equal(t, phoneNumberVerified, *verified)
				assert.Equal(t, label, *kcUserRep.GetAttributeString(constants.AttrbLabel))
				assert.Equal(t, gender, *kcUserRep.GetAttributeString(constants.AttrbGender))
				assert.Equal(t, birthDate, *kcUserRep.GetAttributeString(constants.AttrbBirthDate))
				assert.Equal(t, locale, *kcUserRep.GetAttributeString(constants.AttrbLocale))
				return nil
			}).Times(2)

		newIDDocumentExpiration := "21.12.2030"
		var userAPI = api.UserRepresentation{
			Username:             &username,
			Email:                &email,
			EmailVerified:        &emailVerified,
			FirstName:            &firstName,
			LastName:             &lastName,
			PhoneNumber:          &phoneNumber,
			PhoneNumberVerified:  &phoneNumberVerified,
			Label:                &label,
			Gender:               &gender,
			BirthDate:            &birthDate,
			Locale:               &locale,
			BirthLocation:        &birthLocation,
			IDDocumentExpiration: &newIDDocumentExpiration,
		}

		mockUsersDetailsDBModule.EXPECT().StoreOrUpdateUserDetails(ctx, realmName, gomock.Any()).DoAndReturn(
			func(ctx context.Context, realm string, user dto.DBUser) error {
				assert.Equal(t, id, *user.UserID)
				assert.Equal(t, birthLocation, *user.BirthLocation)
				assert.Equal(t, idDocumentType, *user.IDDocumentType)
				assert.Equal(t, idDocumentNumber, *user.IDDocumentNumber)
				assert.Equal(t, newIDDocumentExpiration, *user.IDDocumentExpiration)
				return nil
			}).Times(1)

		err := managementComponent.UpdateUser(ctx, realmName, id, userAPI)
		assert.Nil(t, err)

		newBirthLocation := "21.12.1988"
		newIDDocumentType := "Permit"
		newIDDocumentNumber := "123456frs"
		userAPI.BirthLocation = &newBirthLocation
		userAPI.IDDocumentType = &newIDDocumentType
		userAPI.IDDocumentNumber = &newIDDocumentNumber

		mockUsersDetailsDBModule.EXPECT().StoreOrUpdateUserDetails(ctx, realmName, gomock.Any()).DoAndReturn(
			func(ctx context.Context, realm string, user dto.DBUser) error {
				assert.Equal(t, id, *user.UserID)
				assert.Equal(t, newBirthLocation, *user.BirthLocation)
				assert.Equal(t, newIDDocumentType, *user.IDDocumentType)
				assert.Equal(t, newIDDocumentNumber, *user.IDDocumentNumber)
				assert.Equal(t, newIDDocumentExpiration, *user.IDDocumentExpiration)
				return nil
			}).Times(1)

		err = managementComponent.UpdateUser(ctx, realmName, id, userAPI)
		assert.Nil(t, err)
	})

	t.Run("Update by locking the user", func(t *testing.T) {
		mockKeycloakClient.EXPECT().GetUser(accessToken, realmName, id).Return(kcUserRep, nil).Times(1)
		mockUsersDetailsDBModule.EXPECT().GetUserDetails(ctx, realmName, id).Return(dbUserRep, nil).Times(1)

		enabled = false
		mockKeycloakClient.EXPECT().UpdateUser(accessToken, realmName, id, gomock.Any()).DoAndReturn(
			func(accessToken, realmName, id string, kcUserRep kc.UserRepresentation) error {
				assert.Equal(t, username, *kcUserRep.Username)
				assert.Equal(t, email, *kcUserRep.Email)
				assert.Equal(t, enabled, *kcUserRep.Enabled)
				assert.Equal(t, emailVerified, *kcUserRep.EmailVerified)
				assert.Equal(t, firstName, *kcUserRep.FirstName)
				assert.Equal(t, lastName, *kcUserRep.LastName)
				assert.Equal(t, phoneNumber, *kcUserRep.GetAttributeString(constants.AttrbPhoneNumber))
				verified, _ := kcUserRep.GetAttributeBool(constants.AttrbPhoneNumberVerified)
				assert.Equal(t, phoneNumberVerified, *verified)
				assert.Equal(t, label, *kcUserRep.GetAttributeString(constants.AttrbLabel))
				assert.Equal(t, gender, *kcUserRep.GetAttributeString(constants.AttrbGender))
				assert.Equal(t, birthDate, *kcUserRep.GetAttributeString(constants.AttrbBirthDate))
				assert.Equal(t, locale, *kcUserRep.GetAttributeString(constants.AttrbLocale))
				return nil
			}).Times(1)

		var userRepLocked = api.UserRepresentation{
			Username:            &username,
			Email:               &email,
			Enabled:             &enabled,
			EmailVerified:       &emailVerified,
			FirstName:           &firstName,
			LastName:            &lastName,
			PhoneNumber:         &phoneNumber,
			PhoneNumberVerified: &phoneNumberVerified,
			Label:               &label,
			Gender:              &gender,
			BirthDate:           &birthDate,
			Locale:              &locale,
		}

		mockEventDBModule.EXPECT().ReportEvent(ctx, "LOCK_ACCOUNT", "back-office", gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)

		err := managementComponent.UpdateUser(ctx, "master", id, userRepLocked)

		assert.Nil(t, err)
	})

	t.Run("Update to unlock the user", func(t *testing.T) {
		mockKeycloakClient.EXPECT().GetUser(accessToken, realmName, id).Return(kcUserRep, nil).Times(1)
		mockUsersDetailsDBModule.EXPECT().GetUserDetails(ctx, realmName, id).Return(dbUserRep, nil).Times(1)

		enabled = true
		mockKeycloakClient.EXPECT().UpdateUser(accessToken, realmName, id, gomock.Any()).Return(nil).Times(1)

		var userRepLocked = api.UserRepresentation{
			Username:            &username,
			Email:               &email,
			Enabled:             &enabled,
			EmailVerified:       &emailVerified,
			FirstName:           &firstName,
			LastName:            &lastName,
			PhoneNumber:         &phoneNumber,
			PhoneNumberVerified: &phoneNumberVerified,
			Label:               &label,
			Gender:              &gender,
			BirthDate:           &birthDate,
			Locale:              &locale,
		}

		mockEventDBModule.EXPECT().ReportEvent(ctx, "UNLOCK_ACCOUNT", "back-office", gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)

		err := managementComponent.UpdateUser(ctx, "master", id, userRepLocked)

		assert.Nil(t, err)
	})

	t.Run("Update by changing the email address", func(t *testing.T) {
		var oldEmail = "toti@elca.ch"
		var oldkcUserRep = kc.UserRepresentation{
			ID:            &id,
			Email:         &oldEmail,
			EmailVerified: &emailVerified,
		}
		mockKeycloakClient.EXPECT().GetUser(accessToken, realmName, id).Return(oldkcUserRep, nil).Times(1)
		mockUsersDetailsDBModule.EXPECT().GetUserDetails(ctx, realmName, id).Return(dbUserRep, nil).Times(1)
		mockKeycloakClient.EXPECT().UpdateUser(accessToken, realmName, id, gomock.Any()).DoAndReturn(
			func(accessToken, realmName, id string, kcUserRep kc.UserRepresentation) error {
				assert.Equal(t, email, *kcUserRep.Email)
				assert.Equal(t, false, *kcUserRep.EmailVerified)
				return nil
			}).Times(1)

		err := managementComponent.UpdateUser(ctx, "master", id, userRep)

		assert.Nil(t, err)
	})

	t.Run("Update by changing the phone number", func(t *testing.T) {
		var oldNumber = "+41789467"
		var oldAttributes = make(kc.Attributes)
		oldAttributes.SetString(constants.AttrbPhoneNumber, oldNumber)
		oldAttributes.SetBool(constants.AttrbPhoneNumberVerified, phoneNumberVerified)
		var oldkcUserRep2 = kc.UserRepresentation{
			ID:         &id,
			Attributes: &oldAttributes,
		}
		mockKeycloakClient.EXPECT().GetUser(accessToken, realmName, id).Return(oldkcUserRep2, nil).Times(1)
		mockUsersDetailsDBModule.EXPECT().GetUserDetails(ctx, realmName, id).Return(dbUserRep, nil).Times(1)
		mockKeycloakClient.EXPECT().UpdateUser(accessToken, realmName, id, gomock.Any()).DoAndReturn(
			func(accessToken, realmName, id string, kcUserRep kc.UserRepresentation) error {
				verified, _ := kcUserRep.GetAttributeBool(constants.AttrbPhoneNumberVerified)
				assert.Equal(t, phoneNumber, *kcUserRep.GetAttributeString(constants.AttrbPhoneNumber))
				assert.Equal(t, false, *verified)
				return nil
			}).Times(1)

		err := managementComponent.UpdateUser(ctx, "master", id, userRep)

		assert.Nil(t, err)
	})

	t.Run("Update without attributes", func(t *testing.T) {
		var userRepWithoutAttr = api.UserRepresentation{
			Username:  &username,
			Email:     &email,
			FirstName: &firstName,
			LastName:  &lastName,
		}

		var oldNumber = "+41789467"
		var oldAttributes = make(kc.Attributes)
		oldAttributes.SetString(constants.AttrbPhoneNumber, oldNumber)
		oldAttributes.SetBool(constants.AttrbPhoneNumberVerified, phoneNumberVerified)
		var oldkcUserRep2 = kc.UserRepresentation{
			ID:         &id,
			Attributes: &oldAttributes,
		}

		mockKeycloakClient.EXPECT().GetUser(accessToken, realmName, id).Return(oldkcUserRep2, nil).Times(1)
		mockUsersDetailsDBModule.EXPECT().GetUserDetails(ctx, realmName, id).Return(dbUserRep, nil).Times(1)
		mockKeycloakClient.EXPECT().UpdateUser(accessToken, realmName, id, gomock.Any()).DoAndReturn(
			func(accessToken, realmName, id string, kcUserRep kc.UserRepresentation) error {
				verified, _ := kcUserRep.GetAttributeBool(constants.AttrbPhoneNumberVerified)
				assert.Equal(t, oldNumber, *kcUserRep.GetAttributeString(constants.AttrbPhoneNumber))
				assert.Equal(t, true, *verified)
				return nil
			}).Times(1)

		err := managementComponent.UpdateUser(ctx, "master", id, userRepWithoutAttr)

		assert.Nil(t, err)
	})

	t.Run("Update user with succces but with error when storing the event in the DB", func(t *testing.T) {
		enabled = true
		var kcUserRep = kc.UserRepresentation{
			ID:       &id,
			Username: &username,
			Enabled:  &enabled,
		}

		var userRep = api.UserRepresentation{
			Username: &username,
			Enabled:  &enabled,
		}

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		ctx = context.WithValue(ctx, cs.CtContextRealm, realmName)
		ctx = context.WithValue(ctx, cs.CtContextUsername, username)

		mockEventDBModule.EXPECT().ReportEvent(ctx, "UNLOCK_ACCOUNT", "back-office", database.CtEventRealmName, realmName, database.CtEventUserID, id, database.CtEventUsername, username).Return(errors.New("error")).Times(1)
		m := map[string]interface{}{"event_name": "UNLOCK_ACCOUNT", database.CtEventRealmName: realmName, database.CtEventUserID: id, database.CtEventUsername: username}
		eventJSON, _ := json.Marshal(m)
		mockLogger.EXPECT().Error(ctx, "err", "error", "event", string(eventJSON))
		mockKeycloakClient.EXPECT().GetUser(accessToken, realmName, id).Return(kcUserRep, nil).Times(1)
		mockUsersDetailsDBModule.EXPECT().GetUserDetails(ctx, realmName, id).Return(dbUserRep, nil).Times(1)
		mockKeycloakClient.EXPECT().UpdateUser(accessToken, realmName, id, gomock.Any()).Return(nil).Times(1)

		err := managementComponent.UpdateUser(ctx, "master", id, userRep)

		assert.Nil(t, err)
	})

	t.Run("Error - get user KC", func(t *testing.T) {
		var id = "1234-79894-7594"
		mockKeycloakClient.EXPECT().GetUser(accessToken, realmName, id).Return(kc.UserRepresentation{}, fmt.Errorf("Unexpected error")).Times(1)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		mockLogger.EXPECT().Warn(ctx, "err", "Unexpected error")
		err := managementComponent.UpdateUser(ctx, "master", id, api.UserRepresentation{})

		assert.NotNil(t, err)
	})

	t.Run("Error - get user info from DB", func(t *testing.T) {
		var id = "1234-79894-7594"
		mockKeycloakClient.EXPECT().GetUser(accessToken, realmName, id).Return(kc.UserRepresentation{}, nil).Times(1)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		mockUsersDetailsDBModule.EXPECT().GetUserDetails(ctx, realmName, id).Return(dto.DBUser{}, fmt.Errorf("SQL Error")).Times(1)

		err := managementComponent.UpdateUser(ctx, "master", id, api.UserRepresentation{})

		assert.NotNil(t, err)
	})

	t.Run("Error - update user KC", func(t *testing.T) {
		var id = "1234-79894-7594"
		var kcUserRep = kc.UserRepresentation{
			ID: &id,
		}
		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)

		mockKeycloakClient.EXPECT().GetUser(accessToken, realmName, id).Return(kcUserRep, nil).AnyTimes()
		mockUsersDetailsDBModule.EXPECT().GetUserDetails(ctx, realmName, id).Return(dbUserRep, nil).Times(1)
		mockKeycloakClient.EXPECT().UpdateUser(accessToken, realmName, id, gomock.Any()).Return(fmt.Errorf("Unexpected error")).Times(1)
		mockLogger.EXPECT().Warn(gomock.Any(), "err", "Unexpected error")

		err := managementComponent.UpdateUser(ctx, "master", id, api.UserRepresentation{})

		assert.NotNil(t, err)
	})

	t.Run("Error - update user info in DB", func(t *testing.T) {
		var id = "1234-79894-7594"
		var kcUserRep = kc.UserRepresentation{
			ID: &id,
		}
		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)

		mockKeycloakClient.EXPECT().GetUser(accessToken, realmName, id).Return(kcUserRep, nil).AnyTimes()
		mockUsersDetailsDBModule.EXPECT().GetUserDetails(ctx, realmName, id).Return(dbUserRep, nil).Times(1)
		mockKeycloakClient.EXPECT().UpdateUser(accessToken, realmName, id, gomock.Any()).Return(nil).Times(1)
		mockUsersDetailsDBModule.EXPECT().StoreOrUpdateUserDetails(ctx, realmName, gomock.Any()).Return(fmt.Errorf("SQL error")).Times(1)
		mockLogger.EXPECT().Warn(gomock.Any(), "msg", "Can't store user details in database", "err", "SQL error")

		var newIDDocumentType = "Visa"
		err := managementComponent.UpdateUser(ctx, realmName, id, api.UserRepresentation{
			IDDocumentExpiration: &newIDDocumentType,
		})

		assert.NotNil(t, err)
	})
}

func TestLockUnlockUser(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()

	var mockKeycloakClient = mock.NewKeycloakClient(mockCtrl)
	var mockEventDBModule = mock.NewEventDBModule(mockCtrl)

	var managementComponent = NewComponent(mockKeycloakClient, nil, mockEventDBModule, nil, nil, log.NewNopLogger())

	var accessToken = "TOKEN=="
	var realmName = "myrealm"
	var userID = "41dbf4a8-32a9-4000-8c17-edc854c31231"
	var anyError = errors.New("any")
	var bTrue = true
	var bFalse = false
	var ctx = context.TODO()
	ctx = context.WithValue(ctx, cs.CtContextAccessToken, accessToken)

	t.Run("GetUser fails", func(t *testing.T) {
		mockKeycloakClient.EXPECT().GetUser(accessToken, realmName, userID).Return(kc.UserRepresentation{}, anyError)
		var err = managementComponent.LockUser(ctx, realmName, userID)
		assert.Equal(t, anyError, err)
	})
	t.Run("Can't lock disabled user", func(t *testing.T) {
		mockKeycloakClient.EXPECT().GetUser(accessToken, realmName, userID).Return(kc.UserRepresentation{Enabled: &bFalse}, nil)
		var err = managementComponent.LockUser(ctx, realmName, userID)
		assert.Nil(t, err)
	})
	t.Run("UpdateUser fails", func(t *testing.T) {
		mockKeycloakClient.EXPECT().GetUser(accessToken, realmName, userID).Return(kc.UserRepresentation{Enabled: &bFalse}, nil)
		mockKeycloakClient.EXPECT().UpdateUser(accessToken, realmName, userID, gomock.Any()).Return(anyError)
		var err = managementComponent.UnlockUser(ctx, realmName, userID)
		assert.Equal(t, anyError, err)
	})
	t.Run("Lock success", func(t *testing.T) {
		mockKeycloakClient.EXPECT().GetUser(accessToken, realmName, userID).Return(kc.UserRepresentation{Enabled: &bTrue}, nil)
		mockKeycloakClient.EXPECT().UpdateUser(accessToken, realmName, userID, gomock.Any()).Return(nil)
		mockEventDBModule.EXPECT().ReportEvent(ctx, "LOCK_ACCOUNT", "back-office", gomock.Any()).Return(nil).Times(1)
		var err = managementComponent.LockUser(ctx, realmName, userID)
		assert.Nil(t, err)
	})
	t.Run("Unlock success", func(t *testing.T) {
		mockKeycloakClient.EXPECT().GetUser(accessToken, realmName, userID).Return(kc.UserRepresentation{Enabled: &bFalse}, nil)
		mockKeycloakClient.EXPECT().UpdateUser(accessToken, realmName, userID, gomock.Any()).Return(nil)
		mockEventDBModule.EXPECT().ReportEvent(ctx, "UNLOCK_ACCOUNT", "back-office", gomock.Any()).Return(nil).Times(1)
		var err = managementComponent.UnlockUser(ctx, realmName, userID)
		assert.Nil(t, err)
	})
}

func TestGetUsers(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()
	var mockKeycloakClient = mock.NewKeycloakClient(mockCtrl)
	var mockUsersDetailsDBModule = mock.NewUsersDetailsDBModule(mockCtrl)
	var mockEventDBModule = mock.NewEventDBModule(mockCtrl)
	var mockConfigurationDBModule = mock.NewConfigurationDBModule(mockCtrl)
	var mockLogger = log.NewNopLogger()
	var allowedTrustIDGroups = []string{"grp1", "grp2"}

	var managementComponent = NewComponent(mockKeycloakClient, mockUsersDetailsDBModule, mockEventDBModule, mockConfigurationDBModule, allowedTrustIDGroups, mockLogger)

	var accessToken = "TOKEN=="
	var realmName = "master"
	var targetRealmName = "DEP"

	// Get user with succces
	{
		var id = "1234-7454-4516"
		var username = "username"
		var email = "toto@elca.ch"
		var enabled = true
		var emailVerified = true
		var firstName = "Titi"
		var lastName = "Tutu"
		var phoneNumber = "+41789456"
		var phoneNumberVerified = true
		var label = "Label"
		var gender = "M"
		var birthDate = "01/01/1988"
		var createdTimestamp = time.Now().UTC().Unix()

		var attributes = make(kc.Attributes)
		attributes.SetString(constants.AttrbPhoneNumber, phoneNumber)
		attributes.SetString(constants.AttrbLabel, label)
		attributes.SetString(constants.AttrbGender, gender)
		attributes.SetString(constants.AttrbBirthDate, birthDate)
		attributes.SetBool(constants.AttrbPhoneNumberVerified, phoneNumberVerified)

		var count = 10
		var kcUserRep = kc.UserRepresentation{
			ID:               &id,
			Username:         &username,
			Email:            &email,
			Enabled:          &enabled,
			EmailVerified:    &emailVerified,
			FirstName:        &firstName,
			LastName:         &lastName,
			Attributes:       &attributes,
			CreatedTimestamp: &createdTimestamp,
		}
		var kcUsersRep = kc.UsersPageRepresentation{
			Count: &count,
			Users: []kc.UserRepresentation{kcUserRep},
		}

		mockKeycloakClient.EXPECT().GetUsers(accessToken, realmName, targetRealmName, "groupId", "123-456-789").Return(kcUsersRep, nil).Times(1)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		ctx = context.WithValue(ctx, cs.CtContextRealm, "master")

		apiUsersRep, err := managementComponent.GetUsers(ctx, "DEP", []string{"123-456-789"})

		var apiUserRep = apiUsersRep.Users[0]
		assert.Nil(t, err)
		assert.Equal(t, username, *apiUserRep.Username)
		assert.Equal(t, email, *apiUserRep.Email)
		assert.Equal(t, enabled, *apiUserRep.Enabled)
		assert.Equal(t, emailVerified, *apiUserRep.EmailVerified)
		assert.Equal(t, firstName, *apiUserRep.FirstName)
		assert.Equal(t, lastName, *apiUserRep.LastName)
		assert.Equal(t, phoneNumber, *apiUserRep.PhoneNumber)
		verified, _ := strconv.ParseBool(((*kcUserRep.Attributes)["phoneNumberVerified"][0]))
		assert.Equal(t, phoneNumberVerified, verified)
		assert.Equal(t, label, *apiUserRep.Label)
		assert.Equal(t, gender, *apiUserRep.Gender)
		assert.Equal(t, birthDate, *apiUserRep.BirthDate)
		assert.Equal(t, createdTimestamp, *apiUserRep.CreatedTimestamp)
	}

	//Error
	{
		mockKeycloakClient.EXPECT().GetUsers(accessToken, realmName, targetRealmName, "groupId", "123-456-789").Return(kc.UsersPageRepresentation{}, fmt.Errorf("Unexpected error")).Times(1)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		ctx = context.WithValue(ctx, cs.CtContextRealm, "master")

		_, err := managementComponent.GetUsers(ctx, "DEP", []string{"123-456-789"})

		assert.NotNil(t, err)
	}
}

func TestGetUserAccountStatus(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()
	var mockKeycloakClient = mock.NewKeycloakClient(mockCtrl)
	var mockUsersDetailsDBModule = mock.NewUsersDetailsDBModule(mockCtrl)
	var mockEventDBModule = mock.NewEventDBModule(mockCtrl)
	var mockConfigurationDBModule = mock.NewConfigurationDBModule(mockCtrl)
	var mockLogger = log.NewNopLogger()
	var allowedTrustIDGroups = []string{"grp1", "grp2"}

	var managementComponent = NewComponent(mockKeycloakClient, mockUsersDetailsDBModule, mockEventDBModule, mockConfigurationDBModule, allowedTrustIDGroups, mockLogger)

	var accessToken = "TOKEN=="
	var realmReq = "master"
	var realmName = "aRealm"
	var userID = "789-789-456"

	// GetUser returns an error
	{
		var userRep kc.UserRepresentation
		mockKeycloakClient.EXPECT().GetUser(accessToken, realmName, userID).Return(userRep, fmt.Errorf("Unexpected error")).Times(1)
		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		_, err := managementComponent.GetUserAccountStatus(ctx, realmName, userID)
		assert.NotNil(t, err)
	}

	// GetUser returns a non-enabled user
	{
		var userRep kc.UserRepresentation
		enabled := false
		userRep.Enabled = &enabled
		mockKeycloakClient.EXPECT().GetUser(accessToken, realmName, userID).Return(userRep, nil).Times(1)
		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		status, err := managementComponent.GetUserAccountStatus(ctx, realmName, userID)
		assert.Nil(t, err)
		assert.False(t, status["enabled"])
	}

	// GetUser returns an enabled user but GetCredentialsForUser fails
	{
		var userRep kc.UserRepresentation
		enabled := true
		userRep.Enabled = &enabled
		mockKeycloakClient.EXPECT().GetUser(accessToken, realmName, userID).Return(userRep, nil).Times(1)
		mockKeycloakClient.EXPECT().GetCredentials(accessToken, realmName, userID).Return(nil, fmt.Errorf("Unexpected error")).Times(1)
		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		ctx = context.WithValue(ctx, cs.CtContextRealm, realmReq)
		_, err := managementComponent.GetUserAccountStatus(ctx, realmName, userID)
		assert.NotNil(t, err)
	}

	// GetUser returns an enabled user but GetCredentialsForUser have no credential
	{
		var userRep kc.UserRepresentation
		enabled := true
		userRep.Enabled = &enabled
		mockKeycloakClient.EXPECT().GetUser(accessToken, realmName, userID).Return(userRep, nil).Times(1)
		mockKeycloakClient.EXPECT().GetCredentials(accessToken, realmName, userID).Return([]kc.CredentialRepresentation{}, nil).Times(1)
		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		ctx = context.WithValue(ctx, cs.CtContextRealm, realmReq)
		status, err := managementComponent.GetUserAccountStatus(ctx, realmName, userID)
		assert.Nil(t, err)
		assert.False(t, status["enabled"])
	}

	// GetUser returns an enabled user and GetCredentialsForUser have credentials
	{
		var userRep kc.UserRepresentation
		var creds1, creds2 kc.CredentialRepresentation
		enabled := true
		userRep.Enabled = &enabled
		mockKeycloakClient.EXPECT().GetUser(accessToken, realmName, userID).Return(userRep, nil).Times(1)
		mockKeycloakClient.EXPECT().GetCredentials(accessToken, realmName, userID).Return([]kc.CredentialRepresentation{creds1, creds2}, nil).Times(1)
		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		ctx = context.WithValue(ctx, cs.CtContextRealm, realmReq)
		status, err := managementComponent.GetUserAccountStatus(ctx, realmName, userID)
		assert.Nil(t, err)
		assert.True(t, status["enabled"])
	}
}

func TestGetClientRolesForUser(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()
	var mockKeycloakClient = mock.NewKeycloakClient(mockCtrl)
	var mockUsersDetailsDBModule = mock.NewUsersDetailsDBModule(mockCtrl)
	var mockEventDBModule = mock.NewEventDBModule(mockCtrl)
	var mockConfigurationDBModule = mock.NewConfigurationDBModule(mockCtrl)
	var mockLogger = log.NewNopLogger()
	var allowedTrustIDGroups = []string{"grp1", "grp2"}

	var managementComponent = NewComponent(mockKeycloakClient, mockUsersDetailsDBModule, mockEventDBModule, mockConfigurationDBModule, allowedTrustIDGroups, mockLogger)

	var accessToken = "TOKEN=="
	var realmName = "master"
	var userID = "789-789-456"
	var clientID = "456-789-147"

	// Get role with succces
	{
		var id = "1234-7454-4516"
		var composite = false
		var containerID = "containerId"
		var description = "description role"
		var clientRole = true
		var name = "client name"

		var kcRoleRep = kc.RoleRepresentation{
			ID:          &id,
			Name:        &name,
			ClientRole:  &clientRole,
			Composite:   &composite,
			ContainerID: &containerID,
			Description: &description,
		}

		var kcRolesRep []kc.RoleRepresentation
		kcRolesRep = append(kcRolesRep, kcRoleRep)

		mockKeycloakClient.EXPECT().GetClientRoleMappings(accessToken, realmName, userID, clientID).Return(kcRolesRep, nil).Times(1)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)

		apiRolesRep, err := managementComponent.GetClientRolesForUser(ctx, "master", userID, clientID)

		var apiRoleRep = apiRolesRep[0]
		assert.Nil(t, err)
		assert.Equal(t, id, *apiRoleRep.ID)
		assert.Equal(t, name, *apiRoleRep.Name)
		assert.Equal(t, clientRole, *apiRoleRep.ClientRole)
		assert.Equal(t, composite, *apiRoleRep.Composite)
		assert.Equal(t, containerID, *apiRoleRep.ContainerID)
		assert.Equal(t, description, *apiRoleRep.Description)
	}

	//Error
	{
		mockKeycloakClient.EXPECT().GetClientRoleMappings(accessToken, realmName, userID, clientID).Return([]kc.RoleRepresentation{}, fmt.Errorf("Unexpected error")).Times(1)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)

		_, err := managementComponent.GetClientRolesForUser(ctx, "master", userID, clientID)

		assert.NotNil(t, err)
	}
}

func TestAddClientRolesToUser(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()
	var mockKeycloakClient = mock.NewKeycloakClient(mockCtrl)
	var mockUsersDetailsDBModule = mock.NewUsersDetailsDBModule(mockCtrl)
	var mockEventDBModule = mock.NewEventDBModule(mockCtrl)
	var mockConfigurationDBModule = mock.NewConfigurationDBModule(mockCtrl)
	var mockLogger = log.NewNopLogger()
	var allowedTrustIDGroups = []string{"grp1", "grp2"}

	var managementComponent = NewComponent(mockKeycloakClient, mockUsersDetailsDBModule, mockEventDBModule, mockConfigurationDBModule, allowedTrustIDGroups, mockLogger)

	var accessToken = "TOKEN=="
	var realmName = "master"
	var userID = "789-789-456"
	var clientID = "456-789-147"

	// Add role with succces
	{
		var id = "1234-7454-4516"
		var composite = false
		var containerID = "containerId"
		var description = "description role"
		var clientRole = true
		var name = "client name"

		mockKeycloakClient.EXPECT().AddClientRolesToUserRoleMapping(accessToken, realmName, userID, clientID, gomock.Any()).DoAndReturn(
			func(accessToken, realmName, userID, clientID string, roles []kc.RoleRepresentation) error {
				var role = roles[0]
				assert.Equal(t, id, *role.ID)
				assert.Equal(t, name, *role.Name)
				assert.Equal(t, clientRole, *role.ClientRole)
				assert.Equal(t, composite, *role.Composite)
				assert.Equal(t, containerID, *role.ContainerID)
				assert.Equal(t, description, *role.Description)
				return nil
			}).Times(1)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)

		var roleRep = api.RoleRepresentation{
			ID:          &id,
			Name:        &name,
			ClientRole:  &clientRole,
			Composite:   &composite,
			ContainerID: &containerID,
			Description: &description,
		}
		var rolesRep []api.RoleRepresentation
		rolesRep = append(rolesRep, roleRep)

		err := managementComponent.AddClientRolesToUser(ctx, "master", userID, clientID, rolesRep)

		assert.Nil(t, err)
	}

	//Error
	{
		mockKeycloakClient.EXPECT().AddClientRolesToUserRoleMapping(accessToken, realmName, userID, clientID, gomock.Any()).Return(fmt.Errorf("Unexpected error")).Times(1)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)

		err := managementComponent.AddClientRolesToUser(ctx, "master", userID, clientID, []api.RoleRepresentation{})

		assert.NotNil(t, err)
	}
}

func TestGetRolesOfUser(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()
	var mockKeycloakClient = mock.NewKeycloakClient(mockCtrl)
	var mockUsersDetailsDBModule = mock.NewUsersDetailsDBModule(mockCtrl)
	var mockEventDBModule = mock.NewEventDBModule(mockCtrl)
	var mockConfigurationDBModule = mock.NewConfigurationDBModule(mockCtrl)
	var mockLogger = log.NewNopLogger()
	var allowedTrustIDGroups = []string{"grp1", "grp2"}

	var managementComponent = NewComponent(mockKeycloakClient, mockUsersDetailsDBModule, mockEventDBModule, mockConfigurationDBModule, allowedTrustIDGroups, mockLogger)

	var accessToken = "TOKEN=="
	var realmName = "master"
	var userID = "789-789-456"

	// Get role with succces
	{
		var id = "1234-7454-4516"
		var composite = false
		var containerID = "containerId"
		var description = "description role"
		var clientRole = false
		var name = "client name"

		var kcRoleRep = kc.RoleRepresentation{
			ID:          &id,
			Name:        &name,
			ClientRole:  &clientRole,
			Composite:   &composite,
			ContainerID: &containerID,
			Description: &description,
		}

		var kcRolesRep []kc.RoleRepresentation
		kcRolesRep = append(kcRolesRep, kcRoleRep)

		mockKeycloakClient.EXPECT().GetRealmLevelRoleMappings(accessToken, realmName, userID).Return(kcRolesRep, nil).Times(1)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)

		apiRolesRep, err := managementComponent.GetRolesOfUser(ctx, "master", userID)

		var apiRoleRep = apiRolesRep[0]
		assert.Nil(t, err)
		assert.Equal(t, id, *apiRoleRep.ID)
		assert.Equal(t, name, *apiRoleRep.Name)
		assert.Equal(t, clientRole, *apiRoleRep.ClientRole)
		assert.Equal(t, composite, *apiRoleRep.Composite)
		assert.Equal(t, containerID, *apiRoleRep.ContainerID)
		assert.Equal(t, description, *apiRoleRep.Description)
	}

	//Error
	{
		mockKeycloakClient.EXPECT().GetRealmLevelRoleMappings(accessToken, realmName, userID).Return([]kc.RoleRepresentation{}, fmt.Errorf("Unexpected error")).Times(1)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)

		_, err := managementComponent.GetRolesOfUser(ctx, "master", userID)

		assert.NotNil(t, err)
	}
}

func TestGetGroupsOfUser(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()
	var mockKeycloakClient = mock.NewKeycloakClient(mockCtrl)
	var mockUsersDetailsDBModule = mock.NewUsersDetailsDBModule(mockCtrl)
	var mockEventDBModule = mock.NewEventDBModule(mockCtrl)
	var mockConfigurationDBModule = mock.NewConfigurationDBModule(mockCtrl)
	var mockLogger = log.NewNopLogger()
	var allowedTrustIDGroups = []string{"grp1", "grp2"}

	var managementComponent = NewComponent(mockKeycloakClient, mockUsersDetailsDBModule, mockEventDBModule, mockConfigurationDBModule, allowedTrustIDGroups, mockLogger)

	var accessToken = "TOKEN=="
	var realmName = "master"
	var userID = "789-789-456"

	// Get groups with succces
	{
		var id = "1234-7454-4516"
		var name = "client name"

		var kcGroupRep = kc.GroupRepresentation{
			ID:   &id,
			Name: &name,
		}

		var kcGroupsRep []kc.GroupRepresentation
		kcGroupsRep = append(kcGroupsRep, kcGroupRep)

		mockKeycloakClient.EXPECT().GetGroupsOfUser(accessToken, realmName, userID).Return(kcGroupsRep, nil).Times(1)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)

		apiGroupsRep, err := managementComponent.GetGroupsOfUser(ctx, "master", userID)

		var apiGroupRep = apiGroupsRep[0]
		assert.Nil(t, err)
		assert.Equal(t, id, *apiGroupRep.ID)
		assert.Equal(t, name, *apiGroupRep.Name)
	}

	//Error
	{
		mockKeycloakClient.EXPECT().GetGroupsOfUser(accessToken, realmName, userID).Return([]kc.GroupRepresentation{}, fmt.Errorf("Unexpected error")).Times(1)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)

		_, err := managementComponent.GetGroupsOfUser(ctx, "master", userID)

		assert.NotNil(t, err)
	}
}

func TestSetGroupsToUser(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()

	var mockKeycloakClient = mock.NewKeycloakClient(mockCtrl)
	var mockUsersDetailsDBModule = mock.NewUsersDetailsDBModule(mockCtrl)
	var mockEventDBModule = mock.NewEventDBModule(mockCtrl)
	var mockConfigurationDBModule = mock.NewConfigurationDBModule(mockCtrl)
	var mockLogger = log.NewNopLogger()

	var accessToken = "a-valid-access-token"
	var realmName = "my-realm"
	var userID = "USER-IDEN-IFIE-R123"
	var allowedTrustIDGroups = []string{"won't be used"}
	var groupID = "user-group-1"
	var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)

	var managementComponent = NewComponent(mockKeycloakClient, mockUsersDetailsDBModule, mockEventDBModule, mockConfigurationDBModule, allowedTrustIDGroups, mockLogger)

	t.Run("AddGroupToUser: KC fails", func(t *testing.T) {
		mockKeycloakClient.EXPECT().AddGroupToUser(accessToken, realmName, userID, groupID).Return(errors.New("kc error"))
		var err = managementComponent.AddGroupToUser(ctx, realmName, userID, groupID)
		assert.NotNil(t, err)
	})
	t.Run("DeleteGroupForUser: KC fails", func(t *testing.T) {
		mockKeycloakClient.EXPECT().DeleteGroupFromUser(accessToken, realmName, userID, groupID).Return(errors.New("kc error"))
		var err = managementComponent.DeleteGroupForUser(ctx, realmName, userID, groupID)
		assert.NotNil(t, err)
	})
	t.Run("AddGroupToUser: Success", func(t *testing.T) {
		mockKeycloakClient.EXPECT().AddGroupToUser(accessToken, realmName, userID, groupID).Return(nil)
		var err = managementComponent.AddGroupToUser(ctx, realmName, userID, groupID)
		assert.Nil(t, err)
	})
	t.Run("DeleteGroupForUser: Success", func(t *testing.T) {
		mockKeycloakClient.EXPECT().DeleteGroupFromUser(accessToken, realmName, userID, groupID).Return(nil)
		var err = managementComponent.DeleteGroupForUser(ctx, realmName, userID, groupID)
		assert.Nil(t, err)
	})
}

func TestGetAvailableTrustIDGroups(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()

	var mockKeycloakClient = mock.NewKeycloakClient(mockCtrl)
	var mockUsersDetailsDBModule = mock.NewUsersDetailsDBModule(mockCtrl)
	var mockEventDBModule = mock.NewEventDBModule(mockCtrl)
	var mockConfigurationDBModule = mock.NewConfigurationDBModule(mockCtrl)
	var mockLogger = log.NewNopLogger()

	var allowedTrustIDGroups = []string{"grp1", "grp2"}
	var realmName = "master"

	var component = NewComponent(mockKeycloakClient, mockUsersDetailsDBModule, mockEventDBModule, mockConfigurationDBModule, allowedTrustIDGroups, mockLogger)

	var res, err = component.GetAvailableTrustIDGroups(context.TODO(), realmName)
	assert.Nil(t, err)
	assert.Len(t, res, len(allowedTrustIDGroups))
}

func TestGetTrustIDGroupsOfUser(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()

	var mockKeycloakClient = mock.NewKeycloakClient(mockCtrl)
	var mockUsersDetailsDBModule = mock.NewUsersDetailsDBModule(mockCtrl)
	var mockEventDBModule = mock.NewEventDBModule(mockCtrl)
	var mockConfigurationDBModule = mock.NewConfigurationDBModule(mockCtrl)
	var mockLogger = log.NewNopLogger()

	var allowedTrustIDGroups = []string{"grp1", "grp2"}
	var groups = []string{"some", "/groups"}
	var accessToken = "TOKEN=="
	var realmName = "master"
	var userID = "789-789-456"
	var attrbs = keycloak.Attributes{constants.AttrbTrustIDGroups: groups}
	var ctx = context.WithValue(context.TODO(), cs.CtContextAccessToken, accessToken)

	var component = NewComponent(mockKeycloakClient, mockUsersDetailsDBModule, mockEventDBModule, mockConfigurationDBModule, allowedTrustIDGroups, mockLogger)

	t.Run("Keycloak fails", func(t *testing.T) {
		mockKeycloakClient.EXPECT().GetUser(accessToken, realmName, userID).Return(kc.UserRepresentation{}, errors.New("kc error"))
		var _, err = component.GetTrustIDGroupsOfUser(ctx, realmName, userID)
		assert.NotNil(t, err)
	})
	t.Run("User without attributes", func(t *testing.T) {
		mockKeycloakClient.EXPECT().GetUser(accessToken, realmName, userID).Return(kc.UserRepresentation{}, nil)
		var res, err = component.GetTrustIDGroupsOfUser(ctx, realmName, userID)
		assert.Nil(t, err)
		assert.Len(t, res, 0)
	})
	t.Run("User has attributes", func(t *testing.T) {
		mockKeycloakClient.EXPECT().GetUser(accessToken, realmName, userID).Return(kc.UserRepresentation{Attributes: &attrbs}, nil)
		var res, err = component.GetTrustIDGroupsOfUser(ctx, realmName, userID)
		assert.Nil(t, err)
		assert.Equal(t, "some", res[0])
		assert.Equal(t, "groups", res[1]) // Without heading slash
	})
}

func TestSetTrustIDGroupsToUser(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()
	var mockKeycloakClient = mock.NewKeycloakClient(mockCtrl)
	var mockUsersDetailsDBModule = mock.NewUsersDetailsDBModule(mockCtrl)
	var mockEventDBModule = mock.NewEventDBModule(mockCtrl)
	var mockConfigurationDBModule = mock.NewConfigurationDBModule(mockCtrl)
	var mockLogger = log.NewNopLogger()
	var allowedTrustIDGroups = []string{"grp1", "grp2"}

	var managementComponent = NewComponent(mockKeycloakClient, mockUsersDetailsDBModule, mockEventDBModule, mockConfigurationDBModule, allowedTrustIDGroups, mockLogger)

	var accessToken = "TOKEN=="

	var username = "user"
	var realmName = "master"
	var userID = "789-1234-5678"

	t.Run("Set groups with success", func(t *testing.T) {
		var kcUserRep = kc.UserRepresentation{
			Username: &username,
		}
		grpNames := []string{"grp1", "grp2"}
		extGrpNames := []string{"/grp1", "/grp2"}
		attrs := make(kc.Attributes)
		attrs.Set(constants.AttrbTrustIDGroups, extGrpNames)
		var kcUserRep2 = kc.UserRepresentation{
			Username:   &username,
			Attributes: &attrs,
		}
		mockKeycloakClient.EXPECT().GetUser(accessToken, realmName, userID).Return(kcUserRep, nil).Times(1)
		mockKeycloakClient.EXPECT().UpdateUser(accessToken, realmName, userID, kcUserRep2).Return(nil).Times(1)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)

		err := managementComponent.SetTrustIDGroupsToUser(ctx, realmName, userID, grpNames)

		assert.Nil(t, err)
	})

	t.Run("Try to set unknown group", func(t *testing.T) {
		grpNames := []string{"grp1", "grp3"}
		attrs := make(map[string][]string)
		attrs["trustIDGroups"] = grpNames

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)

		err := managementComponent.SetTrustIDGroupsToUser(ctx, realmName, userID, grpNames)

		assert.NotNil(t, err)
	})

	t.Run("Error while get user", func(t *testing.T) {
		grpNames := []string{"grp1", "grp2"}
		mockKeycloakClient.EXPECT().GetUser(accessToken, realmName, userID).Return(kc.UserRepresentation{}, fmt.Errorf("Unexpected error")).Times(1)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)

		err := managementComponent.SetTrustIDGroupsToUser(ctx, realmName, userID, grpNames)

		assert.NotNil(t, err)
	})

	t.Run("Error while update user", func(t *testing.T) {
		var kcUserRep = kc.UserRepresentation{
			Username: &username,
		}
		grpNames := []string{"grp1", "grp2"}
		extGrpNames := []string{"/grp1", "/grp2"}
		attrs := make(kc.Attributes)
		attrs.Set(constants.AttrbTrustIDGroups, extGrpNames)
		var kcUserRep2 = kc.UserRepresentation{
			Username:   &username,
			Attributes: &attrs,
		}
		mockKeycloakClient.EXPECT().GetUser(accessToken, realmName, userID).Return(kcUserRep, nil).Times(1)
		mockKeycloakClient.EXPECT().UpdateUser(accessToken, realmName, userID, kcUserRep2).Return(fmt.Errorf("Unexpected error")).Times(1)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)

		err := managementComponent.SetTrustIDGroupsToUser(ctx, realmName, userID, grpNames)

		assert.NotNil(t, err)
	})
}

func TestResetPassword(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()
	var mockKeycloakClient = mock.NewKeycloakClient(mockCtrl)
	var mockUsersDetailsDBModule = mock.NewUsersDetailsDBModule(mockCtrl)
	var mockEventDBModule = mock.NewEventDBModule(mockCtrl)
	var mockConfigurationDBModule = mock.NewConfigurationDBModule(mockCtrl)
	var mockLogger = mock.NewLogger(mockCtrl)
	var allowedTrustIDGroups = []string{"grp1", "grp2"}

	var managementComponent = NewComponent(mockKeycloakClient, mockUsersDetailsDBModule, mockEventDBModule, mockConfigurationDBModule, allowedTrustIDGroups, mockLogger)

	var accessToken = "TOKEN=="
	var realmName = "master"
	var userID = "41dbf4a8-32a9-4000-8c17-edc854c31231"
	var password = "P@ssw0rd"
	var typePassword = "password"
	var username = "username"

	// Change password
	{
		var kcCredRep = kc.CredentialRepresentation{
			Type:  &typePassword,
			Value: &password,
		}

		mockKeycloakClient.EXPECT().ResetPassword(accessToken, realmName, userID, kcCredRep).Return(nil).Times(1)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		ctx = context.WithValue(ctx, cs.CtContextRealm, realmName)
		ctx = context.WithValue(ctx, cs.CtContextUsername, username)

		mockEventDBModule.EXPECT().ReportEvent(ctx, "INIT_PASSWORD", "back-office", gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)

		var passwordRep = api.PasswordRepresentation{
			Value: &password,
		}

		_, err := managementComponent.ResetPassword(ctx, "master", userID, passwordRep)

		assert.Nil(t, err)
	}
	// Change password but with error when storing the DB
	{
		var kcCredRep = kc.CredentialRepresentation{
			Type:  &typePassword,
			Value: &password,
		}

		mockKeycloakClient.EXPECT().ResetPassword(accessToken, realmName, userID, kcCredRep).Return(nil).Times(1)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		ctx = context.WithValue(ctx, cs.CtContextRealm, realmName)
		ctx = context.WithValue(ctx, cs.CtContextUsername, username)

		mockEventDBModule.EXPECT().ReportEvent(ctx, "INIT_PASSWORD", "back-office", database.CtEventRealmName, realmName, database.CtEventUserID, userID).Return(errors.New("error")).Times(1)
		m := map[string]interface{}{"event_name": "INIT_PASSWORD", database.CtEventRealmName: realmName, database.CtEventUserID: userID}
		eventJSON, _ := json.Marshal(m)
		mockLogger.EXPECT().Error(gomock.Any(), "err", "error", "event", string(eventJSON))
		var passwordRep = api.PasswordRepresentation{
			Value: &password,
		}

		_, err := managementComponent.ResetPassword(ctx, "master", userID, passwordRep)

		assert.Nil(t, err)
	}

	// No password offered
	{
		var id = "master_id"
		var keycloakVersion = "4.8.3"
		var realm = "master"
		var displayName = "Master"
		var enabled = true

		var policy = "forceExpiredPasswordChange(365) and specialChars(1) and upperCase(1) and lowerCase(1) and length(4) and digits(1) and notUsername(undefined)"
		var kcRealmRep = kc.RealmRepresentation{
			ID:              &id,
			KeycloakVersion: &keycloakVersion,
			Realm:           &realm,
			DisplayName:     &displayName,
			Enabled:         &enabled,
			PasswordPolicy:  &policy,
		}

		mockKeycloakClient.EXPECT().ResetPassword(accessToken, realmName, userID, gomock.Any()).Return(nil).Times(1)
		mockKeycloakClient.EXPECT().GetRealm(accessToken, realmName).Return(kcRealmRep, nil).AnyTimes()

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		ctx = context.WithValue(ctx, cs.CtContextRealm, realmName)
		ctx = context.WithValue(ctx, cs.CtContextUsername, username)

		mockEventDBModule.EXPECT().ReportEvent(ctx, "INIT_PASSWORD", "back-office", gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)

		var passwordRep = api.PasswordRepresentation{
			Value: nil,
		}

		pwd, err := managementComponent.ResetPassword(ctx, "master", userID, passwordRep)

		assert.Nil(t, err)
		assert.NotNil(t, pwd)
	}

	// No password offered, no keycloak policy
	{
		var id = "master_id"

		var kcRealmRep = kc.RealmRepresentation{
			ID: &id,
		}

		mockKeycloakClient.EXPECT().ResetPassword(accessToken, realmName, userID, gomock.Any()).Return(nil).Times(1)
		mockKeycloakClient.EXPECT().GetRealm(accessToken, realmName).Return(kcRealmRep, nil).AnyTimes()

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		ctx = context.WithValue(ctx, cs.CtContextRealm, realmName)
		ctx = context.WithValue(ctx, cs.CtContextUsername, username)

		mockEventDBModule.EXPECT().ReportEvent(ctx, "INIT_PASSWORD", "back-office", gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)

		var passwordRep = api.PasswordRepresentation{
			Value: nil,
		}

		pwd, err := managementComponent.ResetPassword(ctx, "master", userID, passwordRep)

		assert.Nil(t, err)
		assert.NotNil(t, pwd)
	}
	// Error
	{
		mockKeycloakClient.EXPECT().ResetPassword(accessToken, realmName, userID, gomock.Any()).Return(fmt.Errorf("Invalid input")).Times(1)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)

		var passwordRep = api.PasswordRepresentation{
			Value: &password,
		}
		mockLogger.EXPECT().Warn(gomock.Any(), "err", "Invalid input")
		_, err := managementComponent.ResetPassword(ctx, "master", userID, passwordRep)

		assert.NotNil(t, err)
	}

}

func TestRecoveryCode(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()
	var mockKeycloakClient = mock.NewKeycloakClient(mockCtrl)
	var mockUsersDetailsDBModule = mock.NewUsersDetailsDBModule(mockCtrl)
	var mockEventDBModule = mock.NewEventDBModule(mockCtrl)
	var mockConfigurationDBModule = mock.NewConfigurationDBModule(mockCtrl)
	var mockLogger = mock.NewLogger(mockCtrl)
	var allowedTrustIDGroups = []string{"grp1", "grp2"}

	var managementComponent = NewComponent(mockKeycloakClient, mockUsersDetailsDBModule, mockEventDBModule, mockConfigurationDBModule, allowedTrustIDGroups, mockLogger)

	var accessToken = "TOKEN=="
	var realmName = "master"
	var userID = "41dbf4a8-32a9-4000-8c17-edc854c31231"
	var username = "username"
	var code = "123456"

	// RecoveryCode
	{
		var kcCodeRep = kc.RecoveryCodeRepresentation{
			Code: &code,
		}

		mockKeycloakClient.EXPECT().CreateRecoveryCode(accessToken, realmName, userID).Return(kcCodeRep, nil).Times(1)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		ctx = context.WithValue(ctx, cs.CtContextRealm, realmName)
		ctx = context.WithValue(ctx, cs.CtContextUsername, username)

		mockEventDBModule.EXPECT().ReportEvent(ctx, "CREATE_RECOVERY_CODE", "back-office", gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)

		recoveryCode, err := managementComponent.CreateRecoveryCode(ctx, "master", userID)

		assert.Nil(t, err)
		assert.Equal(t, code, recoveryCode)
	}

	// RecoveryCode already exists
	{
		var err409 = kc.HTTPError{
			HTTPStatus: 409,
			Message:    "Conflict",
		}
		var kcCodeRep = kc.RecoveryCodeRepresentation{}

		mockKeycloakClient.EXPECT().CreateRecoveryCode(accessToken, realmName, userID).Return(kcCodeRep, err409).Times(1)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		ctx = context.WithValue(ctx, cs.CtContextRealm, realmName)
		ctx = context.WithValue(ctx, cs.CtContextUsername, username)

		mockLogger.EXPECT().Warn(gomock.Any(), "err", "409:Conflict")
		_, err := managementComponent.CreateRecoveryCode(ctx, "master", userID)

		assert.NotNil(t, err)
	}

	// Error
	{
		var kcCodeRep = kc.RecoveryCodeRepresentation{}
		mockKeycloakClient.EXPECT().CreateRecoveryCode(accessToken, realmName, userID).Return(kcCodeRep, fmt.Errorf("Error")).Times(1)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		ctx = context.WithValue(ctx, cs.CtContextRealm, realmName)
		ctx = context.WithValue(ctx, cs.CtContextUsername, username)

		mockLogger.EXPECT().Warn(gomock.Any(), "err", "Error")
		_, err := managementComponent.CreateRecoveryCode(ctx, "master", userID)

		assert.NotNil(t, err)
	}

}

func TestExecuteActionsEmail(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()
	var mockKeycloakClient = mock.NewKeycloakClient(mockCtrl)
	var mockUsersDetailsDBModule = mock.NewUsersDetailsDBModule(mockCtrl)
	var mockEventDBModule = mock.NewEventDBModule(mockCtrl)
	var mockConfigurationDBModule = mock.NewConfigurationDBModule(mockCtrl)
	var mockLogger = log.NewNopLogger()
	var allowedTrustIDGroups = []string{"grp1", "grp2"}

	var managementComponent = NewComponent(mockKeycloakClient, mockUsersDetailsDBModule, mockEventDBModule, mockConfigurationDBModule, allowedTrustIDGroups, mockLogger)

	var accessToken = "TOKEN=="
	var realmName = "master"
	var userID = "1245-7854-8963"
	var reqActions = []api.RequiredAction{initPasswordAction, "action1", "action2"}
	var actions = []string{initPasswordAction, "action1", "action2"}
	var key1 = "key1"
	var value1 = "value1"
	var key2 = "key2"
	var value2 = "value2"

	// Send email actions
	{

		mockKeycloakClient.EXPECT().ExecuteActionsEmail(accessToken, realmName, userID, actions, key1, value1, key2, value2).Return(nil).Times(1)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		mockEventDBModule.EXPECT().ReportEvent(ctx, "INIT_PASSWORD", "back-office", gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)
		mockEventDBModule.EXPECT().ReportEvent(ctx, "ACTION_EMAIL", "back-office", gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)

		err := managementComponent.ExecuteActionsEmail(ctx, "master", userID, reqActions, key1, value1, key2, value2)

		assert.Nil(t, err)
	}
	// Error
	{
		mockKeycloakClient.EXPECT().ExecuteActionsEmail(accessToken, realmName, userID, actions).Return(fmt.Errorf("Invalid input")).Times(1)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		mockEventDBModule.EXPECT().ReportEvent(ctx, "ACTION_EMAIL", "back-office", gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)
		mockEventDBModule.EXPECT().ReportEvent(ctx, "INIT_PASSWORD", "back-office", gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)

		err := managementComponent.ExecuteActionsEmail(ctx, "master", userID, reqActions)

		assert.NotNil(t, err)
	}
	// Send email actions, but not sms-password-set
	{

		var actions2 = []string{"action1", "action2"}
		var reqActions2 = []api.RequiredAction{"action1", "action2"}
		mockKeycloakClient.EXPECT().ExecuteActionsEmail(accessToken, realmName, userID, actions2, key1, value1, key2, value2).Return(nil).Times(1)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		mockEventDBModule.EXPECT().ReportEvent(ctx, "ACTION_EMAIL", "back-office", gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)

		err := managementComponent.ExecuteActionsEmail(ctx, "master", userID, reqActions2, key1, value1, key2, value2)

		assert.Nil(t, err)
	}
}

func TestSendNewEnrolmentCode(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()
	var mockKeycloakClient = mock.NewKeycloakClient(mockCtrl)
	var mockUsersDetailsDBModule = mock.NewUsersDetailsDBModule(mockCtrl)
	var mockEventDBModule = mock.NewEventDBModule(mockCtrl)
	var mockConfigurationDBModule = mock.NewConfigurationDBModule(mockCtrl)
	var mockLogger = mock.NewLogger(mockCtrl)
	var allowedTrustIDGroups = []string{"grp1", "grp2"}

	var managementComponent = NewComponent(mockKeycloakClient, mockUsersDetailsDBModule, mockEventDBModule, mockConfigurationDBModule, allowedTrustIDGroups, mockLogger)

	var accessToken = "TOKEN=="
	var realmName = "master"
	var userID = "1245-7854-8963"

	// Send new enrolment code
	{
		var code = "1234"
		mockKeycloakClient.EXPECT().SendNewEnrolmentCode(accessToken, realmName, userID).Return(kc.SmsCodeRepresentation{Code: &code}, nil).Times(1)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)

		mockEventDBModule.EXPECT().ReportEvent(ctx, "SMS_CHALLENGE", "back-office", gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)

		codeRes, err := managementComponent.SendNewEnrolmentCode(ctx, "master", userID)

		assert.Nil(t, err)
		assert.Equal(t, "1234", codeRes)
	}
	// Send new enrolment code but have error when storing the event in the DB
	{
		var code = "1234"
		mockKeycloakClient.EXPECT().SendNewEnrolmentCode(accessToken, realmName, userID).Return(kc.SmsCodeRepresentation{Code: &code}, nil).Times(1)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)

		mockEventDBModule.EXPECT().ReportEvent(ctx, "SMS_CHALLENGE", "back-office", database.CtEventRealmName, realmName, database.CtEventUserID, userID).Return(errors.New("error")).Times(1)
		m := map[string]interface{}{"event_name": "SMS_CHALLENGE", database.CtEventRealmName: realmName, database.CtEventUserID: userID}
		eventJSON, _ := json.Marshal(m)
		mockLogger.EXPECT().Error(gomock.Any(), "err", "error", "event", string(eventJSON))
		codeRes, err := managementComponent.SendNewEnrolmentCode(ctx, "master", userID)

		assert.Nil(t, err)
		assert.Equal(t, "1234", codeRes)
	}
	// Error
	{
		mockKeycloakClient.EXPECT().SendNewEnrolmentCode(accessToken, realmName, userID).Return(kc.SmsCodeRepresentation{}, fmt.Errorf("Invalid input")).Times(1)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		mockLogger.EXPECT().Warn(gomock.Any(), "err", "Invalid input")
		_, err := managementComponent.SendNewEnrolmentCode(ctx, "master", userID)

		assert.NotNil(t, err)
	}
}

func TestSendReminderEmail(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()
	var mockKeycloakClient = mock.NewKeycloakClient(mockCtrl)
	var mockUsersDetailsDBModule = mock.NewUsersDetailsDBModule(mockCtrl)
	var mockEventDBModule = mock.NewEventDBModule(mockCtrl)
	var mockConfigurationDBModule = mock.NewConfigurationDBModule(mockCtrl)
	var mockLogger = log.NewNopLogger()
	var allowedTrustIDGroups = []string{"grp1", "grp2"}

	var managementComponent = NewComponent(mockKeycloakClient, mockUsersDetailsDBModule, mockEventDBModule, mockConfigurationDBModule, allowedTrustIDGroups, mockLogger)

	var accessToken = "TOKEN=="
	var realmName = "master"
	var userID = "1245-7854-8963"

	var key1 = "key1"
	var value1 = "value1"
	var key2 = "key2"
	var value2 = "value2"
	var key3 = "key3"
	var value3 = "value3"

	// Send email
	{

		mockKeycloakClient.EXPECT().SendReminderEmail(accessToken, realmName, userID, key1, value1, key2, value2, key3, value3).Return(nil).Times(1)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)

		err := managementComponent.SendReminderEmail(ctx, "master", userID, key1, value1, key2, value2, key3, value3)

		assert.Nil(t, err)
	}

	// Error
	{
		mockKeycloakClient.EXPECT().SendReminderEmail(accessToken, realmName, userID).Return(fmt.Errorf("Invalid input")).Times(1)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)

		err := managementComponent.SendReminderEmail(ctx, "master", userID)

		assert.NotNil(t, err)
	}
}

func TestResetSmsCounter(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()
	var mockKeycloakClient = mock.NewKeycloakClient(mockCtrl)
	var mockUsersDetailsDBModule = mock.NewUsersDetailsDBModule(mockCtrl)
	var mockEventDBModule = mock.NewEventDBModule(mockCtrl)
	var mockConfigurationDBModule = mock.NewConfigurationDBModule(mockCtrl)
	var mockLogger = log.NewNopLogger()
	var allowedTrustIDGroups = []string{"grp1", "grp2"}

	var managementComponent = NewComponent(mockKeycloakClient, mockUsersDetailsDBModule, mockEventDBModule, mockConfigurationDBModule, allowedTrustIDGroups, mockLogger)

	var accessToken = "TOKEN=="
	var realmName = "master"
	var userID = "1245-7854-8963"
	var id = "1234-7454-4516"
	var username = "username"
	var email = "toto@elca.ch"
	var enabled = true
	var emailVerified = true
	var firstName = "Titi"
	var lastName = "Tutu"
	var phoneNumber = "+41789456"
	var phoneNumberVerified = true
	var label = "Label"
	var gender = "M"
	var birthDate = "01/01/1988"
	var createdTimestamp = time.Now().UTC().Unix()
	var attributes = make(kc.Attributes)
	attributes.SetString(constants.AttrbPhoneNumber, phoneNumber)
	attributes.SetString(constants.AttrbLabel, label)
	attributes.SetString(constants.AttrbGender, gender)
	attributes.SetString(constants.AttrbBirthDate, birthDate)
	attributes.SetBool(constants.AttrbPhoneNumberVerified, phoneNumberVerified)
	attributes.SetInt(constants.AttrbSmsSent, 5)
	attributes.SetInt(constants.AttrbSmsAttempts, 5)

	var kcUserRep = kc.UserRepresentation{
		ID:               &id,
		Username:         &username,
		Email:            &email,
		Enabled:          &enabled,
		EmailVerified:    &emailVerified,
		FirstName:        &firstName,
		LastName:         &lastName,
		Attributes:       &attributes,
		CreatedTimestamp: &createdTimestamp,
	}
	// Reset SMS counter
	{

		mockKeycloakClient.EXPECT().GetUser(accessToken, realmName, userID).Return(kcUserRep, nil).Times(1)
		mockKeycloakClient.EXPECT().UpdateUser(accessToken, realmName, userID, kcUserRep).Return(nil).Times(1)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)

		err := managementComponent.ResetSmsCounter(ctx, "master", userID)

		assert.Nil(t, err)
	}

	// Error at GetUser
	{

		mockKeycloakClient.EXPECT().GetUser(accessToken, realmName, userID).Return(kc.UserRepresentation{}, fmt.Errorf("error")).Times(1)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		err := managementComponent.ResetSmsCounter(ctx, "master", userID)

		assert.NotNil(t, err)
	}

	// Error at UpdateUser
	{

		mockKeycloakClient.EXPECT().GetUser(accessToken, realmName, userID).Return(kcUserRep, nil).Times(1)
		mockKeycloakClient.EXPECT().UpdateUser(accessToken, realmName, userID, kcUserRep).Return(fmt.Errorf("error")).Times(1)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		err := managementComponent.ResetSmsCounter(ctx, "master", userID)

		assert.NotNil(t, err)
	}
}

func TestGetCredentialsForUser(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()
	var mockKeycloakClient = mock.NewKeycloakClient(mockCtrl)
	var mockUsersDetailsDBModule = mock.NewUsersDetailsDBModule(mockCtrl)
	var mockEventDBModule = mock.NewEventDBModule(mockCtrl)
	var mockConfigurationDBModule = mock.NewConfigurationDBModule(mockCtrl)
	var mockLogger = log.NewNopLogger()
	var allowedTrustIDGroups = []string{"grp1", "grp2"}

	var managementComponent = NewComponent(mockKeycloakClient, mockUsersDetailsDBModule, mockEventDBModule, mockConfigurationDBModule, allowedTrustIDGroups, mockLogger)
	var accessToken = "TOKEN=="
	var realmReq = "master"
	var realmName = "otherRealm"
	var userID = "1245-7854-8963"

	// Get credentials for user
	{
		mockKeycloakClient.EXPECT().GetCredentials(accessToken, realmName, userID).Return([]kc.CredentialRepresentation{}, nil).Times(1)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		ctx = context.WithValue(ctx, cs.CtContextRealm, realmReq)

		_, err := managementComponent.GetCredentialsForUser(ctx, realmName, userID)

		assert.Nil(t, err)
	}
}

func TestDeleteCredentialsForUser(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()
	var mockKeycloakClient = mock.NewKeycloakClient(mockCtrl)
	var mockUsersDetailsDBModule = mock.NewUsersDetailsDBModule(mockCtrl)
	var mockEventDBModule = mock.NewEventDBModule(mockCtrl)
	var mockConfigurationDBModule = mock.NewConfigurationDBModule(mockCtrl)
	var mockLogger = mock.NewLogger(mockCtrl)
	var allowedTrustIDGroups = []string{"grp1", "grp2"}

	var managementComponent = NewComponent(mockKeycloakClient, mockUsersDetailsDBModule, mockEventDBModule, mockConfigurationDBModule, allowedTrustIDGroups, mockLogger)
	var accessToken = "TOKEN=="
	var realmReq = "master"
	var realmName = "master"
	var userID = "1245-7854-8963"
	var credential = "987-654-321"
	var pwdID = "51389847-08f4-4a0f-9f9c-694554e626f2"
	var pwd = "password"
	var credKcPwd = kc.CredentialRepresentation{
		ID:   &pwdID,
		Type: &pwd,
	}
	var otpID = "51389847-08f4-4a0f-9f9c-694554e626f3"
	var totp = "totp"
	var credKcOtp = kc.CredentialRepresentation{
		ID:   &otpID,
		Type: &totp,
	}
	var typeCred = "otp-push"

	t.Run("Delete credentials for user", func(t *testing.T) {
		mockKeycloakClient.EXPECT().GetCredentials(accessToken, realmName, userID).Return([]kc.CredentialRepresentation{
			kc.CredentialRepresentation{
				ID:   &credential,
				Type: &typeCred,
			},
		}, nil).Times(1)
		mockKeycloakClient.EXPECT().DeleteCredential(accessToken, realmName, userID, credential).Return(nil).Times(1)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		ctx = context.WithValue(ctx, cs.CtContextRealm, realmReq)

		mockEventDBModule.EXPECT().ReportEvent(ctx, "2ND_FACTOR_REMOVED", "back-office", database.CtEventRealmName, realmName, database.CtEventUserID, userID)

		err := managementComponent.DeleteCredentialsForUser(ctx, realmName, userID, credential)

		assert.Nil(t, err)
	})

	t.Run("Delete credentials for user - error at obtaining the list of credentials", func(t *testing.T) {
		mockKeycloakClient.EXPECT().GetCredentials(accessToken, realmName, userID).Return([]kc.CredentialRepresentation{}, errors.New("error")).Times(1)
		mockLogger.EXPECT().Warn(gomock.Any(), "msg", "Could not obtain list of credentials", "err", "error")

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		ctx = context.WithValue(ctx, cs.CtContextRealm, realmReq)

		err := managementComponent.DeleteCredentialsForUser(ctx, realmName, userID, credential)
		assert.NotNil(t, err)
	})

	t.Run("Delete credentials for user - try to delete credential of another user", func(t *testing.T) {
		mockKeycloakClient.EXPECT().GetCredentials(accessToken, realmName, userID).Return([]kc.CredentialRepresentation{}, nil).Times(1)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		ctx = context.WithValue(ctx, cs.CtContextRealm, realmReq)

		mockLogger.EXPECT().Warn(ctx, "msg", "Try to delete credential of another user", "credId", credential, "userId", userID)

		err := managementComponent.DeleteCredentialsForUser(ctx, realmName, userID, credential)
		assert.NotNil(t, err)
	})

	t.Run("Delete credentials for user - error at deleting the credential", func(t *testing.T) {
		mockKeycloakClient.EXPECT().GetCredentials(accessToken, realmName, userID).Return([]kc.CredentialRepresentation{
			kc.CredentialRepresentation{
				ID:   &credential,
				Type: &typeCred,
			},
		}, nil).Times(1)
		mockKeycloakClient.EXPECT().DeleteCredential(accessToken, realmName, userID, credential).Return(errors.New("error")).Times(1)
		mockLogger.EXPECT().Warn(gomock.Any(), "err", "error")
		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		ctx = context.WithValue(ctx, cs.CtContextRealm, realmReq)

		err := managementComponent.DeleteCredentialsForUser(ctx, realmName, userID, credential)
		assert.NotNil(t, err)
	})

	t.Run("Delete credentials for user", func(t *testing.T) {
		var credsKc []kc.CredentialRepresentation
		credsKc = append(credsKc, credKcPwd)
		credsKc = append(credsKc, credKcOtp)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		ctx = context.WithValue(ctx, cs.CtContextRealm, realmReq)

		mockKeycloakClient.EXPECT().GetCredentials(accessToken, realmName, userID).Return(credsKc, nil).Times(1)
		mockKeycloakClient.EXPECT().DeleteCredential(accessToken, realmName, userID, otpID).Return(nil).Times(1)
		mockEventDBModule.EXPECT().ReportEvent(ctx, "2ND_FACTOR_REMOVED", "back-office", gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)

		err := managementComponent.DeleteCredentialsForUser(ctx, realmName, userID, otpID)

		assert.Nil(t, err)
	})

	t.Run("Delete credentials for user - error at storing the event", func(t *testing.T) {
		var credsKc []kc.CredentialRepresentation
		credsKc = append(credsKc, credKcPwd)
		credsKc = append(credsKc, credKcOtp)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		ctx = context.WithValue(ctx, cs.CtContextRealm, realmReq)

		mockKeycloakClient.EXPECT().GetCredentials(accessToken, realmName, userID).Return(credsKc, nil).Times(1)
		mockKeycloakClient.EXPECT().DeleteCredential(accessToken, realmName, userID, otpID).Return(nil).Times(1)
		mockEventDBModule.EXPECT().ReportEvent(ctx, "2ND_FACTOR_REMOVED", "back-office", gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("error")).Times(1)
		m := map[string]interface{}{"event_name": "2ND_FACTOR_REMOVED", database.CtEventRealmName: realmName, database.CtEventUserID: userID}
		eventJSON, _ := json.Marshal(m)
		mockLogger.EXPECT().Error(gomock.Any(), "err", "error", "event", string(eventJSON))
		err := managementComponent.DeleteCredentialsForUser(ctx, realmName, userID, otpID)

		assert.Nil(t, err)
	})
}

func TestUnlockCredentialForUser(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()

	var mockKeycloakClient = mock.NewKeycloakClient(mockCtrl)
	var mockEventDBModule = mock.NewEventDBModule(mockCtrl)
	var mockConfigurationDBModule = mock.NewConfigurationDBModule(mockCtrl)

	var managementComponent = NewComponent(mockKeycloakClient, nil, mockEventDBModule, mockConfigurationDBModule, nil, log.NewNopLogger())
	var accessToken = "TOKEN=="
	var realmName = "master"
	var userID = "1245-7854-8963"
	var credentialID = "987-654-321"
	var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)

	t.Run("Detect credential type-Keycloak call fails", func(t *testing.T) {
		var kcErr = errors.New("keycloak error")
		mockKeycloakClient.EXPECT().GetCredentials(accessToken, realmName, userID).Return(nil, kcErr)

		var err = managementComponent.ResetCredentialFailuresForUser(ctx, realmName, userID, credentialID)
		assert.Equal(t, kcErr, err)
	})

	t.Run("Detect credential type-Credential not found", func(t *testing.T) {
		mockKeycloakClient.EXPECT().GetCredentials(accessToken, realmName, userID).Return(nil, nil)

		var err = managementComponent.ResetCredentialFailuresForUser(ctx, realmName, userID, credentialID)
		assert.NotNil(t, err)
	})

	var foundCredType = "ctpapercard"
	var credentials = []kc.CredentialRepresentation{
		kc.CredentialRepresentation{ID: &credentialID, Type: &foundCredType},
	}
	mockKeycloakClient.EXPECT().GetCredentials(accessToken, realmName, userID).Return(credentials, nil).AnyTimes()

	t.Run("Detect credential type-Credential found", func(t *testing.T) {
		mockKeycloakClient.EXPECT().ResetPapercardFailures(accessToken, realmName, userID, credentialID).Return(nil)

		var err = managementComponent.ResetCredentialFailuresForUser(ctx, realmName, userID, credentialID)
		assert.Nil(t, err)
	})

	t.Run("Can't unlock paper card", func(t *testing.T) {
		var unlockErr = errors.New("unlock error")
		mockKeycloakClient.EXPECT().ResetPapercardFailures(accessToken, realmName, userID, credentialID).Return(unlockErr)

		var err = managementComponent.ResetCredentialFailuresForUser(ctx, realmName, userID, credentialID)
		assert.Equal(t, unlockErr, err)
	})
}

func TestClearUserLoginFailures(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()
	var mockKeycloakClient = mock.NewKeycloakClient(mockCtrl)
	var mockUsersDetailsDBModule = mock.NewUsersDetailsDBModule(mockCtrl)
	var mockEventDBModule = mock.NewEventDBModule(mockCtrl)
	var mockConfigurationDBModule = mock.NewConfigurationDBModule(mockCtrl)
	var logger = log.NewNopLogger()

	var accessToken = "TOKEN=="
	var realm = "master"
	var userID = "1245-7854-8963"
	var allowedTrustIDGroups = []string{"grp1", "grp2"}
	var ctx = context.WithValue(context.TODO(), cs.CtContextAccessToken, accessToken)
	var component = NewComponent(mockKeycloakClient, mockUsersDetailsDBModule, mockEventDBModule, mockConfigurationDBModule, allowedTrustIDGroups, logger)

	t.Run("Error occured", func(t *testing.T) {
		var expectedError = errors.New("kc error")
		mockKeycloakClient.EXPECT().ClearUserLoginFailures(accessToken, realm, userID).Return(expectedError)
		var err = component.ClearUserLoginFailures(ctx, realm, userID)
		assert.Equal(t, expectedError, err)
	})
	t.Run("Success", func(t *testing.T) {
		mockKeycloakClient.EXPECT().ClearUserLoginFailures(accessToken, realm, userID).Return(nil)
		mockEventDBModule.EXPECT().ReportEvent(ctx, "LOGIN_FAILURE_CLEARED", gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any())
		var err = component.ClearUserLoginFailures(ctx, realm, userID)
		assert.Nil(t, err)
	})
}

func TestGetAttackDetectionStatus(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()
	var mockKeycloakClient = mock.NewKeycloakClient(mockCtrl)
	var mockUsersDetailsDBModule = mock.NewUsersDetailsDBModule(mockCtrl)
	var mockEventDBModule = mock.NewEventDBModule(mockCtrl)
	var mockConfigurationDBModule = mock.NewConfigurationDBModule(mockCtrl)
	var logger = log.NewNopLogger()

	var accessToken = "TOKEN=="
	var realm = "master"
	var userID = "1245-7854-8963"
	var allowedTrustIDGroups = []string{"grp1", "grp2"}
	var ctx = context.WithValue(context.TODO(), cs.CtContextAccessToken, accessToken)
	var component = NewComponent(mockKeycloakClient, mockUsersDetailsDBModule, mockEventDBModule, mockConfigurationDBModule, allowedTrustIDGroups, logger)
	var kcResult = map[string]interface{}{}

	t.Run("Error occured", func(t *testing.T) {
		var expectedError = errors.New("kc error")
		mockKeycloakClient.EXPECT().GetAttackDetectionStatus(accessToken, realm, userID).Return(kcResult, expectedError)
		var _, err = component.GetAttackDetectionStatus(ctx, realm, userID)
		assert.Equal(t, expectedError, err)
	})
	t.Run("Success", func(t *testing.T) {
		var expectedFailures int64 = 57
		kcResult["numFailures"] = expectedFailures
		mockKeycloakClient.EXPECT().GetAttackDetectionStatus(accessToken, realm, userID).Return(kcResult, nil)
		var res, err = component.GetAttackDetectionStatus(ctx, realm, userID)
		assert.Nil(t, err)
		assert.Equal(t, expectedFailures, *res.NumFailures)
	})
}

func TestGetRoles(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()
	var mockKeycloakClient = mock.NewKeycloakClient(mockCtrl)
	var mockUsersDetailsDBModule = mock.NewUsersDetailsDBModule(mockCtrl)
	var mockEventDBModule = mock.NewEventDBModule(mockCtrl)
	var mockConfigurationDBModule = mock.NewConfigurationDBModule(mockCtrl)
	var mockLogger = log.NewNopLogger()
	var allowedTrustIDGroups = []string{"grp1", "grp2"}

	var managementComponent = NewComponent(mockKeycloakClient, mockUsersDetailsDBModule, mockEventDBModule, mockConfigurationDBModule, allowedTrustIDGroups, mockLogger)

	var accessToken = "TOKEN=="
	var realmName = "master"

	// Get roles with succces
	{
		var id = "1234-7454-4516"
		var composite = false
		var containerID = "containerId"
		var description = "description role"
		var clientRole = false
		var name = "name"

		var kcRoleRep = kc.RoleRepresentation{
			ID:          &id,
			Name:        &name,
			ClientRole:  &clientRole,
			Composite:   &composite,
			ContainerID: &containerID,
			Description: &description,
		}

		var kcRolesRep []kc.RoleRepresentation
		kcRolesRep = append(kcRolesRep, kcRoleRep)

		mockKeycloakClient.EXPECT().GetRoles(accessToken, realmName).Return(kcRolesRep, nil).Times(1)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)

		apiRolesRep, err := managementComponent.GetRoles(ctx, "master")

		var apiRoleRep = apiRolesRep[0]
		assert.Nil(t, err)
		assert.Equal(t, id, *apiRoleRep.ID)
		assert.Equal(t, name, *apiRoleRep.Name)
		assert.Equal(t, clientRole, *apiRoleRep.ClientRole)
		assert.Equal(t, composite, *apiRoleRep.Composite)
		assert.Equal(t, containerID, *apiRoleRep.ContainerID)
		assert.Equal(t, description, *apiRoleRep.Description)
	}

	//Error
	{
		mockKeycloakClient.EXPECT().GetRoles(accessToken, realmName).Return([]kc.RoleRepresentation{}, fmt.Errorf("Unexpected error")).Times(1)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)

		_, err := managementComponent.GetRoles(ctx, "master")

		assert.NotNil(t, err)
	}
}

func TestGetRole(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()
	var mockKeycloakClient = mock.NewKeycloakClient(mockCtrl)
	var mockUsersDetailsDBModule = mock.NewUsersDetailsDBModule(mockCtrl)
	var mockEventDBModule = mock.NewEventDBModule(mockCtrl)
	var mockConfigurationDBModule = mock.NewConfigurationDBModule(mockCtrl)
	var mockLogger = log.NewNopLogger()
	var allowedTrustIDGroups = []string{"grp1", "grp2"}

	var managementComponent = NewComponent(mockKeycloakClient, mockUsersDetailsDBModule, mockEventDBModule, mockConfigurationDBModule, allowedTrustIDGroups, mockLogger)

	var accessToken = "TOKEN=="
	var realmName = "master"

	// Get roles with succces
	{
		var id = "1234-7454-4516"
		var composite = false
		var containerID = "containerId"
		var description = "description role"
		var clientRole = false
		var name = "name"

		var kcRoleRep = kc.RoleRepresentation{
			ID:          &id,
			Name:        &name,
			ClientRole:  &clientRole,
			Composite:   &composite,
			ContainerID: &containerID,
			Description: &description,
		}

		mockKeycloakClient.EXPECT().GetRole(accessToken, realmName, id).Return(kcRoleRep, nil).Times(1)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)

		apiRoleRep, err := managementComponent.GetRole(ctx, "master", id)

		assert.Nil(t, err)
		assert.Equal(t, id, *apiRoleRep.ID)
		assert.Equal(t, name, *apiRoleRep.Name)
		assert.Equal(t, clientRole, *apiRoleRep.ClientRole)
		assert.Equal(t, composite, *apiRoleRep.Composite)
		assert.Equal(t, containerID, *apiRoleRep.ContainerID)
		assert.Equal(t, description, *apiRoleRep.Description)
	}

	//Error
	{
		var id = "1234-7454-4516"
		mockKeycloakClient.EXPECT().GetRole(accessToken, realmName, id).Return(kc.RoleRepresentation{}, fmt.Errorf("Unexpected error")).Times(1)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)

		_, err := managementComponent.GetRole(ctx, "master", id)

		assert.NotNil(t, err)
	}
}

func TestGetGroups(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()
	var mockKeycloakClient = mock.NewKeycloakClient(mockCtrl)
	var mockUsersDetailsDBModule = mock.NewUsersDetailsDBModule(mockCtrl)
	var mockEventDBModule = mock.NewEventDBModule(mockCtrl)
	var mockConfigurationDBModule = mock.NewConfigurationDBModule(mockCtrl)
	var mockLogger = log.NewNopLogger()
	var allowedTrustIDGroups = []string{"grp1", "grp2"}

	var managementComponent = NewComponent(mockKeycloakClient, mockUsersDetailsDBModule, mockEventDBModule, mockConfigurationDBModule, allowedTrustIDGroups, mockLogger)

	var accessToken = "TOKEN=="
	var realmName = "master"

	// Get groups with succces, non empty result
	{
		var id = "1234-7454-4516"
		var path = "path_group"
		var name = "group1"
		var realmRoles = []string{"role1"}

		var kcGroupRep = kc.GroupRepresentation{
			ID:         &id,
			Name:       &name,
			Path:       &path,
			RealmRoles: &realmRoles,
		}

		var kcGroupsRep []kc.GroupRepresentation
		kcGroupsRep = append(kcGroupsRep, kcGroupRep)

		mockKeycloakClient.EXPECT().GetGroups(accessToken, realmName).Return(kcGroupsRep, nil).Times(1)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)

		apiGroupsRep, err := managementComponent.GetGroups(ctx, "master")

		var apiGroupRep = apiGroupsRep[0]
		assert.Nil(t, err)
		assert.Equal(t, id, *apiGroupRep.ID)
		assert.Equal(t, name, *apiGroupRep.Name)
	}

	// Get groups with success, empty result
	{
		var kcGroupsRep []kc.GroupRepresentation
		mockKeycloakClient.EXPECT().GetGroups(accessToken, realmName).Return(kcGroupsRep, nil).Times(1)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		apiGroupsRep, err := managementComponent.GetGroups(ctx, "master")

		assert.Nil(t, err)
		assert.NotNil(t, apiGroupsRep)
		assert.Equal(t, 0, len(apiGroupsRep))
	}

	//Error
	{
		mockKeycloakClient.EXPECT().GetGroups(accessToken, realmName).Return([]kc.GroupRepresentation{}, fmt.Errorf("Unexpected error")).Times(1)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)

		_, err := managementComponent.GetGroups(ctx, "master")

		assert.NotNil(t, err)
	}
}

func TestCreateGroup(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()
	var mockKeycloakClient = mock.NewKeycloakClient(mockCtrl)
	var mockUsersDetailsDBModule = mock.NewUsersDetailsDBModule(mockCtrl)
	var mockEventDBModule = mock.NewEventDBModule(mockCtrl)
	var mockConfigurationDBModule = mock.NewConfigurationDBModule(mockCtrl)
	var mockLogger = mock.NewLogger(mockCtrl)
	var allowedTrustIDGroups = []string{"grp1", "grp2"}

	var managementComponent = NewComponent(mockKeycloakClient, mockUsersDetailsDBModule, mockEventDBModule, mockConfigurationDBModule, allowedTrustIDGroups, mockLogger)

	var accessToken = "TOKEN=="
	var username = "username"
	var name = "test"
	var realmName = "master"
	var targetRealmName = "DEP"
	var groupID = "41dbf4a8-32a9-4000-8c17-edc854c31231"
	var locationURL = "http://toto.com/realms/" + groupID

	// Create
	{
		var kcGroupRep = kc.GroupRepresentation{
			Name: &name,
		}

		mockKeycloakClient.EXPECT().CreateGroup(accessToken, targetRealmName, kcGroupRep).Return(locationURL, nil).Times(1)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		ctx = context.WithValue(ctx, cs.CtContextRealm, realmName)
		ctx = context.WithValue(ctx, cs.CtContextUsername, username)

		mockEventDBModule.EXPECT().ReportEvent(ctx, "API_GROUP_CREATION", "back-office", gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)

		var groupRep = api.GroupRepresentation{
			Name: &name,
		}

		location, err := managementComponent.CreateGroup(ctx, targetRealmName, groupRep)

		assert.Nil(t, err)
		assert.Equal(t, locationURL, location)
	}

	//Create with having error when storing the event
	{
		var kcGroupRep = kc.GroupRepresentation{
			Name: &name,
		}

		mockKeycloakClient.EXPECT().CreateGroup(accessToken, targetRealmName, kcGroupRep).Return(locationURL, nil).Times(1)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		ctx = context.WithValue(ctx, cs.CtContextRealm, realmName)
		ctx = context.WithValue(ctx, cs.CtContextUsername, username)

		mockEventDBModule.EXPECT().ReportEvent(ctx, "API_GROUP_CREATION", "back-office", database.CtEventRealmName, targetRealmName, database.CtEventGroupID, groupID, database.CtEventGroupName, name).Return(errors.New("error")).Times(1)
		m := map[string]interface{}{"event_name": "API_GROUP_CREATION", database.CtEventRealmName: targetRealmName, database.CtEventGroupID: groupID, database.CtEventGroupName: name}
		eventJSON, _ := json.Marshal(m)
		mockLogger.EXPECT().Error(ctx, "err", "error", "event", string(eventJSON))

		var groupRep = api.GroupRepresentation{
			Name: &name,
		}

		location, err := managementComponent.CreateGroup(ctx, targetRealmName, groupRep)

		assert.Nil(t, err)
		assert.Equal(t, locationURL, location)

	}

	// Error from KC client
	{
		var kcGroupRep = kc.GroupRepresentation{}

		mockKeycloakClient.EXPECT().CreateGroup(accessToken, targetRealmName, kcGroupRep).Return("", fmt.Errorf("Invalid input")).Times(1)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		ctx = context.WithValue(ctx, cs.CtContextRealm, realmName)

		var groupRep = api.GroupRepresentation{}
		mockLogger.EXPECT().Warn(ctx, "err", "Invalid input")

		location, err := managementComponent.CreateGroup(ctx, targetRealmName, groupRep)

		assert.NotNil(t, err)
		assert.Equal(t, "", location)
	}
}

func TestDeleteGroup(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()
	var mockKeycloakClient = mock.NewKeycloakClient(mockCtrl)
	var mockUsersDetailsDBModule = mock.NewUsersDetailsDBModule(mockCtrl)
	var mockEventDBModule = mock.NewEventDBModule(mockCtrl)
	var mockConfigurationDBModule = mock.NewConfigurationDBModule(mockCtrl)
	var mockLogger = mock.NewLogger(mockCtrl)
	var allowedTrustIDGroups = []string{"grp1", "grp2"}

	var managementComponent = NewComponent(mockKeycloakClient, mockUsersDetailsDBModule, mockEventDBModule, mockConfigurationDBModule, allowedTrustIDGroups, mockLogger)

	var accessToken = "TOKEN=="
	var groupID = "41dbf4a8-32a9-4000-8c17-edc854c31231"
	var groupName = "groupName"
	var targetRealmName = "DEP"
	var realmName = "master"
	var username = "username"

	var group = kc.GroupRepresentation{
		ID:   &groupID,
		Name: &groupName,
	}

	// Delete group with success
	{
		mockKeycloakClient.EXPECT().DeleteGroup(accessToken, targetRealmName, groupID).Return(nil).Times(1)
		mockKeycloakClient.EXPECT().GetGroup(accessToken, targetRealmName, groupID).Return(group, nil).Times(1)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		ctx = context.WithValue(ctx, cs.CtContextRealm, realmName)
		ctx = context.WithValue(ctx, cs.CtContextUsername, username)

		mockConfigurationDBModule.EXPECT().DeleteAllAuthorizationsWithGroup(ctx, targetRealmName, groupName).Return(nil).Times(1)
		mockEventDBModule.EXPECT().ReportEvent(ctx, "API_GROUP_DELETION", "back-office", gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)

		err := managementComponent.DeleteGroup(ctx, targetRealmName, groupID)

		assert.Nil(t, err)
	}

	// Delete group with success but having an error when storing the event in the DB
	{
		mockKeycloakClient.EXPECT().DeleteGroup(accessToken, targetRealmName, groupID).Return(nil).Times(1)
		mockKeycloakClient.EXPECT().GetGroup(accessToken, targetRealmName, groupID).Return(group, nil).Times(1)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		ctx = context.WithValue(ctx, cs.CtContextRealm, realmName)
		ctx = context.WithValue(ctx, cs.CtContextUsername, username)

		mockConfigurationDBModule.EXPECT().DeleteAllAuthorizationsWithGroup(ctx, targetRealmName, groupName).Return(nil).Times(1)
		mockEventDBModule.EXPECT().ReportEvent(ctx, "API_GROUP_DELETION", "back-office", database.CtEventRealmName, targetRealmName, database.CtEventGroupName, groupName).Return(errors.New("error")).Times(1)
		m := map[string]interface{}{"event_name": "API_GROUP_DELETION", database.CtEventRealmName: targetRealmName, database.CtEventGroupName: groupName}
		eventJSON, _ := json.Marshal(m)
		mockLogger.EXPECT().Error(ctx, "err", "error", "event", string(eventJSON))
		err := managementComponent.DeleteGroup(ctx, targetRealmName, groupID)

		assert.Nil(t, err)
	}

	// Error with DB
	{
		mockKeycloakClient.EXPECT().DeleteGroup(accessToken, targetRealmName, groupID).Return(nil).Times(1)
		mockKeycloakClient.EXPECT().GetGroup(accessToken, targetRealmName, groupID).Return(group, nil).Times(1)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		ctx = context.WithValue(ctx, cs.CtContextRealm, realmName)
		ctx = context.WithValue(ctx, cs.CtContextUsername, username)

		mockConfigurationDBModule.EXPECT().DeleteAllAuthorizationsWithGroup(ctx, targetRealmName, groupName).Return(fmt.Errorf("Error")).Times(1)
		mockLogger.EXPECT().Warn(ctx, "err", "Error")

		err := managementComponent.DeleteGroup(ctx, targetRealmName, groupID)

		assert.NotNil(t, err)
	}

	// Error from KC client
	{
		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		mockLogger.EXPECT().Warn(ctx, "err", "Error")
		mockKeycloakClient.EXPECT().GetGroup(accessToken, targetRealmName, groupID).Return(kc.GroupRepresentation{}, errors.New("Error")).Times(1)

		err := managementComponent.DeleteGroup(ctx, targetRealmName, groupID)
		assert.NotNil(t, err)

		mockKeycloakClient.EXPECT().GetGroup(accessToken, targetRealmName, groupID).Return(group, nil).Times(1)
		mockKeycloakClient.EXPECT().DeleteGroup(accessToken, targetRealmName, groupID).Return(fmt.Errorf("Invalid input")).Times(1)
		mockLogger.EXPECT().Warn(ctx, "err", "Invalid input")

		err = managementComponent.DeleteGroup(ctx, targetRealmName, groupID)
		assert.NotNil(t, err)
	}
}

func TestGetAuthorizations(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()
	var mockKeycloakClient = mock.NewKeycloakClient(mockCtrl)
	var mockUsersDetailsDBModule = mock.NewUsersDetailsDBModule(mockCtrl)
	var mockEventDBModule = mock.NewEventDBModule(mockCtrl)
	var mockConfigurationDBModule = mock.NewConfigurationDBModule(mockCtrl)
	var mockLogger = mock.NewLogger(mockCtrl)
	var allowedTrustIDGroups = []string{"grp1", "grp2"}

	var managementComponent = NewComponent(mockKeycloakClient, mockUsersDetailsDBModule, mockEventDBModule, mockConfigurationDBModule, allowedTrustIDGroups, mockLogger)

	var accessToken = "TOKEN=="
	var realmName = "master"
	var targetRealmname = "DEP"
	var groupID = "41dbf4a8-32a9-4000-8c17-edc854c31231"
	var groupName = "groupName"
	var username = "username"
	var action = "action"

	var group = kc.GroupRepresentation{
		ID:   &groupID,
		Name: &groupName,
	}

	// Get authorizations with succces
	{
		var configurationAuthz = []configuration.Authorization{
			configuration.Authorization{
				RealmID:   &realmName,
				GroupName: &groupName,
				Action:    &action,
			},
		}

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		ctx = context.WithValue(ctx, cs.CtContextRealm, realmName)
		ctx = context.WithValue(ctx, cs.CtContextUsername, username)

		mockConfigurationDBModule.EXPECT().GetAuthorizations(ctx, targetRealmname, groupName).Return(configurationAuthz, nil).Times(1)
		mockKeycloakClient.EXPECT().GetGroup(accessToken, targetRealmname, groupID).Return(group, nil).Times(1)

		apiAuthorizationRep, err := managementComponent.GetAuthorizations(ctx, targetRealmname, groupID)

		var matrix = map[string]map[string]map[string]struct{}{
			"action": {},
		}

		var expectedAPIAuthorization = api.AuthorizationsRepresentation{
			Matrix: &matrix,
		}

		assert.Nil(t, err)
		assert.Equal(t, expectedAPIAuthorization, apiAuthorizationRep)
	}

	//Error when retrieving authorizations from DB
	{
		mockKeycloakClient.EXPECT().GetGroup(accessToken, targetRealmname, groupID).Return(group, nil).Times(1)
		mockConfigurationDBModule.EXPECT().GetAuthorizations(gomock.Any(), targetRealmname, groupName).Return([]configuration.Authorization{}, fmt.Errorf("Error")).Times(1)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		mockLogger.EXPECT().Warn(ctx, "err", "Error")

		_, err := managementComponent.GetAuthorizations(ctx, targetRealmname, groupID)

		assert.NotNil(t, err)
	}

	//Error with KC
	{
		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)

		mockKeycloakClient.EXPECT().GetGroup(accessToken, targetRealmname, groupID).Return(kc.GroupRepresentation{}, errors.New("Error")).Times(1)
		mockLogger.EXPECT().Warn(ctx, "err", "Error")
		_, err := managementComponent.GetAuthorizations(ctx, targetRealmname, groupID)
		assert.NotNil(t, err)
	}
}

func TestUpdateAuthorizations(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()
	var mockKeycloakClient = mock.NewKeycloakClient(mockCtrl)
	var mockUsersDetailsDBModule = mock.NewUsersDetailsDBModule(mockCtrl)
	var mockEventDBModule = mock.NewEventDBModule(mockCtrl)
	var mockTransaction = mock.NewTransaction(mockCtrl)
	var mockConfigurationDBModule = mock.NewConfigurationDBModule(mockCtrl)
	var mockLogger = mock.NewLogger(mockCtrl)
	var allowedTrustIDGroups = []string{"grp1", "grp2"}

	var managementComponent = NewComponent(mockKeycloakClient, mockUsersDetailsDBModule, mockEventDBModule, mockConfigurationDBModule, allowedTrustIDGroups, mockLogger)

	var accessToken = "TOKEN=="
	var currentRealmName = "master"
	var realmName = "customer1"
	var targetRealmName = "DEP"
	var ID = "00000-32a9-4000-8c17-edc854c31231"
	var groupID = "41dbf4a8-32a9-4000-8c17-edc854c31231"
	var groupName = "groupName"
	var ID1 = "111111-32a9-4000-8c17-edc854c31231"
	var ID2 = "222222-32a9-4000-8c17-edc854c31231"
	var ID3 = "333333-32a9-4000-8c17-edc854c31231"
	var ID4 = "444444-32a9-4000-8c17-edc854c31231"
	var ID5 = "555555-32a9-4000-8c17-edc854c31231"
	var username = "username"
	var clientID = "realm-management"
	var clientID2 = "backofficeid"

	var realm = kc.RealmRepresentation{
		ID:    &targetRealmName,
		Realm: &targetRealmName,
	}
	var realms = []kc.RealmRepresentation{realm}

	var group = kc.GroupRepresentation{
		ID:   &groupID,
		Name: &groupName,
	}
	var groups = []kc.GroupRepresentation{group}

	var client = kc.ClientRepresentation{
		ID:       &ID,
		ClientID: &clientID,
	}
	var client2 = kc.ClientRepresentation{
		ID:       &ID2,
		ClientID: &clientID2,
	}
	var clients = []kc.ClientRepresentation{client, client2}

	var roleName = []string{"manage-users", "view-clients", "view-realm", "view-users", "other"}
	var roleManageUser = kc.RoleRepresentation{
		ID:   &ID1,
		Name: &roleName[0],
	}
	var roleViewClients = kc.RoleRepresentation{
		ID:   &ID2,
		Name: &roleName[1],
	}
	var roleViewRealm = kc.RoleRepresentation{
		ID:   &ID3,
		Name: &roleName[2],
	}
	var roleViewUsers = kc.RoleRepresentation{
		ID:   &ID4,
		Name: &roleName[3],
	}
	var roleOther = kc.RoleRepresentation{
		ID:   &ID5,
		Name: &roleName[4],
	}

	t.Run("Update authorizations with succces (MGMT_action so KC roles needed)", func(t *testing.T) {
		var action = "MGMT_action"
		var matrix = map[string]map[string]map[string]struct{}{
			action: {},
		}

		var apiAuthorizations = api.AuthorizationsRepresentation{
			Matrix: &matrix,
		}

		var rolesAvailable = []kc.RoleRepresentation{
			roleManageUser,
			roleViewClients,
			roleViewRealm,
			roleViewUsers,
		}
		var rolesCurrent = []kc.RoleRepresentation{
			roleOther,
		}

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		ctx = context.WithValue(ctx, cs.CtContextRealm, currentRealmName)
		ctx = context.WithValue(ctx, cs.CtContextUsername, username)

		mockKeycloakClient.EXPECT().GetGroup(accessToken, targetRealmName, groupID).Return(group, nil).Times(1)
		mockKeycloakClient.EXPECT().GetRealms(accessToken).Return(realms, nil).Times(1)
		mockKeycloakClient.EXPECT().GetGroups(accessToken, targetRealmName).Return(groups, nil).Times(1)
		mockKeycloakClient.EXPECT().GetClients(accessToken, targetRealmName).Return(clients, nil).Times(1)
		mockKeycloakClient.EXPECT().GetAvailableGroupClientRoles(accessToken, targetRealmName, groupID, ID).Return(rolesAvailable, nil).Times(1)
		mockKeycloakClient.EXPECT().GetGroupClientRoles(accessToken, targetRealmName, groupID, ID).Return(rolesCurrent, nil).Times(1)
		mockKeycloakClient.EXPECT().AssignClientRole(accessToken, targetRealmName, groupID, ID, gomock.Any()).Return(nil).Times(1)

		mockConfigurationDBModule.EXPECT().NewTransaction(ctx).Return(mockTransaction, nil).Times(1)
		mockConfigurationDBModule.EXPECT().DeleteAuthorizations(ctx, targetRealmName, groupName).Return(nil).Times(1)
		mockConfigurationDBModule.EXPECT().CreateAuthorization(ctx, gomock.Any()).Return(nil).Times(1)
		mockTransaction.EXPECT().Close().Times(1)
		mockTransaction.EXPECT().Commit().Times(1)

		mockEventDBModule.EXPECT().ReportEvent(ctx, "API_AUTHORIZATIONS_UPDATE", "back-office", gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)

		err := managementComponent.UpdateAuthorizations(ctx, targetRealmName, groupID, apiAuthorizations)
		assert.Nil(t, err)
	})

	t.Run("Update authorizations with succces (action so no KC roles needed)", func(t *testing.T) {
		var action = "action"
		var matrix = map[string]map[string]map[string]struct{}{
			action: {},
		}

		var apiAuthorizations = api.AuthorizationsRepresentation{
			Matrix: &matrix,
		}

		var rolesCurrent = []kc.RoleRepresentation{
			roleManageUser,
			roleViewClients,
			roleViewRealm,
			roleViewUsers,
			roleOther,
		}
		var rolesAvailable = []kc.RoleRepresentation{}

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		ctx = context.WithValue(ctx, cs.CtContextRealm, realmName)
		ctx = context.WithValue(ctx, cs.CtContextUsername, username)

		mockKeycloakClient.EXPECT().GetGroup(accessToken, targetRealmName, groupID).Return(group, nil).Times(1)
		mockKeycloakClient.EXPECT().GetRealms(accessToken).Return(realms, nil).Times(1)
		mockKeycloakClient.EXPECT().GetGroups(accessToken, targetRealmName).Return(groups, nil).Times(1)
		mockKeycloakClient.EXPECT().GetClients(accessToken, targetRealmName).Return(clients, nil).Times(1)
		mockKeycloakClient.EXPECT().GetAvailableGroupClientRoles(accessToken, targetRealmName, groupID, ID).Return(rolesAvailable, nil).Times(1)
		mockKeycloakClient.EXPECT().GetGroupClientRoles(accessToken, targetRealmName, groupID, ID).Return(rolesCurrent, nil).Times(1)
		mockKeycloakClient.EXPECT().RemoveClientRole(accessToken, targetRealmName, groupID, ID, gomock.Any()).Return(nil).Times(1)

		mockConfigurationDBModule.EXPECT().NewTransaction(ctx).Return(mockTransaction, nil).Times(1)
		mockConfigurationDBModule.EXPECT().DeleteAuthorizations(ctx, targetRealmName, groupName).Return(nil).Times(1)
		mockConfigurationDBModule.EXPECT().CreateAuthorization(ctx, gomock.Any()).Return(nil).Times(1)
		mockTransaction.EXPECT().Close().Times(1)
		mockTransaction.EXPECT().Commit().Times(1)

		mockEventDBModule.EXPECT().ReportEvent(ctx, "API_AUTHORIZATIONS_UPDATE", "back-office", gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)

		err := managementComponent.UpdateAuthorizations(ctx, targetRealmName, groupID, apiAuthorizations)

		assert.Nil(t, err)
	})

	t.Run("Errors", func(t *testing.T) {
		var action = "MGMT_action"
		var matrix = map[string]map[string]map[string]struct{}{
			action: {},
		}

		var apiAuthorizations = api.AuthorizationsRepresentation{
			Matrix: &matrix,
		}

		var rolesAvailable = []kc.RoleRepresentation{
			roleManageUser,
			roleViewClients,
			roleViewRealm,
			roleViewUsers,
		}
		var rolesCurrent = []kc.RoleRepresentation{
			roleOther,
		}

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		ctx = context.WithValue(ctx, cs.CtContextRealm, realmName)
		ctx = context.WithValue(ctx, cs.CtContextUsername, username)

		mockKeycloakClient.EXPECT().GetGroup(accessToken, targetRealmName, groupID).Return(kc.GroupRepresentation{}, errors.New("Error")).Times(1)
		mockLogger.EXPECT().Warn(ctx, "err", "Error").Times(1)
		err := managementComponent.UpdateAuthorizations(ctx, targetRealmName, groupID, apiAuthorizations)
		assert.NotNil(t, err)

		mockKeycloakClient.EXPECT().GetGroup(accessToken, targetRealmName, groupID).Return(group, nil).Times(1)
		mockKeycloakClient.EXPECT().GetRealms(accessToken).Return([]kc.RealmRepresentation{}, fmt.Errorf("Unexpected error")).Times(1)
		mockLogger.EXPECT().Warn(ctx, "err", "Unexpected error").Times(1)
		err = managementComponent.UpdateAuthorizations(ctx, targetRealmName, groupID, apiAuthorizations)
		assert.NotNil(t, err)

		mockKeycloakClient.EXPECT().GetGroup(accessToken, targetRealmName, groupID).Return(group, nil).Times(1)
		mockKeycloakClient.EXPECT().GetRealms(accessToken).Return(realms, nil).Times(1)
		mockKeycloakClient.EXPECT().GetGroups(accessToken, targetRealmName).Return([]kc.GroupRepresentation{}, fmt.Errorf("Unexpected error")).Times(1)
		mockLogger.EXPECT().Warn(ctx, "err", "Unexpected error").Times(1)
		err = managementComponent.UpdateAuthorizations(ctx, targetRealmName, groupID, apiAuthorizations)
		assert.NotNil(t, err)

		mockKeycloakClient.EXPECT().GetGroup(accessToken, targetRealmName, groupID).Return(group, nil).Times(1)
		mockKeycloakClient.EXPECT().GetRealms(accessToken).Return(realms, nil).Times(1)
		mockKeycloakClient.EXPECT().GetGroups(accessToken, targetRealmName).Return(groups, nil).Times(1)
		mockKeycloakClient.EXPECT().GetClients(accessToken, targetRealmName).Return([]kc.ClientRepresentation{}, fmt.Errorf("Unexpected error")).Times(1)
		mockLogger.EXPECT().Warn(ctx, "err", "Unexpected error").Times(1)
		err = managementComponent.UpdateAuthorizations(ctx, targetRealmName, groupID, apiAuthorizations)
		assert.NotNil(t, err)

		mockKeycloakClient.EXPECT().GetGroup(accessToken, targetRealmName, groupID).Return(group, nil).Times(1)
		mockKeycloakClient.EXPECT().GetRealms(accessToken).Return(realms, nil).Times(1)
		mockKeycloakClient.EXPECT().GetGroups(accessToken, targetRealmName).Return(groups, nil).Times(1)
		mockKeycloakClient.EXPECT().GetClients(accessToken, targetRealmName).Return(clients, nil).Times(1)
		mockKeycloakClient.EXPECT().GetAvailableGroupClientRoles(accessToken, targetRealmName, groupID, ID).Return([]kc.RoleRepresentation{}, fmt.Errorf("Unexpected error")).Times(1)
		mockLogger.EXPECT().Warn(ctx, "err", "Unexpected error").Times(1)
		err = managementComponent.UpdateAuthorizations(ctx, targetRealmName, groupID, apiAuthorizations)
		assert.NotNil(t, err)

		mockKeycloakClient.EXPECT().GetGroup(accessToken, targetRealmName, groupID).Return(group, nil).Times(1)
		mockKeycloakClient.EXPECT().GetRealms(accessToken).Return(realms, nil).Times(1)
		mockKeycloakClient.EXPECT().GetGroups(accessToken, targetRealmName).Return(groups, nil).Times(1)
		mockKeycloakClient.EXPECT().GetClients(accessToken, targetRealmName).Return(clients, nil).Times(1)
		mockKeycloakClient.EXPECT().GetAvailableGroupClientRoles(accessToken, targetRealmName, groupID, ID).Return(rolesAvailable, nil).Times(1)
		mockKeycloakClient.EXPECT().GetGroupClientRoles(accessToken, targetRealmName, groupID, ID).Return([]kc.RoleRepresentation{}, fmt.Errorf("Unexpected error")).Times(1)
		mockLogger.EXPECT().Warn(ctx, "err", "Unexpected error").Times(1)
		err = managementComponent.UpdateAuthorizations(ctx, targetRealmName, groupID, apiAuthorizations)
		assert.NotNil(t, err)

		mockKeycloakClient.EXPECT().GetGroup(accessToken, targetRealmName, groupID).Return(group, nil).Times(1)
		mockKeycloakClient.EXPECT().GetRealms(accessToken).Return(realms, nil).Times(1)
		mockKeycloakClient.EXPECT().GetGroups(accessToken, targetRealmName).Return(groups, nil).Times(1)
		mockKeycloakClient.EXPECT().GetClients(accessToken, targetRealmName).Return(clients, nil).Times(1)
		mockKeycloakClient.EXPECT().GetAvailableGroupClientRoles(accessToken, targetRealmName, groupID, ID).Return(rolesAvailable, nil).Times(1)
		mockKeycloakClient.EXPECT().GetGroupClientRoles(accessToken, targetRealmName, groupID, ID).Return(rolesCurrent, nil).Times(1)
		mockKeycloakClient.EXPECT().AssignClientRole(accessToken, targetRealmName, groupID, ID, gomock.Any()).Return(fmt.Errorf("Unexpected error")).Times(1)
		mockLogger.EXPECT().Warn(ctx, "err", "Unexpected error").Times(1)
		err = managementComponent.UpdateAuthorizations(ctx, targetRealmName, groupID, apiAuthorizations)
		assert.NotNil(t, err)

		mockKeycloakClient.EXPECT().GetGroup(accessToken, targetRealmName, groupID).Return(group, nil).Times(1)
		mockKeycloakClient.EXPECT().GetRealms(accessToken).Return(realms, nil).Times(1)
		mockKeycloakClient.EXPECT().GetGroups(accessToken, targetRealmName).Return(groups, nil).Times(1)
		mockKeycloakClient.EXPECT().GetClients(accessToken, targetRealmName).Return(clients, nil).Times(1)
		mockKeycloakClient.EXPECT().GetAvailableGroupClientRoles(accessToken, targetRealmName, groupID, ID).Return(rolesAvailable, nil).Times(1)
		mockKeycloakClient.EXPECT().GetGroupClientRoles(accessToken, targetRealmName, groupID, ID).Return(rolesCurrent, nil).Times(1)
		mockKeycloakClient.EXPECT().AssignClientRole(accessToken, targetRealmName, groupID, ID, gomock.Any()).Return(nil).Times(1)
		mockConfigurationDBModule.EXPECT().NewTransaction(ctx).Return(nil, fmt.Errorf("Unexpected error")).Times(1)
		mockLogger.EXPECT().Warn(ctx, "err", "Unexpected error").Times(1)
		err = managementComponent.UpdateAuthorizations(ctx, targetRealmName, groupID, apiAuthorizations)
		assert.NotNil(t, err)

		mockKeycloakClient.EXPECT().GetGroup(accessToken, targetRealmName, groupID).Return(group, nil).Times(1)
		mockKeycloakClient.EXPECT().GetRealms(accessToken).Return(realms, nil).Times(1)
		mockKeycloakClient.EXPECT().GetGroups(accessToken, targetRealmName).Return(groups, nil).Times(1)
		mockKeycloakClient.EXPECT().GetClients(accessToken, targetRealmName).Return(clients, nil).Times(1)
		mockKeycloakClient.EXPECT().GetAvailableGroupClientRoles(accessToken, targetRealmName, groupID, ID).Return(rolesAvailable, nil).Times(1)
		mockKeycloakClient.EXPECT().GetGroupClientRoles(accessToken, targetRealmName, groupID, ID).Return(rolesCurrent, nil).Times(1)
		mockKeycloakClient.EXPECT().AssignClientRole(accessToken, targetRealmName, groupID, ID, gomock.Any()).Return(nil).Times(1)
		mockConfigurationDBModule.EXPECT().NewTransaction(ctx).Return(mockTransaction, nil).Times(1)
		mockTransaction.EXPECT().Close().Times(1)
		mockConfigurationDBModule.EXPECT().DeleteAuthorizations(ctx, targetRealmName, groupName).Return(fmt.Errorf("Unexpected error")).Times(1)
		mockLogger.EXPECT().Warn(ctx, "err", "Unexpected error").Times(1)
		err = managementComponent.UpdateAuthorizations(ctx, targetRealmName, groupID, apiAuthorizations)
		assert.NotNil(t, err)

		mockKeycloakClient.EXPECT().GetGroup(accessToken, targetRealmName, groupID).Return(group, nil).Times(1)
		mockKeycloakClient.EXPECT().GetRealms(accessToken).Return(realms, nil).Times(1)
		mockKeycloakClient.EXPECT().GetGroups(accessToken, targetRealmName).Return(groups, nil).Times(1)
		mockKeycloakClient.EXPECT().GetClients(accessToken, targetRealmName).Return(clients, nil).Times(1)
		mockKeycloakClient.EXPECT().GetAvailableGroupClientRoles(accessToken, targetRealmName, groupID, ID).Return(rolesAvailable, nil).Times(1)
		mockKeycloakClient.EXPECT().GetGroupClientRoles(accessToken, targetRealmName, groupID, ID).Return(rolesCurrent, nil).Times(1)
		mockKeycloakClient.EXPECT().AssignClientRole(accessToken, targetRealmName, groupID, ID, gomock.Any()).Return(nil).Times(1)
		mockConfigurationDBModule.EXPECT().NewTransaction(ctx).Return(mockTransaction, nil).Times(1)
		mockTransaction.EXPECT().Close().Times(1)
		mockConfigurationDBModule.EXPECT().DeleteAuthorizations(ctx, targetRealmName, groupName).Return(nil).Times(1)
		mockConfigurationDBModule.EXPECT().CreateAuthorization(ctx, gomock.Any()).Return(fmt.Errorf("Unexpected error")).Times(1)
		mockLogger.EXPECT().Warn(ctx, "err", "Unexpected error").Times(1)
		err = managementComponent.UpdateAuthorizations(ctx, targetRealmName, groupID, apiAuthorizations)
		assert.NotNil(t, err)

		mockKeycloakClient.EXPECT().GetGroup(accessToken, targetRealmName, groupID).Return(group, nil).Times(1)
		mockKeycloakClient.EXPECT().GetRealms(accessToken).Return(realms, nil).Times(1)
		mockKeycloakClient.EXPECT().GetGroups(accessToken, targetRealmName).Return(groups, nil).Times(1)
		mockKeycloakClient.EXPECT().GetClients(accessToken, targetRealmName).Return(clients, nil).Times(1)
		mockKeycloakClient.EXPECT().GetAvailableGroupClientRoles(accessToken, targetRealmName, groupID, ID).Return(rolesAvailable, nil).Times(1)
		mockKeycloakClient.EXPECT().GetGroupClientRoles(accessToken, targetRealmName, groupID, ID).Return(rolesCurrent, nil).Times(1)
		mockKeycloakClient.EXPECT().AssignClientRole(accessToken, targetRealmName, groupID, ID, gomock.Any()).Return(nil).Times(1)
		mockConfigurationDBModule.EXPECT().NewTransaction(ctx).Return(mockTransaction, nil).Times(1)
		mockConfigurationDBModule.EXPECT().DeleteAuthorizations(ctx, targetRealmName, groupName).Return(nil).Times(1)
		mockConfigurationDBModule.EXPECT().CreateAuthorization(ctx, gomock.Any()).Return(nil).Times(1)
		mockTransaction.EXPECT().Close().Times(1)
		mockTransaction.EXPECT().Commit().Return(fmt.Errorf("Unexpected error")).Times(1)
		mockLogger.EXPECT().Warn(ctx, "err", "Unexpected error").Times(1)
		err = managementComponent.UpdateAuthorizations(ctx, targetRealmName, groupID, apiAuthorizations)
		assert.NotNil(t, err)

		mockKeycloakClient.EXPECT().GetGroup(accessToken, targetRealmName, groupID).Return(group, nil).Times(1)
		mockKeycloakClient.EXPECT().GetRealms(accessToken).Return(realms, nil).Times(1)
		mockKeycloakClient.EXPECT().GetGroups(accessToken, targetRealmName).Return(groups, nil).Times(1)
		mockKeycloakClient.EXPECT().GetClients(accessToken, targetRealmName).Return(clients, nil).Times(1)
		mockKeycloakClient.EXPECT().GetAvailableGroupClientRoles(accessToken, targetRealmName, groupID, ID).Return(rolesAvailable, nil).Times(1)
		mockKeycloakClient.EXPECT().GetGroupClientRoles(accessToken, targetRealmName, groupID, ID).Return(rolesCurrent, nil).Times(1)
		mockKeycloakClient.EXPECT().AssignClientRole(accessToken, targetRealmName, groupID, ID, gomock.Any()).Return(nil).Times(1)
		mockConfigurationDBModule.EXPECT().NewTransaction(ctx).Return(mockTransaction, nil).Times(1)
		mockConfigurationDBModule.EXPECT().DeleteAuthorizations(ctx, targetRealmName, groupName).Return(nil).Times(1)
		mockConfigurationDBModule.EXPECT().CreateAuthorization(ctx, gomock.Any()).Return(nil).Times(1)
		mockTransaction.EXPECT().Close().Times(1)
		mockTransaction.EXPECT().Commit().Times(1)

		mockEventDBModule.EXPECT().ReportEvent(ctx, "API_AUTHORIZATIONS_UPDATE", "back-office", database.CtEventRealmName, targetRealmName, database.CtEventGroupName, groupName).Return(errors.New("error")).Times(1)
		m := map[string]interface{}{"event_name": "API_AUTHORIZATIONS_UPDATE", database.CtEventRealmName: targetRealmName, database.CtEventGroupName: groupName}
		eventJSON, _ := json.Marshal(m)
		mockLogger.EXPECT().Error(ctx, "err", "error", "event", string(eventJSON))
		err = managementComponent.UpdateAuthorizations(ctx, targetRealmName, groupID, apiAuthorizations)
		assert.Nil(t, err)
	})

	t.Run("YAT (yet another test)", func(t *testing.T) {
		var action = "action"
		var matrix = map[string]map[string]map[string]struct{}{
			action: {},
		}

		var apiAuthorizations = api.AuthorizationsRepresentation{
			Matrix: &matrix,
		}

		var rolesCurrent = []kc.RoleRepresentation{
			roleManageUser,
			roleViewClients,
			roleViewRealm,
			roleViewUsers,
			roleOther,
		}
		var rolesAvailable = []kc.RoleRepresentation{}

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		ctx = context.WithValue(ctx, cs.CtContextRealm, realmName)
		ctx = context.WithValue(ctx, cs.CtContextUsername, username)

		mockKeycloakClient.EXPECT().GetGroup(accessToken, targetRealmName, groupID).Return(group, nil).Times(1)
		mockKeycloakClient.EXPECT().GetRealms(accessToken).Return(realms, nil).Times(1)
		mockKeycloakClient.EXPECT().GetGroups(accessToken, targetRealmName).Return(groups, nil).Times(1)
		mockKeycloakClient.EXPECT().GetClients(accessToken, targetRealmName).Return(clients, nil).Times(1)
		mockKeycloakClient.EXPECT().GetAvailableGroupClientRoles(accessToken, targetRealmName, groupID, ID).Return(rolesAvailable, nil).Times(1)
		mockKeycloakClient.EXPECT().GetGroupClientRoles(accessToken, targetRealmName, groupID, ID).Return(rolesCurrent, nil).Times(1)
		mockKeycloakClient.EXPECT().RemoveClientRole(accessToken, targetRealmName, groupID, ID, gomock.Any()).Return(fmt.Errorf("Unexpected error")).Times(1)
		mockLogger.EXPECT().Warn(ctx, "err", "Unexpected error").Times(1)
		err := managementComponent.UpdateAuthorizations(ctx, targetRealmName, groupID, apiAuthorizations)
		assert.NotNil(t, err)
	})

	t.Run("Authorizations provided not valid", func(t *testing.T) {
		var jsonMatrix = `{
			"Action1": {},
			"Action2": {"*": {}, "realm1": {}}
		}`

		var matrix map[string]map[string]map[string]struct{}
		if err := json.Unmarshal([]byte(jsonMatrix), &matrix); err != nil {
			assert.Fail(t, "")
		}

		var apiAuthorizations = api.AuthorizationsRepresentation{
			Matrix: &matrix,
		}

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		ctx = context.WithValue(ctx, cs.CtContextRealm, realmName)
		ctx = context.WithValue(ctx, cs.CtContextUsername, username)

		mockKeycloakClient.EXPECT().GetGroup(accessToken, targetRealmName, groupID).Return(group, nil).Times(1)
		mockKeycloakClient.EXPECT().GetRealms(accessToken).Return(realms, nil).Times(1)
		mockKeycloakClient.EXPECT().GetGroups(accessToken, targetRealmName).Return(groups, nil).Times(1)
		mockLogger.EXPECT().Warn(ctx, "err", gomock.Any()).Times(1)
		err := managementComponent.UpdateAuthorizations(ctx, targetRealmName, groupID, apiAuthorizations)
		assert.NotNil(t, err)
	})
}

func TestUpdateAuthorizationsWithAny(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()
	var mockKeycloakClient = mock.NewKeycloakClient(mockCtrl)
	var mockUsersDetailsDBModule = mock.NewUsersDetailsDBModule(mockCtrl)
	var mockEventDBModule = mock.NewEventDBModule(mockCtrl)
	var mockTransaction = mock.NewTransaction(mockCtrl)
	var mockConfigurationDBModule = mock.NewConfigurationDBModule(mockCtrl)
	var mockLogger = mock.NewLogger(mockCtrl)
	var allowedTrustIDGroups = []string{"grp1", "grp2"}

	var managementComponent = NewComponent(mockKeycloakClient, mockUsersDetailsDBModule, mockEventDBModule, mockConfigurationDBModule, allowedTrustIDGroups, mockLogger)

	var accessToken = "TOKEN=="
	var currentRealmName = "master"
	var targetRealmName = "DEP"
	var targetMasterRealmName = "master"
	var groupID = "41dbf4a8-32a9-4000-8c17-edc854c31231"
	var groupName = "groupName"
	var username = "username"

	var group = kc.GroupRepresentation{
		ID:   &groupID,
		Name: &groupName,
	}
	var groups = []kc.GroupRepresentation{}
	var clients = []kc.ClientRepresentation{}

	var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
	ctx = context.WithValue(ctx, cs.CtContextRealm, currentRealmName)
	ctx = context.WithValue(ctx, cs.CtContextUsername, username)

	// Check * as target realm is forbidden for non master realm
	{
		var action = "action"
		var matrix = map[string]map[string]map[string]struct{}{
			action: {"*": {}},
		}

		var apiAuthorizations = api.AuthorizationsRepresentation{
			Matrix: &matrix,
		}

		var realm = kc.RealmRepresentation{
			ID:    &targetRealmName,
			Realm: &targetRealmName,
		}
		var realms = []kc.RealmRepresentation{realm}

		mockKeycloakClient.EXPECT().GetGroup(accessToken, targetRealmName, groupID).Return(group, nil).Times(1)
		mockKeycloakClient.EXPECT().GetRealms(accessToken).Return(realms, nil).Times(1)
		mockKeycloakClient.EXPECT().GetGroups(accessToken, targetRealmName).Return(groups, nil).Times(1)
		mockLogger.EXPECT().Warn(ctx, "err", gomock.Any()).Times(1)

		err := managementComponent.UpdateAuthorizations(ctx, targetRealmName, groupID, apiAuthorizations)

		assert.NotNil(t, err)
	}

	// Check * as target realm is allowed for master realm
	{
		var action = "action"
		var matrix = map[string]map[string]map[string]struct{}{
			action: {"*": {}},
		}

		var apiAuthorizations = api.AuthorizationsRepresentation{
			Matrix: &matrix,
		}

		var realm = kc.RealmRepresentation{
			ID:    &targetMasterRealmName,
			Realm: &targetMasterRealmName,
		}
		var realms = []kc.RealmRepresentation{realm}

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		ctx = context.WithValue(ctx, cs.CtContextRealm, targetMasterRealmName)
		ctx = context.WithValue(ctx, cs.CtContextUsername, username)

		mockKeycloakClient.EXPECT().GetGroup(accessToken, targetMasterRealmName, groupID).Return(group, nil).Times(1)
		mockKeycloakClient.EXPECT().GetRealms(accessToken).Return(realms, nil).Times(1)
		mockKeycloakClient.EXPECT().GetGroups(accessToken, targetMasterRealmName).Return(groups, nil).Times(1)
		mockKeycloakClient.EXPECT().GetClients(accessToken, targetMasterRealmName).Return(clients, nil).Times(1)

		mockConfigurationDBModule.EXPECT().NewTransaction(ctx).Return(mockTransaction, nil).Times(1)
		mockConfigurationDBModule.EXPECT().DeleteAuthorizations(ctx, targetMasterRealmName, groupName).Return(nil).Times(1)
		mockConfigurationDBModule.EXPECT().CreateAuthorization(ctx, gomock.Any()).Return(nil).Times(1)
		mockTransaction.EXPECT().Close().Times(1)
		mockTransaction.EXPECT().Commit().Times(1)

		mockEventDBModule.EXPECT().ReportEvent(ctx, "API_AUTHORIZATIONS_UPDATE", "back-office", gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil).Times(1)

		err := managementComponent.UpdateAuthorizations(ctx, targetMasterRealmName, groupID, apiAuthorizations)

		assert.Nil(t, err)
	}
}

func TestGetClientRoles(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()
	var mockKeycloakClient = mock.NewKeycloakClient(mockCtrl)
	var mockUsersDetailsDBModule = mock.NewUsersDetailsDBModule(mockCtrl)
	var mockEventDBModule = mock.NewEventDBModule(mockCtrl)
	var mockConfigurationDBModule = mock.NewConfigurationDBModule(mockCtrl)
	var mockLogger = log.NewNopLogger()
	var allowedTrustIDGroups = []string{"grp1", "grp2"}

	var managementComponent = NewComponent(mockKeycloakClient, mockUsersDetailsDBModule, mockEventDBModule, mockConfigurationDBModule, allowedTrustIDGroups, mockLogger)

	var accessToken = "TOKEN=="
	var realmName = "master"
	var clientID = "15436-464-4"

	// Get roles with succces
	{
		var id = "1234-7454-4516"
		var composite = false
		var containerID = "containerId"
		var description = "description role"
		var clientRole = true
		var name = "name"

		var kcRoleRep = kc.RoleRepresentation{
			ID:          &id,
			Name:        &name,
			ClientRole:  &clientRole,
			Composite:   &composite,
			ContainerID: &containerID,
			Description: &description,
		}

		var kcRolesRep []kc.RoleRepresentation
		kcRolesRep = append(kcRolesRep, kcRoleRep)

		mockKeycloakClient.EXPECT().GetClientRoles(accessToken, realmName, clientID).Return(kcRolesRep, nil).Times(1)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)

		apiRolesRep, err := managementComponent.GetClientRoles(ctx, "master", clientID)

		var apiRoleRep = apiRolesRep[0]
		assert.Nil(t, err)
		assert.Equal(t, id, *apiRoleRep.ID)
		assert.Equal(t, name, *apiRoleRep.Name)
		assert.Equal(t, clientRole, *apiRoleRep.ClientRole)
		assert.Equal(t, composite, *apiRoleRep.Composite)
		assert.Equal(t, containerID, *apiRoleRep.ContainerID)
		assert.Equal(t, description, *apiRoleRep.Description)
	}

	//Error
	{
		mockKeycloakClient.EXPECT().GetClientRoles(accessToken, realmName, clientID).Return([]kc.RoleRepresentation{}, fmt.Errorf("Unexpected error")).Times(1)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)

		_, err := managementComponent.GetClientRoles(ctx, "master", clientID)

		assert.NotNil(t, err)
	}
}

func TestCreateClientRole(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()
	var mockKeycloakClient = mock.NewKeycloakClient(mockCtrl)
	var mockUsersDetailsDBModule = mock.NewUsersDetailsDBModule(mockCtrl)
	var mockEventDBModule = mock.NewEventDBModule(mockCtrl)
	var mockConfigurationDBModule = mock.NewConfigurationDBModule(mockCtrl)
	var mockLogger = log.NewNopLogger()
	var allowedTrustIDGroups = []string{"grp1", "grp2"}

	var managementComponent = NewComponent(mockKeycloakClient, mockUsersDetailsDBModule, mockEventDBModule, mockConfigurationDBModule, allowedTrustIDGroups, mockLogger)

	var accessToken = "TOKEN=="
	var realmName = "master"
	var clientID = "456-789-147"

	// Add role with succces
	{
		var id = "1234-7454-4516"
		var composite = false
		var containerID = "containerId"
		var description = "description role"
		var clientRole = true
		var name = "client name"

		var locationURL = "http://location.url"

		mockKeycloakClient.EXPECT().CreateClientRole(accessToken, realmName, clientID, gomock.Any()).DoAndReturn(
			func(accessToken, realmName, clientID string, role kc.RoleRepresentation) (string, error) {
				assert.Equal(t, id, *role.ID)
				assert.Equal(t, name, *role.Name)
				assert.Equal(t, clientRole, *role.ClientRole)
				assert.Equal(t, composite, *role.Composite)
				assert.Equal(t, containerID, *role.ContainerID)
				assert.Equal(t, description, *role.Description)
				return locationURL, nil
			}).Times(1)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)

		var roleRep = api.RoleRepresentation{
			ID:          &id,
			Name:        &name,
			ClientRole:  &clientRole,
			Composite:   &composite,
			ContainerID: &containerID,
			Description: &description,
		}

		location, err := managementComponent.CreateClientRole(ctx, "master", clientID, roleRep)

		assert.Nil(t, err)
		assert.Equal(t, locationURL, location)
	}

	//Error
	{
		mockKeycloakClient.EXPECT().CreateClientRole(accessToken, realmName, clientID, gomock.Any()).Return("", fmt.Errorf("Unexpected error")).Times(1)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)

		_, err := managementComponent.CreateClientRole(ctx, "master", clientID, api.RoleRepresentation{})

		assert.NotNil(t, err)
	}
}

func TestGetRealmCustomConfiguration(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()
	var mockKeycloakClient = mock.NewKeycloakClient(mockCtrl)
	var mockUsersDetailsDBModule = mock.NewUsersDetailsDBModule(mockCtrl)
	var mockEventDBModule = mock.NewEventDBModule(mockCtrl)
	var mockConfigurationDBModule = mock.NewConfigurationDBModule(mockCtrl)
	var mockLogger = log.NewNopLogger()
	var allowedTrustIDGroups = []string{"grp1", "grp2"}

	var managementComponent = NewComponent(mockKeycloakClient, mockUsersDetailsDBModule, mockEventDBModule, mockConfigurationDBModule, allowedTrustIDGroups, mockLogger)

	var accessToken = "TOKEN=="
	var realmID = "master_id"

	// Get existing config
	{
		var id = realmID
		var keycloakVersion = "4.8.3"
		var realm = "master"
		var displayName = "Master"
		var enabled = true

		var kcRealmRep = kc.RealmRepresentation{
			ID:              &id,
			KeycloakVersion: &keycloakVersion,
			Realm:           &realm,
			DisplayName:     &displayName,
			Enabled:         &enabled,
		}

		mockKeycloakClient.EXPECT().GetRealm(accessToken, realmID).Return(kcRealmRep, nil).Times(1)

		var clientID = "ClientID"
		var redirectURI = "http://redirect.url.com/test"

		var realmConfig = configuration.RealmConfiguration{
			DefaultClientID:    &clientID,
			DefaultRedirectURI: &redirectURI,
		}

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		mockConfigurationDBModule.EXPECT().GetConfiguration(ctx, realmID).Return(realmConfig, nil).Times(1)

		configJSON, err := managementComponent.GetRealmCustomConfiguration(ctx, realmID)

		assert.Nil(t, err)
		assert.Equal(t, *configJSON.DefaultClientID, *realmConfig.DefaultClientID)
		assert.Equal(t, *configJSON.DefaultRedirectURI, *realmConfig.DefaultRedirectURI)
	}

	// Get empty config
	{
		var id = realmID
		var keycloakVersion = "4.8.3"
		var realm = "master"
		var displayName = "Master"
		var enabled = true

		var kcRealmRep = kc.RealmRepresentation{
			ID:              &id,
			KeycloakVersion: &keycloakVersion,
			Realm:           &realm,
			DisplayName:     &displayName,
			Enabled:         &enabled,
		}

		mockKeycloakClient.EXPECT().GetRealm(accessToken, realmID).Return(kcRealmRep, nil).Times(1)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		mockConfigurationDBModule.EXPECT().GetConfiguration(ctx, realmID).Return(configuration.RealmConfiguration{}, errorhandler.Error{}).Times(1)

		configJSON, err := managementComponent.GetRealmCustomConfiguration(ctx, realmID)

		assert.Nil(t, err)
		assert.Nil(t, configJSON.DefaultClientID)
		assert.Nil(t, configJSON.DefaultRedirectURI)
	}

	// Unknown realm
	{
		var id = realmID
		var keycloakVersion = "4.8.3"
		var realm = "master"
		var displayName = "Master"
		var enabled = true

		var kcRealmRep = kc.RealmRepresentation{
			ID:              &id,
			KeycloakVersion: &keycloakVersion,
			Realm:           &realm,
			DisplayName:     &displayName,
			Enabled:         &enabled,
		}

		mockKeycloakClient.EXPECT().GetRealm(accessToken, realmID).Return(kcRealmRep, errors.New("error")).Times(1)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)

		_, err := managementComponent.GetRealmCustomConfiguration(ctx, realmID)

		assert.NotNil(t, err)
	}

	// DB error
	{
		var id = realmID
		var keycloakVersion = "4.8.3"
		var realm = "master"
		var displayName = "Master"
		var enabled = true

		var kcRealmRep = kc.RealmRepresentation{
			ID:              &id,
			KeycloakVersion: &keycloakVersion,
			Realm:           &realm,
			DisplayName:     &displayName,
			Enabled:         &enabled,
		}

		mockKeycloakClient.EXPECT().GetRealm(accessToken, realmID).Return(kcRealmRep, nil).Times(1)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		mockConfigurationDBModule.EXPECT().GetConfiguration(ctx, realmID).Return(configuration.RealmConfiguration{}, errors.New("error")).Times(1)

		_, err := managementComponent.GetRealmCustomConfiguration(ctx, realmID)

		assert.NotNil(t, err)
	}
}

func TestUpdateRealmCustomConfiguration(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()
	var mockKeycloakClient = mock.NewKeycloakClient(mockCtrl)
	var mockUsersDetailsDBModule = mock.NewUsersDetailsDBModule(mockCtrl)
	var mockEventDBModule = mock.NewEventDBModule(mockCtrl)
	var mockConfigurationDBModule = mock.NewConfigurationDBModule(mockCtrl)
	var mockLogger = log.NewNopLogger()
	var allowedTrustIDGroups = []string{"grp1", "grp2"}

	var managementComponent = NewComponent(mockKeycloakClient, mockUsersDetailsDBModule, mockEventDBModule, mockConfigurationDBModule, allowedTrustIDGroups, mockLogger)

	var accessToken = "TOKEN=="
	var realmID = "master_id"

	var id = realmID
	var keycloakVersion = "4.8.3"
	var realm = "master"
	var displayName = "Master"
	var enabled = true

	var kcRealmRep = kc.RealmRepresentation{
		ID:              &id,
		KeycloakVersion: &keycloakVersion,
		Realm:           &realm,
		DisplayName:     &displayName,
		Enabled:         &enabled,
	}

	var clients = make([]kc.ClientRepresentation, 2)
	var clientID1 = "clientID1"
	var clientName1 = "clientName1"
	var redirectURIs1 = []string{"https://www.cloudtrust.io/*", "https://www.cloudtrust-old.com/*"}
	var clientID2 = "clientID2"
	var clientName2 = "clientName2"
	var redirectURIs2 = []string{"https://www.cloudtrust2.io/*", "https://www.cloudtrust2-old.com/*"}
	clients[0] = kc.ClientRepresentation{
		ClientID:     &clientID1,
		Name:         &clientName1,
		RedirectUris: &redirectURIs1,
	}
	clients[1] = kc.ClientRepresentation{
		ClientID:     &clientID2,
		Name:         &clientName2,
		RedirectUris: &redirectURIs2,
	}

	var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
	var clientID = "clientID1"
	var redirectURI = "https://www.cloudtrust.io/test"
	var configInit = api.RealmCustomConfiguration{
		DefaultClientID:    &clientID,
		DefaultRedirectURI: &redirectURI,
	}

	// Update config with appropriate values
	{
		mockKeycloakClient.EXPECT().GetRealm(accessToken, realmID).Return(kcRealmRep, nil).Times(1)
		mockKeycloakClient.EXPECT().GetClients(accessToken, realmID).Return(clients, nil).Times(1)
		mockConfigurationDBModule.EXPECT().StoreOrUpdateConfiguration(ctx, realmID, gomock.Any()).Return(nil).Times(1)
		err := managementComponent.UpdateRealmCustomConfiguration(ctx, realmID, configInit)

		assert.Nil(t, err)
	}

	// Update config with unknown client ID
	{
		mockKeycloakClient.EXPECT().GetRealm(accessToken, realmID).Return(kcRealmRep, nil).Times(1)
		mockKeycloakClient.EXPECT().GetClients(accessToken, realmID).Return(clients, nil).Times(1)

		var clientID = "clientID1Nok"
		var redirectURI = "https://www.cloudtrust.io/test"
		var configInit = api.RealmCustomConfiguration{
			DefaultClientID:    &clientID,
			DefaultRedirectURI: &redirectURI,
		}
		err := managementComponent.UpdateRealmCustomConfiguration(ctx, realmID, configInit)

		assert.NotNil(t, err)
		assert.IsType(t, commonhttp.Error{}, err)
		e := err.(commonhttp.Error)
		assert.Equal(t, 400, e.Status)
	}

	// Update config with invalid redirect URI
	{
		mockKeycloakClient.EXPECT().GetRealm(accessToken, realmID).Return(kcRealmRep, nil).Times(1)
		mockKeycloakClient.EXPECT().GetClients(accessToken, realmID).Return(clients, nil).Times(1)

		var clientID = "clientID1"
		var redirectURI = "https://www.cloudtrustnok.io/test"
		var configInit = api.RealmCustomConfiguration{
			DefaultClientID:    &clientID,
			DefaultRedirectURI: &redirectURI,
		}
		err := managementComponent.UpdateRealmCustomConfiguration(ctx, realmID, configInit)

		assert.NotNil(t, err)
		assert.IsType(t, commonhttp.Error{}, err)
		e := err.(commonhttp.Error)
		assert.Equal(t, 400, e.Status)
	}

	// Update config with invalid redirect URI (trying to take advantage of the dots in the url)
	{
		mockKeycloakClient.EXPECT().GetRealm(accessToken, realmID).Return(kcRealmRep, nil).Times(1)
		mockKeycloakClient.EXPECT().GetClients(accessToken, realmID).Return(clients, nil).Times(1)

		var clientID = "clientID1"
		var redirectURI = "https://wwwacloudtrust.io/test"
		var configInit = api.RealmCustomConfiguration{
			DefaultClientID:    &clientID,
			DefaultRedirectURI: &redirectURI,
		}
		err := managementComponent.UpdateRealmCustomConfiguration(ctx, realmID, configInit)

		assert.NotNil(t, err)
		assert.IsType(t, commonhttp.Error{}, err)
		e := err.(commonhttp.Error)
		assert.Equal(t, 400, e.Status)
	}

	// error while calling GetClients
	{
		mockKeycloakClient.EXPECT().GetRealm(accessToken, realmID).Return(kcRealmRep, nil).Times(1)
		mockKeycloakClient.EXPECT().GetClients(accessToken, realmID).Return([]kc.ClientRepresentation{}, errors.New("error")).Times(1)
		err := managementComponent.UpdateRealmCustomConfiguration(ctx, realmID, configInit)

		assert.NotNil(t, err)
	}

	// error while calling GetRealm
	{
		mockKeycloakClient.EXPECT().GetRealm(accessToken, realmID).Return(kc.RealmRepresentation{}, errors.New("error")).Times(1)
		err := managementComponent.UpdateRealmCustomConfiguration(ctx, realmID, configInit)

		assert.NotNil(t, err)
	}
}

func TestGetRealmAdminConfiguration(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()

	var mockKeycloakClient = mock.NewKeycloakClient(mockCtrl)
	var mockUsersDetailsDBModule = mock.NewUsersDetailsDBModule(mockCtrl)
	var mockEventDBModule = mock.NewEventDBModule(mockCtrl)
	var mockConfigurationDBModule = mock.NewConfigurationDBModule(mockCtrl)
	var logger = log.NewNopLogger()

	var allowedTrustIDGroups = []string{"grp1", "grp2"}
	var realmName = "myrealm"
	var realmID = "1234-5678"
	var accessToken = "acce-ssto-ken"
	var expectedError = errors.New("expectedError")
	var dbAdminConfig configuration.RealmAdminConfiguration
	var apiAdminConfig = api.ConvertRealmAdminConfigurationFromDBStruct(dbAdminConfig)
	var ctx = context.WithValue(context.TODO(), cs.CtContextAccessToken, accessToken)

	var component = NewComponent(mockKeycloakClient, mockUsersDetailsDBModule, mockEventDBModule, mockConfigurationDBModule, allowedTrustIDGroups, logger)

	t.Run("Request to Keycloak client fails", func(t *testing.T) {
		mockKeycloakClient.EXPECT().GetRealm(accessToken, realmName).Return(kc.RealmRepresentation{}, expectedError)
		var _, err = component.GetRealmAdminConfiguration(ctx, realmName)
		assert.Equal(t, expectedError, err)
	})
	t.Run("Request to database fails", func(t *testing.T) {
		mockKeycloakClient.EXPECT().GetRealm(accessToken, realmName).Return(kc.RealmRepresentation{ID: &realmID}, nil)
		mockConfigurationDBModule.EXPECT().GetAdminConfiguration(ctx, gomock.Any()).Return(dbAdminConfig, expectedError)
		var _, err = component.GetRealmAdminConfiguration(ctx, realmName)
		assert.Equal(t, expectedError, err)
	})
	t.Run("Success", func(t *testing.T) {
		mockKeycloakClient.EXPECT().GetRealm(accessToken, realmName).Return(kc.RealmRepresentation{ID: &realmID}, nil)
		mockConfigurationDBModule.EXPECT().GetAdminConfiguration(ctx, realmID).Return(dbAdminConfig, nil)
		var res, err = component.GetRealmAdminConfiguration(ctx, realmName)
		assert.Nil(t, err)
		assert.Equal(t, apiAdminConfig, res)
	})
}

func TestUpdateRealmAdminConfiguration(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()

	var mockKeycloakClient = mock.NewKeycloakClient(mockCtrl)
	var mockEventDBModule = mock.NewEventDBModule(mockCtrl)
	var mockUsersDetailsDBModule = mock.NewUsersDetailsDBModule(mockCtrl)
	var mockConfigurationDBModule = mock.NewConfigurationDBModule(mockCtrl)
	var logger = log.NewNopLogger()

	var allowedTrustIDGroups = []string{"grp1", "grp2"}
	var realmName = "myrealm"
	var realmID = "1234-5678"
	var accessToken = "acce-ssto-ken"
	var expectedError = errors.New("expectedError")
	var ctx = context.WithValue(context.TODO(), cs.CtContextAccessToken, accessToken)
	var adminConfig api.RealmAdminConfiguration

	var component = NewComponent(mockKeycloakClient, mockUsersDetailsDBModule, mockEventDBModule, mockConfigurationDBModule, allowedTrustIDGroups, logger)

	t.Run("Request to Keycloak client fails", func(t *testing.T) {
		mockKeycloakClient.EXPECT().GetRealm(accessToken, realmName).Return(kc.RealmRepresentation{}, expectedError)
		var err = component.UpdateRealmAdminConfiguration(ctx, realmName, adminConfig)
		assert.Equal(t, expectedError, err)
	})
	t.Run("Request to database fails", func(t *testing.T) {
		mockKeycloakClient.EXPECT().GetRealm(accessToken, realmName).Return(kc.RealmRepresentation{ID: &realmID}, nil)
		mockConfigurationDBModule.EXPECT().StoreOrUpdateAdminConfiguration(ctx, realmID, gomock.Any()).Return(expectedError)
		var err = component.UpdateRealmAdminConfiguration(ctx, realmName, adminConfig)
		assert.Equal(t, expectedError, err)
	})
	t.Run("Success", func(t *testing.T) {
		mockKeycloakClient.EXPECT().GetRealm(accessToken, realmName).Return(kc.RealmRepresentation{ID: &realmID}, nil)
		mockConfigurationDBModule.EXPECT().StoreOrUpdateAdminConfiguration(ctx, realmID, gomock.Any()).Return(nil)
		var err = component.UpdateRealmAdminConfiguration(ctx, realmName, adminConfig)
		assert.Nil(t, err)
	})
}

func createBackOfficeConfiguration(JSON string) dto.BackOfficeConfiguration {
	var conf dto.BackOfficeConfiguration
	json.Unmarshal([]byte(JSON), &conf)
	return conf
}

func TestRealmBackOfficeConfiguration(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()
	var mockKeycloakClient = mock.NewKeycloakClient(mockCtrl)
	var mockEventDBModule = mock.NewEventDBModule(mockCtrl)
	var mockUsersDetailsDBModule = mock.NewUsersDetailsDBModule(mockCtrl)
	var mockConfigurationDBModule = mock.NewConfigurationDBModule(mockCtrl)
	var mockLogger = log.NewNopLogger()
	var allowedTrustIDGroups = []string{"grp1", "grp2"}

	var component = NewComponent(mockKeycloakClient, mockUsersDetailsDBModule, mockEventDBModule, mockConfigurationDBModule, allowedTrustIDGroups, mockLogger)

	var realmID = "master_id"
	var groupName = "the.group"
	var config = api.BackOfficeConfiguration{}
	var ctx = context.WithValue(context.TODO(), cs.CtContextGroups, []string{"grp1", "grp2"})
	var largeConf = `
		{
			"realm1": {
				"a": [ "grp1" ]
			},
			"realm2": {
				"a": [ "grp1" ],
				"b": [ "grp2" ],
				"c": [ "grp1", "grp2" ]
			}
		}
	`
	var smallConf = `
		{
			"realm2": {
				"a": [ "grp1" ],
				"c": [ "grp2" ]
			}
		}
	`

	t.Run("UpdateRealmBackOfficeConfiguration - db.GetBackOfficeConfiguration fails", func(t *testing.T) {
		var expectedError = errors.New("db error")
		mockConfigurationDBModule.EXPECT().GetBackOfficeConfiguration(ctx, realmID, []string{groupName}).Return(nil, expectedError)
		var err = component.UpdateRealmBackOfficeConfiguration(ctx, realmID, groupName, config)
		assert.Equal(t, expectedError, err)
	})

	t.Run("UpdateRealmBackOfficeConfiguration - remove items", func(t *testing.T) {
		var dbConf = createBackOfficeConfiguration(largeConf)
		var requestConf, _ = api.NewBackOfficeConfigurationFromJSON(smallConf)
		mockConfigurationDBModule.EXPECT().GetBackOfficeConfiguration(ctx, realmID, []string{groupName}).Return(dbConf, nil)
		mockConfigurationDBModule.EXPECT().DeleteBackOfficeConfiguration(ctx, realmID, groupName, "realm1", nil, nil).Return(nil)
		mockConfigurationDBModule.EXPECT().DeleteBackOfficeConfiguration(ctx, realmID, groupName, "realm2", gomock.Not(nil), nil).Return(nil)
		mockConfigurationDBModule.EXPECT().DeleteBackOfficeConfiguration(ctx, realmID, groupName, "realm2", gomock.Not(nil), gomock.Not(nil)).Return(nil)
		var err = component.UpdateRealmBackOfficeConfiguration(ctx, realmID, groupName, requestConf)
		assert.Nil(t, err)
	})

	t.Run("UpdateRealmBackOfficeConfiguration - add items", func(t *testing.T) {
		var dbConf = createBackOfficeConfiguration(smallConf)
		var requestConf, _ = api.NewBackOfficeConfigurationFromJSON(largeConf)
		mockConfigurationDBModule.EXPECT().GetBackOfficeConfiguration(ctx, realmID, []string{groupName}).Return(dbConf, nil)
		mockConfigurationDBModule.EXPECT().InsertBackOfficeConfiguration(ctx, realmID, groupName, "realm1", "a", []string{"grp1"}).Return(nil)
		mockConfigurationDBModule.EXPECT().InsertBackOfficeConfiguration(ctx, realmID, groupName, "realm2", "b", []string{"grp2"}).Return(nil)
		mockConfigurationDBModule.EXPECT().InsertBackOfficeConfiguration(ctx, realmID, groupName, "realm2", "c", []string{"grp1"}).Return(nil)
		var err = component.UpdateRealmBackOfficeConfiguration(ctx, realmID, groupName, requestConf)
		assert.Nil(t, err)
	})

	t.Run("GetRealmBackOfficeConfiguration - error", func(t *testing.T) {
		var dbConf = createBackOfficeConfiguration(smallConf)
		var expectedError = errors.New("db error")
		mockConfigurationDBModule.EXPECT().GetBackOfficeConfiguration(ctx, realmID, []string{groupName}).Return(dbConf, expectedError)
		var res, err = component.GetRealmBackOfficeConfiguration(ctx, realmID, groupName)
		assert.Equal(t, expectedError, err)
		assert.Nil(t, res)
	})

	t.Run("GetRealmBackOfficeConfiguration - success", func(t *testing.T) {
		var dbConf = createBackOfficeConfiguration(smallConf)
		mockConfigurationDBModule.EXPECT().GetBackOfficeConfiguration(ctx, realmID, []string{groupName}).Return(dbConf, nil)
		var res, err = component.GetRealmBackOfficeConfiguration(ctx, realmID, groupName)
		assert.Nil(t, err)
		assert.Equal(t, api.BackOfficeConfiguration(dbConf), res)
	})

	t.Run("GetUserRealmBackOfficeConfiguration - db error", func(t *testing.T) {
		var dbError = errors.New("db error")
		var groups = ctx.Value(cs.CtContextGroups).([]string)
		mockConfigurationDBModule.EXPECT().GetBackOfficeConfiguration(ctx, realmID, groups).Return(nil, dbError)
		var _, err = component.GetUserRealmBackOfficeConfiguration(ctx, realmID)
		assert.Equal(t, dbError, err)
	})

	t.Run("GetUserRealmBackOfficeConfiguration - success", func(t *testing.T) {
		var dbConf = createBackOfficeConfiguration(smallConf)
		var groups = ctx.Value(cs.CtContextGroups).([]string)
		mockConfigurationDBModule.EXPECT().GetBackOfficeConfiguration(ctx, realmID, groups).Return(dbConf, nil)
		var res, err = component.GetUserRealmBackOfficeConfiguration(ctx, realmID)
		assert.Nil(t, err)
		assert.Equal(t, api.BackOfficeConfiguration(dbConf), res)
	})
}

func TestLinkShadowUser(t *testing.T) {
	var mockCtrl = gomock.NewController(t)
	defer mockCtrl.Finish()
	var mockKeycloakClient = mock.NewKeycloakClient(mockCtrl)
	var mockUsersDetailsDBModule = mock.NewUsersDetailsDBModule(mockCtrl)
	var mockEventDBModule = mock.NewEventDBModule(mockCtrl)
	var mockConfigurationDBModule = mock.NewConfigurationDBModule(mockCtrl)
	var mockLogger = mock.NewLogger(mockCtrl)
	var allowedTrustIDGroups = []string{"grp1", "grp2"}

	var managementComponent = NewComponent(mockKeycloakClient, mockUsersDetailsDBModule, mockEventDBModule, mockConfigurationDBModule, allowedTrustIDGroups, mockLogger)

	var accessToken = "TOKEN=="
	var username = "test"
	var realmName = "master"
	var userID = "41dbf4a8-32a9-4000-8c17-edc854c31231"
	var provider = "provider"

	// Create shadow user
	t.Run("Create shadow user successfully", func(t *testing.T) {
		fedIDKC := kc.FederatedIdentityRepresentation{UserName: &username, UserID: &userID}
		fedID := api.FederatedIdentityRepresentation{Username: &username, UserID: &userID}

		mockKeycloakClient.EXPECT().LinkShadowUser(accessToken, realmName, userID, provider, fedIDKC).Return(nil).Times(1)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		ctx = context.WithValue(ctx, cs.CtContextRealm, realmName)
		ctx = context.WithValue(ctx, cs.CtContextUsername, username)

		err := managementComponent.LinkShadowUser(ctx, realmName, userID, provider, fedID)

		assert.Nil(t, err)
	})

	// Error from KC client
	t.Run("Create shadow user - error at KC client", func(t *testing.T) {
		fedIDKC := kc.FederatedIdentityRepresentation{UserName: &username, UserID: &userID}
		fedID := api.FederatedIdentityRepresentation{Username: &username, UserID: &userID}

		mockKeycloakClient.EXPECT().LinkShadowUser(accessToken, realmName, userID, provider, fedIDKC).Return(fmt.Errorf("error")).Times(1)

		var ctx = context.WithValue(context.Background(), cs.CtContextAccessToken, accessToken)
		ctx = context.WithValue(ctx, cs.CtContextRealm, realmName)
		ctx = context.WithValue(ctx, cs.CtContextUsername, username)

		mockLogger.EXPECT().Warn(ctx, "err", "error")
		err := managementComponent.LinkShadowUser(ctx, realmName, userID, provider, fedID)

		assert.NotNil(t, err)
	})
}
