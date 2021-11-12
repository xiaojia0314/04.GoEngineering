package appmanage

import (
	"context"
	"log"

	"golang.org/x/sync/errgroup"
)

type Server interface {
	Serve(ctx context.Context) error
}

type GrpcServer Server
type HttpServer Server

type AppManage struct {
	http HttpServer
	grpc GrpcServer
}

func NewAppManage(http HttpServer, grpc GrpcServer) *AppManage {
	return &AppManage{
		http: http,
		grpc: grpc,
	}
}

func (manage *AppManage) Run(ctx context.Context) {
	group := new(errgroup.Group)
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	group.Go(func() error {
		err := manage.http.Serve(ctx)
		if err != nil {
			cancel()
		}
		return err
	})
	group.Go(func() error {
		err := manage.grpc.Serve(ctx)
		if err != nil {
			cancel()
		}
		return err
	})

	if err := group.Wait(); err != nil {
		log.Printf("Exit Reason: \n\t%v\n", err)
	}
}