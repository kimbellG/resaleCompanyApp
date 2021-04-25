package models

import (
	"testing"
)

func TestEditProviderField(t *testing.T) {
	db := CreateNewDBConnection()
	db.createTables()
	pr := &Provider{Name: "test", UNP: "test12345", TermsOfPayment: "test",
		Address: "test", PhoneNumber: "+375294321232", Email: "test", WebSite: "test"}
	if err := db.InsertProviderInDB(pr); err != nil {
		t.Fatalf("failed insert value: %v", err)
	}

	vendor_code, err := getVendorCodeByName(db, pr.Name)
	if err != nil {
		t.Fatal("get vendor code:", err)
	}

	new_fields := map[string]string{"name": "abc", "unp": "new_unp12", "terms_of_payment": "new",
		"address": "new"}
	if err := db.EditFieldsInProviderTable(vendor_code, new_fields); err != nil {
		t.Fatalf("%v", err)
	}

	new_pr, err := getProviderByVendorCode(db, vendor_code)
	if err != nil {
		t.Fatal(err)
	}

	if new_pr.Name != new_fields["name"] {
		t.Error("name isn't update")
	}

	if new_pr.UNP != new_fields["unp"] {
		t.Errorf("unp isn't update: %v!=%v", new_pr.UNP, new_fields["unp"])
	}

	if new_pr.TermsOfPayment != new_fields["terms_of_payment"] {
		t.Error("terms_of_payment isn't update")
	}

	if new_pr.Address != new_fields["address"] {
		t.Error("address isn't update")
	}

	if err := dltByVendorCode(db, vendor_code); err != nil {
		t.Fatal(err)
	}
}

func getVendorCodeByName(db *DBController, name string) (int, error) {
	stmt, err := db.Prepare("SELECT vendor_code FROM provider WHERE name = $1")
	if err != nil {
		return -1, err
	}

	vendor_code, err := stmt.Query(name)
	if err != nil {
		return -1, err
	}

	var result int
	for vendor_code.Next() {
		if err := vendor_code.Scan(&result); err != nil {
			return -1, err
		}
	}

	return result, nil
}

func getProviderByVendorCode(db *DBController, vendor int) (*Provider, error) {
	stmt, err := db.Prepare("SELECT name, unp, terms_of_payment, address, phone_number, email, web_site FROM provider WHERE vendor_code = $1")
	if err != nil {
		return nil, err
	}

	dbpr, err := stmt.Query(vendor)
	if err != nil {
		return nil, err
	}

	result := &Provider{}
	for dbpr.Next() {
		if err := dbpr.Scan(&result.Name, &result.UNP, &result.TermsOfPayment, &result.Address, &result.PhoneNumber, &result.Email, &result.WebSite); err != nil {
			return nil, err
		}
	}

	return result, nil
}

func dltByVendorCode(db *DBController, vendor int) error {
	stmt, err := db.Prepare("DELETE FROM provider WHERE vendor_code = $1")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(vendor)
	if err != nil {
		return err
	}

	return nil
}
