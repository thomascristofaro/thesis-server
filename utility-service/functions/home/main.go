package main

import (
	"context"
	"encoding/json"
	"thesis/lib/component"
	"thesis/lib/utility"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func main() {
	lambda.Start(
		func(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
			return utility.AmazonHandler(ctx, request, NewHome())
		})
}

type Home struct {
	component.BasePage
}

func NewHome() *Home {
	c := &Home{}
	c.BasePage = *component.NewBasePage(c, "Home Page", component.Home, "")
	c.AddSection(component.NewSection("pie_chart_1", component.PieChart, "Distribuzione degli Ordini per Cliente"))
	sec := component.NewSection("line_chart_1", component.LineChart, "Andamento Mensile degli Ordini")

	x_obj := map[string]interface{}{"min": 0, "max": 11}
	var data []map[string]interface{}
	data = append(data, map[string]interface{}{"label": "MAR", "value": 2})
	data = append(data, map[string]interface{}{"label": "JUN", "value": 5})
	data = append(data, map[string]interface{}{"label": "SEP", "value": 8})
	x_obj["data"] = data
	sec.Options["x"] = x_obj

	y_obj := map[string]interface{}{"min": 0, "max": 6}
	var datay []map[string]interface{}
	datay = append(datay, map[string]interface{}{"label": "10K", "value": 1})
	datay = append(datay, map[string]interface{}{"label": "30K", "value": 3})
	datay = append(datay, map[string]interface{}{"label": "50k", "value": 5})
	y_obj["data"] = datay
	sec.Options["y"] = y_obj

	c.AddSection(sec)
	return c
}

func (p *Home) Get(filters map[string][]string) ([]byte, error) {
	var recordset map[string]interface{}
	var data []map[string]interface{}

	// sarebbe da fare una query a db e prendere le informazioni e ridarle così al client
	data = append(data, map[string]interface{}{"label": "CRONUS IT S.p.A.", "value": 45})
	data = append(data, map[string]interface{}{"label": "Erdman Group", "value": 38})
	data = append(data, map[string]interface{}{"label": "Hand Ltd", "value": 17})

	data = append(data, map[string]interface{}{"x": 0, "y": 3})
	data = append(data, map[string]interface{}{"x": 2.6, "y": 2})
	data = append(data, map[string]interface{}{"x": 4.9, "y": 5})
	data = append(data, map[string]interface{}{"x": 6.8, "y": 3.1})
	data = append(data, map[string]interface{}{"x": 8, "y": 4})
	data = append(data, map[string]interface{}{"x": 9.5, "y": 3})
	data = append(data, map[string]interface{}{"x": 11, "y": 4})

	recordset = make(map[string]interface{})
	recordset["recordset"] = data
	return json.Marshal(recordset)
}
