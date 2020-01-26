package fns

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Example: t=20200115T2110&s=1030.00&fn=9251440300046840&i=29414&fp=1250830908&n=1
func ParseQrStr(str string) (fn string, opType, fd, fpd int, date time.Time, sum float32, err error) {
	params := strings.Split(str, "&")
	if len(params) != 6 {
		err = fmt.Errorf("unexpected params count %v. 6 expected", len(params))
		return
	}

	for _, param := range params {
		kv := strings.Split(param, "=")
		if len(kv) != 2 {
			err = fmt.Errorf("bad kv at %v", param)
			return
		}

		switch kv[0] {
		case "t":
			date, err = time.Parse("20060102T1504", kv[1])
			if err != nil {
				date, err = time.Parse("20060102T150405", kv[1])
				if err != nil {
					err = fmt.Errorf("bad date %v", kv[1])
					return
				}
			}
		case "s":
			s, e := strconv.ParseFloat(kv[1], 32)
			if e != nil {
				err = fmt.Errorf("bad sum %v", kv[1])
				return
			}
			sum = float32(s)
		case "fn":
			fn = kv[1]
		case "i":
			fd, err = strconv.Atoi(kv[1])
			if err != nil {
				err = fmt.Errorf("bad fd %v", kv[1])
				return
			}
		case "fp":
			fpd, err = strconv.Atoi(kv[1])
			if err != nil {
				err = fmt.Errorf("bad fpd %v", kv[1])
				return
			}
		case "n":
			opType, err = strconv.Atoi(kv[1])
			if err != nil {
				err = fmt.Errorf("bad opType %v", kv[1])
				return
			}
		default:
			err = fmt.Errorf("unexpected key %v", kv[0])
			return
		}
	}

	return
}

// CheckReceipt returns nil if receipt exists
func CheckReceipt(ctx context.Context, fn string, opType, fd, fpd int, date time.Time, sum float32) error {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("https://proverkacheka.nalog.ru:9999/v1/ofds/*/inns/*/fss/%v/operations/%v/tickets/%v?fiscalSign=%v&date=%v&sum=%v", fn, opType, fd, fpd, date.Format("2006-01-02T15:04:05"), sum*100), nil)
	if err != nil {
		return err
	}

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("unexpected code %v", resp.StatusCode)
	}

	return nil
}

// GetReceipt returns full receipt info
func GetReceipt(ctx context.Context, phone, pass string, fn string, fd, fpd int) (Receipt, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("https://proverkacheka.nalog.ru:9999/v1/inns/*/kkts/*/fss/%v/tickets/%v?fiscalSign=%v&sendToEmail=no", fn, fd, fpd), nil)
	if err != nil {
		return Receipt{}, err
	}

	req.Header.Add("device-id", "")
	req.Header.Add("device-os", "")
	req.SetBasicAuth(phone, pass)

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return Receipt{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Receipt{}, fmt.Errorf("unexpected code %v", resp.StatusCode)
	}

	dec := json.NewDecoder(resp.Body)
	var prj receiptPrj
	if err := dec.Decode(&prj); err != nil {
		return Receipt{}, fmt.Errorf("unable to decode response: %w", err)
	}

	return NewReceipt(prj)
}
