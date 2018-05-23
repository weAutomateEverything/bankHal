package bankCallout
import (

	kitlog "github.com/go-kit/kit/log"
	kithttp "github.com/go-kit/kit/transport/http"
	gokitjwt "github.com/go-kit/kit/auth/jwt"

	"github.com/gorilla/mux"
	"github.com/weAutomateEverything/go2hal/gokit"
	"github.com/weAutomateEverything/go2hal/machineLearning"
	"net/http"
	"github.com/dgrijalva/jwt-go"
	"context"
	"encoding/json"
	"github.com/weAutomateEverything/go2hal/telegram"
)

/*
MakeHandler returns a HTTP Restul endpoint to handle user requests
*/
func MakeHandler(service Service, logger kitlog.Logger, ml machineLearning.Service) http.Handler {
	opts := gokit.GetServerOpts(logger, ml)

	getCalloutGroup := kithttp.NewServer(gokitjwt.NewParser(gokit.GetJWTKeys(),jwt.SigningMethodHS256,
		telegram.CustomClaimFactory)( makeGetCalloutGroupEndpoint(service)), gokit.DecodeString, gokit.EncodeResponse,
		opts...)

	setCalloutGroup := kithttp.NewServer(gokitjwt.NewParser(gokit.GetJWTKeys(),jwt.SigningMethodHS256,
		telegram.CustomClaimFactory)( makeSetCalloutGroupEndpoint(service)),decodeSetCalloutRequest, gokit.EncodeResponse,
		opts...)


	r := mux.NewRouter()

	r.Handle("/bankcallout/firstcall", setCalloutGroup).Methods("POST")
	r.Handle("/bankcallout/firstcall", getCalloutGroup).Methods("GET")
	//r.Handle("/httpendpoints/{id}", authpoll).Methods("DELETE")

	return r

}

func decodeSetCalloutRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	v := setCalloutRequest{}
	err := json.NewDecoder(r.Body).Decode(&v)
	return v, err
}
