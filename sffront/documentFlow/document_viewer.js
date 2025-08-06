// Variables para PDF.js
let pdfDoc = null,
    pageNum = 1,
    pageRendering = false,
    pageNumPending = null,
    scale = 1.0,
    pdfAspectRatio = 8.5 / 11; 

// Arreglo de marcadores 
let markers = [];
let markerColors = ['#003b65aa', '#861d7daa', '#5c7500aa', '#0086aeaa', '#160064aa', '#611500aa'];

// Elementos del DOM
let pdfCanvas, overlayCanvas, overlayCtx, pdfCtx;
let isDragging = false;
let isResizing = false;
let currentDraggedReviewer = null;
let selectedTextBox = null;
let dragOffsetX = 0, dragOffsetY = 0;

// Propiedades del PDF real (en puntos)
let pdfRealWidth = 612; // Ancho estándar carta en puntos
let pdfRealHeight = 792; // Alto estándar carta en puntos

// Configura PDF.js worker al inicio
pdfjsLib.GlobalWorkerOptions.workerSrc = 'https://cdnjs.cloudflare.com/ajax/libs/pdf.js/3.4.120/pdf.worker.min.js';

// Cargar datos del documento al cargar la página
document.addEventListener('DOMContentLoaded', function() {
    initPDFElements();
    const documentData = JSON.parse(localStorage.getItem('currentDocument'));
    
    if (documentData) {
        // Mostrar información del documento
        const docTitleElement = document.querySelector('.flexbox.min-w-72 p');
        if (docTitleElement) {
            docTitleElement.innerHTML = `${documentData.fileName} 
                <span class="text-secondary text-sm font-normal leading-normal">Uploaded by ${documentData.uploadedBy} on ${documentData.uploadDate}</span>`;
        }
        
        // Mostrar preview del PDF si existe
        if (documentData.fileBase64 && documentData.fileMime === 'application/pdf') {
            loadPDF(documentData.fileBase64);
        }

        // Actualizar tabla de approvals
        updateApprovalTable(documentData.reviewers);
        
        setupDraggables();
        setupCanvasInteractions();
        loadExistingMarkers();
    }
});

function initPDFElements() {
    // Obtener elementos del DOM
    pdfCanvas = document.getElementById('pdf-canvas');
    overlayCanvas = document.getElementById('overlay-canvas');
    
    pdfCtx = pdfCanvas.getContext('2d');
    overlayCtx = overlayCanvas.getContext('2d');
    
    // Configurar eventos de los controles
    document.getElementById('prev-page').addEventListener('click', onPrevPage);
    document.getElementById('next-page').addEventListener('click', onNextPage);
    document.getElementById('view-full').addEventListener('click', openFullPDF);
    
    // Redimensionar cuando cambia el tamaño de la ventana
    window.addEventListener('resize', function() {
        if (pdfDoc) {
            resizeCanvases();
            renderPage(pageNum);
        }
    });
}

function updateApprovalTable(reviewers) {
    const tableBody = document.querySelector('.glass-card table tbody');
    if (tableBody) {
        tableBody.innerHTML = '';
        
        reviewers.forEach(reviewer => {
            const row = document.createElement('tr');
            row.className = 'border-t border-t-white/10';
            
            const statusClass = {
                'Approved': 'bg-green-500/20 text-green-400',
                'Pending': 'bg-yellow-500/20 text-yellow-400',
                'Not Started': 'bg-white/20 text-primary'
            }[reviewer.status] || 'bg-white/20 text-primary';
            
            row.innerHTML = `
                <td class="h-[72px] px-4 py-2 w-[400px] text-primary text-sm font-normal leading-normal">
                    ${reviewer.team}
                </td>
                <td class="h-[72px] px-4 py-2 w-60 text-sm font-normal leading-normal">
                    <button class="flex min-w-[84px] max-w-[480px] cursor-default items-center justify-center overflow-hidden rounded-lg h-8 px-4 ${statusClass} text-sm font-medium leading-normal w-full">
                        <span class="truncate">${reviewer.status}</span>
                    </button>
                </td>
                <td class="h-[72px] px-4 py-2 w-[400px] text-secondary text-sm font-normal leading-normal">
                    ${reviewer.due_date}
                </td>
                <td class="h-[72px] px-4 py-2 w-[400px] text-secondary text-sm font-normal leading-normal">
                    ${reviewer.user}
                </td>
            `;
            
            tableBody.appendChild(row);
        });
    }
}

function loadPDF(base64Data) {
    const byteCharacters = atob(base64Data);
    const byteArray = new Uint8Array(byteCharacters.length);
    for (let i = 0; i < byteCharacters.length; i++) {
        byteArray[i] = byteCharacters.charCodeAt(i);
    }

    const loadingTask = pdfjsLib.getDocument({ data: byteArray });

    loadingTask.promise.then(function(pdf) {
        pdfDoc = pdf;
        document.getElementById('page-info').textContent = `Página 1 de ${pdf.numPages}`;
        
        // Obtener dimensiones reales del PDF
        pdf.getPage(1).then(function(page) {
            const viewport = page.getViewport({scale: 1.0});
            pdfRealWidth = viewport.width;
            pdfRealHeight = viewport.height;
            pdfAspectRatio = pdfRealWidth / pdfRealHeight;
            
            resizeCanvases();
            renderPage(1);
        });
    }).catch(function(error) {
        console.error('Error al cargar PDF:', error);
    });
}

function resizeCanvases() {
    const container = document.getElementById('preview-container');
    const containerWidth = container.clientWidth;
    const containerHeight = container.clientHeight;
    
    // Calcular tamaño manteniendo aspect ratio
    let canvasWidth = containerWidth;
    let canvasHeight = containerWidth / pdfAspectRatio;
    
    // Si el alto calculado es mayor que el contenedor, ajustar por altura
    if (canvasHeight > containerHeight) {
        canvasHeight = containerHeight;
        canvasWidth = containerHeight * pdfAspectRatio;
    }
    
    // Aplicar tamaños a ambos canvas
    [pdfCanvas, overlayCanvas].forEach(canvas => {
        canvas.style.width = `${canvasWidth}px`;
        canvas.style.height = `${canvasHeight}px`;
        canvas.width = canvasWidth;
        canvas.height = canvasHeight;
    });
}

function renderPage(num) {
    pageRendering = true;
    
    pdfDoc.getPage(num).then(function(page) {
        const container = document.getElementById('preview-container');
        const containerWidth = container.clientWidth;
        const containerHeight = container.clientHeight;
        
        // 1. Obtener viewport con scale=1 para conocer el aspect ratio original
        const viewport = page.getViewport({scale: 1});
        const pdfAspectRatio = viewport.width / viewport.height;
        
        // 2. Calcular la escala para ajustar al contenedor
        const containerAspectRatio = containerWidth / containerHeight;
        let scale;
        
        if (pdfAspectRatio > containerAspectRatio) {
            // Ajustar por ancho
            scale = containerWidth / viewport.width;
        } else {
            // Ajustar por alto
            scale = containerHeight / viewport.height;
        }
        
        // 3. Crear viewport con la escala calculada
        const scaledViewport = page.getViewport({scale: scale});
        
        // 4. Ajustar el canvas al tamaño calculado
        pdfCanvas.width = scaledViewport.width;
        pdfCanvas.height = scaledViewport.height;
        
        // 5. Renderizar con el viewport escalado
        const renderContext = {
            canvasContext: pdfCtx,
            viewport: scaledViewport
        };
        
        const renderTask = page.render(renderContext);
        
        renderTask.promise.then(function() {
            pageRendering = false;
            if (pageNumPending !== null) {
                renderPage(pageNumPending);
                pageNumPending = null;
            }
            redrawMarkers();
        });
    });
    
    document.getElementById('page-info').textContent = `Página ${num} de ${pdfDoc.numPages}`;
}

function onPrevPage() {
    if (pageNum <= 1) return;
    pageNum--;
    queueRenderPage(pageNum);
}

function onNextPage() {
    if (pageNum >= pdfDoc.numPages) return;
    pageNum++;
    queueRenderPage(pageNum);
}

function queueRenderPage(num) {
    if (pageRendering) {
        pageNumPending = num;
    } else {
        renderPage(num);
    }
}

function openFullPDF() {
    const documentData = JSON.parse(localStorage.getItem('currentDocument'));
    if (documentData && documentData.fileBase64) {
        const pdfWindow = window.open('', '_blank');
        pdfWindow.document.write(`
            <html>
                <head>
                    <title>${documentData.fileName}</title>
                    <style>
                        body { margin: 0; }
                        embed { width: 100%; height: 100vh; }
                    </style>
                </head>
                <body>
                    <embed src="data:application/pdf;base64,${documentData.fileBase64}#page=${pageNum}&zoom=100%">
                </body>
            </html>
        `);
    }
}

// Configurar interacciones del canvas overlay
function setupCanvasInteractions() {
    // Drag and drop desde la tabla
    overlayCanvas.addEventListener('dragover', (e) => {
        e.preventDefault();
        overlayCanvas.style.cursor = 'copy';
    });

    overlayCanvas.addEventListener('dragleave', () => {
        overlayCanvas.style.cursor = 'default';
    });

    overlayCanvas.addEventListener('drop', (e) => {
        e.preventDefault();
        overlayCanvas.style.cursor = 'default';
        
        if (!currentDraggedReviewer) return;
        
        const rect = overlayCanvas.getBoundingClientRect();
        const x = e.clientX - rect.left;
        const y = e.clientY - rect.top;
    
        // Crear nuevo marcador
        const newMarker = {
            user: currentDraggedReviewer.user,
            x: x,
            y: y,
            width: 100,
            height: 50,
            page: pageNum,
            bgColor: markerColors[markers.length%markerColors.length]
        };
        
        markers.push(newMarker);
        redrawMarkers();
        
        // Enviar a API
        //sendMarkerToAPI(newMarker);
        
        currentDraggedReviewer = null;
    });

    // Interacciones con las cajas de texto
    overlayCanvas.addEventListener('mousedown', handleMouseDown);
    overlayCanvas.addEventListener('mousemove', handleMouseMove);
    overlayCanvas.addEventListener('mouseup', handleMouseUp);
    overlayCanvas.addEventListener('dblclick', handleDoubleClick);
}

function handleMouseDown(e) {
    const rect = overlayCanvas.getBoundingClientRect();
    const x = e.clientX - rect.left;
    const y = e.clientY - rect.top;
    
    const clickedMarker = getMarkerAtPosition(x, y);
    
    if (clickedMarker) {
        selectedTextBox = clickedMarker;
        
        // Verificar si está en el área de redimensionamiento (esquina inferior derecha)
        const isInResizeArea = x >= clickedMarker.x + clickedMarker.width - 10 && 
                              y >= clickedMarker.y + clickedMarker.height - 10;
        
        if (isInResizeArea) {
            isResizing = true;
            overlayCanvas.style.cursor = 'nw-resize';
        } else {
            isDragging = true;
            dragOffsetX = x - clickedMarker.x;
            dragOffsetY = y - clickedMarker.y;
            overlayCanvas.style.cursor = 'move';
        }
    }
}

function handleMouseMove(e) {
    const rect = overlayCanvas.getBoundingClientRect();
    const x = e.clientX - rect.left;
    const y = e.clientY - rect.top;
    
    if (isDragging && selectedTextBox) {
        selectedTextBox.x = Math.max(0, Math.min(x - dragOffsetX, overlayCanvas.width - selectedTextBox.width));
        selectedTextBox.y = Math.max(0, Math.min(y - dragOffsetY, overlayCanvas.height - selectedTextBox.height));
        redrawMarkers();
    } else if (isResizing && selectedTextBox) {
        selectedTextBox.width = Math.max(50, x - selectedTextBox.x);
        selectedTextBox.height = Math.max(25, y - selectedTextBox.y);
        redrawMarkers();
    } else {
        // Cambiar cursor según la posición
        const marker = getMarkerAtPosition(x, y);
        if (marker) {
            const isInResizeArea = x >= marker.x + marker.width - 10 && 
                                  y >= marker.y + marker.height - 10;
            overlayCanvas.style.cursor = isInResizeArea ? 'nw-resize' : 'move';
        } else {
            overlayCanvas.style.cursor = 'default';
        }
    }
}

function handleMouseUp() {
    if (isDragging || isResizing) {
        // Enviar actualización a API
        if (selectedTextBox) {
            sendMarkerToAPI(selectedTextBox);
        }
    }
    
    isDragging = false;
    isResizing = false;
    selectedTextBox = null;
    overlayCanvas.style.cursor = 'default';
}

function handleDoubleClick(e) {
    const rect = overlayCanvas.getBoundingClientRect();
    const x = e.clientX - rect.left;
    const y = e.clientY - rect.top;
    
    const clickedMarker = getMarkerAtPosition(x, y);
    
    if (clickedMarker) {
        // Eliminar marcador
        markers = markers.filter(m => m !== clickedMarker);
        redrawMarkers();
        
        // Enviar eliminación a API
        sendMarkerDeletion(clickedMarker);
    }
}

function getMarkerAtPosition(x, y) {
    // Buscar en orden inverso para obtener el marcador más arriba
    for (let i = markers.length - 1; i >= 0; i--) {
        const marker = markers[i];
        if (marker.page === pageNum &&
            x >= marker.x && x <= marker.x + marker.width &&
            y >= marker.y && y <= marker.y + marker.height) {
            return marker;
        }
    }
    return null;
}

function redrawMarkers() {
    // Limpiar canvas overlay
    overlayCtx.clearRect(0, 0, overlayCanvas.width, overlayCanvas.height);
    
    // Dibujar marcadores de la página actual
    const currentPageMarkers = markers.filter(marker => marker.page === pageNum);
    
    currentPageMarkers.forEach(marker => {
        drawTextBox(marker);
    });
}

function drawTextBox(marker) {
    const ctx = overlayCtx;
    const radius = 5; // Ajusta este valor para cambiar el redondeo

    // Dibujar fondo de la caja con bordes redondeados
    ctx.fillStyle = marker.bgColor;
    ctx.beginPath();
    ctx.roundRect(marker.x, marker.y, marker.width, marker.height, radius);
    ctx.fill();

    // Dibujar borde redondeado
    ctx.strokeStyle = 'rgba(0, 0, 0, 0.1)';
    ctx.lineWidth = 2;
    ctx.beginPath();
    ctx.roundRect(marker.x, marker.y, marker.width, marker.height, radius);
    ctx.stroke();
    
    // Dibujar texto
    ctx.fillStyle = 'white';
    ctx.font = '14px Arial';
    ctx.textAlign = 'center';
    ctx.textBaseline = 'middle';
    
    // Truncar texto si es muy largo
    let displayText = marker.user;
    const maxWidth = marker.width - 10;
    if (ctx.measureText(displayText).width > maxWidth) {
        displayText = displayText.substring(0, Math.floor(displayText.length * maxWidth / ctx.measureText(displayText).width)) + '...';
    }
    
    ctx.fillText(displayText, marker.x + marker.width / 2, marker.y + marker.height / 2);
    
    // Dibujar indicador de redimensionamiento
    if (selectedTextBox === marker) {
        ctx.fillStyle = 'white';
        ctx.fillRect(marker.x + marker.width - 8, marker.y + marker.height - 8, 8, 8);
    }
}

// Funciones para convertir entre coordenadas del canvas y del PDF real
function canvasToRealPDF(canvasX, canvasY, canvasWidth, canvasHeight) {
    const scaleX = pdfRealWidth / overlayCanvas.width;
    const scaleY = pdfRealHeight / overlayCanvas.height;
    
    return {
        x: canvasX * scaleX,
        y: canvasY * scaleY,
        width: canvasWidth * scaleX,
        height: canvasHeight * scaleY
    };
}

function realPDFToCanvas(realX, realY, realWidth, realHeight) {
    const scaleX = overlayCanvas.width / pdfRealWidth;
    const scaleY = overlayCanvas.height / pdfRealHeight;
    
    return {
        x: realX * scaleX,
        y: realY * scaleY,
        width: realWidth * scaleX,
        height: realHeight * scaleY
    };
}

// Hacer elementos arrastrables (filas de la tabla)
function setupDraggables() {
    const table = document.querySelector('.glass-card table');
    const rows = table.querySelectorAll('tbody tr');
    
    rows.forEach(row => {
        row.draggable = true;
        
        row.addEventListener('dragstart', (e) => {
            const user = row.cells[3].textContent.split(',')[0].trim();
            currentDraggedReviewer = { user };
            e.dataTransfer.setData('text/plain', JSON.stringify({ user }));
        });
    });
}

// Cargar marcadores existentes
function loadExistingMarkers() {
    try {
        // Simulación de marcadores existentes
        const existingMarkers = []; 
        markers.push(...existingMarkers);
        redrawMarkers();
    } catch (error) {
        console.error('Error cargando marcadores:', error);
    }
}

// Enviar datos a API
async function sendMarkerToAPI(marker) {
    // Convertir a coordenadas reales del PDF
    const realCoords = canvasToRealPDF(marker.x, marker.y, marker.width, marker.height);
    
    const markerData = {
        document: localStorage.getItem('currentDocumentFileName'),
        reviewer: marker.user,
        position: {
            x: realCoords.x,
            y: realCoords.y,
            width: realCoords.width,
            height: realCoords.height
        },
        page: marker.page
    };
    
    console.log('Enviando marcador a API:', JSON.stringify(markerData));
    
    /*
    try {
        const response = await fetch('https://tu-api.com/markers', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': 'Bearer tu-token'
            },
            body: JSON.stringify(markerData)
        });
        
        if (!response.ok) {
            console.error('Error al enviar marcador:', await response.text());
        }
    } catch (error) {
        console.error('Error de red:', error);
    }
    */
}

async function sendMarkerDeletion(marker) {
    console.log('Eliminando marcador:', marker.user, 'página', marker.page);
    
    /*
    try {
        const response = await fetch(`https://tu-api.com/markers/${marker.id}`, {
            method: 'DELETE',
            headers: {
                'Authorization': 'Bearer tu-token'
            }
        });
        
        if (!response.ok) {
            console.error('Error al eliminar marcador:', await response.text());
        }
    } catch (error) {
        console.error('Error de red:', error);
    }
    */
}