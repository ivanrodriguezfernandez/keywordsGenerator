/*

Check status code for each url and store urls I could not
open in a dedicate array.
Fetch urls concurrently using goroutines.

*/
package main

import (
    "fmt"
    "net/http"
    "io/ioutil"
)

// ------------------------------


const (
    userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) " +
        "AppleWebKit/537.36 (KHTML, like Gecko) " +
        "Chrome/53.0.2785.143 " +
        "Safari/537.36"
)

func fetchUrl(url string, chFailedUrls chan string , chIsFinished chan bool){
  client := &http.Client{}
  req, _ := http.NewRequest("GET", url , nil)
  req.Header.Set("User-Agent", userAgent)
  response, err := client.Do(req)

  fmt.Println("response status:", response.Status)
  defer response.Body.Close()
  body, err := ioutil.ReadAll(response.Body)
  fmt.Printf("%s", body)

  defer func(){
    chIsFinished <- true
  }()

  if err!= nil || response.StatusCode !=200 {
    chFailedUrls <- url
    return
  }
}

func main() {
    fmt.Println("Google!!");

    urlsList := [1]string{
    "https://www.google.com",
    }

    chFailedUrls:= make(chan string)
    chIsFinished:= make(chan bool)

    for _, url := range urlsList{
        go fetchUrl(url, chFailedUrls, chIsFinished)
    }

    failedUrls := make ([]string, 0)
    for i := 0;i<len(urlsList); {
      select {
        case url := <-chFailedUrls:
          failedUrls = append(failedUrls, url)
        case <- chIsFinished:
          i++
        }
    }

    // Print all urls we could not open:
    fmt.Println("Could not fetch these urls: ", failedUrls)

} // end main
