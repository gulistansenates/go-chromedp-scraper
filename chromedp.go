package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

func GetInfo() {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
		chromedp.Flag("disable-gpu", false),
		chromedp.Flag("enable-automation", false),
		chromedp.Flag("disable-extensions", false),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()

	city := "Ankara"
	now := time.Now()
	after2Day := now.AddDate(0, 0, 2)

	if err := chromedp.Run(ctx,
		chromedp.Navigate("https://www.tatilsepeti.com/"),
		chromedp.WaitVisible(`#bolge`),
		chromedp.SetValue(`#bolge`, city),
		chromedp.SetValue(`#giris-tarihi`, now.Format("02.01.2006")),
		chromedp.SetValue(`#cikis-tarihi`, after2Day.Format("02.01.2006")),
		chromedp.Click(`#searchBtn`),
		chromedp.Click(fmt.Sprintf(`a[href="%s"].advisor-link`, strings.ToLower("/"+city+"-otelleri")), chromedp.NodeVisible),
		chromedp.WaitVisible(`article`),
		//chromedp.WaitVisible(`Guliyi bekle`),
	); err != nil {
		fmt.Println("Hata:", err)
		return
	}

	htmlContent := getHTML(ctx)
	fmt.Println("HTML içeriği:")
	fmt.Println(htmlContent)
}

func getHTML(ctx context.Context) string {
	var htmlContent string
	if err := chromedp.Run(ctx, chromedp.OuterHTML(`html`, &htmlContent)); err != nil {
		fmt.Println("HTML içeriği alınamadı:", err)
		return ""
	}
	fmt.Println("HTML CONTENT: ", htmlContent)
	return htmlContent
}
