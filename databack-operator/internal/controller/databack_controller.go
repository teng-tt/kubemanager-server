/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"os"
	"os/exec"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	operatorkubemanagercomv1beta1 "kubemanager.com/operator-databackup/api/v1beta1"
)

var logger *zap.Logger

func init() {
	cfg := zap.Config{
		Encoding:         "console",
		Level:            zap.NewAtomicLevelAt(zapcore.InfoLevel),
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder   // 设置时间格式
	cfg.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder // 显示短路径的调用者信息

	zaplogger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	logger = zaplogger
}

// DatabackReconciler reconciles a Databack object
type DatabackReconciler struct {
	client.Client
	Scheme      *runtime.Scheme
	BackupQueue map[string]operatorkubemanagercomv1beta1.Databack
	Wg          sync.WaitGroup
	Tickers     []*time.Ticker
	lock        sync.RWMutex
}

//+kubebuilder:rbac:groups=operator.kubemanager.com,resources=databacks,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=operator.kubemanager.com,resources=databacks/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=operator.kubemanager.com,resources=databacks/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Databack object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.4/pkg/reconcile
func (r *DatabackReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)
	// 查询资源是否存在 不存在 说明资源被删除
	var databackK8s operatorkubemanagercomv1beta1.Databack
	err := r.Client.Get(ctx, req.NamespacedName, &databackK8s)
	if err != nil {
		if errors.IsNotFound(err) {
			logger.Sugar().Infof("%s 停止", req.Name)
			r.DeleteQueue(databackK8s)
			return ctrl.Result{}, nil
		}
		logger.Sugar().Errorf("%s 异常 %v", req.Name, err)
		return ctrl.Result{}, err
	}
	if lastDataback, ok := r.BackupQueue[databackK8s.Name]; ok {
		isEqual := reflect.DeepEqual(lastDataback.Spec, databackK8s.Spec)
		if isEqual {
			return ctrl.Result{}, nil
		}
	}
	// create/update logic
	r.AddQueue(databackK8s)
	return ctrl.Result{}, nil
}

func (r *DatabackReconciler) AddQueue(databack operatorkubemanagercomv1beta1.Databack) {
	if r.BackupQueue == nil {
		r.BackupQueue = make(map[string]operatorkubemanagercomv1beta1.Databack)
	}
	r.BackupQueue[databack.Name] = databack
	r.StopLoop()
	go r.RunLoop()
}

func (r *DatabackReconciler) DeleteQueue(databack operatorkubemanagercomv1beta1.Databack) {
	delete(r.BackupQueue, databack.Name)
	r.StopLoop()
	go r.RunLoop()
}

// StopLoop 开始循环
func (r *DatabackReconciler) StopLoop() {
	for _, ticker := range r.Tickers {
		if ticker != nil {
			ticker.Stop()
		}
	}
}

// RunLoop 停止循环
func (r *DatabackReconciler) RunLoop() {
	for name, databack := range r.BackupQueue {
		if !databack.Spec.Enable {
			logger.Sugar().Infof("%s 未开启", name)
			databack.Status.Active = false
			r.UpdateStatus(databack)
			continue
		}
		delay := r.getDelaySeconds(databack.Spec.StartTime)
		if delay.Hours() < 1 {
			logger.Sugar().Infof("%s 将在 %.1f分钟后执行", name, delay.Minutes())
		} else {
			logger.Sugar().Infof("%s 将在 %.1f小时后执行", name, delay.Hours())
		}
		// 更新状态
		databack.Status.Active = true
		nextTime := r.getNextTime(delay.Seconds())
		databack.Status.NexTime = nextTime.Unix()
		r.UpdateStatus(databack)
		ticker := time.NewTicker(delay)
		r.Tickers = append(r.Tickers, ticker)
		r.Wg.Add(1)
		go func(databack operatorkubemanagercomv1beta1.Databack) {
			defer r.Wg.Done()
			for {
				<-ticker.C
				// 重置ticker
				ticker.Reset(time.Duration(databack.Spec.Period) * time.Minute)
				logger.Sugar().Infof("%s 将在 %d 分钟后循环执行", databack.Name, databack.Spec.Period)
				databack.Status.Active = true
				databack.Status.NexTime = r.getNextTime(float64(databack.Spec.Period) * 60).Unix()
				// 备份任务
				err := r.DumpWithUploadOss(databack)
				if err != nil {
					logger.Sugar().Errorf("%s同步数据报错:%v", databack.Name, err)
					databack.Status.LastBackupResult = fmt.Sprintf("databack failed %v", err)
				} else {
					logger.Sugar().Infof("%s同步数据成功", databack.Name)
					databack.Status.LastBackupResult = "databack successful"
				}
				// 更新备份状态
				r.UpdateStatus(databack)
			}
		}(databack)
	}
	r.Wg.Wait()
}

// 获取第一次启动的延时时间(秒)
func (r *DatabackReconciler) getDelaySeconds(startTime string) time.Duration {
	// 计算小时和分钟
	times := strings.Split(startTime, ":")
	expectedHour, _ := strconv.Atoi(times[0])
	expectedMin, _ := strconv.Atoi(times[1])
	now := time.Now().Truncate(time.Second)
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	todayEnd := todayStart.Add(24 * time.Hour)
	var seconds int
	// 期望时间点 一天的秒数
	expectedDuration := time.Hour*time.Duration(expectedHour) + time.Minute*time.Duration(expectedMin)
	// 当前时间点 一天的秒数
	curDuration := time.Hour*time.Duration(now.Hour()) + time.Minute*time.Duration(now.Minute())
	// 当前时间已经过了预定时间 明天执行
	if curDuration >= expectedDuration {
		seconds = int(todayEnd.Add(expectedDuration).Sub(now).Seconds())
	} else {
		// 今天执行
		seconds = int(todayStart.Add(expectedDuration).Sub(now).Seconds())
	}
	return time.Second * time.Duration(seconds)
}

// 当前时间x秒后的时间
// 返回标准时间格式
func (r *DatabackReconciler) getNextTime(seconds float64) time.Time {
	currentTime := time.Now()
	return currentTime.Add(time.Second * time.Duration(seconds))
}

func (r *DatabackReconciler) UpdateStatus(backup operatorkubemanagercomv1beta1.Databack) {
	r.lock.Lock()
	defer r.lock.Unlock()
	ctx := context.TODO()
	namespacedName := types.NamespacedName{
		Name:      backup.Name,
		Namespace: backup.Namespace,
	}
	var dataBackupK8s operatorkubemanagercomv1beta1.Databack
	err := r.Get(ctx, namespacedName, &dataBackupK8s)
	if err != nil {
		logger.Sugar().Error(err)
		return
	}
	//状态更新为激活
	dataBackupK8s.Status = backup.Status
	err = r.Client.Status().Update(ctx, &dataBackupK8s)
	if err != nil {
		logger.Sugar().Error(err)
		return
	}
}

// DumpWithUploadOss mysql数据导出+数据上报到OSS
func (r *DatabackReconciler) DumpWithUploadOss(backup operatorkubemanagercomv1beta1.Databack) error {
	defer func() {
		if err := recover(); err != nil {
			logger.Sugar().Errorf("run time panic: %v", err)
		}
	}()
	//dump
	mysqlHost := backup.Spec.Origin.Host
	mysqlPort := backup.Spec.Origin.Port
	mysqlUsername := backup.Spec.Origin.Username
	mysqlPassword := backup.Spec.Origin.Password
	now := time.Now()
	backupDate := fmt.Sprintf("%02d-%02d", now.Month(), now.Day())
	folderPath := fmt.Sprintf("/tmp/%s/%s/", backup.Name, backupDate)
	//创建文件夹
	if _, err := os.Stat(folderPath); err != nil {
		if errx := os.MkdirAll(folderPath, 0700); errx == nil {
			logger.Sugar().Infof("created dir %s", folderPath)
		} else {
			logger.Sugar().Errorf("%s ：%v", backup.Name, errx)
			return errx
		}
	}
	//计算当天同步的文件个数
	files, err := ioutil.ReadDir(folderPath)
	if err != nil {
		logger.Sugar().Errorf("%s ：%v", backup.Name, err)
		return err
	}
	number := len(files) + 1
	filename := fmt.Sprintf("%s#%d.sql", folderPath, number)
	dumpCmd := fmt.Sprintf("mysqldump -h%s -P%d -u%s -p%s --all-databases > %s",
		mysqlHost, mysqlPort, mysqlUsername, mysqlPassword, filename)
	logger.Sugar().Infof("%s %s", backup.Name, dumpCmd)
	command := exec.Command("bash", "-c", dumpCmd)
	_, err = command.Output() // 执行命令并获取输出
	if err != nil {
		logger.Sugar().Errorf("%s %v", backup.Name, err)
		return err
	}
	//upload
	endpoint := backup.Spec.Destination.Endpoint
	accessKey := backup.Spec.Destination.AccessKey
	asscessSecret := backup.Spec.Destination.AccessSecret
	bucketName := backup.Spec.Destination.BucketName
	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, asscessSecret, ""),
		Secure: false,
	})
	if err != nil {
		logger.Sugar().Infof("%s %v", backup.Name, err)
		return err
	}
	object, err := os.Open(filename)
	if err != nil {
		logger.Sugar().Infof("%s %v", backup.Name, err)
		return err
	}
	ctx := context.TODO()
	_, err = minioClient.PutObject(ctx, bucketName, filename, object, -1, minio.PutObjectOptions{})
	if err != nil {
		logger.Sugar().Infof("%s %v", backup.Name, err)
		return err
	}
	return err
}

// SetupWithManager sets up the controller with the Manager.
func (r *DatabackReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&operatorkubemanagercomv1beta1.Databack{}).
		Complete(r)
}
