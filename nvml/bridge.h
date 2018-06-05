#pragma once
#include "nvml.h"
#include <stdlib.h>

typedef nvmlReturn_t (*getNvmlCharProperty) (nvmlDevice_t device , char *buf, unsigned int length);
typedef nvmlReturn_t (*getNvmlIntProperty) (nvmlDevice_t device, unsigned int *value);

nvmlReturn_t bridge_get_text_property(getNvmlCharProperty f, nvmlDevice_t device, char *buf, unsigned int length);
nvmlReturn_t bridge_get_int_property(getNvmlIntProperty f, nvmlDevice_t device, unsigned int *value);