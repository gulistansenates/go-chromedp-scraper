package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"strings"
	"time"
)

func GetHotelInfo() {

	city := "antalya"
	adultCount := 2
	isAvailable := true
	now := time.Now()
	dateFormat := "02.01.2006"
	after2Day := now.AddDate(0, 0, 2)
	todayDate := now.Format(dateFormat)
	after2Date := after2Day.Format(dateFormat)
	fmt.Println(now.Format(dateFormat))
	fmt.Println(after2Day.Format(dateFormat))

	requestUrl := fmt.Sprintf("https://www.tatilsepeti.com/%s-otelleri?ara=oda:%d;tarih:%s,%s;musait:%t", city, adultCount, todayDate, after2Date, isAvailable)
	fmt.Println(requestUrl)

	res, err := http.Get(requestUrl)
	if err != nil {
		fmt.Println("HTTP isteği başarısız:", err)
		return
	}

	if res.StatusCode != 200 {
		fmt.Println("Hata", res.StatusCode)
		return
	}
	// htmlResponse, _ := io.ReadAll(res.Body)
	// fmt.Println(string(x)[0:100])
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		fmt.Println("Belge oluşturma hatası:", err)
		return
	}

	doc.Find("article").Each(func(i int, selection *goquery.Selection) {
		hotelName := selection.Find("a").First().Text()
		hotelScore := selection.Find("div.score__right").Text()
		hotelScore = strings.ReplaceAll(hotelScore, " ", "")
		hotelScore = strings.ReplaceAll(hotelScore, "\n", "")
		hotelPrice := selection.Find("p.currency::before").Text()
		fmt.Printf("Hotel Name : %s, Score : %s, Price : %s\n", hotelName, hotelScore, hotelPrice)

	})

	/*TODO:
	- https://www.tatilsepeti.com/antalya-otelleri?ara=oda:2;tarih:15.04.2024,20.04.2024;musait:true
	- tatilsepeti.com'a giriş yaptırılacak
	- Şehir, bölge, tatil bölgesi placeholder'ına text yazdırılmaya çalışılacak
	- Giriş tarihi bulunduğumuz gün olarak çıkış tarihi bulunduğumuz günden +2 gün sonrası olarak ayarlanacak
	- yetişkin sayısı 2 olsun, çocuk olmasın(yani default olarak kalsın)
	- Otel Ara butonuna click işlemi
	- tüm bu işlemlerden sonra dönen html'i yazdır
	*/
}
