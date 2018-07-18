package main

import "testing"

func TestReservationFlow(t *testing.T) {
	token := NewToken()

	// Make reservation
	reservation := NewReservation(&token)
	if len(reservations) != 1 {
		t.Fatal("Unexpected number of reservations.")
	}
	if reservation.Id == "" {
		t.Fatal("Reservation does not contain an ID.")
	}
	if reservation.Token != &token {
		t.Fatal("Reservation does not reference the token.")
	}

	// Leave reservation
	reservation.Delete()
	if len(reservations) != 0 {
		t.Fatal("Unexpected number of reservations.")
	}
}
func TestCreateReservation(t *testing.T) {
	// Make reservation
}
func TestDeleteReservationWithoutUse(t *testing.T) {
	// Make reservation
}

func BenchmarkNewToken(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewToken()
	}
}
