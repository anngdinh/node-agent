#include "http.c"
#include "postgres.c"
#include "redis.c"
#include "memcached.c"
#include "mysql.c"
#include "mongo.c"
#include "kafka.c"
#include "cassandra.c"
#include "rabbitmq.c"

#define PROTOCOL_UNKNOWN 0
#define PROTOCOL_HTTP 1
#define PROTOCOL_POSTGRES 2
#define PROTOCOL_REDIS 3
#define PROTOCOL_MEMCACHED 4
#define PROTOCOL_MYSQL 5
#define PROTOCOL_MONGO 6
#define PROTOCOL_KAFKA 7
#define PROTOCOL_CASSANDRA 8
#define PROTOCOL_RABBITMQ 9

#define METHOD_UNKNOWN 0
#define METHOD_PRODUCE 1
#define METHOD_CONSUME 2
#define METHOD_STATEMENT_PREPARE 3
#define METHOD_STATEMENT_CLOSE 4

#define TASK_COMM_LEN 16
#define SOCK_ADDR_LEN 16
#define LOST_BYTE 4
#define MAX_PAYLOAD_SIZE 1024

struct l7_event
{
    __u64 fd;
    __u64 connection_timestamp;
    __u32 pid;
    __u32 status;
    __u64 duration;
    __u8 protocol;
    __u8 method;
    __u16 padding;
    __u32 statement_id;
    __u64 id;
    __u64 size;
    char payload[MAX_PAYLOAD_SIZE];
};

struct
{
    __uint(type, BPF_MAP_TYPE_PERCPU_ARRAY);
    __type(key, __u32);
    __type(value, struct l7_event);
    __uint(max_entries, 1);
} l7_event_heap SEC(".maps");

struct
{
    __uint(type, BPF_MAP_TYPE_PERF_EVENT_ARRAY);
    __uint(key_size, sizeof(int));
    __uint(value_size, sizeof(int));
} l7_events SEC(".maps");

struct read_args
{
    __u64 fd;
    __u64 size;
    char *buf;
    __u64 *ret;
    __u64 ns;
    char payload[LOST_BYTE];
};

struct
{
    __uint(type, BPF_MAP_TYPE_HASH);
    __uint(key_size, sizeof(__u64));
    __uint(value_size, sizeof(struct read_args));
    __uint(max_entries, 10240);
} active_reads SEC(".maps");

struct l7_request
{
    __u64 size;
    __u64 fd;
    __u64 ns;
    __u8 protocol;
    __u8 partial;
    __u8 request_type;
    __s32 request_id;
    char payload[MAX_PAYLOAD_SIZE];
};

struct
{
    __uint(type, BPF_MAP_TYPE_PERCPU_ARRAY);
    __type(key, __u32);
    __type(value, struct l7_request);
    __uint(max_entries, 1);
} l7_request_heap SEC(".maps");

struct
{
    __uint(type, BPF_MAP_TYPE_LRU_HASH);
    // __uint(key_size, sizeof(struct socket_key));
    __uint(key_size, sizeof(__u64));
    __uint(value_size, sizeof(struct l7_request));
    __uint(max_entries, 32768);
} active_l7_requests SEC(".maps");

struct trace_event_raw_sys_enter_rw__stub
{
    __u64 unused;
    long int id;
    __u64 fd;
    char *buf;
    __u64 size;
};

struct trace_event_raw_sys_exit_rw__stub
{
    __u64 unused;
    long int id;
    long int ret;
};

SEC("tracepoint/syscalls/sys_enter_recvfrom")
int sys_enter_recvfrom(struct trace_event_raw_sys_enter_rw__stub *ctx)
{
    __u64 id = bpf_get_current_pid_tgid();
    char comm[TASK_COMM_LEN];
    bpf_get_current_comm(&comm, sizeof(comm));
    if (comm[0] != 'c' || comm[1] != 'o' || comm[2] != 'n' || comm[3] != 'n' || comm[4] != 'e' || comm[5] != 'c')
        return 0;
    if (ctx->size <= 4)
        return 0;

    struct read_args args = {};
    args.fd = ctx->fd;
    args.size = ctx->size;
    args.buf = ctx->buf;
    args.ret = 0;
    args.ns = bpf_ktime_get_ns();
    bpf_probe_read(args.payload, LOST_BYTE, ctx->buf);
    bpf_map_update_elem(&active_reads, &id, &args, BPF_ANY);
    return 0;
}

SEC("tracepoint/syscalls/sys_exit_recvfrom")
int sys_exit_recvfrom(struct trace_event_raw_sys_exit_rw__stub *ctx)
{
    __u64 pid_tgid = bpf_get_current_pid_tgid();
    __u32 pid = pid_tgid >> 32;

    struct read_args *args = bpf_map_lookup_elem(&active_reads, &pid_tgid);
    if (!args)
    {
        return 0;
    }
    bpf_map_delete_elem(&active_reads, &pid_tgid);

    int zero = 0;
    struct l7_request *req = bpf_map_lookup_elem(&l7_request_heap, &zero);
    if (!req)
    {
        return 0;
    }

    req->protocol = PROTOCOL_UNKNOWN;
    req->size = args->size + LOST_BYTE;
    req->fd = args->fd;
    req->ns = args->ns;
    bpf_probe_read(req->payload, LOST_BYTE, args->payload);
    bpf_probe_read(req->payload + LOST_BYTE, MAX_PAYLOAD_SIZE - LOST_BYTE, (void *)args->buf);

    if (is_mysql_query(req->payload, req->size, &req->request_type))
    {
        req->protocol = PROTOCOL_MYSQL;
    }

    if (req->protocol == PROTOCOL_UNKNOWN)
    {
        return 0;
    }
    bpf_map_update_elem(&active_l7_requests, &pid_tgid, req, BPF_ANY);
    return 0;
}

/* --------------------------------------------- */
SEC("tracepoint/syscalls/sys_enter_sendto")
int sys_enter_sendto(struct trace_event_raw_sys_enter_rw__stub *ctx)
{
    __u64 pid_tgid = bpf_get_current_pid_tgid();
    __u32 pid = pid_tgid >> 32;
    char comm[TASK_COMM_LEN];
    bpf_get_current_comm(&comm, sizeof(comm));
    if (comm[0] != 'c' || comm[1] != 'o' || comm[2] != 'n' || comm[3] != 'n' || comm[4] != 'e' || comm[5] != 'c')
        return 0;

    struct l7_request *req = bpf_map_lookup_elem(&active_l7_requests, &pid_tgid);
    if (!req)
    {
        return 0;
    }

    int zero = 0;
    struct l7_event *evt = bpf_map_lookup_elem(&l7_event_heap, &zero);
    if (!evt)
    {
        return 0;
    }

    __u64 ns = req->ns;
    evt->size = req->size;
    evt->id = pid_tgid;
    evt->pid = pid;
    evt->fd = req->fd;
    evt->protocol = req->protocol;
    // evt->request_type = req->request_type;
    evt->duration = bpf_ktime_get_ns() - ns;

    // copy payload
    // bpf_probe_read(evt->payload, MAX_PAYLOAD_SIZE, &buf);
    bpf_probe_read(evt->payload, MAX_PAYLOAD_SIZE, req->payload);

    // send to maps
    bpf_perf_event_output(ctx, &l7_events, BPF_F_CURRENT_CPU, evt, sizeof(*evt));
    return 0;
}
