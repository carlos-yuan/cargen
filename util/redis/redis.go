package redisd

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/carlos-yuan/cargen/util/convert"
	"github.com/carlos-yuan/cargen/util/lock"
	"github.com/carlos-yuan/cargen/util/timeUtil"
	"github.com/coocood/freecache"
	"github.com/opentracing/opentracing-go/log"

	"github.com/vmihailenco/msgpack"
)

const (
	Nil = redis.Nil
)

type Config struct {
	Addr             string `yaml:"addr"`
	Passwd           string `yaml:"passwd"`
	LocalCache       bool   `yaml:"localCache"`
	Chanel           string `yaml:"chanel"`
	Size             int    `yaml:"size"`
	LocalExpire      int    `yaml:"localExpire"`
	RedisLocalExpire int    `yaml:"redisLocalExpire"`
	FailSleep        int64  `yaml:"failSleep"`
	FailWaitCount    int64  `yaml:"failWaitCount"`
}

type Client interface {
	redis.Cmdable
	Do(ctx context.Context, args ...interface{}) *redis.Cmd
	Subscribe(ctx context.Context, channels ...string) *redis.PubSub
}

func InitRedis(cnf Config) *Decorator {
	addrs := strings.Split(cnf.Addr, ";")
	ctx := context.Background()
	var client Client
	if len(addrs) > 1 {
		opt := redis.ClusterOptions{
			Addrs: addrs,
		}
		if cnf.Passwd != "" {
			opt.Password = cnf.Passwd
		}
		client = redis.NewClusterClient(&opt)
	} else {
		opt := redis.Options{
			Addr: cnf.Addr,
		}
		if cnf.Passwd != "" {
			opt.Password = cnf.Passwd
		}
		client = redis.NewClient(&opt)
	}
	decorator := Decorator{cnf: cnf, Client: client, locker: lock.NewKeyLock()}
	if cnf.LocalCache {
		decorator.cache = freecache.NewCache(cnf.Size * 1024 * 1024)
		go decorator.subscribeLocalCache(ctx, cnf.Chanel)
	}
	return &decorator
}

type Decorator struct {
	cnf Config
	Client
	cache  *freecache.Cache
	locker *lock.KeyLock
}

func (r *Decorator) Set(ctx context.Context, key string, value string, expire time.Duration) error {
	res := r.Client.Set(ctx, key, value, expire)
	if res.Err() != nil {
		return res.Err()
	}
	if r.cnf.LocalCache {
		if r.cnf.LocalExpire > int(expire) {
			r.cnf.LocalExpire = int(expire)
		}
		err := r.cache.Set(convert.Str2bytes(key), convert.Str2bytes(value), r.cnf.LocalExpire)
		if err != nil {
			return err
		}
	}
	return res.Err()
}

func (r *Decorator) SetObj(ctx context.Context, key string, obj interface{}, expire time.Duration) error {
	if obj != nil {
		data, err := msgpack.Marshal(obj)
		if err != nil {
			return err
		}
		err = r.Set(ctx, key, convert.Bytes2str(data), expire)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Decorator) SetObjRefKey(ctx context.Context, key string, obj interface{}, expire time.Duration) error {
	if obj != nil {
		data, err := msgpack.Marshal(obj)
		if err == nil {
			idx := strings.Index(key, "}:")
			if idx != -1 {
				setKey := key[0:idx+2] + "keySet"
				res, err := r.Eval(ctx,
					"if (redis.call('setex',KEYS[1],ARGV[1],ARGV[2]).ok == 'OK') then "+
						"     redis.call('sadd',KEYS[2],KEYS[1]) "+
						"     return 1 "+
						"  else "+
						"     return 0 "+
						"end ", []string{key, setKey}, []string{strconv.Itoa(int(expire / time.Second)), convert.Bytes2str(data)}).Int64()
				if err == nil && res != 1 {
					err = errors.New("add obj and ref key expire script error")
				}
			} else {
				err = errors.New("keys not support error")
			}
		}
		return err
	} else {
		return errors.New("nil object error")
	}
}

func (r *Decorator) Get(ctx context.Context, key string) (string, error) {
	var b []byte
	var err error
	if r.cnf.LocalCache {
		b, err = r.cache.Get(convert.Str2bytes(key))
		if err != nil && err != freecache.ErrNotFound {
			return "", err
		}
		if err == freecache.ErrNotFound {
			state := r.Client.Get(ctx, key)
			b, err = state.Bytes()
			if err != nil {
				return "", err
			}
			if r.cnf.RedisLocalExpire != 0 {
				err = r.cache.Set(convert.Str2bytes(key), b, r.cnf.RedisLocalExpire)
				if err != nil {
					return "", err
				}
			}
		}
	} else {
		state := r.Client.Get(ctx, key)
		b, err = state.Bytes()
		if err != nil {
			return "", err
		}
	}
	if err == nil || err.Error() == redis.Nil.Error() {
		return string(b), nil
	}
	return "", err
}

func (r *Decorator) GetSet(ctx context.Context, key string, expire time.Duration, getFormDb func() (string, error)) (string, error) {
	str, err := r.Get(ctx, key)
	if err == redis.Nil {
		err = nil
		r.locker.Lock(key)
		defer r.locker.Unlock(key)
		str, err = r.Get(ctx, key)
		if err == redis.Nil {
			str, err = getFormDb()
			if err != nil {
				return str, err
			}
			err = r.Set(ctx, key, str, expire)
		}
	}
	return str, err
}

func (r *Decorator) GetObj(ctx context.Context, key string, obj interface{}) error {
	var b []byte
	var err error
	if r.cnf.LocalCache {
		b, err = r.cache.Get(convert.Str2bytes(key))
		if err != nil && err != freecache.ErrNotFound {
			return err
		}
		if err == freecache.ErrNotFound {
			state := r.Client.Get(ctx, key)
			b, err = state.Bytes()
			if err != nil {
				return err
			}
			if r.cnf.RedisLocalExpire != 0 {
				err = r.cache.Set(convert.Str2bytes(key), b, r.cnf.RedisLocalExpire)
				if err != nil {
					return err
				}
			}
		}
	} else {
		state := r.Client.Get(ctx, key)
		b, err = state.Bytes()
		if err != nil {
			return err
		}
	}
	if err != nil && err.Error() != redis.Nil.Error() {
		return err
	}
	if len(b) > 0 {
		err = msgpack.Unmarshal(b, obj)
	}
	return err
}

func (r *Decorator) GetSetObj(ctx context.Context, key string, obj interface{}, expire time.Duration, getFormDb func() (interface{}, error)) error {
	err := r.GetObj(ctx, key, obj)
	if err == redis.Nil {
		err = nil
		r.locker.Lock(key)
		defer r.locker.Unlock(key)
		err = r.GetObj(ctx, key, obj)
		if err == redis.Nil {
			err = nil
			res, err := getFormDb()
			if err != nil {
				return err
			}
			v := reflect.ValueOf(res)
			o := reflect.ValueOf(obj)
			if o.Kind() != reflect.Pointer {
				panic("data not pointer")
			}
			if v.Kind() == reflect.Pointer {
				if o.Elem().Kind() == reflect.Pointer {
					o.Elem().Set(v)
				} else {
					o.Elem().Set(v.Elem())
				}
			} else {
				if o.Elem().Kind() == reflect.Pointer {
					o.Elem().Elem().Set(v)
				} else {
					o.Elem().Set(v)
				}
			}
			err = r.SetObj(ctx, key, obj, expire)
			if err != nil {
				return err
			}
		}
	}
	return err
}

func (r *Decorator) GetSetObjRefKey(ctx context.Context, key string, obj interface{}, expire time.Duration, getFormDb func() (interface{}, error)) error {
	err := r.GetObj(ctx, key, obj)
	if err == redis.Nil {
		err = nil
		r.locker.Lock(key)
		defer r.locker.Unlock(key)
		err = r.GetObj(ctx, key, obj)
		if err == redis.Nil {
			err = nil
			res, err := getFormDb()
			if err != nil {
				return err
			}
			v := reflect.ValueOf(res)
			if v.Kind() == reflect.Pointer {
				reflect.ValueOf(obj).Elem().Set(v.Elem())
			} else {
				reflect.ValueOf(obj).Elem().Set(v)
			}
			err = r.SetObjRefKey(ctx, key, obj, expire)
			if err != nil {
				return err
			}
		}
	}
	return err
}

func (r *Decorator) Limit(ctx context.Context, key string, duration time.Duration, count int32) (bool, error) {
	res, err := r.Eval(ctx,
		"local times = redis.call('incr',KEYS[1])"+
			" if times == 1 then"+
			"   redis.call('expire',KEYS[1], ARGV[1])"+
			" end"+
			"  if times > tonumber(ARGV[2]) then"+
			"    return 0"+
			"  end"+
			" return 1", []string{key}, []string{strconv.Itoa(int(duration / time.Second)), strconv.Itoa(int(count))}).Int64()
	if err == nil {
		if res == 1 {
			return true, nil
		} else {
			return false, nil
		}
	}
	return false, err
}

func (r *Decorator) SAddEx(ctx context.Context, key, value string, seconds int32) (int64, error) {
	return r.Eval(ctx,
		"if (redis.call('sadd',KEYS[1],ARGV[2]) >= 0) then "+
			"     redis.call('expire',KEYS[1],ARGV[1]) "+
			"     return 1 "+
			"  else "+
			"     return 0 "+
			"end ", []string{key}, []string{strconv.Itoa(int(seconds)), value}).Int64()
}

func (r *Decorator) HSetExEval(ctx context.Context, key, field, value string, seconds int32) error {
	res, err := r.Eval(ctx,
		"if (redis.call('hset',KEYS[1],ARGV[2],ARGV[3]) >= 0) then "+
			"     redis.call('expire',KEYS[1],ARGV[1]) "+
			"     return 1 "+
			"  else "+
			"     return 0 "+
			"end ", []string{key}, []string{strconv.Itoa(int(seconds)), field, value}).Int64()
	if err == nil && res != 1 {
		err = errors.New("hsetex script error")
	}
	return err
}

func (r *Decorator) Del(ctx context.Context, key string) error {
	if r.cnf.LocalCache {
		res := r.Client.Publish(ctx, r.cnf.Chanel, key)
		if res.Err() != nil {
			log.Error(res.Err())
			return res.Err()
		}
	}
	res := r.Client.Del(ctx, key)
	if res.Err() != nil {
		log.Error(res.Err())
	}
	return res.Err()
}

func (r *Decorator) DelKeys(ctx context.Context, key string) {
	pushScript := ""
	if r.cnf.LocalCache {
		pushScript = " redis.call('publish'," + r.cnf.Chanel + ",value) "
	}
	res, err := r.Eval(ctx,
		"local setKey=redis.call('smembers',KEYS[1]) "+
			"for key,value in ipairs(setKey) do "+
			"	redis.call('del',value)  "+
			pushScript+
			"end return 1", []string{key + ":keySet"}).Int64()
	if err == nil && res != 1 {
		err = errors.New("set del keys script error")
	}
}

func (r *Decorator) IncrEx(ctx context.Context, key string, expire, count int32) (int64, error) {
	res, err := r.Eval(ctx,
		"local c=redis.call('incrBy',KEYS[1],ARGV[1]) "+
			"redis.call('expire',KEYS[1],ARGV[2]) "+
			"return c", []string{key}, []string{strconv.Itoa(int(count)), strconv.Itoa(int(expire))}).Int64()
	return res, err
}

// 增加count大于Limit
func (r *Decorator) IncrLimitEx(ctx context.Context, key string, expire, count, limit int32) (int64, error) {
	res, err := r.Eval(ctx,
		"local c=-1 "+
			"if (redis.call('exists',KEYS[1]) == 1) then "+
			"    c=redis.call('incrBy',KEYS[1],ARGV[1]) "+
			"    if ( c < tonumber(ARGV[2])) then "+
			"        redis.call('incrBy',KEYS[1],-ARGV[1]) "+
			"    else "+
			"        redis.call('expire',KEYS[1],ARGV[3]) "+
			"    end "+
			"end "+
			"return c", []string{key}, []string{strconv.Itoa(int(count)), strconv.Itoa(int(limit)), strconv.Itoa(int(expire))}).Int64()
	return res, err
}

// 增加count不能超过than
func (r *Decorator) IncrThanEx(ctx context.Context, key string, expire, count, than int32) (int64, error) {
	res, err := r.Eval(ctx,
		"local c=-1 "+
			"if (redis.call('exists',KEYS[1]) == 1) then "+
			"    c=redis.call('incrBy',KEYS[1],ARGV[1]) "+
			"    if ( c > tonumber(ARGV[2])) then "+
			"        redis.call('incrBy',KEYS[1],-ARGV[1]) "+
			"    else "+
			"        redis.call('expire',KEYS[1],ARGV[3]) "+
			"    end "+
			"end "+
			"return c", []string{key}, []string{strconv.Itoa(int(count)), strconv.Itoa(int(than)), strconv.Itoa(int(expire))}).Int64()
	return res, err
}

func getSetKey(key string) string {
	index := strings.Index(key, "}:")
	if index != -1 {
		return key[0:index+2] + "keySet"
	}
	return ""
}

func (r *Decorator) Lock(ctx context.Context, key string, outTime int32) int64 {
	expire := timeUtil.Milli() + int64(outTime*1000)
	res, err := r.Eval(ctx,
		"if (redis.call('setnx',KEYS[1],ARGV[1]) == 1) then "+
			"   redis.call('expire',KEYS[1],ARGV[2]) "+
			"   return 1 "+
			"end "+
			"if (redis.call('ttl',KEYS[1]) < 0) then "+
			"   if(redis.call('setex',KEYS[1],ARGV[1],ARGV[2]).ok == 'OK') then "+
			"       return 1 "+
			"   end "+
			"end "+
			"return 0", []string{key}, []string{strconv.Itoa(int(expire)), strconv.Itoa(int(outTime))}).Int64()
	if err == nil && res > 0 {
		return expire
	}
	return 0
}

// 单向锁
func (r *Decorator) EitherLock(ctx context.Context, key string, lockBack bool) (string, int64) {
	now := timeUtil.Milli()
	expire := now + timeUtil.Minute
	keys := make([]string, 2)
	if lockBack {
		keys[0] = "{one:way:" + key + "}:back"
		keys[1] = "{one:way:" + key + "}:front"
	} else {
		keys[0] = "{one:way:" + key + "}:front"
		keys[1] = "{one:way:" + key + "}:back"
	}
	res, err := r.Eval(ctx,
		"local setKey=redis.call('smembers',KEYS[1]) "+
			"local isDel=true "+
			"for key,value in ipairs(setKey) do "+
			"    if (tonumber(value)>tonumber(ARGV[1])) then "+
			"        isDel=false "+
			"        break "+
			"    else "+
			"        redis.call('srem',KEYS[1],value) "+
			"    end "+
			"end "+
			"if (isDel) then "+
			"    redis.call('sadd',KEYS[2],ARGV[2]) "+
			"    redis.call('expire',KEYS[2],ARGV[3]) "+
			"    return 1 "+
			"end "+
			"return 0 ", keys, []string{strconv.Itoa(int(now)), strconv.Itoa(int(expire)), strconv.Itoa(int(expire - now/timeUtil.Second))}).Int64()
	if err == nil && res > 0 {
		return keys[1], expire
	}
	return "", 0
}

func (r *Decorator) BfMAdd(ctx context.Context, key string, values ...string) ([]interface{}, error) {
	var args = make([]interface{}, len(values)+2)
	args[0] = "BF.MADD"
	args[1] = key
	for i, val := range values {
		args[i+2] = val
	}
	data, err := r.Do(ctx, args...).Result()
	if err != nil {
		return nil, err
	}
	return data.([]interface{}), err
}

func (r *Decorator) BfMExists(ctx context.Context, key string, values ...string) ([]interface{}, error) {
	var args = make([]interface{}, len(values)+2)
	args[0] = "BF.MEXISTS"
	args[1] = key
	for i, val := range values {
		args[i+2] = val
	}
	data, err := r.Do(ctx, args...).Result()
	if err != nil {
		return nil, err
	}
	return data.([]interface{}), err
}

type XAddArgs struct {
	redis.XAddArgs
	Partition int32 //分区数
}

func (r *Decorator) subscribeLocalCache(ctx context.Context, channels ...string) {
	sub := r.Client.Subscribe(ctx, channels...)
	var failCount int64
	for {
		msg, err := sub.ReceiveMessage(ctx)
		if err != nil {
			log.Error(err)
			failCount++
			if failCount > r.cnf.FailWaitCount {
				time.Sleep(time.Second * time.Duration(r.cnf.FailSleep))
			}
		} else {
			failCount = 0
			r.cache.Del([]byte(msg.Payload))
		}
		re := recover()
		if re != nil {
			log.Error(errors.New(re.(string)))
			failCount++
			if failCount > r.cnf.FailWaitCount {
				time.Sleep(time.Second * time.Duration(r.cnf.FailSleep))
			}
		}
	}
}
