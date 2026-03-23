// Al cargar la página
document.addEventListener('DOMContentLoaded', () => {
    // Configurar la zona de drag & drop
    configurarDragAndDrop();
    
    // Configurar el botón de envío
    document.getElementById('btnSubmit').addEventListener('click', subirArchivos);
});

//===================================== DRAG & DROP =========================================
function configurarDragAndDrop() {
    const dropZone = document.getElementById('dropZone');
    const fileInput = document.getElementById('fileInput');

    // 1. SOLUCIÓN AL DOBLE CLIC:
    // Usamos una bandera para evitar que el evento se propague infinitamente
    dropZone.addEventListener('click', (e) => {
        // Si el clic vino directamente del input, no hacemos nada más
        if (e.target === fileInput) return;
        
        // Evitamos que el clic del div haga cosas raras
        e.preventDefault();
        
        // Antes de abrir, verificamos el tipo de documento
        const documentType = document.getElementById('documentType').value;
        if (!documentType) {
            mostrarMensaje('Por favor, selecciona un tipo de documento primero', 'error');
            return;
        }

        fileInput.click();
    });

    // 2. SOLUCIÓN AL "NO CARGA LA SEGUNDA VEZ":
    fileInput.addEventListener('change', function(e) {
        if (this.files && this.files.length > 0) {
            handleFiles(this.files);
            
            // Limpiamos el valor del input esto permite que si el usuario elige el mismo archivo
            // el evento 'change' se vuelva a disparar siempre.
            this.value = ''; 
        }
    });

    // Prevenir comportamientos por defecto
    ['dragenter', 'dragover', 'dragleave', 'drop'].forEach(eventName => {
        dropZone.addEventListener(eventName, (e) => {
            e.preventDefault();
            e.stopPropagation();
        }, false);
    });

    // Efectos visuales
    ['dragenter', 'dragover'].forEach(eventName => {
        dropZone.addEventListener(eventName, () => dropZone.classList.add('dragover'), false);
    });

    ['dragleave', 'drop'].forEach(eventName => {
        dropZone.addEventListener(eventName, () => dropZone.classList.remove('dragover'), false);
    });

    dropZone.addEventListener('drop', (e) => {
        const dt = e.dataTransfer;
        handleFiles(dt.files);
    }, false);
}

// Archivos seleccionados para subir
let archivosParaSubir = [];

function handleFiles(files) {
    if (files.length === 0) return;
    
    // Verificar que se haya seleccionado un tipo de documento
    const documentType = document.getElementById('documentType').value;
    if (!documentType) {
        mostrarMensaje('Por favor, selecciona un tipo de documento primero', 'error');
        return;
    }
    
    for (let i = 0; i < files.length; i++) {
        const file = files[i];
        
        // Verificar si el archivo ya está en la lista
        if (archivosParaSubir.some(f => f.name === file.name && f.size === file.size)) {
            continue;
        }
        
        // Añadir a la lista
        archivosParaSubir.push(file);
        
        // Mostrar en la interfaz
        mostrarArchivoEnLista(file);
    }
    
    // Habilitar el botón de subir si hay archivos
    document.getElementById('btnSubmit').disabled = archivosParaSubir.length === 0;
}

function mostrarArchivoEnLista(file) {
    const filesList = document.getElementById('filesList');
    const fileId = 'file-' + Date.now() + '-' + Math.random().toString(36).substr(2, 9);
    
    const fileItem = document.createElement('div');
    fileItem.className = 'file-item';
    fileItem.id = fileId;
    
    // Formatear tamaño del archivo
    const fileSize = formatFileSize(file.size);
    
    fileItem.innerHTML = `
        <div class="file-info">
            <i class="fas fa-file file-icon"></i>
            <div class="file-details">
                <div class="file-name" title="${file.name}">${file.name}</div>
                <div class="file-size">${fileSize}</div>
                <div class="progress-bar">
                    <div class="progress" id="progress-${fileId}"></div>
                </div>
            </div>
        </div>
        <div class="file-actions">
            <button class="btn btn-remove" onclick="eliminarArchivo('${fileId}', '${file.name}', ${file.size})">
                <i class="fas fa-times"></i>
                Eliminar
            </button>
        </div>
    `;
    
    filesList.appendChild(fileItem);
}

function formatFileSize(bytes) {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const sizes = ['Bytes', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
}

function eliminarArchivo(fileId, fileName, fileSize) {
    // Eliminar de la lista visual
    const fileElement = document.getElementById(fileId);
    if (fileElement) {
        fileElement.remove();
    }
    
    // Eliminar del array de archivos
    archivosParaSubir = archivosParaSubir.filter(file => 
        !(file.name === fileName && file.size === fileSize)
    );
    
    // Deshabilitar el botón si no hay archivos
    document.getElementById('btnSubmit').disabled = archivosParaSubir.length === 0;
}

//===================================== UPLOAD FILES =========================================
async function subirArchivos() {
    const documentType = document.getElementById('documentType').value;
    if (!documentType) {
        mostrarMensaje('Por favor, selecciona un tipo de documento', 'error');
        return;
    }
    
    if (archivosParaSubir.length === 0) {
        mostrarMensaje('No hay archivos para subir', 'error');
        return;
    }
    
    try {
        mostrarMensaje('Subiendo archivos...', 'info');
        document.getElementById('btnSubmit').disabled = true;
        document.getElementById('btnSubmit').innerHTML = '<span class="loading"></span> Subiendo...';
        
        const formData = new FormData();
        
        // Añadir el tipo de documento
        formData.append('documentType', documentType);
        
        // Añadir todos los archivos (usando el mismo campo 'document' que espera tu backend)
        for (let i = 0; i < archivosParaSubir.length; i++) {
            formData.append('document', archivosParaSubir[i]);
        }
        
        // Realizar la petición
        const response = await fetch('/uploadDocs', {
            method: 'POST',
            credentials: 'include',
            body: formData
        });
        
        if (!response.ok) {
            if (response.status === 401) {  
                window.location.href = '/login';
            }
            const errorText = await response.text();
            throw new Error(`Error ${response.status}: ${errorText}`);
        }
        
        const result = await response.json();
        
        if (result.success) {
            mostrarMensaje(result.message, 'success');
            
            // Limpiar la lista de archivos
            archivosParaSubir = [];
            document.getElementById('filesList').innerHTML = '';
            document.getElementById('btnSubmit').disabled = true;
            document.getElementById('btnSubmit').innerHTML = '<i class="fas fa-upload"></i> Subir documentos';
            document.getElementById('documentType').value = '';
            
            // Redirigir a Mis Documentos después de 1.5s
            setTimeout(() => { window.location.href = '/mydocs'; }, 1500);
            
        } else {
            throw new Error(result.message || 'Error al subir archivos');
        }
    } catch (error) {
        console.error('Error al subir archivos:', error);
        mostrarMensaje('Error al subir archivos: ' + error.message, 'error');
        document.getElementById('btnSubmit').disabled = false;
        document.getElementById('btnSubmit').innerHTML = '<i class="fas fa-upload"></i> Subir documentos';
    }
}

//===================================== SHOW MESSAGES =========================================
function mostrarMensaje(mensaje, tipo) {
    const messageContainer = document.getElementById('messageContainer');
    const messageClass = tipo === 'error' ? 'error' : tipo === 'success' ? 'success' : 'info';
    messageContainer.innerHTML = `<div class="${messageClass} animate-fade-in">${mensaje}</div>`;
    
    // Auto-ocultar después de 5 segundos
    setTimeout(() => {
        messageContainer.innerHTML = '';
    }, 5000);
}