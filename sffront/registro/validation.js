// Mapeo de tipos de documentos a las propiedades de la respuesta
const documentTypes = {
    'acta': 'docActa',
    'poder': 'docPoder',
    'identificacion': 'docIdentidad',
    'domicilio': 'docResidencia'
};

// Estado de documentos
let documentos = {
    acta: { subido: false, archivo: null, nombre: '' },
    poder: { subido: false, archivo: null, nombre: '' },
    identificacion: { subido: false, archivo: null, nombre: '' },
    domicilio: { subido: false, archivo: null, nombre: '' }
};

// Información de la empresa
let empresa = null;

// Al cargar la página
document.addEventListener('DOMContentLoaded', () => {
    // Obtener datos de validación
    obtenerDatosValidacion();
    
    // Configurar event listeners para los inputs de archivo
    configurarEventListeners();
    
    // Configurar el botón de completar validación
    document.getElementById('btnCompletar').addEventListener('click', completarValidacion);
    
});

//===================================== GET DATA VALIDATION =========================================
async function obtenerDatosValidacion() {
    try {
        mostrarMensaje('Cargando datos de validación...', 'info');
        
        const response = await fetch('/getvaldata', {
            method: 'GET',
            credentials: 'include' // Incluir cookies para el JWT
        });
        
        if (!response.ok) {
            throw new Error(`Error ${response.status}: ${response.statusText}`);
        }
        if (!response.ok) {
            if (response.status === 401) {  
                window.location.href = '/login';
            }
            const errorText = await response.text();
            throw new Error(`Error ${response.status}: ${errorText}`);
        }
        const data = await response.json();
        console.log(data);
        // Actualizar la interfaz con los datos recibidos
        actualizarInterfaz(data);
        
        mostrarMensaje('Datos cargados correctamente', 'success');
    } catch (error) {
        console.error('Error al obtener datos de validación:', error);
        mostrarMensaje('Error al cargar los datos: ' + error.message, 'error');
    }
}

//===================================== UPDATE DATA INTERFACE =========================================
function actualizarInterfaz(data) {
    
    // Actualizar información de la empresa
    const inputline = '<div class="mb-4"><label class="info-label">{1}</label><input type="text" id="{0}" class="input-field" placeholder="{1}" value="{2}"></div>';
    
    document.getElementById('empresaInfo').innerHTML = `
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            
                <div>
                    <p class="info-label">Razón Social</p>
                    <p class="info-value">${data.legalName?data.legalName:'No especificado'}</p>
                </div>
                <div>
                    <p class="info-label">Alias</p>
                    <p class="info-value">${data.aliasName?data.aliasName:'No especificado'}</p>
                </div>
                <div>
                    <p class="info-label">RFC</p>
                    <p class="info-value">${data.taxNum?data.taxNum:'No especificado'}</p>
                </div>
                <div>
                    <p class="info-label">Representante Legal</p>
                    <p class="info-value">${data.legalSignupName?data.legalSignupName + ' ' + (data.legalSignupLastname || '') :'No especificado'}</p>
                </div>
            
        </div>
        
        <div class="mt-6 grid grid-cols-1 md:grid-cols-2 gap-6">
            ${inputline.format("street-address", "Calle y Número", data.streetAddress || '')}
            ${inputline.format("neighborhood", "Colonia", data.neighborhood || '')}
        </div>
        <div class="space-y-4">
            ${inputline.format("locality", "Delegación o Municipio", data.locality || '')}
            ${inputline.format("postal-code", "Código Postal", data.postalCode || '')}
        </div>
        
        <div class="mt-8 pt-6 border-t border-white/10">
            <div class="flex justify-end">
                <button class="btn-cyber px-6 py-2 rounded-lg flex items-center gap-2" id="btnGuardarInfo">
                    <span class="material-symbols-outlined text-sm">save</span>
                    Guardar Información
                </button>
            </div>
        </div>
    `;
    
    // Configurar evento para guardar información
    document.getElementById('btnGuardarInfo').addEventListener('click', guardarInformacionTexto);
    
    for (const [key, value] of Object.entries(documentTypes)) {
        if (data[value]) {
            // Documento ya existe, bloquear la subida
            documentos[key].subido = true;
            documentos[key].nombre = data[value];

            // Construir IDs dinámicos
            const capitalized = key.charAt(0).toUpperCase() + key.slice(1);

            // Actualizar UI
            const statusEl = document.getElementById(`status${capitalized}`);
            statusEl.className = 'status status-completo';
            statusEl.textContent = 'Completo';

            // Cambiar el botón
            const buttonEl = document.getElementById(`btn${capitalized}`);
            if (buttonEl) {
                buttonEl.textContent = 'Documento subido';
                buttonEl.classList.remove('btn-cyber');
                buttonEl.classList.add('btn-secondary');
                buttonEl.disabled = true;
            }

            // Deshabilitar el input file real
            document.getElementById(`file${capitalized}`).disabled = true;

            // Mostrar información del archivo
            document.getElementById(`info${capitalized}`).textContent = `Archivo: ${data[value]}`;
        }
    }
    
    // Verificar si todos los documentos están completos
    verificarCompletitud();
}

// Agregar esta función para recopilar datos de los inputs
function recopilarDatosTexto() {
    const datos = {};
    
    // Obtener valores de todos los inputs de texto
    const inputs = document.querySelectorAll('#empresaInfo input');
    inputs.forEach(input => {
        switch(input.id) {
            case 'street-address':
                datos.streetAddress = input.value;
                break;
            case 'postal-code':
                datos.postalCode = input.value;
                break;
            case 'neighborhood':
                datos.neighborhood = input.value;
                break;
            case 'locality':
                datos.locality = input.value;
                break;
        }
    });
    
    return datos;
}

//===================================== UPLOAD DOCS =========================================
function configurarEventListeners() {
    for (const key of Object.keys(documentTypes)) {
        const input = document.getElementById(`file${key.charAt(0).toUpperCase() + key.slice(1)}`);
        input.addEventListener('change', (event) => {
            subirDocumento(key, event.target.files[0]);
        });
    }
}

// Función para subir un documento
async function subirDocumento(tipo, archivo) {
    console.log(tipo);
    console.log(archivo);
    if (!archivo) return;
    
    try {
        mostrarMensaje(`Subiendo ${tipo}...`, 'info');
        
        const formData = new FormData();
        formData.append('documentType', documentTypes[tipo]);
        formData.append('document', archivo); // Este es el key que busca Go
        
        // Debug: ver qué se está enviando
        console.log('DocumentType value:', documentTypes[tipo]);
        console.log('Archivo name:', archivo.name);
        
        // También enviar datos de texto si existen
        const datosTexto = recopilarDatosTexto();
        for (const key in datosTexto) {
            if (datosTexto[key]) {
                formData.append(key, datosTexto[key]);
            }
        }
        
        const response = await fetch('/updatevaldata', {
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
        console.log('Resultado:', result);
        mostrarMensaje('Documento subido exitosamente', 'success');
        
        // Actualizar estado del documento
        documentos[tipo].subido = true;
        documentos[tipo].archivo = archivo;
        documentos[tipo].nombre = archivo.name;
        
        // Actualizar UI
        const capitalized = tipo.charAt(0).toUpperCase() + tipo.slice(1);
        document.getElementById(`status${capitalized}`).className = 'status status-completo';
        document.getElementById(`status${capitalized}`).textContent = 'Completo';
        
        const buttonEl = document.getElementById(`btn${capitalized}`);
        buttonEl.textContent = 'Documento subido';
        buttonEl.classList.remove('btn-cyber');
        buttonEl.classList.add('btn-secondary');
        buttonEl.disabled = true;
        
        document.getElementById(`file${capitalized}`).disabled = true;
        document.getElementById(`info${capitalized}`).textContent = `Archivo: ${archivo.name}`;
        
        // Verificar si todos los documentos están completos
        verificarCompletitud();
        
        mostrarMensaje(`${tipo.charAt(0).toUpperCase() + tipo.slice(1)} subido correctamente`, 'success');
    } catch (error) {
        console.error(`Error al subir ${tipo}:`, error);
        mostrarMensaje(`Error al subir ${tipo}: ${error.message}`, 'error');
        
        // Limpiar el input de archivo
        document.getElementById(`file${tipo.charAt(0).toUpperCase() + tipo.slice(1)}`).value = '';
    }
}

// Función para guardar solo información de texto
async function guardarInformacionTexto() {
    try {
        mostrarMensaje('Guardando información...', 'info');
        
        const datos = recopilarDatosTexto();
        const response = await fetch('/updatevaldata', {
            method: 'POST',
            credentials: 'include',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(datos)
        });
        
        if (!response.ok) {
            if (response.status === 401) {  
                window.location.href = '/login';
            }
            const errorText = await response.text();
            throw new Error(`Error ${response.status}: ${errorText}`);
        }
        
        mostrarMensaje('Información guardada correctamente', 'success');
    } catch (error) {
        console.error('Error al guardar información:', error);
        mostrarMensaje('Error al guardar información: ' + error.message, 'error');
    }
}

//===================================== CHECK COMPLETE =========================================
function verificarCompletitud() {
    const todosCompletos = Object.values(documentos).every(doc => doc.subido);
    document.getElementById('btnCompletar').disabled = !todosCompletos;
}

// Función para completar la validación
async function completarValidacion() {
    try {
        mostrarMensaje('Completando validación...', 'info');
        
        const response = await fetch('/completevalidation', {
            method: 'POST',
            body: JSON.stringify({ status: true }),
            headers: {
                'Content-Type': 'application/json'
            },
            credentials: 'include'
        });
        
        if (!response.ok) {
            if (response.status === 401) {  
                window.location.href = '/login';
            }
            const errorText = await response.text();
            throw new Error(`Error ${response.status}: ${errorText}`);
        }
        
        const result = await response.json();
        
        mostrarMensaje('Validación completada con éxito', 'success');
        
        // Redirigir o mostrar mensaje de éxito
        setTimeout(() => {
            window.location.href = '/waitapprove'; // Cambiar por la URL correcta
        }, 2000);
    } catch (error) {
        console.error('Error al completar validación:', error);
        mostrarMensaje('Error al completar validación: ' + error.message, 'error');
    }
}

//===================================== SHOW MESSAGES =========================================
function mostrarMensaje(mensaje, tipo) {
    const messageContainer = document.getElementById('messageContainer');
    messageContainer.innerHTML = `<div class="message-${tipo}">${mensaje}</div>`;
    
    // Auto-ocultar después de 5 segundos
    setTimeout(() => {
        messageContainer.innerHTML = '';
    }, 5000);
}

String.prototype.format = function () {
    const args = arguments;
    return this.replace(/{(\d+)}/g, function (match, number) {
        return typeof args[number] !== 'undefined' ? args[number] : match;
    });
};