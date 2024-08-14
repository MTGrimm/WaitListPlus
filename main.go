package main

import (
  "fmt"
  "time"
  "regexp"
  "math/rand/v2"
  "net/smtp"
  "github.com/go-rod/rod"
)

func sendEmail(message string) {
  auth := smtp.PlainAuth(
    "",
    "aryan.timilsina195@gmail.com",
    "sracqawdqrrwnooc",
    "smtp.gmail.com",
  )

  err := smtp.SendMail(
    "smtp.gmail.com:587",
    auth,
    "aryan.timilsina195@gmail.com",
    []string{"aryan.timilsina195@gmail.com"},
    []byte(message),
  )

  if err != nil {
    fmt.Println(err)
  }
}

func main() {
  courses := make([]string, 0)
  fmt.Print("PLease enter in courses in the following format 'NAME-101', and exit to exit: ")
  for true {
    var input string;
    fmt.Scanf("%s", &input)
    if input == "exit" {
      break
    }
    courses = append(courses, input)
  }
  
  baseUrl := "https://register.beartracks.ualberta.ca/criteria.jsp?access=0&lang=en&page=results&advice=0&legend=1&term=1890&sort=none&filters=liiiiiiiii&bbs=&ds=&cams=UOFABiOFF_UOFABiMAIN&locs=any&isrts=any&ses=any&pl=&pac=1"
  for i, course := range courses {
    baseUrl += fmt.Sprintf("&course_%d_0=%s", i, course)
  }
  fmt.Println(baseUrl)
  
  browser := rod.New().MustConnect()
	defer browser.MustClose()
  
  firstLoaded := false
  previousStatuses := make([]bool, len(courses))
  for (true) {
    page := browser.MustPage(baseUrl).MustWaitStable()
    time.Sleep(3 * time.Second)

    matches, _ := page.Elements(".cbox-unlocked")

    r, _ := regexp.Compile("All classes are full")

    visibleIndex := 0
    for _, element := range matches {
      if element.MustVisible() {
        warning, err := element.Element(".cbox-warnings")
        isFull := false
        if err == nil {
          html, _ := warning.HTML()
          if r.Match([]byte(html)) {
            fmt.Printf("%s is full\n", courses[visibleIndex])
            isFull = true
          }
        }
        if !isFull {
          fmt.Printf("%s is not full\n", courses[visibleIndex])
        }
        if firstLoaded {
          fmt.Println(previousStatuses)
          if previousStatuses[visibleIndex] != isFull {
            if (isFull) {
              sendEmail(fmt.Sprintf("%s is now full\n", courses[visibleIndex]))
            } else {
              sendEmail(fmt.Sprintf("%s is now NOT FULL\n", courses[visibleIndex]))
            }
          }
        }
        previousStatuses[visibleIndex] = isFull
        visibleIndex += 1
      }
    }
    firstLoaded = true
    time.Sleep(time.Duration((rand.IntN(240) + 120)) * time.Second)
  }
}
