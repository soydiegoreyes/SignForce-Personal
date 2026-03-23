// Variables para PDF.js
let pdfDoc = null,
    pageNum = 1,
    pageRendering = false,
    pageNumPending = null,
    scale = 1.0,
    pdfAspectRatio = 8.5 / 11;

// Marcadores y control
let markers = [];
let markerColors = ['#1D2891aa', '#861d7daa', '#5c7500aa', '#004a61aa', '#160064aa', '#611500aa'];
let currentDraggedReviewer = null;
let selectedTextBox = null;
let isDragging = false, isResizing = false;
let dragOffsetX = 0, dragOffsetY = 0;

// PDF real dimensions
let pdfRealWidth = 612, pdfRealHeight = 792;

// Configurar PDF.js worker
pdfjsLib.GlobalWorkerOptions.workerSrc = 'https://cdnjs.cloudflare.com/ajax/libs/pdf.js/3.4.120/pdf.worker.min.js';

// Control de documentos en sessionStorage
let inviteRequest = JSON.parse(sessionStorage.getItem("inviteRequest")) || {};
let docIds = Object.keys(inviteRequest);
let currentIndex = parseInt(sessionStorage.getItem("signIndex") || "0");
let currentDocId = docIds[currentIndex];
let currentDocRev = inviteRequest[currentDocId];

document.addEventListener('DOMContentLoaded', async function() {
    if (!currentDocRev) {
        sfAlert('No hay más documentos por firmar.');
        sessionStorage.removeItem("signIndex");
        return;
    }

    initPDFElements();

    const title = document.getElementById('docTitle');
    if (title) title.textContent = `${currentDocRev.document.documentName}.${currentDocRev.document.documentExt}`;

    // Cambiar texto del botón según documento actual
    const nextDocBtn = document.getElementById('nextDoc');
    if (nextDocBtn) {
        if (currentIndex + 1 < docIds.length) {
            nextDocBtn.textContent = "Guardar y continuar";
        } else {
            nextDocBtn.textContent = "Finalizar proceso";
        }
    }

    updateApprovalTable(currentDocRev.reviewers);
    setupCanvasInteractions();

    await loadPDFfromServer(currentDocRev.document.idDocument);
});

async function loadPDFfromServer(docId) {
    try {
        const docs = {idDocs: [docId], type: "uploaded"}
        const response = await fetch('/downloadDoc', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(docs)
        });

        if (!response.ok) {
            if (response.status === 401) {  
                window.location.href = '/login';
            }
            throw new Error(`HTTP ${response.status}`);
        }
        const blob = await response.blob();
        const arrayBuffer = await blob.arrayBuffer();
        const pdfData = new Uint8Array(arrayBuffer);

        const loadingTask = pdfjsLib.getDocument({ data: pdfData });
        const pdf = await loadingTask.promise;
        pdfDoc = pdf;

        const firstPage = await pdf.getPage(1);
        const viewport = firstPage.getViewport({ scale: 1.0 });
        pdfRealWidth = viewport.width;
        pdfRealHeight = viewport.height;
        pdfAspectRatio = pdfRealWidth / pdfRealHeight;

        resizeCanvases();
        renderPage(1);

    } catch (err) {
        console.error("Error al cargar PDF:", err);
        sfAlert('No se pudo cargar el documento.');
    }
}

// ==========================
// 📄 Renderizado PDF
// ==========================
function initPDFElements() {
    pdfCanvas = document.getElementById('pdf-canvas');
    overlayCanvas = document.getElementById('overlay-canvas');

    if (!pdfCanvas || !overlayCanvas) {
        console.warn("⚠️ No se encontró el canvas PDF u overlay en el DOM.");
        return;
    }

    pdfCtx = pdfCanvas.getContext('2d');
    overlayCtx = overlayCanvas.getContext('2d');

    // Agregar eventos solo si existen los elementos
    const prevBtn = document.getElementById('prev-page');
    const nextBtn = document.getElementById('next-page');
    const nextDocBtn = document.getElementById('nextDoc');

    if (prevBtn) prevBtn.addEventListener('click', onPrevPage);
    if (nextBtn) nextBtn.addEventListener('click', onNextPage);
    if (nextDocBtn) nextDocBtn.addEventListener('click', saveAndNextDoc);

    window.addEventListener('resize', () => {
        if (pdfDoc) {
            resizeCanvases();
            renderPage(pageNum);
        }
    });
}

function renderPage(num) {
    pageRendering = true;

    pdfDoc.getPage(num).then(function(page) {
        const container = document.getElementById('preview-container');
        const containerWidth = container.clientWidth;
        const containerHeight = container.clientHeight;

        const viewport = page.getViewport({ scale: 1 });
        const containerAspectRatio = containerWidth / containerHeight;
        let scale = (viewport.width / viewport.height > containerAspectRatio)
            ? containerWidth / viewport.width
            : containerHeight / viewport.height;

        const scaledViewport = page.getViewport({ scale });
        pdfCanvas.width = scaledViewport.width;
        pdfCanvas.height = scaledViewport.height;

        const renderTask = page.render({ canvasContext: pdfCtx, viewport: scaledViewport });
        renderTask.promise.then(() => {
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
    if (!pdfDoc || pageNum <= 1) return;
    pageNum--;
    queueRenderPage(pageNum);
}

function onNextPage() {
    if (!pdfDoc || pageNum >= pdfDoc.numPages) return;
    pageNum++;
    queueRenderPage(pageNum);
}
function queueRenderPage(num) { pageRendering ? pageNumPending = num : renderPage(num); }

function resizeCanvases() {
    const container = document.getElementById('preview-container');
    let w = container.clientWidth, h = w / pdfAspectRatio;
    if (h > container.clientHeight) { h = container.clientHeight; w = h * pdfAspectRatio; }
    [pdfCanvas, overlayCanvas].forEach(c => {
        c.width = w; c.height = h; c.style.width = `${w}px`; c.style.height = `${h}px`;
    });
}

// ==========================
// 👥 Firmantes
// ==========================
function updateApprovalTable(reviewers) {
    const body = document.querySelector('.glass-card table tbody');
    body.innerHTML = '';
    reviewers.forEach(rev => {
        const row = document.createElement('tr');
        row.draggable = true;
        row.innerHTML = `
            <td class="px-4 py-2 text-primary text-sm">${rev.user}</td>
            <td class="px-4 py-2 text-secondary text-sm">${rev.role === 1 ? "Firmante" : "Revisor"}</td>
        `;
        row.addEventListener('dragstart', e => {
            currentDraggedReviewer = { user: rev.user };
            e.dataTransfer.setData('text/plain', rev.user);
        });
        body.appendChild(row);
    });
}

// ==========================
// ✍️ Canvas Interactions
// ==========================
function setupCanvasInteractions() {
    overlayCanvas.addEventListener('dragover', e => { e.preventDefault(); overlayCanvas.style.cursor = 'copy'; });
    overlayCanvas.addEventListener('dragleave', () => overlayCanvas.style.cursor = 'default');

    overlayCanvas.addEventListener('drop', e => {
        e.preventDefault();
        overlayCanvas.style.cursor = 'default';
        if (!currentDraggedReviewer) return;
        const rect = overlayCanvas.getBoundingClientRect();
        const x = e.clientX - rect.left;
        const y = e.clientY - rect.top;
        const newMarker = {
            user: currentDraggedReviewer.user,
            x, y, width: 50, height: 50, page: pageNum,
            bgColor: markerColors[markers.length % markerColors.length]
        };
        markers.push(newMarker);
        redrawMarkers();
        currentDraggedReviewer = null;
    });

    overlayCanvas.addEventListener('mousedown', handleMouseDown);
    overlayCanvas.addEventListener('mousemove', handleMouseMove);
    overlayCanvas.addEventListener('mouseup', handleMouseUp);
    overlayCanvas.addEventListener('dblclick', handleDoubleClick);
}

function getMarkerAtPosition(x, y) {
    for (let i = markers.length - 1; i >= 0; i--) {
        const m = markers[i];
        if (m.page === pageNum && x >= m.x && x <= m.x + m.width && y >= m.y && y <= m.y + m.height) return m;
    }
    return null;
}

function handleMouseDown(e) {
    const r = overlayCanvas.getBoundingClientRect();
    const x = e.clientX - r.left, y = e.clientY - r.top;
    const m = getMarkerAtPosition(x, y);
    if (m) {
        selectedTextBox = m;
        const inResize = x >= m.x + m.width - 10 && y >= m.y + m.height - 10;
        if (inResize) { isResizing = true; overlayCanvas.style.cursor = 'nw-resize'; }
        else { isDragging = true; dragOffsetX = x - m.x; dragOffsetY = y - m.y; overlayCanvas.style.cursor = 'move'; }
    }
}

function handleMouseMove(e) {
    const r = overlayCanvas.getBoundingClientRect();
    const x = e.clientX - r.left, y = e.clientY - r.top;
    if (isDragging && selectedTextBox) {
        selectedTextBox.x = Math.max(0, Math.min(x - dragOffsetX, overlayCanvas.width - selectedTextBox.width));
        selectedTextBox.y = Math.max(0, Math.min(y - dragOffsetY, overlayCanvas.height - selectedTextBox.height));
        redrawMarkers();
    } else if (isResizing && selectedTextBox) {
        selectedTextBox.width = Math.max(20, x - selectedTextBox.x);
        selectedTextBox.height = Math.max(20, y - selectedTextBox.y);
        redrawMarkers();
    } else {
        const m = getMarkerAtPosition(x, y);
        if (m) {
            const resize = x >= m.x + m.width - 10 && y >= m.y + m.height - 10;
            overlayCanvas.style.cursor = resize ? 'nw-resize' : 'move';
        } else overlayCanvas.style.cursor = 'default';
    }
}

function handleMouseUp() { isDragging = isResizing = false; selectedTextBox = null; overlayCanvas.style.cursor = 'default'; }

function handleDoubleClick(e) {
    const r = overlayCanvas.getBoundingClientRect();
    const x = e.clientX - r.left, y = e.clientY - r.top;
    const m = getMarkerAtPosition(x, y);
    if (m) { markers = markers.filter(mm => mm !== m); redrawMarkers(); }
}

function redrawMarkers() {
    overlayCtx.clearRect(0, 0, overlayCanvas.width, overlayCanvas.height);
    markers.filter(m => m.page === pageNum).forEach(drawTextBox);
}

function drawTextBox(m) {
    const ctx = overlayCtx;
    ctx.fillStyle = m.bgColor;
    ctx.fillRect(m.x, m.y, m.width, m.height);
    ctx.strokeStyle = 'rgba(0,0,0,0.8)';
    ctx.lineWidth = 2;
    ctx.strokeRect(m.x, m.y, m.width, m.height);
    ctx.fillStyle = 'white';
    ctx.font = '14px Arial';
    ctx.textAlign = 'center';
    ctx.textBaseline = 'middle';
    ctx.fillText(m.user, m.x + m.width / 2, m.y + m.height / 2);
}

// ==========================
// 💾 Guardar y pasar al siguiente documento
// ==========================
function canvasToRealPDF(x, y, w, h) {
    const sx = pdfRealWidth / overlayCanvas.width, sy = pdfRealHeight / overlayCanvas.height;
    return { x: x * sx, y: y * sy, width: w * sx, height: h * sy };
}

async function saveAndNextDoc() {
    // Agrupamos los marcadores por usuario
    const groupedByUser = {};
    markers.forEach(m => {
        const pos = canvasToRealPDF(m.x, m.y, m.width, m.height);
        if (!groupedByUser[m.user]) groupedByUser[m.user] = [];
        groupedByUser[m.user].push({
            page: m.page,
            x: pos.x,
            y: pos.y,
            width: pos.width,
            height: pos.height
        });
    });

    // Insertamos las posiciones dentro de cada reviewer del documento actual
    if (currentDocRev.reviewers && Array.isArray(currentDocRev.reviewers)) {
        currentDocRev.reviewers.forEach(rev => {
            if (groupedByUser[rev.user]) {
                rev.positions = groupedByUser[rev.user];
            }
        });
    }

    // Actualizar el inviteRequest con los cambios
    inviteRequest[currentDocId] = currentDocRev;
    sessionStorage.setItem("inviteRequest", JSON.stringify(inviteRequest));

    console.log("✅ Posiciones guardadas para documento:", currentDocId);

    // Lógica para pasar al siguiente documento
    const nextDocBtn = document.getElementById('nextDoc');
    if (currentIndex + 1 < docIds.length) {
        sessionStorage.setItem("signIndex", currentIndex + 1);
        if (nextDocBtn) nextDocBtn.textContent = "Cargando siguiente documento...";
        nextDocBtn.disabled = true;
        window.location.href = "/addSignatures";
    } else {
        sessionStorage.removeItem("signIndex");
        if (nextDocBtn) {
            nextDocBtn.textContent = "Finalizando...";
            nextDocBtn.disabled = true;
        }
        
        console.log("✅ Paquete final listo para enviar:", inviteRequest);
        
        // Enviar todo el inviteRequest al backend
        try {
            const response = await fetch('/closeInvite', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(inviteRequest)
            });
            
            if (response.ok) {
                sfAlert('Proceso completado. Todos los documentos tienen posiciones de firma.');
                window.location.href = "/myfolders";
            } else {
                if (response.status === 401) {  
                    window.location.href = '/login';
                }
                sfAlert('Error al enviar las invitaciones.');
            }
        } catch (err) {
            console.error("Error al enviar invitaciones:", err);
            sfAlert('Error al enviar las invitaciones.');
        }
    }
}