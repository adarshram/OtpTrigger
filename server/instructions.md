cat service-ac.json | docker login -u \_json_key --password-stdin \
https://us-central1-docker.pkg.dev

docker buildx build --platform linux/amd64 -t us-central1-docker.pkg.dev/otptrigger-69fe8/url-trigger/test .
docker push us-central1-docker.pkg.dev/otptrigger-69fe8/url-trigger/test
