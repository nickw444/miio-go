package capability

import (
	"testing"

	"github.com/nickw444/miio-go/protocol/transport"
	transportMocks "github.com/nickw444/miio-go/protocol/transport/mocks"
	subscriptionMocks "github.com/nickw444/miio-go/subscription/common/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestMiiORGB_GetComponents1(t *testing.T) {
	m := miioRGB(0xffffff)
	r, g, b := m.GetComponents()
	assert.Equal(t, 255, r)
	assert.Equal(t, 255, g)
	assert.Equal(t, 255, b)
}

func TestMiiORGB_GetComponents2(t *testing.T) {
	m := miioRGB(0xff7f0f)
	r, g, b := m.GetComponents()
	assert.Equal(t, 255, r)
	assert.Equal(t, 127, g)
	assert.Equal(t, 15, b)
}

func TestMiiORGB_SetComponents1(t *testing.T) {
	m := miioRGB(0)
	m.SetComponents(255, 255, 255)
	assert.Equal(t, miioRGB(0xffffff), m)
}

func TestMiiORGB_SetComponents2(t *testing.T) {
	m := miioRGB(0)
	m.SetComponents(255, 127, 15)
	assert.Equal(t, miioRGB(0xff7f0f), m)
}

func Light_SetUp() (tt struct {
	light    *Light
	outbound *transportMocks.Outbound
	target   *subscriptionMocks.SubscriptionTarget
}) {
	tt.target = new(subscriptionMocks.SubscriptionTarget)
	tt.outbound = new(transportMocks.Outbound)
	tt.light = NewLight(tt.target, tt.outbound)
	return
}

func TestLight_Update(t *testing.T) {
	tt := Light_SetUp()

	tt.outbound.On("CallAndDeserialize", mock.AnythingOfType("string"), mock.AnythingOfType("[]string"), mock.Anything).
		Return(nil).
		Run(func(args mock.Arguments) {
			resp := args.Get(2).(*transport.Response)
			resp.Result = []interface{}{"100", "3", "12345", "128", "100"}
		})
	tt.target.On("Publish", mock.Anything).Return(nil).Once()

	err := tt.light.Update()
	assert.NoError(t, err)
	tt.target.AssertExpectations(t)
}

func TestLight_UpdateNoChanges(t *testing.T) {
	tt := Light_SetUp()

	tt.outbound.On("CallAndDeserialize", mock.AnythingOfType("string"), mock.AnythingOfType("[]string"), mock.Anything).
		Return(nil).
		Run(func(args mock.Arguments) {
			resp := args.Get(2).(*transport.Response)
			resp.Result = []interface{}{"0", "0", "0", "0", "0"}
		})
	err := tt.light.Update()
	assert.NoError(t, err)
	tt.target.AssertExpectations(t)
}

func TestLight_SetRGB(t *testing.T) {
	tt := Light_SetUp()
	tt.outbound.On("Call", "set_rgb", []interface{}{16777215}).Return(nil, nil)
	tt.target.On("Publish", mock.Anything).Return(nil).Once()

	err := tt.light.SetRGB(255, 255, 255)
	assert.NoError(t, err)
	tt.target.AssertExpectations(t)
}

func TestLight_SetHSV(t *testing.T) {
	tt := Light_SetUp()
	tt.outbound.On("Call", "set_hsv", []interface{}{120, 77}).Return(nil, nil)
	tt.target.On("Publish", mock.Anything).Return(nil).Once()

	err := tt.light.SetHSV(120, 77)
	assert.NoError(t, err)
	tt.target.AssertExpectations(t)
}

func TestLight_SetBrightness(t *testing.T) {
	tt := Light_SetUp()
	tt.outbound.On("Call", "set_bright", []interface{}{55}).Return(nil, nil)
	tt.target.On("Publish", mock.Anything).Return(nil).Once()

	err := tt.light.SetBrightness(55)
	assert.NoError(t, err)
	tt.target.AssertExpectations(t)
}
