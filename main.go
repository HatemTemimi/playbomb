package main

import (
	"bufio"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

func main(){

  var allproxies []string

  //get all http proxies from TheSpeedX List
  resp, err := http.Get("https://raw.githubusercontent.com/TheSpeedX/SOCKS-List/master/http.txt")
  if err != nil {
    panic(err)
  }
  //create new http client and scan previous results
  client := &http.Client{}
  scanner := bufio.NewScanner(resp.Body)

  //for each proxy found in results, we will perform a custom http get request with it
  for scanner.Scan() {
      // Define the proxy URL
      withHttp := "http://" + scanner.Text()
      proxyURL, err := url.Parse(withHttp)
      if err != nil {
        log.Fatal(err)
      }

      // Create a new Transport with the proxy
      transport := &http.Transport{
        Proxy: http.ProxyURL(proxyURL),
        ResponseHeaderTimeout: 10 * time.Second,
      }

      // set the custom Transport with proxy
      client.Transport = transport
      httpPoke, _ := http.NewRequest("GET", "http://www.soundcloud.com", nil)
      httpsPoke, _ := http.NewRequest("GET", "https://www.soundcloud.com", nil) 

      // fool the server into a new agent(lowers forbidden chance)
      httpPoke.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36")
      httpsPoke.Header.Set("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36")


      // perform requests
      httpResp, httpErr  := client.Do(httpPoke)
      httpsResp, httpsErr  := client.Do(httpsPoke)

      if (httpErr != nil && httpsErr != nil){
        //proxy did not respond
        log.Println(withHttp, "is Down.")

      } else {
        if(httpsResp != nil && httpResp != nil){
          //proxy did respond
          log.Println(withHttp, " is Up.")

          if (httpResp.StatusCode != 403 && httpsResp.StatusCode != 403){
            //proxy did respond and is not blocked
            // append it to array and write array to file(writes a new array each time)
            log.Println(withHttp, " is Up & Fresh.")
            allproxies = append(allproxies, withHttp)
            proxyJson, _ := json.Marshal(allproxies)
            os.WriteFile("proxy.json", proxyJson, 0660)
          }
        }
      } 
  }

  //panic within scanner
  if err := scanner.Err(); err != nil {
      panic(err)
  }

}
