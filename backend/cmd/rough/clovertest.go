package main

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/ostafen/clover/v2"
	"github.com/ostafen/clover/v2/document"
	"github.com/ostafen/clover/v2/query"
)

func main() {
	deviceHistory := "devicehistory"
	personDeviceHistory := "persondevicehistory"
	personHistory := "personhistory"

	collectionNames := []string{
		deviceHistory,
		personDeviceHistory,
		personHistory,
	}

	exeDirectory, errAbs := filepath.Abs(filepath.Dir(os.Args[0]))
	if errAbs != nil {
		fmt.Printf("errAbs: %v\n", errAbs)
		return
	}
	fmt.Printf("dir: %v\n", exeDirectory)
	if strings.Contains(exeDirectory, "go-build") {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Println("Error getting current working directory:", err)
			return
		}
		exeDirectory = cwd
	}
	cloverDbDirectory := path.Join(exeDirectory, "clover_data/history")
	fmt.Println("Current working directory:", cloverDbDirectory)
	os.MkdirAll(cloverDbDirectory, os.ModePerm)

	db, dbOpenErr := clover.Open(cloverDbDirectory)
	if dbOpenErr != nil {
		fmt.Printf("dbOpenErr: %v\n", dbOpenErr)
		return
	}
	defer db.Close()
	for _, collectionName := range collectionNames {
		hasDeviceHistory, hasCollectionErr := db.HasCollection(collectionName)
		if hasCollectionErr != nil {
			fmt.Printf("hasCollectionErr: %v\n", hasCollectionErr)
			return
		}
		if !hasDeviceHistory {
			db.CreateCollection(collectionName)
		}
	}

	doc := document.NewDocument()
	doc.Set("Honk honk", "blarge blarge")
	db.InsertOne(personDeviceHistory, doc)

	docs, findAllErr := db.FindAll(query.NewQuery(personDeviceHistory))
	if findAllErr != nil {
		fmt.Printf("findAllErr: %v\n", findAllErr)
		return
	}
	for _, doc := range docs {
		fmt.Printf("doc: %v\n", doc)
	}
}
