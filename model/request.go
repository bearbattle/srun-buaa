package model

import (
	"fmt"
	"github.com/vouv/srun/hash"
	"net/url"
	"strconv"
	"time"
)

func Challenge(username string) url.Values {
	return url.Values{
		"username": {username},
		"ip":       {""},
	}
}

func Login(username, password string, acid int) url.Values {
	return url.Values{
		"action":   {"login"},
		"username": {username},
		"password": {password},
		"ac_id":    {fmt.Sprint(acid)},
		"ip":       {""},
		"info":     {},
		"chksum":   {},
		"n":        {"200"},
		"type":     {"1"},
	}
}

func Logout(ip string, username string) url.Values {
	form := url.Values{
		"ip":       {ip},
		"username": {username},
		"time":     {strconv.FormatInt(time.Now().Unix(), 10)},
		"unbind":   {"0"},
	}
	form.Add("sign", hash.GetLogoutSign(&form))
	return form
}
