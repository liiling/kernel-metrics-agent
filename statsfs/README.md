## Build Linux Kernel with Statsfs Patch on a VM Running Ubuntu18.04

1. `git clone https://github.com/esposem/linux.git`
2. `git fetch origin statsfs-final` (Fetch the branch with example)
3. `git checkout statsfs-final`
4. Install compilers and other tools

```bash
sudo apt-get install build-essential libncurses-dev bison flex libssl-dev libelf-dev
```

5. Clean the kernel tree: `make mrproper`
6. Generate `.config` file. Some possible alternatives:
	- `make localmodconfig` (generate a config from the kernel options currently in use)
	- `make menuconfig` (command line interface for config creation)
	- `cp -v /boot/config-$(uname -r) .config` (copy the boot config of the host machine)
7. Modify `.config` for statsfs example: set `CONFIG_NET_NS=n`
8. Compile the Linux kernel using all available cpu threads: `make -j $(nproc)`
9. Install the Linux kernel modules: `sudo make modules_install`
10. Install the Linux kernel: `sudo make install`. The following files are installed to the `/boot` directory, and the grub configuration is updated.
    - config-5.7.0-rc2+
    - System.map-5.7.0-rc2+
    - vmlinuz-5.7.0-rc2+
11. Reboot the VM: `reboot`
12. After reboot, check that statsfs is supported in the running Linux kernel:
    - check Linux kernel version: `uname -mrs`
    - check statsfs filesystem is supported: `cat /proc/filesystems | grep statsfs`
13. Mount statfs: `sudo mount -t statsfs statsfs /sys/kernel/stats`
14. Check statsfs is mounted: `cat /proc/mounts | grep stats`
15. Change permission of statsfs filesystem to be readable and executable for everyone, but writable only by the owner (root): `sudo chmod -R 755 /sys/kernel/stats`

## Notes

#### procfs

- A filesystem storing information about processes and some other system information.
- Mounted under */proc*.
- Each running process has a directory.

#### sysfs

- A filesystem storing information about kernel device models, e.g. kernel subsystems, device drivers, etc.
- Mounted under */sys*
- One value per file

#### [debugfs](https://lwn.net/Articles/334546/)

- A RAM-based filesystem for debugging purposes
- No rules at all
- No stable user-space ABI or stability constraints
- Require manual deletion of debugfs files

#### [statsfs](https://lkml.org/lkml/2020/5/26/332)

- A RAM-based filesystem designed to expose kernel statistics to user space
- Problems tacked:
    - Remove code duplication resulted from each linux kernel subsystems having to write codes to gather and display statistics to user space
    - Replaces some statistics currently gathered using debugfs, since debugfs is not meant for metrics
    - A generic and stable API for metrics

## Resources

- [Fork with statsfs implementation](https://github.com/esposem/linux)
- [Userspace implementation of statsfs](https://github.com/esposem/statsfs)
- [Linux kernel mailing list archive entry on statfs](https://lkml.org/lkml/2020/5/26/332)