package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
)

func parseData(filePath string) map[string]string {
	f, e := os.Open(filePath)
	if nil != e {
		_, _ = fmt.Fprintf(os.Stderr, "Could not open file %s\n", filePath)
		return nil
	}
	defer func() {
		_ = f.Close()
	}()
	reader := bufio.NewReader(f)
	item := make(map[string]string)
	for {
		bytes, _, e := reader.ReadLine()
		if nil != e && e == io.EOF {
			break
		}
		if nil != e {
			_, _ = fmt.Fprintf(os.Stderr, "Could not read file %s\n", filePath)
			continue
		}
		line := string(bytes)
		line = strings.TrimSpace(line)
		parts := strings.Split(line, "=")
		if len(parts) != 2 {
			continue
		}
		item[parts[0]] = parts[1]
	}
	return item
}

// copy file
func copyFile(src, dest string) bool {
	sf, err := os.Open(src)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to open %s for reading\n", src)
		return false
	}
	defer func() {
		_ = sf.Close()
	}()
	wf, err := os.Create(dest)
	if nil != err {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to open %s for writing\n", dest)
		return false
	}
	defer func() {
		err = wf.Sync()
		if nil != err {
			_, _ = fmt.Fprintf(os.Stderr, "Failed to fush file %s\n", dest)
			return
		}
		_ = wf.Close()
	}()
	_, err = io.Copy(wf, sf)
	if nil != err {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to copy content from %s to %s\n", src, dest)
	}
	return true
}
func processItem(item map[string]string, filePath string) {
	v, ok := item["deleted"]
	//ignore deleted items
	if ok && v == "true" {
		return
	}
	filePathKey := "@BlobStore.blob-name"
	propFilePath := item[filePathKey]
	parts := strings.Split(propFilePath, ".")
	ext := parts[len(parts)-1]
	if ext != "jar" && ext != "pom" {
		return
	}
	folder := path.Dir(propFilePath)
	err := os.MkdirAll(folder, 0755)
	if nil != err {
		_, _ = fmt.Fprintf(os.Stderr, "Could not create folder %s: %v\n", folder, err)
		return
	}
	fileName := path.Base(propFilePath)
	// physical file path
	srcFileNameParts := strings.Split(path.Base(filePath), ".")
	bytesFileName := strings.Join(srcFileNameParts[0:len(srcFileNameParts)-1], ".") + ".bytes"
	bytesFilePath := path.Join(path.Dir(filePath), bytesFileName)
	// target file path
	targetFilePath := path.Join(folder, fileName)
	fmt.Printf("Copying %s to %s\n", bytesFilePath, targetFilePath)
	// copy the file
	copyFile(bytesFilePath, targetFilePath)
	if ext == "jar" {
		folderParts := strings.Split(folder, "/")
		groupName := strings.Join(folderParts[0:len(folderParts)-2], ".")
		artifactName := folderParts[len(folderParts)-2]
		versionNo := folderParts[len(folderParts)-1]
		cmd := fmt.Sprintf("~/nexus-clone/clone-tool-linux-amd64 %s %s %s %s", groupName, artifactName, versionNo, targetFilePath)
		f, e := os.OpenFile("call-generate-nexus-upload-command", os.O_APPEND|os.O_CREATE|os.O_WRONLY|os.O_SYNC, 0644)
		if nil == e {
			_, _ = f.WriteString(cmd + "\n")
			_ = f.Sync()
			_ = f.Close()
		} else {
			_, _ = fmt.Fprintf(os.Stderr, "Failed to open command file for write\n")
		}
	}
}
func main() {
	if len(os.Args) == 1 {
		_, _ = fmt.Fprintf(os.Stderr, "file path needed\n")
		return
	}
	filePath := os.Args[1]
	info, err := os.Stat(filePath)
	if nil != err {
		_, _ = fmt.Fprintf(os.Stderr, "Error read file %v\n", err)
		os.Exit(1)
	}
	if info.IsDir() {
		_, _ = fmt.Fprintf(os.Stderr, "path %v is a directory\n", filePath)
		os.Exit(1)
	}
	m := parseData(filePath)
	if nil != m {
		processItem(m, filePath)
	}
}
