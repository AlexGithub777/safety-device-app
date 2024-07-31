package utils

import (
	"html/template"
	"io"
	"log"
	"math"
	"net"
	"time"

	"github.com/labstack/echo/v4"
)

// TemplateRenderer implements echo.Renderer
type TemplateRenderer struct {
	templates *template.Template
}

// NewTemplateRenderer creates a new TemplateRenderer
func NewTemplateRenderer() *TemplateRenderer {
	funcMap := FuncMap()
	t := template.Must(template.New("").Funcs(funcMap).ParseGlob("templates/*.html"))
	return &TemplateRenderer{templates: t}
}

// Render implements echo.Renderer
func (tr *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return tr.templates.ExecuteTemplate(w, name, data)
}

// FormatDate formats a time.Time to dd-mm-yyyy
func FormatDate(t time.Time) string {
	return t.Format("02-01-2006")
}

func FuncMap() template.FuncMap {
	return template.FuncMap{
		"formatDate": FormatDate,
		"add":        add,
		"sub":        sub,
		"mul":        mul,
		"div":        div,
		"ceil":       math.Ceil,
	}
}

func add(a, b int) int {
	return a + b
}

func sub(a, b int) int {
	return a - b
}

func mul(a, b int) int {
	return a * b
}

// Define the div function
func div(x, y int) int {
	if y == 0 {
		return 0 // Handle division by zero
	}
	return x / y
}

func GetLocalIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddress := conn.LocalAddr().(*net.UDPAddr)

	return localAddress.IP
}
