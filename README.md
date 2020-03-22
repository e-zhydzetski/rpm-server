# RPM server

## HTTP API for RPM repository
* Authorized push
* Free pull
* File server without directory indexing vulnerability
* RPM rewrite protection

## API
### Push new package to target repository
```
POST /api/packages
Authorization: Bearer qwe123
Content-Type: multipart/form-data; ...
...
Content-Disposition: form-data; name="package"; filename="test.rpm"
...
```
Sample: `curl -X POST http://host:port/api/packages --form 'package=@test.rpm' -H "Authorization: Bearer qwe123"`
### File server access to repositories for package managers
```
GET /repos
```
Sample: `curl http://host:port/repos/7/updates/x86_64/test.rpm`

## Configuration
### System
* [createrepo](https://linux.die.net/man/8/createrepo) tool should be installed

### Environment variables
* _LISTEN_ADDR_ - listening interface and port in format _interface:port_  
`LISTEN_ADDR=:8080` - listen 8080 port on all interfaces  
`LISTEN_ADDR=127.0.0.1:8888` - listen 8888 port on localhost only (useful with reverse proxy)
* _ACCESS_TOKEN_ - auth token for push access, any string  
`ACCESS_TOKEN=qwe123` - push requests should contains header `Authorization: Bearer qwe123`
* _PUSH_PATH_ - target repo folder path  
`PUSH_PATH=/opt/repos/7/updates/x86_64` - pushed RPMs will be saved into _/opt/repos/7/updates/x86_64_ folder
* _REPOS_ROOT_ - root folder for file server `/repos` path  
`REPOS_ROOT=/opt/repos` - RPM pushed to _/opt/repos/7/updates/x86_64/test.rpm_ will be accessible by `http://host:port/repos/7/updates/x86_64/test.rpm` 

## Docker
* Build: `docker build -f docker/Dockerfile -t local/rpm-server:dev .`
* Run: `docker run --rm -p "8080:8080" -e "ACCESS_TOKEN=123" --name rpm-server local/rpm-server:dev`

## TODO
* Push support for multiple repos inside _REPOS_ROOT_
* Improve access control