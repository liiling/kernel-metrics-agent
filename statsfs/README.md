## Table of Contents

1. [Create VM Instances in GCP with Nested Virtualization Support](#create-vm-instances-in-gcp-with-nested-virtualization-support)
1. [Build Linux Kernel with Statsfs Patch on a Ubuntu18.04 VM](#build-linux-kernel-with-statsfs-patch-on-ubuntu18.04vm)
1. [Install the user-space agent](#install-the-user-space-agent)
1. [Instrument Statsfs with OpenTelemetry](#instrument-statsfs-with-opentelemetry)
1. [Notes](#notes)
1. [Resources](#resources)

## Create VM Instances in GCP with Nested Virtualization Support

The statsfs implementation is heavily inspired by the KVM code that exposes statistics to debugfs, and one of the main usage examples given by statsfs creators exposes KVM statistics. As such, it is useful to create VMs with KVM enabled for testing purposes.
This project runs on GCP, [enabling nested virtualization](https://cloud.google.com/compute/docs/instances/enable-nested-virtualization-vm-instances#tested_os_versions) is requried to start KVM.

1. Create a boot disk from a public or custom image with an operating system.

    a. Via gcloud: `gcloud compute disks create kvm-disk --image-project debian-cloud --image-family debian-9 --zone europe-west1-c --size=200g`

    b. Via GCP's _Disks_ web interface under _Compute Engine_

1. Create a custom image with nested virtualization enabled:
    ```
    gcloud compute images create nested-vm-image \
    --source-disk kvm-disk \
    --source-disk-zone europe-west1-c \
    --licenses https://compute.googleapis.com/compute/v1/projects/vm-options/global/licenses/enable-vmx
    ```

1. Delete the source disk if it is no longer needed:
    ```
    gcloud compute disks delete kvm-disk --zone europe-west1-c
    ```
    
1. Create a VM instance using the new custom image - note how we're using 16 (!) vCPUs, for a very fast Kernel build:

    ```
    gcloud compute instances create nested-vm --zone europe-west1-c \
    --min-cpu-platform "Intel Haswell" \
    --machine-type=n1-standard-16 \
    --image nested-vm-image
    ```   
1. SSH into the newly created VM instance and check nested virtualizaiton is enabled:
    ```
    gcloud beta compute ssh --zone "europe-west1-c" "nested-vm" --project "open-kernel-monitoring"
    grep -cw vmx /proc/cpuinfo
    # should return non-zero!
    ```


## Build Linux Kernel with Statsfs Patch on Ubuntu18.04VM 

Originally from https://github.com/esposem/linux, now at https://github.com/liiling/linux:

1. `sudo apt install -y --fix-missing git build-essential libncurses-dev bison flex libssl-dev libelf-dev bc`
1. `git clone https://github.com/liiling/linux.git ; cd linux`
1. `git checkout statsfs`

1. Compile the Linux kernel using all available cpu threads: `time make -j $(nproc)`
1. Install the Linux kernel modules: `sudo make modules_install`
1. Install the Linux kernel: `sudo make install`. The following files are installed to the `/boot` directory, and the grub configuration is updated.
    - config-5.7.0-rc2+
    - System.map-5.7.0-rc2+
    - vmlinuz-5.7.0-rc2+
1. Reboot the VM: `reboot`
1. After reboot, check that statsfs is supported in the running Linux kernel:
    - check Linux kernel version: `uname -mrs`
    - check statsfs filesystem is supported: `cat /proc/filesystems | grep statsfs`
1. Mount statfs: `sudo mount -t statsfs statsfs /sys/kernel/stats`
1. Check statsfs is mounted: `cat /proc/mounts | grep stats`
1. Change permission of statsfs filesystem to be readable and executable for everyone, but writable only by the owner (root): `sudo chmod -R 755 /sys/kernel/stats`

The `statsfs` branch on `liiling/linux` includes a Kernel config, so this is not required:

1. Clean the kernel tree: `make mrproper`
1. Generate `.config` file. Some possible alternatives:
	- `make localmodconfig` (generate a config from the kernel options currently in use)
	- `make menuconfig` (command line interface for config creation)
	- `cp -v /boot/config-$(uname -r) .config` (copy the boot config of the host machine)
   
Double check that stats_fs pseudo filesystem is enabled, i.e. set `CONFIG_STATS_FS=y`

## Install the user-space agent

### Install Golang
1. `curl -L https://golang.org/dl/go1.15.2.src.tar.gz > go1.15.2.src.tar.gz`
2. Might require `sudo`: `tar -C /usr/local -xzf go1.15.2.src.tar.gz`
3. `export PATH=$PATH:/usr/local/go/bin`

### Pull the agent from Git and running it
1. `git clone https://github.com/liiling/kernel-metrics-agent.git`
2. `cd kernel-metrics-agent`
3. `git checkout statsfs-metadata`
4. `cd statsfs`
5. Symlink otelstats to GOROOT directory (since it is not published as a Go package)...`ln -s /home/your-username/kernel-metrics-agent/statsfs/otelstats /usr/local/go/src/otelstats`
6. `go run main.go -exporter gcp -statsfspath /sys/kernel/statsfs`
7. Check metrics in Google Cloud Platform's monitoring page

## Instrument Statsfs with OpenTelemetry

### Statsfs filesystem structure

Assume statsfs is mounted at `/sys/kernel/stats/` (statsfs root), each Linux
subsystem with statsfs metrics should appear as a directory under statsfs
root. 
Each device or subdevice should appear as a directory under the subsystem directory, and each metric appears as a file with filename being the metric reported. 
The folder [otelstats/testsys](https://github.com/liiling/kernel-metrics-agent/tree/statsfs/statsfs/otelstats/testsys/) shows a test example of statsfs filesystem.

The current statsfs documentation is found [here](https://github.com/esposem/linux/blob/35624f8292988e2f3189c1b4d2cb503a47230df0/Documentation/filesystems/stats_fs.rst).

Example tree structure:
```bash
sys
└── kernel
    └── stats
        ├── subsys0
        │   ├── dev0
        │   │   ├── m0
        │   │   └── m1
        │   └── dev1
        │       └── m0
        └── subsys1
            ├── dev0
            │   ├── in_all_m
            │   ├── in_both_devs_m
            │   ├── in_top_and_dev0_m
            │   └── only_in_dev0_m
            ├── dev1
            │   ├── in_all_m
            │   └── in_both_devs_m
            ├── in_all_m
            ├── in_top_and_dev0_m
            └── top_level_m
```

### Statsfs to OpenTelemetry Metrics

Metrics in OpenTelemetry are defined by their semantically meaningful and unique names. Each metric could be associated with multiple labels (key-value pairs).
One way to transform statsfs metrics to OpenTelemetry metrics is using `subsys/metric_filename` as OpenTelemetry metric name, and the device path as label.
For example, a statsfs metric file `sys/kernel/stats/subsys0/dev0/subdev0/metric0` denotes a statsfs metric named `metric0` for a device `dev0/subdev0` under a Linux subsystem `subsys0`.
When ported to OpenTelemetry, the metric will appear with name `subsys0/metric0`, and label `device=/dev0/subdev0`.

### Reflections on Combining Statsfs and OpenTelemetry

1. Lack of metadata from statsfs:

    OpenTelemetry provides six types of [metric instruments](https://github.com/open-telemetry/opentelemetry-specification/blob/master/specification/metrics/api.md#metric-instruments), three of which (SumObserver, UpDownSumObserver, ValueObserver) are asynchronous. 
    [Statsfs](https://github.com/esposem/linux/blob/35624f8292988e2f3189c1b4d2cb503a47230df0/Documentation/filesystems/stats_fs.rst) provides two metric value types: cumulative and floating, where cumulative values are forever increasing (maps to SumObserver in OpenTelemetry) and floating value could exhibit different behaviors (maps to UpDownSumObserver or ValueObserver in OpenTelemetry).
    Although the kernel developers could write code to decide the value type of the exposed statsfs metrics, the exposed metric files themselves have no information on the metric type.
    Currently, the demo uses ValueObserver, the most generous, all-encompassing OpenTelemetry instrument for all statsfs metrics. However, this behavior is not ideal. For instance, it allows strictly increasing metrics to exhibit floating behavior.

2. Different ways of viewing metrics:

    Statsfs is organised as `Linux subsystem -> device -> subdevice -> metrics`. The same metrics can be added to many different sources/devices.
    OpenTelemetry is organised as `metrics & a list of labels`. Each metrics has a unique name, with labels specifying details such as devices.

3. Number of I/O operation per metrics:

    In statsfs, the same metric for different devices in the same subsystem are spread across multiple files. 
    This implies that with the current design of the statsfs -> OpenTelemetry demo, exporting one metrics (along with the list of all associated labels) requires many file I/O operations.

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
