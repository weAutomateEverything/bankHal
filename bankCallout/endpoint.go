package bankCallout

import (
	"github.com/go-kit/kit/endpoint"
	"context"
	"github.com/go-kit/kit/auth/jwt"
	"github.com/weAutomateEverything/go2hal/telegram"
)

func makeGetCalloutGroupEndpoint(s Service) endpoint.Endpoint{
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		claim := ctx.Value(jwt.JWTClaimsContextKey).(*telegram.CustomClaims)
		group,err  := s.getGroup(ctx,claim.RoomToken)
		if err != nil {
			return nil, err
		}
		return getCalloutGroupResponse{
			Group:group,
		}, nil

	}
}

func makeSetCalloutGroupEndpoint(s Service) endpoint.Endpoint{
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		claim := ctx.Value(jwt.JWTClaimsContextKey).(*telegram.CustomClaims)
		req := request.(setCalloutRequest)
		name, number, err := s.setGroup(ctx,claim.RoomToken,req.Group)
		if err != nil {
			return nil, err
		}

		return setCalloutResponse{
			Name:name,
			PhoneNumber:number,
		}, nil
	}
}

type getCalloutGroupResponse struct {
	Group string
}

type setCalloutRequest struct {
	Group string
}

type setCalloutResponse struct {
	Name string
	PhoneNumber string
}