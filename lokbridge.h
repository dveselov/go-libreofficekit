#include <stdlib.h>
#include "LibreOfficeKit/LibreOfficeKitInit.h"
#include "LibreOfficeKit/LibreOfficeKit.h"


typedef void (*voidFunc) ();
typedef int (*intFunc) ();
typedef char* (*charFunc) ();
void destroy_bridge(voidFunc f, void* handle) {
    return f(handle);
};
LibreOfficeKitDocument* document_load_bridge(voidFunc f,
                                             LibreOfficeKit* pThis,
                                             const char* pURL) {
    f(pThis, pURL);
};
char* get_error_bridge(charFunc f, LibreOfficeKit* pThis) {
    return f(pThis);
};
int get_document_bridge(intFunc f, LibreOfficeKitDocument* pThis) {
    return f(pThis);
};
void get_document_size_bridge(voidFunc f,
                            LibreOfficeKitDocument* pThis,
                            long* pWidth,
                            long* pHeight) {
    return f(pThis, pWidth, pHeight);
};
void set_document_part_bridge(voidFunc f,
                            LibreOfficeKitDocument* pThis,
                            int nPart) {
    return f(pThis, nPart);
};

char* get_document_part_name_bridge(charFunc f,
                            LibreOfficeKitDocument* pThis,
                            int nPart) {
    return f(pThis, nPart);
};

void initialize_for_rendering_bridge(voidFunc f,
                                    LibreOfficeKitDocument* pThis,
                                    const char* pArguments) {
    return f(pThis, pArguments);
};

int document_save_bridge(intFunc f,
                        LibreOfficeKitDocument* pThis,
                        const char* pUrl,
                        const char* pFormat,
                        const char* pFilterOptions) {
    return f(pThis, pUrl, pFormat, pFilterOptions);
};
