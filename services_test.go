package makeplans

import "testing"

var mockServiceWrapCreate = serviceWrap{
	Service: Service{
		Active:          true,
		BookingCapacity: 1,
		Interval:        30,
		MaxSlots:        100,
		Price:           "115.0",
		SameDay:         true,
		Title:           "Chiropractor",
		BookingTypeID:   1,
	},
}

func TestService_create(t *testing.T) {
	_, client := mockServerClient(t)

	svc := mockServiceWrapCreate.Service
	svc, err := client.ServiceCreate(svc)
	if err != nil {
		t.Fatal(err)
	}
	if e := true; svc.Active != e {
		t.Errorf("got: %t wanted: %t", svc.Active, e)
	}

	if e := 1; svc.BookingTypeID != e {
		t.Errorf("got: %d wanted: %d", svc.BookingTypeID, e)
	}
}

var testServiceDelete = []byte(`{"service":{"active":false,"booking_capacity":1,"booking_type_id":1,"created_at":"2015-09-11T23:01:27-05:00","custom_data":{},"day_booking_specify_time":null,"description":null,"has_day_booking":false,"id":401,"interval":30,"mail_confirmation":null,"mail_verification":null,"max_slots":100,"price":"115.0","same_day":true,"sms_confirmation":null,"sms_reminder":null,"sms_verification":null,"template":null,"title":"Chiropractor","updated_at":"2015-12-07T21:44:11-06:00"}}`)

func TestService_delete(t *testing.T) {
	_, client := mockServerClient(t)
	// client := New(ac.Name, ac.Token)
	svc, err := client.ServiceDelete(401)
	if err != nil {
		t.Fatal(err)
	}

	if e := false; svc.Active != e {
		t.Errorf("got: %t wanted: %t", svc.Active, e)
	}
}

func TestService_list(t *testing.T) {

	_, client := mockServerClient(t)

	svcs, err := client.Services()
	if err != nil {
		t.Error(err)
	}

	if e := 3; len(svcs) != e {
		t.Fatalf("wrong number of services returned got: %d wanted: %d",
			len(svcs), e)
	}

	if e := "20.0"; svcs[0].Price != e {
		t.Fatalf("got: %s wanted: %s", svcs[0].Price, e)
	}

	if e := "Cross Fit Session"; e != svcs[0].Title {
		t.Fatalf("got: %s wanted: %s", svcs[0].Title, e)
	}

}
