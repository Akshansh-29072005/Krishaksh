package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/aarcsx/krisho-backend/internal/observability"
	"github.com/redis/go-redis/v9"
)

type Service interface {
	Get(ctx context.Context, lat, lon string) (Response, string, error)
}

type serviceImpl struct {
	redis       *redis.Client
	httpClient  *http.Client
	geoClient   *http.Client
	cacheTTL    time.Duration
	ghPrecision int
}

func NewService(redisClient *redis.Client) Service {
	precision := 5 // true geohash precision; roughly ~4.9km cells near equator
	if raw := os.Getenv("WEATHER_GEOHASH_PRECISION"); raw != "" {
		if v, err := strconv.Atoi(raw); err == nil && v >= 1 && v <= 8 {
			precision = v
		}
	}
	ttl := 10 * time.Minute
	if raw := os.Getenv("WEATHER_CACHE_TTL_SECONDS"); raw != "" {
		if v, err := strconv.Atoi(raw); err == nil && v > 0 {
			ttl = time.Duration(v) * time.Second
		}
	}
	return &serviceImpl{
		redis:       redisClient,
		httpClient:  &http.Client{Timeout: 10 * time.Second},
		geoClient:   &http.Client{Timeout: 10 * time.Second},
		cacheTTL:    ttl,
		ghPrecision: precision,
	}
}

func (s *serviceImpl) Get(ctx context.Context, lat, lon string) (Response, string, error) {
	latF, err := strconv.ParseFloat(lat, 64)
	if err != nil {
		return Response{}, "", fmt.Errorf("invalid lat")
	}
	lonF, err := strconv.ParseFloat(lon, 64)
	if err != nil {
		return Response{}, "", fmt.Errorf("invalid lon")
	}

	gh := EncodeGeohash(latF, lonF, s.ghPrecision)
	cacheKey := "weather:gh:" + gh
	if cached, err := s.redis.Get(ctx, cacheKey).Result(); err == nil && cached != "" {
		var out Response
		if json.Unmarshal([]byte(cached), &out) == nil {
			observability.M.Inc("weather_cache_hits_total")
			return out, "redis-cache", nil
		}
	}
	observability.M.Inc("weather_cache_misses_total")

	weatherURL := fmt.Sprintf(
		"https://api.open-meteo.com/v1/forecast?latitude=%s&longitude=%s&current=temperature_2m,relative_humidity_2m,weather_code,wind_speed_10m&timezone=auto",
		lat, lon,
	)
	wreq, _ := http.NewRequestWithContext(ctx, http.MethodGet, weatherURL, nil)
	weatherResp, err := s.httpClient.Do(wreq)
	if err != nil {
		return fallback(), "degraded", nil
	}
	defer weatherResp.Body.Close()

	var weatherData struct {
		Current struct {
			Temperature float64 `json:"temperature_2m"`
			Humidity    float64 `json:"relative_humidity_2m"`
			WeatherCode int     `json:"weather_code"`
			WindSpeed   float64 `json:"wind_speed_10m"`
		} `json:"current"`
	}
	if err := json.NewDecoder(weatherResp.Body).Decode(&weatherData); err != nil {
		return fallback(), "degraded", nil
	}

	locationName := "Your Location"
	geoURL := fmt.Sprintf("https://nominatim.openstreetmap.org/reverse?lat=%s&lon=%s&format=json&accept-language=en", lat, lon)
	greq, _ := http.NewRequestWithContext(ctx, http.MethodGet, geoURL, nil)
	greq.Header.Set("User-Agent", "Krishaksh-App/1.0")
	geoResp, err := s.geoClient.Do(greq)
	if err == nil {
		defer geoResp.Body.Close()
		var geoData struct {
			Address struct {
				City, Town, Village, County, State string
			} `json:"address"`
		}
		if json.NewDecoder(geoResp.Body).Decode(&geoData) == nil {
			place := geoData.Address.City
			if place == "" {
				place = geoData.Address.Town
			}
			if place == "" {
				place = geoData.Address.Village
			}
			if place == "" {
				place = geoData.Address.County
			}
			region := geoData.Address.State
			if place != "" && region != "" {
				locationName = place + ", " + region
			} else if place != "" {
				locationName = place
			} else if region != "" {
				locationName = region
			}
		}
	}

	out := Response{
		Temperature:  fmt.Sprintf("%.0f°C", weatherData.Current.Temperature),
		Condition:    mapWMO(weatherData.Current.WeatherCode),
		Humidity:     fmt.Sprintf("%.0f%%", weatherData.Current.Humidity),
		WindSpeed:    fmt.Sprintf("%.0f km/h", weatherData.Current.WindSpeed),
		LocationName: locationName,
	}
	if b, err := json.Marshal(out); err == nil {
		_ = s.redis.Set(ctx, cacheKey, string(b), s.cacheTTL).Err()
	}
	return out, "origin", nil
}

func fallback() Response {
	return Response{
		Temperature:  "--°C",
		Condition:    "Unknown",
		Humidity:     "--%",
		WindSpeed:    "-- km/h",
		LocationName: "Unknown Location",
	}
}

func mapWMO(wc int) string {
	switch {
	case wc == 0:
		return "Clear"
	case wc <= 3:
		return "Cloudy"
	case wc >= 45 && wc <= 48:
		return "Foggy"
	case wc >= 51 && wc <= 57:
		return "Drizzle"
	case wc >= 61 && wc <= 67:
		return "Rain"
	case wc >= 71 && wc <= 77:
		return "Snow"
	case wc >= 80 && wc <= 82:
		return "Showers"
	case wc >= 95:
		return "Thunderstorm"
	default:
		return "Clear"
	}
}
