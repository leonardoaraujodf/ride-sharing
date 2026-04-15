allow_k8s_contexts('kind-kind')

### K8s Config ###

# Uncomment to use secrets
k8s_yaml('./infra/development/k8s/secrets.yaml')
k8s_yaml('./infra/development/k8s/app-config.yaml')

### End of K8s Config ###

### RabbitMQ ###
k8s_yaml('./infra/development/k8s/rabbitmq-deployment.yaml')
k8s_resource('rabbitmq', port_forwards=['5672', '15672'], labels="tooling")
### End RabbitMQ ###

### API Gateway ###

gateway_compile_cmd = 'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/api-gateway ./services/api-gateway'
if os.name == 'nt':
  gateway_compile_cmd = './infra/development/docker/api-gateway-build.bat'

local_resource(
  'api-gateway-compile',
  gateway_compile_cmd,
  deps=['./services/api-gateway', './shared'], labels="compiles")


custom_build(
  'ride-sharing/api-gateway',
  'podman build -t docker.io/$EXPECTED_REF -f ./infra/development/docker/api-gateway.Dockerfile . && podman save docker.io/$EXPECTED_REF | podman exec -i kind-control-plane ctr -n k8s.io images import -',
  ['./build/api-gateway', './shared'],
  skips_local_docker=True,
  live_update=[
    sync('./build/api-gateway', '/app/build/api-gateway'),
    sync('./shared', '/app/shared'),
    run('kill 1'),
  ],
)

k8s_yaml('./infra/development/k8s/api-gateway-deployment.yaml')
k8s_resource('api-gateway', port_forwards=8081,
             resource_deps=['api-gateway-compile', 'rabbitmq'], labels="services")
### End of API Gateway ###
### Trip Service ###

# Uncomment once we have a trip service

trip_compile_cmd = 'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/trip-service ./services/trip-service/cmd/'
if os.name == 'nt':
  trip_compile_cmd = './infra/development/docker/trip-build.bat'

local_resource(
  'trip-service-compile',
  trip_compile_cmd,
  deps=['./services/trip-service', './shared'], labels="compiles")

custom_build(
  'ride-sharing/trip-service',
  'podman build -t docker.io/$EXPECTED_REF -f ./infra/development/docker/trip-service.Dockerfile . && podman save docker.io/$EXPECTED_REF | podman exec -i kind-control-plane ctr -n k8s.io images import -',
  ['./build/trip-service', './shared'],
  skips_local_docker=True,
  live_update=[
    sync('./build/trip-service', '/app/build/trip-service'),
    sync('./shared', '/app/shared'),
    run('kill 1'),
  ],
)

k8s_yaml('./infra/development/k8s/trip-service-deployment.yaml')
k8s_resource('trip-service', resource_deps=['trip-service-compile', 'rabbitmq'], labels="services")

### End of Trip Service ###

### Driver Service ###
driver_compile_cmd = 'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/driver-service ./services/driver-service/'
if os.name == 'nt':
 driver_compile_cmd = './infra/development/docker/driver-build.bat'

local_resource(
  'driver-service-compile',
  driver_compile_cmd,
  deps=['./services/driver-service', './shared'], labels="compiles")

custom_build(
  'ride-sharing/driver-service',
  'podman build -t docker.io/$EXPECTED_REF -f ./infra/development/docker/driver-service.Dockerfile . && podman save docker.io/$EXPECTED_REF | podman exec -i kind-control-plane ctr -n k8s.io images import -',
  ['./build/driver-service', './shared'],
  skips_local_docker=True,
  live_update=[
    sync('./build/driver-service', '/app/build/driver-service'),
    sync('./shared', '/app/shared'),
    run('kill 1'),
  ],
)

k8s_yaml('./infra/development/k8s/driver-service-deployment.yaml')
k8s_resource('driver-service', resource_deps=['driver-service-compile', 'rabbitmq'], labels="services")

### End of Driver Service ###
### Web Frontend ###

custom_build(
  'ride-sharing/web',
  'podman build -t docker.io/$EXPECTED_REF -f ./infra/development/docker/web.Dockerfile . && podman save docker.io/$EXPECTED_REF | podman exec -i kind-control-plane ctr -n k8s.io images import -',
  ['./web'],
  skips_local_docker=True,
)

k8s_yaml('./infra/development/k8s/web-deployment.yaml')
k8s_resource('web', port_forwards=3000, labels="frontend")

### End of Web Frontend ###

### Payment Service ###

payment_compile_cmd = 'CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/payment-service ./services/payment-service/cmd/main.go'
if os.name == 'nt':
  payment_compile_cmd = './infra/development/docker/payment-build.bat'

local_resource(
  'payment-service-compile',
  payment_compile_cmd,
  deps=['./services/payment-service', './shared'], labels="compiles")

custom_build(
  'ride-sharing/payment-service',
  'podman build -t docker.io/$EXPECTED_REF -f ./infra/development/docker/payment-service.Dockerfile . && podman save docker.io/$EXPECTED_REF | podman exec -i kind-control-plane ctr -n k8s.io images import -',
  ['./build/payment-service', './shared'],
  skips_local_docker=True,
  live_update=[
    sync('./build/payment-service', '/app/build/payment-service'),
    sync('./shared', '/app/shared'),
    run("kill 1")
  ],
)

k8s_yaml('./infra/development/k8s/payment-service-deployment.yaml')
k8s_resource('payment-service', resource_deps=['payment-service-compile', 'rabbitmq'], labels="services")

### End of Payment Service ###