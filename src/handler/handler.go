package handler

import (
	"encoding/json"
	"filestore-service/src/mate"
	"filestore-service/src/util"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

// 文件上传
func UploadHandler(response http.ResponseWriter, request *http.Request) {
	if request.Method == "GET" {
		data, err := ioutil.ReadFile("../static/view/index.html")
		if err != nil {
			io.WriteString(response, "internel server error")
			return
		}
		io.WriteString(response, string(data))

	} else if request.Method == "POST" {
		file, head, err := request.FormFile("file")
		if err != nil {
			fmt.Printf("Failed to get data,err:%s\n", err.Error())
			return
		}
		defer file.Close()

		fileMeta := mate.FileMeta{
			FileName: head.Filename,
			Location: "../../tmp/" + head.Filename,
			UploadAt: time.Now().Format("2020-7-7 3:46"),
		}

		newFile, err := os.Create(fileMeta.Location)
		if err != nil {
			fmt.Printf("Failed to create file,err:%s\n", err.Error())
			return
		}
		defer newFile.Close()
		fileMeta.FileSize, err = io.Copy(newFile, file)
		if err != nil {
			fmt.Printf("Failed to save data into file,err:%s\n", err.Error())
			return
		}

		newFile.Seek(0, 0)
		fileMeta.FileSha1 = util.FileSha1(newFile)
		fmt.Print(fileMeta.FileSha1)
		// mate.UpdateFileMeta(fileMeta)
		_ = mate.UpdateFileMetaDB(fileMeta)
		http.Redirect(response, request, "/file/upload/success", http.StatusFound)
	}
}

// 上传已完成
func UploadSuccessHandler(response http.ResponseWriter, request *http.Request) {
	io.WriteString(response, "upload finished")
}

// 查询文件元信息
func GetFileMetaHandler(response http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	filehash := request.Form["filehash"][0]
	// fMeta := mate.GetFileMeta(filehash)
	fMeta, err := mate.GetFileMetaDB(filehash)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(fMeta)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	response.Write(data)
}

// 下载文件
func DownloadHandler(response http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	fsha1 := request.Form.Get("filehash")
	fm := mate.GetFileMeta(fsha1)
	f, err := os.Open(fm.Location)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	response.Header().Set("Content-Type", "application/octect-stream")
	response.Header().Set("Content-Descrption", "attachment;filename=\""+fm.FileName+"\"")
	response.Write(data)
}

// 更新元数据信息
func FileMetaUpdateHandler(response http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	opType := request.Form.Get("op")
	fileSha1 := request.Form.Get("filehash")
	newFileName := request.Form.Get("filename")

	if opType != "0" {
		response.WriteHeader(http.StatusForbidden)
		return
	}

	if request.Method != "POST" {
		response.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	curFileMeta := mate.GetFileMeta(fileSha1)
	curFileMeta.FileName = newFileName
	mate.UpdateFileMeta(curFileMeta)

	data, err := json.Marshal(curFileMeta)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		return
	}
	response.WriteHeader(http.StatusOK)
	response.Write(data)
}

// 删除文件
func FileDeleteHandler(response http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	fileSha1 := request.Form.Get("filehash")

	fMeta := mate.GetFileMeta(fileSha1)
	os.Remove(fMeta.Location)

	mate.RemoveFileMeta(fileSha1)

	response.WriteHeader(http.StatusOK)
}
