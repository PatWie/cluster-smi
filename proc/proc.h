// Author: Patrick Wieschollek, 2018
#ifndef GOCODE_Proc_H
#define GOCODE_Proc_H

void get_cmd(unsigned long  pid, char* cmd);
void clock_ticks(long int *hz);
void time_of_day(float *current_time);
int boot_time(float *uptime, float *idle);
void get_mem(unsigned long *mem_total, unsigned long *mem_free, unsigned long *mem_available);
unsigned long long int read_cpu_tick();
void get_uid_from_pid(unsigned long pid, unsigned long *uid);
void read_pid_info(unsigned long pid, unsigned long *time, unsigned long long *starttime, char *name);
unsigned int num_cores();

#endif // GOCODE_Proc_H