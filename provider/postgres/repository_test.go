package postgres

import (
	"context"
	"cw/dbutil"
	"cw/models"
	"testing"
)

func TestGetRepository(t *testing.T) {
	lib_db := dbutil.InitDB()
	defer lib_db.Close()

	pr_db := NewProviderRepository(lib_db)
	if err := pr_db.AddProvider(context.Background(), &models.Provider{VendorCode: 1, Name: "abcd", UNP: "123456789", TermsOfPayment: "asdf", Address: "zxcvb", PhoneNumber: "+375298475820", Email: "asdf", WebSite: "asdf"}); err != nil {
		t.Error(err)
	}

	result, err := pr_db.GetProviders(context.Background())
	if err != nil {
		t.Errorf("Test get provider: %v", err)
	}

	for _, provider := range result {
		t.Error(provider)
	}
}
