package activity

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"

	"github.com/giladrozner/book_service/pkg/config"
)

type RedisRepository struct {
	client *redis.Client
}

func NewRedisRepository(client *redis.Client) *RedisRepository {
	return &RedisRepository{client: client}
}

var _ Repository = (*RedisRepository)(nil)

func (r *RedisRepository) key(username string) string {
	return fmt.Sprintf(config.ActivityKeyPrefix, username)
}

func (r *RedisRepository) Record(ctx context.Context, username string, action Action) error {
	body, err := json.Marshal(action)
	if err != nil {
		return fmt.Errorf("marshal action: %w", err)
	}
	bg := context.Background()
	key := r.key(username)
	pipe := r.client.Pipeline()
	pipe.LPush(bg, key, string(body))
	pipe.LTrim(bg, key, 0, config.ActivityMaxItems-1)
	pipe.Expire(bg, key, config.ActivityTTL)
	if _, err = pipe.Exec(bg); err != nil {
		return fmt.Errorf("record activity for %s: %w", username, err)
	}
	return nil
}

func (r *RedisRepository) Recent(ctx context.Context, username string, n int) ([]Action, error) {
	raws, err := r.client.LRange(context.Background(), r.key(username), 0, int64(n-1)).Result()
	if err != nil {
		return nil, fmt.Errorf("fetch activity for %s: %w", username, err)
	}
	if len(raws) == 0 {
		return []Action{}, nil
	}
	out := make([]Action, 0, len(raws))
	for _, raw := range raws {
		var a Action
		if err := json.Unmarshal([]byte(raw), &a); err != nil {
			return nil, fmt.Errorf("parse action: %w", err)
		}
		out = append(out, a)
	}
	return out, nil
}
