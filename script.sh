docker image build . -t spa:1.0.0 -f Dockerfile

docker container run --rm --publish 8080:8080/tcp spa:1.0.0