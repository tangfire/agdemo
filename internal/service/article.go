package service

import (
	"agdemo/internal/biz"
	"agdemo/internal/code"
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "agdemo/api/blog/v1"
)

func NewBlogService(article *biz.ArticleUsecase, logger log.Logger) *BlogService {
	return &BlogService{
		article: article,
		log:     log.NewHelper(logger),
	}
}

func (s *BlogService) CreateArticle(ctx context.Context, req *pb.CreateArticleRequest) (*pb.CreateArticleReply, error) {
	// 1. 安全日志
	s.log.WithContext(ctx).Infof(
		"CreateArticle title_len:%d content_len:%d",
		len(req.Title),
		len(req.Content),
	)

	// 2. 创建并校验领域对象
	article := &biz.Article{
		Title:   req.Title,
		Content: req.Content,
	}

	// 3. 执行业务逻辑
	if err := s.article.Create(ctx, article); err != nil {
		return nil, status.Error(codes.Internal, "创建文章失败")
	}

	// 4. 转换响应
	return &pb.CreateArticleReply{
		Article: article.ToProto(),
	}, nil
}

func (s *BlogService) UpdateArticle(ctx context.Context, req *pb.UpdateArticleRequest) (*pb.UpdateArticleReply, error) {
	s.log.Infof("input data %v", req)
	article := biz.Article{
		Title:   req.Title,
		Content: req.Content,
	}
	if err := article.Validate(); err != nil {
		return nil, err
	}
	err := s.article.Update(ctx, req.Id, &article)
	return &pb.UpdateArticleReply{}, err
}

func (s *BlogService) DeleteArticle(ctx context.Context, req *pb.DeleteArticleRequest) (*pb.DeleteArticleReply, error) {
	s.log.Infof("input data %v", req)
	err := s.article.Delete(ctx, req.Id)
	return &pb.DeleteArticleReply{}, err
}

func (s *BlogService) GetArticle(ctx context.Context, req *pb.GetArticleRequest) (*pb.GetArticleReply, error) {
	if req.Id < 1 {
		return nil, code.InvalidId
	}
	tr := otel.Tracer("api")
	ctx, span := tr.Start(ctx, "GetArticle")
	defer span.End()
	p, err := s.article.Get(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &pb.GetArticleReply{Article: &pb.Article{Id: p.Id, Title: p.Title, Content: p.Content, Like: p.Like}}, nil
}

func (s *BlogService) ListArticle(ctx context.Context, req *pb.ListArticleRequest) (*pb.ListArticleReply, error) {
	ps, err := s.article.List(ctx)
	reply := &pb.ListArticleReply{}
	for _, p := range ps {
		reply.Results = append(reply.Results, &pb.Article{
			Id:      p.Id,
			Title:   p.Title,
			Content: p.Content,
		})
	}
	return reply, err
}

func (s *BlogService) ArticleCastJson(ctx context.Context, req *pb.ArticleCastJsonRequest) (*pb.ArticleCastJsonReply, error) {
	article := &biz.Article{
		Id:      req.Id,
		Title:   req.Title,
		Content: req.Content,
	}
	json, err := s.article.CastJson(ctx, article)
	if err != nil {
		return nil, status.Error(codes.Internal, "转换文章json失败")
	}
	return &pb.ArticleCastJsonReply{Json: json}, nil
}
