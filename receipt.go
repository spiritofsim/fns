package fns

import (
	"fmt"
	"time"
)

type ReceiptItem struct {
	CalculationSubjectSign int
	Quantity               int
	NdsRate                int
	NdsSum                 float32
	Name                   string
	Sum                    float32
	Price                  float32
}

type Receipt struct {
	PrepaymentSum        float32
	FiscalDocumentNumber int
	CounterSubmissionSum float32
	TaxationType         int
	Code                 int
	FiscalDriveNumber    string
	RawData              string
	Items                []ReceiptItem
	Nds18                float32
	FiscalSign           int
	Operator             string
	OperationType        int
	RequestNumber        int
	PostpaymentSum       float32
	ShiftNumber          int
	OperatorInn          string
	EcashTotalSum        float32
	ProtocolVersion      int
	DateTime             time.Time
	TotalSum             float32
	UserInn              string
	KktRegID             string
	CashTotalSum         float32
}

func NewReceipt(prj receiptPrj) (Receipt, error) {
	dt, err := time.Parse("2006-01-02T15:04:05", prj.Document.Receipt.DateTime)
	if err != nil {
		return Receipt{}, fmt.Errorf("bad date %v", prj.Document.Receipt.DateTime)
	}

	items := make([]ReceiptItem, len(prj.Document.Receipt.Items))
	for i, item := range prj.Document.Receipt.Items {
		items[i] = ReceiptItem{
			CalculationSubjectSign: item.CalculationSubjectSign,
			Quantity:               item.Quantity,
			NdsRate:                item.NdsRate,
			NdsSum:                 float32(item.NdsSum) / 100,
			Name:                   item.Name,
			Sum:                    float32(item.Sum) / 100,
			Price:                  float32(item.Price) / 100,
		}
	}

	return Receipt{
		PrepaymentSum:        float32(prj.Document.Receipt.PrepaymentSum) / 100,
		FiscalDocumentNumber: prj.Document.Receipt.FiscalDocumentNumber,
		CounterSubmissionSum: float32(prj.Document.Receipt.CounterSubmissionSum) / 100,
		TaxationType:         prj.Document.Receipt.TaxationType,
		Code:                 prj.Document.Receipt.ReceiptCode,
		FiscalDriveNumber:    prj.Document.Receipt.FiscalDriveNumber,
		RawData:              prj.Document.Receipt.RawData,
		Items:                items,
		Nds18:                float32(prj.Document.Receipt.Nds18) / 100,
		FiscalSign:           prj.Document.Receipt.FiscalSign,
		Operator:             prj.Document.Receipt.Operator,
		OperationType:        prj.Document.Receipt.OperationType,
		RequestNumber:        prj.Document.Receipt.RequestNumber,
		PostpaymentSum:       float32(prj.Document.Receipt.PostpaymentSum) / 100,
		ShiftNumber:          prj.Document.Receipt.ShiftNumber,
		OperatorInn:          prj.Document.Receipt.OperatorInn,
		EcashTotalSum:        float32(prj.Document.Receipt.EcashTotalSum) / 100,
		ProtocolVersion:      prj.Document.Receipt.ProtocolVersion,
		DateTime:             dt,
		TotalSum:             float32(prj.Document.Receipt.TotalSum) / 100,
		UserInn:              prj.Document.Receipt.UserInn,
		KktRegID:             prj.Document.Receipt.KktRegID,
		CashTotalSum:         float32(prj.Document.Receipt.CashTotalSum) / 100,
	}, nil
}

type receiptPrj struct {
	Document struct {
		Receipt struct {
			PrepaymentSum        int    `json:"prepaymentSum"`
			FiscalDocumentNumber int    `json:"fiscalDocumentNumber"`
			CounterSubmissionSum int    `json:"counterSubmissionSum"`
			TaxationType         int    `json:"taxationType"`
			ReceiptCode          int    `json:"receiptCode"`
			FiscalDriveNumber    string `json:"fiscalDriveNumber"`
			RawData              string `json:"rawData"`
			Items                []struct {
				CalculationSubjectSign int    `json:"calculationSubjectSign"`
				Quantity               int    `json:"quantity"`
				NdsRate                int    `json:"ndsRate"`
				NdsSum                 int    `json:"ndsSum"`
				Name                   string `json:"name"`
				Sum                    int    `json:"sum"`
				Price                  int    `json:"price"`
			} `json:"items"`
			Nds18           int    `json:"nds18"`
			FiscalSign      int    `json:"fiscalSign"`
			Operator        string `json:"operator"`
			OperationType   int    `json:"operationType"`
			RequestNumber   int    `json:"requestNumber"`
			PostpaymentSum  int    `json:"postpaymentSum"`
			ShiftNumber     int    `json:"shiftNumber"`
			OperatorInn     string `json:"operatorInn"`
			EcashTotalSum   int    `json:"ecashTotalSum"`
			ProtocolVersion int    `json:"protocolVersion"`
			DateTime        string `json:"dateTime"`
			TotalSum        int    `json:"totalSum"`
			UserInn         string `json:"userInn"`
			KktRegID        string `json:"kktRegId"`
			CashTotalSum    int    `json:"cashTotalSum"`
		} `json:"receipt"`
	} `json:"document"`
}
