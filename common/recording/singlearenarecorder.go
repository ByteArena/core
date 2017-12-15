package recording

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/bytearena/core/common/types/mapcontainer"
	"github.com/bytearena/core/common/utils"
)

type SingleArenaRecorder struct {
	filename           string
	tempBaseFilename   string
	recordFile         *os.File
	recordMetadataFile *os.File
}

func MakeSingleArenaRecorder(filename string) *SingleArenaRecorder {
	tempBaseFilename := os.TempDir() + "/" + filename

	f, err := os.OpenFile(tempBaseFilename, os.O_RDWR|os.O_CREATE, 0600)
	utils.Check(err, "Could not open file")

	return &SingleArenaRecorder{
		filename:         filename,
		tempBaseFilename: tempBaseFilename,
		recordFile:       f,
	}
}

func (r *SingleArenaRecorder) Stop() {

	err := os.Remove(r.tempBaseFilename + ".meta")
	if err != nil {
		log.Println("Could not remove record temporary meta file: " + err.Error())
	}

	err = os.Remove(r.tempBaseFilename)
	if err != nil {
		log.Println("Could not remove record temporary file: " + err.Error())
	}
}

func (r *SingleArenaRecorder) Close(UUID string) {
	files := make([]ArchiveFile, 0)

	files = append(files, ArchiveFile{
		Name: "RecordMetadata",
		Fd:   r.recordMetadataFile,
	})

	files = append(files, ArchiveFile{
		Name: "Record",
		Fd:   r.recordFile,
	})

	err, _ := MakeArchive(r.filename, files)
	utils.CheckWithFunc(err, func() string {
		return "could not create record archive: " + err.Error()
	})

	r.recordFile.Close()

	utils.Debug("SingleArenaRecorder", "write record archive")
}

func (r *SingleArenaRecorder) RecordMetadata(UUID string, mapcontainer *mapcontainer.MapContainer) error {
	filename := r.tempBaseFilename + ".meta"

	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	utils.Check(err, "Could not open RecordMetadata temporary file")

	metadata := RecordMetadata{
		MapContainer: mapcontainer,
		Date:         time.Now().Format(time.RFC3339),
	}

	data, err := json.Marshal(metadata)
	utils.Check(err, "could not marshall RecordMetadata")

	_, err = file.Write(data)

	err = file.Sync()
	utils.Check(err, "could not flush Record to disk")

	utils.Debug("SingleArenaRecorder", "wrote record metadata for game "+UUID)

	r.recordMetadataFile = file

	return nil
}

func (r *SingleArenaRecorder) Record(UUID string, msg string) error {
	_, err := r.recordFile.WriteString(msg + "\n")
	utils.Check(err, "could write record entry")

	err = r.recordFile.Sync()
	utils.Check(err, "could not flush Record to disk")

	return err
}

func (r *SingleArenaRecorder) GetFilePathForUUID(UUID string) string {
	return ""
}

func (r *SingleArenaRecorder) RecordExists(UUID string) bool {
	recordFile := r.GetFilePathForUUID(UUID)
	_, err := os.Stat(recordFile)

	return !os.IsNotExist(err)
}
