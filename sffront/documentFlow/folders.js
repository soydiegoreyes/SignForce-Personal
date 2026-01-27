// Variables globales
let currentFolders = {};
let currentFolderType = 'user';
let currentPage = 1;
let pageSize = 10;
let totalFolders = 0;
let selectedFolderId = null;

document.addEventListener('DOMContentLoaded', () => {
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
        
        moveBubble('folders');
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
    document.getElementById("logoutBtn").addEventListener("click", logout)
    

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
    
    // Cambio de tipo de folders
    document.querySelectorAll('[data-folder-type]').forEach(tab => {
        tab.addEventListener('click', (e) => {
            e.preventDefault();
            
            // Actualizar clases activas
            document.querySelectorAll('[data-folder-type]').forEach(t => {
                t.classList.remove('text-white', 'bg-gradient-to-r', 'from-blue-500', 'to-purple-600');
                t.classList.add('text-secondary', 'hover:bg-white/10');
            });
            
            tab.classList.remove('text-secondary', 'hover:bg-white/10');
            tab.classList.add('text-white', 'bg-gradient-to-r', 'from-blue-500', 'to-purple-600');
            
            // Actualizar título de sección
            const folderType = tab.getAttribute('data-folder-type');
            currentFolderType = folderType;
            
            const titles = {
                'user': 'Mis Folders',
                'shared': 'Folders Compartidos conmigo'
            };
            
            document.getElementById('sectionTitle').textContent = titles[folderType] || 'Folders';
            
            // Cargar folders del tipo seleccionado
            loadFoldersData(folderType);
        });
    });
    
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

    // Hacer las funciones globales
    window.viewDocument = viewDocument;

    // Cargar los folders al iniciar (mis folders por defecto)
    loadFoldersData('user');

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
// Cargar folders desde la API
async function loadFoldersData(folderType, page = 1) {
    try {
        // Construir el body según el tipo de folder
        const requestBody = {
            Page: page,
            PageSize: pageSize,
            OrderBy: 'lastModified',
            OrderDir: 'DESC'
        };
        
        // Configurar parámetros según el tipo de folder
        switch(folderType) {
            case 'user':
                requestBody.OnlyUser = true;
                break;
            case 'shared':
                requestBody.OnlyShared = true;
                break;
        }

        const response = await fetch('/getfolders', {
            method: 'POST',
            credentials: 'include',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(requestBody)
        });

        if (!response.ok) {
            throw new Error(`Error HTTP: ${response.status}`);
        }

        const responseData = await response.json();
        currentFolders = responseData.folders || {};
        currentPage = responseData.page || page;
        pageSize = responseData.page_size || 10;
        totalFolders = responseData.total || Object.keys(currentFolders).length;

        populateFoldersTable(currentFolders);
        renderPagination();
        
        // Limpiar sidebar cuando se cambia de tipo de folder
        clearSidebar();
    } catch (error) {
        console.error('Error al cargar los folders:', error);
        const gridContainer = document.getElementById('foldersGrid');
        if (gridContainer) {
            gridContainer.innerHTML = `
                <div class="text-center py-8 text-red-400 col-span-full">
                    Error al cargar los folders. Intente nuevamente.
                </div>`;
        }
    }
}

// Poblar la tabla con los folders
function populateFoldersTable(folders) {
    const gridContainer = document.getElementById('foldersGrid');
    if (!gridContainer) {
        console.warn('⚠️ No se encontró el contenedor #foldersGrid');
        return;
    }

    gridContainer.innerHTML = '';

    if (!folders || Object.keys(folders).length === 0) {
        gridContainer.innerHTML = `
            <div class="text-center py-8 text-secondary col-span-full">
                No hay folders en esta categoría.
            </div>`;
        return;
    }

    Object.keys(folders).forEach(folderId => {
        const folder = folders[folderId] || {};

        // 🧩 Valores por defecto seguros
        const numDocs = parseInt(folder.numDocs) || 0;
        const numDocsSign = parseInt(folder.numDocsSign) || 0;
        const progressPercent = numDocs > 0 ? Math.round((numDocsSign / numDocs) * 100) : 0;

        const statusText = folder.deletedAt
            ? 'Eliminado'
            : folder.closedAt
            ? 'Cerrado'
            : progressPercent === 100
            ? 'Completado'
            : 'En Proceso';

        const statusColor =
            statusText === 'Eliminado'
                ? 'text-red-400'
                : statusText === 'Cerrado'
                ? 'text-gray-400'
                : statusText === 'Completado'
                ? 'text-green-400'
                : 'text-yellow-400';

        const expirationDate = folder.expirationDate
            ? new Date(folder.expirationDate).toLocaleDateString('es-MX', {
                year: 'numeric',
                month: 'short',
                day: 'numeric',
            })
            : 'Sin fecha límite';

        const emisorName =
            folder.userEmisor?.nameUserEmisor ||
            folder.userEmisor?.nameInstEmisor ||
            'Usuario desconocido';

        const desc = folder.description?.trim() || 'Sin descripción';

        const card = document.createElement('div');
        card.className = 'feature-card glass-card p-6 rounded-xl cursor-pointer fade-in';
        card.addEventListener('click', () => selectFolder(folderId, folder));

        card.innerHTML = `
            <div class="flex items-center justify-between mb-3">
                <div class="flex items-center gap-2">
                    <span class="material-symbols-outlined text-white bg-blue-600 p-2 rounded-full">folder</span>
                    <h3 class="text-lg font-semibold text-primary">Folder ${folderId}</h3>
                </div>
                <span class="${statusColor} text-sm font-semibold">${statusText}</span>
            </div>

            <p class="text-sm text-secondary mb-2">${desc}</p>
            <p class="text-sm text-secondary mb-1"><strong>Emisor:</strong> ${emisorName}</p>
            <p class="text-sm text-secondary mb-1"><strong>Fecha límite:</strong> ${expirationDate}</p>
            <p class="text-sm text-secondary mb-3"><strong>Documentos:</strong> ${numDocsSign}/${numDocs} firmados</p>

            <div class="w-full bg-white/10 rounded-full h-2 mt-2">
                <div class="bg-blue-500 h-2 rounded-full transition-all duration-500" style="width:${progressPercent}%"></div>
            </div>
            <p class="text-xs text-secondary mt-1">${progressPercent}% completado</p>
        `;

        gridContainer.appendChild(card);
    });
}

// Seleccionar un folder y cargar sus documentos
async function selectFolder(folderId, folderData) {
    selectedFolderId = folderId;
    
    // Actualizar la información del folder en el sidebar
    document.getElementById('sidebarName').textContent = `Folder ${folderId}`;
    document.getElementById('sidebarOwner').textContent = `Creado por: ${folderData.userEmisor?.nameUserEmisor || 'Usuario'}`;
    
    const modifiedDate = folderData.lastModified ? 
        new Date(folderData.lastModified).toLocaleDateString('es-MX', {
            year: 'numeric',
            month: 'long',
            day: 'numeric',
            hour: '2-digit',
            minute: '2-digit'
        }) : 'Fecha no disponible';
        
    document.getElementById('sidebarDate').textContent = `Modificado: ${modifiedDate}`;
    
    // Actualizar barra de progreso
    const numDocs = parseInt(folderData.numDocs) || 0;
    const numDocsSign = parseInt(folderData.numDocsSign) || 0;
    const progressPercent = numDocs > 0 ? Math.round((numDocsSign / numDocs) * 100) : 0;
    
    document.getElementById('progressFill').style.width = `${progressPercent}%`;
    document.getElementById('progressText').textContent = `${progressPercent}% completado (${numDocsSign}/${numDocs} documentos)`;
    
    // Cargar documentos del folder
    await loadFolderDocuments(folderData.documents || []);
}

// Cargar documentos de un folder específico
async function loadFolderDocuments(documentIds) {
    const sidebarDocuments = document.getElementById('sidebarDocuments');
    
    if (documentIds.length === 0) {
        sidebarDocuments.innerHTML = '<p class="text-secondary text-center py-4">No hay documentos en este folder</p>';
        return;
    }

    try {
        const requestBody = {
            idDocs: documentIds.map(item => item.idDocument),
            page: 1,
            page_size: 50, // Número alto para obtener todos los documentos
            order_by: 'createdAtDoc',
            order_dir: 'ASC',
            type: "pdf"
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
        const documents = responseData.data || {};
        
        // Mostrar documentos en el sidebar
        sidebarDocuments.innerHTML = '';
        
        Object.keys(documents).forEach(docId => {
            const doc = documents[docId] || {};

            const uploadDate = doc.createdAtDoc
                ? new Date(doc.createdAtDoc).toLocaleDateString('es-MX', {
                    year: 'numeric',
                    month: 'short',
                    day: 'numeric',
                })
                : 'Fecha no disponible';

            const modifiedDate = doc.lastModifiedDoc
                ? new Date(doc.lastModifiedDoc).toLocaleDateString('es-MX', {
                    year: 'numeric',
                    month: 'short',
                    day: 'numeric',
                })
                : 'N/A';

            const ext = doc.documentExt?.toUpperCase() || 'N/A';
            const hashFull = doc.documentHash || '—';
            const isActive = doc.activeDoc === "1";

            const statusColor = isActive ? 'bg-green-500' : 'bg-yellow-500';
            const statusText = isActive ? 'Activo' : 'Inactivo';

            const abstract = (doc.abstractDoc && doc.abstractDoc.trim()) 
                ? doc.abstractDoc 
                : 'Sin descripción disponible.';

            const docCard = document.createElement('div');
            docCard.className =
                'feature-card glass-card p-4 rounded-xl cursor-pointer hover:bg-white/10 fade-in flex flex-col gap-2';
            docCard.addEventListener('click', () =>
                viewDocument(doc.documentName, docId, doc.documentPath)
            );

            docCard.innerHTML = `
                <div class="flex items-center justify-between mb-2">
                    <div class="flex items-center gap-2 w-full">
                        <span class="material-symbols-outlined text-blue-400 flex-shrink-0">description</span>
                        <div class="flex flex-col w-full">
                            <h4 class="text-sm font-semibold text-primary doc-name truncate-2">${doc.documentName || `Documento ${docId}`}</h4>
                            <p class="text-xs text-secondary">${uploadDate}</p>
                        </div>
                    </div>
                    
                </div>
                <div class="text-xs text-secondary space-y-1">
                    <div class="flex items-center gap-1">
                        <span class="w-3 h-3 rounded-full ${statusColor}"></span>
                        <span class="text-xs text-secondary">${statusText}</span>
                    </div>
                    <p><strong>Extensión:</strong> ${ext}</p>
                    <p><strong>Última mod.:</strong> ${modifiedDate}</p>
                    <p class="break-all"><strong>Hash:</strong> <span class="text-primary">${hashFull}</span></p>
                </div>

                <div class="mt-2 bg-white/5 p-2 rounded-lg">
                    <p class="text-xs text-secondary mb-1 font-semibold">Descripción:</p>
                    <p class="text-xs text-primary whitespace-pre-line">${abstract}</p>
                </div>
            `;

            sidebarDocuments.appendChild(docCard);
        });

    } catch (error) {
        console.error('Error al cargar documentos del folder:', error);
        sidebarDocuments.innerHTML = '<p class="text-red-400 text-center py-4">Error al cargar los documentos</p>';
    }
}

// Limpiar sidebar
function clearSidebar() {
    document.getElementById('sidebarName').textContent = 'Seleccione un folder';
    document.getElementById('sidebarOwner').textContent = '';
    document.getElementById('sidebarDate').textContent = '';
    document.getElementById('progressFill').style.width = '0%';
    document.getElementById('progressText').textContent = '0% completado';
    document.getElementById('sidebarDocuments').innerHTML = '<p class="text-secondary text-center py-4">Seleccione un folder para ver sus documentos</p>';
    document.getElementById('sidebarActions').style.display = 'none';
}

// Paginación
function renderPagination() {
    const container = document.getElementById('paginationContainer');
    container.innerHTML = '';

    if (totalFolders <= pageSize) return;

    const totalPages = Math.ceil(totalFolders / pageSize);

    const createButton = (text, page, disabled = false, active = false) => {
        const btn = document.createElement('button');
        btn.textContent = text;
        btn.className = `
            px-3 py-1 rounded-lg text-sm font-medium transition
            ${active ? 'bg-blue-600 text-white' : 'bg-white/10 text-secondary hover:bg-white/20'}
            ${disabled ? 'opacity-50 cursor-not-allowed' : ''}
        `;
        if (!disabled && !active) {
            btn.addEventListener('click', () => loadFoldersData(currentFolderType, page));
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

async function viewDocument(name, docId, path) {
    const isMobile = /Android|iPhone|iPad|iPod/i.test(navigator.userAgent);

    if (isMobile) {
        // En móvil NO usamos blob
        window.open(`/downloadDocStream?id=${docId}&type=uploaded`, '_blank');
        return;
    }

    // === Desktop (como ya lo tienes) ===
    const modal = document.getElementById('documentModal');
    const modalTitle = document.getElementById('modalTitle');
    const modalLoading = document.getElementById('modalLoading');
    const documentViewer = document.getElementById('documentViewer');

    modalTitle.textContent = `Visualizar: ${name}`;
    modalLoading.style.display = 'flex';
    documentViewer.style.display = 'none';
    modal.style.display = 'block';

    const response = await fetch('/downloadDoc', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ id: docId, type: "uploaded", path })
    });

    const blob = await response.blob();
    documentViewer.src = URL.createObjectURL(blob);

    modalLoading.style.display = 'none';
    documentViewer.style.display = 'block';
}

async function logout() {
    if (confirm('¿Cerrar sesión como administrador?')) {
        try {
            const response = await fetch('/logoutUser', {
                method: 'GET',
                headers: { 'Content-Type': 'application/json' },
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
