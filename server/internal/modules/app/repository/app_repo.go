package repository

import (
    "context"
    "fmt"
    "strconv"

    "github.com/aarcsx/krisho-backend/internal/database"
    "github.com/aarcsx/krisho-backend/internal/models"
)

type AppRepository struct {
    db *database.DB
}

func NewAppRepository(db *database.DB) *AppRepository {
    return &AppRepository{db: db}
}

func (r *AppRepository) GetActiveCrops(ctx context.Context) ([]*models.Crop, error) {
    rows, err := r.db.Pool.Query(ctx, `SELECT name, emoji FROM crops WHERE is_active = true ORDER BY display_order ASC`)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var crops []*models.Crop
    for rows.Next() {
        crop := &models.Crop{}
        if err := rows.Scan(&crop.Name, &crop.Emoji); err != nil {
            return nil, err
        }
        crops = append(crops, crop)
    }
    if err := rows.Err(); err != nil {
        return nil, err
    }
    return crops, nil
}

func (r *AppRepository) GetAppConfig(ctx context.Context) (*models.AppConfig, error) {
    rows, err := r.db.Pool.Query(ctx, `SELECT key, value FROM app_settings`)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    settings := make(map[string]string)
    for rows.Next() {
        var key, value string
        if err := rows.Scan(&key, &value); err != nil {
            return nil, err
        }
        settings[key] = value
    }
    if err := rows.Err(); err != nil {
        return nil, err
    }

    minimumValue, ok := settings["minimum_version_code"]
    if !ok {
        return nil, fmt.Errorf("app config is missing minimum_version_code")
    }

    minimumVersionCode, err := strconv.Atoi(minimumValue)
    if err != nil {
        return nil, fmt.Errorf("invalid minimum_version_code: %w", err)
    }

    return &models.AppConfig{
        MinimumVersionCode: minimumVersionCode,
        LatestVersionName:  settings["latest_version_name"],
        UpdateURL:          settings["update_url"],
        Message:            settings["message"],
    }, nil
}
