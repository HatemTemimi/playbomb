package main

import (
	"bufio"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"
)

func checkProxy(client *http.Client, withHttp string, workingProxies *[]string, wg *sync.WaitGroup) {
	defer wg.Done()
	proxyURL, err := url.Parse(withHttp)
	if err != nil {
		log.Fatal(err)
	}

	transport := &http.Transport{
		Proxy:                 http.ProxyURL(proxyURL),
		ResponseHeaderTimeout: 20 * time.Second,
		DialContext: (&net.Dialer{
			Timeout:   90 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	client.Transport = transport

	httpPoke, _ := http.NewRequest("GET", "http://www.google.com", nil)
	httpsPoke, _ := http.NewRequest("GET", "https://www.google.com", nil)

	httpPoke.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36")
	httpsPoke.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36")

	httpResp, httpErr := client.Do(httpPoke)
	httpsResp, httpsErr := client.Do(httpsPoke)

	if httpErr != nil && httpsErr != nil {
		log.Println(httpErr)
		log.Println(withHttp, "is Down.")
	} else {
		if httpsResp != nil && httpResp != nil {
			log.Println(withHttp, " is Up.")
			if httpResp.StatusCode != 403 && httpsResp.StatusCode != 403 {
				log.Println(withHttp, " is Up & Fresh.")
				*workingProxies = append(*workingProxies, withHttp)
				proxyJson, _ := json.Marshal(*workingProxies)
				os.WriteFile("proxies.json", proxyJson, 0660)
			}
		}
	}
}

func GetFromSpeedx() {

	var allproxies []string
	var workingProxies []string

	resp1, err := http.Get("https://api.proxyscrape.com/v2/?request=displayproxies&protocol=http&timeout=5000&country=all&ssl=all&anonymity=all")
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(resp1.Body)

	for scanner.Scan() {
		allproxies = append(allproxies, scanner.Text())
	}

	resp, err := http.Get("https://raw.githubusercontent.com/TheSpeedX/SOCKS-List/master/http.txt")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	scanner = bufio.NewScanner(resp.Body)

	client := &http.Client{}

	var wg sync.WaitGroup

	for scanner.Scan() {
		allproxies = append(allproxies, scanner.Text())
	}

	log.Println("checking: ", len(allproxies), "proxies..")

	for _, e := range allproxies {
		wg.Add(1)
		go checkProxy(client, "http://"+e, &workingProxies, &wg)
	}

	wg.Wait()

	if err := scanner.Err(); err != nil {
		panic(err)
	}

}
