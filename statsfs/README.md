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