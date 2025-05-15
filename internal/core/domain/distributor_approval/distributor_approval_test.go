package distributorApproval

import (
	"context"
	"os"
	"testing"
)

var testContainer TestContainer
var distributorId int

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func setup() {
	testContainer = NewTestContainer()
}

func teardown() {
	testContainer.Teardown()
}

func Test_Approve_HappyPath(t *testing.T) {
	t.Cleanup(teardown)
	ctx := context.Background()
	err := testContainer.DistributorApprovalService.Approve(ctx, &ApprovalRequest{
		DistributorId: distributorId,
	})
	if err != nil {
		t.Fatalf("failed to approve distributor: %v", err)
	}

	//check
	approvalStatus, err := testContainer.DistributorApprovalService.GetApprovalStatus(ctx, distributorId)
	if err != nil {
		t.Fatalf("Failed to get err: %v", err)
	}
	wantApprovalStatus := true
	if approvalStatus != wantApprovalStatus {
		t.Errorf("Expected approval status: %v Got: %v", wantApprovalStatus, approvalStatus)
	}
}

func Test_Approve_UnhappyPath(t *testing.T) {
	t.Run("alreadyApproved", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		err := testContainer.DistributorApprovalService.Approve(ctx, &ApprovalRequest{
			DistributorId: distributorId,
		})
		if err != nil {
			t.Fatalf("failed to approve distributor: %v", err)
		}

		err = testContainer.DistributorApprovalService.Approve(ctx, &ApprovalRequest{
			DistributorId: distributorId,
		})
		wantErr := ErrAlreadyApproved
		if err != wantErr {
			t.Errorf("Expected err: %v Got: %v", wantErr, err)
		}
	})
}

func Test_Reject_HappyPath(t *testing.T) {
	t.Cleanup(teardown)
	ctx := context.Background()
	err := testContainer.DistributorApprovalService.Reject(ctx, &RejectionRequest{
		DistributorId: distributorId,
		Comment:       "test",
	})
	if err != nil {
		t.Fatalf("failed to approve distributor: %v", err)
	}

	//check
	approvalStatus, err := testContainer.DistributorApprovalService.GetApprovalStatus(ctx, distributorId)
	if err != nil {
		t.Fatalf("Failed to get err: %v", err)
	}
	wantApprovalStatus := false
	if approvalStatus != wantApprovalStatus {
		t.Errorf("Expected approval status: %v Got: %v", wantApprovalStatus, approvalStatus)
	}
}

func Test_Reject_UnhappyPath(t *testing.T) {
	t.Run("alreadyRejected", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		err := testContainer.DistributorApprovalService.Reject(ctx, &RejectionRequest{
			DistributorId: distributorId,
			Comment:       "test",
		})
		if err != nil {
			t.Fatalf("failed to approve distributor: %v", err)
		}

		err = testContainer.DistributorApprovalService.Reject(ctx, &RejectionRequest{
			DistributorId: distributorId,
			Comment:       "test",
		})
		wantErr := ErrAlreadyRejected
		if err != wantErr {
			t.Errorf("Expected err: %v Got: %v", wantErr, err)
		}
	})

	t.Run("commentMandatory", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		err := testContainer.DistributorApprovalService.Reject(ctx, &RejectionRequest{
			DistributorId: distributorId,
		})

		wantErr := ErrCommentMandatory
		if err != wantErr {
			t.Errorf("Expected err: %v Got: %v", wantErr, err)
		}
	})
}

func Test_GetApprovalStatus_HappyPath(t *testing.T) {
	t.Cleanup(teardown)
	ctx := context.Background()
	err := testContainer.DistributorApprovalService.Reject(ctx, &RejectionRequest{
		DistributorId: distributorId,
		Comment:       "test",
	})
	if err != nil {
		t.Fatalf("failed to approve distributor: %v", err)
	}

	isApproved, err := testContainer.DistributorApprovalService.GetApprovalStatus(ctx, distributorId)
	if err != nil {
		t.Fatalf("failed to get approval status err: %v", err)
	}
	wantStatus := false
	if isApproved != wantStatus {
		t.Errorf("Expected status: %v Got: %v", wantStatus, isApproved)
	}
}

func Test_GetApprovalStatus_UnappyPath(t *testing.T) {
	t.Run("noReviews", func(t *testing.T) {
		t.Cleanup(teardown)
		ctx := context.Background()
		_, err := testContainer.DistributorApprovalService.GetApprovalStatus(ctx, distributorId)
		wantErr := ErrNoReview
		if err != wantErr {
			t.Errorf("Expected err: %v Got: %v", wantErr, err)
		}
	})
}

// func Test_GetReviewReport_HappyPath(t *testing.T) {
// }

// func Test_GetReviewReport_UnappyPath(t *testing.T) {
// 	t.Run("emptyGetContent", func(t *testing.T) {

// 	})
// }
