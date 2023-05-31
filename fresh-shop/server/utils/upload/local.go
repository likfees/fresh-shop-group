package upload

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path"
	"strings"
	"time"

	"fresh-shop/server/global"
	"fresh-shop/server/utils"
	"go.uber.org/zap"
)

type Local struct{}

//@author: [likfees](https://github.com/likfees)
//@author: [ccfish86](https://github.com/ccfish86)
//@author: [SliverHorn](https://github.com/SliverHorn)
//@object: *Local
//@function: UploadFile
//@description: 上传文件
//@param: file *multipart.FileHeader
//@return: string, string, error

func (*Local) UploadFile(file *multipart.FileHeader) (string, string, error) {
	// 读取文件后缀
	ext := path.Ext(file.Filename)
	// 读取文件名并加密
	name := strings.TrimSuffix(file.Filename, ext)
	name = utils.MD5V([]byte(name))
	// 拼接新文件名
	filename := name + "_" + time.Now().Format("20060102150405") + ext
	// 尝试创建此路径
	mkdirErr := os.MkdirAll(global.Config.Local.StorePath, os.ModePerm)
	if mkdirErr != nil {
		global.Log.Error("function os.MkdirAll() Filed", zap.Any("err", mkdirErr.Error()))
		return "", "", errors.New("function os.MkdirAll() Filed, err:" + mkdirErr.Error())
	}
	// 拼接路径和文件名
	p := global.Config.Local.StorePath + "/" + filename
	filepath := global.Config.Local.Path + "/" + filename

	f, openError := file.Open() // 读取文件
	if openError != nil {
		global.Log.Error("function file.Open() Filed", zap.Any("err", openError.Error()))
		return "", "", errors.New("function file.Open() Filed, err:" + openError.Error())
	}
	defer f.Close() // 创建文件 defer 关闭

	out, createErr := os.Create(p)
	if createErr != nil {
		global.Log.Error("function os.Create() Filed", zap.Any("err", createErr.Error()))

		return "", "", errors.New("function os.Create() Filed, err:" + createErr.Error())
	}
	defer out.Close() // 创建文件 defer 关闭

	_, copyErr := io.Copy(out, f) // 传输（拷贝）文件
	if copyErr != nil {
		global.Log.Error("function io.Copy() Filed", zap.Any("err", copyErr.Error()))
		return "", "", errors.New("function io.Copy() Filed, err:" + copyErr.Error())
	}
	return filepath, filename, nil
}

//@author: [likfees](https://github.com/likfees)
//@author: [ccfish86](https://github.com/ccfish86)
//@author: [SliverHorn](https://github.com/SliverHorn)
//@object: *Local
//@function: UploadFile
//@description: 上传文件
//@param: file *[]byte
//@return: string, string, error

func (*Local) UploadFileByBytes(file *[]byte, ext string) (string, string, error) {
	// 读取文件名并加密
	// 拼接新文件名
	now := time.Now()
	filename := "import_goods_" + now.Format("20060102150405") + fmt.Sprintf("%09d", now.Nanosecond()) + ext
	// 尝试创建此路径
	mkdirErr := os.MkdirAll(global.Config.Local.StorePath, os.ModePerm)
	if mkdirErr != nil {
		global.Log.Error("function os.MkdirAll() Filed", zap.Any("err", mkdirErr.Error()))
		return "", "", errors.New("function os.MkdirAll() Filed, err:" + mkdirErr.Error())
	}
	// 拼接路径和文件名
	p := global.Config.Local.StorePath + "/" + filename
	filepath := global.Config.Local.Path + "/" + filename

	reader := bytes.NewReader(*file) // 读取文件

	out, createErr := os.Create(p)
	if createErr != nil {
		global.Log.Error("function os.Create() Filed", zap.Any("err", createErr.Error()))

		return "", "", errors.New("function os.Create() Filed, err:" + createErr.Error())
	}
	defer out.Close() // 创建文件 defer 关闭

	_, copyErr := io.Copy(out, reader) // 传输（拷贝）文件
	if copyErr != nil {
		global.Log.Error("function io.Copy() Filed", zap.Any("err", copyErr.Error()))
		return "", "", errors.New("function io.Copy() Filed, err:" + copyErr.Error())
	}
	return filepath, filename, nil
}

//@author: [likfees](https://github.com/likfees)
//@author: [ccfish86](https://github.com/ccfish86)
//@author: [SliverHorn](https://github.com/SliverHorn)
//@object: *Local
//@function: DeleteFile
//@description: 删除文件
//@param: key string
//@return: error

func (*Local) DeleteFile(key string) error {
	p := global.Config.Local.StorePath + "/" + key
	if strings.Contains(p, global.Config.Local.StorePath) {
		if err := os.Remove(p); err != nil {
			return errors.New("本地文件删除失败, err:" + err.Error())
		}
	}
	return nil
}
