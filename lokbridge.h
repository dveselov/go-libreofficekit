#include <stdlib.h>
#include "LibreOfficeKit/LibreOfficeKitInit.h"
#include "LibreOfficeKit/LibreOfficeKit.h"


void destroy_office(LibreOfficeKit* pThis) {
    return pThis->pClass->destroy(pThis);
};

char* get_error(LibreOfficeKit* pThis) {
    return pThis->pClass->getError(pThis);
};

void free_error(LibreOfficeKit* pThis, char* message) {
    return pThis->pClass->freeError(message);
}

char* get_filter_types(LibreOfficeKit* pThis) {
    return pThis->pClass->getFilterTypes(pThis);
};

LibreOfficeKitDocument* document_load(LibreOfficeKit* pThis, const char* pURL) {
    return pThis->pClass->documentLoad(pThis, pURL);
};

void destroy_document(LibreOfficeKitDocument* pThis) {
    return pThis->pClass->destroy(pThis);
};

void get_document_size(LibreOfficeKitDocument* pThis, long* pWidth, long* pHeight) {
    return pThis->pClass->getDocumentSize(pThis, pWidth, pHeight);
};
void set_document_part(LibreOfficeKitDocument* pThis, int nPart) {
    return pThis->pClass->setPart(pThis, nPart);
};

int get_document_type(LibreOfficeKitDocument* pThis) {
    return pThis->pClass->getDocumentType(pThis);
};

int get_document_parts(LibreOfficeKitDocument* pThis) {
    return pThis->pClass->getParts(pThis);
};

int get_document_part(LibreOfficeKitDocument* pThis) {
    return pThis->pClass->getPart(pThis);
};

char* get_document_part_name(LibreOfficeKitDocument* pThis, int nPart) {
    return pThis->pClass->getPartName(pThis, nPart);
};

void initialize_for_rendering(LibreOfficeKitDocument* pThis, const char* pArguments) {
    return pThis->pClass->initializeForRendering(pThis, pArguments);
};

int document_save(LibreOfficeKitDocument* pThis, const char* pUrl, const char* pFormat, const char* pFilterOptions) {
    return pThis->pClass->saveAs(pThis, pUrl, pFormat, pFilterOptions);
};

int create_view(LibreOfficeKitDocument* pThis) {
    return pThis->pClass->createView(pThis);
};

int get_view(LibreOfficeKitDocument* pThis) {
    return pThis->pClass->getView(pThis);
};

int get_views(LibreOfficeKitDocument* pThis) {
    return pThis->pClass->getViews(pThis);
};

void paint_tile(LibreOfficeKitDocument* pThis, unsigned char* pBuffer, const int nCanvasWidth, const int nCanvasHeight, const int nTilePosX, const int nTilePosY, const int nTileWidth, const int nTileHeight) {
    return pThis->pClass->paintTile(pThis, pBuffer, nCanvasWidth, nCanvasHeight, nTilePosX, nTilePosY, nTileWidth, nTileHeight);
};

int get_tile_mode(LibreOfficeKitDocument* pThis) {
    return pThis->pClass->getTileMode(pThis);
};

char* get_part_page_rectangles(LibreOfficeKitDocument* pThis) {
    return pThis->pClass->getPartPageRectangles(pThis);
};

void set_text_selection(LibreOfficeKitDocument* pThis, int nType, int nX, int nY) {
    return pThis->pClass->setTextSelection(pThis, nType, nX, nY);
};

void reset_selection(LibreOfficeKitDocument* pThis) {
    return pThis->pClass->resetSelection(pThis);
};

char* get_text_selection(LibreOfficeKitDocument* pThis, const char* pMimeType) {
    return pThis->pClass->getTextSelection(pThis, pMimeType, NULL);
};
