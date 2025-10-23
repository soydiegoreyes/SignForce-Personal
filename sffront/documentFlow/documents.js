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

    // Cargar documentos desde la API
    async function loadDocumentsData(tab) {
        try {
            // Para documentos subidos, el endpoint espera un body que puede ir vacío
            // o con parámetros específicos de búsqueda
            const requestBody = {
                page: 1,
                page_size: 100,
                order_by: 'createdAtDoc',
                order_dir: 'DESC'
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
            
            // La respuesta tiene la estructura: { data: {...}, page: 1, page_size: 10, total: 1 }
            currentData = responseData.data || {};
            populateTable(currentData, tab);
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
    }

    // Poblar la tabla con los documentos
    function populateTable(data, tab) {
        const tableBody = document.getElementById('documentsTableBody');
        tableBody.innerHTML = '';

        if (!data || Object.keys(data).length === 0) {
            tableBody.innerHTML = `
                <tr>
                    <td colspan="5" class="text-center py-8 text-secondary">
                        No hay documentos en esta categoría.
                    </td>
                </tr>
            `;
            return;
        }

        // Ordenar documentos por fecha de creación (más recientes primero)
        const sortedData = Object.keys(data).sort((a, b) => {
            const docA = data[a];
            const docB = data[b];
            
            const dateA = new Date(docA.createdAtDoc || 0);
            const dateB = new Date(docB.createdAtDoc || 0);
            
            return dateB - dateA;
        });

        sortedData.forEach((docId, index) => {
            const doc = data[docId];
            const delayClass = `delay-${(index % 3) + 1}`;
            
            const row = document.createElement('tr');
            row.className = `cursor-pointer glass-card fade-in ${delayClass} document-card`;
            row.dataset.documentId = docId;
            
            // Formatear fecha
            const uploadDate = doc.createdAtDoc ? 
                new Date(doc.createdAtDoc).toLocaleDateString('es-MX', {
                    year: 'numeric',
                    month: 'short',
                    day: 'numeric',
                    hour: '2-digit',
                    minute: '2-digit'
                }) : 'Fecha no disponible';
            
            // Determinar el estado del documento
            let statusClass, statusText;
            if (doc.deletedAtDoc && doc.deletedAtDoc !== "") {
                statusClass = 'status-canceled';
                statusText = 'Eliminado';
            } else if (doc.activeDoc === "1") {
                statusClass = 'status-completed';
                statusText = 'Activo';
            } else {
                statusClass = 'status-in-progress';
                statusText = 'Inactivo';
            }
            
            // Obtener extensión del archivo
            const fileExt = doc.documentExt || 'file';
            
            row.innerHTML = `
                <td class="px-6 py-4 whitespace-nowrap">
                    <div class="flex items-center">
                        <div class="flex-shrink-0 h-10 w-10 bg-gray-600 rounded-full flex items-center justify-center">
                            <span class="material-symbols-outlined text-white">description</span>
                        </div>
                        <div class="ml-4">
                            <div class="text-sm font-medium text-primary">${doc.documentName || `Documento ${docId}`}</div>
                            <div class="text-sm text-secondary">${uploadDate}</div>
                        </div>
                    </div>
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                    <div class="text-sm text-primary">Tipo: ${fileExt.toUpperCase()}</div>
                    <div class="text-sm text-secondary">${doc.abstractDoc || 'Sin descripción'}</div>
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                    <div class="text-sm text-primary">${doc.lastModifiedDoc ? 
                        new Date(doc.lastModifiedDoc).toLocaleDateString('es-MX', {
                            year: 'numeric',
                            month: 'short',
                            day: 'numeric'
                        }) : 'N/A'}</div>
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                    <div class="flex items-center gap-2">
                        <button class="text-blue-500 hover:text-blue-400" onclick="viewDocument('${doc.documentName}', '${docId}')">
                            <span class="material-symbols-outlined">visibility</span>
                        </button>
                        <button class="text-green-500 hover:text-green-400" onclick="downloadDocument('${doc.documentName}', '${docId}')">
                            <span class="material-symbols-outlined">download</span>
                        </button>
                    </div>
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                    <span class="document-status ${statusClass}">${statusText}</span>
                </td>
            `;
            
            row.addEventListener('click', (e) => {
                // No seleccionar si se hizo clic en un botón
                if (!e.target.closest('button')) {
                    selectDocument(docId, doc);
                }
            });
            
            tableBody.appendChild(row);
        });
    }

    // Seleccionar un documento y mostrar sus detalles
    function selectDocument(docId, documentData) {
        selectedDocumentId = docId;
        
        // Actualizar la información del documento en el sidebar
        document.getElementById('sidebarName').textContent = documentData.documentName || `Documento ${docId}`;
        document.getElementById('sidebarOwner').textContent = `Hash: ${documentData.documentHash.substring(0, 20)}...`;
        
        const uploadDate = documentData.createdAtDoc ? 
            new Date(documentData.createdAtDoc).toLocaleDateString('es-MX', {
                year: 'numeric',
                month: 'long',
                day: 'numeric',
                hour: '2-digit',
                minute: '2-digit'
            }) : 'Fecha no disponible';
        
        document.getElementById('sidebarDate').textContent = `Subido: ${uploadDate}`;
        
        // Ocultar barra de progreso (no aplica para documentos subidos)
        document.getElementById('progressFill').style.width = `100%`;
        document.getElementById('progressText').textContent = `Documento subido`;
        
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
        
        const deletedInfo = documentData.deletedAtDoc && documentData.deletedAtDoc !== "" ?
            `<div class="col-span-2">
                <p class="text-sm text-secondary">Fecha de eliminación:</p>
                <p class="text-sm text-red-400">${new Date(documentData.deletedAtDoc).toLocaleDateString('es-MX')}</p>
            </div>
            ${documentData.deletedReasonDoc ? `<div class="col-span-2">
                <p class="text-sm text-secondary">Razón de eliminación:</p>
                <p class="text-sm text-red-400">${documentData.deletedReasonDoc}</p>
            </div>` : ''}` : '';
        
        detailsContainer.innerHTML = `
            <div class="grid grid-cols-2 gap-4">
                <div>
                    <p class="text-sm text-secondary">Extensión:</p>
                    <p class="text-sm text-primary">${documentData.documentExt || 'N/A'}</p>
                </div>
                <div>
                    <p class="text-sm text-secondary">Estado:</p>
                    <p class="text-sm text-primary">${documentData.activeDoc === "1" ? 'Activo' : 'Inactivo'}</p>
                </div>
                <div class="col-span-2">
                    <p class="text-sm text-secondary">Ruta del archivo:</p>
                    <p class="text-sm text-primary break-all">${documentData.documentPath || 'N/A'}</p>
                </div>
                <div class="col-span-2">
                    <p class="text-sm text-secondary">Última modificación:</p>
                    <p class="text-sm text-primary">${modifiedDate}</p>
                </div>
                <div class="col-span-2">
                    <p class="text-sm text-secondary">Hash del documento:</p>
                    <p class="text-sm text-primary break-all">${documentData.documentHash || 'N/A'}</p>
                </div>
                ${documentData.abstractDoc ? `<div class="col-span-2">
                    <p class="text-sm text-secondary">Descripción:</p>
                    <p class="text-sm text-primary">${documentData.abstractDoc}</p>
                </div>` : ''}
                ${deletedInfo}
            </div>
            <div class="mt-6 flex gap-2">
                <button class="flex-1 py-2 px-4 rounded-lg bg-blue-600 hover:bg-blue-700 text-white font-semibold flex items-center justify-center gap-2" onclick="viewDocument('${documentData.documentPath}', '${documentData.documentName}', '${docId}')">
                    <span class="material-symbols-outlined">visibility</span>
                    Ver Documento
                </button>
                <button class="flex-1 py-2 px-4 rounded-lg bg-green-600 hover:bg-green-700 text-white font-semibold flex items-center justify-center gap-2" onclick="downloadDocument('${documentData.documentPath}', '${documentData.documentName}', '${docId}')">
                    <span class="material-symbols-outlined">download</span>
                    Descargar
                </button>
            </div>
        `;
        
        // Limpiar el historial de firmas (no aplica para documentos subidos)
        const signaturesContainer = document.getElementById('sidebarSignatures');
        signaturesContainer.innerHTML = '<p class="text-secondary text-center py-4">Este es un documento subido sin proceso de firma</p>';
        
        // Ocultar botones de acción (no aplica para documentos subidos)
        document.getElementById('sidebarActions').style.display = 'none';
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
    loadDocumentsData('uploaded');
});