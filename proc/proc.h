// Author: Patrick Wieschollek, 2018
#ifndef GOCODE_Proc_H
#define GOCODE_Proc_H

void get_mem(unsigned long *mem_total, unsigned long *mem_free, unsigned long *mem_available);
unsigned long long int read_cpu_tick();
void get_uid_from_pid(unsigned long pid, unsigned long *uid);
void read_time_and_name_from_pid(unsigned long pid, unsigned long *time, char *name);
unsigned int num_cores();

#endif // GOCODE_Proc_H