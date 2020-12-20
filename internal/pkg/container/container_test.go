package container_test

import (
	"context"
	"errors"
	"testing"

	"github.com/fajardm/gobackend-server/internal/pkg/container"
	"github.com/fajardm/gobackend-server/internal/pkg/container/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type ContainerTestSuite struct {
	suite.Suite
	ctx         context.Context
	mockCtrl    *gomock.Controller
	mockService *mock.MockService
	container   container.Container
}

func TestContainerTestSuite(t *testing.T) {
	suite.Run(t, new(ContainerTestSuite))
}

func (t *ContainerTestSuite) SetupTest() {
	t.ctx = context.Background()
	t.mockCtrl = gomock.NewController(t.T())
	t.mockService = mock.NewMockService(t.mockCtrl)
	t.container = container.New()
}

func (t *ContainerTestSuite) TearDownTest() {
	t.mockCtrl.Finish()
}

func (t *ContainerTestSuite) TestReady_StartupError() {
	mockError := errors.New("unexpected")

	t.mockService.EXPECT().Startup().Return(mockError)

	t.container.RegisterService("xxxx", t.mockService)
	err := t.container.Ready()

	t.Equal(mockError, err)
}

func (t *ContainerTestSuite) TestReady_NoError() {
	t.mockService.EXPECT().Startup().Return(nil)

	t.container.RegisterService("xxxx", t.mockService)
	err := t.container.Ready()

	t.NoError(err)
}

func (t *ContainerTestSuite) TestShutdown_ShutdownError() {
	mockError := errors.New("unexpected")

	t.mockService.EXPECT().Startup().Return(nil)
	t.mockService.EXPECT().Shutdown().Return(mockError)

	t.container.RegisterService("xxxx", t.mockService)
	t.container.Ready()
	t.container.Shutdown()
}

func (t *ContainerTestSuite) TestShutdown_NoError() {
	t.mockService.EXPECT().Startup().Return(nil)
	t.mockService.EXPECT().Shutdown().Return(nil)

	t.container.RegisterService("xxxx", t.mockService)
	t.container.Ready()
	t.container.Shutdown()
}
