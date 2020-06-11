package connect_test

import (
	"encoding/json"
	"github.com/axetroy/go-server/internal/app/customer_service"
	"github.com/axetroy/go-server/internal/app/customer_service/ws"
	"github.com/axetroy/go-server/internal/library/util"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/token"
	"github.com/axetroy/go-server/tester"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestUserRouter(t *testing.T) {
	// Create test server with the echo handler.
	s := httptest.NewServer(customer_service.CustomerServiceRouter)
	defer s.Close()

	// Convert http://127.0.0.1 to ws://127.0.0.
	u := "ws" + strings.TrimPrefix(s.URL, "http") + "/v1/ws/connect/user"

	// Connect to the server
	socket, _, err := websocket.DefaultDialer.Dial(u, nil)
	assert.Nil(t, err)
	defer socket.Close()

	profile, err := tester.CreateUser()

	assert.Nil(t, err)

	defer func() {
		tester.DeleteUserByUid(profile.Id)
	}()

	type TestCase struct {
		Send    ws.Message // 发送的消息
		Receive ws.Message // 接收的消息
	}

	TestCases := []TestCase{
		{
			Send:    ws.Message{Type: string(ws.TypeRequestUserMessageText)},
			Receive: ws.Message{Type: string(ws.TypeResponseUserError)},
		},
	}

	// 没有登录发送消息
	for _, testCase := range TestCases {
		body, _ := json.Marshal(testCase.Send)
		// 发送
		assert.Nil(t, socket.WriteMessage(websocket.TextMessage, body))
		_, b, err := socket.ReadMessage()
		assert.Nil(t, err)

		var msg ws.Message

		assert.Nil(t, json.Unmarshal(b, &msg))

		assert.Equal(t, testCase.Receive.Type, msg.Type)
		assert.Equal(t, "", msg.From)
		assert.NotNil(t, msg.To)
		assert.NotNil(t, msg.Date)
	}

	// 请求身份认证
	{
		reqMsg := ws.Message{Type: string(ws.TypeRequestUserAuth), Payload: map[string]interface{}{
			"token": token.JoinPrefixToken(profile.Token),
		}}

		body, _ := json.Marshal(reqMsg)
		// 发送
		assert.Nil(t, socket.WriteMessage(websocket.TextMessage, body))

		_, b, err := socket.ReadMessage()
		assert.Nil(t, err)

		var msg ws.Message

		assert.Nil(t, json.Unmarshal(b, &msg))

		assert.Equal(t, string(ws.TypeResponseUserAuthSuccess), msg.Type)
		assert.Equal(t, "", msg.From)
		assert.NotNil(t, msg.To)
		assert.NotNil(t, msg.Date)

		var publicProfile schema.ProfilePublic
		assert.Nil(t, util.Decode(&publicProfile, msg.Payload))

		assert.Equal(t, profile.Id, publicProfile.Id)
		assert.Equal(t, profile.Nickname, publicProfile.Nickname)
		assert.Equal(t, profile.Username, publicProfile.Username)
		assert.Equal(t, profile.Avatar, publicProfile.Avatar)
	}

	// 请求连接，此时没有客服,应该会排队
	{
		reqMsg := ws.Message{Type: string(ws.TypeRequestUserConnect)}

		body, _ := json.Marshal(reqMsg)
		// 发送
		assert.Nil(t, socket.WriteMessage(websocket.TextMessage, body))

		_, b, err := socket.ReadMessage()
		assert.Nil(t, err)

		var msg ws.Message

		assert.Nil(t, json.Unmarshal(b, &msg))

		assert.Equal(t, string(ws.TypeResponseUserConnectQueue), msg.Type)
		assert.Equal(t, "", msg.From)
		assert.NotNil(t, msg.To)
		assert.NotNil(t, msg.Date)

		var payload struct {
			Location int `json:"location"`
		}
		assert.Nil(t, util.Decode(&payload, msg.Payload))

		assert.Equal(t, 0, payload.Location)
	}

	// 连接客服
	{
		// Create test server with the echo handler.
		waiterServer := httptest.NewServer(customer_service.CustomerServiceRouter)
		defer waiterServer.Close()

		// Convert http://127.0.0.1 to ws://127.0.0.
		u := "ws" + strings.TrimPrefix(waiterServer.URL, "http") + "/v1/ws/connect/waiter"

		// Connect to the server
		socket, _, err := websocket.DefaultDialer.Dial(u, nil)
		assert.Nil(t, err)

		defer socket.Close()

		//socket.WriteJSON(ws)
		assert.Nil(t, socket.WriteJSON(ws.Message{}))
	}
}
