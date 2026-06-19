package handler

import (
	"encoding/json"
	"net/http"
	"shorter-url/internal/domain"
	"time"

	"github.com/julienschmidt/httprouter"
)

type ShortUrlResponse struct {
	Id          int64     `json:"id"`
	ShortCode   string    `json:"short_code"`
	OriginalUrl string    `json:"original_url"`
	ExpiredAt   time.Time `json:"expired_at"`
	CreatedAt   time.Time `json:"created_at"`
}

type ShortUrlHandler struct {
	Service domain.ShortUrlsService
}

func NewShortUrlHandler(service domain.ShortUrlsService) *ShortUrlHandler {
	return &ShortUrlHandler{
		Service: service,
	}
}

func (s *ShortUrlHandler) Create(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	var inputData struct {
		OriginalUrl string `json:"original_url"`
	}

	inputRequest := json.NewDecoder(r.Body)
	inputRequest.DisallowUnknownFields()
	if inputRequest.Decode(&inputData) != nil {
		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(http.StatusBadRequest)

		errorResponse := map[string]any{
			"message": "JSON format wrong!",
		}

		json.NewEncoder(w).Encode(errorResponse)

		return
	}

	ctx := r.Context()

	userId := int64(99)
	expiredAt := time.Now().AddDate(0, 1, 0)

	result, err := s.Service.CreateShortUrl(ctx, userId, inputData.OriginalUrl, expiredAt)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(http.StatusBadRequest)

		errorResponse := map[string]any{
			"message": "something wrong when create short code!",
		}

		json.NewEncoder(w).Encode(errorResponse)

		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusCreated)

	data := &ShortUrlResponse{
		Id:          result.Id,
		ShortCode:   result.ShortCode,
		OriginalUrl: result.OriginalUrl,
		ExpiredAt:   result.ExpiredAt,
		CreatedAt:   result.CreatedAt,
	}

	json.NewEncoder(w).Encode(map[string]any{
		"message": "Short code created successfuly",
		"data":    data,
	})
}

func (s *ShortUrlHandler) GetCode(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	ctx := r.Context()

	shortCode := p.ByName("shortCode")
	if shortCode == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]any{
			"message": "Short code cannot be empty",
		})

		return
	}

	_, err := s.Service.GetShortUrlByShortCode(ctx, shortCode)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]any{
			"message": "Short code not found",
		})

		return
	}

	// originalUrl := result.OriginalUrl
	originalUrl := "https://app.notion.com/p/rabs30/Golang-Web-367bc69b480380cb9f37cc1696814776#36bbc69b480380f88aacddc5db48f472"

	http.Redirect(w, r, originalUrl, http.StatusFound)
}
