package service

import (
    "context"
    "testing"

    "filmhub/internal/models"
)

// stubFilmRepo is an in-memory implementation of FilmRepository used in tests.
type stubFilmRepo struct {
    nextID int
    films  map[int]models.Film
}

func newStubFilmRepo() *stubFilmRepo {
    return &stubFilmRepo{
        nextID: 1,
        films:  make(map[int]models.Film),
    }
}

func (s *stubFilmRepo) CreateFilm(_ context.Context, req *models.FilmRequest) (int, error) {
    id := s.nextID
    s.nextID++
    s.films[id] = models.Film{
        ID:          id,
        Title:       req.Title,
        Description: req.Description,
        ReleaseDate: req.ReleaseDate,
    }
    return id, nil
}

func (s *stubFilmRepo) GetFilmByID(_ context.Context, id int) (*models.Film, error) {
    f, ok := s.films[id]
    if !ok {
        return nil, context.Canceled // placeholder error
    }
    return &f, nil
}

func (s *stubFilmRepo) SearchFilms(_ context.Context, query string) ([]models.Film, error) {
    var result []models.Film
    for _, f := range s.films {
        if query == "" || contains(f.Title, query) || contains(f.Description, query) {
            result = append(result, f)
        }
    }
    return result, nil
}

// simplistic substring match helper
func contains(haystack, needle string) bool {
    return len(needle) == 0 || stringContainsCaseInsensitive(haystack, needle)
}

func stringContainsCaseInsensitive(a, b string) bool {
    // naive implementation without extra deps
    A, B := []rune(a), []rune(b)
    for i := 0; i+len(B) <= len(A); i++ {
        match := true
        for j, rb := range B {
            if toLower(A[i+j]) != toLower(rb) {
                match = false
                break
            }
        }
        if match {
            return true
        }
    }
    return false
}

func toLower(r rune) rune {
    if 'A' <= r && r <= 'Z' {
        return r + ('a' - 'A')
    }
    return r
}

func TestFilmService_BasicFlow(t *testing.T) {
    repo := newStubFilmRepo()
    svc := NewFilmService(repo)
    ctx := context.Background()

    req := &models.FilmRequest{Title: "Matrix", Description: "Sci-fi",}
    id, err := svc.CreateFilm(ctx, req)
    if err != nil {
        t.Fatalf("create film failed: %v", err)
    }
    f, err := svc.GetFilm(ctx, id)
    if err != nil || f == nil {
        t.Fatalf("get film failed: %v", err)
    }
    if f.Title != req.Title {
        t.Errorf("expected title %s got %s", req.Title, f.Title)
    }

    res, err := svc.SearchFilms(ctx, "matrix")
    if err != nil || len(res) == 0 {
        t.Fatalf("search failed: %v", err)
    }
} 