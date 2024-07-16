package databases

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/ostafen/clover/v2"
)

func NewHistoryDatabase() *clover.DB {
	deviceHistory := "devicehistory"
	personDeviceHistory := "persondevicehistory"
	personHistory := "personhistory"

	exeDirectory, errAbs := filepath.Abs(filepath.Dir(os.Args[0]))
	if errAbs != nil {
		log.Fatal(errAbs)
		return nil
	}
	if strings.Contains(exeDirectory, "go-build") {
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatal("Error getting current working directory:", err)
			return nil
		}
		exeDirectory = cwd
	}
	cloverDbDirectory := path.Join(exeDirectory, "clover_data/history")
	fmt.Println("Current working directory:", cloverDbDirectory)
	os.MkdirAll(cloverDbDirectory, os.ModePerm)

	db, dbOpenErr := clover.Open(cloverDbDirectory)
	if dbOpenErr != nil {
		log.Fatal(dbOpenErr)
		return nil
	}
	for _, collectionName := range []string{
		deviceHistory,
		personDeviceHistory,
		personHistory,
	} {
		hasDeviceHistory, hasCollectionErr := db.HasCollection(collectionName)
		if hasCollectionErr != nil {
			log.Fatal(hasCollectionErr)
			return nil
		}
		if !hasDeviceHistory {
			db.CreateCollection(collectionName)
		}
	}
	return db
}
