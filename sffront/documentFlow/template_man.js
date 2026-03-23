// Variables globales
let currentTemplates = [];
let currentPage = 1;
let pageSize = 9;
let totalTemplates = 0;
let currentFilter = 'all';
let currentSearch = '';
let selectedTemplateId = null;
let currentIADocumentsType = 'templates';
let previewContainer, listContainer, renderedContent;

// Inicialización
document.addEventListener('DOMContentLoaded', () => {
    // Inicializar los contenedores cuando el HTML ya existe
    previewContainer = document.getElementById('docx-preview-container');
    listContainer = document.getElementById('iaDocumentsList');
    renderedContent = document.getElementById('docx-rendered-content');

    // Cargar plantillas
    loadTemplates();
    
    // Configurar eventos
    setupEventListeners();
    
    // Configurar búsqueda
    setupSearch();
});

// Configurar event listeners
function setupEventListeners() {
    // Botón para crear plantilla
    document.getElementById('createTemplateBtn').addEventListener('click', openAIModal);
    document.getElementById("logoutBtn").addEventListener("click", logout)
    
    // Botones de filtro
    document.querySelectorAll('.filter-btn').forEach(btn => {
        btn.addEventListener('click', () => {
            const filter = btn.dataset.filter;
            setActiveFilter(filter);
        });
    });
    
    // Botón generar plantilla en modal
    document.getElementById('generateTemplateBtn').addEventListener('click', generateTemplate);
    
    // Botón regenerar en modal
    document.getElementById('regenerateBtn').addEventListener('click', regenerateTemplate);
    
    // Botón crear PDF
    document.getElementById('createPdfBtn').addEventListener('click', createPdfFromTemplate);
    
    // Botón guardar plantilla
    document.getElementById('saveTemplateBtn').addEventListener('click', saveTemplate);
    
    // Botón generar documento en modal de llenado
    document.getElementById('generateFillBtn').addEventListener('click', generateDocumentFromTemplate);
    
    // Botón regenerar en modal de llenado
    document.getElementById('regenerateFillBtn').addEventListener('click', regenerateFill);
    
    // Botón descargar PDF
    document.getElementById('downloadPdfBtn').addEventListener('click', downloadGeneratedPdf);
    // Botón para ver documentos IA
    document.getElementById('viewIADocumentsBtn').addEventListener('click', openIADocumentsModal);
    
    // Tabs en el modal de documentos IA
    document.querySelectorAll('.doc-tab').forEach(tab => {
        tab.addEventListener('click', () => {
            const type = tab.dataset.type;
            switchIADocumentsTab(type);
        });
    });
    
    // Cerrar modales
    document.querySelectorAll('.close-modal').forEach(closeBtn => {
        closeBtn.addEventListener('click', () => {
            document.querySelectorAll('.modal').forEach(modal => {
                modal.style.display = 'none';
            });
        });
    });
    
    // Cerrar modal al hacer clic fuera
    window.addEventListener('click', (event) => {
        document.querySelectorAll('.modal').forEach(modal => {
            if (event.target === modal) {
                modal.style.display = 'none';
            }
        });
    });
}

// Abrir modal de documentos IA
function openIADocumentsModal() {
    const modal = document.getElementById('iaDocumentsModal');
    modal.style.display = 'block';
    
    // Cargar documentos
    loadIADocuments(currentIADocumentsType);
}

// Cerrar modal de documentos IA
function closeIADocumentsModal() {
    document.getElementById('iaDocumentsModal').style.display = 'none';
}

// Cambiar tab en documentos IA
function switchIADocumentsTab(type) {
    currentIADocumentsType = type;
    
    // Actualizar clases de tabs
    document.querySelectorAll('.doc-tab').forEach(tab => {
        if (tab.dataset.type === type) {
            tab.classList.add('active');
        } else {
            tab.classList.remove('active');
        }
    });
    
    // Cargar documentos del tipo seleccionado
    loadIADocuments(type);
}

// Cargar documentos IA
async function loadIADocuments(type = 'templates') {
    const container = document.getElementById('iaDocumentsList');
    
    // Mostrar estado de carga
    container.innerHTML = `
        <div class="documents-loading">
            <div class="spinner mx-auto mb-4"></div>
            <p class="text-secondary">Cargando ${type === 'templates' ? 'plantillas' : 'documentos'}...</p>
        </div>
    `;
    
    try {
        // Hacer fetch para obtener la lista de documentos
        // Esto depende de cómo esté configurado tu servidor
        // Necesitarías un endpoint que liste los archivos en el directorio
        const response = await fetch(`/iatemplates?type=${type}`);
        
        if (!response.ok) {
            // Si no hay endpoint específico, intentar cargar desde /iatemplates
            // mostrando un mensaje de que necesitamos listar los archivos
            container.innerHTML = `
                <div class="documents-error">
                    <span class="material-symbols-outlined text-3xl mb-2">error</span>
                    <p>No se pudo cargar la lista de documentos.</p>
                </div>
            `;
            return;
        }
        
        const documents = await response.json();
        renderIADocuments(documents, type);
        
    } catch (error) {
        console.error('Error al cargar documentos IA:', error);
        container.innerHTML = `
            <div class="documents-error">
                <span class="material-symbols-outlined text-3xl mb-2">error</span>
                <p>Error al cargar los documentos: ${error.message}</p>
                <button onclick="loadIADocuments('${type}')" class="mt-4 px-4 py-2 bg-white/10 hover:bg-white/20 rounded-lg">
                    Reintentar
                </button>
            </div>
        `;
    }
}

// Renderizar documentos IA
function renderIADocuments(documents, type) {
    const container = document.getElementById('iaDocumentsList');
    
    if (!documents || documents.length === 0) {
        container.innerHTML = `
            <div class="empty-documents">
                <span class="material-symbols-outlined text-4xl mb-2">folder_open</span>
                <p>No hay ${type === 'templates' ? 'plantillas' : 'documentos'} generados por IA.</p>
                <p class="text-sm mt-2">Crea tu primer documento usando el botón "Nueva Plantilla con IA"</p>
            </div>
        `;
        return;
    }
    
    container.innerHTML = '';
    
    documents.forEach(doc => {
        const docElement = createIADocumentElement(doc, type);
        container.appendChild(docElement);
    });
}


// Crear elemento de documento IA
function createIADocumentElement(doc, type) {
    const div = document.createElement('div');
    div.className = 'document-item';
    
    // Determinar icono y tipo
    const isTemplate = type === 'templates' || doc.name.endsWith('.docx') || doc.type === 'docx';
    const icon = isTemplate ? 'description' : 'picture_as_pdf';
    const fileType = isTemplate ? 'Template' : 'PDF';
    const fileExt = isTemplate ? '.docx' : '.pdf';
    
    // Formatear fecha
    const createdDate = doc.created ? 
        new Date(doc.created).toLocaleDateString('es-MX', {
            day: 'numeric',
            month: 'short',
            year: 'numeric',
            hour: '2-digit',
            minute: '2-digit'
        }) : 'Fecha desconocida';
    
    // Formatear tamaño
    const sizeFormatted = doc.size ? formatBytes(doc.size) : 'Tamaño desconocido';
    
    // Nombre del archivo sin extensión para mostrar
    const displayName = doc.name ? 
        doc.name.replace(/\.[^/.]+$/, '') : 
        'Documento sin nombre';
    
    div.innerHTML = `
        <div class="flex items-center flex-1">
            <span class="material-symbols-outlined document-icon">${icon}</span>
            <div class="document-info">
                <div class="document-name">${displayName}</div>
                <div class="document-meta">
                    <span>${createdDate}</span>
                    <span>${sizeFormatted}</span>
                    <span class="doc-type-badge ${isTemplate ? 'doc-template' : 'doc-pdf'}">
                        ${fileType}
                    </span>
                </div>
            </div>
        </div>
        <div class="doc-actions">
            <button class="px-3 py-1 bg-white/5 hover:bg-white/10 text-secondary rounded text-sm transition-colors view-doc-btn" 
                    data-path="${doc.path}" data-name="${doc.name}" data-type="${type}">
                <span class="material-symbols-outlined text-sm align-middle">visibility</span>
                Ver
            </button>
            <button class="px-3 py-1 bg-blue-600/20 hover:bg-blue-600/30 text-blue-400 rounded text-sm transition-colors download-doc-btn"
                    data-path="${doc.path}" data-name="${doc.name || displayName + fileExt}" data-type="${type}">
                <span class="material-symbols-outlined text-sm align-middle">download</span>
                Descargar
            </button>
        </div>
    `;
    
    // Agregar event listeners a los botones
    const viewBtn = div.querySelector('.view-doc-btn');
    const downloadBtn = div.querySelector('.download-doc-btn');
    
    viewBtn.addEventListener('click', (e) => {
        e.stopPropagation();
        const path = e.currentTarget.dataset.path;
        const type = e.currentTarget.dataset.type;
        viewIADocument(path, type);
    });
    
    downloadBtn.addEventListener('click', (e) => {
        e.stopPropagation();
        const path = e.currentTarget.dataset.path;
        const name = e.currentTarget.dataset.name;
        const type = e.currentTarget.dataset.type;
        downloadIADocument(path, name, type);
    });
    
    return div;
}


// Descargar documento IA
function downloadIADocument(path, name, type) {
    // path ya es el nombre del archivo según el JSON del back
    const url = `/iatemplates/${encodeURIComponent(path)}?type=${type}&download=1`;
    
    const a = document.createElement('a');
    a.href = url;
    a.download = name; // Sugerencia de nombre para el navegador
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);
}
// Refrescar documentos IA
function refreshIADocuments() {
    loadIADocuments(currentIADocumentsType);
}

// Ver documento IA
async function viewIADocument(path, type) {
    // Ocultar el modal de documentos IA temporalmente
    document.getElementById('iaDocumentsModal').style.display = 'none';
    
    // Asegurarnos de que los elementos existen
    if (!previewContainer) previewContainer = document.getElementById('docx-preview-container');
    if (!renderedContent) renderedContent = document.getElementById('docx-rendered-content');
    
    const url = `/iatemplates/${encodeURIComponent(path)}?type=${type}`;
    
    if (path.endsWith('.pdf')) {
        window.open(url, '_blank');
        return;
    }
    
    if (path.endsWith('.docx')) {
        try {
            showMessage('Cargando vista previa...', 'info');
            
            // Mostrar el contenedor de vista previa
            previewContainer.style.display = 'block';
            
            // Agregar overlay para cerrar
            const overlay = document.createElement('div');
            overlay.id = 'preview-overlay';
            overlay.style.cssText = `
                position: fixed;
                top: 0;
                left: 0;
                width: 100%;
                height: 100%;
                background: rgba(0,0,0,0);
                z-index: 9999;
                display: flex;
                align-items: center;
                justify-content: center;
            `;
            
            // Configurar el contenedor de vista previa
            previewContainer.style.cssText = `
                display: block !important;
                background: rgba(0,0,0,0);
                padding: 20px;
                border-radius: 8px;
                max-height: 90vh;
                width: 90%;
                max-width: 1000px;
                overflow-y: auto;
                z-index: 10000;
                position: relative;
                color: #000;
            `;
            
            renderedContent.innerHTML = `
                <div class="text-center py-8">
                    <div class="spinner mx-auto mb-4"></div>
                    <p>Cargando documento...</p>
                </div>
            `;
            
            overlay.appendChild(previewContainer);
            document.body.appendChild(overlay);
            
            const response = await fetch(url);
            
            if (!response.ok) {
                const errorText = await response.text();
                throw new Error(`Error al cargar el documento: ${response.status} - ${errorText}`);
            }
            
            const arrayBuffer = await response.arrayBuffer();
            
            // Verificar que el arrayBuffer tenga datos
            if (!arrayBuffer || arrayBuffer.byteLength === 0) {
                throw new Error('El documento está vacío o no se pudo cargar');
            }
            
            // Limpiar contenido anterior
            renderedContent.innerHTML = '';
            
            // Verificar que la librería docx esté disponible
            if (!window.docx) {
                throw new Error('La librería de vista previa no está disponible');
            }
            
            // Renderizar el documento
            await window.docx.renderAsync(arrayBuffer, renderedContent, null, {
                className: "docx", // className for the document container
                inWrapper: true, // enable wrapping the document in a div
                ignoreWidth: false,
                ignoreHeight: false,
                ignoreFonts: false,
                breakPages: true,
                ignoreLastRenderedPageBreak: true,
                experimental: true,
            });
            
            showMessage('Documento cargado correctamente', 'success');
            
        } catch (error) {
            console.error('Error al previsualizar:', error);
            
            // Mostrar error en el contenedor
            if (renderedContent) {
                renderedContent.innerHTML = `
                    <div class="text-center py-12 text-red-600">
                        <span class="material-symbols-outlined text-4xl mb-4">error</span>
                        <p class="text-lg mb-2">Error al cargar el documento</p>
                        <p class="text-sm mb-4">${error.message}</p>
                        <button onclick="closePreview()" class="px-4 py-2 bg-blue-600 text-white rounded">
                            Volver
                        </button>
                    </div>
                `;
            }
            
            showMessage(`Error: ${error.message}`, 'error');
        }
    }
}

// Función para volver corregida
function closePreview() {
    // Cerrar el overlay de vista previa
    const overlay = document.getElementById('preview-overlay');
    overlay.remove();
    
    
    // Mostrar nuevamente el modal de documentos IA
    document.getElementById('iaDocumentsModal').style.display = 'block';
    
    // Resetear el previewContainer
    if (previewContainer) {
        previewContainer.style.cssText = 'display: none;';
        document.body.appendChild(previewContainer);
    }
    
    if (renderedContent) {
        renderedContent.innerHTML = '';
    }
}


// Configurar búsqueda
function setupSearch() {
    const searchInput = document.getElementById('searchTemplates');
    let searchTimeout;
    
    searchInput.addEventListener('input', (e) => {
        clearTimeout(searchTimeout);
        searchTimeout = setTimeout(() => {
            currentSearch = e.target.value.trim();
            currentPage = 1;
            loadTemplates();
        }, 500);
    });
}

// Establecer filtro activo
function setActiveFilter(filter) {
    currentFilter = filter;
    currentPage = 1;
    
    // Actualizar clases de botones
    document.querySelectorAll('.filter-btn').forEach(btn => {
        if (btn.dataset.filter === filter) {
            btn.classList.add('bg-blue-500/20', 'text-blue-400');
            btn.classList.remove('bg-white/5', 'text-secondary');
        } else {
            btn.classList.remove('bg-blue-500/20', 'text-blue-400');
            btn.classList.add('bg-white/5', 'text-secondary');
        }
    });
    
    loadTemplates();
}

// Cargar plantillas
async function loadTemplates() {
    const grid = document.getElementById('templatesGrid');
    grid.innerHTML = `
        <div class="col-span-3 text-center py-12 text-secondary">
            <div class="spinner mx-auto mb-4"></div>
            <p>Cargando plantillas...</p>
        </div>
    `;
    
    try {
        const requestBody = {
            page: currentPage,
            page_size: pageSize,
            order_by: 'createdAtDoc',
            order_dir: 'ASC',
            type: 'docx'
        };
        const response = await fetch('/statusDocs', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(requestBody)
        });
        
        if (!response.ok) {
            if (resp.status === 401) {  
                window.location.href = '/login';
            }
            throw new Error(`Error HTTP: ${response.status}`);
        }
        
        const data = await response.json();
        currentTemplates = currentTemplates = Object.entries(data.data || {}).map(
            ([id, template]) => ({id, ...template})
        );
        totalTemplates = data.total || currentTemplates.length;
        
        renderTemplates();
        updateStats();
        renderPagination();
        
    } catch (error) {
        console.error('Error al cargar plantillas:', error);
        grid.innerHTML = `
            <div class="col-span-3 text-center py-12 text-red-400">
                <span class="material-symbols-outlined text-4xl mb-4">error</span>
                <p>Error al cargar las plantillas.</p>
                <button onclick="loadTemplates()" class="mt-4 px-4 py-2 bg-white/10 hover:bg-white/20 rounded-lg">
                    Reintentar
                </button>
            </div>
        `;
    }
}

// Renderizar plantillas en grid
function renderTemplates() {
    const grid = document.getElementById('templatesGrid');
    
    if (currentTemplates.length === 0) {
        grid.innerHTML = `
            <div class="col-span-3 text-center py-12 text-secondary">
                <span class="material-symbols-outlined text-4xl mb-4">description</span>
                <p class="text-xl mb-2">No hay plantillas aún</p>
                <p class="mb-6">Crea tu primera plantilla con IA</p>
                <button onclick="openAIModal()" class="btn-ai">
                    <span class="material-symbols-outlined">add</span>
                    Crear Primera Plantilla
                </button>
            </div>
        `;
        return;
    }
    
    grid.innerHTML = '';
    
    currentTemplates.forEach(template => {
        const templateCard = createTemplateCard(template);
        grid.appendChild(templateCard);
    });
}

// Crear tarjeta de plantilla
function createTemplateCard(template) {
    const card = document.createElement('div');
    card.className = 'glass-card rounded-xl overflow-hidden fade-in';
    
    // Formatear fecha
    const createdDate = template.createdAtDoc ? 
        new Date(template.createdAtDoc).toLocaleDateString('es-MX', {
            day: 'numeric',
            month: 'short',
            year: 'numeric'
        }) : 'Fecha no disponible';
    
    // Tamaño formateado
    const sizeFormatted = formatBytes(template.sizeB || 0);
    
    card.innerHTML = `
        <div class="p-6">
            <div class="flex justify-between items-start mb-4">
                <div class="flex-1">
                    <h3 class="text-lg font-bold text-primary mb-1 truncate">${template.documentName || 'Plantilla sin nombre'}</h3>
                    <p class="text-secondary text-sm">${template.documentExt || 'Sin extensión'} • ${sizeFormatted}</p>
                </div>
                <div>
                    ${template.authUseStatus === "1"? '<span class="status-badge status-public">Pública</span>' : 
                        template.authUseStatus === "2" ? '<span class="status-badge status-private">Privada</span>' : 
                      '<span class="status-badge status-inactive">Inactiva</span>'}
                </div>
            </div>
            
            <div class="mb-4">
                <p class="text-secondary text-sm line-clamp-3">
                    ${template.abstractDoc || 'Descripción no disponible'}
                </p>
            </div>
            
            <div class="flex items-center justify-between text-sm text-secondary mb-6">
                <span>Creada: ${createdDate}</span>
            </div>
            
            <div class="flex gap-2">
                <button class="flex-1 px-3 py-2 bg-white/5 hover:bg-white/10 text-secondary rounded-lg transition-colors flex items-center justify-center gap-2 fill-template-btn" data-id="${template.id}">
                    <span class="material-symbols-outlined text-sm">edit_document</span>
                    Llenar
                </button>
                <button class="px-3 py-2 bg-red-600/20 hover:bg-red-600/30 text-red-400 rounded-lg transition-colors flex items-center justify-center gap-1 convert-pdf-btn" data-id="${template.id}" title="Convertir a PDF">
                    <span class="material-symbols-outlined text-sm">picture_as_pdf</span>
                    PDF
                </button>
                <button class="px-3 py-2 bg-blue-600/20 hover:bg-blue-600/30 text-blue-400 rounded-lg transition-colors flex items-center justify-center gap-1" onclick="downloadTemplate('${template.documentPath}', '${template.documentName+"."+template.documentExt}', '${template.id}')">
                    <span class="material-symbols-outlined text-sm">download</span>
                </button>
            </div>
        </div>
    `;
    
    // Agregar event listener al botón de llenar
    const fillBtn = card.querySelector('.fill-template-btn');
    fillBtn.addEventListener('click', (e) => {
        e.stopPropagation();
        openFillModal(e.currentTarget.dataset.id);
    });

    // Agregar event listener al botón de convertir a PDF
    const pdfBtn = card.querySelector('.convert-pdf-btn');
    pdfBtn.addEventListener('click', (e) => {
        e.stopPropagation();
        convertTemplateToPdf(e.currentTarget.dataset.id);
    });

    return card;
}

// Actualizar estadísticas
function updateStats() {
    // En una implementación real, esto vendría de la API
    // Por ahora, calculamos de los datos locales
    
    const total = totalTemplates;
    const published = currentTemplates.filter(t => t.authUseStatus === '1').length;
    const priv = currentTemplates.filter(t => !t.authUseStatus || t.authUseStatus === '2').length;
    
    document.getElementById('totalTemplates').textContent = total;
    document.getElementById('publishedTemplates').textContent = published;
    document.getElementById('privateTemplates').textContent = priv;
    
}

// Renderizar paginación
function renderPagination() {
    const container = document.getElementById('paginationContainer');
    container.innerHTML = '';
    
    if (totalTemplates <= pageSize) return;
    
    const totalPages = Math.ceil(totalTemplates / pageSize);
    
    const createButton = (text, page, disabled = false, active = false) => {
        const btn = document.createElement('button');
        btn.textContent = text;
        btn.className = `
            px-3 py-1 rounded-lg text-sm font-medium transition
            ${active ? 'bg-blue-600 text-white' : 'bg-white/10 text-secondary hover:bg-white/20'}
            ${disabled ? 'opacity-50 cursor-not-allowed' : ''}
        `;
        if (!disabled && !active) {
            btn.addEventListener('click', () => {
                currentPage = page;
                loadTemplates();
            });
        }
        return btn;
    };
    
    // Botón anterior
    container.appendChild(createButton('←', currentPage - 1, currentPage === 1));
    
    // Páginas
    for (let i = 1; i <= totalPages; i++) {
        if (i === 1 || i === totalPages || (i >= currentPage - 1 && i <= currentPage + 1)) {
            container.appendChild(createButton(i, i, false, i === currentPage));
        } else if (i === currentPage - 2 || i === currentPage + 2) {
            const dots = document.createElement('span');
            dots.textContent = '...';
            dots.className = 'text-secondary px-2';
            container.appendChild(dots);
        }
    }
    
    // Botón siguiente
    container.appendChild(createButton('→', currentPage + 1, currentPage === totalPages));
}

// Abrir modal de IA
function openAIModal() {
    const modal = document.getElementById('aiTemplateModal');
    modal.style.display = 'block';
    
    // Resetear campos
    document.getElementById('templateName').value = '';
    document.getElementById('aiPrompt').value = '';
    document.getElementById('aiResponseContainer').classList.add('hidden');
    document.getElementById('aiLoading').classList.add('hidden');
    
    // Enfocar el textarea
    setTimeout(() => {
        document.getElementById('aiPrompt').focus();
    }, 100);
}

// Cerrar modal de IA
function closeAIModal() {
    document.getElementById('aiTemplateModal').style.display = 'none';
}

// Generar plantilla con IA
async function generateTemplate() {
    const userPrompt = document.getElementById('aiPrompt').value.trim();
    const templateName = document.getElementById('templateName').value.trim();

    if (!userPrompt) {
        showMessage('Por favor, describe el documento que quieres crear.', 'error');
        return;
    }

    const fullPrompt = templateName
        ? `Nombre del documento: ${templateName}. ${userPrompt}`
        : userPrompt;

    document.getElementById('aiLoading').classList.remove('hidden');
    document.getElementById('generateTemplateBtn').disabled = true;

    try {
        const response = await fetch('/interactDoc', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                type: 'agent',
                idDocument: '',
                prompt: fullPrompt,
            })
        });

        if (response.status === 401) { window.location.href = '/login'; return; }
        if (!response.ok) throw new Error(`Error HTTP: ${response.status}`);

        const data = await response.json();

        document.getElementById('aiLoading').classList.add('hidden');
        document.getElementById('aiResponseContainer').classList.remove('hidden');
        document.getElementById('aiResponse').textContent = data.message || 'Plantilla generada exitosamente.';

        // Refrescar lista de plantillas para que aparezca la recién creada
        loadTemplates();

    } catch (error) {
        console.error('Error al generar plantilla:', error);
        document.getElementById('aiLoading').classList.add('hidden');
        showMessage('Error al generar la plantilla. Intenta nuevamente.', 'error');
    } finally {
        document.getElementById('generateTemplateBtn').disabled = false;
    }
}

// Regenerar plantilla
async function regenerateTemplate() {
    const prompt = document.getElementById('aiPrompt').value.trim();

    if (!prompt) {
        showMessage('Por favor, modifica el prompt antes de regenerar.', 'error');
        return;
    }

    document.getElementById('aiLoading').classList.remove('hidden');

    try {
        const response = await fetch('/interactDoc', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ type: 'agent', idDocument: '', prompt })
        });
        if (response.status === 401) { window.location.href = '/login'; return; }
        if (!response.ok) throw new Error(`Error HTTP: ${response.status}`);
        const data = await response.json();
        document.getElementById('aiLoading').classList.add('hidden');
        document.getElementById('aiResponse').textContent = data.message;
        loadTemplates();
    } catch (error) {
        console.error('Error al regenerar:', error);
        document.getElementById('aiLoading').classList.add('hidden');
        showMessage('Error al regenerar la plantilla.', 'error');
    }
}

// Cerrar modal de creación y ver plantillas (la plantilla ya fue guardada por el agente)
async function createPdfFromTemplate() {
    closeAIModal();
    showMessage('La plantilla ya fue guardada. Usa el botón "→ PDF" en la tarjeta para convertirla.', 'info');
}

// Convertir plantilla existente a PDF
async function convertTemplateToPdf(templateId) {
    showMessage('Convirtiendo a PDF…', 'info');
    try {
        const response = await fetch('/interactDoc', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                type: 'agent',
                idDocument: templateId,
                prompt: 'Convierte esta plantilla a PDF sin modificar su contenido ni los placeholders.',
            })
        });
        if (response.status === 401) { window.location.href = '/login'; return; }
        if (!response.ok) throw new Error(`Error HTTP: ${response.status}`);
        const data = await response.json();
        showMessage(data.message || 'PDF generado exitosamente.', 'success');
    } catch (error) {
        console.error('Error al convertir a PDF:', error);
        showMessage(`Error: ${error.message}`, 'error');
    }
}

// Subir PDF generado
async function uploadGeneratedPdf(pdfContent, filename) {
    try {
        // Convertir base64 a blob si es necesario
        let blob;
        if (pdfContent.startsWith('data:')) {
            // Es data URL
            const response = await fetch(pdfContent);
            blob = await response.blob();
        } else {
            // Asumir que es base64
            const byteCharacters = atob(pdfContent);
            const byteNumbers = new Array(byteCharacters.length);
            for (let i = 0; i < byteCharacters.length; i++) {
                byteNumbers[i] = byteCharacters.charCodeAt(i);
            }
            const byteArray = new Uint8Array(byteNumbers);
            blob = new Blob([byteArray], { type: 'application/pdf' });
        }
        
        // Crear FormData
        const formData = new FormData();
        formData.append('document', blob, filename);
        formData.append('documentType', 'template');
        
        // Subir a /uploadDocs
        const response = await fetch('/uploadDocs', {
            method: 'POST',
            body: formData
        });
        
        if (!response.ok) {
            if (resp.status === 401) {  
                window.location.href = '/login';
            }
            throw new Error(`Error HTTP: ${response.status}`);
        }
        
        const result = await response.json();
        showMessage('Documento subido exitosamente.', 'success');
        
    } catch (error) {
        console.error('Error al subir PDF:', error);
        throw error;
    }
}

// La plantilla ya fue guardada automáticamente por el agente al generarla
function saveTemplate() {
    closeAIModal();
    showMessage('Plantilla guardada. Puedes verla en la lista.', 'success');
    loadTemplates();
}

// Abrir modal para llenar plantilla existente
function openFillModal(templateId) {
    const modal = document.getElementById('fillTemplateModal');
    const template = currentTemplates.find(t => t.id === templateId || t.documentHash === templateId);
    
    if (!template) {
        showMessage('Plantilla no encontrada.', 'error');
        return;
    }
    
    selectedTemplateId = templateId;
    
    // Configurar modal
    document.getElementById('fillModalTitle').textContent = `Llenar: ${template.documentName}`;
    document.getElementById('selectedTemplateName').textContent = template.documentName;
    document.getElementById('selectedTemplateDescription').textContent = template.abstractDoc || 'Sin descripción';
    document.getElementById('fillPrompt').value = '';
    document.getElementById('fillResponseContainer').classList.add('hidden');
    document.getElementById('fillLoading').classList.add('hidden');
    
    modal.style.display = 'block';
    
    setTimeout(() => {
        document.getElementById('fillPrompt').focus();
    }, 100);
}

// Cerrar modal de llenado
function closeFillModal() {
    document.getElementById('fillTemplateModal').style.display = 'none';
}

// Generar documento desde plantilla (llena placeholders y convierte a PDF)
async function generateDocumentFromTemplate() {
    const prompt = document.getElementById('fillPrompt').value.trim();

    if (!prompt) {
        showMessage('Por favor, proporciona instrucciones para llenar la plantilla.', 'error');
        return;
    }

    document.getElementById('fillLoading').classList.remove('hidden');
    document.getElementById('generateFillBtn').disabled = true;

    try {
        const response = await fetch('/interactDoc', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                type: 'agent',
                idDocument: selectedTemplateId,
                prompt: prompt + ' Rellena los placeholders, guarda el documento y conviértelo a PDF.',
            })
        });

        if (response.status === 401) { window.location.href = '/login'; return; }
        if (!response.ok) throw new Error(`Error HTTP: ${response.status}`);

        const data = await response.json();

        document.getElementById('fillLoading').classList.add('hidden');
        document.getElementById('fillResponseContainer').classList.remove('hidden');
        document.getElementById('fillResponse').textContent = data.message || 'Documento generado.';

    } catch (error) {
        console.error('Error al generar documento:', error);
        document.getElementById('fillLoading').classList.add('hidden');
        showMessage('Error al generar el documento.', 'error');
    } finally {
        document.getElementById('generateFillBtn').disabled = false;
    }
}

// Regenerar llenado
async function regenerateFill() {
    const prompt = document.getElementById('fillPrompt').value.trim();

    if (!prompt) {
        showMessage('Por favor, modifica las instrucciones.', 'error');
        return;
    }

    document.getElementById('fillLoading').classList.remove('hidden');

    try {
        const response = await fetch('/interactDoc', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                type: 'agent',
                idDocument: selectedTemplateId,
                prompt: prompt + ' Rellena los placeholders, guarda el documento y conviértelo a PDF.',
            })
        });
        if (response.status === 401) { window.location.href = '/login'; return; }
        if (!response.ok) throw new Error(`Error HTTP: ${response.status}`);
        const data = await response.json();
        document.getElementById('fillLoading').classList.add('hidden');
        document.getElementById('fillResponse').textContent = data.message;
    } catch (error) {
        console.error('Error al regenerar:', error);
        document.getElementById('fillLoading').classList.add('hidden');
        showMessage('Error al regenerar el documento.', 'error');
    }
}

// Descargar PDF generado — cierra el modal y muestra el resultado en Borradores
function downloadGeneratedPdf() {
    closeFillModal();
    showMessage('El PDF fue generado y guardado. Encuéntralo en la sección "Borradores".', 'success');
}


async function downloadTemplate(path, name, id) {
    try {
        const docs = {idDocs: [id], type: "template", paths: [path]}
        const response = await fetch('/downloadDoc', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(docs)
        });
        
        if (!response.ok) {
            if (resp.status === 401) {  
                window.location.href = '/login';
            }
            throw new Error(`Error HTTP: ${response.status}`);
        }
        
        const blob = await response.blob();
        const url = URL.createObjectURL(blob);
        
        const a = document.createElement('a');
        a.href = url;
        a.download = name || 'plantilla';
        document.body.appendChild(a);
        a.click();
        document.body.removeChild(a);
        URL.revokeObjectURL(url);
        
        showMessage('Plantilla descargada.', 'success');
        
    } catch (error) {
        console.error('Error al descargar plantilla:', error);
        showMessage('Error al descargar la plantilla.', 'error');
    }
}

function formatBytes(bytes) {
    bytes = Number(bytes);
    if (isNaN(bytes)) return "0 B";

    const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
    const i = bytes === 0 ? 0 : Math.floor(Math.log(bytes) / Math.log(1024));
    const value = (bytes / Math.pow(1024, i)).toFixed(2);
    return `${value} ${sizes[i]}`;
}

function showMessage(message, type = 'info') {
    // Crear elemento de mensaje
    const messageDiv = document.createElement('div');
    messageDiv.className = `fixed top-24 right-6 z-50 px-6 py-3 rounded-lg border text-sm font-medium flex items-center gap-2 animate-fade-in ${
        type === 'success' ? 'bg-green-500/10 border-green-500/50 text-green-400' :
        type === 'error' ? 'bg-red-500/10 border-red-500/50 text-red-400' :
        'bg-blue-500/10 border-blue-500/50 text-blue-400'
    }`;
    
    messageDiv.innerHTML = `
        <span class="material-symbols-outlined text-base">
            ${type === 'success' ? 'check_circle' : type === 'error' ? 'error' : 'info'}
        </span>
        ${message}
    `;
    
    // Agregar al body
    document.body.appendChild(messageDiv);
    
    // Remover después de 5 segundos
    setTimeout(() => {
        messageDiv.remove();
    }, 5000);
}

async function logout() {
    var ok = await sfConfirm({title:'Cerrar sesión',message:'¿Seguro que deseas cerrar tu sesión actual?',type:'warn',confirmText:'Cerrar sesión',confirmClass:'sf-modal-btn-danger'}); if(ok){
        try {
            const response = await fetch('/logoutUser', {
                method: 'GET',
                headers: { 'Content-Type': 'application/json' }}),
                credentials: 'include'
            });

            const data = await response.json();

            if (response.ok) { 
                // BORRAR TOKEN EN SESSION STORAGE
                sessionStorage.removeItem('aut');
                window.location.href = '/login';
                
            } else {
                sessionStorage.removeItem('aut');
                throw new Error(data.message || 'Error al cerrar sesión');
            }

        } catch (err) {
            console.error(err);
        }
    }
}

// Exportar funciones al scope global
window.openAIModal = openAIModal;
window.closeAIModal = closeAIModal;
window.openFillModal = openFillModal;
window.closeFillModal = closeFillModal;
window.downloadTemplate = downloadTemplate;
window.convertTemplateToPdf = convertTemplateToPdf;
window.openIADocumentsModal = openIADocumentsModal;
window.closeIADocumentsModal = closeIADocumentsModal;
window.refreshIADocuments = refreshIADocuments;