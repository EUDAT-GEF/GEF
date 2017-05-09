# Create VMs for the Swarm
for N in 1 2 3; do docker-machine create --driver virtualbox swarm$N; done

# SSH to the VM
docker-machine ssh swarm1

# Initialize a swarm
docker swarm init --advertise-addr 192.168.99.106

# Retrieve a manager token
docker swarm join-token manager

# Retrieve a worker token
docker swarm join-token worker

# Running multiple instances of a container (container is launched and detached)
e=80;for N in $(seq 1 1000); do docker run -d -p $((N+e)):80 nginx; done