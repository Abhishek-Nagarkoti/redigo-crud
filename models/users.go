package models

import (
	// "github.com/Abhishek-Nagarkoti/redigo-crud/lib"
	"github.com/garyburd/redigo/redis"
	"log"
	"reflect"
	"regexp"
	// "strconv"
)

type User struct {
	Id        string `redis:"_id" json:"_id"`
	FirstName string `redis:"first_name" json:"first_name"`
	LastName  string `redis:"last_name" json:"last_name"`
	Gender    string `redis:"gender" json:"gender"`
	Yoo       string `redis:"yoo" json:"yoo"`
	Hoo       string `redis:"hoo" json:"hoo"`
}

func (u User) Create(DB redis.Conn) (error, User) {
	sv := reflect.ValueOf(u)
	st := reflect.TypeOf(u)
	var err error
	err = nil
	for i := 0; i < sv.NumField(); i++ {
		tag := st.Field(i).Tag.Get("redis")
		val := sv.Field(i).Interface().(string)
		if i == 0 {
			_, err = DB.Do("HMSET", "user:"+u.Id, tag, val)
			if err != nil {
				_, _ = DB.Do("DEL", "user:"+u.Id)
				break
			}
		} else {
			_, err = DB.Do("HSETNX", "user:"+u.Id, tag, val)
			if err != nil {
				_, _ = DB.Do("DEL", "user:"+u.Id)
				break
			}
		}
	}
	if err == nil {
		return nil, u
	} else {
		return err, u
	}
}

func (u User) Get(DB redis.Conn, id string) (error, User) {
	reply, err := redis.Values(DB.Do("HGETALL", "user:"+id))
	if err != nil {
		return err, u
	} else {
		if err = redis.ScanStruct(reply, &u); err != nil {
			return err, u
		} else {
			return nil, u
		}
	}
}

func (u User) Update(DB redis.Conn) (error, User) {
	_, err := redis.String(DB.Do("HMSET", "user:"+u.Id, "_id", u.Id, "first_name", u.FirstName, "last_name", u.LastName, "gender", u.Gender))
	if err != nil {
		return err, u
	} else {
		reply, err := redis.Values(DB.Do("HGETALL", "user:"+u.Id))
		if err != nil {
			return err, u
		} else {
			if err = redis.ScanStruct(reply, &u); err != nil {
				return err, u
			} else {
				return nil, u
			}
		}
	}
}

func (u User) Delete(DB redis.Conn) error {
	_, err := DB.Do("DEL", "user:"+u.Id)
	if err != nil {
		return err
	} else {
		return nil
	}
}

func (u User) GetALL(DB redis.Conn) (error, []User) {
	users := []User{}
	values, err := redis.Values(DB.Do("KEYS", "user:*"))
	if err != nil {
		return err, users
	} else {
		for i := 0; i < len(values); i += 1 {
			reply, err := redis.Values(DB.Do("HGETALL", values[i]))
			if err != nil {
				//Do nothing
			} else {
				if err = redis.ScanStruct(reply, &u); err != nil {
					//Do nothing
				} else {
					users = append(users, u)
				}
			}
		}
		return nil, users
	}
}

func (u User) Find(DB redis.Conn) (error, []User) {
	users := []User{}
	values, err := redis.Values(DB.Do("KEYS", "user:*"))
	if err != nil {
		return err, users
	} else {
		for i := 0; i < len(values); i += 1 {
			reply, err := redis.Values(DB.Do("HGETALL", values[i]))
			if err != nil {
				//Do nothing
			} else {

				user := User{}
				if err = redis.ScanStruct(reply, &user); err != nil {
					//Do nothing
				} else {
					matched, err := regexp.MatchString(u.FirstName+".*", user.FirstName)
					if err != nil {
						//Do Nothing
					} else {
						if matched {
							users = append(users, user)
						} else {
							//Do Nothing
						}
					}
				}
			}
		}
		return nil, users
	}
}

func (u User) Automigration(DB redis.Conn) error {
	values, err := redis.Values(DB.Do("KEYS", "user:*"))
	if err != nil {
		return err
	} else {
		if len(values) > 0 {
			tags := u.Got()
			var notfound []string
			for _, num := range tags {
				val, _ := redis.Int(DB.Do("HEXISTS", values[0], num))
				if val == 0 {
					notfound = append(notfound, num)
				}
			}
			if len(notfound) > 0 {
				u.NotFound(values, notfound, DB)
			}
			useless := u.Useless(values[0], tags, DB)
			err = u.RemoveUseless(values, useless, DB)
			if err == nil {
				return nil
			} else {
				return err
			}
		} else {
			return nil
		}
	}
}

func (u User) Got() []string {
	var tags []string
	sv := reflect.ValueOf(u)
	st := reflect.TypeOf(u)
	for i := 0; i < sv.NumField(); i++ {
		tags = append(tags, st.Field(i).Tag.Get("redis"))
	}
	return tags
}

func (u User) NotFound(values []interface{}, data []string, DB redis.Conn) {
	for _, value := range values {
		for _, num := range data {
			_, _ = DB.Do("HSETNX", value, num, "")
		}
	}
}

func (u User) Useless(value interface{}, data []string, DB redis.Conn) []string {
	var keys []string
	reply, err := redis.Values(DB.Do("HGETALL", value))
	if err != nil {
		//Do nothing
	} else {
		if err = redis.ScanStruct(reply, &u); err != nil {
			//Do nothing
		} else {
			for i := 0; i < len(reply); i += 2 {
				abort := false
				b := make([]byte, len(reply[i].([]uint8)))
				for i, v := range reply[i].([]uint8) {
					b[i] = byte(v)
				}
				for _, tag := range data {
					if tag == string(b) {
						abort = true
					}
				}
				if !abort {
					keys = append(keys, string(b))
				}
			}
		}
	}
	return keys
}

func (u User) RemoveUseless(values []interface{}, useless []string, DB redis.Conn) error {
	for i := 0; i < len(values); i += 1 {
		for _, val := range useless {
			_, _ = redis.Int(DB.Do("HDEL", values[i], val))
		}
	}
	return nil
}
