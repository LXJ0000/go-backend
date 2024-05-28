package bootstrap

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/LXJ0000/go-backend/internal/domain"
	repository "github.com/LXJ0000/go-backend/internal/repository"
	"github.com/LXJ0000/go-backend/internal/usecase"
	"github.com/LXJ0000/go-backend/pkg/cache"
	"github.com/LXJ0000/go-backend/pkg/orm"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/robfig/cron/v3"
)

type Job interface {
	Name() string
	Run() error
}

type RankJob struct {
	rankUsecase domain.RankUsecase
	timeout     time.Duration
}

func NewRankJob(rankUsecase domain.RankUsecase, timeout time.Duration) *RankJob {
	return &RankJob{rankUsecase: rankUsecase, timeout: timeout}
}

func (j *RankJob) Name() string {
	return `post_rank`
}

func (j *RankJob) Run() error {
	ctx, cancel := context.WithTimeout(context.Background(), j.timeout)
	defer cancel()
	return j.rankUsecase.TopN(ctx)
}

type CronRankJob struct {
	job        *RankJob
	prometheus prometheus.Summary
}

func NewCronRankJob(job *RankJob) *CronRankJob {
	p := prometheus.NewSummary(prometheus.SummaryOpts{
		Namespace: "jannan",
		Subsystem: "go_backend",
		Help:      "summary of job",
		Name:      "cron_job",
		ConstLabels: map[string]string{
			"name": job.Name(),
		}})
	prometheus.MustRegister(p)
	return &CronRankJob{job: job, prometheus: p}
}

var rankJobMutex = sync.Mutex{}

func (c *CronRankJob) Run() {
	rankJobMutex.Lock() // 避免第一个任务还未计算完毕就开始第二个任务的计算
	defer rankJobMutex.Unlock()
	slog.Debug("定时任务 启动！", "job", c.job.Name())
	begin := time.Now()
	defer func() {
		duration := time.Since(begin).Milliseconds()
		c.prometheus.Observe(float64(duration))
		slog.Debug("定时任务 结束！", "job", c.job.Name(), "duration", time.Since(begin))

	}()
	if err := c.job.Run(); err != nil {
		slog.Error("Cronjob Run Error", "Error", err.Error(), "Name", c.job.Name())
	}
}

func InitCronRankJob(rankJob *RankJob) *cron.Cron {
	expr := cron.New(cron.WithSeconds())

	cronRankJob := NewCronRankJob(rankJob)
	go cronRankJob.Run()
	if _, err := expr.AddJob("@every 5m", cronRankJob); err != nil {
		panic(err)
	}

	return expr
}

// NewCron 
func NewCron(localCache cache.LocalCache, redisCache cache.RedisCache, dao orm.Database) *cron.Cron {
	timeout := time.Minute
	interactionRepository := repository.NewInteractionRepository(dao, redisCache)
	postRepository := repository.NewPostRepository(dao, redisCache)
	postRankRepository := repository.NewPostRankRepository(localCache, redisCache)
	postRankUsecase := usecase.NewPostRankUsecase(interactionRepository, postRepository, postRankRepository, timeout)
	rankJob := NewRankJob(postRankUsecase, timeout)
	return InitCronRankJob(rankJob)
}
