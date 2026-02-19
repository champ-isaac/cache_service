package redis

import (
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	vinRef  = "vin"
	uuidRef = "uuid"
	minNum  = 1
	maxNum  = 999999
	vinFmt  = "JTI%06d"
	uuidFmt = "GUID%06d"
)

type CacheData struct {
	Raw       string    `redis:"raw_value"`
	Token     string    `redis:"token_value"`
	Active    bool      `redis:"active"`
	Timestamp time.Time `redis:"timestamp"`
}

func TestClient_HSet_N_HGetAllNScan(t *testing.T) {
	for i := 0; i < 100; i++ {
		raw := fmt.Sprintf(vinFmt, randInt(minNum, maxNum))
		value := CacheData{
			Raw:       raw,
			Token:     sha512HashToBase64(raw),
			Active:    true,
			Timestamp: time.Now(),
		}
		key := fmt.Sprintf("%s:%s", vinRef, raw)
		_, err := client.HSet(ctx, key, value)
		assert.NoError(t, err)

		var out CacheData
		err = client.HGetAllNScan(ctx, key, &out)
		assert.NoError(t, err)

		assert.Equal(t, value.Raw, out.Raw)
		assert.Equal(t, value.Token, out.Token)
		assert.True(t, value.Active, out.Active)
	}

	for i := 0; i < 50; i++ {
		raw := fmt.Sprintf(uuidFmt, randInt(minNum, maxNum))
		value := CacheData{
			Raw:       raw,
			Token:     sha512HashToBase64(raw),
			Active:    false,
			Timestamp: time.Now(),
		}
		key := fmt.Sprintf("%s:%s", uuidRef, raw)
		_, err := client.HSet(ctx, key, value)
		assert.NoError(t, err)

		var out CacheData
		err = client.HGetAllNScan(ctx, key, &out)
		assert.NoError(t, err)

		assert.Equal(t, value.Raw, out.Raw)
		assert.Equal(t, value.Token, out.Token)
		assert.False(t, value.Active, out.Active)
	}
}

func TestClient_KeyNotExist_HGetNScan(t *testing.T) {
	var out map[string]any
	err := client.HGetAllNScan(ctx, "Not_Exist", &out)
	assert.NoError(t, err)
	assert.Equal(t, map[string]any{}, out)
}

func sha512HashToBase64(value string) string {
	sum := sha512.Sum512([]byte(value))
	return base64.StdEncoding.EncodeToString(sum[:])
}

func randInt(min int, max int) int {
	rand.NewSource(time.Now().UnixNano())
	return min + rand.Intn(max-min)
}
