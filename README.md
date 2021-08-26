# Curl Monitor

## Introduction

Curl Monitor is responsible for running curl commands, store and detect changes in results. It's build as an example project to improve development cycle in microservice environment.

It in active development because I want to try different tools on not so basic projects as "hello-world".

## Preview
```shell
❯ kubectl get pod,service,deployment,hpa,job
NAME                                            READY   STATUS      RESTARTS   AGE
pod/curl-monitor-command-run-5587c5d5f6-5dzrx   1/1     Running     0          86m
pod/curl-monitor-job-856c4c9ff4-szq99           1/1     Running     0          5m52s
pod/curl-monitor-result-588986fd6b-qdsh7        1/1     Running     0          5m47s
pod/curl-monitor-trigger-9588cd7fc-v9k2z        1/1     Running     0          5m50s
pod/migrations-bhqb7                            0/1     Completed   0          178m
pod/postgres-postgresql-0                       1/1     Running     0          12m
pod/rabbitmq-0                                  1/1     Running     8          3d3h

NAME                                   TYPE        CLUSTER-IP       EXTERNAL-IP   PORT(S)                                 AGE
service/curl-monitor-command-run       NodePort    10.107.225.72    <none>        8080:30939/TCP                          147m
service/curl-monitor-job               NodePort    10.100.232.34    <none>        8080:30200/TCP                          178m
service/curl-monitor-result            NodePort    10.111.253.138   <none>        8080:32081/TCP                          169m
service/curl-monitor-trigger           NodePort    10.109.69.207    <none>        8080:31874/TCP                          138m
service/kubernetes                     ClusterIP   10.96.0.1        <none>        443/TCP                                 10d
service/postgres-postgresql            NodePort    10.111.66.49     <none>        5432:30432/TCP                          3d4h
service/postgres-postgresql-headless   ClusterIP   None             <none>        5432/TCP                                3d4h
service/rabbitmq                       ClusterIP   10.106.81.184    <none>        5672/TCP,4369/TCP,25672/TCP,15672/TCP   3d3h
service/rabbitmq-headless              ClusterIP   None             <none>        4369/TCP,5672/TCP,25672/TCP,15672/TCP   3d3h

NAME                                       READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/curl-monitor-command-run   1/1     1            1           147m
deployment.apps/curl-monitor-job           1/1     1            1           178m
deployment.apps/curl-monitor-result        1/1     1            1           169m
deployment.apps/curl-monitor-trigger       1/1     1            1           138m

NAME                                                           REFERENCE                             TARGETS                        MINPODS   MAXPODS   REPLICAS   AGE
horizontalpodautoscaler.autoscaling/curl-monitor-command-run   Deployment/curl-monitor-command-run   <unknown>/80%, <unknown>/80%   1         5         1          147m
horizontalpodautoscaler.autoscaling/curl-monitor-job           Deployment/curl-monitor-job           <unknown>/80%, <unknown>/80%   1         5         1          178m
horizontalpodautoscaler.autoscaling/curl-monitor-result        Deployment/curl-monitor-result        <unknown>/80%, <unknown>/80%   1         5         1          169m
horizontalpodautoscaler.autoscaling/curl-monitor-trigger       Deployment/curl-monitor-trigger       <unknown>/80%, <unknown>/80%   1         5         1          138m

NAME                   COMPLETIONS   DURATION   AGE
job.batch/migrations   1/1           2s         178m
```

## Requirements
1. Go https://golang.org/ - base language
2. Sqlc https://sqlc.dev/ - to generate database structs
3. Docker https://www.docker.com/ 
4. Kubernetes https://kubernetes.io/
5. Helm https://helm.sh/ - to manage packages in Kubernetes


## Deployment
1. Generate database structs `make generate`
2. Build Docker image `make docker-image`
3. Deploy infrastructure `make infra`

## Development
1. Build Docker image `make docker-image`
2. Redeploy needed service, for instance: `make infra-job`

### Missing
1. Skaffold and image auto rebuild (hot reload?)  


## Infrastructure
All resources have configs in `infrastructure` directory. You can easily use `make infra` to setup whole application with its dependencies or you cen got throught each component. 

## Database
To create database we need some configmap, persistant volume, running pods and service to expose to outer world. All of this is in config file `infrastructure/postgress.yaml` 
```shell
make infra-database
```
