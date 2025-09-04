package service

import (
	pb "agdemo/api/blog/v1"
	"agdemo/internal/biz"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewGreeterService, NewBlogService)

type BlogService struct {
	pb.UnimplementedBlogServiceServer

	article *biz.ArticleUsecase

	log *log.Helper
}
