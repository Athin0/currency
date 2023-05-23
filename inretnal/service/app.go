package service

import (
	"bytes"
	"context"
	"currency/inretnal/adapters"
	"currency/inretnal/model"
	"encoding/xml"
	"errors"
	"fmt"
	"golang.org/x/text/encoding/charmap"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type App struct {
	db *adapters.PostgresDB
}

func NewApp(db *adapters.PostgresDB) *App {
	return &App{db: db}
}

func (a *App) GetCurrencies(date string) (*model.Result, error) {
	var result model.Result
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://www.cbr.ru/scripts/XML_daily_eng.asp?date_req="+date, nil)
	req.Header.Add("Accept", `text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8`)
	req.Header.Add("User-Agent", `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_7_5) AppleWebKit/537.11 (KHTML, like Gecko) Chrome/23.0.1271.64 Safari/537.11`)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("bad http response status:" + resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	decoder := xml.NewDecoder(bytes.NewReader(body))
	decoder.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		switch charset {
		case "windows-1251":
			return charmap.Windows1251.NewDecoder().Reader(input), nil
		default:
			return nil, fmt.Errorf("unknown charset: %s", charset)
		}
	}
	err = decoder.Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (a *App) LastDaysAddToDB(ctx context.Context, nDays int) error {
	lastDate, err := a.db.GetLastDate(ctx)
	if err != nil {
		log.Println(err)
	}
	end := time.Now()
	start := end.AddDate(0, 0, -1*nDays)
	if lastDate.After(start) {
		start = lastDate.AddDate(0, 0, 1)
	}
	for d := start; d.After(end) == false; d = d.AddDate(0, 0, 1) {
		date := d.Format("02/01/2006")
		result, err := a.GetCurrencies(date)
		if err != nil {
			return err
		}
		err = a.db.InsertCurrencies(ctx, result.Currencies, d)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *App) GetAvg(ctx context.Context) ([]model.ResponseCurrencyAvg, error) {
	return a.db.GetAvg(ctx)
}

func (a *App) GetMax(ctx context.Context) (*model.ResponseCurrency, error) {
	return a.db.GetMax(ctx)
}

func (a *App) GetMin(ctx context.Context) (*model.ResponseCurrency, error) {
	return a.db.GetMin(ctx)
}
