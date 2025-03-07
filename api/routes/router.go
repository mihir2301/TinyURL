package routes

import (
	"tinyurl/constants"
	"tinyurl/controllers"
)

var healthcheck = Router{
	Route{"Healthcheck", "GET", constants.HealthCheck, controllers.HealthCheck},
}

var Verification = Router{
	Route{"EmailVerificaton", "POST", constants.VerifyEmail, controllers.VerifyEmail},
	Route{"VerifyOtp", "POST", constants.VerifyOtp, controllers.VerifyOtp},
	Route{"ResendOTP", "POST", constants.ResendOTP, controllers.Resend},
}
