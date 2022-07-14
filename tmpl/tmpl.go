package tmpl

import (
	"embed"
	"errors"
	"fmt"
	"net/url"
	"path"
	"strings"
	"text/template"
	"time"

	"github.com/sirupsen/logrus"
)

//go:embed templates/*
var fs embed.FS
var embedTemplates map[string]*template.Template
var customTemplates map[string]*template.Template
var funcMap template.FuncMap

func init() {
	// func
	funcMap = template.FuncMap{
		"date": func(dt time.Time, zone string) string {
			loc, err := time.LoadLocation(zone)
			if err != nil {
				logrus.Error(err)
				return err.Error()
			}
			dt = dt.In(loc)
			return dt.Format("2006-01-02 15:04:05")
		},
		"isNonZeroDate": func(dt time.Time) bool {
			return !(dt == time.Time{})
		},
		"in": func(m map[string]string, key string) bool {
			_, ok := m[key]
			return ok
		},
		"toUpper": strings.ToUpper,
		"toLink": func(s string) string {
			return fmt.Sprintf("[%s](%s)", s, s)
		},
		"displayKV": func(k, v string) string {
			_, err := url.ParseRequestURI(v)
			if err != nil {
				return fmt.Sprintf("%s:%s", k, v)
			}
			return fmt.Sprintf("[%s](%s)", k, v)
		},
		"displayLabels": func(labels map[string]string) string {
			s := ""
			for k, v := range labels {
				switch k {
				case "alertname", "cluster", "prometheus", "stage", "uid", "instance", "reason", "service", "endpoint":
					continue
				default:
					s += fmt.Sprintf("%s:%s, ", k, v)
				}
			}
			return s
		},
		"contains": strings.Contains,
		"escapeQuotes": func(s string) string {
			s = fmt.Sprintf("%#v", s)
			s = strings.TrimPrefix(s, "\"")
			s = strings.TrimSuffix(s, "\"")
			return s
		},
	}

	// embed
	dir, err := fs.ReadDir("templates")
	if err != nil {
		panic(err)
	}

	embedTemplates = make(map[string]*template.Template)
	for _, entry := range dir {
		if entry.IsDir() {
			continue
		}

		filename := entry.Name()
		if !strings.HasSuffix(filename, ".tmpl") {
			continue
		}

		t, err := template.New(filename).Funcs(funcMap).ParseFS(fs, "templates/"+filename)
		if err != nil {
			panic(err)
		}

		embedTemplates[t.Name()] = t
	}

	// custom
	customTemplates = make(map[string]*template.Template)
}

func GetEmbedTemplate(filename string) (*template.Template, error) {
	if t, ok := embedTemplates[filename]; ok {
		return t, nil
	}

	return nil, errors.New("template not found")
}

func GetCustomTemplate(filepath string) (*template.Template, error) {
	if t, ok := customTemplates[filepath]; ok {
		return t, nil
	}

	t, err := template.New(path.Base(filepath)).Funcs(funcMap).ParseFiles(filepath)
	if err != nil {
		return nil, err
	}
	customTemplates[filepath] = t

	return t, nil
}
