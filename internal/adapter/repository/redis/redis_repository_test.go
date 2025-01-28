package redis

import (
	"context"
	"testing"

	"github.com/go-redis/redismock/v9"
)

func TestAuthRepository(t *testing.T) {
	db, mock := redismock.NewClientMock()

	t.Run("Ping", func(t *testing.T) {
		expected := "PONG"
		mock.ExpectPing().SetVal(expected)

		ping, err := db.Ping(context.Background()).Result()
		if err != nil {
			// because t.Fatalf is equivalent to Log and FailNow!!!!
			// this should not fail at all.
			t.Fatalf("err pinging db: %v", err.Error())
		}

		if ping != expected {
			t.Fatalf("ping was expected: %v", err.Error())
		}
	})

	t.Run("SaveUser", func(t *testing.T) {
		t.Run("success", func(t *testing.T) {
		})
	})

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("Expectations were not met: %v", err)
	}
}
