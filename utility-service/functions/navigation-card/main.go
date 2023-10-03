package main

import (
	"context"
	"thesis/lib/component"
	"thesis/lib/utility"
	"thesis/utility-service/models"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type NavigationList struct {
	component.BasePage
	Model *models.Navigation
}

func NewNavigationList() *NavigationList {
	c := &NavigationList{
		Model: models.NewNavigation(),
	}
	c.BasePage = *component.NewBasePage(c, "Navigation Card", component.PageCard, "")
	c.AddModel(c.Model)
	c.AddKey("ID")
	c.AddSection(component.NewSection(
		"Group1",
		component.Group,
		"General",
		component.NewField("ID", "ID", &c.Model.ID),
		component.NewField("PageId", "Page ID", &c.Model.PageId),
		component.NewField("Caption", "Caption", &c.Model.Caption),
		component.NewField("URL", "URL", &c.Model.URL),
		component.NewField("Icon", "Icon", &c.Model.Icon),
		component.NewField("SelectedIcon", "Selected Icon", &c.Model.SelectedIcon),
		component.NewField("Tooltip", "Tooltip", &c.Model.Tooltip),
	))
	return c
}

func main() {
	lambda.Start(
		func(ctx context.Context, request events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
			return utility.AmazonHandler(ctx, request, NewNavigationList())
		})
}
