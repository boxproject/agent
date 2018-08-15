package commands

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path"
	"path/filepath"
	"runtime"

	"github.com/AlecAivazis/survey"
	logger "github.com/alecthomas/log4go"
	"github.com/boxproject/agent/config"
)

const (
	ServiceName_KeyStore              = "/agent/keystore"        //公钥
	ServiceName_Operate               = "/agent/operate"         //操作
	ServiceName_Status                = "/agent/status"          //查询状态
	ServiceName_Allow                 = "/agent/allow"           //同意
	ServiceName_DisAllow              = "/agent/disallow"        //拒绝
	ServiceName_Regist_Add            = "/agent/registadd"       //添加注册
	ServiceName_Regist_Aproval        = "/agent/registaproval"   //注册审批
	ServiceName_Regist_List           = "/agent/registlist"      //查询注册信息
	ServiceName_Approval_Add          = "/agent/approvaladd"     //审批流新增
	ServiceName_Approval_Invalid      = "/agent/approvalinvalid" //作废审批
	ServiceName_Approval_List         = "/agent/approvallist"    //审批流列表
	ServiceName_Approval_Detail       = "/agent/approvaldetail"  //审批流详情
	ServiceName_Approval_Operate_List = "/agent/approvaloplist"  //审批流操作列表
	//ServiceName_Hash_Add        = "/agent/hashadd"        //查询hash list
	//ServiceName_Hash_List       = "/agent/hashlist"       //查询hash list
	ServiceName_Token_Add     = "/agent/tokenedit"     //添加token
	ServiceName_Token_Del     = "/agent/tokendel"      //删除token
	ServiceName_Token_List    = "/agent/tokenlist"     //查询token list
	ServiceName_Coin          = "/agent/coin"          //coin
	ServiceName_Coin_List     = "/agent/coinlist"      //coin list
	ServiceName_Wtihdraw      = "/agent/wtihdraw"      //提现
	ServiceName_Manager_Info  = "/agent/msinfo"        //管理端信息
	ServiceName_Assets        = "/agent/assets"        // 查询余额
	ServiceName_Trade_History = "/agent/trade/history" // 查询交易流水
)

var (
	rootPath string
	filePath string

	qs = []*survey.Question{
		{
			Name: "password",
			Prompt: &survey.Password{
				Message: "Input password: ",
			},
			Validate: survey.Required,
		},
		{
			Name: "passwordConfirm",
			Prompt: &survey.Password{
				Message: "Input password again: ",
			},
			Validate: survey.Required,
		},
	}

	ErrAESTextSize = errors.New("ciphertext is not a multiple of the block size")
	ErrAESPadding  = errors.New("cipher padding size error")
)

type answers struct {
	Passphrase        string `survey:"passphrase"`
	PassphraseConfirm string `survey:"passphraseConfirm"`
	Password          string `survey:"password"`
	Confirm           string `survey:"passwordConfirm"`
}

func init() {
	main, _ := exec.LookPath(os.Args[0])
	file, _ := filepath.Abs(main)
	rootPath = path.Dir(file)
}

func GetFilePath() string {
	return filePath
}

func DefaultConfigDir() string {
	home := homeDir()
	if home != "" {
		if runtime.GOOS == "darwin" {
			return filepath.Join(home, ".bcmonitor")
		} else if runtime.GOOS == "windows" {
			return filepath.Join(home, "AppData", "Roaming", "bcmonitor")
		} else {
			return filepath.Join(home, ".bcmonitor")
		}
	}

	return ""
}

func LoadConfig(configPath, defaultFileName string) (*config.Config, error) {
	configPath = GetConfigFilePath(configPath, defaultFileName)

	logger.Debug("config path: %s", configPath)
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var cfg config.Config
	if err = json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// configPath 不为空时，不检查fileName
func GetConfigFilePath(configPath, defaultFileName string) string {
	for i := 0; i < 3; i++ {
		if configPath != "" {
			if _, err := os.Stat(configPath); !os.IsNotExist(err) {
				break
			}
		}

		if i == 0 {
			configPath = path.Join(GetFilePath(), defaultFileName)
		} else if i == 1 {
			configPath = path.Join(DefaultConfigDir(), defaultFileName)
		}
	}

	return configPath
}

func homeDir() string {
	if home := os.Getenv("HOME"); home != "" {
		return home
	}

	if usr, err := user.Current(); err == nil {
		return usr.HomeDir
	}

	return ""
}

func InitLogger() {
	logFile := path.Join(rootPath, "log.xml")
	for i := 0; i < 3; i++ {
		if _, err := os.Stat(logFile); !os.IsNotExist(err) {
			break
		}
		if i == 0 {
			logFile = path.Join(filePath, "log.xml")
		} else if i == 1 {
			logFile = path.Join(DefaultConfigDir(), "log.xml")
		}
	}
	logger.LoadConfiguration(logFile)
}

// AES解密
func aesDecrypt(password, src []byte) ([]byte, error) {
	// 长度不能小于aes.Blocksize
	if len(src) < aes.BlockSize*2 || len(src)%aes.BlockSize != 0 {
		return nil, ErrAESTextSize
	}

	padLen := aes.BlockSize - (len(password) % aes.BlockSize)
	for i := 0; i < padLen; i++ {
		password = append(password, byte(padLen))
	}

	aesBlock, err := aes.NewCipher(password)
	if err != nil {
		return nil, err
	}

	srcLen := len(src) - aes.BlockSize
	decryptText := make([]byte, srcLen)
	iv := src[srcLen:]

	mode := cipher.NewCBCDecrypter(aesBlock, iv)
	mode.CryptBlocks(decryptText, src[:srcLen])
	paddingLen := int(decryptText[srcLen-1])

	if paddingLen > 16 {
		return nil, ErrAESPadding
	}

	return decryptText[:srcLen-paddingLen], nil
}

// AES加密
func aesEncrypt(password, src []byte) ([]byte, error) {
	padLen := aes.BlockSize - (len(src) % aes.BlockSize)
	for i := 0; i < padLen; i++ {
		src = append(src, byte(padLen))
	}

	padLen = aes.BlockSize - (len(password) % aes.BlockSize)
	for i := 0; i < padLen; i++ {
		password = append(password, byte(padLen))
	}

	aesBlock, err := aes.NewCipher(password)
	if err != nil {
		fmt.Printf("aes new cipher error: %v\n", err)
		return nil, err
	}

	srcLen := len(src)
	encryptText := make([]byte, srcLen+aes.BlockSize)
	iv := encryptText[srcLen:]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	mode := cipher.NewCBCEncrypter(aesBlock, iv)
	mode.CryptBlocks(encryptText[:srcLen], src)

	return encryptText, nil
}
