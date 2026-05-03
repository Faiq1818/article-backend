package article

import "article/internal/apperror"

func (s *Service) DeleteArticle(slug string) error {
	err := s.Repo.DeleteArticle(slug)
	if err != nil {
		s.Logger.Error("failed delete article", "error", err)
		return &apperror.AppError{
			Message: "Failed to delete article",
			Code:    500,
			Err:     err,
		}
	}

	s.Logger.Info("Successfully deleted article")
	return nil
}
