package component

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type PageType int

const (
	PageList PageType = 0
	PageCard          = 1
	Home              = 2
)

type SectionType int

const (
	Repeater  SectionType = 0
	Group                 = 1
	PieChart              = 2
	LineChart             = 3
	SubPage               = 4
)

type FieldType int

const (
	TextType     FieldType = 0
	IntType                = 1
	DecimalType            = 2
	BooleanType            = 3
	DateType               = 4
	TimeType               = 5
	DatetimeType           = 6
)

type Page interface {
	GetSchema() ([]byte, error)
	Get(filters map[string][]string) ([]byte, error)
	Post(body []byte) ([]byte, error)
	Patch(body []byte) ([]byte, error)
	Delete(filters map[string][]string) ([]byte, error)
	Button(queryParams map[string][]string) ([]byte, error)
	GetId() string
}

func PageEntryPoint(w http.ResponseWriter, r *http.Request, page Page) error {
	//se ha il path button
	//altrimenti GET/POST/DELETE/PATCH sono riferiti alla tabella

	path := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if (path[0] != page.GetId()) || (len(path) > 2) {
		return errors.New("404 - Not Found")
	}

	if len(path) == 2 {
		if r.Method == http.MethodGet {
			switch path[1] {
			case "schema":
				s, err := page.GetSchema()
				if err != nil {
					return err
				}
				return responseByValue(w, s)
			case "button":
				s, err := page.Button(r.URL.Query())
				if err != nil {
					return err
				}
				return responseByValue(w, s)
			default:
				http.Error(w, "405 - Method Not Allowed", http.StatusMethodNotAllowed)
				return nil
			}
		} else {
			// http.Error(w, "405 - Method Not Allowed", http.StatusMethodNotAllowed)
			// return nil
			return errors.New("405 - Method Not Allowed")
		}
	}

	switch r.Method {
	case http.MethodGet:
		s, err := page.Get(r.URL.Query())
		if err != nil {
			return err
		}
		return responseByValue(w, s)
	case http.MethodPost:
		body, err := io.ReadAll(r.Body)
		if err != nil {
			return err
		}
		s, err := page.Post(body)
		if err != nil {
			return err
		}
		return responseByValue(w, s)
	case http.MethodDelete:
		s, err := page.Delete(r.URL.Query())
		if err != nil {
			return err
		}
		return responseByValue(w, s)
	case http.MethodPatch:
		body, err := io.ReadAll(r.Body)
		if err != nil {
			return err
		}
		s, err := page.Patch(body)
		if err != nil {
			return err
		}
		return responseByValue(w, s)
	default:
		http.Error(w, "405 - Method Not Allowed", http.StatusMethodNotAllowed)
	}
	return nil
}

func responseByValue(w http.ResponseWriter, b []byte) error {
	if _, err := fmt.Fprint(w, string(b)); err != nil {
		return err
	}
	return nil
}
