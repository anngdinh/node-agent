// #include "http.c"
// #include "postgres.c"
// #include "redis.c"
// #include "memcached.c"
// #include "mysql.c"
// #include "mongo.c"
// #include "kafka.c"
// #include "cassandra.c"
// #include "rabbitmq.c"

#include "sys_enter_recvfrom.c"
// #include "sys_enter_sendto.c"

// #define PROTOCOL_UNKNOWN    0
// #define PROTOCOL_HTTP	    1
// #define PROTOCOL_POSTGRES	2
// #define PROTOCOL_REDIS	    3
// #define PROTOCOL_MEMCACHED  4
// #define PROTOCOL_MYSQL      5
// #define PROTOCOL_MONGO      6
// #define PROTOCOL_KAFKA      7
// #define PROTOCOL_CASSANDRA  8
// #define PROTOCOL_RABBITMQ   9

// #define METHOD_UNKNOWN           0
// #define METHOD_PRODUCE           1
// #define METHOD_CONSUME           2
// #define METHOD_STATEMENT_PREPARE 3
// #define METHOD_STATEMENT_CLOSE   4

// #define MAX_PAYLOAD_SIZE 1024

// struct l7_event {
//     __u64 fd;
//     __u64 connection_timestamp;
//     __u32 pid;
//     __u32 status;
//     __u64 duration;
//     __u8 protocol;
//     __u8 method;
//     __u16 padding;
//     __u32 statement_id;
//     char payload[MAX_PAYLOAD_SIZE];
//     __u8 partial;
//     __u64 ns;
//     __u8 request_type;
//     __u64 size;
//     __s32 request_id;
// };

// struct {
//      __uint(type, BPF_MAP_TYPE_PERCPU_ARRAY);
//      __type(key, __u32);
//      __type(value, struct l7_event);
//      __uint(max_entries, 1);
// } l7_event_heap SEC(".maps");

// struct {
//     __uint(type, BPF_MAP_TYPE_PERF_EVENT_ARRAY);
//     __uint(key_size, sizeof(int));
//     __uint(value_size, sizeof(int));
// } l7_events SEC(".maps");

// struct rw_args_t {
//     __u64 fd;
//     char* buf;
//     __u64 size;
// };

// struct {
//     __uint(type, BPF_MAP_TYPE_HASH);
//     __uint(key_size, sizeof(__u64));
//     __uint(value_size, sizeof(struct rw_args_t));
//     __uint(max_entries, 10240);
// } active_reads SEC(".maps");

// struct socket_key {
//     __u64 fd;
//     __u32 pid;
//     __s16 stream_id;
// };

// struct l7_request {
//     __u64 ns;
//     __u8 protocol;
//     __u8 partial;
//     __u8 request_type;
//     __u64 size;
//     __s32 request_id;
//     char payload[MAX_PAYLOAD_SIZE];
// };

// struct {
//      __uint(type, BPF_MAP_TYPE_PERCPU_ARRAY);
//      __type(key, __u32);
//      __type(value, struct l7_request);
//      __uint(max_entries, 1);
// } l7_request_heap SEC(".maps");

// struct {
//     __uint(type, BPF_MAP_TYPE_LRU_HASH);
//     __uint(key_size, sizeof(struct socket_key));
//     __uint(value_size, sizeof(struct l7_request));
//     __uint(max_entries, 32768);
// } active_l7_requests SEC(".maps");

// struct message_event {
//     __u32 messsage; // EVENT_TYPE_PROCESS_START | EVENT_TYPE_PROCESS_EXIT
//     // __u32 pid;
//     // __u32 reason; // 0 | EVENT_REASON_OOM_KILL
// };

// struct {
// 	__uint(type, BPF_MAP_TYPE_PERF_EVENT_ARRAY);
// 	__uint(key_size, sizeof(int));
// 	__uint(value_size, sizeof(int));
// } message_events SEC(".maps");

// struct trace_event_raw_sys_exit_rw__stub_new {
//     __u64 unused;
//     __u64 addr_len;
//     __u64 fd; // 3
//     void * ubuf; // 4
//     __u64 size; // 5
//     __u64 flags; // 6
//     struct sockaddr * addr;

    
//     // int fd
//     // void * ubuf
//     // size_t size
//     // unsigned int flags
//     // struct sockaddr * addr
//     // int * addr_len
// };
// // 4 - 139723741988784

// SEC("tracepoint/syscalls/sys_enter_recvfrom")
// int sys_enter_read(struct trace_event_raw_sys_exit_rw__stub_new* ctx) {
//     __u64 id = bpf_get_current_pid_tgid();
//     // if (id != 3478){
//     //     return 0;
//     // }
    
    

//     int zero = 0;
//     // struct l7_request *req = bpf_map_lookup_elem(&l7_request_heap, &zero);
//     struct l7_event *e = bpf_map_lookup_elem(&l7_event_heap, &zero);
//     if (!e) {
//         return 0;
//     }

//     if (!ctx || !ctx->size) return 0;
    
//     // req->size = 55;
//     e->pid = id;
//     e->size = 5;
//     e->protocol = PROTOCOL_UNKNOWN;
//     // req->partial = 0;
//     // req->request_id = 0;
//     // req->ns = 0;
//     // struct socket_key k = {};
//     // k.pid = id >> 32;
//     // k.fd = 5;
//     // k.stream_id = -1;

//     if (is_mysql_query((char *)ctx->ubuf, ctx->size, &e->request_type)) {
//         // if (req->request_type == MYSQL_COM_STMT_CLOSE) {
//         //     struct l7_event *e = bpf_map_lookup_elem(&l7_event_heap, &zero);
//         //     if (!e) {
//         //         return 0;
//         //     }
//         //     e->protocol = PROTOCOL_MYSQL;
//         //     e->fd = k.fd;
//         //     e->pid = k.pid;
//         //     e->method = METHOD_STATEMENT_CLOSE;
//         //     e->status = 200;
//         //     e->connection_timestamp = get_connection_timestamp(k.pid, k.fd);
//         //     bpf_probe_read(e->payload, MAX_PAYLOAD_SIZE, (void *)ctx->buf);
//         //     bpf_perf_event_output(ctx, &l7_events, BPF_F_CURRENT_CPU, e, sizeof(*e));
//         //     return 0;
//         // }
//         e->protocol = PROTOCOL_MYSQL;
//     } 
//     if (e->protocol == PROTOCOL_UNKNOWN) {
//         return 0;
//     }
//     // if (req->ns == 0) {
//     //     req->ns = bpf_ktime_get_ns();
//     // }
//     bpf_probe_read(e->payload, MAX_PAYLOAD_SIZE, (char *)ctx->ubuf);
//     // bpf_probe_read(req->payload, MAX_PAYLOAD_SIZE, (ctx->addr)[0].sa_data);
//     // bpf_map_update_elem(&active_l7_requests, &k, req, BPF_ANY);
    
//     struct message_event m = {};
//     m.messsage = 97;
//     // bpf_perf_event_output(ctx, &message_events, BPF_F_CURRENT_CPU, &m, sizeof(m)); 
//     // bpf_perf_event_output(ctx, &l7_events, BPF_F_CURRENT_CPU, &m, sizeof(m));
//     bpf_perf_event_output(ctx, &l7_events, BPF_F_CURRENT_CPU, e, sizeof(*e));
//     return 0;
// }