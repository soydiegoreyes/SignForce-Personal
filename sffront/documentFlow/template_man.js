// Variables globales
let currentTemplates = [];
let currentPage = 1;
let pageSize = 9;
let totalTemplates = 0;
let currentFilter = 'all';
let currentSearch = '';
let selectedTemplateId = null;

// Inicialización
document.addEventListener('DOMContentLoaded', () => {
    // Cargar plantillas
    loadTemplates();
    
    // Configurar eventos
    setupEventListeners();
    
    // Configurar tema
    setupThemeToggle();
    
    // Configurar búsqueda
    setupSearch();
});

// Configurar event listeners
function setupEventListeners() {
    // Botón para crear plantilla
    document.getElementById('createTemplateBtn').addEventListener('click', openAIModal);
    
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

// Configurar tema
function setupThemeToggle() {
    const themeToggle = document.getElementById('themeToggle');
    if (themeToggle) {
        themeToggle.addEventListener('click', () => {
            document.body.classList.toggle('light-theme');
            
            const icon = themeToggle.querySelector('i');
            if (document.body.classList.contains('light-theme')) {
                icon.classList.remove('fa-moon');
                icon.classList.add('fa-sun');
            } else {
                icon.classList.remove('fa-sun');
                icon.classList.add('fa-moon');
            }
        });
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
            throw new Error(`Error HTTP: ${response.status}`);
        }
        
        const data = await response.json();
        currentTemplates = Object.values(data.data || {});
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
    
    // Determinar estado
    const isDraft = template.status === 'draft' || !template.status;
    const isPublished = template.status === 'published';
    const isProcessing = template.status === 'processing';
    
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
                    <p class="text-secondary text-sm">${template.documentExt?.toUpperCase() || 'DOCX'} • ${sizeFormatted}</p>
                </div>
                <div>
                    ${isDraft ? '<span class="status-badge status-draft">Borrador</span>' : 
                      isPublished ? '<span class="status-badge status-published">Publicada</span>' : 
                      '<span class="status-badge status-processing">Procesando</span>'}
                </div>
            </div>
            
            <div class="mb-4">
                <p class="text-secondary text-sm line-clamp-3">
                    ${template.abstractDoc || 'Descripción no disponible'}
                </p>
            </div>
            
            <div class="flex items-center justify-between text-sm text-secondary mb-6">
                <span>Creada: ${createdDate}</span>
                ${template.aiGenerated ? '<span class="template-badge"><span class="material-symbols-outlined text-xs">auto_awesome</span>IA</span>' : ''}
            </div>
            
            <div class="flex gap-2">
                <button class="flex-1 px-3 py-2 bg-white/5 hover:bg-white/10 text-secondary rounded-lg transition-colors flex items-center justify-center gap-2 fill-template-btn" data-id="${template.id || template.documentHash}">
                    <span class="material-symbols-outlined text-sm">edit_document</span>
                    Llenar
                </button>
                <button class="flex-1 px-3 py-2 bg-blue-600/20 hover:bg-blue-600/30 text-blue-400 rounded-lg transition-colors flex items-center justify-center gap-2" onclick="downloadTemplate('${template.documentPath}', '${template.documentName}', '${template.id || template.documentHash}')">
                    <span class="material-symbols-outlined text-sm">download</span>
                    Descargar
                </button>
                <button class="px-3 py-2 bg-white/5 hover:bg-white/10 text-secondary rounded-lg transition-colors" onclick="viewTemplate('${template.documentPath}', '${template.id || template.documentHash}')">
                    <span class="material-symbols-outlined text-sm">visibility</span>
                </button>
            </div>
        </div>
    `;
    
    // Agregar event listener al botón de llenar
    const fillBtn = card.querySelector('.fill-template-btn');
    fillBtn.addEventListener('click', (e) => {
        e.stopPropagation();
        const templateId = e.currentTarget.dataset.id;
        openFillModal(templateId);
    });
    
    return card;
}

// Actualizar estadísticas
function updateStats() {
    // En una implementación real, esto vendría de la API
    // Por ahora, calculamos de los datos locales
    
    const total = totalTemplates;
    const published = currentTemplates.filter(t => t.status === 'published').length;
    const draft = currentTemplates.filter(t => !t.status || t.status === 'draft').length;
    const aiGenerated = currentTemplates.filter(t => t.aiGenerated).length;
    
    document.getElementById('totalTemplates').textContent = total;
    document.getElementById('publishedTemplates').textContent = published;
    document.getElementById('draftTemplates').textContent = draft;
    document.getElementById('aiGenerated').textContent = aiGenerated;
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
    const prompt = document.getElementById('aiPrompt').value.trim();
    const templateName = document.getElementById('templateName').value.trim();
    
    if (!prompt) {
        showMessage('Por favor, describe el documento que quieres crear.', 'error');
        return;
    }
    
    // Mostrar estado de carga
    document.getElementById('aiLoading').classList.remove('hidden');
    document.getElementById('generateTemplateBtn').disabled = true;
    
    try {
        // Llamar al endpoint de IA
        const response = await fetch('/interactDoc', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                type: 'agent',
                prompt: prompt,
                templateName: templateName || undefined,
                context: 'Estoy creando una plantilla de documento desde cero.'
            })
        });
        
        if (!response.ok) {
            throw new Error(`Error HTTP: ${response.status}`);
        }
        
        const data = await response.json();
        
        // Ocultar carga y mostrar respuesta
        document.getElementById('aiLoading').classList.add('hidden');
        document.getElementById('aiResponseContainer').classList.remove('hidden');
        
        // Mostrar respuesta
        document.getElementById('aiResponse').textContent = data.content || data.message || 'Plantilla generada exitosamente.';
        
        // Guardar el ID de la plantilla generada si existe
        if (data.templateId) {
            selectedTemplateId = data.templateId;
        }
        
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
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                type: 'agent',
                prompt: prompt,
                previousContent: document.getElementById('aiResponse').textContent
            })
        });
        
        const data = await response.json();
        document.getElementById('aiLoading').classList.add('hidden');
        document.getElementById('aiResponse').textContent = data.content || data.message;
        
    } catch (error) {
        console.error('Error al regenerar:', error);
        document.getElementById('aiLoading').classList.add('hidden');
        showMessage('Error al regenerar la plantilla.', 'error');
    }
}

// Crear PDF desde plantilla generada
async function createPdfFromTemplate() {
    const templateContent = document.getElementById('aiResponse').textContent;
    
    if (!templateContent || templateContent.trim() === '') {
        showMessage('No hay contenido para generar el PDF.', 'error');
        return;
    }
    
    try {
        // Mostrar carga
        document.getElementById('createPdfBtn').disabled = true;
        document.getElementById('createPdfBtn').innerHTML = '<div class="ai-loading"></div> Generando...';
        
        // Llamar al endpoint para crear PDF
        const response = await fetch('/interactDoc', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                type: 'agent',
                prompt: templateContent,
                templateId: selectedTemplateId || undefined
            })
        });
        
        if (!response.ok) {
            throw new Error(`Error HTTP: ${response.status}`);
        }
        
        const data = await response.json();
        
        // Subir el PDF generado
        if (data.pdfContent) {
            await uploadGeneratedPdf(data.pdfContent, data.filename || 'documento_generado.pdf');
        } else if (data.pdfUrl) {
            // Descargar el PDF
            const a = document.createElement('a');
            a.href = data.pdfUrl;
            a.download = data.filename || 'documento_generado.pdf';
            document.body.appendChild(a);
            a.click();
            document.body.removeChild(a);
        }
        
        showMessage('PDF generado exitosamente.', 'success');
        
        // Cerrar modal y recargar plantillas
        closeAIModal();
        loadTemplates();
        
    } catch (error) {
        console.error('Error al crear PDF:', error);
        showMessage('Error al generar el PDF.', 'error');
    } finally {
        document.getElementById('createPdfBtn').disabled = false;
        document.getElementById('createPdfBtn').innerHTML = '<span class="material-symbols-outlined text-sm mr-1">download</span> Crear PDF';
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
            throw new Error(`Error HTTP: ${response.status}`);
        }
        
        const result = await response.json();
        showMessage('Documento subido exitosamente.', 'success');
        
    } catch (error) {
        console.error('Error al subir PDF:', error);
        throw error;
    }
}

// Guardar plantilla
async function saveTemplate() {
    const templateName = document.getElementById('templateName').value.trim() || 
                         `Plantilla_${new Date().toISOString().slice(0, 10)}`;
    const templateContent = document.getElementById('aiResponse').textContent;
    
    if (!templateContent || templateContent.trim() === '') {
        showMessage('No hay contenido para guardar.', 'error');
        return;
    }
    
    try {
        document.getElementById('saveTemplateBtn').disabled = true;
        document.getElementById('saveTemplateBtn').innerHTML = '<div class="ai-loading"></div> Guardando...';
        
        // Guardar como documento de tipo template
        const response = await fetch('/saveTemplate', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                name: templateName,
                content: templateContent,
                type: 'template',
                aiGenerated: true
            })
        });
        
        if (!response.ok) {
            throw new Error(`Error HTTP: ${response.status}`);
        }
        
        const data = await response.json();
        showMessage('Plantilla guardada exitosamente.', 'success');
        
        // Cerrar modal y recargar
        closeAIModal();
        loadTemplates();
        
    } catch (error) {
        console.error('Error al guardar plantilla:', error);
        showMessage('Error al guardar la plantilla.', 'error');
    } finally {
        document.getElementById('saveTemplateBtn').disabled = false;
        document.getElementById('saveTemplateBtn').innerHTML = '<span class="material-symbols-outlined text-sm mr-1">save</span> Guardar Plantilla';
    }
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

// Generar documento desde plantilla
async function generateDocumentFromTemplate() {
    const prompt = document.getElementById('fillPrompt').value.trim();
    
    if (!prompt) {
        showMessage('Por favor, proporciona instrucciones para llenar la plantilla.', 'error');
        return;
    }
    
    document.getElementById('fillLoading').classList.remove('hidden');
    document.getElementById('generateFillBtn').disabled = true;
    
    try {
        // Obtener la plantilla seleccionada
        const template = currentTemplates.find(t => t.id === selectedTemplateId || t.documentHash === selectedTemplateId);
        
        if (!template) {
            throw new Error('Plantilla no encontrada');
        }
        
        // Llamar a la IA para llenar la plantilla
        const response = await fetch('/interactDoc', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                type: 'agent',
                templateId: selectedTemplateId,
                templateContent: template.content || template.abstractDoc,
                prompt: prompt,
                context: 'Estoy llenando una plantilla existente con información específica.'
            })
        });
        
        if (!response.ok) {
            throw new Error(`Error HTTP: ${response.status}`);
        }
        
        const data = await response.json();
        
        // Mostrar resultados
        document.getElementById('fillLoading').classList.add('hidden');
        document.getElementById('fillResponseContainer').classList.remove('hidden');
        document.getElementById('fillResponse').textContent = data.content || data.message;
        
        // Mostrar vista previa si está disponible
        if (data.previewHtml) {
            document.getElementById('documentPreview').innerHTML = data.previewHtml;
        }
        
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
        const template = currentTemplates.find(t => t.id === selectedTemplateId || t.documentHash === selectedTemplateId);
        
        const response = await fetch('/interactDoc', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                type: 'agent',
                templateId: selectedTemplateId,
                templateContent: template.content || template.abstractDoc,
                prompt: prompt,
                previousContent: document.getElementById('fillResponse').textContent
            })
        });
        
        const data = await response.json();
        document.getElementById('fillLoading').classList.add('hidden');
        document.getElementById('fillResponse').textContent = data.content || data.message;
        
    } catch (error) {
        console.error('Error al regenerar:', error);
        document.getElementById('fillLoading').classList.add('hidden');
        showMessage('Error al regenerar el documento.', 'error');
    }
}

// Descargar PDF generado
async function downloadGeneratedPdf() {
    try {
        document.getElementById('downloadPdfBtn').disabled = true;
        document.getElementById('downloadPdfBtn').innerHTML = '<div class="ai-loading"></div> Generando...';
        
        const content = document.getElementById('fillResponse').textContent;
        
        const response = await fetch('/interactDoc', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                type: 'agent',
                prompt: content,
                templateId: selectedTemplateId
            })
        });
        
        if (!response.ok) {
            throw new Error(`Error HTTP: ${response.status}`);
        }
        
        const data = await response.json();
        
        // Subir o descargar el PDF
        if (data.pdfUrl) {
            const a = document.createElement('a');
            a.href = data.pdfUrl;
            a.download = data.filename || `documento_${new Date().toISOString().slice(0, 10)}.pdf`;
            document.body.appendChild(a);
            a.click();
            document.body.removeChild(a);
        } else if (data.pdfContent) {
            await uploadGeneratedPdf(data.pdfContent, data.filename || `documento_${new Date().toISOString().slice(0, 10)}.pdf`);
        }
        
        showMessage('PDF descargado exitosamente.', 'success');
        
        // Cerrar modal
        closeFillModal();
        
    } catch (error) {
        console.error('Error al descargar PDF:', error);
        showMessage('Error al descargar el PDF.', 'error');
    } finally {
        document.getElementById('downloadPdfBtn').disabled = false;
        document.getElementById('downloadPdfBtn').innerHTML = '<span class="material-symbols-outlined text-sm mr-1">download</span> Descargar PDF';
    }
}

// Funciones auxiliares
function viewTemplate(path, id) {
    // Usar la misma función viewDocument de documents.js
    if (window.viewDocument) {
        window.viewDocument(path.split('/').pop() || 'plantilla', id);
    } else {
        // Implementación alternativa
        window.open(`/downloadDoc?id=${id}&type=template`, '_blank');
    }
}

async function downloadTemplate(path, name, id) {
    try {
        const response = await fetch('/downloadDoc', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ 
                id: id,
                type: "template",
                path: path,
            })
        });
        
        if (!response.ok) {
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

// Exportar funciones al scope global
window.openAIModal = openAIModal;
window.closeAIModal = closeAIModal;
window.openFillModal = openFillModal;
window.closeFillModal = closeFillModal;
window.viewTemplate = viewTemplate;
window.downloadTemplate = downloadTemplate;