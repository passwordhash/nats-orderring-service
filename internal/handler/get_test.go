package handler

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"nats_server/internal/entity"
	"nats_server/internal/service"
	mock_service "nats_server/internal/service/mocks"
	"net/http/httptest"
	"testing"
)

func TestHandler_get(t *testing.T) {
	type mockBehavior func(r *mock_service.MockOrder, orderID entity.Order)

	tests := []struct {
		name                 string
		orderID              string
		order                entity.Order
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:    "Valid",
			orderID: "b563feb7b2b84b6test",
			order: entity.Order{
				OrderUID:    "b563feb7b2b84b6test",
				TrackNumber: "WBILMTESTTRACK",
				Entry:       "WBIL",
				Delivery: entity.Delivery{
					Name:    "Test Testov",
					Phone:   "+9720000000",
					Zip:     "2639809",
					City:    "Kiryat Mozkin",
					Address: "Ploshad Mira 15",
					Region:  "Kraiot",
					Email:   "some@email.ru",
				},
				Payment: entity.Payment{
					Transaction:  "b563feb7b2b84b6test",
					RequestID:    "",
					Currency:     "USD",
					Provider:     "wbpay",
					Amount:       1817,
					PaymentDt:    1637907727,
					Bank:         "alpha",
					DeliveryCost: 1500,
					GoodsTotal:   317,
					CustomFee:    0,
				},
				Items: []entity.Item{
					{
						ChrtID:      9934930,
						TrackNumber: "WBILMTESTTRACK",
						Price:       453,
						RID:         "ab4219087a764ae0btest",
						Name:        "Mascaras",
						Sale:        30,
						Size:        "0",
						TotalPrice:  317,
						NmID:        2389212,
						Brand:       "Vivienne Sabo",
						Status:      202,
					},
				},
			},
			mockBehavior: func(r *mock_service.MockOrder, order entity.Order) {
				r.EXPECT().Get(order.OrderUID).Return(order, nil).AnyTimes()
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"order_uid":"b563feb7b2b84b6test","track_number":"WBILMTESTTRACK","entry":"WBIL","delivery":{"name":"Test Testov","phone":"+9720000000","zip":"2639809","city":"Kiryat Mozkin","address":"Ploshad Mira 15","region":"Kraiot","email":"some@email.ru","delivery_dt":0,"delivery_type":"","delivery_cost":0,"delivery_status":0,"delivery_status_dt":0,"delivery_status_desc":"","delivery_status_desc_en":"","delivery_status_desc_he":""},"payment":{"transaction":"b563feb7b2b84b6test","request_id":"","currency":"USD","provider":"wbpay","amount":1817,"payment_dt":1637907727,"bank":"alpha","delivery_cost":1500,"goods_total":317,"custom_fee":0},"items":[{"chrt_id":9934930,"track_number":"WBILMTESTTRACK","price":453,"rid":"ab4219087a764ae0btest","name":"Mascaras","sale":30,"size":"0","total_price":317,"nm_id":2389212,"brand":"Vivienne Sabo","status":202}]}`,
		},
		{
			name:    "Order not found",
			orderID: "b563feb7b2b84b6test",
			order: entity.Order{
				OrderUID: "b563feb7b2b84b6test",
			},
			mockBehavior: func(r *mock_service.MockOrder, order entity.Order) {
				r.EXPECT().Get(order.OrderUID).Return(entity.Order{}, service.OrderNotFoundErr).AnyTimes()
			},
			expectedStatusCode:   400,
			expectedResponseBody: `{"error":"order not found"}`,
		},
		{
			name:    "Internal error",
			orderID: "b563feb7b2b84b6test",
			order: entity.Order{
				OrderUID: "b563feb7b2b84b6test",
			},
			mockBehavior: func(r *mock_service.MockOrder, order entity.Order) {
				r.EXPECT().Get(order.OrderUID).Return(entity.Order{}, errors.New("some internal error")).AnyTimes()
			},
			expectedStatusCode:   500,
			expectedResponseBody: `{"error":"some internal error"}`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			repo := mock_service.NewMockOrder(c)
			test.mockBehavior(repo, test.order)

			services := &service.Service{
				Order: repo,
			}
			handler := NewHandler(services)

			r := gin.New()
			r.GET("/get/:id", handler.get)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/get/"+test.orderID, nil)

			r.ServeHTTP(w, req)

			assert.Equal(t, test.expectedStatusCode, w.Code)

			var expectedOrder entity.Order
			_ = json.Unmarshal([]byte(test.expectedResponseBody), &expectedOrder)

			var actOrder entity.Order
			_ = json.Unmarshal(w.Body.Bytes(), &actOrder)

			assert.Equal(t, expectedOrder, actOrder)
		})
	}
}
