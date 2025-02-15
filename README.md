# Music Catalog API

API for searching and managing music catalog using Spotify API.

## Features

- Authentication (Sign up, Login, Logout)
- Song search
- Track activity tracking (like/unlike)
- Song recommendations based on liked tracks

## Technology Stack

- Go 1.21
- Gin Web Framework
- GORM
- PostgreSQL
- JWT Authentication
- Spotify API

## Requirements

- Go 1.21+
- PostgreSQL
- [Spotify Developer Account](https://developer.spotify.com/) (for Client ID and Secret)

## Installation

1. Clone repository

```
git clone https://github.com/username/music-catalog.git
cd music-catalog
```

2. Copy configuration file

```
cp internal/configs/config.example.yaml internal/configs/config.yaml
```

3. Adjust configuration in `config.yaml`:

```yaml
database:
  databaseSourceName: "change with your db credentials"

service:
  port: ":9876"
  secretJWT: "your-jwt-secret"

spotifyConfig:
  clientID: "your-spotify-client-id"
  clientSecret: "your-spotify-client-secret"
```

4. Run database migrations

```
go run cmd/main.go
```

5. Run the application

```
go run cmd/main.go
```

## API Endpoints

### Authentication

- `POST /memberships/sign-up` - Register new user
- `POST /memberships/login` - User login
- `POST /memberships/logout` - User logout (requires token)

### Track

- `GET /tracks/search` - Search songs (requires token)
- `POST /tracks/track-activity` - Save track activity (requires token)
- `GET /tracks/recommendations` - Get song recommendations (requires token)

## Testing

To run unit tests:

```
go test ./...
```
