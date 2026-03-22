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
    const v = (val) => val || '<span style="color:rgba(255,255,255,.25);font-style:italic">No especificado</span>';
    const legalRep = [data.legalSignupName, data.legalSignupLastname].filter(Boolean).join(' ') || null;

    document.getElementById('empresaInfo').innerHTML = `
        <div style="display:grid;grid-template-columns:1fr 1fr;gap:.5rem 2rem;margin-bottom:1.5rem">
            <div style="padding:.75rem 0;border-bottom:1px solid rgba(255,255,255,.06)">
                <div style="font-size:.68rem;font-weight:700;color:rgba(255,255,255,.3);text-transform:uppercase;letter-spacing:.08em;margin-bottom:.35rem">Razón Social</div>
                <div style="font-size:.92rem;font-weight:600;color:#fff">${v(data.legalName)}</div>
            </div>
            <div style="padding:.75rem 0;border-bottom:1px solid rgba(255,255,255,.06)">
                <div style="font-size:.68rem;font-weight:700;color:rgba(255,255,255,.3);text-transform:uppercase;letter-spacing:.08em;margin-bottom:.35rem">RFC</div>
                <div style="font-size:.92rem;font-weight:600;color:#fff;font-family:monospace;letter-spacing:.04em">${v(data.taxNum)}</div>
            </div>
            <div style="padding:.75rem 0;border-bottom:1px solid rgba(255,255,255,.06)">
                <div style="font-size:.68rem;font-weight:700;color:rgba(255,255,255,.3);text-transform:uppercase;letter-spacing:.08em;margin-bottom:.35rem">Representante Legal</div>
                <div style="font-size:.92rem;font-weight:600;color:#fff">${v(legalRep)}</div>
            </div>
            <div style="padding:.75rem 0;border-bottom:1px solid rgba(255,255,255,.06)">
                <div style="font-size:.68rem;font-weight:700;color:rgba(255,255,255,.3);text-transform:uppercase;letter-spacing:.08em;margin-bottom:.35rem">Alias</div>
                <div style="font-size:.92rem;font-weight:600;color:#fff">${v(data.aliasName)}</div>
            </div>
        </div>

        <div style="margin-bottom:.5rem;font-size:.72rem;font-weight:700;color:rgba(255,255,255,.25);text-transform:uppercase;letter-spacing:.08em">Domicilio Fiscal</div>
        <div style="display:grid;grid-template-columns:1fr 1fr;gap:.75rem;margin-bottom:.75rem">
            <div style="grid-column:1/-1">
                <input type="text" id="street-address" placeholder="Calle y Número"
                    value="${data.streetAddress || ''}"
                    style="width:100%;padding:.65rem 1rem;border-radius:10px;font-size:.85rem;background:rgba(255,255,255,.04);border:1px solid rgba(255,255,255,.08);color:#fff;font-family:'DM Sans',sans-serif;outline:none;box-sizing:border-box;transition:border-color .2s"
                    onfocus="this.style.borderColor='rgba(181,196,19,.5)'" onblur="this.style.borderColor='rgba(255,255,255,.08)'">
            </div>
            <input type="text" id="neighborhood" placeholder="Colonia"
                value="${data.neighborhood || ''}"
                style="width:100%;padding:.65rem 1rem;border-radius:10px;font-size:.85rem;background:rgba(255,255,255,.04);border:1px solid rgba(255,255,255,.08);color:#fff;font-family:'DM Sans',sans-serif;outline:none;box-sizing:border-box;transition:border-color .2s"
                onfocus="this.style.borderColor='rgba(181,196,19,.5)'" onblur="this.style.borderColor='rgba(255,255,255,.08)'">
            <input type="text" id="locality" placeholder="Delegación o Municipio"
                value="${data.locality || ''}"
                style="width:100%;padding:.65rem 1rem;border-radius:10px;font-size:.85rem;background:rgba(255,255,255,.04);border:1px solid rgba(255,255,255,.08);color:#fff;font-family:'DM Sans',sans-serif;outline:none;box-sizing:border-box;transition:border-color .2s"
                onfocus="this.style.borderColor='rgba(181,196,19,.5)'" onblur="this.style.borderColor='rgba(255,255,255,.08)'">
            <input type="text" id="postal-code" placeholder="Código Postal" maxlength="5"
                value="${data.postalCode || ''}"
                style="width:100%;padding:.65rem 1rem;border-radius:10px;font-size:.85rem;background:rgba(255,255,255,.04);border:1px solid rgba(255,255,255,.08);color:#fff;font-family:'DM Sans',sans-serif;outline:none;box-sizing:border-box;transition:border-color .2s"
                onfocus="this.style.borderColor='rgba(181,196,19,.5)'" onblur="this.style.borderColor='rgba(255,255,255,.08)'">
        </div>
        <button id="btnGuardarInfo"
            style="padding:.55rem 1.25rem;border-radius:9px;font-size:.82rem;font-weight:600;cursor:pointer;border:1px solid rgba(181,196,19,.35);background:rgba(181,196,19,.08);color:#B5C413;font-family:'DM Sans',sans-serif;display:inline-flex;align-items:center;gap:6px;transition:all .2s"
            onmouseover="this.style.background='rgba(181,196,19,.16)'" onmouseout="this.style.background='rgba(181,196,19,.08)'">
            <span class="material-symbols-outlined" style="font-size:15px;">save</span> Guardar domicilio
        </button>
    `;

    // Responsive en pantallas chicas
    if (window.innerWidth < 540) {
        document.querySelectorAll('#empresaInfo input').forEach(el => {
            el.style.gridColumn = '1/-1';
        });
    }

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
        const statusDone = document.getElementById('status' + capitalized);
        if (statusDone) { statusDone.textContent = 'Listo'; statusDone.className = 'doc-status done'; }
        const btnDone = document.getElementById('btn' + capitalized);
        if (btnDone) { btnDone.classList.add('done'); btnDone.innerHTML = '<span class="material-symbols-outlined" style="font-size:15px;">check_circle</span> Subido'; btnDone.style.pointerEvents = 'none'; }
        const fileDone = document.getElementById('file' + capitalized);
        if (fileDone) fileDone.disabled = true;
        const infoDone = document.getElementById('info' + capitalized);
        if (infoDone) infoDone.textContent = archivo.name;
        const itemDone = document.getElementById('docItem' + capitalized);
        if (itemDone) itemDone.classList.add('doc-done');
        
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
    const total = Object.keys(documentos).length;
    const done = Object.values(documentos).filter(d => d.subido).length;
    const todosCompletos = done === total;
    const pct = 25 + Math.round((done / total) * 75);
    const bar = document.getElementById('progressBar');
    const lbl = document.getElementById('progressPct');
    if (bar) bar.style.width = pct + '%';
    if (lbl) lbl.textContent = pct + '%';
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
    const c = document.getElementById('messageContainer');
    const map = { info: 'msg-info', success: 'msg-success', error: 'msg-error' };
    c.innerHTML = '<div class="msg ' + (map[tipo] || 'msg-info') + '">' + mensaje + '</div>';
    if (tipo !== 'error') setTimeout(() => { c.innerHTML = ''; }, 5000);
}

String.prototype.format = function () {
    const args = arguments;
    return this.replace(/{(\d+)}/g, function (match, number) {
        return typeof args[number] !== 'undefined' ? args[number] : match;
    });
};