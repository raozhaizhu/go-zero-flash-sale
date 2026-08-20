package testutil

import (
	"fmt"
	"math/rand/v2"
	"strings"

	"github.com/google/uuid"
)

/** ====================================================================================
 * 🏁 RandomString
 * =====================================================================================
 */

// alphabet 小写字母表
const alphabet = "abcdefghijklmnopqrstuvwxyz"

// RandomString 返回长度为 n 的随机字符串
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.IntN(k)] // 随机字符
		sb.WriteByte(c)
	}

	return sb.String()
}

// RandomUsername 返回随机用户名
func RandomUsername() string {
	return RandomString(6)
}

func RandomPassword() string {
	return RandomString(8)
}

func RandomEmail() string {
	return fmt.Sprintf("%s@example.com", RandomUsername())
}

func DeriveEmail(username string) string {
	return fmt.Sprintf("%s@example.com", username)
}

func RandomUUID() string {
	uuid := uuid.New()
	return uuid.String()
}
