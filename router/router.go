package router

import (
	"go-rest-api/controller"
	"net/http"
	"os"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewRouter(uc controller.IUserController, tc controller.ITaskController) *echo.Echo {
	// Echoインスタンスを作成
	e := echo.New()

	// CORSミドルウェアの設定
	// Cross-Origin Resource Sharing (CORS)を設定し、異なるオリジン間でのリソース共有を制御
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		// 許可するオリジン（フロントエンドのURL）を指定
		AllowOrigins: []string{"http://localhost:3000", os.Getenv("FE_URL")},
		// CSRFトークンを含むヘッダーを許可するための設定
		// HeaderXCSRFTokenを含めることで、フロントエンドからCSRFトークンを送信可能に
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAccessControlAllowHeaders, echo.HeaderXCSRFToken},
		AllowMethods:     []string{echo.GET, echo.POST, echo.PUT, echo.DELETE},
		AllowCredentials: true,
	}))

	// CSRFミドルウェアの設定
	// Cross-Site Request Forgery (CSRF)攻撃から保護するための設定
	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		// CSRFトークンを含むCookieの設定
		// CookiePath: Cookieが有効なパスを指定（"/"で全てのパスで有効）
		CookiePath: "/",
		// CookieDomain: Cookieが有効なドメインを指定
		CookieDomain: os.Getenv("API_DOMAIN"),
		// SameSite属性の設定：クロスサイトリクエスト時のCookieの送信方法を制御
		// DefaultModeでは、クロスサイトリクエストでCookieは送信されるが、一部の制限あり
		CookieSameSite: http.SameSiteDefaultMode, // postManでの検証用
		// CookkieSameSite: http.SameSiteNoneMode, // 本番環境ではSameSiteNoneModeにする
		// CookkieMaxAge: 60 // 60秒後にCookieを削除
	}))
	e.POST("/signup", uc.SignUp)
	e.POST("/login", uc.LogIn)
	e.POST("/logout", uc.LogOut)
	e.GET("/csrf", uc.CsrfToken)
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
