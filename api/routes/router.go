package routes

import (
	"tinyurl/constants"
	"tinyurl/controllers"
)

var healthcheck = Router{
	Route{"Healthcheck", "GET", constants.HealthCheck, controllers.HealthCheck},
}

var User = Router{
	Route{"EmailVerificaton", "POST", constants.VerifyEmail, controllers.VerifyEmail},
	Route{"VerifyOtp", "POST", constants.VerifyOtp, controllers.VerifyOtp},
	Route{"ResendOTP", "POST", constants.ResendOTP, controllers.Resend},
	Route{"UserLogin", "POST", constants.Login, controllers.UserLogin},
	Route{"UserRegister", "POST", constants.Register, controllers.RegisterUser},
}

var Shortner = Router{
	Route{"UrlShortner", "POST", constants.Shortner, controllers.Shorten},
	Route{"DirectUrl", "GET", constants.DirectUrl, controllers.DirectUrl},
}
