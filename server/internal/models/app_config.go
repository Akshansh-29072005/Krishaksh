package models

type Crop struct {
    Name  string `json:"name"`
    Emoji string `json:"emoji"`
}

type AppConfig struct {
    MinimumVersionCode int    `json:"minimum_version_code"`
    LatestVersionName  string `json:"latest_version_name,omitempty"`
    UpdateURL          string `json:"update_url,omitempty"`
    Message            string `json:"message,omitempty"`
}
