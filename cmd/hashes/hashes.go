package hashes

import (
	"fmt"

	"github.com/PuerkitoBio/gred/cmd"
	"github.com/PuerkitoBio/gred/srv"
	"github.com/PuerkitoBio/gred/types"
)

func init() {
	cmd.Register("hdel", hdel)
	cmd.Register("hexists", hexists)
	cmd.Register("hget", hget)
	cmd.Register("hgetall", hgetall)
	cmd.Register("hincrby", hincrby)
	cmd.Register("hincrbyfloat", hincrbyfloat)
	cmd.Register("hkeys", hkeys)
	cmd.Register("hlen", hlen)
	cmd.Register("hmget", hmget)
	cmd.Register("hmset", hmset)
	cmd.Register("hset", hset)
	cmd.Register("hsetnx", hsetnx)
	cmd.Register("hvals", hvals)
}

var hdel = cmd.NewDBCmd(
	&cmd.ArgDef{
		MinArgs: 2,
		MaxArgs: -1,
	},
	hdelFn)

func hdelFn(db srv.DB, args []string, ints []int64, floats []float64) (interface{}, error) {
	// Since HDEL may delete the key (if the hash is empty), must get an exclusive
	// DB lock right away (can't think of a sane way to upgrade the lock without restartint
	// the whole operation).
	k, unl := db.XLockGetKey(args[0], srv.NoKeyDefaultVal)
	defer unl()

	// Lock the key
	k.Lock()
	defer k.Unlock()

	// Delete the values
	v := k.Val()
	if v, ok := v.(types.Hash); ok {
		ret := v.HDel(args[1:]...)
		if ret > 0 {
			// Is it now an empty hash?
			if v.HLen() == 0 {
				db.DelKey(args[0])
			}
		}
		return ret, nil
	}
	return nil, cmd.ErrInvalidValType
}

var hexists = cmd.NewSingleKeyCmd(
	&cmd.ArgDef{
		MinArgs: 2,
		MaxArgs: 2,
	},
	srv.NoKeyDefaultVal,
	hexistsFn)

func hexistsFn(k srv.Key, args []string, ints []int64, floats []float64) (interface{}, error) {
	k.RLock()
	defer k.RUnlock()

	v := k.Val()
	if v, ok := v.(types.Hash); ok {
		return v.HExists(args[1]), nil
	}
	return nil, cmd.ErrInvalidValType
}

var hget = cmd.NewSingleKeyCmd(
	&cmd.ArgDef{
		MinArgs: 2,
		MaxArgs: 2,
	},
	srv.NoKeyDefaultVal,
	hgetFn)

func hgetFn(k srv.Key, args []string, ints []int64, floats []float64) (interface{}, error) {
	k.RLock()
	defer k.RUnlock()

	v := k.Val()
	if v, ok := v.(types.Hash); ok {
		ret, ok := v.HGet(args[1])
		if ok {
			return ret, nil
		}
		return nil, nil
	}
	return nil, cmd.ErrInvalidValType
}

var hgetall = cmd.NewSingleKeyCmd(
	&cmd.ArgDef{
		MinArgs: 1,
		MaxArgs: 1,
	},
	srv.NoKeyDefaultVal,
	hgetallFn)

func hgetallFn(k srv.Key, args []string, ints []int64, floats []float64) (interface{}, error) {
	k.RLock()
	defer k.RUnlock()

	v := k.Val()
	if v, ok := v.(types.Hash); ok {
		return v.HGetAll(), nil
	}
	return nil, cmd.ErrInvalidValType
}

var hincrby = cmd.NewSingleKeyCmd(
	&cmd.ArgDef{
		MinArgs:    3,
		MaxArgs:    3,
		IntIndices: []int{2},
	},
	srv.NoKeyCreateHash,
	hincrbyFn)

func hincrbyFn(k srv.Key, args []string, ints []int64, floats []float64) (interface{}, error) {
	k.Lock()
	defer k.Unlock()

	v := k.Val()
	if v, ok := v.(types.IncHash); ok {
		val, ok := v.HIncrBy(args[1], ints[0])
		if ok {
			return val, nil
		}
		return nil, cmd.ErrHashFieldNotInt
	}
	return nil, cmd.ErrInvalidValType
}

var hincrbyfloat = cmd.NewSingleKeyCmd(
	&cmd.ArgDef{
		MinArgs:      3,
		MaxArgs:      3,
		FloatIndices: []int{2},
	},
	srv.NoKeyCreateHash,
	hincrbyfloatFn)

func hincrbyfloatFn(k srv.Key, args []string, ints []int64, floats []float64) (interface{}, error) {
	k.Lock()
	defer k.Unlock()

	v := k.Val()
	if v, ok := v.(types.IncHash); ok {
		val, ok := v.HIncrByFloat(args[1], floats[0])
		if ok {
			return val, nil
		}
		return nil, cmd.ErrHashFieldNotFloat
	}
	return nil, cmd.ErrInvalidValType
}

var hkeys = cmd.NewSingleKeyCmd(
	&cmd.ArgDef{
		MinArgs: 1,
		MaxArgs: 1,
	},
	srv.NoKeyDefaultVal,
	hkeysFn)

func hkeysFn(k srv.Key, args []string, ints []int64, floats []float64) (interface{}, error) {
	k.RLock()
	defer k.RUnlock()

	v := k.Val()
	if v, ok := v.(types.Hash); ok {
		return v.HKeys(), nil
	}
	return nil, cmd.ErrInvalidValType
}

var hlen = cmd.NewSingleKeyCmd(
	&cmd.ArgDef{
		MinArgs: 1,
		MaxArgs: 1,
	},
	srv.NoKeyDefaultVal,
	hlenFn)

func hlenFn(k srv.Key, args []string, ints []int64, floats []float64) (interface{}, error) {
	k.RLock()
	defer k.RUnlock()

	v := k.Val()
	if v, ok := v.(types.Hash); ok {
		return v.HLen(), nil
	}
	return nil, cmd.ErrInvalidValType
}

var hmget = cmd.NewSingleKeyCmd(
	&cmd.ArgDef{
		MinArgs: 2,
		MaxArgs: -1,
	},
	srv.NoKeyDefaultVal,
	hmgetFn)

func hmgetFn(k srv.Key, args []string, ints []int64, floats []float64) (interface{}, error) {
	k.RLock()
	defer k.RUnlock()

	v := k.Val()
	if v, ok := v.(types.Hash); ok {
		return v.HMGet(args[1:]...), nil
	}
	return nil, cmd.ErrInvalidValType
}

var hmset = cmd.NewSingleKeyCmd(
	&cmd.ArgDef{
		MinArgs: 3,
		MaxArgs: -1,
		ValidateFn: func(args []string, ints []int64, floats []float64) error {
			// Must have odd number of args (key + any number of tuples field-value)
			if len(args)%2 == 0 {
				return fmt.Errorf(cmd.WrongNumberOfArgsFmt, "hmset")
			}
			return nil
		},
	},
	srv.NoKeyCreateHash,
	hmsetFn)

func hmsetFn(k srv.Key, args []string, ints []int64, floats []float64) (interface{}, error) {
	k.Lock()
	defer k.Unlock()

	v := k.Val()
	if v, ok := v.(types.Hash); ok {
		v.HMSet(args[1:]...)
		return cmd.OKVal, nil
	}
	return nil, cmd.ErrInvalidValType
}

var hset = cmd.NewSingleKeyCmd(
	&cmd.ArgDef{
		MinArgs: 3,
		MaxArgs: 3,
	},
	srv.NoKeyCreateHash,
	hsetFn)

func hsetFn(k srv.Key, args []string, ints []int64, floats []float64) (interface{}, error) {
	k.Lock()
	defer k.Unlock()

	v := k.Val()
	if v, ok := v.(types.Hash); ok {
		return v.HSet(args[1], args[2]), nil
	}
	return nil, cmd.ErrInvalidValType
}

var hsetnx = cmd.NewSingleKeyCmd(
	&cmd.ArgDef{
		MinArgs: 3,
		MaxArgs: 3,
	},
	srv.NoKeyCreateHash,
	hsetnxFn)

func hsetnxFn(k srv.Key, args []string, ints []int64, floats []float64) (interface{}, error) {
	k.Lock()
	defer k.Unlock()

	v := k.Val()
	if v, ok := v.(types.Hash); ok {
		return v.HSetNx(args[1], args[2]), nil
	}
	return nil, cmd.ErrInvalidValType
}

var hvals = cmd.NewSingleKeyCmd(
	&cmd.ArgDef{
		MinArgs: 1,
		MaxArgs: 1,
	},
	srv.NoKeyDefaultVal,
	hvalsFn)

func hvalsFn(k srv.Key, args []string, ints []int64, floats []float64) (interface{}, error) {
	k.RLock()
	defer k.RUnlock()

	v := k.Val()
	if v, ok := v.(types.Hash); ok {
		return v.HVals(), nil
	}
	return nil, cmd.ErrInvalidValType
}
