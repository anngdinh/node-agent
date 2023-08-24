# name: sys_enter_sendto
# ID: 1493
# format:
# 	field:unsigned short common_type;	offset:0;	size:2;	signed:0;
# 	field:unsigned char common_flags;	offset:2;	size:1;	signed:0;
# 	field:unsigned char common_preempt_count;	offset:3;	size:1;	signed:0;
# 	field:int common_pid;	offset:4;	size:4;	signed:1;

# 	field:int __syscall_nr;	offset:8;	size:4;	signed:1;
# 	field:int fd;	offset:16;	size:8;	signed:0;
# 	field:void * buff;	offset:24;	size:8;	signed:0;
# 	field:size_t len;	offset:32;	size:8;	signed:0;
# 	field:unsigned int flags;	offset:40;	size:8;	signed:0;
# 	field:struct sockaddr * addr;	offset:48;	size:8;	signed:0;
# 	field:int addr_len;	offset:56;	size:8;	signed:0;
sendto:
	/home/stackops/bpftrace -e 'tracepoint:syscalls:sys_enter_sendto { printf("fd: %d, buff: 0x%08lx, len: %d, flags: %d, addr: 0x%08lx, addr_len: %d, addr.sa_data: %s\n", ((args->fd)), ((args->buff)), ((args->len)), ((args->flags)), ((args->addr)), ((args->addr_len)), ((args->addr->sa_data))); }'

# name: sys_enter_recvfrom
# ID: 1491
# format:
# 	field:unsigned short common_type;	offset:0;	size:2;	signed:0;
# 	field:unsigned char common_flags;	offset:2;	size:1;	signed:0;
# 	field:unsigned char common_preempt_count;	offset:3;	size:1;	signed:0;
# 	field:int common_pid;	offset:4;	size:4;	signed:1;

# 	field:int __syscall_nr;	offset:8;	size:4;	signed:1;
# 	field:int fd;	offset:16;	size:8;	signed:0;
# 	field:void * ubuf;	offset:24;	size:8;	signed:0;
# 	field:size_t size;	offset:32;	size:8;	signed:0;
# 	field:unsigned int flags;	offset:40;	size:8;	signed:0;
# 	field:struct sockaddr * addr;	offset:48;	size:8;	signed:0;
# 	field:int * addr_len;	offset:56;	size:8;	signed:0;

recvfrom:
	/home/stackops/bpftrace -e 'tracepoint:syscalls:sys_enter_recvfrom { printf("fd: %d, ubuf: 0x%08lx, size: %d, flags: %d, addr: 0x%08lx, addr_len: %d\n", ((args->fd)), ((args->ubuf)), ((args->size)), ((args->flags)), ((args->addr)), ((args->addr_len))); }'


arg:
	# /home/stackops/bpftrace -lv tracepoint:syscalls:sys_enter_openat
	# /home/stackops/bpftrace -lv kprobe:tcp_sendmsg
	# /home/stackops/bpftrace -lv tracepoint:syscalls:sys_enter_recvfrom
	/home/stackops/bpftrace -lv tracepoint:syscalls:sys_enter_close
	/home/stackops/bpftrace -lv tracepoint:syscalls:sys_enter_connect
	# /home/stackops/bpftrace -lv tracepoint:syscalls:sys_enter_read
	# /home/stackops/bpftrace -lv "struct path"
	# /home/stackops/bpftrace -lv "struct sockaddr"
	# /home/stackops/bpftrace -lv "struct msghdr"
	# /home/stackops/bpftrace -lv "struct iov_iter"
	# /home/stackops/bpftrace -lv "struct page"
	# /home/stackops/bpftrace -lv "struct iovec"
	# /home/stackops/bpftrace -lv "struct size_t"

	# /home/stackops/bpftrace -e 'tracepoint:syscalls:sys_enter_openat { printf("%s %s\n", comm, str(args->filename)); }'

	# /home/stackops/bpftrace -e 'BEGIN { pid & 1 ? printf("Odd\n") : printf("Even\n"); exit(); }'
	# /home/stackops/bpftrace -e 'BEGIN { printf("hello\n"); }'
	# /home/stackops/bpftrace -e 'tracepoint:syscalls:sys_enter_sendto { printf("Datagram bytes: %r\n", buf(args.buff, args.len)); }'
	# /home/stackops/bpftrace -e 'tracepoint:syscalls:sys_enter_recvfrom { printf("Datagram bytes: %r, %d\n", buf(args.ubuf, args.size), args.size); }'
	
	# bpftrace -l tracepoint:syscalls:sys_enter_*

bpftrace:
	# ./bpftrace ./ebpftracer/ebpf/l7/sys_enter_recvfrom.bt
	# ./bpftrace ./ebpftracer/ebpf/l7/tcp_recvmsg.bt
	./bpftrace test.bt

.PHONY: bpftrace