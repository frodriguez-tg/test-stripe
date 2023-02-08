package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/charge"
	"github.com/stripe/stripe-go/v74/checkout/session"
	"github.com/stripe/stripe-go/v74/customer"
	"github.com/stripe/stripe-go/v74/token"
)

const customerID = ""

func main() {
	// Set your secret key. Remember to switch to your live secret key in production.
	// See your keys here: https://dashboard.stripe.com/apikeys
	stripe.Key = ""

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Static("/", "")

	e.POST("/create-checkout-session", createCheckoutSession)
	e.POST("/charge", func(c echo.Context) error {
		// create a card token
		tokenParams := &stripe.TokenParams{
			Card: &stripe.CardParams{
				Number:   stripe.String("4242424242424242"),
				ExpMonth: stripe.String("09"),
				ExpYear:  stripe.String("2026"),
				CVC:      stripe.String("333"),
			},
		}
		token, err := token.New(tokenParams)
		if err != nil {
			return err
		}
		fmt.Println(token.ID)
		tokenID := stripe.String(token.ID)

		cus, err := customer.Update(customerID, &stripe.CustomerParams{
			Source: tokenID,
		})
		if err != nil {
			return err
		}
		customerId := stripe.String(cus.ID)

		// Charge the Customer instead of the card:
		chargeParams := &stripe.ChargeParams{
			Amount:   stripe.Int64(10000),
			Currency: stripe.String(string(stripe.CurrencyMXN)),
			Customer: customerId,
		}

		ch, err := charge.New(chargeParams)
		if err != nil {
			return err
		}

		fmt.Println(ch)

		// When it's time to charge the customer again, retrieve the customer ID.
		params := &stripe.ChargeParams{
			Amount:   stripe.Int64(15000), // $150.00 this time
			Currency: stripe.String(string(stripe.CurrencyMXN)),
			Customer: customerId,
		}
		ch, err = charge.New(params)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, "success")
	})

	e.Logger.Fatal(e.Start("localhost:4242"))
}

func createCheckoutSession(c echo.Context) (err error) {
	params := &stripe.CheckoutSessionParams{
		PaymentIntentData: &stripe.CheckoutSessionPaymentIntentDataParams{
			SetupFutureUsage: stripe.String("off_session"),
		},
		Customer: stripe.String(customerID),
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
