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
	distributorId = 1
}

func teardown() {
	testContainer.Teardown()
}

func Test_Init_HappyPath(t *testing.T) {
	t.Cleanup(teardown)
	ctx := context.Background()
	err := testContainer.DistributorApprovalService.CreateApprovalProcess(ctx, distributorId)
	if err != nil {
		t.Fatalf("failed to create distributor approval: %v", err)
	}

	//check
	approvalStatus, err := testContainer.DistributorApprovalService.GetApprovalStatus(ctx, distributorId)
	if err != nil {
		t.Fatalf("Failed to get err: %v", err)
	}
	wantApprovalStatus := PENDING_STATUS
	if approvalStatus.Status != wantApprovalStatus {
		t.Errorf("Expected approval status: %v Got: %v", wantApprovalStatus, approvalStatus)
	}
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
	wantApprovalStatus := APPROVED_STATUS
	if approvalStatus.Status != wantApprovalStatus {
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
	wantApprovalStatus := REJECTED_STATUS
	if approvalStatus.Status != wantApprovalStatus {
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

	status, err := testContainer.DistributorApprovalService.GetApprovalStatus(ctx, distributorId)
	if err != nil {
		t.Fatalf("failed to get approval status err: %v", err)
	}
	wantStatus := REJECTED_STATUS
	if status.Status != wantStatus {
		t.Errorf("Expected status: %v Got: %v", wantStatus, status.Status)
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
