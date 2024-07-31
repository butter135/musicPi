package main

import (
	"io"
	"log"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	lockFilePath := "/home/musicPi/.lock"
	// ファイルロックのためのファイルを作成
	lockFile, err := os.Create(lockFilePath)
	if err != nil {
		log.Fatal("Another instance is already running.")
	}
	defer lockFile.Close()

	// ファイルロックを取得
	err = syscall.Flock(int(lockFile.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if err != nil {
		log.Fatal("Another instance is already running.")
		return
	}
	defer func() {
		// ロック解除
		syscall.Flock(int(lockFile.Fd()), syscall.LOCK_UN)
		// ロックファイルの削除
		os.Remove(lockFilePath)
	}()
	for {
		url, filePath := dequeue()
		// コマンド名と引数を指定
		cmd := exec.Command("mpv", "--no-video", "--ao=pulse", "--ytdl-format=bestaudio", url)

		// コマンドの実行
		err := cmd.Run()

		if err != nil {
			if _, ok := err.(*exec.ExitError); ok {
				os.Remove(filePath)
				continue
			}
			log.Fatal("コマンドの実行時にエラーが発生しました:", err)
		}
		os.Remove(filePath)
	}
}

func dequeue() (string, string) {
	// ディレクトリのパスを指定
	dirPath := "/home/musicPi/.queue/"

	// ディレクトリを開く
	dir, err := os.Open(dirPath)
	if err != nil {
		log.Fatal(err)
	}
	defer dir.Close()

	// ディレクトリ内の最初のファイルを取得
	files, err := dir.Readdir(0)
	if err != nil {
		log.Fatal(err)
	}

	// 最初のパスを取得
	if len(files) > 0 {
		filePath := dirPath + files[0].Name()
		file, err := os.Open(filePath)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		// ファイルの中身を読み取る
		content, err := io.ReadAll(file)
		if err != nil {
			log.Fatal(err)
		}
		return string(content), filePath

	} else {
		log.Fatal("No files found in", dirPath)
	}

	return "none", "none"
}
