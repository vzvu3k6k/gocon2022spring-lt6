package main

import (
	"os"
	"time"

	"github.com/go-rod/rod"
)

func main() {
	// Chromeのウィンドウを用意してログインページを開く
	page := rod.New().MustConnect().MustPage("https://connpass.com/login/")
	page.MustWaitLoad()

	// ログインする
	loginForm := page.MustElement("#login_form")
	loginForm.MustElement(`[name="username"]`).MustInput(os.Getenv("CONNPASS_ID"))
	loginForm.MustElement(`[name="password"]`).MustInput(os.Getenv("CONNPASS_PASSWORD"))
	loginForm.MustElement(`[type="submit"]`).MustClick()
	page.MustWaitLoad()

	// イベント作成ボタンを押して、イベント名を入力して作成ボタンを押す
	page.WaitElementsMoreThan("#DashboardCalendar > *", 0)
	page.MustElement(".EventCreate").MustClick()
	popup := page.MustElement(".popup")
	popup.MustElement(`[name="title"]`).MustInput("Gopherの会")
	popup.MustElement(".EventCreateSubmit").MustClick()
	// （イベントの編集ページに遷移する）

	// 説明文を入力して保存ボタンを押す
	page.WaitElementsMoreThan(".JoinOptions > *", 0)
	fieldDescripion := page.MustElement("#FieldDescription")
	fieldDescripion.MustClick()
	fieldDescripion.MustElement(`[name="description_input"]`).MustInput("みんな来てくれ")
	fieldDescripion.MustElement(`[type="submit"]`).MustClick()
	page.WaitRequestIdle(1*time.Second, nil, nil)
}
