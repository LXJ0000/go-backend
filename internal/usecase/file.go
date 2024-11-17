package usecase

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"path/filepath"
	"sync"
	"time"

	"github.com/LXJ0000/go-backend/pkg/file"
	"github.com/LXJ0000/go-backend/utils/fileutil"
	"github.com/LXJ0000/go-backend/utils/md5util"
	snowflake "github.com/LXJ0000/go-backend/utils/snowflakeutil"
	"gorm.io/gorm"

	"github.com/LXJ0000/go-backend/internal/domain"
)

type fileUsecase struct {
	repo           domain.FileRepository
	contextTimeout time.Duration
	localStaticUrl string
	urlStaticUrl   string

	minioClient file.FileStorage
}

func NewFileUsecase(repo domain.FileRepository,
	contextTimeout time.Duration, localStaticUrl string, urlStaticUrl string,
	minioClient file.FileStorage,
) domain.FileUsecase {
	return &fileUsecase{repo: repo,
		contextTimeout: contextTimeout,
		localStaticUrl: localStaticUrl,
		urlStaticUrl:   urlStaticUrl,
		minioClient:    minioClient,
	}
}

func (f *fileUsecase) Upload(c context.Context, sourceFile *multipart.FileHeader) (domain.File, error) {
	ctx, cancel := context.WithTimeout(c, f.contextTimeout)
	defer cancel()

	// 计算文件 MD5
	fileHash, err := f.hash(sourceFile)
	if err != nil {
		return domain.File{}, err
	}

	// 检查文件是否已存在
	exists, file, err := f.checkFileMd5(ctx, fileHash)
	if err != nil {
		return domain.File{}, err
	}
	if exists {
		return file, nil
	}

	file = domain.File{
		FileID: snowflake.GenID(),
		Name:   sourceFile.Filename,
		Type:   sourceFile.Header.Get("Content-Type"),
		Hash:   fileHash,
	}

	// 文件不存在 存储文件 默认存储到 minio 后续改造可以选择存储本地与 minio
	var path string
	if true { // 判断文件上传到本地还是 minio
		file.Source = domain.FileSourceMinio
		path, err = f.upload2Minio(ctx, sourceFile)
		if err != nil {
			return domain.File{}, err
		}
	} else {
		file.Source = domain.FileSourceLocal
		path, err = f.upload2Local(sourceFile)
		if err != nil {
			return domain.File{}, err
		}
	}

	// 存储到数据库
	file.Path = path
	if err := f.repo.Upload(ctx, &file); err != nil { // 存储失败的一些操作
		go func() {
			if true { // 判断文件上传到本地还是 minio

			} else {
				_ = fileutil.RemoveFile(path)
			}
		}() // 上传失败删除文件
		return domain.File{}, nil
	}

	return file, nil
}

func (f *fileUsecase) Uploads(c context.Context, sourceFiles []*multipart.FileHeader) (domain.FileUploadsResponse, error) {
	resp := domain.FileUploadsResponse{
		Data: make(map[string]interface{}, len(sourceFiles)),
	}
	g := sync.WaitGroup{}
	g.Add(len(sourceFiles))
	for _, sourceFile := range sourceFiles {
		sourceFile := sourceFile
		go func() {
			defer g.Done()
			sourceFileName := sourceFile.Filename
			ctx, cancel := context.WithTimeout(c, f.contextTimeout)
			defer cancel()
			fileHash, err := f.hash(sourceFile)
			if err != nil {
				resp.Data[sourceFileName] = err.Error()
				return
			}
			exists, file, err := f.checkFileMd5(ctx, fileHash)
			if err != nil {
				resp.Data[sourceFileName] = err.Error()
				return
			}
			if exists {
				resp.Data[sourceFileName] = file
				return
			}

			file = domain.File{
				FileID: snowflake.GenID(),
				Name:   sourceFile.Filename,
				Type:   sourceFile.Header.Get("Content-Type"),
				Hash:   fileHash,
			}

			// 文件不存在 存储文件 默认存储到 minio 后续改造可以选择存储本地与 minio
			var path string
			if true { // 判断文件上传到本地还是 minio
				file.Source = domain.FileSourceMinio
				path, err = f.upload2Minio(ctx, sourceFile)
				if err != nil {
					resp.Data[sourceFileName] = err.Error()
					return
				}
			} else {
				file.Source = domain.FileSourceLocal
				path, err = f.upload2Local(sourceFile)
				if err != nil {
					resp.Data[sourceFileName] = err.Error()
					return
				}
			}

			// 存储到数据库
			file.Path = path
			if err := f.repo.Upload(ctx, &file); err != nil { // 存储失败的一些操作
				go func() {
					if true { // 判断文件上传到本地还是 minio

					} else {
						_ = fileutil.RemoveFile(path)
					}
				}() // 上传失败删除文件
				resp.Data[sourceFileName] = err.Error()
				return
			}
			resp.Data[sourceFileName] = file
		}()
	}
	g.Wait()

	return resp, nil
}

func (f *fileUsecase) FileList(c context.Context, fileType, fileSource string, page, size int) ([]domain.File, int, error) {
	ctx, cancel := context.WithTimeout(c, f.contextTimeout)
	defer cancel()
	return f.repo.FileList(ctx, fileType, fileSource, page, size)
}

func (f *fileUsecase) uniqueFileName(file *multipart.FileHeader) string {
	return fmt.Sprintf("%d%s", time.Now().UnixNano(), filepath.Ext(file.Filename))
}

func (f *fileUsecase) checkFileMd5(c context.Context, fileHash string) (bool, domain.File, error) {
	gotFile, err := f.repo.FindByHash(c, fileHash)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, domain.File{}, nil
		}
		slog.Error("Find FileHash Fail", "Error", err.Error())
		return false, domain.File{}, err
	}
	return true, gotFile, nil
}

func (f *fileUsecase) hash(file *multipart.FileHeader) (string, error) {
	// 读取文件内容
	fileObj, err := file.Open()
	if err != nil {
		slog.Error("Open File Fail", "Error", err.Error())
		return "", err
	}
	byteData, err := io.ReadAll(fileObj)
	if err != nil {
		slog.Error("Read File Fail", "Error", err.Error())
		return "", err
	}
	// 获取文件MD5
	fileHash := md5util.Md5(byteData)
	return fileHash, nil
}

func (f *fileUsecase) upload2Minio(c context.Context, sourceFile *multipart.FileHeader) (string, error) {
	ctx, cancel := context.WithTimeout(c, f.contextTimeout)
	defer cancel()
	byteData, err := fileutil.FileHeaderToBytes(sourceFile)
	if err != nil {
		return "", err
	}
	if err := f.minioClient.UploadFile(ctx, domain.FileBucket, sourceFile.Filename, byteData); err != nil {
		return "", err
	}
	// 获取 minio 文件路径
	return f.minioClient.GetFilePath(ctx, domain.FileBucket, sourceFile.Filename)
}

func (f *fileUsecase) upload2Local(sourceFile *multipart.FileHeader) (string, error) {
	sourceFile.Filename = f.uniqueFileName(sourceFile)
	path := filepath.Join(f.localStaticUrl, "file", sourceFile.Filename)
	if err := fileutil.SaveUploadedFile(sourceFile, path); err != nil {
		return "", err
	}
	return path, nil
}
