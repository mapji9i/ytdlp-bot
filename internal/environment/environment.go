package environment

import (
	"log"

	"github.com/kelseyhightower/envconfig"
)

type environment struct {
	TelegramToken string `envconfig:"telegram_token" required:"true"`
	GoogleToken   string `envconfig:"google_token" required:"true"`
	ChatId        int64  `envconfig:"chat_id" required:"true"`
	DBFilePath    string `envconfig:"db_file" required:"true"`
	DownloadRoot  string `envconfig:"download_root" required:"true"`
	WorkingDir    string `envconfig:"working_dir" required:"true"`
	RunScriptName string `envconfig:"run_script_name" required:"true"`
	RootUserId    int64  `envconfig:"root_user_id" required:"true"`
	Stage         string `envconfig:"stage" required:"true" default:"DEBUG"`
}

var Environment environment

func GetEnvironment() environment {
	err := envconfig.Process("myapp", &Environment)
	if err != nil {
		log.Fatal(err.Error())
	}
	return Environment
}
