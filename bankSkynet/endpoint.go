package bankSkynet

import (
	"context"
	"github.com/go-kit/kit/endpoint"
	"github.com/weAutomateEverything/go2hal/gokit"
)

type SkynetRebuildRequest struct {
	NodeName string `json:"Nodename"`
	User     string `json:"User"`
}

func makeSkynetRebuildEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(SkynetRebuildRequest)
		s.RecreateNode(ctx,gokit.GetChatId(ctx), req.NodeName, req.User)
		return nil, err
	}
}

func makeSkynetAlertEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(string)
		s.sendSkynetAlert(ctx,gokit.GetChatId(ctx), req)
		return nil, nil
	}
}
