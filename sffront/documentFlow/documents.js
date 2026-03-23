window.signDocument = signDocument;
window.cancelProcess = cancelProcess;
window.acceptShared = acceptShared;
window.viewHistory = viewHistory;

// Chat globals
let chatHistory = [];
let chatDocId = null;

// Variables globales
let currentData = {};
let selectedDocuments = {};
/*==================================================== */
let currentTab = 'uploaded';
let currentSearchText = '';
let searchTimer = null;
const SEARCH_DELAY = 500; // ms de espera después de escribir
/*==================================================== */

// cargar documentos
let totalDocs = 0;
let currentPage = 1;
let pageSize = 10;
const searchParams = {
    page: 1,
    page_size: pageSize,
    order_by: 'createdAtDoc',
    order_dir: 'DESC',
    type: 'pdf'
};
/*====================================================*/




document.addEventListener('DOMContentLoaded', () => {
    // Variables globales locales al scope
    document.getElementById("logoutBtn").addEventListener("click", logout)
    
    // Configuración del buscador
    const searchInput = document.getElementById('searchInput');
    const clearSearchBtn = document.getElementById('clearSearch');
    
    if (searchInput) {
        // Buscar al escribir (con debounce)
        searchInput.addEventListener('input', (e) => {
            const text = e.target.value;
            
            // Mostrar/ocultar botón de limpiar
            if (clearSearchBtn) {
                if (text.trim()) {
                    clearSearchBtn.classList.remove('hidden');
                } else {
                    clearSearchBtn.classList.add('hidden');
                }
            }
            
            // Limpiar timer anterior
            if (searchTimer) {
                clearTimeout(searchTimer);
            }
            
            // Nuevo timer
            searchTimer = setTimeout(() => {
                performSearch(text);
            }, SEARCH_DELAY);
        });
        
        // Buscar al presionar Enter
        searchInput.addEventListener('keydown', (e) => {
            if (e.key === 'Enter') {
                if (searchTimer) {
                    clearTimeout(searchTimer);
                }
                performSearch(e.target.value);
            }
        });
    }
    
    if (clearSearchBtn) {
        clearSearchBtn.addEventListener('click', () => {
            searchInput.value = '';
            clearSearchBtn.classList.add('hidden');
            currentSearchText = '';
            delete searchParams["tags"];
            performSearch(''); // Recargar sin filtros
        });
    }

    // Configuración del input del chat para enviar con Enter
    const chatInput = document.getElementById('docChatInput');
    if (chatInput) {
        chatInput.addEventListener('keydown', (e) => {
            if (e.key === 'Enter') {
                e.preventDefault();
                sendDocChatMessage();
            }
        });
    }

    // Efecto liquid glass para la burbuja del menú
    const liquidBubble = document.getElementById('liquidBubble');
    const navItems = document.querySelectorAll('.nav-item');
    
    const itemPositions = {};
    
    function calculatePositions() {
        navItems.forEach(item => {
            const rect = item.getBoundingClientRect();
            const headerRect = document.getElementById('mainHeader').getBoundingClientRect();
            
            itemPositions[item.dataset.item] = {
                left: rect.left - headerRect.left,
                width: rect.width
            };
        });
        
        moveBubble('documents');
    }
    
    function moveBubble(itemName) {
        const item = itemPositions[itemName];
        if (item && liquidBubble) {
            liquidBubble.style.left = `${item.left - 10}px`;
            liquidBubble.style.width = `${item.width + 20}px`;
        }
    }
    
    navItems.forEach(item => {
        item.addEventListener('mouseenter', () => {
            moveBubble(item.dataset.item);
        });
    });

    window.addEventListener('load', calculatePositions);
    window.addEventListener('resize', calculatePositions);

    // Intersection Observer para animaciones
    const observerOptions = {
        threshold: 0.1
    };

    const observer = new IntersectionObserver((entries) => {
        entries.forEach(entry => {
            if (entry.isIntersecting) {
                entry.target.classList.add('fade-in');
                observer.unobserve(entry.target);
            }
        });
    }, observerOptions);

    document.querySelectorAll('.fade-in').forEach(element => {
        observer.observe(element);
    });
    
    // Cambio de pestañas
    document.querySelectorAll('[data-tab]').forEach(tab => {
        tab.addEventListener('click', (e) => {
            e.preventDefault();
            
            // Actualizar clases activas
            document.querySelectorAll('[data-tab]').forEach(t => {
                t.classList.remove('text-white', 'bg-gradient-to-r', 'from-blue-500', 'to-purple-600');
                t.classList.add('text-secondary', 'hover:bg-white/10');
            });
            
            tab.classList.remove('text-secondary', 'hover:bg-white/10');
            tab.classList.add('text-white', 'bg-gradient-to-r', 'from-blue-500', 'to-purple-600');
            
            // Actualizar título de sección
            const tabName = tab.getAttribute('data-tab');
            currentTab = tabName;
            
            const titles = {
                'uploaded': 'Documentos Subidos',
                'inprocess': 'Documentos en Proceso de Firma',
                'shared': 'Documentos Compartidos',
                'finished': 'Documentos Finalizados o Cancelados'
            };
            
            document.getElementById('sectionTitle').textContent = titles[tabName] || 'Documentos';
            
            // Limpiar búsqueda al cambiar de pestaña
            if (searchInput) {
                searchInput.value = '';
                currentSearchText = '';
                delete searchParams["tags"];
                if (clearSearchBtn) {
                    clearSearchBtn.classList.add('hidden');
                }
            }
            
            // Ocultar botón de consulta y chat al cambiar de pestaña
            document.getElementById('consultDocBtnContainer').classList.add('hidden');
            document.getElementById('leftSidebarChat').classList.add('hidden');

            // Cargar documentos de la pestaña seleccionada
            loadDocumentsData(tabName);
        });
    });

    // Función para realizar búsqueda
    async function performSearch(searchText = "", page = 1) {
        currentSearchText = (searchText || "").trim();

        // Mostrar indicador de búsqueda
        const tableBody = document.getElementById('documentsTableBody');
        if (tableBody) {
            tableBody.innerHTML = `
                <tr>
                    <td colspan="5" class="text-center py-8">
                        <div class="loading">
                            <div class="spinner"></div>
                            <div class="mt-2 text-secondary">Buscando...</div>
                        </div>
                    </td>
                </tr>
            `;
        }

        // Ajustar parámetros de búsqueda
        searchParams.page = page;
        // Mantener page_size definido en pageSize
        pageSize = searchParams.page_size || pageSize;

        // Si hay texto de búsqueda, enviarlo como tags
        if (currentSearchText) {
            const keywords = currentSearchText
                .split(/\s+/)
                .filter(word => word.length > 2)
                .map(word => word.toLowerCase());

            if (keywords.length > 0) {
                searchParams.tags = keywords;
            } else {
                delete searchParams.tags;
            }
        } else {
            delete searchParams.tags;
        }

        try {
            const response = await fetch('/statusDocs', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(searchParams)
            });

            if (!response.ok) {
                if (response.status === 401) {  
                    window.location.href = '/login';
                }
                throw new Error(`Error HTTP: ${response.status}`);
            }

            const responseData = await response.json();
            currentData = responseData.data || {};
            currentPage = responseData.page || page;
            pageSize = responseData.page_size || pageSize;
            totalDocs = responseData.total ?? Object.keys(currentData).length;

            populateTable(currentData, currentTab);
            renderPagination();

            // Mensaje de resultados (igual que antes)
            const messageContainer = document.getElementById('messageContainer');
            if (currentSearchText) {
                const count = Object.keys(currentData).length;
                messageContainer.innerHTML = `
                    <div class="bg-blue-500 text-white px-4 py-2 rounded-lg">
                        ${count} resultado${count !== 1 ? 's' : ''} para "${currentSearchText}"
                    </div>
                `;

                setTimeout(() => {
                    if (messageContainer.innerHTML.includes(currentSearchText)) {
                        messageContainer.innerHTML = '';
                    }
                }, 3000);
            } else {
                messageContainer.innerHTML = '';
            }

        } catch (error) {
            console.error('Error en la búsqueda:', error);
            const messageContainer = document.getElementById('messageContainer');
            if (messageContainer) {
                messageContainer.innerHTML = `
                    <div class="bg-red-500 text-white px-4 py-2 rounded-lg">
                        Error en la búsqueda. Intente nuevamente.
                    </div>
                `;
            }
        }
    }
    /*====================================================*/

    // Cargar documentos desde la API
    async function loadDocumentsData(tab, page = 1) {
        const tableBody = document.getElementById('documentsTableBody');
        if (tableBody) {
            tableBody.innerHTML = `
                <tr>
                    <td colspan="10" class="text-center py-6 text-secondary">
                        Cargando documentos...
                    </td>
                </tr>
            `;
        }

        if (window.loadingDocuments) return;
        window.loadingDocuments = true;

        // Actualizamos el parámetro de búsqueda (page) y llamamos al backend
        try {
            await performSearch("", page);
        } finally {
            window.loadingDocuments = false;
        }
    }

    function populateTable(data, tab) {
        const tableBody = document.getElementById('documentsTableBody');
        tableBody.innerHTML = '';

        if (!data || Object.keys(data).length === 0) {
            tableBody.innerHTML = `
                <tr>
                    <td colspan="6" class="text-center py-8 text-secondary">
                        <div class="btn-secondary">
                            <a href="/upload">Subir documentos</a>
                        </div>
                    </td>
                </tr>`;
            return;
        }

        const sortedData = Object.keys(data).sort((a, b) => {
            const docA = data[a], docB = data[b];
            return new Date(docB.createdAtDoc || 0) - new Date(docA.createdAtDoc || 0);
        });

        sortedData.forEach((docId) => {
            const doc = data[docId];
            const checked = docId in selectedDocuments? 'checked' : '';

            const uploadDate = doc.createdAtDoc
                ? new Date(doc.createdAtDoc).toLocaleDateString('es-MX', {
                    year: 'numeric',
                    month: 'short',
                    day: 'numeric',
                    hour: '2-digit',
                    minute: '2-digit',
                })
                : 'N/A';

            const statusText = doc.activeDoc === "1" ? "Activo" : "Inactivo";

            const row = document.createElement('tr');
            row.addEventListener('click', (e) => {
                // evitar conflicto con checkbox
                if (e.target.classList.contains('doc-checkbox')) return;

                const docData = data[docId];
                selectDocument(docId, docData);
            });
            row.className = "glass-card fade-in cursor-pointer";
            row.innerHTML = `
                <td class="px-4 text-center">
                    <input type="checkbox" class="doc-checkbox" data-id="${docId}" data-hash="${doc.documentHash || ''}" ${checked}>
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                    <div class="flex items-center">
                        <div class="flex-shrink-0 h-10 w-10 bg-gray-600 rounded-full flex items-center justify-center">
                            <span class="material-symbols-outlined text-white" onclick="viewDocument('${doc.documentName}', '${docId}')">description</span>
                        </div>
                        <div class="ml-4">
                            <div class="text-sm font-medium text-primary">${doc.documentName}</div>
                            <div class="text-sm text-secondary">${uploadDate}</div>
                        </div>
                    </div>
                </td>
                <td class="px-6 py-4 text-sm text-primary">${doc.documentExt || 'N/A'}</td>
                <td class="px-6 py-4 text-sm text-secondary">
                    ${formatBytes(doc.sizeB)}
                </td>
                <td class="px-6 py-4">${statusText}</td>
            `;

            tableBody.appendChild(row);
            if (selectedDocuments[docId]) {
                const checkbox = row.querySelector('.doc-checkbox');
                if (checkbox) checkbox.checked = true;
            }
        });

        document.querySelectorAll('.doc-checkbox').forEach(chk => {
            chk.addEventListener('change', (e) => {
                const id = e.target.dataset.id;
                const hash = e.target.dataset.hash || null;

                if (e.target.checked) {
                    selectedDocuments[id] = hash;
                } else {
                    delete selectedDocuments[id];
                }

                toggleGlobalSignButton();
            });
        });

        toggleGlobalSignButton();
    }

    function toggleGlobalSignButton() {
        
        let btnContainer = document.getElementById('globalSignBtnContainer');
        
        if (!btnContainer) {
            const section = document.getElementById('sectionHeader') || document.querySelector('h2');
            btnContainer = document.createElement('div');
            btnContainer.id = 'globalSignBtnContainer';
            btnContainer.className = 'mb-4 text-right';
        }

        const selectedCount = Object.keys(selectedDocuments).length;

        if (selectedCount > 0) {
            btnContainer.innerHTML = `
                <button id="createSignFolder"
                    class="py-2 px-4 rounded-lg bg-yellow-500 hover:bg-yellow-600 text-white font-semibold">
                    Iniciar proceso de firma (${selectedCount})
                </button>
            `;
            document.getElementById('createSignFolder').onclick = createSignFolder;
        } else {
            btnContainer.innerHTML = '';
        }
    }

    // --------------------
    // FUNCIÓN PRINCIPAL DE SELECCIÓN
    // --------------------
    function selectDocument(docId, documentData) {
        
        // 1. Verificar cambio de documento para reiniciar chat
        if (chatDocId !== docId) {
            chatDocId = docId;
            chatHistory = []; // Limpiar historial
            
            // Limpiar UI del chat
            const chatMessages = document.getElementById('docChatMessages');
            chatMessages.innerHTML = `
                <div class="text-center text-xs text-secondary mt-2">
                    Haz una pregunta sobre "${documentData.documentName || 'el documento'}".
                </div>
            `;
            
            // Cerrar el chat si estaba abierto con otro doc y mostrar solo el botón
            document.getElementById('leftSidebarChat').classList.add('hidden');
            document.getElementById('consultDocBtnContainer').classList.remove('hidden');
            
            // Actualizar título del chat
            const chatTitle = document.getElementById('chatDocTitle');
            if(chatTitle) chatTitle.textContent = documentData.documentName || 'Chat';
        }

        selectedDocumentId = docId;

        // Actualizar la información del documento en el sidebar DERECHO
        document.getElementById('sidebarName').textContent = documentData.documentName || `Documento ${docId}`;
        const sidebarState = document.getElementById('sidebarState');
        if (sidebarState) {
            const isActive = documentData.activeDoc === "1";
            const color = isActive ? 'bg-green-500' : 'bg-yellow-500';
            const text = isActive ? 'Activo' : 'Inactivo';

            sidebarState.innerHTML = `
                <div class="flex items-center gap-2">
                    <span class="w-3 h-3 rounded-full ${color}"></span>
                    <span class="text-primary text-sm">${text}</span>
                </div>
            `;
        }
        const uploadDate = documentData.createdAtDoc ?
            new Date(documentData.createdAtDoc).toLocaleDateString('es-MX', {
                year: 'numeric',
                month: 'long',
                day: 'numeric',
                hour: '2-digit',
                minute: '2-digit'
            }) : 'Fecha no disponible';

        document.getElementById('sidebarDate').textContent = `Subido: ${uploadDate}`;

        if (currentTab === 'inprocess') {
            const percent = documentData.progressPercent ? Math.min(100, Math.max(0, Number(documentData.progressPercent))) : 0;
            document.getElementById('progressFill').style.width = `${percent}%`;
            document.getElementById('progressText').textContent = `${percent}% completado`;
        } else {
            document.getElementById('progressFill').style.width = `100%`;
            document.getElementById('progressText').textContent = currentTab === 'uploaded' ? 'Documento subido' : (currentTab === 'finished' ? 'Proceso finalizado' : 'Disponible para firmar');
        }

        const detailsContainer = document.getElementById('sidebarDetailsContent');
        const modifiedDate = documentData.lastModifiedDoc ?
            new Date(documentData.lastModifiedDoc).toLocaleDateString('es-MX', {
                year: 'numeric',
                month: 'long',
                day: 'numeric',
                hour: '2-digit',
                minute: '2-digit'
            }) : 'N/A';

        const deletedInfo = documentData.deletedAtDoc && documentData.deletedAtDoc !== "" ? `
            <div class="col-span-2">
                <p class="text-sm text-secondary">Fecha de eliminación:</p>
                <p class="text-sm text-red-400">${new Date(documentData.deletedAtDoc).toLocaleDateString('es-MX')}</p>
            </div>
            ${documentData.deletedReasonDoc ? `<div class="col-span-2">
                <p class="text-sm text-secondary">Razón de eliminación:</p>
                <p class="text-sm text-red-400">${documentData.deletedReasonDoc}</p>
            </div>` : ''}
        ` : '';

        detailsContainer.innerHTML = `
            <div class="grid grid-cols-2 gap-4">
                <div>
                    <p class="text-sm text-secondary">Extensión:</p>
                    <p class="text-sm text-primary">${documentData.documentExt || 'N/A'}</p>
                </div>
                <div class="col-span-2">
                    <p class="text-sm text-secondary">Última modificación:</p>
                    <p class="text-sm text-primary">${modifiedDate}</p>
                </div>
                <div class="col-span-2">
                    <p class="text-sm text-secondary">Hash del documento:</p>
                    <p class="text-sm text-primary break-all">${documentData.documentHash || 'N/A'}</p>
                </div>
                ${documentData.abstractDoc ? `
                    <div class="col-span-2">
                        <p class="text-sm text-secondary mb-1">Descripción:</p>
                        <div class="text-sm text-primary bg-white/5 p-2 rounded-lg"
                            style="max-height: 150px; overflow-y: auto; scrollbar-width: thin;">
                            ${documentData.abstractDoc}
                        </div>
                    </div>` : ''}
                ${deletedInfo}
            </div>
        `;

        const signaturesContainer = document.getElementById('sidebarSignatures');
        if (currentTab === 'finished') {
            signaturesContainer.innerHTML = '<p class="text-secondary text-center py-4">Proceso finalizado. Ver historial para más detalles.</p>';
        } else if (currentTab === 'inprocess') {
            signaturesContainer.innerHTML = '<p class="text-secondary text-center py-4">Firmas en proceso. Use las acciones para firmar o cancelar.</p>';
        } else if (currentTab === 'shared') {
            signaturesContainer.innerHTML = '<p class="text-secondary text-center py-4">Documento compartido para su firma. Revise y firme si corresponde.</p>';
        } else {
            signaturesContainer.innerHTML = '<p class="text-secondary text-center py-4">Este es un documento subido sin proceso de firma</p>';
        }

        const sidebarActions = document.getElementById('sidebarActions');
        sidebarActions.innerHTML = ''; 

        // Botones según pestaña
        if (currentTab === 'uploaded') {
            sidebarActions.style.display = 'flex';
            const count = Object.keys(selectedDocuments).length;
            const disabled = count === 0 ? 'disabled' : '';

            sidebarActions.innerHTML = `
                <button id="createFolderSignSidebar"
                    class="w-full py-3 px-4 rounded-lg bg-yellow-500 hover:bg-yellow-600 text-white font-semibold ${disabled ? 'opacity-60 cursor-not-allowed' : ''}"
                    ${disabled}>
                    Iniciar Proceso de Firma ${count > 0 ? `(${count})` : ''}
                </button>
            `;
            const startBtn = document.getElementById('createFolderSignSidebar');
            startBtn.addEventListener('click', async () => {
                if (Object.keys(selectedDocuments).length === 0) {
                    sfAlert('Selecciona al menos un documento antes de iniciar el proceso.');
                    return;
                }
                createSignFolder();
            });
        } else if (currentTab === 'inprocess') {
            sidebarActions.style.display = 'flex';
            sidebarActions.innerHTML = `
                <button class="w-1/2 py-3 px-4 rounded-lg bg-blue-600 hover:bg-blue-700 text-white font-semibold" onclick="signDocument('${docId}')">Firmar / Aprobar</button>
                <button class="w-1/2 py-3 px-4 rounded-lg bg-red-600 hover:bg-red-700 text-white font-semibold" onclick="cancelProcess('${docId}')">Cancelar</button>
            `;
        } else if (currentTab === 'shared') {
            sidebarActions.style.display = 'flex';
            sidebarActions.innerHTML = `
                <button class="w-1/2 py-3 px-4 rounded-lg bg-blue-600 hover:bg-blue-700 text-white font-semibold" onclick="acceptShared('${docId}')">Firmar (Aceptar)</button>
                <button class="w-1/2 py-3 px-4 rounded-lg bg-gradient-to-r from-blue-500 to-purple-600 text-white font-semibold" onclick="viewDocument('${documentData.documentName}', '${docId}')">Ver</button>
            `;
        } else if (currentTab === 'finished') {
            sidebarActions.style.display = 'flex';
            sidebarActions.innerHTML = `
                <button class="w-1/2 py-3 px-4 rounded-lg bg-gray-700 hover:bg-gray-800 text-white font-semibold" onclick="viewHistory('${docId}')">Historial</button>
                <button class="w-1/2 py-3 px-4 rounded-lg bg-gradient-to-r from-blue-500 to-purple-600 text-white font-semibold" onclick="downloadDocument('${documentData.documentName}', '${docId}')">Descargar</button>
            `;
        } else {
            sidebarActions.style.display = 'none';
        }
        
        // NOTA: El botón de chat se maneja en el sidebar izquierdo ahora, 
        // ya no lo añadimos aquí a sidebarActions.
    }

    function renderPagination() {
        const container = document.getElementById('paginationContainer');
        if (!container) return;
        container.innerHTML = '';

        // Asegurarnos de tener pageSize
        pageSize = pageSize || searchParams.page_size || 10;

        if (totalDocs <= pageSize) return;

        const totalPages = Math.max(1, Math.ceil(totalDocs / pageSize));

        const createButton = (text, page, disabled = false, active = false) => {
            const btn = document.createElement('button');
            btn.textContent = text;
            btn.className = `
                px-3 py-1 rounded-lg text-sm font-medium transition
                ${active ? 'bg-blue-600 text-white' : 'bg-white/10 text-secondary hover:bg-white/20'}
                ${disabled ? 'opacity-50 cursor-not-allowed' : ''}
            `;
            if (!disabled && !active) {
                btn.addEventListener('click', () => loadDocumentsData(currentTab, page));
            }
            return btn;
        };

        container.appendChild(createButton('←', Math.max(1, currentPage - 1), currentPage === 1));

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

        container.appendChild(createButton('→', Math.min(totalPages, currentPage + 1), currentPage === totalPages));
    }

    // Ver documento en modal
    async function viewDocument(name, docId) {
        const modal = document.getElementById('documentModal');
        const modalTitle = document.getElementById('modalTitle');
        const modalLoading = document.getElementById('modalLoading');
        const documentViewer = document.getElementById('documentViewer');
        
        modalTitle.textContent = `Visualizar: ${name}`;
        modalLoading.style.display = 'flex';
        documentViewer.style.display = 'none';
        modal.style.display = 'block';
        
        try {
            const docs = {idDocs: [docId], type: "uploaded"}
            const response = await fetch('/downloadDoc', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(docs)
            });
            
            if (!response.ok) {
                if (response.status === 401) {  
                    window.location.href = '/login';
                }
                throw new Error(`Error HTTP: ${response.status}`);
            }
            
            const blob = await response.blob();
            const url = URL.createObjectURL(blob);
            
            documentViewer.src = url;
            modalLoading.style.display = 'none';
            documentViewer.style.display = 'block';
            
        } catch (error) {
            console.error('Error al cargar el documento:', error);
            modalLoading.innerHTML = `
                <div class="text-red-400 text-center">
                    <span class="material-symbols-outlined text-4xl mb-2">error</span>
                    <p>Error al cargar el documento. Intente nuevamente.</p>
                </div>
            `;
        }
    }

    async function downloadDocument(name, docId) {
        try {
            const docs = {idDocs: [docId], type: "uploaded", paths: []}
            const response = await fetch('/downloadDoc', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(docs)
            });
            
            if (!response.ok) {
                if (response.status === 401) {  
                    window.location.href = '/login';
                }
                throw new Error(`Error HTTP: ${response.status}`);
            }
            
            const blob = await response.blob();
            const url = URL.createObjectURL(blob);
            
            const a = document.createElement('a');
            a.href = url;
            a.download = name || 'documento';
            document.body.appendChild(a);
            a.click();
            document.body.removeChild(a);
            URL.revokeObjectURL(url);
            
        } catch (error) {
            console.error('Error al descargar el documento:', error);
            sfAlert('Error al descargar el documento. Intente nuevamente.');
        }
    }

    window.viewDocument = viewDocument;
    window.downloadDocument = downloadDocument;

    document.querySelector('.close-modal').addEventListener('click', () => {
        document.getElementById('documentModal').style.display = 'none';
    });

    window.addEventListener('click', (event) => {
        const modal = document.getElementById('documentModal');
        if (event.target === modal) {
            modal.style.display = 'none';
        }
    });
});

// Acciones: stubs / llamadas al backend
async function createSignFolder() {
    let docs = {idDocs: [], hashDocs: []}

    Object.entries(selectedDocuments).forEach(([id, hash]) => {
        docs.idDocs.push(id);
        docs.hashDocs.push(hash);
    });

    try {
        const resp = await fetch('/newSignFolder', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(docs)
        });

        if (!resp.ok) {
            if (resp.status === 401) {  
                window.location.href = '/login';
            }
            throw new Error(`HTTP ${resp.status}`);
        }
        const result = await resp.json();

        sessionStorage.setItem("folder", JSON.stringify(result));

        if (result.redirect_url) {
            window.location.href = result.redirect_url;
        } else {
            window.location.href = "/addSigners"; 
        }
    } catch (err) {
        console.error('Error creando el proceso:', err);
        sfAlert('No se pudo crear el proceso de firma.');
    }
}

async function signDocument(docId) {
    // mostrarMensaje es local en DOMContentLoaded, aquí usaremos alert o una global si es necesario
    // Para simplificar, asumiremos que existe o usamos alert
    sfAlert('Firmando documento...');
    try {
        const resp = await fetch('/signDoc', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ id: docId })
        });
        if (!resp.ok) {
            if (resp.status === 401) {  
                window.location.href = '/login';
            }
            throw new Error('Error al firmar');
        }
        sfAlert('Documento firmado correctamente');
        // Necesitaríamos recargar, pero loadDocumentsData está dentro del scope. 
        // Idealmente refactorizar para que loadDocumentsData sea global, o recargar página.
        window.location.reload(); 
    } catch (err) {
        console.error(err);
        sfAlert('Error al firmar el documento');
    }
}

async function cancelProcess(docId) {
    var _ok2 = await sfConfirm({title:'Cancelar firma',message:'¿Estás seguro que deseas cancelar el proceso de firma?',type:'danger',confirmText:'Cancelar proceso',confirmClass:'sf-modal-btn-danger'});if(!_ok2) return;
    try {
        const resp = await fetch('/cancelProcess', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ id: docId })
        });
        if (!resp.ok) {
            if (resp.status === 401) {  
                window.location.href = '/login';
            }
            throw new Error('Error al cancelar');
        }
        sfAlert('Proceso cancelado');
        window.location.reload();
    } catch (err) {
        console.error(err);
        sfAlert('No se pudo cancelar el proceso');
    }
}

async function acceptShared(docId) {
    try {
        const resp = await fetch('/acceptShared', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ id: docId })
        });
        if (!resp.ok) {
            if (resp.status === 401) {  
                window.location.href = '/login';
            }
            throw new Error('Error al aceptar');
        }
        sfAlert('Documento aceptado para firma');
        window.location.reload();
    } catch (err) {
        console.error(err);
        sfAlert('No se pudo aceptar el documento');
    }
}

function viewHistory(docId) {
    sfAlert('Abriendo historial...');
}

function formatBytes(bytes) {
    bytes = Number(bytes);
    if (isNaN(bytes)) return "N/A";

    const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
    const i = bytes === 0 ? 0 : Math.floor(Math.log(bytes) / Math.log(1024));
    const value = (bytes / Math.pow(1024, i)).toFixed(2);
    return `${value} ${sizes[i]}`;
}

// ---------------------------
// LÓGICA DEL CHAT ACTUALIZADA
// ---------------------------

// Abre el chat en el sidebar izquierdo y oculta el botón grande
window.openDocChat = function() {
    if (!chatDocId) return;
    document.getElementById('consultDocBtnContainer').classList.add('hidden');
    document.getElementById('leftSidebarChat').classList.remove('hidden');
    document.getElementById('docChatInput').focus();
};

// Cierra el chat y muestra de nuevo el botón
window.closeDocChat = function() {
    document.getElementById('leftSidebarChat').classList.add('hidden');
    document.getElementById('consultDocBtnContainer').classList.remove('hidden');
};

window.sendDocChatMessage = async function() {
    const input = document.getElementById('docChatInput');
    const question = input.value.trim();
    if (!question) return;

    input.value = "";

    // Agregar mensaje del usuario a la vista
    chatHistory.push({ role: "user", text: question });
    renderChatMessages();

    try {
        // Llamar backend
        const response = await fetch('/interactDoc', {
            method: 'POST',
            credentials: "include",
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                idDocument: chatDocId,
                prompt: question,
                type: "interact",
                history: chatHistory    // opcional si tu backend lo soporta
            })
        });
        if (!response.ok) {
            if (resp.status === 401) {  
                window.location.href = '/login';
            }
            throw new Error('Error al aceptar');
        }
        const data = await response.json();

        // Agregar respuesta del servidor
        chatHistory.push({
            role: "assistant",
            text: data.message || "Sin respuesta."
        });

        renderChatMessages();
    } catch (error) {
        console.error("Error chat:", error);
        chatHistory.push({
            role: "assistant",
            text: "Error de conexión con el asistente."
        });
        renderChatMessages();
    }
};

function renderChatMessages() {
    const container = document.getElementById('docChatMessages');
    container.innerHTML = "";

    chatHistory.forEach(msg => {
        const div = document.createElement('div');
        div.className = msg.role === "user"
            ? "text-right"
            : "text-left";

        div.innerHTML = `
            <div class="inline-block px-2 py-1.5 rounded-lg mb-1 max-w-[90%] break-words
                ${msg.role === "user"
                ? "bg-purple-600 text-white"
                : "bg-white/10 text-primary border border-gray-600"}">
                ${msg.text}
            </div>
        `;

        container.appendChild(div);
    });

    // Scroll al final
    container.scrollTop = container.scrollHeight;
}
async function logout() {
    sfConfirm({title:'Cerrar sesión',message:'¿Seguro que deseas cerrar tu sesión actual?',type:'warn',confirmText:'Cerrar sesión',confirmClass:'sf-modal-btn-danger'}).then(function(ok){if(ok){
        try {
            const response = await fetch('/logoutUser', {
                method: 'GET',
                headers: { 'Content-Type': 'application/json' }}),
                credentials: 'include'
            });
            if (!response.ok) {
                if (resp.status === 401) {  
                    window.location.href = '/login';
                }
                throw new Error('Error al aceptar');
            }
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