package main

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/checkout/session"
)

func main() {
	// Set your secret key. Remember to switch to your live secret key in production.
	// See your keys here: https://dashboard.stripe.com/apikeys
	stripe.Key = ""

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Static("/", "")

	e.POST("/create-checkout-session", createCheckoutSession)

	e.Logger.Fatal(e.Start("localhost:4242"))
}

func createCheckoutSession(c echo.Context) (err error) {
	params := &stripe.CheckoutSessionParams{
		PaymentIntentData: &stripe.CheckoutSessionPaymentIntentDataParams{
			SetupFutureUsage: stripe.String("off_session"),
		},
		Customer: stripe.String(""),
		Mode:     stripe.String(string(stripe.CheckoutSessionModePayment)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String("MXN"),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String("Another T-shirt"),
					},
					UnitAmount: stripe.Int64(10000),
				},
				Quantity: stripe.Int64(1),
			},
		},
		SuccessURL: stripe.String("http://localhost:4242/success.html"),
		CancelURL:  stripe.String("http://localhost:4242/cancel"),
	}

	s, _ := session.New(params)

	if err != nil {
		return err
	}

	return c.Redirect(http.StatusSeeOther, s.URL)
}
