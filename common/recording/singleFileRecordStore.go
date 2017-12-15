package recording

type singleFileRecordStore struct {
	filepath string
}

func NewSingleFileRecordStore(filepath string) *singleFileRecordStore {
	return &singleFileRecordStore{
		filepath: filepath,
	}
}

func (r *singleFileRecordStore) GetFilePathForUUID(UUID string) string {
	return r.filepath
}

func (r *singleFileRecordStore) RecordExists(UUID string) bool {
	return true
}
