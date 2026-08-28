# HTTPEagle

A lightweight HTTP service for serving Eagle images. Built with Go.

## API

### Get Item

```
GET /api/items?id=<item_id>&thumbnail=true
```

**Parameters:**
- `id` (required): Item identifier
- `thumbnail` (optional): specify `true` to require thumbnail

**Examples:**
```bash
curl "http://localhost:41596/api/items?id=MTCQADG5S1BFL"
```

## Deployment with Docker

### Docker Run

**1. Build the image:**
```bash
docker build -t httpeagle:latest .
```

**2. Run the container:**
```bash
docker run -d \
  --name httpeagle-server \
  -p 41596:41596 \
  -v /mnt/c/eagle/personal.library/images:/items:ro \
  httpeagle:latest
```

### Compose
```yaml
services:
  httpeagle:
    build: .
    container_name: httpeagle-server
    ports:
      - "41596:41596"
    volumes:
      - /mnt/c/eagle/personal.library/images:/items:ro
      # - /mnt/c/certs/localhost/localhost+1.pem:/certs/cert.pem:ro
      # - /mnt/c/certs/localhost/localhost+1-key.pem:/certs/key.pem:ro
    restart: unless-stopped
```

### HTTPS Setup

Map your certificate files to `/certs/`:
- `/certs/`
    - `cert.pem`
    - `key.pem`
