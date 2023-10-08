package main

import (
	"fmt"

	"github.com/go-rod/rod"
)

func GetFromProxyCZ() {

	//page := rod.New().MustConnect().MustPage("https://www.wikipedia.org/")
	page := rod.New().MustConnect().MustPage("http://free-proxy.cz/en/proxylist/country/all/https/uptime/all")

	fmt.Println(page)

	page.MustWaitStable().MustScreenshot("a.png")

	el := page.MustElement("#proxy_list")
	fmt.Println(el)
	fmt.Println(el.MustText())
}
