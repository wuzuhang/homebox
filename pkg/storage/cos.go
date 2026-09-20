package storage

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sync"

	"github.com/tencentyun/cos-go-sdk-v5"
)

// COSClient 封装腾讯云 COS 客户端
type COSClient struct {
	Client *cos.Client
	Region string
}

var (
	// GlobalCOS 全局暴露的 COS 实例，供全项目直接调用
	GlobalCOS *COSClient
	once      sync.Once
)

// InitCOS 全局初始化函数（在 main.go 启动时只需调用一次）
func InitCOS(bucketName, region, secretID, secretKey string) {
	once.Do(func() {
		GlobalCOS = NewCOSClient(bucketName, region, secretID, secretKey)
	})
}

// NewCOSClient 初始化 COS 客户端
func NewCOSClient(bucketName, region, secretID, secretKey string) *COSClient {
	// 拼接存储桶 URL
	bucketURLStr := fmt.Sprintf("https://%s.cos.%s.myqcloud.com", bucketName, region)
	u, _ := url.Parse(bucketURLStr)

	baseURL := &cos.BaseURL{BucketURL: u}

	// 使用密钥初始化客户端
	client := cos.NewClient(baseURL, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  secretID,
			SecretKey: secretKey,
		},
	})

	return &COSClient{
		Client: client,
		Region: region,
	}
}

// CreateBucket 创建存储桶
func (c *COSClient) CreateBucket(ctx context.Context) error {
	opt := &cos.BucketPutOptions{
		// 私有读写：private，公有读私有写：public-read
		XCosACL: "private",
	}

	_, err := c.Client.Bucket.Put(ctx, opt)
	if err != nil {
		return fmt.Errorf("创建存储桶失败: %w", err)
	}
	return nil
}

// UploadFile 上传文件到存储桶
func (c *COSClient) UploadFile(ctx context.Context, cosPath, localPath string) (string, error) {
	// 本地文件上传
	if c == nil || c.Client == nil {
		return "", fmt.Errorf("COS 客户端未初始化，请先调用 InitCOS()")
	}
	_, _, err := c.Client.Object.Upload(ctx, cosPath, localPath, nil)
	if err != nil {
		// ⚠️ 注意：返回 2 个值 (空字符串, error)
		return "", fmt.Errorf("文件上传失败: %w", err)
	}
	// 2. 获取上传成功后的文件访问 URL
	fileURL := c.Client.Object.GetObjectURL(cosPath).String()
	return fileURL, nil
}

// UploadStream 上传流式数据到存储桶
func (c *COSClient) UploadStream(ctx context.Context, cosPath string, data io.Reader) (string, error) {
	if c == nil || c.Client == nil {
		return "", fmt.Errorf("COS 客户端未初始化，请先调用 InitCOS()")
	}
	_, err := c.Client.Object.Put(ctx, cosPath, data, nil)
	if err != nil {
		return "", fmt.Errorf("流式数据上传失败: %w", err)
	}
	// 获取上传成功后的文件访问 URL
	fileURL := c.Client.Object.GetObjectURL(cosPath).String()
	return fileURL, nil
}

// DownloadFile 从存储桶下载文件
func (c *COSClient) DownloadFile(ctx context.Context, cosPath, localPath string) error {
	resp, err := c.Client.Object.Get(ctx, cosPath, nil)
	if err != nil {
		return fmt.Errorf("读取 COS 对象失败: %w", err)
	}
	defer resp.Body.Close()

	// 创建本地文件并写入
	out, err := os.Create(localPath)
	if err != nil {
		return fmt.Errorf("创建本地文件失败: %w", err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("写入本地文件失败: %w", err)
	}
	return nil
}
