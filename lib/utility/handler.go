package utility

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"thesis/lib/component"

	"github.com/aws/aws-lambda-go/events"
	"github.com/awslabs/aws-lambda-go-api-proxy/core"
)

// se vuoi loggare l'errore utilizza l'error di risposta ma questa non intacca la risposta http
// questa devi comporla tu, quindi con 500 e body il testo dell'errore
// {"message":"Internal Server Error"}
func AmazonHandler(ctx context.Context, request events.APIGatewayV2HTTPRequest, page component.Page) (events.APIGatewayV2HTTPResponse, error) {
	r := core.RequestAccessorV2{}
	httpRequest, err := r.EventToRequestWithContext(ctx, request)
	if err != nil {
		return AmazonHandlerError(err)
	}
	responseWriter := core.NewProxyResponseWriterV2()

	err = PageEntryPoint(responseWriter, httpRequest, page)
	if err != nil {
		return AmazonHandlerError(err)
	}

	response, err := responseWriter.GetProxyResponse()
	if err != nil {
		return AmazonHandlerError(err)
	}
	return response, nil
}

func AmazonHandlerError(err error) (events.APIGatewayV2HTTPResponse, error) {
	return events.APIGatewayV2HTTPResponse{
		StatusCode: 500,
		Body:       fmt.Sprintf("{\"message\":\"%s\"}", err.Error()),
	}, err // rilancio l'errore il body non viene utilizzato, per il momento va bene così
}

func PageEntryPoint(w http.ResponseWriter, r *http.Request, page component.Page) error {
	//se ha il path button
	//altrimenti GET/POST/DELETE/PATCH sono riferiti alla tabella

	path := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if (path[0] != page.GetId()) || (len(path) > 2) {
		return errors.New("404 - Not Found")
	}

	if len(path) == 2 {
		if (path[1] == "schema") && (r.Method == http.MethodGet) {
			s, err := page.GetSchema()
			if err != nil {
				return err
			}
			return responseByValue(w, s)
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
