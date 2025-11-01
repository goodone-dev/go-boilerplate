package html

import (
	"context"
	"fmt"
	"html/template"
	"io"
	"strings"
	"time"

	"github.com/goodone-dev/go-boilerplate/internal/infrastructure/logger"
)

func NewTemplate() *template.Template {
	tmpl, err := template.New("baseTemplate").Funcs(template.FuncMap{
		"FormatNumber": formatNumber,
		"FormatDate":   formatDate,
	}).ParseGlob("./templates/**/*.html")
	if err != nil {
		logger.Fatal(context.Background(), err, "❌ Failed to parse template")
		return nil
	}

	return tmpl
}

func ExecuteTemplate(wr io.Writer, file string, data any) (err error) {
	err = NewTemplate().ExecuteTemplate(wr, file, data)
	if err != nil {
		return
	}

	return
}

func formatNumber(amount float64) string {
	sign := ""
	if amount < 0 {
		sign = "-"
		amount = -amount
	}

	rounded := fmt.Sprintf("%.0f", amount)

	var formatted []string
	for i := len(rounded); i > 0; i -= 3 {
		if i-3 > 0 {
			formatted = append([]string{rounded[i-3 : i]}, formatted...)
		} else {
			formatted = append([]string{rounded[:i]}, formatted...)
		}
	}

	return sign + strings.Join(formatted, ".")
}

func formatDate(date time.Time) string {
	return date.Format("02 January 2006")
}
