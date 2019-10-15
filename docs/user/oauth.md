### oAuth 认证登陆

[GET] /v1/oauth2/:provider

前端跳转到这个 URL 进行认证，对应的 `provider` 有

| Provider | 说明                     |
| -------- | ------------------------ |
| github   | 使用 `Github` 帐号登陆   |
| gitlab   | 使用 `Gitlab` 帐号登陆   |
| twitter  | 使用 `Twitter` 帐号登陆  |
| facebook | 使用 `Facebook` 帐号登陆 |
| google   | 使用 `Google` 帐号登陆   |

### oAuth 认证成功的回调地址

[GET] /v1/oauth2/:provider/callback

前端页面跳转到 `/v1/oauth2/:provider` 接口认证成功后，跳转会来的 URL

这个一般不需要前端设置，只需要在对应的 `provider` 那里进行设置回调地址
