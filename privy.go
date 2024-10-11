package main

import (
	"fmt"
	"io"
	"strings"

	http "github.com/bogdanfinn/fhttp"
	"github.com/data-harvesters/goapify"
	goapifytls "github.com/data-harvesters/goapify-tls"
	"github.com/skateboard/ajson"
	"github.com/spf13/cast"
)

type privy struct {
	actor *goapify.Actor
	input *input

	client *goapifytls.TlsClient
}

func newPrivy(input *input, actor *goapify.Actor) (*privy, error) {
	tlsClient, err := goapifytls.NewTlsClient(actor, goapifytls.DefaultOptions())
	if err != nil {
		return nil, err
	}

	return &privy{
		actor:  actor,
		input:  input,
		client: tlsClient,
	}, nil
}

func (p *privy) Run() {
	if p.input.UseCatchAll {
		fmt.Printf("using catch-all email %s for %d raffle-signups\n", p.input.CatchAllEmail, p.input.CatchAllLimit)

		for i := 0; i < p.input.CatchAllLimit; i++ {
			first, last, err := getRandomName()
			if err != nil {
				continue
			}
			email := fmt.Sprintf("%s.%s1%s", first, last, p.input.CatchAllEmail)
			fmt.Printf("signing-up email: %s\n", email)

			err = p.signUpEmail(first, last, email)
			if err != nil {
				fmt.Printf("%s: failed to signup email: %v\n", email, err)
				continue
			}

			err = p.actor.Output(map[string]string{
				"email": email,
			})
			if err != nil {
				fmt.Printf("failed to send output: %v\n", err)
			}

			fmt.Printf("signed-up email: %s\n", email)
		}
		fmt.Printf("succesfully signed-up %d emails\n", p.input.CatchAllLimit)
		return
	}

	if p.input.Emails == nil {
		fmt.Printf("you need to provide emails if you arent going to use a catch-all!\n")
		return
	}
	emails := *p.input.Emails

	if len(emails) == 0 {
		fmt.Printf("you need to provide emails if you arent going to use a catch-all!\n")
		return
	}

	for _, email := range emails {
		first, last, err := getRandomName()
		if err != nil {
			continue
		}
		fmt.Printf("signing-up email: %s\n", email)

		err = p.signUpEmail(first, last, email)
		if err != nil {
			fmt.Printf("%s: failed to signup email: %v\n", email, err)
			continue
		}

		err = p.actor.Output(map[string]string{
			"email": email,
		})
		if err != nil {
			fmt.Printf("failed to send output: %v\n", err)
		}

		fmt.Printf("signed-up email: %s\n", email)
	}

	fmt.Printf("succesfully signed-up %d emails\n", len(emails))
}

func (p *privy) signUpEmail(first, last, email string) error {
	url := fmt.Sprintf("https://api.privy.com/businesses/%v/campaigns/%v/transactions", p.input.BuisnessID,
		p.input.CampaignID)

	payload := fmt.Sprintf(`
{
    "_method": "post",
    "customer_attributes": {
        "first_name": "%v",
        "last_name": "%v",
        "email": "%v",
        "country_dropdown": "US",
        "agree_to_terms": [
            "agree to terms and conditions"
        ],
        "form_name": "%v",
        "url": "url",
        "source": "Privy",
        "country": "{{country_code}}"
    },
    "utf8": "✓",
    "captcha": "",
    "context": "embedded_form",
    "display_id": %v,
    "browsing_data": {
        "sessions_count": 0,
        "pageviews_all_time": 0,
        "pageviews_this_session": 0,
        "utm_medium": "unknown"
    }
}`, first, last, email, p.input.FormName, cast.ToInt(p.input.DisplayID))

	req, err := http.NewRequest("POST", url, strings.NewReader(payload))
	if err != nil {
		return err
	}

	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/112.0")
	req.Header.Add("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Add("Accept-Language", "en-US,en;q=0.5")
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req.Header.Add("Referer", "https://promotions.lpage.co/")
	req.Header.Add("Content-Type", "text/plain")
	req.Header.Add("Origin", "https://promotions.lpage.co")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Sec-Fetch-Dest", "empty")
	req.Header.Add("Sec-Fetch-Mode", "cors")
	req.Header.Add("Sec-Fetch-Site", "cross-site")
	req.Header.Add("Pragma", "no-cache")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("TE", "trailers")

	res, err := p.client.ProxiedClient().Do(req)
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("failed to sign-up: %d", res.StatusCode)
	}

	return nil
}

func getRandomName() (string, string, error) {
	res, err := http.Get("http://random.disqualifier.me/profile/G5d5HuSb0snwon77")
	if err != nil {
		return "", "", err
	}
	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return "", "", err
	}
	j := ajson.Parse(string(b))

	return j.Get("first").String(), j.Get("last").String(), nil
}
