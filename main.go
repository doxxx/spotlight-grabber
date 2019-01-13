package main

import (
	"image"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	_ "image/jpeg"
	_ "image/png"
)

func main() {
	targetDir := "."
	if len(os.Args) > 1 {
		targetDir = os.Args[1]
	}

	localAppData := os.Getenv("LOCALAPPDATA")
	assetsDir := filepath.Join(localAppData, "Packages", "Microsoft.Windows.ContentDeliveryManager_cw5n1h2txyewy", "LocalState", "Assets")

	d, err := os.Open(assetsDir)
	if err != nil {
		log.Fatalln("ERROR: could not open Spotlight assets directory:", err)
	}
	filenames, err := d.Readdirnames(0)
	if err != nil {
		log.Fatalln("ERROR: could not read Spotlight assets directory contents:", err)
	}
	for _, filename := range filenames {
		filename = filepath.Join(assetsDir, filename)
		ok, w, h, ext, err := detectImage(filename)
		if err != nil {
			log.Fatalln("ERROR: could not detect file content type:", err)
		}
		if ok && h > 400 && w > h {
			destFilename := filepath.Join(targetDir, filepath.Base(filename)) + "." + ext
			if _, err := os.Stat(destFilename); os.IsNotExist(err) {
				log.Printf("Copying %s", filepath.Base(filename))
				err = copyFile(destFilename, filename)
				if err != nil {
					log.Fatalln("ERROR: could not copy file:", err)
				}
			}
		}
	}
}

var formatExtensions = map[string]string{
	"jpeg": "jpg",
}

func detectImage(filename string) (ok bool, width int, height int, ext string, err error) {
	var f *os.File
	f, err = os.Open(filename)
	if err != nil {
		return
	}
	defer f.Close()

	data := make([]byte, 512)
	_, err = f.Read(data)
	if err != nil {
		return
	}

	contentType := http.DetectContentType(data)
	if strings.HasPrefix(contentType, "image/") {
		f.Seek(0, os.SEEK_SET)

		var c image.Config
		c, ext, err = image.DecodeConfig(f)
		if err != nil {
			return
		}

		if newExt, found := formatExtensions[ext]; found {
			ext = newExt
		}

		ok = true
		width = c.Width
		height = c.Height
	}

	return
}

func copyFile(dest, src string) error {
	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return err
	}

	return nil
}
