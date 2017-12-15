package recording

import (
	"github.com/bytearena/core/common/types/mapcontainer"
)

type EmptyRecorder struct{}

func MakeEmptyRecorder() EmptyRecorder {
	return EmptyRecorder{}
}

func (r EmptyRecorder) Record(UUID string, msg string) error {
	return nil
}

func (r EmptyRecorder) RecordMetadata(UUID string, mapcontainer *mapcontainer.MapContainer) error {
	return nil
}

func (r EmptyRecorder) Close(UUID string) {}
func (r EmptyRecorder) Stop()             {}

func (r EmptyRecorder) GetFilePathForUUID(UUID string) string {
	return ""
}

func (r EmptyRecorder) RecordExists(UUID string) bool {
	return false
}
