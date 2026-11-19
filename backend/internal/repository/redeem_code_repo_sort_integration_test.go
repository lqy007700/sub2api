//go:build integration

package repository

import (
	"fmt"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/Wei-Shaw/sub2api/internal/service"
)

func (s *RedeemCodeRepoSuite) TestListWithFilters_SortByValueAsc() {
	s.Require().NoError(s.repo.Create(s.ctx, &service.RedeemCode{Code: "VALUE-20", Type: service.RedeemTypeBalance, Value: 20, Status: service.StatusUnused}))
	s.Require().NoError(s.repo.Create(s.ctx, &service.RedeemCode{Code: "VALUE-10", Type: service.RedeemTypeBalance, Value: 10, Status: service.StatusUnused}))

	codes, _, err := s.repo.ListWithFilters(s.ctx, pagination.PaginationParams{
		Page:      1,
		PageSize:  10,
		SortBy:    "value",
		SortOrder: "asc",
	}, "", "", "")
	s.Require().NoError(err)
	s.Require().Len(codes, 2)
	s.Require().Equal("VALUE-10", codes[0].Code)
	s.Require().Equal("VALUE-20", codes[1].Code)
}

func (s *RedeemCodeRepoSuite) TestUserRedemptionsExcludePaymentCodesAndAdjustmentsBeforePagination() {
	user := s.createUser("orders@example.com")
	other := s.createUser("other-orders@example.com")
	now := time.Now().UTC()
	createCode := func(code string, codeType string, userID int64, offset time.Duration) {
		_, err := s.client.RedeemCode.Create().SetCode(code).SetType(codeType).
			SetStatus(service.StatusUsed).SetUsedBy(userID).SetUsedAt(now.Add(offset)).Save(s.ctx)
		s.Require().NoError(err)
	}
	createCode("PAYMENT-CODE", service.RedeemTypeBalance, user.ID, 3*time.Minute)
	createCode("ADMIN-CREDIT", "admin_balance", user.ID, 2*time.Minute)
	createCode("USER-REDEEM-BALANCE", service.RedeemTypeBalance, user.ID, time.Minute)
	createCode("USER-REDEEM-SUBSCRIPTION", service.RedeemTypeSubscription, user.ID, 0)
	createCode("OTHER-REDEEM", service.RedeemTypeConcurrency, other.ID, 4*time.Minute)
	_, err := s.client.PaymentOrder.Create().SetUserID(user.ID).SetUserEmail(user.Email).
		SetUserName("test").SetAmount(10).SetPayAmount(10).SetPaymentType("stripe").
		SetPaymentTradeNo("").SetRechargeCode("PAYMENT-CODE").SetExpiresAt(now).
		SetClientIP("127.0.0.1").SetSrcHost("localhost").Save(s.ctx)
	s.Require().NoError(err)

	for i, code := range []string{"USER-REDEEM-BALANCE", "USER-REDEEM-SUBSCRIPTION"} {
		rows, page, err := s.repo.ListUserRedemptions(s.ctx, user.ID, pagination.PaginationParams{Page: i + 1, PageSize: 1})
		s.Require().NoError(err)
		s.Require().Len(rows, 1)
		s.Equal(code, rows[0].Code)
		s.EqualValues(2, page.Total)
	}
}

func (s *RedeemCodeRepoSuite) TestHistoryPaginationIsolationAndOrdering() {
	user := s.createUser("history@example.com")
	other := s.createUser("other-history@example.com")
	now := time.Now().UTC().Truncate(time.Second)
	var ids []int64
	types := []string{service.RedeemTypeBalance, service.RedeemTypeConcurrency, service.RedeemTypeSubscription, "admin_balance", "admin_concurrency"}
	for i := 0; i < 105; i++ {
		usedAt := now
		if i == 0 {
			usedAt = now.Add(time.Hour)
		}
		code, err := s.client.RedeemCode.Create().SetCode(fmt.Sprintf("HISTORY-%d", i)).SetType(types[i%len(types)]).SetStatus(service.StatusUsed).SetUsedBy(user.ID).SetUsedAt(usedAt).Save(s.ctx)
		s.Require().NoError(err)
		ids = append(ids, code.ID)
	}
	_, err := s.client.RedeemCode.Create().SetCode("OTHER-HISTORY").SetType(service.RedeemTypeSubscription).SetStatus(service.StatusUsed).SetUsedBy(other.ID).SetUsedAt(now.Add(2 * time.Hour)).Save(s.ctx)
	s.Require().NoError(err)
	expected := []int64{ids[0]}
	for i := len(ids) - 1; i > 0; i-- {
		expected = append(expected, ids[i])
	}
	for _, size := range []int{20, 50, 100} {
		var got []int64
		for page := 1; page <= (105+size-1)/size+1; page++ {
			codes, result, err := s.repo.ListByUserPaginated(s.ctx, user.ID, pagination.PaginationParams{Page: page, PageSize: size}, "")
			s.Require().NoError(err)
			s.Equal(int64(105), result.Total)
			s.LessOrEqual(len(codes), size)
			for _, code := range codes {
				s.Equal(user.ID, *code.UsedBy)
				got = append(got, code.ID)
			}
		}
		s.Equal(expected, got)
	}
	codes, result, err := s.repo.ListByUserPaginated(s.ctx, other.ID, pagination.PaginationParams{Page: 1, PageSize: 20}, "")
	s.Require().NoError(err)
	s.Equal(int64(1), result.Total)
	s.Require().Len(codes, 1)
	s.Equal("OTHER-HISTORY", codes[0].Code)
	empty := s.createUser("empty-history@example.com")
	codes, result, err = s.repo.ListByUserPaginated(s.ctx, empty.ID, pagination.PaginationParams{Page: 1, PageSize: 20}, "")
	s.Require().NoError(err)
	s.Empty(codes)
	s.Zero(result.Total)
}
