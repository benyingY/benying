## deploy

Local dev with Kubernetes:
- build image: docker build -t python-service:dev services/python-service-template
- apply manifest: kubectl apply -f deploy/k8s-dev.yaml
- port-forward: kubectl -n benying-dev port-forward svc/python-service 8000:8000
