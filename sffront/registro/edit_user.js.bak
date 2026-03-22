/* ============================================================
VARIABLES GLOBALES
============================================================ */
let existUser = null;
let inviteId = null;
let inviteEmail = null;
let faceModel = null;
let faceVector = null;

let videoStream = null;
let isProcessing = false;
let currentFacingMode = 'user';

const passRegex = /^(?:[^\s',"/$\\%&(){}|[\]\^`]){12,}$/;
const complexityRegex = /^(?=.*[A-Z])(?=.*\d).+$/;

/* ============================================================
INIT
============================================================ */
document.addEventListener('DOMContentLoaded', async () => {
    const params = new URLSearchParams(window.location.search);
    inviteId = params.get('id');

    if (!inviteId) {
        alert('Invitación inválida');
        return;
    }

    await loadFaceModel();
    fetchInviteData(inviteId);
    setupCameraUI();
    ['oldPass', 'newPass', 'confirmPass'].forEach(id => {
        document.getElementById(id).addEventListener('input', validateForm);
    });
});


/* ============================================================
LÓGICA DE CONTRASEÑAS Y UI
============================================================ */
function updatePassFeedback(pass) {
    const reqLength = document.getElementById('reqLength');
    const reqChars = document.getElementById('reqChars');
    const reqComplexity = document.getElementById('reqComplexity');

    // 1. Validar Longitud
    if (pass.length >= 12) {
        setReqStatus(reqLength, true);
    } else {
        setReqStatus(reqLength, false);
    }

    // 2. Validar Caracteres prohibidos (usando tu lógica de Regex)
    if (passRegex.test(pass)) {
        setReqStatus(reqChars, true);
    } else {
        setReqStatus(reqChars, false);
    }

    // 3. Validar Complejidad (Mayúscula y Número)
    if (complexityRegex.test(pass)) {
        setReqStatus(reqComplexity, true);
    } else {
        setReqStatus(reqComplexity, false);
    }
}

function setReqStatus(el, isValid) {
    const icon = el.querySelector('.material-symbols-outlined');
    if (isValid) {
        el.classList.replace('text-gray-500', 'text-green-400');
        icon.textContent = 'check_circle';
    } else {
        el.classList.replace('text-green-400', 'text-gray-500');
        icon.textContent = 'circle';
    }
}

function setupPasswordUI() {
    const section = document.getElementById('passwordInputs');
    const oldGroup = document.getElementById('oldPassGroup');
    const newGroup = document.getElementById('newPassGroup');
    const confirmGroup = document.getElementById('confirmPassGroup');
    const btnToggle = document.getElementById('btnTogglePass');
    const btnContainer = document.getElementById('changePassBtnContainer');

    // Reset de visibilidad
    [oldGroup, newGroup, confirmGroup, btnContainer].forEach(el => el.classList.add('hidden'));

    if (existUser === "0" || existUser === "1") {
        // Caso 0 (Nuevo) y 1 (Root primera vez): Forzar nueva contraseña
        newGroup.classList.remove('hidden');
        confirmGroup.classList.remove('hidden');
        document.getElementById('newPassLabel').textContent = existUser === "1" ? "Establecer Nueva Contraseña" : "Crear Contraseña";
    } 
    else if (existUser === "-1") {
        // Caso -1 (Editar): Oculto tras botón
        btnContainer.classList.remove('hidden');
        btnToggle.onclick = () => {
            const isHidden = oldGroup.classList.contains('hidden');
            [oldGroup, newGroup, confirmGroup].forEach(el => el.classList.toggle('hidden'));
            btnToggle.querySelector('span:last-child').textContent = isHidden ? "Cancelar cambio de contraseña" : "Cambiar contraseña actual";
            validateForm();
        };
    }
}

function validateForm() {
    const required = ['nameUser', 'lastNameUser', 'aliasUser', 'phoneUser', 'taxNumUser', 'pobUidUser'];
    const fieldsOK = required.every(id => document.getElementById(id).value.trim() !== '');
    const faceOK = Array.isArray(faceVector) && faceVector.length === 256;
    const AliasRegex = /^[A-Za-z\d\S]{5,20}$/;
    if (!AliasRegex.test(businessAlias)) {
        const errLabel = document.getElementById('passError');
        errLabel.textContent = "ALIAS solo debe contener letras y números con al menos 5 caracteres";
        errLabel.classList.remove('hidden');
        return;
    }

    const nPass = document.getElementById('newPass').value;
    const cPass = document.getElementById('confirmPass').value;
    const oPass = document.getElementById('oldPass').value;
    
    // Actualizar visual de requisitos
    updatePassFeedback(nPass);
    
    // La contraseña es válida si cumple tu Regex Y la complejidad
    const isPassStrong = passRegex.test(nPass) && complexityRegex.test(nPass);
    const isConfirmMatch = nPass === cPass && nPass !== '';

    let passOK = true;

    if (existUser === "0" || existUser === "1") {
        passOK = isPassStrong && isConfirmMatch;
    } else if (existUser === "-1") {
        const isChanging = !document.getElementById('newPassGroup').classList.contains('hidden');
        if (isChanging) {
            passOK = oPass.length > 0 && isPassStrong && isConfirmMatch;
        }
    }

    // Error visual si no coinciden
    const errLabel = document.getElementById('passError');
    if (nPass && cPass && nPass !== cPass) {
        errLabel.textContent = "Las contraseñas no coinciden.";
        errLabel.classList.remove('hidden');
    } else if (nPass && !isPassStrong) {
        errLabel.textContent = "La contraseña no cumple con los requisitos de seguridad.";
        errLabel.classList.remove('hidden');
    } else {
        errLabel.classList.add('hidden');
    }

    document.getElementById('btnSubmit').disabled = !(fieldsOK && faceOK && passOK);
}
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
        alert('Error cargando el modelo facial');
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

/* ============================================================
CAMARA
============================================================ */
function setupCameraUI() {
    document.getElementById('btnTakePhoto').onclick = startCamera;
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
        alert('No se pudo acceder a la cámara');
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

        // Preview visual
        document.getElementById('photoPreview').src = canvas.toDataURL('image/jpeg');
        document.getElementById('photoPreviewContainer').classList.remove('hidden');
        document.getElementById('photoStatus').textContent = 'Rostro capturado ✓';
        document.getElementById('photoStatus').style.color = '#10b981';

        stopCamera();
        validateForm();

    } catch (err) {
        console.error(err);
        alert('Error procesando el rostro');
    } finally {
        isProcessing = false;
    }
}

/* ============================================================
FORM VALIDATION
============================================================ */
function validateForm() {
    const required = [
        'nameUser',
        'lastNameUser',
        'aliasUser',
        'phoneUser',
        'taxNumUser',
        'pobUidUser'
    ];

    const fieldsOK = required.every(id =>
        document.getElementById(id).value.trim() !== ''
    );

    const faceOK = Array.isArray(faceVector) && faceVector.length === 256;

    document.getElementById('btnSubmit').disabled = !(fieldsOK && faceOK);
}

// --- Funciones Lógicas Existentes (modificadas) ---

async function fetchInviteData(id) {
    try {
        const response = await fetch(`/getinviteuser?id=${id}`);
        if (!response.ok) throw new Error('Error obteniendo invitación');
        
        const data = await response.json();
        const hostData = data.host;
        const inviteData = data[id];

        // Asignamos el estado del usuario según el backend
        // hostData.exists debe devolver "0", "1" o "-1"
        existUser = String(hostData.exists); 
        
        if (existUser !== "0") {
            document.getElementById('nameUser').value = hostData.nameUser || '';
            document.getElementById('lastNameUser').value = hostData.lastNameUser || '';
            document.getElementById('aliasUser').value = hostData.aliasUser || '';
            document.getElementById('phoneUser').value = hostData.phoneUser || '';
        }

        setupPasswordUI(); // <--- Inicializar UI de passwords

        // Llenado de información visual (labels/textos)
        inviteEmail = inviteData.emailDest;
        document.getElementById('loadingInfo').classList.add('hidden');
        document.getElementById('inviteInfoContent').classList.remove('hidden');

        const instName = hostData.legalNameInst || hostData.aliasNameInst || 'Institución Sin Nombre';
        document.getElementById('hostName').textContent = `${hostData.nameUser} ${hostData.lastNameUser}`;
        document.getElementById('hostInstitution').textContent = instName;
        document.getElementById('userEmailDisplay').textContent = inviteData.emailDest;
        
        // Formateo de fecha
        const expDate = new Date(inviteData.expirationDate);
        document.getElementById('expireDate').textContent = !isNaN(expDate) 
            ? expDate.toLocaleDateString('es-MX', { year: 'numeric', month: 'long', day: 'numeric' }) 
            : 'Sin fecha';
        
        // Manejo del Rol
        const roletype = {1: "Root", 2: "Admin", 3: "Usuario"};
        const roleInput = document.getElementById('roleApp');
        roleInput.textContent = roletype[inviteData.roleApp] || 'Sin rol';
        roleInput.value = inviteData.roleApp;

    } catch (err) {
        console.error(err);
        showErrorGlobal(err.message);
        document.getElementById('btnSubmit').disabled = true;
        const inputs = document.querySelectorAll('input');
        inputs.forEach(i => i.disabled = true);
    }
}

async function handleRegister(e) {
    e.preventDefault();
    const btn = document.getElementById('btnSubmit');
    const msgDiv = document.getElementById('formMessage');
    
    const payload = {
        idInvite: inviteId,
        name: document.getElementById('nameUser').value.trim(),
        lastname: document.getElementById('lastNameUser').value.trim(),
        alias: document.getElementById('aliasUser').value.trim(),
        email: inviteEmail,
        phone: document.getElementById('phoneUser').value.trim(),
        rfc: document.getElementById('taxNumUser').value.toUpperCase(),
        curp: document.getElementById('pobUidUser').value.toUpperCase(),
        role: document.getElementById('roleApp').value,
        alive: true,
        facevector: faceVector
    };

    // Agregar passwords según el caso
    const nPass = document.getElementById('newPass').value;
    const oPass = document.getElementById('oldPass').value;

    if (nPass) payload.newpass = nPass;
    if (existUser === "-1" && oPass) payload.oldpass = oPass;

    btn.disabled = true;
    btn.innerHTML = `<div class="animate-spin h-5 w-5 border-2 border-white border-t-transparent rounded-full"></div> Procesando...`;

    let url = (existUser === "0") ? '/createuser' : '/updateUserStatus';
    let msj = (existUser === "0") ? "¡Registro exitoso!" : "¡Datos actualizados!";

    try {
        const response = await fetch(url, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            credentials: "include",
            body: JSON.stringify(payload)
        });

        if (response.ok) {
            msgDiv.textContent = msj;
            msgDiv.className = 'rounded-lg p-4 text-sm font-medium bg-green-500/20 text-green-400 border border-green-500/30';
            setTimeout(() => window.location.href = '/login', 2000);
        } else {
            const errData = await response.text();
            throw new Error(errData || 'Error en la operación');
        }
    } catch (err) {
        msgDiv.textContent = `Error: ${err.message}`;
        msgDiv.className = 'rounded-lg p-4 text-sm font-medium bg-red-500/20 text-red-400 border border-red-500/30';
        btn.disabled = false;
        btn.innerHTML = `<span>Reintentar</span> <span class="material-symbols-outlined">refresh</span>`;
    }
}

function showErrorGlobal(msg) {
    const container = document.getElementById('loadingInfo');
    container.innerHTML = `
        <div class="text-center p-6 bg-red-500/10 border border-red-500/30 rounded-xl">
            <span class="material-symbols-outlined text-red-400 text-3xl mb-2">error</span>
            <h3 class="text-red-400 font-bold">Error</h3>
            <p class="text-gray-400 text-sm mt-1">${msg}</p>
        </div>
    `;
    document.getElementById('registerForm').style.opacity = '0.5';
    document.getElementById('registerForm').style.pointerEvents = 'none';
}

// Añadir event listeners para validación en tiempo real
document.querySelectorAll('.input-tech').forEach(input => {
    input.addEventListener('input', validateForm);
});