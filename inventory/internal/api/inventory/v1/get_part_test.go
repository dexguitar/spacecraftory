package v1

import (
	"errors"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/dexguitar/spacecraftory/inventory/internal/converter"
	"github.com/dexguitar/spacecraftory/inventory/internal/model"
	inventoryV1 "github.com/dexguitar/spacecraftory/shared/pkg/proto/inventory/v1"
)

func (s *APISuite) TestGetPartSuccess() {
	validUUID := "123e4567-e89b-12d3-a456-426614174000"

	servicePart := s.serviceMockData[validUUID]

	s.inventoryService.On("GetPart", s.ctx, validUUID).
		Return(servicePart, nil).Once()

	req := &inventoryV1.GetPartRequest{Uuid: validUUID}
	resp, err := s.api.GetPart(s.ctx, req)

	protoPart := converter.ToProtoPart(servicePart)

	s.Require().NoError(err)
	assert.Equal(s.T(), protoPart, resp.Part)
}

func (s *APISuite) TestGetPartError() {
	testCases := []struct {
		name             string
		uuid             string
		serviceError     error
		expectedCode     codes.Code
		expectedMsgParts []string
	}{
		{
			name:             "Part not found error",
			uuid:             "non-existent-uuid",
			serviceError:     model.ErrPartNotFound,
			expectedCode:     codes.NotFound,
			expectedMsgParts: []string{"not found", "non-existent-uuid"},
		},
		{
			name:             "Internal error",
			uuid:             "123e4567-e89b-12d3-a456-426614174000",
			serviceError:     errors.New("database connection failed"),
			expectedCode:     codes.Internal,
			expectedMsgParts: []string{"Internal server error"},
		},
	}

	for _, tc := range testCases {
		s.Run(tc.name, func() {
			s.inventoryService.On("GetPart", s.ctx, tc.uuid).
				Return(nil, tc.serviceError).Once()

			req := &inventoryV1.GetPartRequest{Uuid: tc.uuid}
			resp, err := s.api.GetPart(s.ctx, req)

			s.Require().Error(err)
			s.Require().Nil(resp)

			st, ok := status.FromError(err)
			s.Require().True(ok)
			assert.Equal(s.T(), tc.expectedCode, st.Code())
			for _, msgPart := range tc.expectedMsgParts {
				assert.Contains(s.T(), st.Message(), msgPart)
			}
		})
	}
}
