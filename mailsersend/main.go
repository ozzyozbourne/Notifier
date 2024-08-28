package main

import (
	"context"
	"fmt"
	"github.com/mailersend/mailersend-go"
	"time"
)

const APIKEY = "mlsn.59df36400d6c5f1da30ea37c24e8f64056778acee04f88d2ab1261db01cb6325"

func main() {
	ms := mailersend.NewMailersend(APIKEY)
	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	message := ms.Email.NewMessage()

	message.SetFrom(mailersend.Recipient{
		Name:  "Twitcher",
		Email: "twicher@trial-351ndgw957r4zqx8.mlsender.net",
	})

	message.SetRecipients([]mailersend.Recipient{
		{
			Name:  "Suvaid",
			Email: "khansuvaid@yahoo.com",
		},
	})

	message.SetSubject("Testing mailersend API")
	message.SetHTML("I am your father luke")
	message.SetText("whiskey in the ja oor")
	message.SetTags([]string{"foo", "bar"})
	res, err := ms.Email.Send(ctx, message)
	if err != nil {
		fmt.Printf("Error %v\n", err)
	}
	fmt.Printf("%v\n", res)

}
