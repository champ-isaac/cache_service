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

func TestClient_HSet_N_HGetAllNScan(t *testing.T) {
	for i := 0; i < 100; i++ {
		key := fmt.Sprintf(vinFmt, randInt(minNum, maxNum))
		value := map[string]any{
			"raw_value":   key,
			"token_value": sha512HashToBase64(key),
			"active":      true,
			"timestamp":   time.Now().Format("2006-01-02 15:04:05"),
		}
		_, err := client.HSet(ctx, fmt.Sprintf("%s:%s", vinRef, key), value)
		assert.NoError(t, err)

		var out map[string]any
		err = client.HGetAllNScan(ctx, key, &out)
		assert.NoError(t, err)

		assert.Equal(t, key, out["raw_value"].(string))
		assert.Equal(t, value, out["token_value"].(string))
		assert.True(t, out["active"].(bool))
	}

	for i := 0; i < 50; i++ {
		key := fmt.Sprintf(uuidFmt, randInt(minNum, maxNum))
		value := map[string]any{
			"raw_value":   key,
			"token_value": sha512HashToBase64(key),
			"active":      false,
			"timestamp":   time.Now().Format("2006-01-02 15:04:05"),
		}
		_, err := client.HSet(ctx, fmt.Sprintf("%s:%s", uuidRef, key), value)
		assert.NoError(t, err)

		var out map[string]any
		err = client.HGetAllNScan(ctx, key, &out)
		assert.NoError(t, err)

		assert.Equal(t, key, out["raw_value"].(string))
		assert.Equal(t, value, out["token_value"].(string))
		assert.False(t, out["active"].(bool))
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
