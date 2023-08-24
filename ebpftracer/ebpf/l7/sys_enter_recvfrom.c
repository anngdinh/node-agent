#include "http.c"
#include "postgres.c"
#include "redis.c"
#include "memcached.c"
#include "mysql.c"
#include "mongo.c"
#include "kafka.c"
#include "cassandra.c"
#include "rabbitmq.c"

#define PROTOCOL_UNKNOWN    0
#define PROTOCOL_HTTP	    1
#define PROTOCOL_POSTGRES	2
#define PROTOCOL_REDIS	    3
#define PROTOCOL_MEMCACHED  4
#define PROTOCOL_MYSQL      5
#define PROTOCOL_MONGO      6
#define PROTOCOL_KAFKA      7
#define PROTOCOL_CASSANDRA  8
#define PROTOCOL_RABBITMQ   9

#define METHOD_UNKNOWN           0
#define METHOD_PRODUCE           1
#define METHOD_CONSUME           2
#define METHOD_STATEMENT_PREPARE 3
#define METHOD_STATEMENT_CLOSE   4

#define MAX_PAYLOAD_SIZE 1024

#define TASK_COMM_LEN 16
#define SOCK_ADDR_LEN 16
#define LOST_BYTE 4

struct {
    __uint(type, BPF_MAP_TYPE_PERF_EVENT_ARRAY);
    __uint(key_size, sizeof(int));
    __uint(value_size, sizeof(int));
} l7_events SEC(".maps");

struct l7_request {  // this should map byte by byte with l7Event in trace.go
    // __u8 protocol;
    // __u8 request_type;
    __u64 id;
    __u64 pid;
    __u64 size;
    __u64 flags;
    char payload[MAX_PAYLOAD_SIZE];
    char comm[TASK_COMM_LEN];
};

struct {
     __uint(type, BPF_MAP_TYPE_PERCPU_ARRAY);
     __type(key, __u32);
     __type(value, struct l7_request);
     __uint(max_entries, 1);
} l7_request_heap SEC(".maps");


struct trace_event_raw_sys_exit_rw__stub_new {
    __u64 unused;
    __u64 id;
    __u64 fd;
    char* buf;
    __u64 size;
    __u64 flags; // 6
    struct sockaddr * addr;

    // __u64 unused;
    // __u64 addr_len;
    // __u64 fd; // 3
    // char * ubuf; // 4
    // __u64 size; // 5
    // __u64 unused2;
    // __u64 flags; // 6
    // struct sockaddr * addr;

    // int fd
    // void * ubuf
    // size_t size
    // unsigned int flags
    // struct sockaddr * addr
    // int * addr_len
};
// 4 - 139723741988784

struct trace_event_raw_sys_exit_rw__stub {
    __u64 unused;
    long int id;
    long int ret;
};

struct sockaddr {
    __u64 sa_family;
    char sa_data[14];
};

struct read_args {
    __u64 fd;
    __u64 size;
    __u64 flags;
    char* buf;
    __u64* ret;
    char lost_byte[LOST_BYTE];
    struct sockaddr * addr;
};

struct {
    __uint(type, BPF_MAP_TYPE_HASH);
    __uint(key_size, sizeof(__u64));
    __uint(value_size, sizeof(struct read_args));
    __uint(max_entries, 10240);
} active_reads SEC(".maps");

SEC("tracepoint/syscalls/sys_enter_recvfrom")
int sys_enter_recvfrom(struct trace_event_raw_sys_exit_rw__stub_new* ctx) {
    __u64 id = bpf_get_current_pid_tgid();
    char comm[TASK_COMM_LEN];
    bpf_get_current_comm(&comm, sizeof(comm));
    if (comm[0] != 'c' || comm[1] != 'o' || comm[2] != 'n' || comm[3] != 'n' || comm[4] != 'e' || comm[5] != 'c') return 0;
    // if (ctx->size <= 4) return 0;

    struct read_args args = {};
    args.fd = ctx->fd;
    args.size = ctx->size;
    args.flags = ctx->flags;
    args.buf = ctx->buf;
    args.addr = ctx->addr;
    args.ret = 0;
    bpf_probe_read(args.lost_byte, LOST_BYTE, ctx->buf);
    bpf_map_update_elem(&active_reads, &id, &args, BPF_ANY);
    return 0;
}


SEC("tracepoint/syscalls/sys_exit_recvfrom")
int sys_exit_recvfrom(struct trace_event_raw_sys_exit_rw__stub* ctx) {
    __u64 pid_tgid = bpf_get_current_pid_tgid();
    __u32 pid = pid_tgid >> 32;
    // return trace_exit_read(ctx, pid_tgid, pid, 0, ctx->ret);

    struct read_args *args = bpf_map_lookup_elem(&active_reads, &pid_tgid);
    if (!args) {
        return 0;
    }
    bpf_map_delete_elem(&active_reads, &pid_tgid);

    if (args->size <= 4) return 0;
    


    int zero = 0;
    struct l7_request *req = bpf_map_lookup_elem(&l7_request_heap, &zero);
    if (!req) {
        return 0;
    }

    req->size = args->size;
    req->id = pid_tgid;
    req->pid = pid;
    // req->protocol = PROTOCOL_UNKNOWN;
    
    // char comm[TASK_COMM_LEN];
    // bpf_get_current_comm(&comm, sizeof(comm));

    // if (comm[0] != 'c' || comm[1] != 'o' || comm[2] != 'n' || comm[3] != 'n' || comm[4] != 'e' || comm[5] != 'c') return 0;
    bpf_probe_read(req->comm, TASK_COMM_LEN, (void *)(args->addr)->sa_data);

    // if (is_mysql_query((char *)ctx->ubuf, ctx->size, &req->request_type)) {
    //     req->protocol = PROTOCOL_MYSQL;
    // } 
    // if (req->protocol == PROTOCOL_UNKNOWN) {
    //     return 0;
    // }

    // get comm (ex: systemd, ps, connection,...)
    

    // copy payload
    bpf_probe_read(req->payload, LOST_BYTE, args->lost_byte);
    bpf_probe_read(req->payload + 4, MAX_PAYLOAD_SIZE - LOST_BYTE, args->buf);

    // send to maps
    bpf_perf_event_output(ctx, &l7_events, BPF_F_CURRENT_CPU, req, sizeof(*req));
    return 0;
}
