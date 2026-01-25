package http

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	customerhttp "marketplace/internal/core/customer/api/http"
	referralhttp "marketplace/internal/core/referral/api/http"
	referralrepo "marketplace/internal/core/referral/repository"
	referralservice "marketplace/internal/core/referral/service"
	"marketplace/internal/infra/idempotency"
	pkgmiddleware "marketplace/pkg/http/middleware"
	testutil "marketplace/test"

	"github.com/stretchr/testify/require"
)

const referralAffiliatePath = "/referral/affiliate"
const referralCodePath = "/referral/code"
const referralValidatePath = "/referral/validate"
const referralRelationshipPath = "/referral/relationship"
const referralCommissionPath = "/referral/commission"

func TestReferralHTTPCreateAffiliateAndGet(t *testing.T) {
	server, cleanup := newReferralHTTPServer(t)
	defer cleanup()

	createPayload := referralhttp.CreateAffiliateRequest{
		FullName:    "John Doe",
		Email:       "john@example.com",
		PhoneNumber: "0912345678",
		Status:      "active",
		Platform:    "Instagram",
	}

	createdAffiliate := doCreateAffiliate(t, server, createPayload, http.StatusCreated)

	getURL := fmt.Sprintf("%s%s/%d", server.URL, referralAffiliatePath, createdAffiliate.ID)
	resp := doJSONRequest(t, server.Client(), http.MethodGet, getURL, nil)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	var fetched referralhttp.AffiliateResponse
	err := json.NewDecoder(resp.Body).Decode(&fetched)
	require.NoError(t, err)

	require.Equal(t, createdAffiliate.ID, fetched.ID)
	require.Equal(t, createPayload.FullName, fetched.FullName)
	require.Equal(t, createPayload.Email, fetched.Email)
	require.Equal(t, "active", fetched.Status)
	require.True(t, fetched.IsActive)
}

func TestReferralHTTPListAffiliates(t *testing.T) {
	server, cleanup := newReferralHTTPServer(t)
	defer cleanup()

	payloads := []referralhttp.CreateAffiliateRequest{
		{
			FullName:    "Affiliate One",
			Email:       "one@example.com",
			PhoneNumber: "0911111111",
			Status:      "active",
		},
		{
			FullName:    "Affiliate Two",
			Email:       "two@example.com",
			PhoneNumber: "0922222222",
			Status:      "inactive",
		},
	}

	for _, payload := range payloads {
		doCreateAffiliate(t, server, payload, http.StatusCreated)
	}

	listURL := fmt.Sprintf("%s%s?page=1&limit=10", server.URL, referralAffiliatePath)
	resp := doJSONRequest(t, server.Client(), http.MethodGet, listURL, nil)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	var listResponse referralhttp.AffiliateListResponse
	err := json.NewDecoder(resp.Body).Decode(&listResponse)
	require.NoError(t, err)

	require.GreaterOrEqual(t, listResponse.Total, 2)
	require.GreaterOrEqual(t, len(listResponse.Items), 2)
}

func TestReferralHTTPCreateReferralCode(t *testing.T) {
	server, cleanup := newReferralHTTPServer(t)
	defer cleanup()

	// First create an affiliate
	affiliatePayload := referralhttp.CreateAffiliateRequest{
		FullName:    "Test Affiliate",
		Email:       "affiliate@example.com",
		PhoneNumber: "0912345678",
		Status:      "active",
	}
	affiliate := doCreateAffiliate(t, server, affiliatePayload, http.StatusCreated)

	// Create referral code
	codePayload := referralhttp.CreateReferralCodeRequest{
		AffiliateID: affiliate.ID,
		CustomCode:  "TEST123",
		Status:      "active",
	}

	resp := doJSONRequest(t, server.Client(), http.MethodPost, server.URL+referralCodePath, codePayload)
	defer resp.Body.Close()

	require.Equal(t, http.StatusCreated, resp.StatusCode)

	var created referralhttp.ReferralCodeResponse
	err := json.NewDecoder(resp.Body).Decode(&created)
	require.NoError(t, err)
	require.NotZero(t, created.ID)
	require.Equal(t, "TEST123", created.Code)
	require.Equal(t, affiliate.ID, created.AffiliateID)
	require.True(t, created.IsCustom)
}

func TestReferralHTTPCreateReferralCodeAutoGenerate(t *testing.T) {
	server, cleanup := newReferralHTTPServer(t)
	defer cleanup()

	// First create an affiliate
	affiliatePayload := referralhttp.CreateAffiliateRequest{
		FullName:    "Test Affiliate",
		Email:       "affiliate2@example.com",
		PhoneNumber: "0912345679",
		Status:      "active",
	}
	affiliate := doCreateAffiliate(t, server, affiliatePayload, http.StatusCreated)

	// Create referral code without custom code (auto-generate)
	codePayload := referralhttp.CreateReferralCodeRequest{
		AffiliateID: affiliate.ID,
		Status:      "active",
	}

	resp := doJSONRequest(t, server.Client(), http.MethodPost, server.URL+referralCodePath, codePayload)
	defer resp.Body.Close()

	require.Equal(t, http.StatusCreated, resp.StatusCode)

	var created referralhttp.ReferralCodeResponse
	err := json.NewDecoder(resp.Body).Decode(&created)
	require.NoError(t, err)
	require.NotZero(t, created.ID)
	require.NotEmpty(t, created.Code)
	require.False(t, created.IsCustom)
}

func TestReferralHTTPCreateReferralRelationship(t *testing.T) {
	db, cleanupDB := testutil.SetupTestDB(t)
	defer cleanupDB()

	// Create customer server to create a customer first
	customerServer := newCustomerHTTPServerWithDB(t, db)
	customerPayload := customerhttp.CreateCustomerRequest{
		FullName:    "Test Customer",
		Email:       "customer@example.com",
		PhoneNumber: "0911111111",
		Status:      "active",
	}
	customer := doCreateCustomer(t, customerServer, customerPayload, http.StatusCreated)

	// Create referral server
	server := newReferralHTTPServerWithDB(t, db)

	// Create affiliate
	affiliatePayload := referralhttp.CreateAffiliateRequest{
		FullName:    "Test Affiliate",
		Email:       "affiliate3@example.com",
		PhoneNumber: "0912345680",
		Status:      "active",
	}
	affiliate := doCreateAffiliate(t, server, affiliatePayload, http.StatusCreated)

	// Create referral code
	codePayload := referralhttp.CreateReferralCodeRequest{
		AffiliateID: affiliate.ID,
		CustomCode:  "REL123",
		Status:      "active",
	}
	codeResp := doJSONRequest(t, server.Client(), http.MethodPost, server.URL+referralCodePath, codePayload)
	defer codeResp.Body.Close()

	var code referralhttp.ReferralCodeResponse
	err := json.NewDecoder(codeResp.Body).Decode(&code)
	require.NoError(t, err)

	// Create referral relationship
	relationshipPayload := referralhttp.CreateReferralRelationshipRequest{
		AffiliateID:    affiliate.ID,
		CustomerID:     customer.ID,
		ReferralCodeID: code.ID,
		Source:        "Instagram",
		Campaign:      "Summer2024",
		UTMSource:     "instagram",
		UTMMedium:     "social",
		UTMCampaign:   "summer2024",
	}

	relResp := doJSONRequest(t, server.Client(), http.MethodPost, server.URL+referralRelationshipPath, relationshipPayload)
	defer relResp.Body.Close()

	require.Equal(t, http.StatusCreated, relResp.StatusCode)

	var relationship referralhttp.ReferralRelationshipResponse
	err = json.NewDecoder(relResp.Body).Decode(&relationship)
	require.NoError(t, err)
	require.NotZero(t, relationship.ID)
	require.Equal(t, affiliate.ID, relationship.AffiliateID)
	require.Equal(t, customer.ID, relationship.CustomerID)
	require.Equal(t, "Instagram", relationship.Source)
}

func TestReferralHTTPCreateAffiliateEmailConflict(t *testing.T) {
	server, cleanup := newReferralHTTPServer(t)
	defer cleanup()

	payload := referralhttp.CreateAffiliateRequest{
		FullName:    "First Affiliate",
		Email:       "duplicate@example.com",
		PhoneNumber: "0911111111",
		Status:      "active",
	}

	doCreateAffiliate(t, server, payload, http.StatusCreated)

	// Try to create another with same email
	payload2 := referralhttp.CreateAffiliateRequest{
		FullName:    "Second Affiliate",
		Email:       "duplicate@example.com",
		PhoneNumber: "0922222222",
		Status:      "active",
	}

	resp := doJSONRequest(t, server.Client(), http.MethodPost, server.URL+referralAffiliatePath, payload2)
	defer resp.Body.Close()

	require.Equal(t, http.StatusConflict, resp.StatusCode)
}

func TestReferralHTTPGetReferralCode(t *testing.T) {
	server, cleanup := newReferralHTTPServer(t)
	defer cleanup()

	// Create affiliate
	affiliatePayload := referralhttp.CreateAffiliateRequest{
		FullName:    "Test Affiliate",
		Email:       "affiliate4@example.com",
		PhoneNumber: "0912345681",
		Status:      "active",
	}
	affiliate := doCreateAffiliate(t, server, affiliatePayload, http.StatusCreated)

	// Create referral code
	codePayload := referralhttp.CreateReferralCodeRequest{
		AffiliateID: affiliate.ID,
		CustomCode:  "GET123",
		Status:      "active",
	}
	codeResp := doJSONRequest(t, server.Client(), http.MethodPost, server.URL+referralCodePath, codePayload)
	defer codeResp.Body.Close()

	var created referralhttp.ReferralCodeResponse
	err := json.NewDecoder(codeResp.Body).Decode(&created)
	require.NoError(t, err)

	// Get referral code
	getURL := fmt.Sprintf("%s%s/%d", server.URL, referralCodePath, created.ID)
	resp := doJSONRequest(t, server.Client(), http.MethodGet, getURL, nil)
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode)

	var fetched referralhttp.ReferralCodeResponse
	err = json.NewDecoder(resp.Body).Decode(&fetched)
	require.NoError(t, err)

	require.Equal(t, created.ID, fetched.ID)
	require.Equal(t, "GET123", fetched.Code)
}

func newReferralHTTPServer(t *testing.T) (*httptest.Server, func()) {
	t.Helper()

	db, cleanupDB := testutil.SetupTestDB(t)

	repo := referralrepo.NewReferralRepository(db)
	service := referralservice.NewReferralService(repo)
	handler := referralhttp.NewReferralHandler(service)
	idempoStore := idempotency.NewStore(db)
	idempoMiddleware := pkgmiddleware.NewIdempotencyMiddleware(idempoStore)

	mux := http.NewServeMux()
	referralhttp.RegisterHTTPRoutes(mux, handler)

	server := httptest.NewServer(idempoMiddleware.Handle(mux))

	cleanup := func() {
		server.Close()
		cleanupDB()
	}

	return server, cleanup
}

func newReferralHTTPServerWithDB(t *testing.T, db *sql.DB) *httptest.Server {
	t.Helper()

	repo := referralrepo.NewReferralRepository(db)
	service := referralservice.NewReferralService(repo)
	handler := referralhttp.NewReferralHandler(service)
	idempoStore := idempotency.NewStore(db)
	idempoMiddleware := pkgmiddleware.NewIdempotencyMiddleware(idempoStore)

	mux := http.NewServeMux()
	referralhttp.RegisterHTTPRoutes(mux, handler)

	return httptest.NewServer(idempoMiddleware.Handle(mux))
}

func doCreateAffiliate(t *testing.T, server *httptest.Server, payload referralhttp.CreateAffiliateRequest, expectedStatus int) referralhttp.AffiliateResponse {
	t.Helper()

	resp := doJSONRequest(t, server.Client(), http.MethodPost, server.URL+referralAffiliatePath, payload)
	defer resp.Body.Close()

	if resp.StatusCode != expectedStatus {
		bodyBytes := make([]byte, 1024)
		n, _ := resp.Body.Read(bodyBytes)
		t.Logf("Response body: %s", string(bodyBytes[:n]))
	}
	require.Equal(t, expectedStatus, resp.StatusCode, "expected status %d, got %d", expectedStatus, resp.StatusCode)

	var created struct {
		Message   string                          `json:"message"`
		Affiliate referralhttp.AffiliateResponse `json:"affiliate"`
	}
	err := json.NewDecoder(resp.Body).Decode(&created)
	require.NoError(t, err)
	require.NotZero(t, created.Affiliate.ID)

	return created.Affiliate
}

