package pika

import "github.com/knvi/pika"

var prefixes = []pika.PikaPrefixDefinition{
	{
		Prefix:      "user",
		Description: "User prefix",
		Secure:      false,
	},
	{
		Prefix:      "jti",
		Description: "JWT ID Prefix. For use ONLY as jti field.",
		Secure:      true,
	},
}

var P = pika.NewPika(prefixes, pika.PikaInitOptions{
	NodeID:           622,
	DisableLowercase: true,
})

func Get() *pika.Pika {
	return P
}

func Gen(prefix string) string {
	return P.Gen(prefix)
}
