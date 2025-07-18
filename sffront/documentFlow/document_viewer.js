// Variables para PDF.js
let pdfDoc = null,
    pageNum = 1,
    pageRendering = false,
    pageNumPending = null,
    scale = 1.0,
    pdfAspectRatio = 8.5 / 11; 

// Arreglo de marcadores (ahora serán cajas de texto)
let markers = [];
// Elementos del DOM
let pdfCanvas, textElementsContainer;

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
        const tableBody = document.querySelector('.glass-card table tbody');
        if (tableBody) {
            tableBody.innerHTML = ''; // Limpiar datos de ejemplo
            
            documentData.reviewers.forEach(reviewer => {
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
        
        setupDraggables();
        setupDropZone();
        loadExistingMarkers();
    }
});

function initPDFElements() {
    // Obtener elementos del DOM
    pdfCanvas = document.getElementById('pdf-canvas');
    textElementsContainer = document.getElementById('text-elements-container');
    pdfCtx = pdfCanvas.getContext('2d');
    
    // Configurar eventos de los controles
    document.getElementById('prev-page').addEventListener('click', onPrevPage);
    document.getElementById('next-page').addEventListener('click', onNextPage);
    document.getElementById('view-full').addEventListener('click', openFullPDF);
    
    // Redimensionar cuando cambia el tamaño de la ventana
    window.addEventListener('resize', function() {
        if (pdfDoc) {
            resizeCanvas();
            renderPage(pageNum);
        }
    });
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
        resizeCanvas();
        renderPage(1);
    }).catch(function(error) {
        console.error('Error al cargar PDF:', error);
    });
}

function resizeCanvas() {
    const container = document.getElementById('preview-container');
    const width = container.clientWidth;
    const height = container.clientHeight;
    
    // Mantener relación de aspecto carta (8.5x11 pulgadas)
    const aspectRatio = pdfAspectRatio;
    const canvasWidth = width;
    const canvasHeight = width/aspectRatio;
    
    // Si el alto calculado es mayor que el contenedor, ajustar
    if (canvasHeight > height) {
        const newHeight = height;
        const newWidth = height * aspectRatio;
        pdfCanvas.style.width = `${newWidth}px`;
        pdfCanvas.style.height = `${newHeight}px`;
    } else {
        pdfCanvas.style.width = `${canvasWidth}px`;
        pdfCanvas.style.height = `${canvasHeight}px`;
    }
    
    // Ajustar tamaño real del canvas
    pdfCanvas.width = pdfCanvas.clientWidth;
    pdfCanvas.height = pdfCanvas.clientHeight;
    
    // Ajustar el contenedor de elementos de texto para que coincida
    textElementsContainer.style.width = `${pdfCanvas.clientWidth}px`;
    textElementsContainer.style.height = `${pdfCanvas.clientHeight}px`;
}

function renderPage(num) {
    pageRendering = true;
    
    pdfDoc.getPage(num).then(function(page) {
        const viewport = page.getViewport({scale: scale});
        // Calcular y guardar el aspect ratio real
        pdfAspectRatio = viewport.width / viewport.height;

        // Ajustar el canvas al tamaño de la página
        pdfCanvas.height = viewport.height;
        pdfCanvas.width = viewport.width;
        
        // Ajustar el contenedor de texto al mismo tamaño
        textElementsContainer.style.height = `${viewport.height}px`;
        textElementsContainer.style.width = `${viewport.width}px`;
        
        // Renderizar página PDF
        const renderContext = {
            canvasContext: pdfCtx,
            viewport: viewport
        };
        
        const renderTask = page.render(renderContext);
        
        renderTask.promise.then(function() {
            pageRendering = false;
            if (pageNumPending !== null) {
                renderPage(pageNumPending);
                pageNumPending = null;
            }
            
            // Redibujar marcadores (cajas de texto)
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

// Dibujar marcador en canvas
function drawMarker(x, y, user) {
  const size = 30;
  
  // Dibujar círculo
  ctx.beginPath();
  ctx.arc(x, y, size / 2, 0, Math.PI * 2);
  ctx.fillStyle = 'rgba(59, 130, 246, 0.7)';
  ctx.fill();
  
  // Dibujar inicial
  ctx.font = 'bold 14px Arial';
  ctx.fillStyle = 'white';
  ctx.textAlign = 'center';
  ctx.textBaseline = 'middle';
  ctx.fillText(user.charAt(0), x, y);
}

// Configurar zona de drop
function setupDropZone() {
    pdfCanvas.addEventListener('dragover', (e) => {
        e.preventDefault();
        pdfCanvas.style.cursor = 'copy';
    });

    pdfCanvas.addEventListener('dragleave', () => {
        pdfCanvas.style.cursor = 'default';
    });

    pdfCanvas.addEventListener('drop', (e) => {
        e.preventDefault();
        pdfCanvas.style.cursor = 'default';
        
        if (!currentDraggedReviewer) return;
        
        // Calcular posición relativa al PDF
        const rect = pdfCanvas.getBoundingClientRect();
        const x = e.clientX - rect.left;
        const y = e.clientY - rect.top;
        
        // Crear caja de texto
        const textBox = createTextBox(currentDraggedReviewer.user, x, y);
        textElementsContainer.appendChild(textBox);
        
        // Guardar marcador
        updateMarkerPosition(textBox);
        
        currentDraggedReviewer = null;
    });
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



// Función para redibujar marcadores al cambiar de página
function redrawMarkers() {
    // Limpiar contenedor
    textElementsContainer.innerHTML = '';
    
    // Filtrar marcadores de la página actual
    const currentPageMarkers = markers.filter(marker => marker.page === pageNum);
    
    // Crear cajas de texto para cada marcador
    currentPageMarkers.forEach(marker => {
        const rect = pdfCanvas.getBoundingClientRect();
        const x = marker.x * rect.width;
        const y = marker.y * rect.height;
        
        const textBox = createTextBox(marker.user, x, y);
        textElementsContainer.appendChild(textBox);
        
        // Aplicar tamaño si existe
        if (marker.width && marker.height) {
            textBox.style.width = `${marker.width * rect.width}px`;
            textBox.style.height = `${marker.height * rect.height}px`;
        }
    });
}

// Cargar marcadores existentes
function loadExistingMarkers() {
    try {
        // Simulación de marcadores existentes (en una aplicación real, esto vendría de una API)
        const existingMarkers = []; 
        
        // Limpiar contenedor
        textElementsContainer.innerHTML = '';
        
        // Filtrar marcadores de la página actual
        const currentPageMarkers = existingMarkers.filter(m => m.page === pageNum);
        
        // Crear cajas de texto para cada marcador
        currentPageMarkers.forEach(marker => {
            const rect = pdfCanvas.getBoundingClientRect();
            const x = marker.x * rect.width;
            const y = marker.y * rect.height;
            
            const textBox = createTextBox(marker.user, x, y);
            textElementsContainer.appendChild(textBox);
            
            // Aplicar tamaño si existe
            if (marker.width && marker.height) {
                textBox.style.width = `${marker.width * rect.width}px`;
                textBox.style.height = `${marker.height * rect.height}px`;
            }
        });
        
        // Agregar a la lista de marcadores
        markers.push(...currentPageMarkers);
        
    } catch (error) {
        console.error('Error cargando marcadores:', error);
    }
}


// Enviar datos a API

async function sendMarkerToAPI(markerData) {
  console.log(JSON.stringify(markerData))
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
  }*/
}

// Crear caja de texto para un revisor
function createTextBox(user, x, y) {
    const textBox = document.createElement('div');
    textBox.className = 'text-box';
    textBox.textContent = user;
    textBox.dataset.user = user;
    textBox.style.left = `${x}px`;
    textBox.style.top = `${y}px`;
    
    makeDraggableResizable(textBox, textElementsContainer);
    
    return textBox;
}

// Hacer las cajas de texto arrastrables y redimensionables
function makeDraggableResizable(element, container) {
    let offsetX = 0, offsetY = 0, isDragging = false;

    element.addEventListener("mousedown", (e) => {
        // Ignorar si el clic fue en un área que permite redimensionar
        const style = window.getComputedStyle(element);
        const isResizing = style.resize !== "none" && (
            e.offsetX > element.clientWidth - 16 && e.offsetY > element.clientHeight - 16
        );
        if (isResizing) return;

        e.preventDefault();
        isDragging = true;
        offsetX = e.offsetX;
        offsetY = e.offsetY;
        element.style.zIndex = "10";
    });

    document.addEventListener("mousemove", (e) => {
        if (!isDragging) return;

        const containerRect = container.getBoundingClientRect();
        let newLeft = e.clientX - containerRect.left - offsetX;
        let newTop = e.clientY - containerRect.top - offsetY;

        newLeft = Math.max(0, Math.min(newLeft, container.clientWidth - element.offsetWidth));
        newTop = Math.max(0, Math.min(newTop, container.clientHeight - element.offsetHeight));

        element.style.left = newLeft + "px";
        element.style.top = newTop + "px";
        element.style.maxWidth = container.clientWidth + "px";
        element.style.maxHeight = container.clientHeight + "px";
        
        updateMarkerPosition(element);
    });

    document.addEventListener("mouseup", () => {
        if (isDragging) {
            isDragging = false;
            element.style.zIndex = "1";
            updateMarkerPosition(element);
        }
    });

    // Actualizar posición/tamaño también al terminar redimensionamiento manual
    const resizeObserver = new ResizeObserver(() => updateMarkerPosition(element));
    resizeObserver.observe(element);
}

function updateMarkerPosition(element) {
    const container = textElementsContainer;
    const user = element.dataset.user;
    const rect = container.getBoundingClientRect();
    
    const x = element.offsetLeft / rect.width;
    const y = element.offsetTop / rect.height;
    
    // Actualizar marcador existente o crear uno nuevo
    const existingMarker = markers.find(m => m.user === user && m.page === pageNum);
    if (existingMarker) {
        existingMarker.x = x;
        existingMarker.y = y;
        existingMarker.width = element.offsetWidth / rect.width;
        existingMarker.height = element.offsetHeight / rect.height;
    } else {
        markers.push({
            user,
            x,
            y,
            width: element.offsetWidth / rect.width,
            height: element.offsetHeight / rect.height,
            page: pageNum
        });
    }
    
    // Enviar a API
    sendMarkerToAPI({
        document: localStorage.getItem('currentDocumentFileName'),
        reviewer: user,
        position: { x, y, width: element.offsetWidth / rect.width, height: element.offsetHeight / rect.height },
        page: pageNum
    });
}

// Modificar queueRenderPage para redibujar marcadores
function queueRenderPage(num) {
    if (pageRendering) {
        pageNumPending = num;
    } else {
        renderPage(num);
        redrawMarkers();
    }
}