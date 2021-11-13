package store

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/moby/term"
	log "github.com/sirupsen/logrus"
	"github.com/vouv/srun/model"
	"io"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

const accountFileName = "account.json"

var RootPath string

func GetAccount() (account *model.Account, err error) {
	account, err = ReadAccount()
	if err == nil {
		return
	}
	fmt.Println("没有发现账号配置信息，您可以为本次登录提供账号信息，密码将不会被保存")
	fmt.Println("如需储存账号信息，请使用 `srun config`")
	return ReadAccountFromConsole()
}

func ReadAccountFromConsole() (account *model.Account, err error) {
	in := os.Stdin
	fmt.Print("校园网账号:\n>")
	username := readInput(in)

	// 终端API
	fd, _ := term.GetFdInfo(in)
	oldState, err := term.SaveState(fd)
	if err != nil {
		return nil, err
	}
	fmt.Print("校园网密码(隐私输入):\n>")

	// read in stdin
	_ = term.DisableEcho(fd, oldState)
	pwd := readInput(in)
	_ = term.RestoreTerminal(fd, oldState)

	fmt.Println()

	// trim
	username = strings.TrimSpace(username)
	pwd = strings.TrimSpace(pwd)
	account = &model.Account{Username: username, Password: pwd}
	return
}

func readInput(in io.Reader) string {
	reader := bufio.NewReader(in)
	line, _, err := reader.ReadLine()
	if err != nil {
		panic(err)
	}
	return string(line)
}

func ReadAccount() (account *model.Account, err error) {
	file, err := OpenAccountFile(os.O_RDONLY)
	if err != nil {
		log.Debugf("打开账号文件错误, %s,", err)
		return
	}
	defer file.Close()

	err = json.NewDecoder(base64.NewDecoder(base64.RawStdEncoding, file)).Decode(&account)
	if account.Password == "" {
		err = errors.New("密码为空")
	}
	return
}

func OpenAccountFile(flag int) (file *os.File, err error) {
	accountFilename, err := getAccountFilename()
	if err != nil {
		return
	}
	return os.OpenFile(accountFilename, flag, 0600)
}

func WriteAccount(account *model.Account) (err error) {
	file, err := OpenAccountFile(os.O_CREATE | os.O_TRUNC | os.O_WRONLY)
	if err != nil {
		log.Debugf("打开账号文件错误, %s", err)
		return
	}

	defer file.Close()

	enc := base64.NewEncoder(base64.RawStdEncoding, file)
	err = json.NewEncoder(enc).Encode(account)
	if err != nil {
		return err
	}
	return enc.Close()
}

func getAccountFilename() (fileSrc string, err error) {
	storageDir := filepath.Join(RootPath, ".srun")
	if _, sErr := os.Stat(storageDir); sErr != nil {
		if mErr := os.MkdirAll(storageDir, 0755); mErr != nil {
			log.Debugf("mkdir `%s` error, %s", storageDir, mErr)
			return
		}
	}
	fileSrc = filepath.Join(storageDir, accountFileName)
	return
}

func init() {
	curUser, gErr := user.Current()
	if gErr != nil {
		log.Fatalln("无法读取账户信息, 请重新设置账户信息")
	} else {
		RootPath = curUser.HomeDir
	}
}
