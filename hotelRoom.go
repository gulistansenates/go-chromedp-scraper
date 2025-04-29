package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func GetHotelRoomInfo(hotelName string, checkInDate, checkOutDate string) {
	adultCount := 2
	isAvailable := true
	dateFormat := "02.01.2006"

	// Tarihlerin doğru formatta olduğundan emin ol
	_, err := time.Parse(dateFormat, checkInDate)
	if err != nil {
		fmt.Println("Giriş tarihini doğru formatta giriniz (GG.AA.YYYY).")
		return
	}

	_, err = time.Parse(dateFormat, checkOutDate)
	if err != nil {
		fmt.Println("Çıkış tarihini doğru formatta giriniz (GG.AA.YYYY).")
		return
	}

	// İstenen tarihler için URL oluştur
	requestUrl := fmt.Sprintf("https://www.tatilsepeti.com/%s?ara=oda:%d;tarih:%s,%s;musait:%t", hotelName, adultCount, checkInDate, checkOutDate, isAvailable)
	fmt.Println("Request URL:", requestUrl)

	res, err := http.Get(requestUrl)
	if err != nil {
		fmt.Println("HTTP isteği başarısız:", err)
		return
	}

	if res.StatusCode != 200 {
		fmt.Println("Hata", res.StatusCode)
		return
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println("Belge oluşturma hatası:", err)
		return
	}

	found := false

	doc.Find("article").Each(func(i int, selection *goquery.Selection) {
		hotel := selection.Find("a").First().Text()
		if strings.Contains(strings.ToLower(hotel), strings.ToLower(hotelName)) {
			found = true
			roomType := selection.Find("h1.Header--Title").Text()
			price := selection.Find("span.Prices--Price").Text()
			price = strings.TrimSpace(price)
			if price == "" {
				price = "Bilgi Bulunamadı"
			}
			fmt.Printf("Otel: %s\nOda Türü: %s\nİndirimli Fiyat: %s\n", hotel, roomType, price)

			availability := selection.Find("p").Text()
			if strings.Contains(availability, "Bu otelde") {
				fmt.Println(availability)
			}
		}
	})

	if !found {
		// Belirtilen otel bulunamadıysa veya belirtilen tarihlerde müsait oda bulunmuyorsa, HTML'den mesajı alıp yazdır
		doc.Find("p").Each(func(i int, selection *goquery.Selection) {
			if strings.Contains(selection.Text(), "Bu otelde") {
				fmt.Println(selection.Text())
				found = true
			}
		})
	}

	if !found {
		fmt.Println("Belirtilen otel bulunamadı veya belirtilen tarihlerde müsait oda bulunmamaktadır.")
	}
}
