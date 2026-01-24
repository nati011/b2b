package service

import (
	"context"
	roleDomain "marketplace/internal/infra/authz/role/domain"
	"marketplace/internal/infra/user/domain"
	"marketplace/pkg/pagination"
	"errors"
	"fmt"
	"testing"
)

type stubRoleReader struct{}

func (stubRoleReader) FindByID(ctx context.Context, id string) (*roleDomain.Role, error) {
	return &roleDomain.Role{ID: id}, nil
}

func (stubRoleReader) GetByName(ctx context.Context, name string) (*roleDomain.Role, error) {
	return nil, errors.New("role not found")
}

// Benchmark tests for UserService

const (
	benchUserID     = "bench-user"
	benchExternalID = "ext-bench"
	benchUserName   = "Benchmark User"
	benchUserEmail  = "benchmark@example.com"
)

func BenchmarkUserServiceCreate(b *testing.B) {
	ctx := context.Background()
	mockRepo := newMockUserRepository()
	service := NewService(mockRepo, &stubRoleReader{}, nil, nil)
	email, _ := domain.NewEmail(benchUserEmail)
	var phoneNumber domain.PhoneNumber

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result, _ := service.Create(ctx, benchExternalID, email, phoneNumber, benchUserName, domain.UserTypeSelfService)
		_ = result
	}
}

func BenchmarkUserServiceGet(b *testing.B) {
	ctx := context.Background()
	mockRepo := newMockUserRepository()
	service := NewService(mockRepo, &stubRoleReader{}, nil, nil)

	// Pre-populate with test user
	user := createTestUser(benchUserID, benchExternalID, benchUserEmail, benchUserName, domain.UserTypeSelfService)
	mockRepo.users[benchUserID] = user

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.Get(ctx, benchUserID)
	}
}

func BenchmarkUserServiceFindByEmail(b *testing.B) {
	ctx := context.Background()
	mockRepo := newMockUserRepository()
	service := NewService(mockRepo, &stubRoleReader{}, nil, nil)

	// Pre-populate with test user
	user := createTestUser(benchUserID, benchExternalID, benchUserEmail, benchUserName, domain.UserTypeSelfService)
	mockRepo.users[benchUserID] = user
	mockRepo.usersByEmail[benchUserEmail] = user

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.FindByEmail(ctx, benchUserEmail)
	}
}

func BenchmarkUserServiceUpdate(b *testing.B) {
	ctx := context.Background()
	mockRepo := newMockUserRepository()
	service := NewService(mockRepo, &stubRoleReader{}, nil, nil)

	// Pre-populate with test user
	user := createTestUser(benchUserID, benchExternalID, benchUserEmail, benchUserName, domain.UserTypeSelfService)
	mockRepo.users[benchUserID] = user
	mockRepo.usersByEmail[benchUserEmail] = user

	email, _ := domain.NewEmail("updated@example.com")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		emptyStr := ""
		_, _ = service.Update(ctx, benchUserID, &emptyStr, email, "Updated Name")
	}
}

func BenchmarkUserServiceDelete(b *testing.B) {
	ctx := context.Background()
	mockRepo := newMockUserRepository()
	service := NewService(mockRepo, &stubRoleReader{}, nil, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Create user before each delete
		user := createTestUser(benchUserID, benchExternalID, benchUserEmail, benchUserName, domain.UserTypeSelfService)
		mockRepo.users[benchUserID] = user
		mockRepo.usersByEmail[benchUserEmail] = user

		_ = service.Delete(ctx, benchUserID)
	}
}

func BenchmarkUserServiceActivate(b *testing.B) {
	ctx := context.Background()
	mockRepo := newMockUserRepository()
	service := NewService(mockRepo, &stubRoleReader{}, nil, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Create inactive user before each activation
		user := createTestUser(benchUserID, benchExternalID, benchUserEmail, benchUserName, domain.UserTypeSelfService)
		user.Status = domain.UserStatusInactive
		mockRepo.users[benchUserID] = user

		_, _ = service.Activate(ctx, benchUserID)
	}
}

func BenchmarkUserServiceDeactivate(b *testing.B) {
	ctx := context.Background()
	mockRepo := newMockUserRepository()
	service := NewService(mockRepo, &stubRoleReader{}, nil, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Create active user before each deactivation
		user := createTestUser(benchUserID, benchExternalID, benchUserEmail, benchUserName, domain.UserTypeSelfService)
		user.Status = domain.UserStatusActive
		mockRepo.users[benchUserID] = user

		_, _ = service.Deactivate(ctx, benchUserID)
	}
}

func BenchmarkUserServiceSuspend(b *testing.B) {
	ctx := context.Background()
	mockRepo := newMockUserRepository()
	service := NewService(mockRepo, &stubRoleReader{}, nil, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Create active user before each suspension
		user := createTestUser(benchUserID, benchExternalID, benchUserEmail, benchUserName, domain.UserTypeSelfService)
		user.Status = domain.UserStatusActive
		mockRepo.users[benchUserID] = user

		_, _ = service.Suspend(ctx, benchUserID)
	}
}

func BenchmarkUserServiceFindPage(b *testing.B) {
	ctx := context.Background()
	mockRepo := newMockUserRepository()
	service := NewService(mockRepo, &stubRoleReader{}, nil, nil)

	// Pre-populate with multiple users
	for i := 0; i < 100; i++ {
		userID := fmt.Sprintf("user-%d", i)
		user := createTestUser(
			userID,
			fmt.Sprintf("ext-%d", i),
			fmt.Sprintf("user%d@example.com", i),
			fmt.Sprintf("User %d", i),
			domain.UserTypeSelfService,
		)
		mockRepo.users[user.ID] = user
		mockRepo.usersByEmail[user.Email.String()] = user
	}

	pageReq := pagination.NewPageRequest(1, 10)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.FindPage(ctx, pageReq)
	}
}

func BenchmarkUserServiceFindPageLargeDataset(b *testing.B) {
	ctx := context.Background()
	mockRepo := newMockUserRepository()
	service := NewService(mockRepo, &stubRoleReader{}, nil, nil)

	// Pre-populate with many users
	for i := 0; i < 1000; i++ {
		userID := fmt.Sprintf("user-%d", i)
		user := createTestUser(
			userID,
			fmt.Sprintf("ext-%d", i),
			fmt.Sprintf("user%d@example.com", i),
			fmt.Sprintf("User %d", i),
			domain.UserTypeSelfService,
		)
		mockRepo.users[user.ID] = user
		mockRepo.usersByEmail[user.Email.String()] = user
	}

	pageReq := pagination.NewPageRequest(1, 10)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.FindPage(ctx, pageReq)
	}
}
