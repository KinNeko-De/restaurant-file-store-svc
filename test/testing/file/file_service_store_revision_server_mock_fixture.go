package file

import (
	"testing"

	apiRestaurantFile "github.com/kinneko-de/api-contract/golang/kinnekode/restaurant/file/v1"
	"github.com/stretchr/testify/mock"
)

func (mockStream *MockFileService_StoreRevisionServer) SetupSendAndClose(t *testing.T) func() *apiRestaurantFile.StoreFileResponse {
	var actualResponse *apiRestaurantFile.StoreFileResponse
	mockStream.EXPECT().SendAndClose(mock.Anything).Run(func(response *apiRestaurantFile.StoreFileResponse) {
		actualResponse = response
	}).Return(nil).Times(1)

	return func() *apiRestaurantFile.StoreFileResponse {
		return actualResponse
	}
}
