CREATE TABLE IF NOT EXISTS crops (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) UNIQUE NOT NULL,
    emoji VARCHAR(10) NOT NULL DEFAULT '',
    slug VARCHAR(100) UNIQUE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    display_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS app_settings (
    id SERIAL PRIMARY KEY,
    key VARCHAR(100) UNIQUE NOT NULL,
    value TEXT NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

INSERT INTO crops (name, emoji, slug, display_order) VALUES
    ('Rice', '🌾', 'rice', 1),
    ('Brinjal', '🍆', 'brinjal', 2),
    ('Maize', '🌽', 'maize', 3),
    ('Tomato', '🍅', 'tomato', 4),
    ('Potato', '🥔', 'potato', 5),
    ('Wheat', '🌾', 'wheat', 6)
ON CONFLICT (name) DO NOTHING;

INSERT INTO app_settings (key, value) VALUES
    ('minimum_version_code', '1'),
    ('latest_version_name', '1.0'),
    ('update_url', 'https://play.google.com/store/apps/details?id=com.aarcsx.krisho'),
    ('message', 'A newer Krisho version is required for security patches and bug fixes. Please update to continue.')
ON CONFLICT (key) DO NOTHING;
