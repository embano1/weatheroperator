## About
This is a demo Kubernetes operator (weather service) using the [operator-sdk](https://github.com/operator-framework/operator-sdk). Based on custom resource events (or every 60s, hardcoded resync interval), this controller will query the current temperature for cities.

```bash
CITY          TEMP    METRIC
Hamburg,DE    8.36    Celsius
Leipzig       9       Celsius
New York,US   36.18   Fahrenheit
```

## Requirements (build and run)

Requirements if you just want to **run** this operator demo in a Kubernetes cluster:

- git (for cloning the repo)
- Kubernetes environment with CRD support (a recent version >=v1.10 is recommended)

Requirements for **building** the operator and if you want to run the operator outside the Kubernetes cluster (testing, dev, etc.):

- [dep](https://golang.github.io/dep/docs/installation.html) version v0.5.0+.
- [git](https://git-scm.com/downloads)
- [go](https://golang.org/dl/) version v1.10+.
- [docker](https://docs.docker.com/install/) version 17.03+.
- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/) version v1.11.0+.
- Access to a kubernetes v.1.11.0+ cluster.
- Then follow the steps to build and run the [operator-sdk](https://github.com/operator-framework/operator-sdk) which was used to build the operator

## Obtain a TOKEN (API key) from OpenWeatherMap

To query the OpenWeatherMap API we need a token. [Sign up](https://home.openweathermap.org/users/sign_up) for the service (free for our purposes as of writing this README) and then [create a custom token (API key)](https://home.openweathermap.org/api_keys) to be used by the operator.

## How to run the operator demo in a Kubernetes cluster

Clone the repo:

```bash
git clone https://github.com/embano1/weatheroperator && cd weatheroperator
```

Prepare the Kubernetes namespace and register the custom resource definition (CRD) for the weather service:

```bash
kubectl create ns weather
kubectl create -f deploy/crds/weatherservice_v1alpha1_weatherservice_crd.yaml
# Verify the CRD was correctly registered before we move on
kubectl get crds
NAME                                        CREATED AT
weatherservices.weatherservice.vmware.com   2018-11-14T09:47:46Z
```

Create a secret containing the API key:

```bash
kubectl -n weather create secret generic appid --from-literal=appid=APIKEY
```

Deploy the operator (incl. RBAC rules):

```bash
kubectl -n weather create -f deploy/
```

Deploy the example locations (custom resources):

```bash
kubectl -n weather create -f deploy/crds/cities/
```

Open another terminal to watch the operator in action:

```bash
watch -n 1 "kubectl -n weather get ws -o=custom-columns=CITY:.spec.city,TEMP:.status.temperature,METRIC:.spec.unit"

CITY          TEMP    METRIC
Hamburg,DE    8.36    Celsius
Leipzig       9       Celsius
New York,US   36.18   Fahrenheit
```

## Add your city

Adding your city to the weather operator is as easy as creating a custom resource:

```yaml
apiVersion: weatherservice.vmware.com/v1alpha1
kind: WeatherService
metadata:
  name: yourcityname
spec:
  # Add fields here
  city: "City,COUNTRYCODE"
  unit: "Fahrenheit" # or "Celsius"
```

Save this as a file (`yourcity.yaml`) and then submit to Kubernetes:

```bash
kubectl -n weather create -f yourcity.yaml
```

## How to run the operator outside the Kubernetes cluster

Make sure you have installed the required software listed above (see "Requirements for building the operator"). Create a file containing the OpenWeatherMap API key. The commands assume you've cloned and changed into the repository as described above.

```bash
# Surpress newline with "-n"
echo -n "APIKEY" > app.id
```

Fetch dependencies:

```bash
dep ensure -vendor-only
```

Prepare the Kubernetes namespace and register the custom resource definition (CRD) for the weather service:

```bash
kubectl create ns weather
kubectl create -f deploy/crds/weatherservice_v1alpha1_weatherservice_crd.yaml
# Verify the CRD was correctly registered before we move on
kubectl get crds
NAME                                        CREATED AT
weatherservices.weatherservice.vmware.com   2018-11-14T09:47:46Z
```

Start the operator in stand-alone mode:

```bash
operator-sdk up local --operator-flags "-c app.id"
2018/11/14 10:51:10 Go Version: go1.11.2
2018/11/14 10:51:10 Go OS/Arch: darwin/amd64
2018/11/14 10:51:10 operator-sdk Version: v0.1.0+git
2018/11/14 10:51:10 Registering Components.
2018/11/14 10:51:10 Starting the Cmd.
```

Deploy the example locations (custom resources):

```bash
kubectl -n weather create -f deploy/crds/cities/
```

Open another terminal to watch the operator in action:

```bash
watch -n 1 "kubectl -n weather get ws -o=custom-columns=CITY:.spec.city,TEMP:.status.temperature,METRIC:.spec.unit"
```