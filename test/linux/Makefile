all: test

ioctl_checker: nvme_ioctl.c nvme_ioctl.h
	@gcc -o ioctl_checker nvme_ioctl.c

test: ioctl_checker
	@./ioctl_checker

clean:
	@rm -f ioctl_checker