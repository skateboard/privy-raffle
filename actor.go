package main

import (
	"fmt"
	"os"

	"github.com/data-harvesters/goapify"
)

type input struct {
	*goapify.ProxyConfigurationOptions `json:"proxyConfiguration"`

	BuisnessID string `json:"buisnessID"`
	CampaignID string `json:"campaignID"`
	FormName   string `json:"formName"`

	UseCatchAll   bool   `json:"useCatchAll"`
	CatchAllEmail string `json:"catchAllEmail"`

	Emails []string `json:"emails"`
}

func main() {
	a := goapify.NewActor(
		os.Getenv("APIFY_DEFAULT_KEY_VALUE_STORE_ID"),
		os.Getenv("APIFY_TOKEN"),
		os.Getenv("APIFY_DEFAULT_DATASET_ID"),
	)

	i := new(input)

	err := a.Input(i)
	if err != nil {
		fmt.Printf("failed to decode input: %v\\n", err)
		panic(err)
	}

	if i.ProxyConfigurationOptions != nil {
		err = a.CreateProxyConfiguration(i.ProxyConfigurationOptions)
		if err != nil {
			panic(err)
		}
	}

	p, err := newPrivy(i, a)
	if err != nil {
		fmt.Printf("failed to make privy: %v\\n", err)
		panic(err)
	}

	p.Run()
}
