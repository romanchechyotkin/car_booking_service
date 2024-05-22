package httpsrv

import (
	"net/http"

	"github.com/romanchechyotkin/car_booking_service/pkg/config"
	minioClient "github.com/romanchechyotkin/car_booking_service/pkg/minio"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"go.uber.org/zap"
)

type MinioClient interface {
	Images() *minio.Client
}

type Server struct {
	log    *zap.Logger
	base   *http.Server
	router *gin.Engine
	minio  MinioClient
}

func NewServer(cfg *config.Config, log *zap.Logger, minioClient *minioClient.Client) (*Server, error) {
	instance := Server{
		log:    log,
		router: gin.New(),
		minio:  minioClient,
	}

	instance.base = &http.Server{
		Addr:    cfg.HTTP.Host + ":" + cfg.HTTP.Port,
		Handler: instance.router,
	}

	instance.log.Debug("http server configuration", zap.Any("cfg", cfg))

	return &instance, nil
}
