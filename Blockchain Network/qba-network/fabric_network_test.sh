#!/bin/bash

# Hyperledger Fabric Network Health Check Script

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color


log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_separator() {
    echo "=================================================="
}


ORDERER_CONTAINERS=("orderer0.example.com" "orderer1.example.com" "orderer2.example.com" "orderer3.example.com" "orderer4.example.com")
PEER_CONTAINERS=("peer0.org1.example.com")
CLI_CONTAINER="cli"
ALL_CONTAINERS=("${ORDERER_CONTAINERS[@]}" "${PEER_CONTAINERS[@]}" "$CLI_CONTAINER")


declare -A CONTAINER_PORTS
CONTAINER_PORTS["orderer0.example.com"]="6050 6070"
CONTAINER_PORTS["orderer1.example.com"]="6051 6071"
CONTAINER_PORTS["orderer2.example.com"]="6052 6072"
CONTAINER_PORTS["orderer3.example.com"]="6053 6073"
CONTAINER_PORTS["orderer4.example.com"]="6054 6074"
CONTAINER_PORTS["peer0.org1.example.com"]="7051"


NETWORK_NAME="solonet"
NETWORK_SUBNET="172.22.0.0/24"


main() {
    log_info "Starting Hyperledger Fabric network health check..."
    echo "Check time: $(date)"
    print_separator
    
  
    check_docker_service
    

    check_container_status
    
  
    check_container_logs
    

    check_port_usage
    

    check_network_connectivity
    

    check_inter_node_communication
    
 
    check_channel_status
    

    #check_resource_usage
    

    #generate_summary_report
    
    log_info "ok!"
}


check_docker_service() {
    print_separator
    log_info "Checking Docker service..."
    
    if systemctl is-active --quiet docker; then
        log_success "Docker is running"
    else
        log_error "Docker is not running"
        return 1
    fi
    

    if command -v docker-compose &> /dev/null; then
        log_success "Docker Compose is installed: $(docker-compose --version)"
    else
        log_error "Docker Compose not found"
    fi
}


check_container_status() {
    print_separator
    log_info "Checking container status..."
    
    local running_count=0
    local total_count=${#ALL_CONTAINERS[@]}
    
    for container in "${ALL_CONTAINERS[@]}"; do
        if docker ps --format "{{.Names}}" | grep -q "^${container}$"; then
            local status=$(docker inspect --format='{{.State.Status}}' "$container" 2>/dev/null)
            local health=$(docker inspect --format='{{.State.Health.Status}}' "$container" 2>/dev/null)
            
            if [ "$status" = "running" ]; then
                log_success " $container is running"
                if [ "$health" != "" ] && [ "$health" != "null" ]; then
                    log_info "  Health status: $health"
                fi
                ((running_count++))
                
             
                local start_time=$(docker inspect --format='{{.State.StartedAt}}' "$container")
                local ip_addr=$(docker inspect --format='{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' "$container")
                log_info "  Start time: $start_time"
                log_info "  IP address: $ip_addr"
            else
                log_error "Container $container has abnormal status: $status"
            fi
        else
            log_error "Container $container not found or not running"
        fi
    done
    
    log_info "Container status summary: $running_count/$total_count containers running"
}


check_container_logs() {
    print_separator
    log_info "Checking container logs..."
    
    for container in "${ALL_CONTAINERS[@]}"; do
        if docker ps --format "{{.Names}}" | grep -q "^${container}$"; then
            log_info "Checking logs for container $container..."
            
            # Get last 100 lines of logs
            local logs=$(docker logs --tail 100 "$container" 2>&1)
            
            # Check for error keywords
            local error_patterns=("ERROR" "FATAL" "panic" "failed" "error" "exception")
            local warning_patterns=("WARN" "WARNING")
            
            local error_count=0
            local warning_count=0
            
            for pattern in "${error_patterns[@]}"; do
                local count=$(echo "$logs" | grep -i "$pattern" | wc -l)
                error_count=$((error_count + count))
            done
            
            for pattern in "${warning_patterns[@]}"; do
                local count=$(echo "$logs" | grep -i "$pattern" | wc -l)
                warning_count=$((warning_count + count))
            done
            
            if [ $error_count -gt 0 ]; then
                log_error "  Found $error_count error logs"
                log_error "  Recent error messages:"
                echo "$logs" | grep -i -E "(ERROR|FATAL|panic|failed|error|exception)" | tail -5 | while read line; do
                    echo -e "${RED}    $line${NC}"
                done
            else
                log_success "  No error logs found"
            fi
            
            if [ $warning_count -gt 0 ]; then
                log_warning "  Found $warning_count warning logs"
            fi
            

            local restart_count=$(docker inspect --format='{{.RestartCount}}' "$container")
            if [ "$restart_count" -gt 0 ]; then
                log_warning "  Container has restarted $restart_count times"
            fi
        fi
    done
}


check_port_usage() {
    print_separator
    log_info "Checking port usage..."
    

    PURPLE='\033[0;35m'
    CYAN='\033[0;36m'
    
    echo -e "\n${PURPLE}Orderer node ports:${NC}"
    for orderer in "${ORDERER_CONTAINERS[@]}"; do
        echo -e "\n${CYAN}$orderer:${NC}"
        
        if docker ps --format "{{.Names}}" | grep -q "^${orderer}$"; then
            local ports="${CONTAINER_PORTS[$orderer]}"
            for port in $ports; do
    
                if netstat -tuln 2>/dev/null | grep -q ":$port "; then
                    local process=$(netstat -tulnp 2>/dev/null | grep ":$port " | awk '{print $7}' | head -1)
                    log_success "  Port $port: In use ($process)"
                    
        
                    local docker_process=$(ps aux | grep -v grep | grep "docker.*$orderer" | head -1)
                    if [ -n "$docker_process" ]; then
                        log_info "    ✓ Used by Docker container $orderer"
                    fi
                else
                    log_error "  Port $port: Not in use"
                fi
                

                if docker exec "$orderer" netstat -tuln 2>/dev/null | grep -q ":$port "; then
                    log_success "  Internal port $port: Listening"
                else
                    log_warning "  Internal port $port: Not listening"
                fi
                
         
                local connectivity_test=$(timeout 3 bash -c "echo >/dev/tcp/localhost/$port" 2>/dev/null)
                if [ $? -eq 0 ]; then
                    log_success "  Port $port: Connectivity OK"
                else
                    log_warning "  Port $port: Connectivity test failed"
                fi
            done
        else
            log_error "  Container $orderer not running, cannot check ports"
        fi
    done
    
    echo -e "\n${PURPLE}Peer node ports:${NC}"
    for peer in "${PEER_CONTAINERS[@]}"; do
        echo -e "\n${CYAN}$peer:${NC}"
        
        if docker ps --format "{{.Names}}" | grep -q "^${peer}$"; then
            local ports="${CONTAINER_PORTS[$peer]}"
            for port in $ports; do
         
                if netstat -tuln 2>/dev/null | grep -q ":$port "; then
                    local process=$(netstat -tulnp 2>/dev/null | grep ":$port " | awk '{print $7}' | head -1)
                    log_success "  Port $port: In use ($process)"
                    
           
                    local docker_process=$(ps aux | grep -v grep | grep "docker.*$peer" | head -1)
                    if [ -n "$docker_process" ]; then
                        log_info "    ✓ Used by Docker container $peer"
                    fi
                else
                    log_error "  Port $port: Not in use"
                fi
                

                if docker exec "$peer" netstat -tuln 2>/dev/null | grep -q ":$port "; then
                    log_success "  Internal port $port: Listening"
                else
                    log_warning "  Internal port $port: Not listening"
                fi
                
       
                local connectivity_test=$(timeout 3 bash -c "echo >/dev/tcp/localhost/$port" 2>/dev/null)
                if [ $? -eq 0 ]; then
                    log_success "  Port $port: Connectivity OK"
                else
                    log_warning "  Port $port: Connectivity test failed"
                fi
                
         
                if [ "$port" = "7051" ]; then
                    local chaincode_port="7052"
                    echo -e "    ${CYAN}Checking Chaincode port $chaincode_port:${NC}"
                    
                    if docker exec "$peer" netstat -tuln 2>/dev/null | grep -q ":$chaincode_port "; then
                        log_success "    Chaincode port $chaincode_port: Listening"
                    else
                        log_warning "    Chaincode port $chaincode_port: Not listening"
                    fi
                fi
            done
        else
            log_error "  Container $peer not running, cannot check ports"
        fi
    done
    
    # Additional port conflict check
    echo -e "\n${PURPLE}Port conflict check:${NC}"
    local all_fabric_ports=()
    
    # Collect all Fabric ports
    for container in "${!CONTAINER_PORTS[@]}"; do
        for port in ${CONTAINER_PORTS[$container]}; do
            all_fabric_ports+=("$port")
        done
    done
    
    # Check if non-Fabric processes are using these ports
    for port in "${all_fabric_ports[@]}"; do
        local processes=$(netstat -tulnp 2>/dev/null | grep ":$port " | awk '{print $7}' | grep -v "docker" | head -1)
        if [ -n "$processes" ] && [ "$processes" != "-" ]; then
            log_warning "  Port $port used by non-Docker process: $processes"
        fi
    done
    
    # Check system common ports
    echo -e "\n${PURPLE}System port check:${NC}"
    local system_ports=("22" "80" "443" "3000" "8080")
    for port in "${system_ports[@]}"; do
        if netstat -tuln 2>/dev/null | grep -q ":$port "; then
            local process=$(netstat -tulnp 2>/dev/null | grep ":$port " | awk '{print $7}' | head -1)
            log_info "  System port $port: In use ($process)"
        fi
    done
}

# Check network connectivity
check_network_connectivity() {
    print_separator
    log_info "Checking network connectivity..."
    
    # Check Docker network
    if docker network ls | grep -q "$NETWORK_NAME"; then
        log_success "Docker network $NETWORK_NAME exists"
        
        # Check network details
        local network_info=$(docker network inspect "$NETWORK_NAME")
        log_info "Network subnet: $(echo "$network_info" | jq -r '.[0].IPAM.Config[0].Subnet' 2>/dev/null || echo "Unable to get")"
        
        # Check containers in network
        local containers_in_network=$(docker network inspect "$NETWORK_NAME" | jq -r '.[0].Containers | keys[]' 2>/dev/null)
        if [ -n "$containers_in_network" ]; then
            log_info "Containers in network:"
            echo "$containers_in_network" | while read container_id; do
                local container_name=$(docker inspect --format='{{.Name}}' "$container_id" | sed 's/\///')
                local container_ip=$(docker network inspect "$NETWORK_NAME" | jq -r ".[0].Containers[\"$container_id\"].IPv4Address" | cut -d'/' -f1)
                log_info "  $container_name: $container_ip"
            done
        fi
    else
        log_error "Docker network $NETWORK_NAME does not exist"
    fi
    
    # Check inter-container connectivity
    log_info "Testing inter-container connectivity..."
    
    for container in "${ALL_CONTAINERS[@]}"; do
        if docker ps --format "{{.Names}}" | grep -q "^${container}$"; then
            # Test connectivity to other containers
            for target_container in "${ALL_CONTAINERS[@]}"; do
                if [ "$container" != "$target_container" ] && docker ps --format "{{.Names}}" | grep -q "^${target_container}$"; then
                    local ping_result=$(docker exec "$container" ping -c 1 -W 3 "$target_container" 2>/dev/null)
                    if [ $? -eq 0 ]; then
                        log_success "  $container -> $target_container: Connected"
                    else
                        log_error "  $container -> $target_container: Not connected"
                    fi
                fi
            done
        fi
    done
}

# Check inter-node communication
check_inter_node_communication() {
    print_separator
    log_info "Checking inter-node communication..."
    
    # Check PBFT communication between Orderer nodes
    log_info "Checking Orderer node PBFT communication..."
    for orderer in "${ORDERER_CONTAINERS[@]}"; do
        if docker ps --format "{{.Names}}" | grep -q "^${orderer}$"; then
            # Check PBFT port connectivity
            local pbft_port=""
            case $orderer in
                "orderer0.example.com") pbft_port="6070" ;;
                "orderer1.example.com") pbft_port="6071" ;;
                "orderer2.example.com") pbft_port="6072" ;;
                "orderer3.example.com") pbft_port="6073" ;;
                "orderer4.example.com") pbft_port="6074" ;;
            esac
            
            log_info "  Checking $orderer PBFT port $pbft_port..."
            
            # Test PBFT connection to other orderers
            for target_orderer in "${ORDERER_CONTAINERS[@]}"; do
                if [ "$orderer" != "$target_orderer" ] && docker ps --format "{{.Names}}" | grep -q "^${target_orderer}$"; then
                    local target_pbft_port=""
                    case $target_orderer in
                        "orderer0.example.com") target_pbft_port="6070" ;;
                        "orderer1.example.com") target_pbft_port="6071" ;;
                        "orderer2.example.com") target_pbft_port="6072" ;;
                        "orderer3.example.com") target_pbft_port="6073" ;;
                        "orderer4.example.com") target_pbft_port="6074" ;;
                    esac
                    
                    local telnet_result=$(docker exec "$orderer" timeout 3 bash -c "echo > /dev/tcp/$target_orderer/$target_pbft_port" 2>/dev/null)
                    if [ $? -eq 0 ]; then
                        log_success "    $orderer -> $target_orderer:$target_pbft_port: Connected"
                    else
                        log_error "    $orderer -> $target_orderer:$target_pbft_port: Not connected"
                    fi
                fi
            done
        fi
    done
    
    # Check Peer to Orderer communication
    log_info "Checking Peer to Orderer communication..."
    for peer in "${PEER_CONTAINERS[@]}"; do
        if docker ps --format "{{.Names}}" | grep -q "^${peer}$"; then
            for orderer in "${ORDERER_CONTAINERS[@]}"; do
                if docker ps --format "{{.Names}}" | grep -q "^${orderer}$"; then
                    local orderer_port=""
                    case $orderer in
                        "orderer0.example.com") orderer_port="6050" ;;
                        "orderer1.example.com") orderer_port="6051" ;;
                        "orderer2.example.com") orderer_port="6052" ;;
                        "orderer3.example.com") orderer_port="6053" ;;
                        "orderer4.example.com") orderer_port="6054" ;;
                    esac
                    
                    local telnet_result=$(docker exec "$peer" timeout 3 bash -c "echo > /dev/tcp/$orderer/$orderer_port" 2>/dev/null)
                    if [ $? -eq 0 ]; then
                        log_success "  $peer -> $orderer:$orderer_port: Connected"
                    else
                        log_error "  $peer -> $orderer:$orderer_port: Not connected"
                    fi
                fi
            done
        fi
    done
}

# Check channel status
check_channel_status() {
    print_separator
    log_info "Checking channel status..."
    
    if docker ps --format "{{.Names}}" | grep -q "^${CLI_CONTAINER}$"; then
        # Set environment variables (using correct docker exec -e syntax)
        local peer_env="CORE_PEER_LOCALMSPID=Org1MSP"
        local tls_env="CORE_PEER_TLS_ENABLED=false"
        local addr_env="CORE_PEER_ADDRESS=peer0.org1.example.com:7051"
        local msp_env="CORE_PEER_MSPCONFIGPATH=/opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp"
        
        # Check if peer can list channels
        log_info "Attempting to list joined channels..."
        local channel_list_output=$(docker exec \
            -e "$peer_env" \
            -e "$tls_env" \
            -e "$addr_env" \
            -e "$msp_env" \
            "$CLI_CONTAINER" peer channel list 2>&1)
        
        if [ $? -eq 0 ] && echo "$channel_list_output" | grep -q "Channels peers has joined"; then
            log_success "Successfully retrieved channel list"
            
            # Extract channel names (skip header line)
            local channels=$(echo "$channel_list_output" | sed -n '/Channels peers has joined/,$p' | tail -n +2 | grep -v '^$')
            
            if [ -n "$channels" ]; then
                echo "$channels" | while IFS= read -r channel; do
                    # Clean channel name (remove spaces and special chars)
                    channel=$(echo "$channel" | tr -d ' \t\r\n')
                    
                    if [ -n "$channel" ]; then
                        log_info "  Joined channel: $channel"
                        
                        # Check channel details
                        log_info "  Checking details for channel $channel..."
                        local channel_info=$(docker exec \
                            -e "$peer_env" \
                            -e "$tls_env" \
                            -e "$addr_env" \
                            -e "$msp_env" \
                            "$CLI_CONTAINER" peer channel getinfo -c "$channel" 2>&1)
                        
                        if [ $? -eq 0 ] && echo "$channel_info" | grep -q "Blockchain info"; then
                            # Extract block height
                            local height=$(echo "$channel_info" | grep -o '"height":[0-9]*' | cut -d: -f2)
                            if [ -z "$height" ]; then
                                # Try alternative format
                                height=$(echo "$channel_info" | grep -o "height:[0-9]*" | cut -d: -f2)
                            fi
                            
                            if [ -n "$height" ]; then
                                log_success "    Channel height: $height"
                            else
                                log_success "    Channel info retrieved but unable to parse height"
                            fi
                            
                            # Extract other useful info if available
                            local current_hash=$(echo "$channel_info" | grep -o '"currentBlockHash":"[^"]*"' | cut -d'"' -f4)
                            if [ -n "$current_hash" ]; then
                                log_info "    Current block hash: ${current_hash:0:16}..."
                            fi
                        else
                            log_error "    Unable to get info for channel $channel"
                            log_info "    Error details: $channel_info"
                        fi
                    fi
                done
            else
                log_warning "Peer is connected but has not joined any channels"
            fi
        else
            log_warning "Unable to get channel list"
            log_info "Command output: $channel_list_output"
            
            # If command fails, check peer connection status
            log_info "Checking peer connection status..."
            local peer_status=$(docker exec \
                -e "$peer_env" \
                -e "$tls_env" \
                -e "$addr_env" \
                -e "$msp_env" \
                "$CLI_CONTAINER" peer node status 2>&1)
            
            if [ $? -eq 0 ]; then
                log_info "Peer node status: Running normally"
            else
                log_error "Peer node status check failed: $peer_status"
            fi
        fi
        
    else
        log_error "CLI container not running, cannot check channel status"
        log_info "Please ensure CLI container is started and running normally"
    fi
}

# Check resource usage
check_resource_usage() {
    print_separator
    log_info "Checking container resource usage..."
    
    # Check container CPU and memory usage
    log_info "Container resource usage statistics:"
    printf "%-25s %-10s %-10s %-10s %-10s\n" "Container Name" "CPU%" "Mem Usage" "Mem Limit" "Net I/O"
    echo "--------------------------------------------------------------------------------"
    
    for container in "${ALL_CONTAINERS[@]}"; do
        if docker ps --format "{{.Names}}" | grep -q "^${container}$"; then
            local stats=$(docker stats --no-stream --format "table {{.Container}}\t{{.CPUPerc}}\t{{.MemUsage}}\t{{.NetIO}}" "$container" | tail -n 1)
            printf "%-25s %s\n" "$container" "$stats"
        fi
    done
    
    # Check disk usage
    log_info "Checking disk usage..."
    local total_size=$(docker system df | grep "Images" | awk '{print $3}')
    local containers_size=$(docker system df | grep "Containers" | awk '{print $3}')
    local volumes_size=$(docker system df | grep "Local Volumes" | awk '{print $3}')
    
    log_info "Docker disk usage:"
    log_info "  Images: $total_size"
    log_info "  Containers: $containers_size"
    log_info "  Volumes: $volumes_size"
}

# Generate summary report
generate_summary_report() {
    print_separator
    log_info "Generating health check summary report..."
    
    local report_file="fabric_health_report_$(date +%Y%m%d_%H%M%S).txt"
    
    {
        echo "Hyperledger Fabric Network Health Check Report"
        echo "================================================"
        echo "Check time: $(date)"
        echo "Network configuration: 5 Orderer nodes (PBFT consensus), 1 Peer node, 1 CLI container"
        echo ""
        
        echo "Container status summary:"
        for container in "${ALL_CONTAINERS[@]}"; do
            if docker ps --format "{{.Names}}" | grep -q "^${container}$"; then
                echo "  ✓ $container: Running"
            else
                echo "  ✗ $container: Not running"
            fi
        done
        echo ""
        
        echo "Port listening status:"
        for container in "${!CONTAINER_PORTS[@]}"; do
            if docker ps --format "{{.Names}}" | grep -q "^${container}$"; then
                for port in ${CONTAINER_PORTS[$container]}; do
                    if netstat -tuln | grep -q ":${port} "; then
                        echo "  ✓ $container:$port - Listening"
                    else
                        echo "  ✗ $container:$port - Not listening"
                    fi
                done
            fi
        done
        echo ""
        
        echo "Network connectivity test:"
        echo "  Detailed connectivity test results please see above output"
        echo ""
        
        echo "Recommendations:"
        echo "1. Regularly check container logs for errors"
        echo "2. Monitor container resource usage"
        echo "3. Ensure all inter-node network connectivity is normal"
        echo "4. Regularly backup important data and configuration files"
        echo ""
        
    } > "$report_file"
    
    log_success "Health check report saved to: $report_file"
}

# Script usage help
show_help() {
    echo "Hyperledger Fabric Network Health Check Script"
    echo ""
    echo "Usage: $0 [options]"
    echo ""
    echo "Options:"
    echo "  -h, --help      Show this help message"
    echo "  -v, --verbose   Verbose output mode"
    echo "  -q, --quiet     Silent mode"
    echo ""
    echo "Examples:"
    echo "  $0              # Run full health check"
    echo "  $0 -v           # Run in verbose mode"
    echo "  $0 -q           # Run in silent mode"
}

# Parameter handling
VERBOSE=false
QUIET=false

while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            show_help
            exit 0
            ;;
        -v|--verbose)
            VERBOSE=true
            shift
            ;;
        -q|--quiet)
            QUIET=true
            shift
            ;;
        *)
            log_error "Unknown parameter: $1"
            show_help
            exit 1
            ;;
    esac
done

# Silent mode handling
if [ "$QUIET" = true ]; then
    exec > /dev/null 2>&1
fi

# Check required tools
check_dependencies() {
    local tools=("docker" "docker-compose" "netstat" "jq")
    for tool in "${tools[@]}"; do
        if ! command -v "$tool" &> /dev/null; then
            log_error "Required tool $tool not found, please install and try again"
            exit 1
        fi
    done
}

# Run script
check_dependencies
main
