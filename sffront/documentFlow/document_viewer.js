// Variables globales
let canvas, ctx;
let currentDraggedReviewer = null;
const markers = [];

// Cargar datos del documento al cargar la página
document.addEventListener('DOMContentLoaded', function() {
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
            const previewContainer = document.getElementById('preview-container');
            previewContainer.innerHTML = '';
            
            const wrapper = document.createElement('div');
            wrapper.style.position = 'relative';
            wrapper.style.width = '100%';
            wrapper.style.height = '100%';

            const iframe = document.createElement('iframe');
            iframe.src = `data:application/pdf;base64,${documentData.fileBase64}#page=1&zoom=80%`;
            iframe.style.width = '100%';
            iframe.style.height = '100%';
            iframe.style.border = 'none';

            // Crear el canvas
            const canvas = document.createElement('canvas');
            canvas.id = 'pdfCanvas';
            canvas.style.position = 'absolute';
            canvas.style.top = '0';
            canvas.style.left = '0';
            canvas.style.width = '90%';
            canvas.style.height = '80%';
            canvas.style.pointerEvents = 'none';
            canvas.style.zIndex = '5';

            const overlay = document.createElement('div');
            overlay.innerText = 'Ver completo';
            overlay.style.position = 'absolute';
            overlay.style.bottom = '10px';
            overlay.style.right = '10px';
            overlay.style.background = 'rgba(0,0,0,0.6)';
            overlay.style.color = '#fff';
            overlay.style.padding = '6px 10px';
            overlay.style.borderRadius = '8px';
            overlay.style.cursor = 'pointer';
            overlay.style.fontSize = '12px';
            overlay.style.zIndex = '10';

            overlay.onclick = () => {
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
                            <embed src="data:application/pdf;base64,${documentData.fileBase64}#page=1&zoom=100%">
                        </body>
                    </html>
                `);
            };

            wrapper.appendChild(iframe);
            wrapper.appendChild(canvas);
            wrapper.appendChild(overlay);
            previewContainer.appendChild(wrapper);
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
        // Nuevas inicializaciones
        initCanvas();
        setupDraggables();
        setupDropZone();
        
        // Cargar marcadores existentes
        loadExistingMarkers();
    }
});

document.getElementById('fileInput').addEventListener('change', function (event) {
  let fileURL = null;
  const previewContainer = document.getElementById('preview-container');

  const file = event.target.files[0];
  if (file && file.type === 'application/pdf') {
    document.getElementById('docName').value= file.name;
    const reader = new FileReader();
    reader.onload = function (e) {
      if (fileURL) URL.revokeObjectURL(fileURL);
      fileURL = URL.createObjectURL(file);

      previewContainer.innerHTML = '';

      // Wrapper relativo
      const wrapper = document.createElement('div');
      wrapper.style.position = 'relative';
      wrapper.style.width = '100%';
      wrapper.style.height = '100%';

      // Iframe con eventos habilitados
      const iframe = document.createElement('iframe');
      iframe.src = `${fileURL}#page=1&zoom=25%`;
      iframe.style.width = '100%';
      iframe.style.height = '100%';
      iframe.style.border = 'none';

      // Div flotante para click
      const overlay = document.createElement('div');
      overlay.innerText = 'Ver completo';
      overlay.style.position = 'absolute';
      overlay.style.bottom = '10px';
      overlay.style.right = '10px';
      overlay.style.background = 'rgba(0,0,0,0.6)';
      overlay.style.color = '#fff';
      overlay.style.padding = '6px 10px';
      overlay.style.borderRadius = '8px';
      overlay.style.cursor = 'pointer';
      overlay.style.fontSize = '12px';
      overlay.style.zIndex = '10';

      overlay.onclick = () => window.open(fileURL, '_blank');

      wrapper.appendChild(iframe);
      wrapper.appendChild(overlay);
      previewContainer.appendChild(wrapper);
    };
    reader.readAsArrayBuffer(file);
  }
});

// Inicializar canvas
function initCanvas() {
  canvas = document.getElementById('pdfCanvas');
  ctx = canvas.getContext('2d');
  canvas.width = canvas.clientWidth;
  canvas.height = canvas.clientHeight;
  canvas.style.pointerEvents = 'auto'; // Permitir interacción
}

// Hacer elementos arrastrables
function setupDraggables() {
  const table = document.querySelector('.glass-card table');
  const rows = table.querySelectorAll('tbody tr');
  
  rows.forEach(row => {
    row.draggable = true;
    
    row.addEventListener('dragstart', (e) => {
      const team = row.cells[0].textContent;
      const user = row.cells[3].textContent.split(',')[0].trim();
      currentDraggedReviewer = { team, user };
      
      e.dataTransfer.setData('text/plain', JSON.stringify({
        team,
        user
      }));
    });
  });
}

// Configurar zona de drop
function setupDropZone() {
  canvas.addEventListener('dragover', (e) => {
    e.preventDefault();
    canvas.style.cursor = 'copy';
  });

  canvas.addEventListener('dragleave', () => {
    canvas.style.cursor = 'default';
  });

  canvas.addEventListener('drop', (e) => {
    e.preventDefault();
    canvas.style.cursor = 'default';
    
    if (!currentDraggedReviewer) return;
    
    // Calcular posición relativa al PDF
    const rect = canvas.getBoundingClientRect();
    const x = e.clientX - rect.left;
    const y = e.clientY - rect.top;
    const normalizedX = x / rect.width;
    const normalizedY = y / rect.height;
    
    // Guardar marcador
    markers.push({
      ...currentDraggedReviewer,
      x: normalizedX,
      y: normalizedY
    });
    
    // Dibujar marcador
    drawMarker(x, y, currentDraggedReviewer.user);
    
    // Enviar a API
    sendMarkerToAPI({
      document: localStorage.getItem('currentDocumentFileName'),
      reviewer: currentDraggedReviewer.user,
      position: { x: normalizedX, y: normalizedY }
    });
    
    currentDraggedReviewer = null;
  });
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


// Cargar marcadores guardados

async function loadExistingMarkers() {
  try {
    const response = await fetch('https://tu-api.com/markers?document=' + 
                               encodeURIComponent(localStorage.getItem('currentDocumentFileName')));
    
    if (response.ok) {
      const existingMarkers = await response.json();
      
      existingMarkers.forEach(marker => {
        // Convertir coordenadas normalizadas a píxeles
        const rect = canvas.getBoundingClientRect();
        const x = marker.position.x * rect.width;
        const y = marker.position.y * rect.height;
        
        markers.push(marker);
        drawMarker(x, y, marker.reviewer);
      });
    }
  } catch (error) {
    console.error('Error cargando marcadores:', error);
  }
}