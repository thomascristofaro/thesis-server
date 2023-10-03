package utility

import (
	"context"
	"fmt"
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

	err = component.PageEntryPoint(responseWriter, httpRequest, page)
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
