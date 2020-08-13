export HOST_IP='192.168.67.11'
export RPT_VER='0.0.16'

docker run -it \
    --name rpt \
    -p 5000:5000 \
    -e RPT_PRIMARY_HOST=postgres:$HOST_IP \
    -e RPT_PRIMARY_PORT=5432 \
    -e RPT_PRIMARY_USER=postgres \
    -e RPT_PRIMARY_PASS=mysecretpassword \
    -e RPT_PRIMARY_SSLMODE=disable \
    -e RPT_SECONDARY_HOST=postgres:$HOST_IP \
    -e RPT_SECONDARY_PORT=5432 \
    -e RPT_SECONDARY_USER=postgres \
    -e RPT_SECONDARY_PASS=mysecretpassword \
    -e RPT_SECONDARY_SSLMODE=disable \
    -e RPT_API=TRUE \
    -e RPT_API_BASEPATH='/api' \
    -e RPT_API_LISTEN_ADDR=':5000' \
    chrishaylesnortal/rpt:$RPT_VER