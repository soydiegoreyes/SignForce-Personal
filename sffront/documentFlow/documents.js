window.signDocument = signDocument;
window.cancelProcess = cancelProcess;
window.acceptShared = acceptShared;
window.viewHistory = viewHistory;

// Variables globales
let currentData = {};
let selectedDocuments = {};
let currentTab = 'uploaded';

document.addEventListener('DOMContentLoaded', () => {
    // Variables globales
    let currentData = {};
    let selectedDocumentId = null;
    let currentTab = 'uploaded';

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

    // Toggle de tema claro/oscuro
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
            
            // Cargar documentos de la pestaña seleccionada
            loadDocumentsData(tabName);
        });
    });
    
    function mostrarMensaje(mensaje, tipo) {
        const messageContainer = document.getElementById('messageContainer');
        const alertClass = tipo === 'error' ? 'bg-red-500' : 'bg-green-500';
        messageContainer.innerHTML = `<div class="${alertClass} text-white px-4 py-2 rounded-lg">${mensaje}</div>`;
        
        // Auto-ocultar después de 5 segundos
        setTimeout(() => {
            messageContainer.innerHTML = '';
        }, 5000);
    }
    // cargar documentos
    let currentPage = 1;
    let pageSize = 10;
    let totalDocs = 0;
    // Cargar documentos desde la API
    async function loadDocumentsData(tab, page = 1) {
        // --- FIX: limpiar antes de cargar (previene páginas congeladas) ---
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

        // --- FIX: bloquear peticiones mientras una ya está en curso ---
        if (window.loadingDocuments) return;
        window.loadingDocuments = true;
        try {
            const requestBody = {
                page: page,
                page_size: pageSize,
                order_by: 'createdAtDoc',
                order_dir: 'ASC'
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

            const responseData = await response.json();
            currentData = responseData.data || {};
            currentPage = responseData.page || page;
            pageSize = responseData.page_size || 10;
            totalDocs = responseData.total || Object.keys(currentData).length;

            populateTable(currentData, tab);
            renderPagination();
        } catch (error) {
            console.error('Error al cargar los documentos:', error);
            document.getElementById('documentsTableBody').innerHTML = `
                <tr>
                    <td colspan="5" class="text-center py-8 text-red-400">
                        Error al cargar los documentos. Intente nuevamente.
                    </td>
                </tr>
            `;
        }
        window.loadingDocuments = false;
    }

    // Poblar la tabla con los documentos
    // --------------------
    // --- dentro de populateTable ---
    function populateTable(data, tab) {
        const tableBody = document.getElementById('documentsTableBody');
        tableBody.innerHTML = '';

        if (!data || Object.keys(data).length === 0) {
            tableBody.innerHTML = `
                <tr>
                    <td colspan="6" class="text-center py-8 text-secondary">
                        No hay documentos en esta categoría.
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
                            <span class="material-symbols-outlined text-white" onclick="viewDocument('${doc.documentPath || doc.documentName}', '${docId}')">description</span>
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

        // Escuchar cambios de checkboxes
        document.querySelectorAll('.doc-checkbox').forEach(chk => {
            chk.addEventListener('change', (e) => {
                const id = e.target.dataset.id;
                const hash = e.target.dataset.hash || null;

                if (e.target.checked) {
                    selectedDocuments[id] = hash;  // añade o actualiza
                } else {
                    delete selectedDocuments[id];  // elimina si se desmarca
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
            section.parentNode.insertBefore(btnContainer, section.nextSibling);
        }

        if (selectedDocuments.keys > 0) {
            btnContainer.innerHTML = `
                <button id="createSignFolder"
                    class="py-2 px-4 rounded-lg bg-yellow-500 hover:bg-yellow-600 text-white font-semibold">
                    Iniciar proceso de firma (${selectedDocuments.size})
                </button>
            `;
            document.getElementById('createSignFolder').onclick = createSignFolder;
        } else {
            btnContainer.innerHTML = '';
        }
        if (document.getElementById('createFolderSignSidebar')) selectDocument(selectedDocumentId, currentData[selectedDocumentId]);
    }
    // --------------------
    // Reemplazo: selectDocument
    // --------------------
    function selectDocument(docId, documentData) {
        selectedDocumentId = docId;

        // Actualizar la información del documento en el sidebar
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

        // Barra de progreso para pestañas que la requieran
        if (currentTab === 'inprocess') {
            const percent = documentData.progressPercent ? Math.min(100, Math.max(0, Number(documentData.progressPercent))) : 0;
            document.getElementById('progressFill').style.width = `${percent}%`;
            document.getElementById('progressText').textContent = `${percent}% completado`;
        } else {
            document.getElementById('progressFill').style.width = `100%`;
            document.getElementById('progressText').textContent = currentTab === 'uploaded' ? 'Documento subido' : (currentTab === 'finished' ? 'Proceso finalizado' : 'Disponible para firmar');
        }

        // Actualizar detalles del documento
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

        // Historial/simple mensaje según pestaña
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

        // Sidebar actions dinámicas (ahora SÓLO en el sidebar)
        const sidebarActions = document.getElementById('sidebarActions');
        sidebarActions.innerHTML = ''; // limpiar

        // Determinar si el usuario es propietario (ajusta según tu API; si el campo no existe asumimos false)
        const isOwner = !!documentData.isOwner || !!documentData.ownerIsMe;

        // Construir botones según la pestaña
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

            // Escuchar clic en el botón (crea el folder con todos los documentos seleccionados)
            const startBtn = document.getElementById('createFolderSignSidebar');
            startBtn.addEventListener('click', async () => {
                if (Object.keys(selectedDocuments).length === 0) {
                    alert('Selecciona al menos un documento antes de iniciar el proceso.');
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
                <button class="w-1/2 py-3 px-4 rounded-lg bg-gradient-to-r from-blue-500 to-purple-600 text-white font-semibold" onclick="viewDocument('${documentData.documentPath || documentData.documentName}', '${docId}')">Ver</button>
            `;
        } else if (currentTab === 'finished') {
            sidebarActions.style.display = 'flex';
            sidebarActions.innerHTML = `
                <button class="w-1/2 py-3 px-4 rounded-lg bg-gray-700 hover:bg-gray-800 text-white font-semibold" onclick="viewHistory('${docId}')">Historial</button>
                <button class="w-1/2 py-3 px-4 rounded-lg bg-gradient-to-r from-blue-500 to-purple-600 text-white font-semibold" onclick="downloadDocument('${documentData.documentPath}', '${documentData.documentName}', '${docId}')">Descargar</button>
            `;
        } else {
            sidebarActions.style.display = 'none';
        }
    }

    function renderPagination() {
        const container = document.getElementById('paginationContainer');
        container.innerHTML = '';

        if (totalDocs <= pageSize) return;

        const totalPages = Math.ceil(totalDocs / pageSize);

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

        // Botón anterior
        container.appendChild(createButton('←', currentPage - 1, currentPage === 1));

        // Botones numéricos
        for (let i = 1; i <= totalPages; i++) {
            // Solo mostrar los primeros, últimos y cercanos a la página actual
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
            // Descargar el documento usando el endpoint /downloadDoc
            // Descargar el documento
            const response = await fetch('/downloadDoc', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ 
                    id: docId,
                    type: "uploaded",
                    path: "",

                })
            });
            
            if (!response.ok) {
                throw new Error(`Error HTTP: ${response.status}`);
            }
            
            // Crear un blob a partir de la respuesta
            const blob = await response.blob();
            const url = URL.createObjectURL(blob);
            
            // Mostrar el documento en el iframe
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

    // Descargar documento
    async function downloadDocument(path, name, docId) {
        try {
            const response = await fetch('/downloadDoc', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ 
                    id: docId,
                    type: "uploaded",
                    path: path,
                })
            });
            
            if (!response.ok) {
                throw new Error(`Error HTTP: ${response.status}`);
            }
            
            const blob = await response.blob();
            const url = URL.createObjectURL(blob);
            
            // Crear un enlace temporal para descargar el archivo
            const a = document.createElement('a');
            a.href = url;
            a.download = name || 'documento';
            document.body.appendChild(a);
            a.click();
            document.body.removeChild(a);
            URL.revokeObjectURL(url);
            
        } catch (error) {
            console.error('Error al descargar el documento:', error);
            alert('Error al descargar el documento. Intente nuevamente.');
        }
    }

    // Hacer las funciones globales para que puedan ser llamadas desde los botones
    window.viewDocument = viewDocument;
    window.downloadDocument = downloadDocument;

    // Cerrar modal
    document.querySelector('.close-modal').addEventListener('click', () => {
        document.getElementById('documentModal').style.display = 'none';
    });

    // Cerrar modal al hacer clic fuera del contenido
    window.addEventListener('click', (event) => {
        const modal = document.getElementById('documentModal');
        if (event.target === modal) {
            modal.style.display = 'none';
        }
    });

    // Botones de acción - Eliminados ya que no aplican para documentos subidos
    // Si en el futuro necesitas estos botones para otras pestañas, puedes restaurarlos

    // Cargar los datos al iniciar (documentos subidos por defecto)
    //loadDocumentsData('uploaded');
});

// Acciones: stubs / llamadas al backend
// --------------------
async function createSignFolder() {
    const docs = Object.entries(selectedDocuments).map(([id, hash]) => ({
        idDoc: id,
        hashDoc: hash
    }));

    try {
        const resp = await fetch('/newSignFolder', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(docs)
        });

        if (!resp.ok) throw new Error(`HTTP ${resp.status}`);
        const result = await resp.json();

        // Guardar en sessionStorage para usar en add_signers.js
        sessionStorage.setItem("folder", JSON.stringify(result));

        if (result.redirect_url) {
            window.location.href = result.redirect_url;
        } else {
            window.location.href = "/addSigners"; // ruta por defecto si no hay redirect_url
        }
    } catch (err) {
        console.error('Error creando el proceso:', err);
        alert('No se pudo crear el proceso de firma.');
    }
}

async function signDocument(docId) {
    mostrarMensaje('Firmando documento...', 'info');
    try {
        const resp = await fetch('/signDoc', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ id: docId })
        });
        if (!resp.ok) throw new Error('Error al firmar');
        mostrarMensaje('Documento firmado correctamente', 'success');
        loadDocumentsData(currentTab, currentPage || 1);
    } catch (err) {
        console.error(err);
        mostrarMensaje('Error al firmar el documento', 'error');
    }
}

async function cancelProcess(docId) {
    if (!confirm('¿Estás seguro que deseas cancelar el proceso de firma?')) return;
    mostrarMensaje('Cancelando proceso...', 'info');
    try {
        const resp = await fetch('/cancelProcess', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ id: docId })
        });
        if (!resp.ok) throw new Error('Error al cancelar');
        mostrarMensaje('Proceso cancelado', 'success');
        loadDocumentsData(currentTab, currentPage || 1);
    } catch (err) {
        console.error(err);
        mostrarMensaje('No se pudo cancelar el proceso', 'error');
    }
}

async function acceptShared(docId) {
    mostrarMensaje('Aceptando documento compartido...', 'info');
    try {
        const resp = await fetch('/acceptShared', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ id: docId })
        });
        if (!resp.ok) throw new Error('Error al aceptar');
        mostrarMensaje('Documento aceptado para firma', 'success');
        loadDocumentsData(currentTab, currentPage || 1);
    } catch (err) {
        console.error(err);
        mostrarMensaje('No se pudo aceptar el documento', 'error');
    }
}

function viewHistory(docId) {
    // Mostrar modal o ir a ruta de historial
    mostrarMensaje('Abriendo historial...', 'info');
    // Aquí podrías abrir un modal o navegar a /document/:id/history
    // window.location.href = `/document/${docId}/history`;
}

function formatBytes(bytes) {
    bytes = Number(bytes);
    if (isNaN(bytes)) return "N/A";

    const sizes = ['B', 'KB', 'MB', 'GB', 'TB'];
    const i = bytes === 0 ? 0 : Math.floor(Math.log(bytes) / Math.log(1024));
    const value = (bytes / Math.pow(1024, i)).toFixed(2);
    return `${value} ${sizes[i]}`;
}