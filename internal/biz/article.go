package biz

import (
	pb "agdemo/api/blog/v1"
	"context"
	"errors"
	"github.com/go-kratos/kratos/v2/encoding"
	"github.com/go-kratos/kratos/v2/log"
	"time"
)

type Article struct {
	Id        int64
	Title     string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
	Like      int64
}

func (a *Article) ToProto() *pb.Article {
	return &pb.Article{
		Id:      a.Id,
		Title:   a.Title,
		Content: a.Content,
		Like:    a.Like, // 确保不遗漏字段
	}
}

// biz/article.go
func (a *Article) Validate() error {
	if len(a.Title) < 5 {
		return errors.New("标题至少5个字符")
	}
	return nil
}

func (a *Article) IsPopular() bool {
	return a.Like > 1000
}

type ArticleRepo interface {
	// db
	ListArticle(ctx context.Context) ([]*Article, error)
	GetArticle(ctx context.Context, id int64) (*Article, error)
	CreateArticle(ctx context.Context, article *Article) error
	UpdateArticle(ctx context.Context, id int64, article *Article) error
	DeleteArticle(ctx context.Context, id int64) error

	// redis
	GetArticleLike(ctx context.Context, id int64) (rv int64, err error)
	IncArticleLike(ctx context.Context, id int64) error
}

type ArticleUsecase struct {
	repo ArticleRepo
}

func NewArticleUsecase(repo ArticleRepo, logger log.Logger) *ArticleUsecase {
	return &ArticleUsecase{repo: repo}
}

func (uc *ArticleUsecase) List(ctx context.Context) (ps []*Article, err error) {
	ps, err = uc.repo.ListArticle(ctx)
	if err != nil {
		return
	}
	return
}

func (uc *ArticleUsecase) Get(ctx context.Context, id int64) (p *Article, err error) {
	p, err = uc.repo.GetArticle(ctx, id)
	if err != nil {
		return
	}
	err = uc.repo.IncArticleLike(ctx, id)
	if err != nil {
		return
	}
	p.Like, err = uc.repo.GetArticleLike(ctx, id)
	if err != nil {
		return
	}
	return
}

func (uc *ArticleUsecase) Create(ctx context.Context, article *Article) error {
	return uc.repo.CreateArticle(ctx, article)
}

func (uc *ArticleUsecase) Update(ctx context.Context, id int64, article *Article) error {
	return uc.repo.UpdateArticle(ctx, id, article)
}

func (uc *ArticleUsecase) Delete(ctx context.Context, id int64) error {
	return uc.repo.DeleteArticle(ctx, id)
}

func (uc *ArticleUsecase) CastJson(ctx context.Context, article *Article) (string, error) {
	jsonCodec := encoding.GetCodec("json")
	bytes, err := jsonCodec.Marshal(article)
	if err != nil {
		log.Context(ctx).Errorf("CastJson|Marshal err:%v", err)
		return "", err
	}
	return string(bytes), nil
}
