package main

import (
	"fmt"
	"runtime"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/vouv/srun/core"
	"github.com/vouv/srun/store"
)

func Login(cmd *cobra.Command, args []string) {
	err := LoginE(cmd, args)
	if err != nil {
		log.Error(err)
	}
}

func LoginE(cmd *cobra.Command, args []string) error {
	account, err := store.GetAccount()
	if err != nil {
		return err
	}
	log.Info("尝试登录...")

	if err = core.Login(account); err != nil {
		return err
	}
	log.Info("登录成功!")

	return nil
}

func Logout(cmd *cobra.Command, args []string) {
	err := LogoutE(cmd, args)
	if err != nil {
		log.Error(err)
	}
}

func LogoutE(cmd *cobra.Command, args []string) error {
	info, err := core.Info()
	if err != nil {
		return err
	}
	username := info.UserName
	if username == "" {
		log.Info("账号未登录")
	} else {
		_ = core.Logout(username)
		log.Info("注销成功!")
	}
	return nil
}

func Info(cmd *cobra.Command, args []string) {
	err := InfoE(cmd, args)
	if err != nil {
		log.Error(err)
	}
}

func InfoE(cmd *cobra.Command, args []string) error {
	info, err := core.Info()
	if err != nil {
		return err
	}
	fmt.Println(info.String())
	return nil
}

func Config(cmd *cobra.Command, args []string) {
	err := ConfigE(cmd, args)
	if err != nil {
		log.Error(err)
	}
}

func ConfigE(cmd *cobra.Command, args []string) error {
	account, err := store.ReadAccountFromConsole()
	if err != nil {
		return err
	}
	if err := store.WriteAccount(account); err != nil {
		return err
	}
	log.Info("账号密码已被保存")
	return nil
}

func VersionString() string {
	return fmt.Sprintln("System:") +
		fmt.Sprintf("\tOS:%s ARCH:%s GO:%s\n", runtime.GOOS, runtime.GOARCH, runtime.Version()) +
		fmt.Sprintln("About:") +
		fmt.Sprintf("\tVersion: %s\n", Version) +
		fmt.Sprintln("\n\t</> with ❤ By vouv")
}

func Update(cmd string, params ...string) {
	ok, v, d := HasUpdate()
	if !ok {
		log.Info("当前已是最新版本:", Version)
		return
	}
	log.Info("发现新版本: ", v, "当前版本: ", Version)
	log.Info("打开链接下载: ", d)
}
