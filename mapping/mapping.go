package mapping

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/cloudquery/plugin-sdk/schema"
	"github.com/thoas/go-funk"
)

type Transform func(ctx context.Context, v any) (any, error)

type key string

const (
	MetaKey     key = "meta"
	ResourceKey key = "resource"
	ColumnKey   key = "column"
)

func Apply(transforms ...Transform) schema.ColumnResolver {
	return func(ctx context.Context, meta schema.ClientMeta, r *schema.Resource, c schema.Column) error {
		var (
			value any
			err   error
		)

		ctx = context.WithValue(ctx, MetaKey, meta)
		ctx = context.WithValue(ctx, ResourceKey, r)
		ctx = context.WithValue(ctx, ColumnKey, c)

		for _, transform := range transforms {
			value, err = transform(ctx, value)
			if err != nil {
				return err
			}
		}
		return r.Set(c.Name, value)
	}
}

func GetObjectField(path string) Transform {
	return func(ctx context.Context, _ any) (any, error) {
		v := ctx.Value(ResourceKey)
		if v == nil {
			return nil, errors.New("no resource info in context")
		}
		if r, ok := v.(*schema.Resource); ok {
			value := funk.Get(r.Item, path, funk.WithAllowZero())
			return value, nil
		}
		return nil, errors.New("")
	}
}

func GetMapEntry[K comparable](key K) Transform {
	return func(ctx context.Context, v any) (any, error) {
		if v != nil {
			// log.Printf("v not nil: '%v': %T", v, v)
			switch value := v.(type) {
			case map[K]any:
				// log.Printf("v is map")
				if t, ok := value[key]; ok {
					return t, nil
				}
				// default:
				// 	return nil, fmt.Errorf("invalid type: expected map, got %T", v)
			}
		}
		// log.Printf("returning %q", result)
		return nil, nil
	}
}

func TrimString() Transform {
	return func(ctx context.Context, v any) (any, error) {
		if v == nil {
			return nil, nil
		}
		if s, ok := v.(string); ok {
			return strings.TrimSpace(s), nil
		}
		return nil, fmt.Errorf("invalid type: expected string, got %T", v)
	}
}

func NilIfZero() Transform {
	return func(ctx context.Context, v any) (any, error) {
		if v == nil {
			return nil, nil
		}
		if funk.IsZero(v) {
			return nil, nil
		}
		return v, nil
	}
}

func ToString() Transform {
	return func(ctx context.Context, v any) (any, error) {
		if v == nil {
			return nil, nil
		}
		switch v := v.(type) {
		case string:
			return v, nil
		default:
			return fmt.Sprintf("%v", v), nil
		}
	}
}

func ToInt() Transform {
	return func(ctx context.Context, v any) (any, error) {
		switch v := v.(type) {
		case int:
			return v, nil
		case int8:
			return int(v), nil
		case int16:
			return int(v), nil
		case int32:
			return int(v), nil
		case int64:
			return int(v), nil
		case uint:
			return v, nil
		case uint8:
			return int(v), nil
		case uint16:
			return int(v), nil
		case uint32:
			return int(v), nil
		case uint64:
			return int(v), nil
		case float32:
			return int(v), nil
		case float64:
			return int(v), nil
		case string:
			if strings.TrimSpace(v) != "" {
				return strconv.ParseInt(v, 10, 32)
			}
		default:
			// log.Printf("ToInt: unsupported type: %T", v)
			return nil, fmt.Errorf("unsupported value type: %T", v)
		}
		return nil, nil
	}
}

func OrDefault(value any) Transform {
	return func(ctx context.Context, v any) (any, error) {
		if v == nil {
			return value, nil
		}
		return v, nil
	}
}

func RemapValue[K comparable, V any](remap map[K]V) Transform {
	return func(ctx context.Context, v any) (any, error) {
		if key, ok := v.(K); ok {
			if value, ok := remap[key]; ok {
				return value, nil
			} else {
				return v, nil
			}
		}
		var t K
		return nil, fmt.Errorf("invalid data type: expected %T, got %T", t, v)
	}
}
