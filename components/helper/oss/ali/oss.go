package ali

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"
	"time"
)

var _oss = new(Oss)

type Oss struct {
	conf Conf
}

func NewAilOss(conf Conf) *Oss {
	this := new(Oss)
	this.conf = conf
	if this.conf.Dir != "" {
		this.conf.Dir = strings.Trim(this.conf.Dir, "/") + "/"
	}
	return _oss
}

func (v *Oss) Client() (*oss.Client, error) {
	return oss.New(v.conf.Endpoint, v.conf.AccessKeyId, v.conf.AccessKeySecret)
}

func (v *Oss) GetBucket() (*oss.Bucket, error) {
	client, err := v.Client()
	if err != nil {
		return nil, err
	}

	return client.Bucket(v.conf.Bucket)
}

// 流下载
//______________________________________________________________________
func (v *Oss) GetFileStream(objectName string) ([]byte, error) {
	bucket, err := v.GetBucket()
	if err != nil {
		return []byte{}, err
	}
	body, err := bucket.GetObject(objectName)
	defer body.Close()

	return ioutil.ReadAll(body)
}

// 下载文件到本地
//______________________________________________________________________
func (v *Oss) GetFile(ObjectName string, localFile string) error {
	bucket, err := v.GetBucket()
	if err != nil {
		return err
	}

	return bucket.GetObjectToFile(ObjectName, localFile)
}

// 上传字符串
//______________________________________________________________________
func (v *Oss) UploadFileString(objectName string, file string) error {
	bucket, err := v.GetBucket()
	if err != nil {
		return err
	}
	return bucket.PutObject(objectName, strings.NewReader(file))
}

// 上传Byte数组
//______________________________________________________________________
func (v *Oss) UploadFileSilce(objectName string, data []byte) error {
	bucket, err := v.GetBucket()
	if err != nil {
		return err
	}

	return bucket.PutObject(objectName, bytes.NewReader(data))
}

// 上传文件流
//______________________________________________________________________
func (v *Oss) UploadFileStream(objectName string, localFilePath string) error {
	bucket, err := v.GetBucket()
	if err != nil {
		return err
	}
	fd, err := os.Open(localFilePath)
	if err != nil {
		return err
	}
	defer fd.Close()

	return bucket.PutObject(objectName, fd)
}

// 上传本地文件
//______________________________________________________________________
func (v *Oss) UploadLocalFile(objectName string, localFilePath string) error {
	bucket, err := v.GetBucket()
	if err != nil {
		return err
	}

	return bucket.PutObjectFromFile(objectName, localFilePath)
}

//
//______________________________________________________________________
func (v *Oss) UploadRemoteFile(urlstring string, objectName string) error {
	resp, err := http.Get(urlstring)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return v.UploadFileSilce(objectName, body)
}

// 删除
//______________________________________________________________________
func (v *Oss) Del(objectName string) error {
	bucket, err := v.GetBucket()
	if err != nil {
		return err
	}
	return bucket.DeleteObject(objectName)
}

// 图片处理持久化：缩放。将目标图片存放到相同的bucket中。
// https://help.aliyun.com/document_detail/93499.html?spm=a2c4g.11186623.2.23.3efb31684CJLxz#concept-knm-w3d-lfb
//______________________________________________________________________
func (v *Oss) FileResize(sourceObjectName string, targetObjectName string, style string) error {
	bucket, err := v.GetBucket()
	if err != nil {
		return err
	}
	process := fmt.Sprintf("%s|sys/saveas,o_%v", style, base64.URLEncoding.EncodeToString([]byte(targetObjectName)))
	result, err := bucket.ProcessObject(sourceObjectName, process)
	if err != nil {
		return err
	}
	if strings.ToLower(result.Status) != "ok" {
		return errors.New("异常错误")
	}

	return nil
}

// sx原图宽  sy原图高  d二维数组 0 x>=y   1:x<y   [[1500,1125],[1125,1500]]
//______________________________________________________________________
func (v *Oss) FileResizeStyle(sx, sy int, d [2][2]int) (style string, w int, h int) {

	if sx >= sy {
		return fmt.Sprintf("image/resize,m_pad,color_000000,w_%d,h_%d", d[0][0], d[0][1]), d[0][0], d[0][1]
	} else {
		return fmt.Sprintf("image/resize,m_pad,color_000000,w_%d,h_%d", d[1][0], d[1][1]), d[1][0], d[1][1]

	}
}

// 图片处理：裁剪。将目标图片存放到相同的bucket中。
// https://help.aliyun.com/document_detail/93499.html?spm=a2c4g.11186623.2.23.3efb31684CJLxz#concept-knm-w3d-lfb
//______________________________________________________________________
func (v *Oss) FileCrop(sourceObjectName string, targetObjectName string, style string) error {
	bucket, err := v.GetBucket()
	if err != nil {
		return err
	}
	process := fmt.Sprintf("%s|sys/saveas,o_%v", style, base64.URLEncoding.EncodeToString([]byte(targetObjectName)))

	result, err := bucket.ProcessObject(sourceObjectName, process)
	if err != nil {
		return err
	}
	if strings.ToLower(result.Status) != "ok" {
		return errors.New("异常错误")
	}

	return nil
}

// sx原图宽  sy原图高  dx目标图宽 dy目标图高
//______________________________________________________________________
func (v *Oss) FileCropStyle(sx, sy, dx, dy int) string {

	if sx >= dx && sy >= dy {
		return fmt.Sprintf("image/crop,w_%d,h_%d,x_%d,y_%d", dx, dy, (sx-dx)/2, (sy-dy)/2)
	} else {
		if sx < dx && sy >= dy {
			return fmt.Sprintf("image/crop,w_%d,h_%d,x_%d,y_%d", sx, dy, 0, (sy-dy)/2)
		} else if sx >= dx && sy < dy {
			return fmt.Sprintf("image/crop,w_%d,h_%d,x_%d,y_%d", dx, sy, (sx-dx)/2, 0)
		} else {
			return fmt.Sprintf("image/crop,w_%d,h_%d,x_%d,y_%d", sx, sy, 0, 0)
		}
	}
}

// 获取文件元信息
//______________________________________________________________________
func (v *Oss) GetFileMeta(objectName string) (http.Header, error) {
	bucket, err := v.GetBucket()
	if err != nil {
		return nil, err
	}

	return bucket.GetObjectDetailedMeta(objectName)
}

// 获取图片宽高
//______________________________________________________________________
func (v *Oss) GetFileInfo(urlString string) (map[string]int, error) {
	u, err := url.Parse(urlString)
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(u.Scheme + "://" + u.Host + u.Path + "?x-oss-process=image/info")
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var res map[string]map[string]string
	var res1 map[string]int
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, err
	}

	res1 = make(map[string]int, 0)
	res1["w"], err = strconv.Atoi(res["ImageWidth"]["value"])
	if err != nil {
		return nil, err
	}
	res1["h"], err = strconv.Atoi(res["ImageHeight"]["value"])
	if err != nil {
		return nil, err
	}

	return res1, nil

}

// 文件后缀
//______________________________________________________________________
func (v *Oss) GetSuffix(filename string) string {
	return strings.Split(path.Ext(filename), "?")[0]
}

// 获取文件前缀
//______________________________________________________________________
func (v *Oss) GetPrefix() string {
	now := time.Now()
	return strings.Trim(fmt.Sprintf("%s/%s/%s/", strings.Trim(v.conf.Dir, "/"), now.Format("200601"), now.Format("02")), "/") + "/"
}

// 获取文件前缀
//______________________________________________________________________
func (v *Oss) GetPrefixByDir(dir string) string {
	now := time.Now()
	return strings.Trim(fmt.Sprintf("%s/%s/%s/%s/", strings.Trim(v.conf.Dir, "/"), strings.Trim(dir, "/"), now.Format("200601"), now.Format("02")), "/") + "/"
}

// 获取文件前缀
//______________________________________________________________________
func (v *Oss) GetPrefixById(uid string) string {
	now := time.Now()
	return strings.Trim(fmt.Sprintf("%s/user/%s/%s/", strings.Trim(v.conf.Dir, "/"), uid, now.Format("200601")), "/") + "/"
}

// 获取新的完整名字
//______________________________________________________________________
func (v *Oss) NewObjectName(filename string, uuid string) string {
	return v.GetPrefix() + uuid + v.GetSuffix(filename)
}

// 获取新的完整名字
//______________________________________________________________________
func (v *Oss) NewObjectNameByUid(filename string, uid string, uuid string) string {
	return v.GetPrefixById(uid) + uuid + v.GetSuffix(filename)
}

func (v *Oss) GetObjectNameFromUrl(urlstring string) (string, error) {
	u, err := url.Parse(urlstring)
	if err != nil {
		return "", err
	}
	return strings.Trim(u.Path, "/"), nil
}
