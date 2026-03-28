package article

import (
	"article/internal/models"
	"article/internal/request_type"
	"errors"
	"log/slog"
	"testing"
)

var assertError = errors.New("test error")

type mockArticleRepository struct {
	getManyArticleFunc func(limit int, offset int) ([]models.Article, int, error)
}

func (m *mockArticleRepository) GetManyArticle(limit int, offset int) ([]models.Article, int, error) {
	if m.getManyArticleFunc != nil {
		return m.getManyArticleFunc(limit, offset)
	}
	return nil, 0, nil
}

func (m *mockArticleRepository) GetArticleBySlug(slug string) (models.Article, error) {
	return models.Article{}, nil
}

func (m *mockArticleRepository) SaveArticle(req requesttype.SaveArticleRequest, imgUrl string, slugGenerate string) error {
	return nil
}

func (m *mockArticleRepository) PutArticle(req requesttype.PutArticleRequest, imgUrl string, slugGenerate string, oldSlug string) error {
	return nil
}

func (m *mockArticleRepository) DeleteArticle(slug string) error {
	return nil
}

func (m *mockArticleRepository) AdminGetManyArticle(limit int, offset int) ([]models.Article, int, error) {
	return nil, 0, nil
}

func TestGetArticles(t *testing.T) {
	logger := slog.New(slog.DiscardHandler)

	tests := []struct {
		name         string
		page         int
		limit        int
		mockReturn   []models.Article
		mockTotal    int
		mockError    error
		expectedMeta models.PaginationMeta
		expectError  bool
	}{
		{
			name:  "successful get articles with valid page and limit",
			page:  1,
			limit: 10,
			mockReturn: []models.Article{
				{ID: "1", Title: "Article 1"},
				{ID: "2", Title: "Article 2"},
			},
			mockTotal: 25,
			mockError: nil,
			expectedMeta: models.PaginationMeta{
				CurrentPage: 1,
				Limit:       10,
				TotalItems:  25,
				TotalPages:  3,
				HasNext:     true,
				HasPrev:     false,
			},
			expectError: false,
		},
		{
			name:  "successful get articles on second page",
			page:  2,
			limit: 10,
			mockReturn: []models.Article{
				{ID: "11", Title: "Article 11"},
			},
			mockTotal: 25,
			mockError: nil,
			expectedMeta: models.PaginationMeta{
				CurrentPage: 2,
				Limit:       10,
				TotalItems:  25,
				TotalPages:  3,
				HasNext:     true,
				HasPrev:     true,
			},
			expectError: false,
		},
		{
			name:  "successful get articles on last page",
			page:  3,
			limit: 10,
			mockReturn: []models.Article{
				{ID: "21", Title: "Article 21"},
			},
			mockTotal: 25,
			mockError: nil,
			expectedMeta: models.PaginationMeta{
				CurrentPage: 3,
				Limit:       10,
				TotalItems:  25,
				TotalPages:  3,
				HasNext:     false,
				HasPrev:     true,
			},
			expectError: false,
		},
		{
			name:  "default page and limit when invalid",
			page:  0,
			limit: -5,
			mockReturn: []models.Article{
				{ID: "1", Title: "Article 1"},
			},
			mockTotal: 15,
			mockError: nil,
			expectedMeta: models.PaginationMeta{
				CurrentPage: 1,
				Limit:       10,
				TotalItems:  15,
				TotalPages:  2,
				HasNext:     true,
				HasPrev:     false,
			},
			expectError: false,
		},
		{
			name:        "repository returns error",
			page:        1,
			limit:       10,
			mockReturn:  nil,
			mockTotal:   0,
			mockError:   assertError,
			expectError: true,
		},
		{
			name:       "empty articles list",
			page:       1,
			limit:      10,
			mockReturn: []models.Article{},
			mockTotal:  0,
			mockError:  nil,
			expectedMeta: models.PaginationMeta{
				CurrentPage: 1,
				Limit:       10,
				TotalItems:  0,
				TotalPages:  0,
				HasNext:     false,
				HasPrev:     false,
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &mockArticleRepository{
				getManyArticleFunc: func(limit int, offset int) ([]models.Article, int, error) {
					return tt.mockReturn, tt.mockTotal, tt.mockError
				},
			}

			service := &Service{
				Repo:   mockRepo,
				Logger: logger,
			}

			articles, meta, err := service.GetArticles(tt.page, tt.limit)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if meta.CurrentPage != tt.expectedMeta.CurrentPage {
				t.Errorf("CurrentPage = %d, want %d", meta.CurrentPage, tt.expectedMeta.CurrentPage)
			}

			if meta.Limit != tt.expectedMeta.Limit {
				t.Errorf("Limit = %d, want %d", meta.Limit, tt.expectedMeta.Limit)
			}

			if meta.TotalItems != tt.expectedMeta.TotalItems {
				t.Errorf("TotalItems = %d, want %d", meta.TotalItems, tt.expectedMeta.TotalItems)
			}

			if meta.TotalPages != tt.expectedMeta.TotalPages {
				t.Errorf("TotalPages = %d, want %d", meta.TotalPages, tt.expectedMeta.TotalPages)
			}

			if meta.HasNext != tt.expectedMeta.HasNext {
				t.Errorf("HasNext = %v, want %v", meta.HasNext, tt.expectedMeta.HasNext)
			}

			if meta.HasPrev != tt.expectedMeta.HasPrev {
				t.Errorf("HasPrev = %v, want %v", meta.HasPrev, tt.expectedMeta.HasPrev)
			}

			if len(articles) != len(tt.mockReturn) {
				t.Errorf("articles length = %d, want %d", len(articles), len(tt.mockReturn))
			}
		})
	}
}
