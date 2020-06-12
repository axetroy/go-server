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
	var (
		userUUID   string // 用户 Socket 的 UUID
		waiterUUID string // 客服 Socket 的 UUID
	)

	// Create test server with the echo handler.
	s := httptest.NewServer(customer_service.CustomerServiceRouter)
	defer s.Close()

	// Convert http://127.0.0.1 to ws://127.0.0.
	u := "ws" + strings.TrimPrefix(s.URL, "http") + "/v1/ws/connect/user"

	// Connect to the server
	socket, _, err := websocket.DefaultDialer.Dial(u, nil)
	assert.Nil(t, err)
	defer socket.Close()

	userProfile, err := tester.CreateUser()

	assert.Nil(t, err)

	defer func() {
		tester.DeleteUserByUid(userProfile.Id)
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
			"token": token.JoinPrefixToken(userProfile.Token),
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

		assert.Equal(t, userProfile.Id, publicProfile.Id)
		assert.Equal(t, userProfile.Nickname, publicProfile.Nickname)
		assert.Equal(t, userProfile.Username, publicProfile.Username)
		assert.Equal(t, userProfile.Avatar, publicProfile.Avatar)
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

		userUUID = msg.To

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
		waiterSocket, _, err := websocket.DefaultDialer.Dial(u, nil)
		assert.Nil(t, err)

		defer waiterSocket.Close()

		waiterInfo, err := tester.CreateWaiter()

		defer func() {
			tester.DeleteUserByUid(waiterInfo.Id)
		}()

		assert.Nil(t, err)

		// 发起认证
		assert.Nil(t, waiterSocket.WriteJSON(ws.Message{
			Type: string(ws.TypeRequestWaiterAuth),
			Payload: map[string]interface{}{
				"token": token.JoinPrefixToken(waiterInfo.Token),
			},
		}))

		// 读取下一条消息，应该是认证成功
		{
			_, b, err := waiterSocket.ReadMessage()
			assert.Nil(t, err)

			var msg ws.Message

			assert.Nil(t, json.Unmarshal(b, &msg))

			assert.Equal(t, string(ws.TypeResponseWaiterAuthSuccess), msg.Type)
			assert.Equal(t, "", msg.From)
			assert.NotNil(t, msg.To)
			assert.NotNil(t, msg.Date)

			waiterUUID = msg.To

			var waiterProfile schema.ProfilePublic

			assert.Nil(t, util.Decode(&waiterProfile, msg.Payload))
			assert.Equal(t, waiterInfo.Id, waiterProfile.Id)
			assert.Equal(t, waiterInfo.Username, waiterProfile.Username)
			assert.Equal(t, waiterInfo.Nickname, waiterProfile.Nickname)
			assert.Equal(t, waiterInfo.Avatar, waiterProfile.Avatar)
		}

		// 发起 ready
		assert.Nil(t, waiterSocket.WriteJSON(ws.Message{
			Type: string(ws.TypeRequestWaiterReady),
		}))

		// 再读取吓一跳消息，应该是连接成功，应为楼上已经有在排队的用户
		{
			_, b, err := waiterSocket.ReadMessage()
			assert.Nil(t, err)

			var msg ws.Message

			assert.Nil(t, json.Unmarshal(b, &msg))

			assert.Equal(t, string(ws.TypeResponseWaiterNewConnection), msg.Type)
			assert.Equal(t, userUUID, msg.From)
			assert.Equal(t, waiterUUID, msg.To)
			assert.NotNil(t, msg.Date)

			var userProfile schema.ProfilePublic

			assert.Nil(t, util.Decode(&userProfile, msg.Payload))
			assert.Equal(t, userProfile.Id, userProfile.Id)
			assert.Equal(t, userProfile.Username, userProfile.Username)
			assert.Equal(t, userProfile.Nickname, userProfile.Nickname)
			assert.Equal(t, userProfile.Avatar, userProfile.Avatar)
		}

		// 用户读取消息，应该可以读取到客服的信息
		{
			_, b, err := socket.ReadMessage()
			assert.Nil(t, err)

			var msg ws.Message

			assert.Nil(t, json.Unmarshal(b, &msg))

			assert.Equal(t, string(ws.TypeResponseUserConnectSuccess), msg.Type)
			assert.Equal(t, waiterUUID, msg.From)
			assert.Equal(t, userUUID, msg.To)
			assert.NotNil(t, msg.Date)

			var p schema.ProfilePublic

			assert.Nil(t, util.Decode(&p, msg.Payload))
			assert.Equal(t, waiterInfo.Id, p.Id)
			assert.Equal(t, waiterInfo.Username, p.Username)
			assert.Equal(t, waiterInfo.Nickname, p.Nickname)
			assert.Equal(t, waiterInfo.Avatar, p.Avatar)
		}

		// 用户发送消息
		{
			assert.Nil(t, socket.WriteJSON(ws.Message{
				Type: string(ws.TypeRequestUserMessageText),
				Payload: map[string]interface{}{
					"message": "Hello world!",
				},
			}))

			// 读取消息回执
			_, b, err := socket.ReadMessage()
			assert.Nil(t, err)

			var msg ws.Message

			assert.Nil(t, json.Unmarshal(b, &msg))

			assert.Equal(t, string(ws.TypeResponseUserMessageTextSuccess), msg.Type)
			assert.Equal(t, userUUID, msg.From)
			assert.Equal(t, waiterUUID, msg.To)
			assert.NotNil(t, msg.Date)
			assert.Equal(t, map[string]interface{}{
				"message": "Hello world!",
			}, msg.Payload)

			// 客服读取消息，应该会收到
			{
				// 读取消息回执
				_, b, err := waiterSocket.ReadMessage()
				assert.Nil(t, err)

				var msg ws.Message

				assert.Nil(t, json.Unmarshal(b, &msg))

				assert.Equal(t, string(ws.TypeResponseWaiterMessageText), msg.Type)
				assert.Equal(t, userUUID, msg.From)
				assert.Equal(t, waiterUUID, msg.To)
				assert.NotNil(t, msg.Date)
				assert.Equal(t, map[string]interface{}{
					"message": "Hello world!",
				}, msg.Payload)
			}
		}

		// 客服反会一条消息
		{
			assert.Nil(t, waiterSocket.WriteJSON(ws.Message{
				Type: string(ws.TypeRequestUserMessageText),
				To:   userUUID,
				Payload: map[string]interface{}{
					"message": "你好!",
				},
			}))

			// 读取消息回执
			_, b, err := waiterSocket.ReadMessage()
			assert.Nil(t, err)

			var msg ws.Message

			assert.Nil(t, json.Unmarshal(b, &msg))

			assert.Equal(t, string(ws.TypeResponseUserMessageTextSuccess), msg.Type)
			assert.Equal(t, waiterUUID, msg.From)
			assert.Equal(t, userUUID, msg.To)
			assert.NotNil(t, msg.Date)
			assert.Equal(t, map[string]interface{}{
				"message": "你好!",
			}, msg.Payload)

			// 用户读取消息，应该会收到
			{
				// 读取消息回执
				_, b, err := socket.ReadMessage()
				assert.Nil(t, err)

				var msg ws.Message

				assert.Nil(t, json.Unmarshal(b, &msg))

				assert.Equal(t, string(ws.TypeResponseUserMessageText), msg.Type)
				assert.Equal(t, waiterUUID, msg.From)
				assert.Equal(t, userUUID, msg.To)
				assert.NotNil(t, msg.Date)
				assert.Equal(t, map[string]interface{}{
					"message": "你好!",
				}, msg.Payload)
			}
		}
	}
}
