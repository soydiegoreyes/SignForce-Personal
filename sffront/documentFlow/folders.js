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
                'team': 'Folders del Equipo',
                'shared': 'Folders Compartidos conmigo'
            };
            
            document.getElementById('sectionTitle').textContent = titles[folderType] || 'Folders';
            
            // Cargar folders del tipo seleccionado
            loadFoldersData(folderType);
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
                case 'team':
                    requestBody.OnlyTeam = true;
                    break;
                case 'shared':
                    requestBody.OnlyShared = true;
                    break;
            }

            const response = await fetch('/getfolder', {
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
            document.getElementById('foldersTableBody').innerHTML = `
                <tr>
                    <td colspan="5" class="text-center py-8 text-red-400">
                        Error al cargar los folders. Intente nuevamente.
                    </td>
                </tr>
            `;
        }
    }

    // Poblar la tabla con los folders
    function populateFoldersTable(folders) {
        const tableBody = document.getElementById('foldersTableBody');
        tableBody.innerHTML = '';

        if (!folders || Object.keys(folders).length === 0) {
            tableBody.innerHTML = `
                <tr>
                    <td colspan="5" class="text-center py-8 text-secondary">
                        No hay folders en esta categoría.
                    </td>
                </tr>`;
            return;
        }

        Object.keys(folders).forEach(folderId => {
            const folder = folders[folderId];
            
            // Calcular progreso
            const numDocs = parseInt(folder.numDocs) || 0;
            const numDocsSign = parseInt(folder.numDocsSign) || 0;
            const progressPercent = numDocs > 0 ? Math.round((numDocsSign / numDocs) * 100) : 0;
            
            // Determinar estado
            let statusText = 'En Proceso';
            let statusClass = 'status-in-progress';
            
            if (folder.closedAt && folder.closedAt !== "") {
                statusText = 'Cerrado';
                statusClass = 'status-completed';
            } else if (folder.deletedAt && folder.deletedAt !== "") {
                statusText = 'Eliminado';
                statusClass = 'status-canceled';
            } else if (progressPercent === 100) {
                statusText = 'Completado';
                statusClass = 'status-completed';
            }
            
            // Formatear fecha de expiración
            const expirationDate = folder.expirationDate ? 
                new Date(folder.expirationDate).toLocaleDateString('es-MX', {
                    year: 'numeric',
                    month: 'short',
                    day: 'numeric'
                }) : 'Sin fecha límite';
            
            // Nombre del emisor
            const emisorName = folder.userEmisor?.nameUserEmisor || 
                              folder.userEmisor?.nameTeamEmisor || 
                              folder.userEmisor?.nameInstEmisor || 
                              'Usuario';

            const row = document.createElement('tr');
            row.addEventListener('click', () => {
                selectFolder(folderId, folder);
            });
            row.className = "glass-card fade-in cursor-pointer";
            row.innerHTML = `
                <td class="px-6 py-4 whitespace-nowrap">
                    <div class="flex items-center">
                        <div class="flex-shrink-0 h-10 w-10 bg-gray-600 rounded-full flex items-center justify-center">
                            <span class="material-symbols-outlined text-white">folder</span>
                        </div>
                        <div class="ml-4">
                            <div class="text-sm font-medium text-primary">Folder ${folderId}</div>
                            <div class="text-sm text-secondary">${folder.description || 'Sin descripción'}</div>
                        </div>
                    </div>
                </td>
                <td class="px-6 py-4 text-sm text-primary">${emisorName}</td>
                <td class="px-6 py-4 text-sm text-primary">${expirationDate}</td>
                <td class="px-6 py-4 text-sm text-secondary">${numDocsSign}/${numDocs} firmados</td>
                <td class="px-6 py-4">
                    <span class="document-status ${statusClass}">${statusText}</span>
                </td>
            `;

            tableBody.appendChild(row);
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
                IdDocs: documentIds,
                Type: "uploaded",
                Page: 1,
                PageSize: 50, // Número alto para obtener todos los documentos
                OrderBy: 'createdAtDoc',
                OrderDir: 'ASC'
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
                const doc = documents[docId];
                
                const docElement = document.createElement('div');
                docElement.className = 'glass-card rounded-lg p-4 cursor-pointer hover:bg-white/10';
                docElement.addEventListener('click', () => {
                    viewDocument(doc.documentName, docId, doc.documentPath);
                });
                
                const uploadDate = doc.createdAtDoc ? 
                    new Date(doc.createdAtDoc).toLocaleDateString('es-MX', {
                        year: 'numeric',
                        month: 'short',
                        day: 'numeric'
                    }) : 'Fecha no disponible';
                
                docElement.innerHTML = `
                    <div class="flex items-center justify-between">
                        <div class="flex items-center">
                            <span class="material-symbols-outlined text-primary mr-2">description</span>
                            <div>
                                <p class="text-sm font-medium text-primary">${doc.documentName}</p>
                                <p class="text-xs text-secondary">${uploadDate}</p>
                            </div>
                        </div>
                        <span class="material-symbols-outlined text-secondary">visibility</span>
                    </div>
                `;
                
                sidebarDocuments.appendChild(docElement);
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

    // Ver documento en modal
    async function viewDocument(name, docId, path) {
        const modal = document.getElementById('documentModal');
        const modalTitle = document.getElementById('modalTitle');
        const modalLoading = document.getElementById('modalLoading');
        const documentViewer = document.getElementById('documentViewer');
        
        modalTitle.textContent = `Visualizar: ${name}`;
        modalLoading.style.display = 'flex';
        documentViewer.style.display = 'none';
        modal.style.display = 'block';
        
        try {
            const response = await fetch('/downloadDoc', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ 
                    id: docId,
                    type: "uploaded",
                    path: path || "",
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