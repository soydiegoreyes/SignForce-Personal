// Variables globales
let currentSignatureData = null;
let currentDocumentId = null;

// Inicialización
document.addEventListener('DOMContentLoaded', async () => {
    const urlParams = new URLSearchParams(window.location.search);
    const signatureId = urlParams.get('id');

    if (!signatureId) {
        showError('Error: No se encontró ID de firma en la URL');
        return;
    }

    await loadSignatureData(signatureId);
});

// Cargar datos de la firma
async function loadSignatureData(signatureId) {
    try {
        const response = await fetch(`/viewSign?id=${encodeURIComponent(signatureId)}`, {
            method: 'GET',
            credentials: 'include'
        });
        
        if (!response.ok) throw new Error('Error al obtener datos de la firma');
        
        const data = await response.json();
        currentSignatureData = data;
        currentDocumentId = data.document?.idDocument;
        
        renderSignatureData(data);
        
    } catch (err) {
        console.error(err);
        showError('Error cargando datos de la firma: ' + err.message);
    }
}


// Renderizar datos en la interfaz
function renderSignatureData(data) {
    // ID de la firma
    document.getElementById('sigIdDisplay').textContent = `ID: ${data.idSignature}`;

    // Firmante - Email como enlace
    const signer = data.signer || {};
    document.getElementById('signerName').textContent = signer.name || 'Desconocido';
    document.getElementById('signerStatus').textContent = signer.active ? 'Activo' : 'Inactivo';
    document.getElementById('signerStatus').className = signer.active ? 'info-value text-green-400' : 'info-value text-red-400';
    // Renderizar email del firmante como enlace
    renderEmailLink('signerEmailContainer', signer.email, signer.name || 'Firmante');
    renderWhatsLink('signerPhoneContainer', signer.phone, signer.name || 'Firmante');
    
    // Propietario - Email como enlace
    const owner = data.owner || {};
    document.getElementById('ownerName').textContent = owner.name || 'Desconocido';
    document.getElementById('ownerStatus').textContent = owner.active ? 'Activo' : 'Inactivo';
    document.getElementById('ownerStatus').className = owner.active ? 'info-value text-green-400' : 'info-value text-red-400';
    // Renderizar email del propietario como enlace
    renderEmailLink('ownerEmailContainer', owner.email, owner.name || 'Propietario');
    renderWhatsLink('ownerPhoneContainer', owner.phone, owner.name || 'Propietario');

    // Institución - Email como enlace
    const inst = data.institution || {};
    document.getElementById('instLegalName').textContent = inst.legalName || 'No especificado';
    document.getElementById('instRepresentative').textContent = `${inst.legalSignupName || ''} ${inst.legalSignupLastname || ''}`.trim() || 'No especificado';
    // Renderizar email de la institución como enlace
    renderEmailLink('instEmailContainer', inst.email || inst.contactEmailInst, inst.legalName || 'Institución');
    renderWhatsLink('instPhoneContainer', inst.phone, inst.legalName || 'Institución');

    // Llave de firma
    const keys = data.keys || {};
    document.getElementById('keyOwner').textContent = keys.owner || 'Desconocido';
    document.getElementById('keyUniqueId').textContent = keys.subjectUniqueId || 'No especificado';
    document.getElementById('keyCertName').textContent = keys.nameCer || 'No especificado';
    document.getElementById('keyExpiration').textContent = formatDate(keys.expiration);
    const issuerdetails = document.createElement('div');
    issuerdetails.className = 'key-details';
    issuerdetails.innerHTML = formatIssuer(keys.issuerData);
    document.getElementById('keyCardContainer').appendChild(issuerdetails);

    // Documento
    const doc = data.document || {};
    document.getElementById('docName').textContent = doc.documentName || doc.documentFullName || 'Sin nombre';
    document.getElementById('docType').textContent = (doc.documentExt || 'file').toUpperCase();
    document.getElementById('docSize').textContent = formatFileSize(doc.sizeB);
    document.getElementById('docCreated').textContent = formatDate(doc.createdAtDoc);
    
    const abstractText = doc.abstractDoc || 'No hay resumen disponible para este documento.';
    document.getElementById('docAbstract').textContent = abstractText;
    updateAbstractStats(abstractText);

    // Fecha de firma
    const sig = data.signature || {};
    const sigDate = new Date(sig.genTimeSign);
    document.getElementById('sigDate').textContent = formatDateTime(sigDate);
    document.getElementById('sigTimeAgo').textContent = getTimeAgo(sigDate);

    // Detalles técnicos
    document.getElementById('sigAlgo').textContent = sig.signatureAlgoSign_fk || 'N/A';
    document.getElementById('sigDigestAlgo').textContent = sig.digestAlgoSign_fk || 'N/A';
    document.getElementById('sigDigest').textContent = sig.digestValueSign || 'N/A';
    document.getElementById('sigNonce').textContent = sig.nonceSign || 'N/A';
    document.getElementById('sigValue').textContent = sig.signatureValueSign || 'N/A';
    
}

// Funciones de utilidad
function formatDate(dateString) {
    if (!dateString) return 'N/A';
    const date = new Date(dateString);
    if (isNaN(date)) return dateString;
    return date.toLocaleDateString('es-MX', {
        year: 'numeric',
        month: 'long',
        day: 'numeric'
    });
}

function formatDateTime(date) {
    if (!date || isNaN(date)) return 'N/A';
    return date.toLocaleString('es-MX', {
        year: 'numeric',
        month: 'long',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit',
        second: '2-digit'
    });
}

function formatFileSize(bytes) {
    if (!bytes) return 'N/A';
    const units = ['B', 'KB', 'MB', 'GB'];
    let size = Number(bytes);
    let unitIndex = 0;
    while (size >= 1024 && unitIndex < units.length - 1) {
        size /= 1024;
        unitIndex++;
    }
    return `${size.toFixed(2)} ${units[unitIndex]}`;
}

function formatIssuer(issuerRFC4514) {
    if (!issuerRFC4514) {
        return '';
    }

    // 1. Dividir el string por las comas que separan los atributos
    const attributes = issuerRFC4514.split(',');
    
    // 2. Procesar cada atributo
    const budgetsHTML = attributes.map(attribute => {
        // Expresión Regular para encontrar y eliminar el prefijo
        // El patrón es: cualquier cosa que no sea "=" (.[^=]*) seguida de "="
        // y se usa replace para dejar solo el valor (la parte derecha del "=").
        // También se limpia el espacio al inicio del valor si existe (como en ", O=...")
        const value = attribute
            .replace(/^.[^=]*=/, '') // Elimina el prefijo como 'CN=' o 'OID.2.5.4.45='
            .trim(); // Elimina espacios en blanco al inicio o al final

        // Opcionalmente, puedes obtener el nombre del atributo (el prefijo) para mostrarlo
        const match = attribute.match(/^.[^=]*=/);
        let name = '';
        if (match) {
            // Limpia el '=' final y, en el caso de OIDs, solo deja OID
            name = match[0].replace('=', '').replace(/OID\..*/, 'OID').trim();
        }

        // 3. Crear el HTML para el "budget" (badge)
        // Usamos una clase de Bootstrap para darle estilo
        // Se puede personalizar la clase según tu necesidad (bg-primary, bg-success, etc.)
        // Concatenamos el nombre del atributo y su valor para mayor contexto
        return `<span class="badge badge-neon bg-secondary me-2 mb-1" title="${name}: ${value}">${value}</span>`;
    }).join(''); // Unir todos los elementos en un solo string

    // 4. Devolver el div contenedor con todos los budgets
    // El 'd-flex flex-wrap' es útil para que se acomoden bien si hay muchos.
    return `<div class="d-flex flex-wrap"><p class="info-label">Emisor de certificado</p>${budgetsHTML}</div>`;
}

function getTimeAgo(date) {
    const now = new Date();
    const diffMs = now - date;
    const diffMins = Math.floor(diffMs / 60000);
    const diffHours = Math.floor(diffMs / 3600000);
    const diffDays = Math.floor(diffMs / 86400000);

    if (diffMins < 1) return 'Hace unos segundos';
    if (diffMins < 60) return `Hace ${diffMins} ${diffMins === 1 ? 'minuto' : 'minutos'}`;
    if (diffHours < 24) return `Hace ${diffHours} ${diffHours === 1 ? 'hora' : 'horas'}`;
    if (diffDays < 30) return `Hace ${diffDays} ${diffDays === 1 ? 'día' : 'días'}`;
    return `Hace ${Math.floor(diffDays / 30)} ${Math.floor(diffDays / 30) === 1 ? 'mes' : 'meses'}`;
}
// Función auxiliar para renderizar emails como enlaces
function renderEmailLink(containerId, email, name) {
    const container = document.getElementById(containerId);
    container.innerHTML = ''; // Limpiar contenido
    
    if (email && email.trim() !== '') {
        // Crear enlace de email con estilo
        const emailLink = document.createElement('a');
        emailLink.href = `mailto:${email}`;
        emailLink.target = '_blank';
        emailLink.className = 'inline-flex items-center gap-2 px-3 py-1 rounded-lg bg-blue-500/10 hover:bg-blue-500/20 border border-blue-500/20 hover:border-blue-500/30 text-blue-400 hover:text-blue-300 transition-all duration-200 hover:scale-[1.02] cursor-pointer group';
        emailLink.title = `Enviar correo a ${name} (${email})`;
        
        emailLink.innerHTML = `
            <span class="material-symbols-outlined text-sm group-hover:animate-pulse">mail</span>
            <span class="text-sm font-medium truncate max-w-[200px]">${email}</span>
        `;
        
        container.appendChild(emailLink);
    } else {
        // Mostrar texto si no hay email
        container.innerHTML = '<span class="text-gray-500 italic">No especificado</span>';
    }
}
function renderWhatsLink(containerId, phone, name) {
    const container = document.getElementById(containerId);
    container.innerHTML = ''; // Limpiar contenido
    
    if (phone && phone.trim() !== '') {
        // Crear enlace de email con estilo
        const phoneLink = document.createElement('a');
        phoneLink.href = `https://api.whatsapp.com/send?phone=${phone}`;
        phoneLink.target = '_blank';
        phoneLink.className = 'inline-flex items-center gap-2 px-3 py-1 rounded-lg bg-green-500/10 hover:bg-green-500/20 border border-green-500/20 hover:border-green-500/30 text-green-400 hover:text-green-300 transition-all duration-200 hover:scale-[1.02] cursor-pointer group';
        phoneLink.title = `Enviar mensaje a ${name} (${phone})`;
        
        phoneLink.innerHTML = `
            <span class="material-symbols-outlined text-sm group-hover:animate-pulse">phone</span>
            <span class="text-sm font-medium truncate max-w-[200px]">${phone}</span>
        `;
        
        container.appendChild(phoneLink);
    } else {
        // Mostrar texto si no hay email
        container.innerHTML = '<span class="text-gray-500 italic">No especificado</span>';
    }
}
// Funciones de acción
async function viewDocument() {
    if (!currentDocumentId) {
        showError('No hay documento disponible para visualizar');
        return;
    }

    const modal = document.getElementById('documentModal');
    const viewer = document.getElementById('documentViewer');
    const loading = document.getElementById('modalLoading');
    
    modal.classList.remove('hidden');
    viewer.classList.add('hidden');
    loading.style.display = 'flex';

    try {
        const docs = {idDocs: [currentDocumentId], type: "uploaded"}
        const response = await fetch('/downloadDoc', {
            method: 'POST',
            headers: {'Content-Type': 'application/json'},
            body: JSON.stringify(docs)
        });

        if (!response.ok) {
            if (response.status === 401) {  
                window.location.href = '/login';
            }
            const errorText = await response.text();
            throw new Error(`Error ${response.status}: ${errorText}`);
        }
        
        const blob = await response.blob();
        viewer.src = URL.createObjectURL(blob);
        loading.style.display = 'none';
        viewer.classList.remove('hidden');
    } catch (err) {
        loading.innerHTML = `<p class="text-red-400">Error cargando documento: ${err.message}</p>`;
    }
}

async function downloadDocument() {
    if (!currentDocumentId) {
        showError('No hay documento disponible para descargar');
        return;
    }
    if (!currentSignatureData) {
        showError('no hay datos de la firma');
        return 
    }
    const sig = currentSignatureData.signature || {};
    const spl = sig.pathSign.split("/");
    let idInst = 0;
    let idFolder=0;
    let idDoc=0;
    let hashD = "";
    if (spl[0] === "") {
        idInst = spl[1];
        idFolder = spl[2]; 
        idDoc=spl[3];
        hashD = sig.digestValueSign;
    }
    
    try {
        const query = `/getAsice?f=${encodeURIComponent(idFolder)}&d=${encodeURIComponent(idDoc)}&h=${encodeURIComponent(hashD)}`;
        const response = await fetch(query, {
            method: 'GET',
            credentials: 'include'
        });

        if (!response.ok) throw new Error('Error descargando documento');
        
        const blob = await response.blob();
        const url = URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        
        // Obtener nombre del documento
        const docName = currentSignatureData.document?.documentName || 
                        currentSignatureData.document?.documentFullName || 
                        'documento_firmado';
        a.download = `${docName}.zip`;
        
        document.body.appendChild(a);
        a.click();
        a.remove();
        URL.revokeObjectURL(url);
        
        showMessage('Documento descargado exitosamente', 'success');
        
    } catch (err) {
        showError('Error al descargar documento: ' + err.message);
    }
}

function closeModal() {
    const modal = document.getElementById('documentModal');
    const viewer = document.getElementById('documentViewer');
    
    modal.classList.add('hidden');
    viewer.src = '';
    viewer.classList.add('hidden');
}

function updateAbstractStats(abstractText) {
    if (!abstractText) return;
    
    const charCount = abstractText.length;
    const lineCount = abstractText.split('\n').length;
}

// Funciones de mensajes
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

function showError(message) {
    showMessage(message, 'error');
}