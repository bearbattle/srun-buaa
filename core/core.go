package core

import (
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/vouv/srun/hash"
	"github.com/vouv/srun/model"
	"github.com/vouv/srun/utils"
	"net/url"
	"strings"
)

const (
	baseAddr     = "https://gw.buaa.edu.cn"
	challengeUrl = "/cgi-bin/get_challenge"
	portalUrl    = "/cgi-bin/srun_portal"
	userDmUrl    = "/cgi-bin/rad_user_dm"
	succeedUrl   = "/cgi-bin/rad_user_info"
)

// 获取acid等
func Prepare() (int, error) {
	return 1, nil
	// first, err := get(baseAddr)
	// if err != nil {
	// 	return 1, err
	// }
	// second, err := get(first.Header.Get("Location"))
	// if err != nil {
	// 	return 1, err
	// }
	// target := second.Header.Get("location")
	// query, _ := url.Parse(baseAddr + target)
	// return strconv.Atoi(query.Query().Get("ac_id"))
}

// api Login
// step 1: prepare & get acid
// step 2: get challenge
// step 3: do login
func Login(account *model.Account) (err error) {
	log.Debug("Username: ", account.Username)

	// 先获取acid
	acid, err := Prepare()
	if err != nil {
		log.Debug("prepare error:", err)
		return
	}

	username := account.Username

	// 创建登录表单
	formLogin := model.Login(username, account.Password, acid)

	//	get token
	rc, err := getChallenge(username)
	if err != nil {
		log.Debug("get challenge error:", err)
		return
	}

	token := rc.Challenge
	ip := rc.ClientIp

	formLogin.Set("ip", ip)
	formLogin.Set("info", hash.GenInfo(formLogin, token))
	formLogin.Set("password", hash.PwdHmd5("", token))
	formLogin.Set("chksum", hash.Checksum(formLogin, token))

	// response
	ra := model.ActionResp{}

	if err = utils.GetJson(baseAddr+portalUrl, formLogin, &ra); err != nil {
		log.Debug("request error", err)
		return
	}

	if ra.Res != "ok" {
		log.Debug("response msg is not 'ok'")
		if strings.Contains(ra.ErrorMsg, "Arrearage users") {
			err = errors.New("已欠费")
		} else {
			err = errors.New(fmt.Sprint(ra))
		}
		return
	}

	account.AccessToken = token
	account.Acid = acid
	return
}

// api info
func Info() (info *model.InfoResp, err error) {

	// 余量查询
	err = utils.GetJson(baseAddr+succeedUrl, url.Values{}, &info)
	if err != nil {
		return nil, err
	}
	return
}

// api logout
func Logout(ip string, username string) (err error) {

	q := model.Logout(ip, username)
	ra := model.ActionResp{}
	if err = utils.GetJson(baseAddr+userDmUrl, q, &ra); err != nil {
		log.Debug(err)
		err = ErrRequest
		return
	}
	if ra.Error != "ok" {
		log.Debug(ra)
		err = ErrRequest
	}
	return
}

func getChallenge(username string) (res model.ChallengeResp, err error) {
	qc := model.Challenge(username)
	err = utils.GetJson(baseAddr+challengeUrl, qc, &res)
	return
}
