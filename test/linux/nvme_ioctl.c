#include <stdio.h>

#include "nvme_ioctl.h"

#define PRINT_IOCTL_CMD_CODE(CODE) \
    printf("%-24s 0x%lx\n", #CODE":", (long unsigned int)CODE);

int main(int argc, char* argv[])
{
    printf("sizeof(unsigned long int):   %lu\n", sizeof(unsigned long int));
    printf("sizeof(nvme_user_io):        %lu\n", sizeof(struct nvme_user_io));
    printf("sizeof(nvme_passthur_cmd):   %lu\n", sizeof(struct nvme_passthru_cmd));
    printf("sizeof(nvme_passthur_cmd64): %lu\n", sizeof(struct nvme_passthru_cmd64));
    printf("sizeof(nvme_admin_cmd):      %lu\n", sizeof(struct nvme_admin_cmd));

    PRINT_IOCTL_CMD_CODE(NVME_IOCTL_ID)
    PRINT_IOCTL_CMD_CODE(NVME_IOCTL_ADMIN_CMD)
    PRINT_IOCTL_CMD_CODE(NVME_IOCTL_SUBMIT_IO)
    PRINT_IOCTL_CMD_CODE(NVME_IOCTL_IO_CMD)
    PRINT_IOCTL_CMD_CODE(NVME_IOCTL_RESET)
    PRINT_IOCTL_CMD_CODE(NVME_IOCTL_SUBSYS_RESET)
    PRINT_IOCTL_CMD_CODE(NVME_IOCTL_RESCAN)
    PRINT_IOCTL_CMD_CODE(NVME_IOCTL_ADMIN64_CMD)
    PRINT_IOCTL_CMD_CODE(NVME_IOCTL_IO64_CMD)

    return 0;
}