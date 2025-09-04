package service

import (
	"agdemo/internal/biz"
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"go.opentelemetry.io/otel"

	pb "agdemo/api/blog/v1"
)

func NewBlogService(article *biz.ArticleUsecase, logger log.Logger) *BlogService {
	return &BlogService{
		article: article,
		log:     log.NewHelper(logger),
	}
}

func (s *BlogService) CreateArticle(ctx context.Context, req *pb.CreateArticleRequest) (*pb.CreateArticleReply, error) {
	s.log.Infof("input data %v", req)
	article := biz.Article{
		Title:   req.Title,
		Content: req.Content,
	}
	if err := article.Validate(); err != nil {
		return nil, err
	}

	err := s.article.Create(ctx, &article)
	if err != nil {
		return nil, err
	}
	// 3. 转换为响应
	return &pb.CreateArticleReply{
		Article: &pb.Article{
			Id:      article.Id,
			Title:   article.Title,
			Content: article.Content,
		},
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
		return nil, pb.ErrorBlogInvalidId("invalid article id")
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
