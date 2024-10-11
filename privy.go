package main

import (
	"io"
	"net/http"

	"github.com/data-harvesters/goapify"
	goapifytls "github.com/data-harvesters/goapify-tls"
	"github.com/skateboard/ajson"
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
