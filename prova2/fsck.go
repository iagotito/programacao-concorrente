// Tirei 8 com essa solução

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"sync"
	"time"
)

func check(err error) {
    if err != nil {
        log.Fatal(err)
    }
}

var filesCh chan string
var dirsCh chan string

var wg sync.WaitGroup

var fsckedDirsCount int
var fsckedFilesCount int
var damagedDirsCount int
var damagedFilesCount int

var fsckedDirsCountCh chan struct{}
var fsckedFilesCountCh chan struct{}
var damagedDirsCountCh chan struct{}
var damagedFilesCountCh chan struct{}

func main() {
    args := os.Args[1:]
    dirName := args[0]

    dirPath, err := filepath.Abs(dirName)
    check(err)

    fsckedFilesCountCh = make(chan struct{}, 1)
    damagedFilesCountCh = make(chan struct{}, 1)
    fsckedDirsCountCh = make(chan struct{}, 1)
    damagedDirsCountCh = make(chan struct{}, 1)
    go func() {
        for {
            select {
            case <- fsckedFilesCountCh:
                fsckedFilesCount++
                wg.Done()
            case <- damagedFilesCountCh:
                damagedFilesCount++
            case <- fsckedDirsCountCh:
                fsckedDirsCount++
                wg.Done()
            case <- damagedDirsCountCh:
                damagedDirsCount++
            }
        }
    }()

    go func() {
        for {
            time.Sleep(time.Second)
            fmt.Printf("fscked_files %d damaged_files %d fscked_dirs %d damaged_dirs %d\n", fsckedFilesCount, damagedFilesCount, fsckedDirsCount, damagedDirsCount)
        }
    }()

    var files []string
    loadFiles(dirPath, &files)
    wg.Add(len(files))

    filesCh = make(chan string, 1)
    dirsCh = make(chan string, 1)

    go handleFiles(&files)
    go handleDirs(dirPath)

    wg.Wait()
    fmt.Printf("fscked_files %d damaged_files %d fscked_dirs %d damaged_dirs %d\n", fsckedFilesCount, damagedFilesCount, fsckedDirsCount, damagedDirsCount)
}

func handleFiles(files *[]string) {
    fsckingFilesCh := make(chan struct{}, 8)
    for _, file := range *files {
        fsckingFilesCh <- struct{}{}
        go handleHandleFile(file, fsckingFilesCh)
    }
}

// perdão pelo nome, não pensei em nada melhor
func handleHandleFile(file string, fsckingFilesCh chan struct{}) {
    damaged := fsckFile(file)

    if damaged {
        damagedFilesCountCh <- struct{}{}
        wg.Add(1)
        dirsCh <- parent(file)
    }

    fsckedFilesCountCh <- struct{}{}
    <- fsckingFilesCh
}

func handleDirs(rootDirPath string) {
    fsckingDirsCh := make(chan struct{}, 8)
    fsckedDirsSet := make(map[string]bool)

    for dir := range dirsCh {
        _, ok := fsckedDirsSet[dir]
        // adiciona aos checados
        fsckedDirsSet[dir] = true
        // se o dir ainda não foi checado, cheque
        if !ok {
            fsckingDirsCh <- struct{}{}
            go handleHandleDir(dir, rootDirPath, fsckingDirsCh)
        } else {
            wg.Done()
        }
    }
}

// perdão pelo nome, não pensei em nada melhor
func handleHandleDir(dir, rootDirPath string, fsckingDirsCh chan struct{}) {
    damaged := fsckFile(dir)

    // tem que tirar desse canal antes de adicionar outro dir no canal de
    // diretorios, mas o fsckFile pro diretório já finalizou nesse ponto, então
    // não terão mais do que 8
    <- fsckingDirsCh

    if damaged {
        damagedDirsCountCh <- struct{}{}
        // não olhar acima do root passado
        if dir != rootDirPath {
            wg.Add(1)
            dirsCh <- parent(dir)
        }
    }

    fsckedDirsCountCh <- struct{}{}
}

func fsckFile(filePath string) bool {
    //rSleep := rand.Intn(3)
    rSleep := 0
    //rSleep := 1
    time.Sleep(time.Duration(rSleep) * time.Second)

    return isDamaged()
}

// 1/3 dos arquivos estarão corrompidos
func isDamaged() bool {
    // forçando todos serem corrompidos para debugar
    //return true
    rand.Seed(time.Now().UnixNano())
    num := rand.Intn(2)
    if num == 0 {
        return true
    }
    return false
}

func parent(filePath string) string {
    return filepath.Dir(filePath)
}

// DFS no diretório
func loadFiles(dirPath string, files *[]string) {
    dirFiles, err := ioutil.ReadDir(dirPath)
    check(err)

    for _, file := range(dirFiles) {
        filePath := fmt.Sprintf("%v/%v", dirPath, file.Name())
        if file.IsDir() {
            loadFiles(filePath, files)
        } else {
            *files = append(*files, filePath)
        }
    }
}
