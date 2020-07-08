package mate

import "filestore-service/src/db"

// 文件元信息结构
type FileMeta struct {
	FileSha1 string
	FileName string
	FileSize int64
	Location string
	UploadAt string
}

var fileMetas map[string]FileMeta

func init() {
	fileMetas = make(map[string]FileMeta, 16)
}

// 新增/更新文件元信息
func UpdateFileMeta(fmeta FileMeta) {
	fileMetas[fmeta.FileSha1] = fmeta
}

// 新增/更新文件元信息(数据库)
func UpdateFileMetaDB(fmeta FileMeta) bool {
	return db.OnFileUploadFinished(fmeta.FileSha1, fmeta.FileName, fmeta.FileSize, fmeta.Location)
}

// 通过sha1获取文件的元信息
func GetFileMeta(fileSha1 string) FileMeta {
	return fileMetas[fileSha1]
}

// 通过sha1获取文件的元信息 (数据库)
func GetFileMetaDB(fileSha1 string) (FileMeta, error) {
	tfile, err := db.GetFileMeta(fileSha1)
	if err != nil {
		return FileMeta{}, nil
	}
	fmeta := FileMeta{
		FileSha1: tfile.FileHash,
		FileName: tfile.FileName.String,
		FileSize: tfile.FileSize.Int64,
		Location: tfile.FileAddr.String}
	return fmeta, nil
}

// 删除原信息
func RemoveFileMeta(fileSha1 string) {
	delete(fileMetas, fileSha1)
}
