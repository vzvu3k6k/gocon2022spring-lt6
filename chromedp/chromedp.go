package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/chromedp/chromedp"
)

// prepareContext prepares context for chromedp with timeout and headful option.
func prepareContext() (context.Context, func()) {
	ctx, cancelAlloc := chromedp.NewExecAllocator(
		context.Background(),
		append(chromedp.DefaultExecAllocatorOptions[:], chromedp.Flag("headless", false))...,
	)
	ctx, cancelCdp := chromedp.NewContext(ctx)
	ctx, cancelCtx := context.WithTimeout(ctx, 15*time.Second)

	return ctx, func() {
		cancelCtx()
		cancelCdp()
		cancelAlloc()
	}
}

// createEvent returns chromedp.Tasks to create an event from popup.
func createEvent(title string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Sleep(3 * time.Second),
		chromedp.Query("#DashboardCalendar > *", chromedp.ByQuery),
		chromedp.Click(".EventCreate", chromedp.ByQuery),
		chromedp.SendKeys(`.popup [name="title"]`, title),
		chromedp.Click(".EventCreateSubmit", chromedp.ByQuery),
	}
}

// createEvent returns chromedp.Tasks to set description to event in event edit page.
func setDescription(desc string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Sleep(3 * time.Second),
		chromedp.Click("#FieldDescription", chromedp.ByQuery),
		chromedp.SendKeys(`[name="description_input"]`, desc, chromedp.ByQuery),
		chromedp.Click(`#FieldDescription [type="submit"]`),
	}
}

func main() {
	// chromedpのcontextを用意する
	ctx, cancel := prepareContext()
	defer cancel()

	username := os.Getenv("CONNPASS_ID")
	password := os.Getenv("CONNPASS_PASSWORD")

	// connpassにログインする
	err := chromedp.Run(
		ctx,
		chromedp.Navigate("https://connpass.com/login/"),
		chromedp.SendKeys(`#login_form [name="username"]`, username, chromedp.ByQuery),
		chromedp.SendKeys(`#login_form [name="password"]`, password, chromedp.ByQuery),
		chromedp.Submit("#login_form"),
		chromedp.WaitVisible("body"),
	)
	if err != nil {
		log.Fatal("login failed:\n", err)
	}

	// 「Gopherの会」というイベントを作成する
	if err := chromedp.Run(ctx, createEvent("Gopherの会")); err != nil {
		log.Fatal("createEvent failed:\n", err)
	}

	// イベントの説明文を「参加してね」にする
	if err := chromedp.Run(ctx, setDescription("参加してね")); err != nil {
		log.Fatal("setDescription failed:\n", err)
	}

	if err := chromedp.Run(ctx, chromedp.Sleep(3*time.Second)); err != nil {
		log.Fatal("can't sleep:\n", err)
	}

}
