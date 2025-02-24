package router

import (
	"go-rest-api/controller"
	"os"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func NewRouter(uc controller.IUserController, tc controller.ITaskController) *echo.Echo {
	e := echo.New()
	e.POST("/signup", uc.SignUp)
	e.POST("/login", uc.LogIn)
	e.POST("/logout", uc.LogOut)
	// タスク関連のエンドポイントをグループ化
	// "/tasks"以下のパスを持つエンドポイントをまとめて管理する
	t := e.Group("/tasks")

	// JWTミドルウェアの設定
	// 以下のエンドポイントにアクセスする際は、JWT認証が必要になる
	t.Use(echojwt.WithConfig(echojwt.Config{
		// 環境変数からJWTの署名キーを取得
		// このキーを使用してトークンの検証を行う
		SigningKey: []byte(os.Getenv("SECRET")),

		// トークンの検索場所を指定
		// この設定では、Cookieの"token"フィールドからJWTを探す
		// フロントエンドからのリクエストには、Cookieにトークンが含まれている必要がある
		TokenLookup: "cookie:token",
	}))
	t.GET("", tc.GetAllTasks)
	t.GET("/:taskId", tc.GetTaskById)
	t.POST("", tc.CreateTask)
	t.PUT("/:taskId", tc.UpdateTask)
	t.DELETE("/:taskId", tc.DeleteTask)
	return e
}
