// --- Variables Globales ---
let userKeys = [];
let selectedKeyId = null;
let currentInvite = null;
let currentFolder = null;
let aut = sessionStorage.getItem('aut');

let faceModel = null;
let faceVector = null;

let videoStream = null;
let isProcessing = false;
let currentFacingMode = 'user';
let isBiometryCaptured = false;
let isSamePerson =false;
// --- Inicialización ---
document.addEventListener('DOMContentLoaded', async () => {
    const urlParams = new URLSearchParams(window.location.search);
    const inviteId = urlParams.get('idInvite');

    if (!inviteId) {
        mostrarMensaje('Error: No se encontró idInvite en la URL', 'error');
        return;
    }

    await loadInviteData(inviteId);
});

/* ============================================================
MODELO FACIAL
============================================================ */
async function loadFaceModel() {
    try {
        console.log('⏳ Cargando modelo facial...');
        faceModel = await tf.loadGraphModel('/facevector/model.json');
        console.log('✅ Modelo facial cargado');
    } catch (err) {
        console.error(err);
        sfAlert('Error cargando el modelo facial');
    }
}

/* ============================================================
PREPROCESAMIENTO (105,105,3)
============================================================ */
function preprocessCanvas(canvas) {
    return tf.tidy(() =>
        tf.browser.fromPixels(canvas)
            .resizeBilinear([105, 105])
            .toFloat()
            .div(255.0)
            .expandDims(0) // [1,105,105,3]
    );
}

/* ============================================================
EXTRACCIÓN DEL VECTOR FACIAL (256)
============================================================ */
async function extractFaceVector(canvas) {
    if (!faceModel) throw new Error('Modelo no cargado');

    const input = preprocessCanvas(canvas);
    const embedding = faceModel.predict(input);
    const vector = await embedding.data(); // Float32Array(256)
    tf.dispose([input, embedding]);

    return Array.from(vector);
}

async function compareFaces(vector) {
    if (!faceVector) throw new Error('No se ha tomado una fotografía');
    try {
        const res = await fetch('/validateFace', {
            method: 'POST',
            credentials: 'include',
            headers: {'Content-Type':'application/json'},
            body: JSON.stringify({ facevector: vector })
        });
        
        if (!res.ok) {
            if (res.status === 401) {  
                window.location.href = '/login';
            }
            const errorText = await res.text();
            throw new Error(`Error ${res.status}: ${errorText}`);
        }
        const response = await res.json();
        mostrarMensaje("Acc", response.distance);
        return response.match;
    } catch (err) {
        mostrarMensaje(err.message, 'error');
        return false;
    }
}

/* ============================================================
CAMARA
============================================================ */
function setupCameraUI() {
    document.getElementById('btnCancelCapture').onclick = stopCamera;
    document.getElementById('btnCapture').onclick = capturePhoto;
    document.getElementById('btnSwitchCamera').onclick = switchCamera;
}

async function startCamera() {
    try {
        const video = document.getElementById('videoElement');
        document.getElementById('photoOverlay').style.display = 'flex';
        document.body.style.overflow = 'hidden';

        videoStream = await navigator.mediaDevices.getUserMedia({
            video: { facingMode: currentFacingMode },
            audio: false
        });

        video.srcObject = videoStream;
    } catch (err) {
        sfAlert('No se pudo acceder a la cámara');
        stopCamera();
    }
}

function stopCamera() {
    if (videoStream) {
        videoStream.getTracks().forEach(t => t.stop());
        videoStream = null;
    }

    document.getElementById('photoOverlay').style.display = 'none';
    document.body.style.overflow = 'auto';
}

function switchCamera() {
    currentFacingMode = currentFacingMode === 'user' ? 'environment' : 'user';
    stopCamera();
    setTimeout(startCamera, 300);
}

/* ============================================================
CAPTURA + VECTOR
============================================================ */
async function capturePhoto() {
    if (isProcessing) return;
    isProcessing = true;

    try {
        const video = document.getElementById('videoElement');
        const canvas = document.createElement('canvas');
        const ctx = canvas.getContext('2d');

        canvas.width = 105;
        canvas.height = 105;
        ctx.drawImage(video, 0, 0, 105, 105);

        faceVector = await extractFaceVector(canvas);
        isSamePerson = await compareFaces(faceVector);
        
        // Marcar biometría como capturada
        if (isSamePerson){
            isBiometryCaptured = true;
            updateBiometryButton();
            mostrarMensaje('Prueba de vida capturada exitosamente', 'success');
        } else {
            mostrarMensaje('Coincidencia facial insuficiente. Intente de nuevo.', 'error');
        }
        stopCamera();
    } catch (err) {
        console.error(err);
        sfAlert('Error procesando el rostro');
    } finally {
        isProcessing = false;
    }
}

/* ============================================================
ACTUALIZAR BOTÓN DE BIOMETRÍA
============================================================ */
function updateBiometryButton() {
    const headerBio = document.getElementById('headerBio');
    if (!headerBio) return;
    
    if (isBiometryCaptured) {
        headerBio.innerHTML = `
            <button class="gem-button captured" onclick="startCamera()" title="Prueba de vida capturada - Click para recapturar">
                <span class="material-symbols-outlined text-white text-sm">check</span>
            </button>
        `;
    } else {
        headerBio.innerHTML = `
            <button class="gem-button" onclick="startCamera()" title="Capturar prueba de vida">Requerida</button>
        `;
    }
}

// --- Funciones de Datos ---

async function loadInviteData(inviteId) {
    try {
        const response = await fetch('/getinvite', {
            method: 'POST',
            credentials: 'include',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ idInvite: inviteId })
        });
        
        if (!response.ok) {
            if (response.status === 401) {  
                window.location.href = '/login';
            }
            const errorText = await response.text();
            throw new Error(`Error ${response.status}: ${errorText}`);
        }
        const rawData = await response.json();
        
        if (!rawData || !rawData.folder) throw new Error('Estructura de datos inválida');
        
        currentFolder = rawData.folder;
        currentInvite = currentFolder.invites.find(i => i.idInvite === inviteId) || currentFolder.invites[0];
        
        // Cargar modelo facial si se requiere biometría
        if (currentInvite.userDest.alive_proof) {
            await loadFaceModel();
            setupCameraUI();
        }
        
        // 1. Renderizar Header (Info de Emisor e Invitación)
        renderInviteHeader(currentFolder, currentInvite);

        // 2. Renderizar Grid de Documentos
        renderDocumentCards(currentInvite.invitedocs);

        // 3. Cargar y Renderizar Llaves
        await loadUserKeys();

        // 4. Configurar Botones de Acción (Firmar/Aceptar)
        updateActionButtons();

    } catch (err) {
        console.error(err);
        mostrarMensaje('Error cargando invitación: ' + err.message, 'error');
    }
}

async function loadUserKeys() {
    try {
        const response = await fetch('/statusk', { method: 'GET', credentials: 'include' });
        if (response.ok) {
            const keysData = await response.json();
            userKeys = keysData || [];
            const selectedKey = userKeys.find(k => k.selected);
            if (selectedKey) selectedKeyId = selectedKey.idKey;
        } else {
            userKeys = [];
            if (response.status === 401) {  
                window.location.href = '/login';
            }
            const errorText = await response.text();
            throw new Error(`Error ${response.status}: ${errorText}`);
        }
        renderSidebarKeys(userKeys);
    } catch (err) {
        console.error(err);
        document.getElementById('sidebarKeys').innerHTML = '<p class="text-red-400 text-xs">Error cargando llaves</p>';
    }
}

// --- Funciones de Renderizado ---

function renderInviteHeader(folder, invite) {
    // Datos del Emisor
    document.getElementById('headerEmisorName').textContent = folder.userEmisor.nameUserEmisor || 'Desconocido';
    document.getElementById('headerInstName').textContent = folder.userEmisor.nameInstEmisor || 'Organización no especificada';
    document.getElementById('headerIdInvite').textContent = `ID: ${invite.idInvite}`;

    // Fechas
    const sentDate = new Date(invite.sentAt);
    document.getElementById('headerSentDate').textContent = !isNaN(sentDate) ? sentDate.toLocaleDateString('es-MX', { day: 'numeric', month: 'long', year: 'numeric' }) : '-';

    // Cálculo de fecha de expiración más próxima
    let earliestExp = null;
    if (invite.invitedocs) {
        earliestExp = invite.invitedocs.reduce((acc, curr) => {
            const d = new Date(curr.expirationDate);
            return (!acc || d < acc) ? d : acc;
        }, null);
    }
    document.getElementById('headerExpDate').textContent = earliestExp ? earliestExp.toLocaleDateString('es-MX') : 'Sin vigencia';
    
    // Biometría
    const bioElem = document.getElementById('headerBio');
    const bioIcon = document.getElementById('iconBio');
    if (invite.userDest.alive_proof) {
        // Mostrar botón de gema en lugar de texto
        bioElem.innerHTML = `
            <button class="gem-button" onclick="startCamera()" title="Capturar prueba de vida">Requerida</button>
        `;
        bioIcon.classList.add("text-blue-400");
    } else {
        bioElem.textContent = "No requerida";
    }

    // Contador circular (Total Docs vs Firmados - lógica simple por ahora 0/Total)
    const totalDocs = invite.invitedocs ? invite.invitedocs.length : 0;
    document.getElementById('docCountDisplay').textContent = `0/${totalDocs}`;
}

function renderDocumentCards(invitedocs) {
    const grid = document.getElementById('documentsGrid');
    grid.innerHTML = '';

    if (!invitedocs || invitedocs.length === 0) {
        grid.innerHTML = `<div class="col-span-full text-center text-gray-500">No hay documentos disponibles.</div>`;
        return;
    }

    invitedocs.forEach(item => {
        const doc = item.document;
        const isSign = item.forSign;
        const idSignature = item.idSign;
        // Determinar icono según extensión
        let iconName = 'description';
        let iconColor = 'text-gray-400';
        if(doc.documentExt === 'pdf') {iconName = 'visibility'; iconColor = 'text-blue-500';}
        else {iconName = 'indeterminate_question_box'; iconColor = 'text-orange-400';}

        const card = document.createElement('div');
        card.className = 'tech-card p-5 group';
        card.innerHTML = `
            <div class="flex justify-between items-start mb-4">
                <div class="p-3 bg-white/5 rounded-lg group-hover:bg-white/10 transition-colors flex items-center gap-2">
                    <div class="p-1 bg-white/20 rounded-md flex items-center">
                        <button onclick="viewDocument('${doc.documentFullName}', '${doc.idDocument}')">
                            <span class="material-symbols-outlined ${iconColor} text-2xl shrink-0">${iconName}</span>
                        </button>
                    </div>
                    
                    <h4 class="text-lg font-semibold text-white line-clamp-2" title="${doc.documentFullName}">${doc.documentName || doc.documentFullName}</h4>
                </div>
            </div>
            <div class="text-sm text-gray-400 mb-4 p-2 h-10 overflow-y-auto h-[15vh]">
                ${doc.abstract || item.comment || 'Sin descripción disponible para este documento.'}
                <div class="absolute bottom-0 left-0 right-0 h-4 bg-gradient-to-t from-[#151c2d] to-transparent"></div>
            </div>

            <div class="mt-auto pt-4 border-t border-white/5 flex items-center justify-between">
                <div class="text-xs text-gray-500">
                    ${(doc.documentExt || 'FILE').toUpperCase()} • ${(item.document.documentSize || '1.2 MB')} 
                </div>
                
                ${isSign 
                    ? idSignature===""? `<span class="badge-neon badge-sign-req"><span class="material-symbols-outlined text-[10px]">edit</span>Firmar</span>`
                    : `<a href="https://larue-unchargeable-aden.ngrok-free.dev/viewSignature?id=${idSignature}" class="badge-neon badge-info-only"><span class="material-symbols-outlined text-[10px]">visibility</span>Ver Firma</a>`:
                    `<span class="badge-neon badge-info-only"><span class="material-symbols-outlined text-[10px]">visibility</span>Info</span>`
                }
            </div>
        `;
        grid.appendChild(card);
    });
}

function renderSidebarKeys(keys) {
    const container = document.getElementById('sidebarKeys');
    container.innerHTML = '';

    if (keys.length === 0) {
        container.innerHTML = `
            <div class="text-center py-6 border border-dashed border-gray-700 rounded-xl">
                <span class="material-symbols-outlined text-gray-600 text-3xl mb-2">no_accounts</span>
                <p class="text-sm text-gray-400">No hay llaves disponibles</p>
            </div>`;
        return;
    }

    keys.forEach(k => {
        const isSelected = k.selected || k.idKey === selectedKeyId;
        const isExpired = new Date(k.expiration) < new Date();
        
        const div = document.createElement('div');
        div.className = `key-card-item ${isSelected ? 'selected' : ''}`;
        div.onclick = () => selectKey(k.idKey);
        
        div.innerHTML = `
            <div class="flex items-center gap-3">
                <div class="size-8 rounded-full bg-gradient-to-br ${isSelected ? 'from-blue-500 to-blue-700' : 'from-gray-700 to-gray-800'} flex items-center justify-center text-xs font-bold text-white">
                    ${k.owner.charAt(0)}
                </div>
                <div class="overflow-hidden">
                    <p class="text-sm font-bold text-gray-200 truncate w-40">${k.owner}</p>
                    <p class="text-[10px] text-gray-500">Expira: ${new Date(k.expiration).toLocaleDateString('es-MX')}</p>
                </div>
            </div>
            ${isExpired ? '<div class="absolute bottom-2 right-2 text-[10px] text-red-400 font-bold">EXPIRADA</div>' : ''}
        `;
        container.appendChild(div);
    });
}

// --- Lógica de Negocio ---

async function selectKey(idKey) {
    document.querySelectorAll('.key-card-item').forEach(el => el.classList.remove('selected'));
    
    try {
        const res = await fetch('/updatek', {
            method: 'POST',
            credentials: 'include',
            headers: {'Content-Type':'application/json'},
            body: JSON.stringify({ idKeyUpdate: idKey })
        });
        
        if (!res.ok) {
            if (res.status === 401) {  
                window.location.href = '/login';
            }
            const errorText = await res.text();
            throw new Error(`Error ${res.status}: ${errorText}`);
        }
        
        selectedKeyId = idKey;
        await loadUserKeys(); 
        mostrarMensaje('Llave seleccionada', 'success');
        logoutk();
        
    } catch (err) {
        mostrarMensaje(err.message, 'error');
    }
}

function updateActionButtons() {
    const container = document.getElementById('sidebarActions');
    const hasSignDocs = currentInvite.invitedocs.some(d => d.forSign);

    if (hasSignDocs) {
        container.innerHTML = `
            <button onclick="signAllDocuments()" class="btn-cyber w-full py-3 rounded-xl text-white font-bold flex items-center justify-center gap-2 shadow-lg shadow-blue-900/20">
                <span class="material-symbols-outlined">ink_pen</span>
                Firmar Documentos
            </button>
            <p class="text-[10px] text-center text-gray-500 mt-2">
                Se aplicará tu firma digital a todos los documentos requeridos.
            </p>
        `;
    } else {
        container.innerHTML = `
            <button onclick="acceptInvite()" class="w-full py-3 rounded-xl border border-green-500/30 bg-green-500/10 text-green-400 font-bold hover:bg-green-500/20 transition-all flex items-center justify-center gap-2">
                <span class="material-symbols-outlined">check_circle</span>
                Dar Visto Bueno
            </button>
        `;
    }
}

// 1. Función principal llamada por el botón "Firmar"
function signAllDocuments() {
    // Validaciones básicas
    if (!selectedKeyId) {
        mostrarMensaje('Debes seleccionar una llave primero', 'error');
        return;
    }
    
    // Validar biometría si es requerida
    if (currentInvite.userDest.alive_proof) {
        if (!isBiometryCaptured || !faceVector || !isSamePerson || faceVector.length !== 256) {
            mostrarMensaje('Debes capturar la prueba de vida primero', 'error');
            startCamera();
            return;
        }
    }

    // REVISAR SI EXISTE EL TOKEN (aut)
    let currentAut = sessionStorage.getItem('aut');

    if (!currentAut || currentAut === 'undefined' || currentAut === '') {
        // Si no hay token, pedimos password
        openAuthModal();
    } else {
        // Si hay token, firmamos directo
        executeSigning();
    }
}

// 4. Lógica final de firma (Recibe el token 'aut' como parámetro)
async function executeSigning() {
    const signDocs = currentInvite.invitedocs
        .filter(d => d.forSign)
        .map(d => ({
            idDocument: d.document.idDocument,
            documentHash: d.document.documentHash
        }));

    mostrarMensaje('Firmando documentos...', 'info');
    try {
        const headers= {
            'Content-Type': 'application/json',
            'X-Signer-Timezone': Intl.DateTimeFormat().resolvedOptions().timeZone,
            'X-Signer-Browser': navigator.userAgentData?.brands
                ?.map(b => `${b.brand} ${b.version}`)
                .join(', ') || navigator.userAgent
            }
        const body = {
            inviteId: currentInvite.idInvite,
            keyId: selectedKeyId,
            folderId: currentFolder.idFolder,
            signDocuments: signDocs,
        };

        const res = await fetch('/signDocument', {
            method: 'POST',
            headers: headers,
            body: JSON.stringify(body),
            credentials: 'include'
        });

        const data = await res.json();

        if (res.status === 200 || res.status === 201) {
            
            // Verificamos que signed exista y sea un objeto
            if (data.signed) {
                Object.entries(data.signed).forEach(([idDoc, path]) => {
                    console.log(`ID: ${idDoc}, Path: ${path}`);
                    mostrarMensaje(`Documento firmado exitosamente.`, 'success'); 
                });
                
                const totalFirmados = Object.keys(data.signed).length;
                mostrarMensaje(`Proceso finalizado. Total firmados: ${totalFirmados}`, 'success');
            }

        } else {
            if (res.status === 401 || res.status === 403 || (data.message && data.message.includes('token'))) {
                logoutk();
                mostrarMensaje('Tu sesión expiró. Intenta firmar de nuevo.', 'error');
                const errorText = await res.text();
                if (res.status === 401) {  
                    window.location.href = '/login';
                }
                throw new Error(`Error ${res.status}: ${errorText}`);
            } else {
                mostrarMensaje('Error: ' + (data.message || 'Error desconocido'), 'error');
            }
        }
    } catch (err) {
        console.error(err);
        mostrarMensaje('Error de conexión al firmar', 'error');
    }
}

// 2. Funciones del Modal
function openAuthModal() {
    const modal = document.getElementById('authModal');
    const input = document.getElementById('authPasswordInput');
    modal.classList.remove('hidden');
    setTimeout(() => input.focus(), 100);
}

function closeAuthModal() {
    document.getElementById('authModal').classList.add('hidden');
    document.getElementById('authPasswordInput').value = '';
}

// 3. Manejar el envío del formulario (Obtener Token)
async function handleAuthSubmit(e) {
    e.preventDefault();
    const password = document.getElementById('authPasswordInput').value;
    const btnSubmit = document.getElementById('btnAuthSubmit');
    const originalText = btnSubmit.innerHTML;

    if (!password) {
        mostrarMensaje('Ingresa tu contraseña', 'error');
        return;
    }

    // UI Loading state
    btnSubmit.disabled = true;
    btnSubmit.innerHTML = `<div class="spinner w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin"></div> Procesando...`;

    try {
        const response = await fetch('/logink', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ password: password }),
            credentials: 'include'
        });

        const data = await response.json();

        if (response.ok && data.token) {
            sessionStorage.setItem('aut', data.token);
            
            closeAuthModal();
            mostrarMensaje('Autorización exitosa', 'success');
            executeSigning(data.token);
        } else {
            logoutk();
            if (response.status === 401) {  
                window.location.href = '/login';
            }
            throw new Error(data.message || 'Contraseña incorrecta');
        }

    } catch (err) {
        console.error(err);
        mostrarMensaje(err.message || 'Error de autenticación', 'error');
    } finally {
        btnSubmit.disabled = false;
        btnSubmit.innerHTML = originalText;
    }
}

async function acceptInvite() {
    var _ok3 = await sfConfirm({title:'Visto bueno',message:'¿Confirmar visto bueno en este documento?',type:'success',confirmText:'Confirmar',confirmClass:'sf-modal-btn-confirm'});if(!_ok3) return;
    try {
        const res = await fetch('/acceptDocument', { 
            method: 'POST', 
            headers: {'Content-Type':'application/json'}, 
            body: JSON.stringify({ 
                inviteId: currentInvite.idInvite, 
                folderId: currentFolder.idFolder 
            }) 
        });
        if (!res.ok) {
            if (res.status === 401) {  
                window.location.href = '/login';
            }
            const errorText = await res.text();
            throw new Error(`Error ${res.status}: ${errorText}`);
        }
        const j = await res.json();
        if (j.success) { 
            mostrarMensaje('Visto bueno registrado', 'success'); 
            setTimeout(() => location.reload(), 1500);
        } else {
                mostrarMensaje('Error al aceptar', 'error');
        }
    } catch(e) { console.error(e); }
}

// --- Utilidades ---

function mostrarMensaje(msg, type) {
    const container = document.getElementById('messageContainer');
    const colors = {
        error: 'bg-red-500/10 border-red-500/50 text-red-400',
        success: 'bg-green-500/10 border-green-500/50 text-green-400',
        info: 'bg-blue-500/10 border-blue-500/50 text-blue-400'
    };
    
    container.className = `mb-4 px-4 py-3 rounded-lg border text-xs font-medium flex items-center gap-2 ${colors[type] || colors.info} animate-fade-in`;
    container.innerHTML = `
        <span class="material-symbols-outlined text-base">info</span>
        ${msg}
    `;
    container.style.display = 'flex';
    
    setTimeout(() => {
        container.style.display = 'none';
    }, 5000);
}

async function viewDocument(name, docId) {
    if (currentInvite.userDest.alive_proof) {
        if (!isBiometryCaptured || !faceVector || !isSamePerson || faceVector.length !== 256) {
            mostrarMensaje('Debes capturar la prueba de vida primero', 'error');
            return;
        }
    }
    const modal = document.getElementById('documentModal');
    const viewer = document.getElementById('documentViewer');
    const loading = document.getElementById('modalLoading');
    
    document.getElementById('modalTitle').innerHTML = `<span class="material-symbols-outlined text-gray-400">description</span> ${name}`;
    modal.classList.remove('hidden');
    viewer.classList.add('hidden');
    loading.style.display = 'flex';

    try {
        const docs = {idDocs: [docId], type: "uploaded"}
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
        loading.innerHTML = `<p class="text-red-400">Error cargando documento</p>`;
    }
}

async function logoutk() {
    try {
        const response = await fetch('/logoutk', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ idKey: selectedKeyId }),
            credentials: 'include'
        });

        const data = await response.json();

        if (response.ok) {
            sessionStorage.removeItem('aut');
            mostrarMensaje('Cierre de sesión de firma', 'success');
        } else {
            throw new Error(data.message || 'Error al cerrar sesión');
        }

    } catch (err) {
        console.error(err);
        mostrarMensaje(err.message || 'Error de cierre de sesión', 'error');
    }
}

function closeModal() {
    document.getElementById('documentModal').classList.add('hidden');
    document.getElementById('documentViewer').src = '';
}

// Función mock para subir llaves (puedes reutilizar tu lógica original aquí)
function showKeyUploadForm() {
    const sidebarActions = document.getElementById('sidebarActions');
    sidebarActions.innerHTML = `
        <div class="w-full">
            <form id="credentialsForm" class="space-y-4" autocomplete="off" novalidate>
                <div>
                    <label class="block text-sm text-secondary mb-1">Llave privada (.pem/.key)</label>
                    <input type="file" id="keyFile" accept=".key,.pem,.der" class="w-full p-2 bg-white/10 rounded text-sm text-primary">
                </div>
                <div>
                    <label class="block text-sm text-secondary mb-1">Certificado (.cer/.crt/.pem)</label>
                    <input type="file" id="certFile" accept=".cer,.crt,.pem,.der" class="w-full p-2 bg-white/10 rounded text-sm text-primary">
                </div>
                <div>
                    <label class="block text-sm text-secondary mb-1">Contraseña</label>
                    <input type="password" id="keyPassword" class="w-full p-2 bg-white/10 rounded text-sm text-primary">
                </div>
                <button class="futuristic-btn w-full" type="submit">
                    <span class="material-symbols-outlined text-sm mr-2">cloud_upload</span>
                    Subir Llaves
                </button>
            </form>
        </div>
    `;

    document.getElementById('credentialsForm').addEventListener('submit', async (e) => {
        e.preventDefault();
        mostrarMensaje('Subiendo llaves...', 'info');
        const keyFile = document.getElementById('keyFile').files[0];
        const certFile = document.getElementById('certFile').files[0];
        const keyPassword = document.getElementById('keyPassword').value;

        if (!keyFile || !certFile || !keyPassword) {
        sfAlert('Por favor, sube la llave, el certificado y escribe la contraseña.');
        return;
        }

        try {
        const formData = new FormData();
        formData.append('keyFile', keyFile);
        formData.append('certFile', certFile);
        formData.append('passKey', keyPassword);

        const response = await fetch('/uploadk', {
            method: 'POST',
            body: formData,
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
        console.log(result);
        if (result.valid) {
            await loadUserKeys();
            mostrarMensaje('Llaves subidas correctamente', 'success');
        } else {
            alert('Error: ' + result.message);
        }
        } catch (err) {
            console.error('Error subiendo archivos:', err);
            alert('Error al subir los archivos: ' + err.message);
        }
    });
}