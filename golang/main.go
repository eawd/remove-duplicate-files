package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

var hashes = make(map[string]string);
var lock = sync.RWMutex{};
// Fill in with a folder to keep duplicate files.
var temporaryFolder = "";
// Preferred location to keep.
var desiredLocations = []string {
};

var folderToScan = "";
var counter = 0;

func isDesired(filePath string) bool {
	for _, location := range desiredLocations {
		if matched, _ := regexp.Match(location, []byte(filePath)); matched {
			return true;
		}
	}

	return false;
}

func removeFile(filePath string) {
	fmt.Printf("Removing: %s\n", filePath);
	newPath := path.Join(temporaryFolder, strings.Replace(filePath, "F:\\", "", 1));
	parentPath := filepath.Dir(newPath);
	if err := os.MkdirAll(parentPath, 0777); err != nil {
		fmt.Println(err);
		return;
	}

	err := os.Rename(filePath, newPath);
	if err != nil {
		fmt.Println(err);
		return;
	}
}

func checkFile(wg *sync.WaitGroup,filePath string) {
	// wg.Add(1);
	// defer wg.Done();

	file, err := os.Open(filePath);
	if err != nil {
		fmt.Println(err);
		return;
	}

	hasher := md5.New();
	if _, err := io.Copy(hasher, file); err != nil {
		fmt.Println(err);
		file.Close();
		return;
	}

	file.Close();

	hash := hex.EncodeToString(hasher.Sum(nil)[:16]);
	isCurrentFileDesired := isDesired(filePath);

	// lock map writing.
	lock.Lock();
	defer lock.Unlock();
	if previousFile, ok := hashes[hash]; ok {
		fileToRemove := filePath;
		fmt.Printf("Found duplicate: %s, \t %s\n", filePath, previousFile);

		if (isCurrentFileDesired) {
			fileToRemove = previousFile;
			hashes[hash] = filePath;
		}

		removeFile(fileToRemove);
	} else {
		hashes[hash] = filePath;
	}

	if counter++; counter % 200 == 0 {
		fmt.Printf("Scanned %d files, currently on: %s\n", counter, filePath);
	}
}

func scanFolder(wg *sync.WaitGroup,directory string) {
	wg.Add(1);
	defer wg.Done();

	files, err := os.ReadDir(directory);
	if err != nil {
		fmt.Println(err);
		return;
	}

	for _, f := range files {
		path := path.Join(directory, f.Name());
		if (f.IsDir()) {
			go scanFolder(wg, path);
		} else {
			checkFile(wg, path);
		}
	}
}

func main() {
	var wg sync.WaitGroup;

	scanFolder(&wg, folderToScan);

	wg.Wait();
	fmt.Println("Done");
}
