package usecase

import (
	"context"
	"errors"
	"fmt"
	"github.com/LXJ0000/go-backend/utils/fileutil"
	"github.com/LXJ0000/go-backend/utils/md5util"
	snowflake "github.com/LXJ0000/go-backend/utils/snowflakeutil"
	"gorm.io/gorm"
	"io"
	"log/slog"
	"mime/multipart"
	"path/filepath"
	"time"

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

	exists, file, err := f.checkFileMd5(ctx, sourceFile)
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

	file = domain.File{
		FileID: snowflake.GenID(),
		Name:   sourceFile.Filename,
		Path:   dst,
		Type:   domain.FileTypeUnknown,
		Source: domain.FileSourceLocal,
	}
	if err := f.repo.Upload(ctx, file); err != nil {
		go func() { _ = fileutil.RemoveFile(dst) }()
		return domain.File{}, nil
	}

	return file, nil
}

func (f *fileUsecase) FileList(c context.Context, fileType, fileSource string, page, size int) ([]domain.File, int, error) {
	ctx, cancel := context.WithTimeout(c, f.contextTimeout)
	defer cancel()
	return f.repo.FileList(ctx, fileType, fileSource, page, size)
}

func (f *fileUsecase) uniqueFileName(file *multipart.FileHeader) string {
	return fmt.Sprintf("%d%s", time.Now().UnixNano(), filepath.Ext(file.Filename))
}

func (f *fileUsecase) checkFileMd5(c context.Context, file *multipart.FileHeader) (bool, domain.File, error) {
	//读取文件内容
	fileObj, err := file.Open()
	if err != nil {
		slog.Error("Open File Fail", "Error", err.Error())
		return false, domain.File{}, err
	}
	byteData, err := io.ReadAll(fileObj)
	if err != nil {
		slog.Error("Read File Fail", "Error", err.Error())
		return false, domain.File{}, err
	}
	//获取文件MD5
	fileHash := md5util.Md5(byteData)
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
