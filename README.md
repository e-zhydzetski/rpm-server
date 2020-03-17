# RPM server

## Docker
1. `docker build -f docker/Dockerfile -t local/rpm-server:dev .`
2. `docker run --rm -p "8080:8080" -e "ACCESS_TOKEN=123" --name rpm-server local/rpm-server:dev`