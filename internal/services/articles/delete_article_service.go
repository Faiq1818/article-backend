package article

import "article/internal/pkg"

func (s *Service) DeleteArticle(slug string) error {
	err := s.Repo.DeleteArticle(slug)
	if err != nil {
		s.Logger.Error("failed delete article", "error", err)
		return &pkg.AppError{
			Message: "Gagal menghapus artikel",
			Code:    500,
			Err:     err,
		}
	}

	s.Logger.Info("Successfully deleted article")
	return nil
}
