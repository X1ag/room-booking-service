package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"
)

type tokenResponse struct {
	Token string `json:"token"`
}

type roomDTO struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	Capacity    *int    `json:"capacity"`
	CreatedAt   string  `json:"createdAt"`
}

type createRoomResponse struct {
	Room roomDTO `json:"room"`
}

type scheduleDTO struct {
	ID         string `json:"id"`
	RoomID     string `json:"roomId"`
	StartTime  string `json:"startTime"`
	EndTime    string `json:"endTime"`
	DaysOfWeek []int  `json:"daysOfWeek"`
}

type createScheduleResponse struct {
	Schedule scheduleDTO `json:"schedule"`
}

type slotDTO struct {
	ID     string `json:"id"`
	RoomID string `json:"roomId"`
	Start  string `json:"start"`
	End    string `json:"end"`
}

type slotsResponse struct {
	Slots []slotDTO `json:"slots"`
}

type bookingDTO struct {
	ID             string  `json:"id"`
	SlotID         string  `json:"slotId"`
	UserID         string  `json:"userId"`
	Status         string  `json:"status"`
	ConferenceLink *string `json:"conferenceLink"`
	CreatedAt      string  `json:"createdAt"`
}

type createBookingResponse struct {
	Booking bookingDTO `json:"booking"`
}

type errorEnvelope struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func TestBookingFlow(t *testing.T) {
	server := newTestServer(t)
	client := server.Client()
	client.Timeout = 10 * time.Second

	statusCode, body := doRequest(t, client, server.URL, http.MethodPost, "/dummyLogin", "", map[string]any{
		"role": "admin",
	})
	if statusCode != http.StatusOK {
		t.Fatalf("dummyLogin admin failed with status %d: %s", statusCode, string(body))
	}

	var adminLogin tokenResponse
	if err := json.Unmarshal(body, &adminLogin); err != nil {
		t.Fatalf("decode admin token response: %v", err)
	}
	if adminLogin.Token == "" {
		t.Fatal("expected admin token")
	}

	statusCode, body = doRequest(t, client, server.URL, http.MethodPost, "/rooms/create", adminLogin.Token, map[string]any{
		"name":        fmt.Sprintf("room-e2e-%d", time.Now().UTC().UnixNano()),
		"description": "e2e room",
		"capacity":    6,
	})
	if statusCode != http.StatusCreated {
		t.Fatalf("create room failed with status %d: %s", statusCode, string(body))
	}

	var createdRoom createRoomResponse
	if err := json.Unmarshal(body, &createdRoom); err != nil {
		t.Fatalf("decode create room response: %v", err)
	}
	if createdRoom.Room.ID == "" {
		t.Fatal("expected room id")
	}

	scheduleDays := []int{1, 2, 3, 4, 5}
	statusCode, body = doRequest(t, client, server.URL, http.MethodPost, "/rooms/"+createdRoom.Room.ID+"/schedule/create", adminLogin.Token, map[string]any{
		"daysOfWeek": scheduleDays,
		"startTime":  "09:00",
		"endTime":    "11:00",
	})
	if statusCode != http.StatusCreated {
		t.Fatalf("create schedule failed with status %d: %s", statusCode, string(body))
	}

	var createdSchedule createScheduleResponse
	if err := json.Unmarshal(body, &createdSchedule); err != nil {
		t.Fatalf("decode create schedule response: %v", err)
	}
	if createdSchedule.Schedule.RoomID != createdRoom.Room.ID {
		t.Fatalf("expected schedule room id %s, got %s", createdRoom.Room.ID, createdSchedule.Schedule.RoomID)
	}

	statusCode, body = doRequest(t, client, server.URL, http.MethodPost, "/dummyLogin", "", map[string]any{
		"role": "user",
	})
	if statusCode != http.StatusOK {
		t.Fatalf("dummyLogin user failed with status %d: %s", statusCode, string(body))
	}

	var userLogin tokenResponse
	if err := json.Unmarshal(body, &userLogin); err != nil {
		t.Fatalf("decode user token response: %v", err)
	}
	if userLogin.Token == "" {
		t.Fatal("expected user token")
	}

	targetDate := nextDateForWeekdays(scheduleDays, "09:00")
	statusCode, body = doRequest(t, client, server.URL, http.MethodGet, "/rooms/"+createdRoom.Room.ID+"/slots/list?date="+targetDate, userLogin.Token, nil)
	if statusCode != http.StatusOK {
		t.Fatalf("get slots failed with status %d: %s", statusCode, string(body))
	}

	var slots slotsResponse
	if err := json.Unmarshal(body, &slots); err != nil {
		t.Fatalf("decode slots response: %v", err)
	}
	if len(slots.Slots) == 0 {
		t.Fatalf("expected at least one slot for date %s", targetDate)
	}

	statusCode, body = doRequest(t, client, server.URL, http.MethodPost, "/bookings/create", userLogin.Token, map[string]any{
		"slotId":               slots.Slots[0].ID,
		"createConferenceLink": true,
	})
	if statusCode != http.StatusCreated {
		t.Fatalf("create booking failed with status %d: %s", statusCode, string(body))
	}

	var createdBooking createBookingResponse
	if err := json.Unmarshal(body, &createdBooking); err != nil {
		t.Fatalf("decode create booking response: %v", err)
	}
	if createdBooking.Booking.ID == "" {
		t.Fatal("expected booking id")
	}
	if createdBooking.Booking.SlotID != slots.Slots[0].ID {
		t.Fatalf("expected slot id %s, got %s", slots.Slots[0].ID, createdBooking.Booking.SlotID)
	}
	if createdBooking.Booking.Status != "active" {
		t.Fatalf("expected booking status active, got %s", createdBooking.Booking.Status)
	}
	if createdBooking.Booking.ConferenceLink == nil || *createdBooking.Booking.ConferenceLink == "" {
		t.Fatal("expected conference link to be generated")
	}

	statusCode, body = doRequest(t, client, server.URL, http.MethodPost, "/bookings/create", userLogin.Token, map[string]any{
		"slotId":               slots.Slots[0].ID,
		"createConferenceLink": false,
	})
	if statusCode != http.StatusConflict {
		t.Fatalf("expected 409 on duplicate booking, got %d: %s", statusCode, string(body))
	}

	var apiError errorEnvelope
	if err := json.Unmarshal(body, &apiError); err != nil {
		t.Fatalf("decode duplicate booking error: %v", err)
	}
	if apiError.Error.Code != "SLOT_ALREADY_BOOKED" {
		t.Fatalf("expected SLOT_ALREADY_BOOKED, got %s", apiError.Error.Code)
	}
}
