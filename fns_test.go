package fns

import (
	"context"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

const (
	testFn    = "9251440300046840"
	testFd    = 29414
	testFpd   = 1250830908
	testSum   = float32(1030)
	testQrStr = "t=20200115T2110&s=1030.00&fn=9251440300046840&i=29414&fp=1250830908&n=1"
)

var testDate = time.Date(2020, 1, 15, 21, 10, 0, 0, time.UTC)

func TestGetReceipt(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	receipt, err := GetReceipt(ctx, os.Getenv("fns_phone"), os.Getenv("fns_pass"), testFn, testFd, testFpd)
	require.NoError(t, err)

	require.Equal(t, float32(0), receipt.PrepaymentSum)
	require.Equal(t, testFd, receipt.FiscalDocumentNumber)
	require.Equal(t, float32(0), receipt.CounterSubmissionSum)
	require.Equal(t, 1, receipt.TaxationType)
	require.Equal(t, 3, receipt.Code)
	require.Equal(t, os.Getenv("fns.fn"), receipt.FiscalDriveNumber)
	require.Equal(t, 3, receipt.Code)
	require.Equal(t, float32(171.66), receipt.Nds18)
	require.Equal(t, testFpd, receipt.FiscalSign)
	require.Equal(t, "КАССИР: Вербицкая Анаст", receipt.Operator)
	require.Equal(t, 1, receipt.OperationType)
	require.Equal(t, 93, receipt.RequestNumber)
	require.Equal(t, float32(0), receipt.PostpaymentSum)
	require.Equal(t, 179, receipt.ShiftNumber)
	require.Equal(t, "            ", receipt.OperatorInn)
	require.Equal(t, testSum, receipt.EcashTotalSum)
	require.Equal(t, 2, receipt.ProtocolVersion)
	require.Equal(t, testDate, receipt.DateTime)
	require.Equal(t, testSum, receipt.TotalSum)
	require.Equal(t, os.Getenv("fns.inn"), receipt.UserInn)
	require.Equal(t, float32(0), receipt.CashTotalSum)

	require.Len(t, receipt.Items, 3)
	require.Equal(t, 1, receipt.Items[0].CalculationSubjectSign)
	require.Equal(t, 1, receipt.Items[0].Quantity)
	require.Equal(t, 1, receipt.Items[0].NdsRate)
	require.Equal(t, float32(104.50), receipt.Items[0].NdsSum)
	require.Equal(t, "A-936-. 0  Ролик для пресса Double-wheeled exerciser . р.0", receipt.Items[0].Name)
	require.Equal(t, float32(627.00), receipt.Items[0].Sum)
	require.Equal(t, float32(627.00), receipt.Items[0].Price)
}

func TestGetReceiptReturnErrOnUnableToCreateReq(t *testing.T) {
	_, err := GetReceipt(nil, os.Getenv("fns_phone"), os.Getenv("fns_pass"), testFn, testFd, testFpd)
	require.EqualError(t, err, "net/http: nil Context")
}

func TestNewReceiptReturnErrOnBadData(t *testing.T) {
	prj := receiptPrj{}
	prj.Document.Receipt.DateTime = "bad_date"
	_, err := NewReceipt(prj)
	require.EqualError(t, err, "bad date bad_date")
}

func TestGetReceiptReturnErrOnBadData(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := GetReceipt(ctx, os.Getenv("fns_phone"), os.Getenv("fns_pass"), "bad", testFd, testFpd)
	require.EqualError(t, err, "unexpected code 451")
}

func TestCheckReceipt(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := CheckReceipt(ctx, testFn, 1, testFd, testFpd, testDate, testSum)
	require.NoError(t, err)
}

func TestCheckReceiptReturnErrOnUnableToCreateReq(t *testing.T) {
	err := CheckReceipt(nil, testFn, 1, testFd, testFpd, testDate, testSum)
	require.EqualError(t, err, "net/http: nil Context")
}

func TestCheckReceiptReturnErrOnBadSum(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := CheckReceipt(ctx, testFn, 1, testFd, testFpd, testDate, 1000.00)
	require.EqualError(t, err, "unexpected code 406")
}

func TestParseQrStr(t *testing.T) {
	fn, tp, fd, fpd, date, sum, err := ParseQrStr(testQrStr)
	require.NoError(t, err)
	require.Equal(t, testFn, fn)
	require.Equal(t, 1, tp)
	require.Equal(t, testFd, fd)
	require.Equal(t, testFpd, fpd)
	require.Equal(t, testDate, date)
	require.Equal(t, testSum, sum)

	_, _, _, _, _, _, err = ParseQrStr("s=1030.00&t=20200115T2110&fn=1234567890123456&i=12345&fp=1234567890&n=1")
	require.NoError(t, err)

	// time with seconds
	fn, tp, fd, fpd, date, sum, err = ParseQrStr("t=20200115T211000&s=1030.00&fn=1234567890123456&i=12345&fp=1234567890&n=1")
	require.NoError(t, err)
}

func TestParseQrStrReturnErrOnBadParamsCount(t *testing.T) {
	_, _, _, _, _, _, err := ParseQrStr("t=20200115T2110&s=1030.00&fn=1234567890123456&i=12345&fp=1234567890&n=1&a=2")
	require.EqualError(t, err, "unexpected params count 7. 6 expected")
}

func TestParseQrStrReturnErrOnBadKv(t *testing.T) {
	_, _, _, _, _, _, err := ParseQrStr("t=20200115T2110&s=1030.00=20&fn=1234567890123456&i=12345&fp=1234567890&n=1")
	require.EqualError(t, err, "bad kv at s=1030.00=20")
}

func TestParseQrStrReturnErrOnBadParams(t *testing.T) {
	_, _, _, _, _, _, err := ParseQrStr("t=bad_date&s=1030.00&fn=1234567890123456&i=12345&fp=1234567890&n=1")
	require.EqualError(t, err, "bad date bad_date")

	_, _, _, _, _, _, err = ParseQrStr("t=20200115T2110&s=bad_sum&fn=1234567890123456&i=12345&fp=1234567890&n=1")
	require.EqualError(t, err, "bad sum bad_sum")

	_, _, _, _, _, _, err = ParseQrStr("t=20200115T2110&s=1030.00&fn=1234567890123456&i=bad_fd&fp=1234567890&n=1")
	require.EqualError(t, err, "bad fd bad_fd")

	_, _, _, _, _, _, err = ParseQrStr("t=20200115T2110&s=1030.00&fn=1234567890123456&i=12345&fp=bad_fp&n=1")
	require.EqualError(t, err, "bad fpd bad_fp")

	_, _, _, _, _, _, err = ParseQrStr("t=20200115T2110&s=1030.00&fn=1234567890123456&i=12345&fp=1234567890&n=bad_optype")
	require.EqualError(t, err, "bad opType bad_optype")
}

func TestParseQrStrReturnErrOnUnexpectedParam(t *testing.T) {
	_, _, _, _, _, _, err := ParseQrStr("t=20200115T2110&s=1030.00&fn=1234567890123456&i=12345&fp=1234567890&x=1")
	require.EqualError(t, err, "unexpected key x")
}
