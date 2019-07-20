# go_prometheus_setup
A sample app to demonstrate instrumentation with go and prometheus


### Using Docker Swarm mode and prometheus go client
https://github.com/prometheus/client_golang

Google SRE Best practice
https://landing.google.com/sre/sre-book/chapters/monitoring-distributed-systems/


## Run outside the cluster
`docker-compose up`
- Prometheus runs on `localhost:9090`
- Graphana runs on `localhost:3000` { user: admin, password: pass}
- cadvisor runs on `localhost:8080` 

### build image
`docker build -t fraugster_server .`

### Log in this CLI session using your Docker credentials
`docker login`

### Tag <image> for upload to registry
`docker tag fraugster_server username/fraugsterserver:latest`

### Upload tagged image to registry
`docker push username/fraugsterserver:latest`

### image 
`samuelvarejames/fraugster_server:latest`

### start your single node cluster using docker swarm mode
`docker swarm init`

### deploy app into node
`docker stack deploy -c docker-compose.yml obsvchallenge`

### get the service ID for the one service in our application (app within the node)
`docker service ls`

```
bash
ID                  NAME                       MODE                REPLICAS            IMAGE                                     PORTS
u0bysconq9ge        obsvchallenge_cadvisor     replicated          1/1                 google/cadvisor:latest                    *:8080->8080/tcp
ilu3i6cic8e3        obsvchallenge_db           replicated          1/1                 postgres:latest                           
im9h6eprkuia        obsvchallenge_grafana      replicated          1/1                 grafana/grafana:3.0.0-beta7               *:3000->3000/tcp
1w3eexsj2vop        obsvchallenge_prometheus   replicated          1/1                 prom/prometheus:latest                    *:9090->9090/tcp
wputadz4elbg        obsvchallenge_server       replicated          20/20               samuelvarejames/fraugster_server:latest   *:8000->8000/tcp
```

### list the replicas deployed for the server
`docker service ps obsvchallenge_server`

### Tear down an application obsvchallenge
`docker stack rm obsvchallenge`

### Take down a single node swarm from the manager
`docker swarm leave --force`



-----------------------------------------------------
### Using K8s and  go prometheus
-----------------------------------------------------


### start minikube
`minikube start`

### check cluster info
`kubectl cluster-info` 

### create a ServiceAccount and associate it with the ClusterRole, use a ClusterRoleBinding
```shell
kubectl create serviceaccount -n kube-system tiller
kubectl create clusterrolebinding tiller-cluster-rule --clusterrole=cluster-admin --serviceaccount=kube-system:tiller
```

### Initialize Helm 
`helm init --history-max 200 --service-account tiller`

### confirm tiller works
`kubectl --namespace kube-system get pods | grep tiller`

### bind the prometheus service account
`kubectl apply -f crbinding.yaml`

### Install Prometheus and Postgres into the cluster, using a Helm chart
```shell
helm install --name prometheus-release --namespace monitoring stable/prometheus -f values.yaml
helm install --name postgres-release stable/postgresql
```

### Install graphana
```shell
helm install --name prometheus-graphana --namespace monitoring stable/graphana -f values.yaml
```

### Get the Prometheus server URL by running these commands in the same shell:
```shell
export POD_NAME=$(kubectl get pods --namespace monitoring -l "app=prometheus,component=server" -o jsonpath="{.items[0].metadata.name}")
kubectl --namespace monitoring port-forward $POD_NAME 9090
```

### install the server helm charts to the cluster
`helm install --name server-release docker-compose.kompose`

### view the kubernetes dashboard
`minikube dashboard` 

### delete the release(s)
```shell
helm delete server-release
helm delete prometheus-release
helm delete postgres-release
```

### purge
```shell
helm del --purge server-release
helm del --purge prometheus-release
helm del --purge postgres-release
```

### stop the local k8s cluster
`minikube stop`

### delete local k8s cluster
`minikube delete`
