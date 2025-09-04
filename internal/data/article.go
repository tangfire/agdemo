package data

import (
	"agdemo/internal/biz"
	"context"
	"github.com/go-kratos/kratos/v2/log"
	"time"
)

// data/article.go
// 将数据库模型改为私有（首字母小写）
type article struct { // 注意首字母小写
	Id        int64     `gorm:"primaryKey"`
	Title     string    `gorm:"size:100"`
	Content   string    `gorm:"type:text"`
	LikeCount int64     `gorm:"column:like_count"`
	CreatedAt time.Time `gorm:"column:created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
}

// 实现TableName接口（可选）
func (article) TableName() string {
	return "article"
}

type articleRepo struct {
	data *Data
	log  *log.Helper
}

// data/article.go
// 转换方法整合到Repo中
func (r *articleRepo) toDomain(a *article) *biz.Article {
	return &biz.Article{
		Id:        a.Id,
		Title:     a.Title,
		Content:   a.Content,
		Like:      a.LikeCount,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}
}

func (r *articleRepo) toModel(a *biz.Article) *article {
	return &article{
		Id:        a.Id,
		Title:     a.Title,
		Content:   a.Content,
		LikeCount: a.Like,
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
	}
}

func NewArticleRepo(data *Data, logger log.Logger) biz.ArticleRepo {
	return &articleRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *articleRepo) ListArticle(ctx context.Context) ([]*biz.Article, error) {
	var list []*article
	if err := r.data.db.WithContext(ctx).Find(&list).Error; err != nil {
		r.log.Errorf("List error: %v", err)
		return nil, err
	}

	result := make([]*biz.Article, 0, len(list))
	for _, item := range list {
		result = append(result, r.toDomain(item))
	}
	return result, nil
}

func (r *articleRepo) GetArticle(ctx context.Context, id int64) (*biz.Article, error) {
	var a article
	if err := r.data.db.WithContext(ctx).First(&a, id).Error; err != nil {
		r.log.Errorf("Get error: %v", err)
		return nil, err
	}
	return r.toDomain(&a), nil
}

func (r *articleRepo) CreateArticle(ctx context.Context, a *biz.Article) error {
	model := r.toModel(a)
	err := r.data.db.WithContext(ctx).Create(&model).Error
	if err != nil {
		r.log.Errorf("Create error: %v", err)
		return err
	}
	a.Id = model.Id
	return nil
}

func (r *articleRepo) UpdateArticle(ctx context.Context, id int64, a *biz.Article) error {
	m := r.toModel(a)
	m.Id = id // 确保更新目标ID正确
	return r.data.db.WithContext(ctx).Updates(m).Error
}

func (r *articleRepo) DeleteArticle(ctx context.Context, id int64) error {
	return r.data.db.WithContext(ctx).Delete(&article{}, id).Error
}
