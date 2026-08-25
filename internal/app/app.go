package app

import (
	"database/sql"
	"net"
	"time"

	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"

	"hamrahTask1/internal/adaptors/handler"
	"hamrahTask1/internal/adaptors/repository/postgress"
	cache "hamrahTask1/internal/adaptors/repository/redis"
	"hamrahTask1/internal/service"
	"hamrahTask1/pkg/config"
	"hamrahTask1/pkg/logger"
)

type App struct {
	cfg        *config.Config
	grpcServer *grpc.Server
	db         *sql.DB
	redis      *redis.Client
	log        *logger.Logger
}

func New(cfg *config.Config) *App {
	logg := logger.New(nil)

	db, err := sql.Open("pgx", cfg.PostgresDSN)
	if err != nil {
		logg.Fatal("Failed to open database", err)
	}

	var pingErr error
	for i := 0; i < 5; i++ {
		pingErr = db.Ping()
		if pingErr == nil {
			break
		}
		logg.Info("Waiting for database to become ready...")
		time.Sleep(3 * time.Second)
	}
	if pingErr != nil {
		logg.Fatal("Database is unreachable after retries", pingErr)
	}

	logRepo := postgress.NewPostgresLogRepository(db)
	logg = logger.New(logRepo)

	logg.Info("Successfully connected to PostgreSQL.")

	rdb := redis.NewClient(&redis.Options{Addr: cfg.RedisAddr})
	logg.Info("Successfully connected to Redis.")

	userRepo := postgress.NewPostgresUserRepository(db)
	chatRepo := postgress.NewPostgresChatRepository(db)
	sessionCache := cache.NewRedisCache(rdb)

	chatService := service.NewChatService(userRepo, chatRepo, sessionCache, logg)
	grpcHandler := handler.NewGRPCHandler(chatService, logg)

	grpcServer := grpc.NewServer()
	handler.RegisterChatbotServiceServer(grpcServer, grpcHandler)

	return &App{
		cfg:        cfg,
		grpcServer: grpcServer,
		db:         db,
		redis:      rdb,
		log:        logg,
	}
}

func (a *App) Start() {
	lis, err := net.Listen("tcp", ":"+a.cfg.AppPort)
	if err != nil {
		a.log.Fatal("Failed to listen", err)
	}
	a.log.Info("gRPC server listening on port " + a.cfg.AppPort)

	if err := a.grpcServer.Serve(lis); err != nil {
		a.log.Fatal("Failed to serve", err)
	}
}

func (a *App) Stop() {
	a.log.Info("Shutting down application...")
	a.grpcServer.GracefulStop()
	if a.db != nil {
		a.db.Close()
	}
	if a.redis != nil {
		a.redis.Close()
	}
}

// graceful shutdown
// career lader entekhab yek zel
