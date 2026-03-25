package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"
)

type cancelBookingResponse struct {
	Booking bookingDTO `json:"booking"`
}

type userBookingsResponse struct {
	Bookings []bookingDTO `json:"bookings"`
}

func TestCancelBookingFlow(t *testing.T) {
	server := newTestServer(t)
	client := server.Client()
	client.Timeout = 10 * time.Second

	resp, body := doRequest(t, client, server.URL, http.MethodPost, "/dummyLogin", "", map[string]any{
		"role": "admin",
	})
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("dummyLogin admin failed with status %d: %s", resp.StatusCode, string(body))
	}

	var adminLogin tokenResponse
	if err := json.Unmarshal(body, &adminLogin); err != nil {
		t.Fatalf("decode admin token response: %v", err)
	}
	if adminLogin.Token == "" {
		t.Fatal("expected admin token")
	}

	resp, body = doRequest(t, client, server.URL, http.MethodPost, "/rooms/create", adminLogin.Token, map[string]any{
		"name":        fmt.Sprintf("room-cancel-e2e-%d", time.Now().UTC().UnixNano()),
		"description": "cancel e2e room",
		"capacity":    6,
	})
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("create room failed with status %d: %s", resp.StatusCode, string(body))
	}

	var createdRoom createRoomResponse
	if err := json.Unmarshal(body, &createdRoom); err != nil {
		t.Fatalf("decode create room response: %v", err)
	}
	if createdRoom.Room.ID == "" {
		t.Fatal("expected room id")
	}

	scheduleDays := []int{1, 2, 3, 4, 5}
	resp, body = doRequest(t, client, server.URL, http.MethodPost, "/rooms/"+createdRoom.Room.ID+"/schedule/create", adminLogin.Token, map[string]any{
		"daysOfWeek": scheduleDays,
		"startTime":  "09:00",
		"endTime":    "11:00",
	})
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("create schedule failed with status %d: %s", resp.StatusCode, string(body))
	}

	resp, body = doRequest(t, client, server.URL, http.MethodPost, "/dummyLogin", "", map[string]any{
		"role": "user",
	})
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("dummyLogin user failed with status %d: %s", resp.StatusCode, string(body))
	}

	var userLogin tokenResponse
	if err := json.Unmarshal(body, &userLogin); err != nil {
		t.Fatalf("decode user token response: %v", err)
	}
	if userLogin.Token == "" {
		t.Fatal("expected user token")
	}

	targetDate := nextDateForWeekdays(scheduleDays, "09:00")
	resp, body = doRequest(t, client, server.URL, http.MethodGet, "/rooms/"+createdRoom.Room.ID+"/slots/list?date="+targetDate, userLogin.Token, nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("get slots failed with status %d: %s", resp.StatusCode, string(body))
	}

	var slots slotsResponse
	if err := json.Unmarshal(body, &slots); err != nil {
		t.Fatalf("decode slots response: %v", err)
	}
	if len(slots.Slots) == 0 {
		t.Fatalf("expected at least one slot for date %s", targetDate)
	}

	resp, body = doRequest(t, client, server.URL, http.MethodPost, "/bookings/create", userLogin.Token, map[string]any{
		"slotId":               slots.Slots[0].ID,
		"createConferenceLink": true,
	})
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("create booking failed with status %d: %s", resp.StatusCode, string(body))
	}

	var createdBooking createBookingResponse
	if err := json.Unmarshal(body, &createdBooking); err != nil {
		t.Fatalf("decode create booking response: %v", err)
	}
	if createdBooking.Booking.ID == "" {
		t.Fatal("expected booking id")
	}

	resp, body = doRequest(t, client, server.URL, http.MethodPost, "/bookings/"+createdBooking.Booking.ID+"/cancel", userLogin.Token, nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("cancel booking failed with status %d: %s", resp.StatusCode, string(body))
	}

	var cancelledBooking cancelBookingResponse
	if err := json.Unmarshal(body, &cancelledBooking); err != nil {
		t.Fatalf("decode cancel booking response: %v", err)
	}
	if cancelledBooking.Booking.Status != "cancelled" {
		t.Fatalf("expected cancelled status, got %s", cancelledBooking.Booking.Status)
	}

	resp, body = doRequest(t, client, server.URL, http.MethodPost, "/bookings/"+createdBooking.Booking.ID+"/cancel", userLogin.Token, nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("second cancel booking failed with status %d: %s", resp.StatusCode, string(body))
	}

	var cancelledBookingAgain cancelBookingResponse
	if err := json.Unmarshal(body, &cancelledBookingAgain); err != nil {
		t.Fatalf("decode second cancel booking response: %v", err)
	}
	if cancelledBookingAgain.Booking.Status != "cancelled" {
		t.Fatalf("expected cancelled status after second cancel, got %s", cancelledBookingAgain.Booking.Status)
	}

	resp, body = doRequest(t, client, server.URL, http.MethodGet, "/bookings/my", userLogin.Token, nil)
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("get user bookings failed with status %d: %s", resp.StatusCode, string(body))
	}

	var userBookings userBookingsResponse
	if err := json.Unmarshal(body, &userBookings); err != nil {
		t.Fatalf("decode user bookings response: %v", err)
	}

	for _, userBooking := range userBookings.Bookings {
		if userBooking.ID == createdBooking.Booking.ID {
			t.Fatalf("cancelled booking %s should not be returned from /bookings/my", createdBooking.Booking.ID)
		}
	}
}
