package domain

import "net/url"

type RedirectionRule struct {
	Key string
	Url *url.URL
}
