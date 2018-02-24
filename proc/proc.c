// Author: Patrick Wieschollek, 2018
#include <stdio.h>
#include <stdbool.h>
#include <unistd.h>
#include <stdlib.h>
#include <string.h>
#include <sys/time.h>

#define MAX_NAME 128

void clock_ticks(long int *hz) {
  *hz = sysconf(_SC_CLK_TCK);
}

void get_cmd(unsigned long  pid, char* cmd) {

  char path[40];
  snprintf(path, 40, "/proc/%ld/cmdline", pid);

  printf("read from %s\n", path);

  char line[100], *p;
  FILE* fp;

  fp = fopen(path, "r");
  if (!fp)
    return;

  size_t len = 0;
  ssize_t read;
  while ((read = getline(&cmd, &len, fp)) != -1) {
    printf("Retrieved line of length %zu :\n", read);
    printf("%s", line);
  }
  // bool success = fscanf(fp, "%s", cmd) != EOF;
  printf("found %s\n", cmd);
  fclose(fp);
}

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

void time_of_day(float *current_time) {
  struct timeval tv;
  gettimeofday(&tv, 0);
  *current_time = (float) tv.tv_sec;
}

int boot_time(float *uptime, float *idle) {
  FILE *fp;
  fp = fopen("/proc/uptime", "r");
  if (fp != NULL) {
    if (fscanf(fp, "%f %f", uptime, idle) == EOF) {
      fclose(fp);
      return 0;
    } else {
      fclose(fp);
      return 1;
    }
  }
  return 0;
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
void read_pid_info(unsigned long pid, unsigned long *time, unsigned long long *starttime, char *name) {

  char fn[MAX_NAME + 1];
  snprintf(fn, sizeof fn, "/proc/%ld/stat", pid);

  unsigned long utime = 0;
  unsigned long stime = 0;

  *time = 0;

  FILE * fp;
  fp = fopen(fn, "r");
  if (fp != NULL) {
    /*
    // https://kb.novaordis.com/index.php/Linux_Process_Information
    (1) pid           %d  The process ID
    (2) comm          %s  The filename of the executable, in parentheses.
    (3) state         %c
    (4) ppid          %d  The PID of the parent of this process.
    (5) pgrp          %d  The process group ID of the process.

    (6) session       %d  The session ID of the process.
    (7) tty_nr        %d  The controlling terminal of the process.
    (8) tpgid         %d  The ID of the foreground process group
    (9) flags         %u  The kernel flags word of the process.

    (10) minflt       %lu The number of minor faults the process has made
    (11) cminflt      %lu The number of minor faults that the process's waited-for children have made.
    (12) majflt       %lu The number of major faults the process has made
    (13) cmajflt      %lu The number of major faults that the process's
    (14) utime        %lu Amount of time that this process has been scheduled in user mode
    (15) stime        %lu Amount of time that this process has been scheduled in kernel mode

    (16) cutime       %ld Amount of time that this process's waited-for chil
    (17) cstime       %ld Amount of time that this process's waited-for chil
    (18) priority     %ld
    (19) nice         %ld The nice value (see setpriority(2))
    (20) num_threads  %ld Number of threads in this process
    (21) itrealvalue  %ld The time in jiffies before the next SIGALRM
    (22) starttime    %llu The time the process started after system boot (clock ticks (divide by sysconf(_SC_CLK_TCK)).

    ...
    */

    // extract                   1   2   3   4   5   6   7   8   9   10   11   12   13  14  15   16   17  18  19  20  21    22
    bool success = fscanf(fp, "%*d (%s %*c %*d %*d %*d %*d %*d %*u %*lu %*lu %*lu %*lu %lu %lu %*ld %*ld %*d %*d %*d %*d %llu",
                          name, &utime, &stime, starttime) != EOF;
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

