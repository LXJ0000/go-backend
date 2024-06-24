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
}

func NewFileUsecase(repo domain.FileRepository, contextTimeout time.Duration, localStaticUrl string, urlStaticUrl string) domain.FileUsecase {
	return &fileUsecase{repo: repo,
		contextTimeout: contextTimeout,
		localStaticUrl: localStaticUrl,
		urlStaticUrl:   urlStaticUrl,
	}
}

func (f *fileUsecase) Upload(c context.Context, sourceFile *multipart.FileHeader) (domain.File, error) {
	ctx, cancel := context.WithTimeout(c, f.contextTimeout)
	defer cancel()

	fileHash, err := f.hash(sourceFile)
	if err != nil {
		return domain.File{}, err
	}

	exists, file, err := f.checkFileMd5(ctx, fileHash)
	if err != nil {
		return domain.File{}, err
	}
	if exists {
		return file, nil
	}

	sourceFile.Filename = f.uniqueFileName(sourceFile)
	dst := filepath.Join(f.localStaticUrl, "file", sourceFile.Filename)
	if err := fileutil.SaveUploadedFile(sourceFile, dst); err != nil {
		return domain.File{}, err
	}

	now := time.Now().UnixMicro()
	file = domain.File{
		FileID: snowflake.GenID(),
		Name:   sourceFile.Filename,
		Path:   dst,
		Type:   sourceFile.Header.Get("Content-Type"),
		Source: domain.FileSourceLocal,
		Hash:   fileHash,
	}
	file.CreatedAt = now
	file.UpdatedAt = now
	if err := f.repo.Upload(ctx, file); err != nil {
		go func() { _ = fileutil.RemoveFile(dst) }()
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

			sourceFile.Filename = f.uniqueFileName(sourceFile)
			dst := filepath.Join(f.localStaticUrl, "file", sourceFile.Filename)
			if err := fileutil.SaveUploadedFile(sourceFile, dst); err != nil {
				resp.Data[sourceFileName] = err.Error()
				return
			}

			now := time.Now().UnixMicro()
			file = domain.File{
				FileID: snowflake.GenID(),
				Name:   sourceFile.Filename,
				Path:   dst,
				Type:   sourceFile.Header.Get("Content-Type"),
				Source: domain.FileSourceLocal,
				Hash:   fileHash,
			}
			file.CreatedAt = now
			file.UpdatedAt = now
			if err := f.repo.Upload(ctx, file); err != nil {
				resp.Data[sourceFileName] = err.Error()
				go func() {
					_ = fileutil.RemoveFile(dst)
				}()
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
