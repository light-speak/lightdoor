package main

import (
	token "github.com/light-speak/lightdoor/security/kitex_gen/token/securityservice"
	"log"
)

func main() {
	svr := token.NewServer(new(SecurityServiceImpl))

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
