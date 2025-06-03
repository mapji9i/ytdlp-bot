package executors

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
	"ytdlp-bot/internal/data"
	"ytdlp-bot/internal/environment"
	YtBot "ytdlp-bot/internal/ytBot"

	"github.com/thinhdanggroup/executor"
)

var (
	executorService *executor.Executor
	tryCounters     map[string]int8
	messageId       atomicInt
	mapInWork       map[string]string
	downloadRoot    string
	workDir         string
)

const tryLimit = 3

func StartReader() {
	editDownloadList()
	data.ResetInWorkForAllEntities()

	tryCounters = make(map[string]int8)
	mapInWork = make(map[string]string)
	if os.Getenv("USER") == "root" {
		exec.Command("wg-quick", "up", "yt-dlp").Start()
	}
	ctx := YtBot.GetCtx()
	//go scheduleCleanMessage(ctx)

	downloadRoot = environment.Environment.DownloadRoot
	workDir = environment.Environment.WorkingDir

	for {
		entity, err := data.GetNotUsedEntityInWork()
		_, ok := tryCounters[entity.VideoID]
		if !ok {
			tryCounters[entity.VideoID] = 0
		}
		if err != nil {
			time.Sleep(time.Second)
			continue
		}

		executorService.Publish(func() {

			mapInWork[entity.VideoTitle] = "⏬"
			editDownloadList()

			tryCounters[entity.VideoID] = tryCounters[entity.VideoID] + 1

			path := createRequiredDirs(entity)

			downloadFileName := path + entity.VideoTitle + ".mp4"

			cmd := exec.Command("yt-dlp", "-t", "mp4", "-U", "-o", path+"%(title&"+entity.VideoTitle+")s.%(ext)s",
		
				"--cookies", "./cookies.txt", "https://www.youtube.com/watch?v="+entity.VideoID)
			log.Print(cmd.String())
			logFilePattern := "%s/logs/" + time.Now().Local().Format("2006:01:02 15:04:05") + "_%s_%s.txt"

			outFilePath := fmt.Sprintf(logFilePattern, workDir, entity.VideoID, "out")

			errFilePath := fmt.Sprintf(logFilePattern, workDir, entity.VideoID, "err")

			outPipe, err := cmd.StdoutPipe()
			if err != nil {
				log.Panic("Out pipe creation error")
			}
			errPipe, err := cmd.StderrPipe()
			if err != nil {
				log.Panic("Error pipe creation error")
			}

			go startLogThread(ctx, outFilePath, outPipe)

			go startLogThread(ctx, errFilePath, errPipe)

			cmd.Run()

			isSuccess := checkSuccessDownloadFile(downloadFileName, entity, errFilePath)

			if isSuccess && strings.Contains(entity.PlaylistTitle, "Shorts") {
				YtBot.SendVideo(entity.VideoTitle, entity.VideoTitle+".mp4", downloadFileName)
			}
		})
	}

}

func createRequiredDirs(entity data.DatabaseEntity) string {
	var path string
	if entity.PlaylistTitle == "Shorts" {
		path = fmt.Sprintf("%s/Youtube/%s/", downloadRoot, "!Shorts")
	} else {
		path = fmt.Sprintf("%s/Youtube/%s/%s/", downloadRoot, entity.ChannelTitle, entity.PlaylistTitle)
	}
	makeDirIfNotExists(path)

	logDirPath := workDir + "/logs/"
	makeDirIfNotExists(logDirPath)
	return path
}

func GetService() *executor.Executor {
	var err error
	if executorService == nil {
		executorService, err = executor.New(executor.Config{
			ReqPerSeconds: 0,
			QueueSize:     8,
			NumWorkers:    runtime.NumCPU() - 1,
		})
		if err != nil {
			log.Panic(err)
		}
	}
	return executorService
}

func checkSuccessDownloadFile(fileName string, entity data.DatabaseEntity, errFilePath string) bool {
	var err error
	_, err = os.Stat(fileName)

	if err == nil && !errors.Is(err, fs.ErrNotExist) {
		log.Print(entity.VideoTitle + " Downloaded normal")
		data.RemoveEntityFromDB(entity)
		mapInWork[entity.VideoTitle] = "✅"
		editDownloadList()
		return true
	} else if tryLimit == int(tryCounters[entity.VideoID]) {
		data.RemoveEntityFromDB(entity)
		delete(tryCounters, entity.VideoID)
		log.Print(entity.VideoTitle + " No downloaded, try limit achieved")
		mapInWork[entity.VideoTitle] = "❌"
		YtBot.SendDocument("Ошибка скачивания видео"+entity.VideoTitle, "Stderr.txt", errFilePath)
		editDownloadList()
	} else {
		data.ResetInWorkForEntity(entity)
		log.Print(entity.VideoTitle + " No downloaded, retry")
	}
	return false
}
func trimMapOfTasks() {
	for key, value := range mapInWork {
		if value != "⏬" {
			delete(mapInWork, key)
		}
		if len(mapInWork) == 10 {
			break
		}
	}
}
func editDownloadList() {
	messageOldId := messageId.GetValue()
	if len(mapInWork) > 10 {
		trimMapOfTasks()
	}
	counter := 0
	textWithList := strings.Builder{}
	textWithList.WriteString("Список загрузки:" + "\n")
	if len(mapInWork) == 0 {
		textWithList.WriteString("_пусто_")
	} else {
		for key, value := range mapInWork {
			counter++
			textWithList.WriteString(fmt.Sprintf("%d", counter) + ". " + key + " " + value + "\n")
			if counter == 10 {
				break
			}
		}
	}
	messageId.SetValue(YtBot.EditMessageCreateIfNotExist(messageId.GetValue(), textWithList.String()))
	messageNewId := messageId.GetValue()
	if messageOldId != messageNewId {
		YtBot.PinMessage(messageNewId)
	}
}

func startLogThread(ctx context.Context, filePath string, outStream io.ReadCloser) {

	outfile, err := os.Create(filePath)

	if err != nil {
		log.Printf("Ошибка создания файла %s", filePath)
	}
	defer outfile.Close()

	writerOut := bufio.NewWriter(outfile)

	io.Copy(writerOut, outStream)
}

func makeDirIfNotExists(path string) {
	_, err := os.Stat(path)
	if err != nil && errors.Is(err, fs.ErrNotExist) {
		os.MkdirAll(path, os.ModePerm)
	}
}
