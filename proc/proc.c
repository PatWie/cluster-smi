// Author: Patrick Wieschollek, 2018
#include <stdio.h>
#include <stdbool.h>
#include <unistd.h>
#include <stdlib.h>
#include <string.h>

#define MAX_NAME 128

void get_mem(unsigned long *mem_total, unsigned long *mem_free, unsigned long *mem_available) {
  char line[100], *p;
  FILE* statusf;


  statusf = fopen("/proc/meminfo", "r");
  if (!statusf)
    return;


  fgets(line, 100, statusf);
  sscanf(line, "%*s %lu %*s", mem_total);
  fgets(line, 100, statusf);
  sscanf(line, "%*s %lu %*s", mem_free);
  fgets(line, 100, statusf);
  sscanf(line, "%*s %lu %*s", mem_available);


  fclose(statusf);
}

// read total cpu time
unsigned long long int read_cpu_tick() {
  unsigned long long int usertime, nicetime, systemtime, idletime;
  unsigned long long int ioWait, irq, softIrq, steal, guest, guestnice;
  usertime = nicetime = systemtime = idletime = 0;
  ioWait = irq = softIrq = steal = guest = guestnice = 0;

  FILE *fp;
  fp = fopen("/proc/stat", "r");
  if (fp != NULL) {
    if (fscanf(fp,   "cpu  %16llu %16llu %16llu %16llu %16llu %16llu %16llu %16llu %16llu %16llu",
               &usertime, &nicetime, &systemtime, &idletime,
               &ioWait, &irq, &softIrq, &steal, &guest, &guestnice) == EOF) {
      fclose(fp);
      return 0;
    } else {
      fclose(fp);
      return usertime + nicetime + systemtime + idletime + ioWait + irq + softIrq + steal + guest + guestnice;
    }
  } else {
    return 0;
  }
}

void get_uid_from_pid(unsigned long pid, unsigned long *uid) {
  char path[40], line[100], *p;
  FILE* statusf;

  snprintf(path, 40, "/proc/%ld/status", pid);

  statusf = fopen(path, "r");
  if (!statusf)
    return;

  while (fgets(line, 100, statusf)) {
    if (strncmp(line, "Uid:", 4) != 0)
      continue;
    // Uid: 1000    1000    1000    1000
    sscanf(line, "%*s %lu %*s", uid);
    break;
  }
  fclose(statusf);
}

// read cpu tick for a specific process
void read_time_and_name_from_pid(unsigned long pid, unsigned long *time, char *name) {

  char fn[MAX_NAME + 1];
  snprintf(fn, sizeof fn, "/proc/%ld/stat", pid);

  unsigned long utime = 0;
  unsigned long stime = 0;

  *time = 0;

  FILE * fp;
  fp = fopen(fn, "r");
  if (fp != NULL) {
    /*
    (1) pid       %d  The process ID
    (2) comm      %s  The filename of the executable, in parentheses.
    (3) state     %c
    (4) ppid      %d  The PID of the parent of this process.
    (5) pgrp      %d  The process group ID of the process.
    (6) session   %d  The session ID of the process.
    (7) tty_nr    %d  The controlling terminal of the process.
    (8) tpgid     %d  The ID of the foreground process group 
    (9) flags     %u  The kernel flags word of the process. 
    (10) minflt   %lu The number of minor faults the process has made
    (11) cminflt  %lu The number of minor faults that the process's waited-for children have made.
    (12) majflt   %lu The number of major faults the process has made
    (13) cmajflt  %lu The number of major faults that the process's
    (14) utime    %lu Amount of time that this process has been scheduled in user mode
    (15) stime    %lu Amount of time that this process has been scheduled in kernel mode
    ...
    */

    // extract
    bool success = fscanf(fp, "%*d (%s %*c %*d %*d %*d %*d %*d %*u %*lu %*lu %*lu %*lu %lu"
                          "%lu %*ld %*ld %*d %*d %*d %*d %*u %*lu %*ld",
                          name, &utime, &stime) != EOF;
    fclose(fp);

    if (!success) {
      // something went wrong
      return;
    }

    // remove ")" suffix
    if (strlen(name) > 2) {
      name[strlen(name) - 1] = 0;
    }

    *time = utime + stime;

  }
}

// return number of cores
unsigned int num_cores() {
  return sysconf(_SC_NPROCESSORS_ONLN);
}

