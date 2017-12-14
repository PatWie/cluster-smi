#include "bridge.h"

nvmlReturn_t bridge_get_text_property(getNvmlCharProperty f, nvmlDevice_t device, char *buf, unsigned int length)
{
    nvmlReturn_t ret;
    ret = f(device, buf, length);
    return ret;
}

nvmlReturn_t bridge_get_int_property(getNvmlIntProperty f, nvmlDevice_t device, unsigned int *value)
{
    nvmlReturn_t ret;
    ret = f(device, value);
    return ret;
}