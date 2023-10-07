package main

import (
	"bufio"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"
)

func checkProxy(client *http.Client, withHttp string, allproxies *[]string, wg *sync.WaitGroup) {
	defer wg.Done()

	proxyURL, err := url.Parse(withHttp)
	if err != nil {
		log.Fatal(err)
	}

	transport := &http.Transport{
		Proxy:                 http.ProxyURL(proxyURL),
		ResponseHeaderTimeout: 10 * time.Second,
	}

	client.Transport = transport

	httpPoke, _ := http.NewRequest("GET", "http://www.soundcloud.com", nil)
	httpsPoke, _ := http.NewRequest("GET", "https://www.soundcloud.com", nil)

	httpPoke.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36")
	httpsPoke.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36")

	httpResp, httpErr := client.Do(httpPoke)
	httpsResp, httpsErr := client.Do(httpsPoke)

	if httpErr != nil && httpsErr != nil {
		log.Println(withHttp, "is Down.")
	} else {
		if httpsResp != nil && httpResp != nil {
			log.Println(withHttp, " is Up.")
			if httpResp.StatusCode != 403 && httpsResp.StatusCode != 403 {
				log.Println(withHttp, " is Up & Fresh.")
				*allproxies = append(*allproxies, withHttp)
				proxyJson, _ := json.Marshal(*allproxies)
				os.WriteFile("proxy.json", proxyJson, 0660)
			}
		}
	}
}

func main() {

	var allproxies []string

	resp, err := http.Get("https://raw.githubusercontent.com/TheSpeedX/SOCKS-List/master/http.txt")
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)

	client := &http.Client{}

	var wg sync.WaitGroup

	for scanner.Scan() {
		wg.Add(1)
		go checkProxy(client, "http://"+scanner.Text(), &allproxies, &wg)
	}

	wg.Wait()

	if err := scanner.Err(); err != nil {
		panic(err)
	}

}
