package fileutil

//
//import (
//	"fmt"
//	"github.com/LXJ0000/go-backend/utils/md5util"
//	"github.com/LXJ0000/go-backend/utils/uuidutil"
//	"io"
//	"log/slog"
//	"mime/multipart"
//	"path/filepath"
//	"slices"
//	"strings"
//)
//
//const (
//	staticPath = "assets/"
//	fileSize   = 20
//)
//
//// WhiteImageList 文件白名单
//var WhiteImageList = []string{
//	".jpg", ".png", ".jpeg", ".ico", ".tiff", ".gif", ".svg", ".webp",
//}
//
//type FileUploadResponse struct {
//	FilePath  string `json:"file_name"`
//	IsSuccess bool   `json:"is_success"` // 是否上传成功
//	Msg       string `json:"msg"`
//}
//
//// UploadFile 文件上传
//func UploadFile(file *multipart.FileHeader, fileType string) (fileUploadResponse FileUploadResponse) {
//	natureName := file.Filename
//	basePath := staticPath + fileType
//	fileName := uuidutil.UUID(10)
//	filePath := filepath.Join(basePath, fileName)
//	//默认上传失败
//	fileUploadResponse = FileUploadResponse{FilePath: filePath, IsSuccess: false, Msg: "UploadFile Fail"}
//	//文件名后缀
//	suffix := strings.ToLower(filepath.Ext(fileName))
//	//文件白名单判断
//	if !slices.Contains(WhiteImageList, suffix) {
//		//上传失败
//		slog.Warn("UploadFile Fail", "msg", "文件格式有误")
//		fileUploadResponse.Msg = "文件格式有误"
//		return
//	}
//	//文件大小判断
//	if size := float64(file.Size) / float64(1024*1024); size > fileSize {
//		//上传失败
//		slog.Warn("UploadFile Fail", "msg", fmt.Sprintf("文件当前大小为：%.2fMb 超出限定大小：%dMb", size, fileSize))
//		fileUploadResponse.Msg = fmt.Sprintf("文件当前大小为：%.2fMb 超出限定大小：%dMb", size, fileSize)
//		return
//	}
//
//	//读取文件内容
//	fileObj, err := file.Open()
//	if err != nil {
//		slog.Error("UploadFile Fail", "Error", err.Error())
//		return
//	}
//	byteData, err := io.ReadAll(fileObj)
//	if err != nil {
//		slog.Error("UploadFile Fail", "Error", err.Error())
//		return
//	}
//	//获取文件MD5
//	imageHash := md5util.Md5(byteData)
//	//查询数据库是否存在对应Hash
//	var bannerModel models.BannerModel
//	if err := global.DB.Where("hash = ?", imageHash).First(&bannerModel).Error; err == nil {
//		//	找到了 不需要存入数据库
//		fileUploadResponse.FilePath = bannerModel.Path
//		fileUploadResponse.Msg = "图片已存在"
//		return
//	}
//
//	//文件存储类型
//	fileType := ctype.Local
//	fileUploadResponse.Msg = "图片上传成功 ~本地"
//	if global.Config.QiNiu.Enable {
//		//开启七牛云存储
//		filePath, err = qiniu.UploadImage(byteData, fileName, global.Config.QiNiu.Prefix)
//		if err != nil {
//			global.Log.Error(err.Error())
//			fileUploadResponse.Msg = err.Error()
//			return
//		}
//		//入库
//		fileUploadResponse.FilePath = filePath
//		fileUploadResponse.Msg = "图片上传成功 ~七牛云"
//		fileType = ctype.QiNiu
//
//		//global.DB.Create(&models.BannerModel{
//		//	Path:      filePath,
//		//	Hash:      imageHash,
//		//	Name:      fileName,
//		//	ImageType: fileType,
//		//})
//		//return
//	} else {
//		filePath = "/" + filePath
//	}
//	//if err := c.SaveUploadedFile(file, filePath); err != nil {
//	//	global.Log.Info(err.Error())
//	//	//上传失败
//	//	fileUploadResponse.IsSuccess = false
//	//	fileUploadResponse.Msg = err.Error()
//	//}
//
//	//入库
//	fileUploadResponse.IsSuccess = true
//	global.DB.Create(&models.BannerModel{
//		Path:      filePath,
//		Hash:      imageHash,
//		Name:      natureName,
//		ImageType: fileType,
//	})
//	return
//}
