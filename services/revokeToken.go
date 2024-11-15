package services

import (
	"context"
	"fmt"
	"main-module/initializers"
	"time"
)

func RevokeToken(jti string, exp time.Time) error {
	ttl := time.Until(exp)
	ctx := context.Background()
	err := initializers.RedisClient.Set(ctx, jti, "revoked", ttl).Err()
	if err != nil {
		fmt.Println("Error storing token in Redis:", err)
		return err
	}
	fmt.Println("Token revoked successfully:", jti)
	return nil
}
