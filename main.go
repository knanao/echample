package main

import (
	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

var templates map[string]*template.Template

type Template struct {
}

// Render はHTMLテンプレートにデータを埋め込んだ結果をWriterに書き込みます
func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return templates[name].ExecuteTemplate(w, "layout.html", data)
}

func main() {
	e := echo.New()

	// テンプレートを利用するためのRendererの設定
	t := &Template{}
	e.Renderer = t

	/* ミドルウェアを設定
	HTTP-Request時に動く
	*/
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// 静的ファイルのパスを設定
	e.Static("/public/css/", "./public/css")
	e.Static("/public/js/", "./public/js")
	e.Static("/public/img/", "./public/img")

	e.GET("/", HandleIndexGet)
	e.GET("/hello", HandleHelloGet)
	e.POST("/hello", HandleHelloPost)
	e.GET("/hello_form", HandleHelloFormGet)

	e.GET("/api/hello", HandleAPIHelloGet)
	e.POST("/api/hello", HandleAPIHelloPost)

	e.Logger.Fatal(e.Start(":3000"))
}

func init() {
	loadTemplates()
}

// 各HTMLテンプレートに共通レイアウトを適用した結果を保存
func loadTemplates() {
	var baseTemplate = "templates/layout.html"
	templates = make(map[string]*template.Template)
	templates["hello"] = template.Must(
		template.ParseFiles(baseTemplate, "templates/hello.html"))
	templates["hello_form"] = template.Must(
		template.ParseFiles(baseTemplate, "templates/hello_form.html"))
}

func HandleIndexGet(c echo.Context) error {
	return c.Render(http.StatusOK, "index", "World")
}

func HandleHelloGet(c echo.Context) error {
	greetingto := c.QueryParam("greetingto")
	return c.Render(http.StatusOK, "hello", greetingto)
}

func HandleHelloPost(c echo.Context) error {
	greetingto := c.FormValue("greetingto")
	return c.Render(http.StatusOK, "hello", greetingto)
}

func HandleHelloFormGet(c echo.Context) error {
	return c.Render(http.StatusOK, "hello_form", nil)
}

// /api/hello のGet時のJSONデータ生成処理
func HandleAPIHelloGet(c echo.Context) error {
	greetingto := c.QueryParam("greetingto")
	return c.JSON(http.StatusOK, map[string]interface{}{"hello": greetingto})
}

// /api/hello が受けとるJSONパラメータを定義
type HelloParam struct {
	GreetingTo string `json:"hello"`
}

// /api/hello のPost時のJSONデータ生成処理
func HandleAPIHelloPost(c echo.Context) error {
	param := new(HelloParam)
	if err := c.Bind(param); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"hello": param.GreetingTo})
}
