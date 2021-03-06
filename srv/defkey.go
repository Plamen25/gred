package srv

import (
	"time"

	"github.com/PuerkitoBio/gred/types"
)

var dv defVal

var empty = []string{}

var (
	_ Key          = (*defKey)(nil)
	_ types.String = (*defVal)(nil)
	_ types.Hash   = (*defVal)(nil)
	_ types.List   = (*defVal)(nil)
	_ types.Set    = (*defVal)(nil)
)

type defKey string

func (d defKey) Lock()                                 {}
func (d defKey) Unlock()                               {}
func (d defKey) RLock()                                {}
func (d defKey) RUnlock()                              {}
func (d defKey) Expire(_ time.Duration, _ func()) bool { return true }
func (d defKey) TTL() time.Duration                    { return 0 }
func (d defKey) Abort() bool                           { return true }
func (d defKey) Val() types.Value                      { return dv }

func (d defKey) Name() string { return string(d) }

type defVal struct{}

func (d defVal) Type() string { panic("Type called on defKey value") }

// String implementation
func (d defVal) Append(_ string) int64            { return 0 }
func (d defVal) Get() string                      { return "" }
func (d defVal) GetRange(_, _ int64) string       { return "" }
func (d defVal) GetSet(_ string) string           { return "" }
func (d defVal) Set(_ string)                     {}
func (d defVal) SetRange(_ int64, _ string) int64 { return 0 }
func (d defVal) StrLen() int64                    { return 0 }

// Hashes implementation
func (d defVal) HDel(_ ...string) int64               { return 0 }
func (d defVal) HExists(_ string) bool                { return false }
func (d defVal) HGet(_ string) (string, bool)         { return "", false }
func (d defVal) HGetAll() []string                    { return empty }
func (d defVal) HKeys() []string                      { return empty }
func (d defVal) HLen() int64                          { return 0 }
func (d defVal) HMGet(fields ...string) []interface{} { return make([]interface{}, len(fields)) }
func (d defVal) HMSet(_ ...string)                    {}
func (d defVal) HSet(_, _ string) bool                { return false }
func (d defVal) HSetNx(_, _ string) bool              { return false }
func (d defVal) HVals() []string                      { return empty }

// Lists implementation
func (d defVal) LIndex(_ int64) (string, bool)   { return "", false }
func (d defVal) LInsertBefore(_, _ string) int64 { return 0 }
func (d defVal) LInsertAfter(_, _ string) int64  { return 0 }
func (d defVal) LLen() int64                     { return 0 }
func (d defVal) LPop() (string, bool)            { return "", false }
func (d defVal) LPush(_ ...string) int64         { return 0 }
func (d defVal) LRange(_, _ int64) []string      { return empty }
func (d defVal) LRem(_ int64, _ string) int64    { return 0 }
func (d defVal) LSet(_ int64, _ string) bool     { return false }
func (d defVal) LTrim(_, _ int64)                {}
func (d defVal) RPop() (string, bool)            { return "", false }
func (d defVal) RPush(_ ...string) int64         { return 0 }

// Sets implementation
func (d defVal) SAdd(_ ...string) int64         { return 0 }
func (d defVal) SCard() int64                   { return 0 }
func (d defVal) SDiff(_ ...types.Set) []string  { return empty }
func (d defVal) SInter(_ ...types.Set) []string { return empty }
func (d defVal) SIsMember(_ string) bool        { return false }
func (d defVal) SMembers() []string             { return empty }
func (d defVal) SRem(_ ...string) int64         { return 0 }
func (d defVal) SUnion(_ ...types.Set) []string { return empty }
