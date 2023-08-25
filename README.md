# Node agent by eBPF

## Run in server

### Install mysql ubuntu

```bash
sudo apt update
sudo apt install mysql-server
sudo systemctl start mysql.service
```

### Change bind address

```bash
vi /etc/mysql/mysql.conf.d/mysqld.cnf
# change bind address 127.0.0.1 -> 0.0.0.0
sudo systemctl restart mysql.service
```

### Create user

```bash
# mysql
CREATE USER 'annd2'@'localhost' IDENTIFIED BY 'password';
GRANT ALL PRIVILEGES ON *.* TO 'annd2'@'localhost' WITH GRANT OPTION;
CREATE USER 'annd2'@'%' IDENTIFIED BY 'password';
GRANT ALL PRIVILEGES ON *.* TO 'annd2'@'%' WITH GRANT OPTION;
FLUSH PRIVILEGES;
```

### Get agent and run

```bash
export OTEL_EXPORTER_OTLP_TRACES_ENDPOINT=http://localhost:4318/v1/traces
```

```bash
curl -L https://github.com/anngdinh/node-agent/releases/download/1.0.0/coroot-node-agent > coroot-node-agent
./coroot-node-agent
```

## Sysbench

```bash
# Write
sysbench oltp_read_only --mysql-host=xxxxxxxxxxxxxxxx --mysql-user=xxxxxxxxxxxxxxxx --mysql-password=xxxxxxxxxxxxxxxx --mysql-port=3306 --mysql-db=xxxxxxxxxxxxxxxx --threads=40 --tables=40 --table-size=2000000 prepare

# Read
Read
sysbench oltp_read_only --mysql-host=xxxxxxxxxxxxxxxx --mysql-user=xxxxxxxxxxxxxxxx --mysql-password=xxxxxxxxxxxxxxxx --mysql-port=3306 --mysql-db=xxxxxxxxxxxxxxxx --threads=40 --tables=40 --table-size=2000000 --range_selects=off --db-ps-mode=disable --report-interval=1 --time=300 run

```

## Note

panic: Something in this program imports go4.org/unsafe/assume-no-moving-gc to declare that it assumes a non-moving garbage collector, but your version of go4.org/unsafe/assume-no-moving-gc hasn't been updated to assert that it's safe against the go1.20 runtime. If you want to risk it, run with environment variable ASSUME_NO_MOVING_GC_UNSAFE_RISK_IT_WITH=go1.20 set. Notably, if go1.20 adds a moving garbage collector, this program is unsafe to use.

```bash
export ASSUME_NO_MOVING_GC_UNSAFE_RISK_IT_WITH=go1.20
```

```bash
export PATH=$PATH:/usr/local/go/bin
```

```bash
# list all trace point event
cd /sys/kernel/debug/tracing/events

# detected a new container" pid=1387 cg="/system.slice/mysql.service" id=/system.slice/mysql.service
```

```bash
export QUERY='SELECT * FROM sbtest1'
export URL='annd2:password@tcp(localhost:3306)/mysql'

# build bpftrace
sudo apt-get update
sudo apt-get install -y \
  bison \
  cmake \
  flex \
  g++ \
  git \
  libelf-dev \
  zlib1g-dev \
  libfl-dev \
  systemtap-sdt-dev \
  binutils-dev \
  libcereal-dev \
  llvm-dev \
  llvm-runtime \
  libclang-dev \
  clang \
  libpcap-dev \
  libgtest-dev \
  libgmock-dev \
  asciidoctor \
  libdw-dev \
  pahole
git clone https://github.com/iovisor/bpftrace --recurse-submodules
mkdir bpftrace/build; cd bpftrace/build;
../build-libs.sh
cmake -DCMAKE_BUILD_TYPE=Release ..
make -j8
sudo make install
```