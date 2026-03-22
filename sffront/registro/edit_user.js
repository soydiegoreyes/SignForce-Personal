/* ============================================================
NOTIFICACIONES INLINE (reemplaza alert())
============================================================ */
function showError(msg, autohide = 6000) {
    const banner = document.getElementById('sf-error-banner');
    const text = document.getElementById('sf-error-text');
    if (!banner || !text) { console.error(msg); return; }
    text.textContent = msg;
    banner.style.display = 'block';
    document.getElementById('sf-info-banner').style.display = 'none';
    if (autohide) setTimeout(() => { banner.style.display = 'none'; }, autohide);
}

function showInfo(msg, autohide = 5000) {
    const banner = document.getElementById('sf-info-banner');
    const text = document.getElementById('sf-info-text');
    if (!banner || !text) { console.log(msg); return; }
    text.textContent = msg;
    banner.style.display = 'block';
    document.getElementById('sf-error-banner').style.display = 'none';
    if (autohide) setTimeout(() => { banner.style.display = 'none'; }, autohide);
}
/* ============================================================
HELPERS: CAPITALIZE + ALIAS SUGGESTION + PHONE PREFIX
============================================================ */
function capFirst(input) {
    const v = input.value;
    if (v.length > 0) {
        input.value = v.charAt(0).toUpperCase() + v.slice(1);
    }
}

function suggestAlias() {
    const name = document.getElementById('nameUser')?.value.trim() || '';
    const last = document.getElementById('lastNameUser')?.value.trim() || '';
    if (name.length < 2 || last.length < 1) return;
    const prefix = name.substring(0, 2);
    const surfix = last.split(' ')[0];
    const suggestion = (prefix + surfix).toLowerCase().replace(/[^a-z0-9]/g, '');
    const aliasField = document.getElementById('aliasUser');
    if (aliasField && (aliasField.value === '' || aliasField.dataset.autosuggested === 'true')) {
        aliasField.value = suggestion;
        aliasField.dataset.autosuggested = 'true';
    }
}



/* ============================================================
VARIABLES GLOBALES
============================================================ */
let existUser = null;
let inviteId = null;
let inviteEmail = null;
// faceModel eliminado - ahora usa kycFaceLandmarker de MediaPipe
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
        showError('Invitación inválida o expirada');
        return;
    }


    fetchInviteData(inviteId);
    setupKycUI();
    // Alias: cuando el usuario escribe, desactivar autosuggested
    const aliasField = document.getElementById('aliasUser');
    if (aliasField) {
        aliasField.addEventListener('input', () => {
            aliasField.dataset.autosuggested = 'false';
        });
    }

    ['oldPass', 'newPass', 'confirmPass'].forEach(id => {
        document.getElementById(id).addEventListener('input', validateForm);
    });
});


/* ============================================================
LÓGICA DE CONTRASEÑAS Y UI
============================================================ */
function updatePassFeedback(pass) {
    const bar   = document.getElementById('strengthBar');
    const fill  = document.getElementById('strengthFill');
    const label = document.getElementById('strengthLabel');
    const reqs  = document.getElementById('passReqs');
    if (!bar) return;

    // Mostrar controles si hay texto
    if (pass.length > 0) {
        bar.classList.remove('hidden');
        label.classList.remove('hidden');
        reqs.classList.remove('hidden');
    } else {
        bar.classList.add('hidden');
        label.classList.add('hidden');
        reqs.classList.add('hidden');
        return;
    }

    // Criterios
    const checks = {
        len:   pass.length >= 8,
        upper: /[A-Z]/.test(pass),
        lower: /[a-z]/.test(pass),
        num:   /[0-9]/.test(pass),
        spec:  /[^A-Za-z0-9]/.test(pass),
    };
    const score = Object.values(checks).filter(Boolean).length; // 0-5

    // Actualizar dots de requisitos
    Object.entries(checks).forEach(([key, ok]) => {
        const dot = document.getElementById('dot-' + key);
        const li  = document.getElementById('req-' + key);
        if (!dot || !li) return;
        dot.style.background = ok ? '#4ade80' : '#4b5563';
        li.style.color = ok ? '#4ade80' : '#9ca3af';
    });

    // Barra de fuerza
    const levels = [
        { pct: '20%', color: '#ef4444', text: 'Muy débil',  textColor: '#ef4444' },
        { pct: '40%', color: '#f97316', text: 'Débil',      textColor: '#f97316' },
        { pct: '60%', color: '#eab308', text: 'Regular',    textColor: '#eab308' },
        { pct: '80%', color: '#84cc16', text: 'Fuerte',     textColor: '#84cc16' },
        { pct: '100%',color: '#22c55e', text: '¡Excelente!',textColor: '#22c55e' },
    ];
    const lvl = levels[Math.max(0, score - 1)];
    fill.style.width   = lvl.pct;
    fill.style.background = lvl.color;
    label.textContent  = lvl.text;
    label.style.color  = lvl.textColor;
}

function setReqStatus(el, isValid) {
    // Compatibilidad - no se usa con el nuevo sistema
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
    const nPass = document.getElementById('newPass').value;
    const cPass = document.getElementById('confirmPass').value;
    const oPass = document.getElementById('oldPass').value;

    // Actualizar visual de requisitos SIEMPRE (no puede quedar bloqueado por alias)
    updatePassFeedback(nPass);
    
    // La contraseña es válida si cumple tu Regex Y la complejidad
    const isPassStrong = nPass.length >= 8
            && /[A-Z]/.test(nPass)
            && /[a-z]/.test(nPass)
            && /[0-9]/.test(nPass)
            && /[^A-Za-z0-9]/.test(nPass);
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
KYC MODERNO — face-api.js (vladmandic) + Liveness Detection
Steps: Centro → Izquierda → Derecha → Captura
============================================================ */

const FACEAPI_MODEL_URL = 'https://cdn.jsdelivr.net/npm/@vladmandic/face-api/model';
let kycModelsLoaded = false;
let kycStream = null;
let kycAnimFrame = null;
let kycStep = 0; // 0=centro, 1=izq, 2=der
let kycStepFrames = 0;
const FRAMES_TO_CONFIRM = 10;

const KYC_STEPS = [
  { text: 'Centra tu rostro en el óvalo', sub: 'Mira directamente a la cámara', dot: 0 },
  { text: 'Gira la cabeza a la IZQUIERDA  ←', sub: 'Despacio, mantén hasta que desaparezca la flecha', dot: 1 },
  { text: 'Gira la cabeza a la DERECHA  →', sub: 'Despacio, mantén hasta que desaparezca la flecha', dot: 2 },
];

async function initKycModel() {
  if (kycModelsLoaded) return true;
  try {
    await faceapi.nets.tinyFaceDetector.loadFromUri(FACEAPI_MODEL_URL);
    await faceapi.nets.faceLandmark68TinyNet.loadFromUri(FACEAPI_MODEL_URL);
    kycModelsLoaded = true;
    return true;
  } catch (e) {
    console.error('face-api error:', e);
    showError('No se pudo cargar el sistema de verificación. Verifica tu conexión.');
    return false;
  }
}

function estimateYaw(landmarks) {
  // Usa landmarks 68-point para estimar yaw (rotación horizontal)
  const pts = landmarks.positions;
  const leftEyeOuter  = pts[36]; // esquina externa ojo izq
  const rightEyeOuter = pts[45]; // esquina externa ojo der
  const noseTip       = pts[30]; // punta de la nariz
  const eyeMidX = (leftEyeOuter.x + rightEyeOuter.x) / 2;
  const eyeWidth = Math.abs(rightEyeOuter.x - leftEyeOuter.x);
  if (eyeWidth < 1) return 0;
  // offset normalizado: positivo = nariz hacia derecha = cabeza girada derecha
  return (noseTip.x - eyeMidX) / eyeWidth;
}

function updateKycStep(step) {
  kycStep = step;
  kycStepFrames = 0;
  const s = KYC_STEPS[step];
  if (!s) return;
  document.getElementById('kycInstructionText').textContent = s.text;
  document.getElementById('kycInstructionSub').textContent = s.sub;
  [0,1,2].forEach(i => {
    const dot = document.getElementById('dot'+i);
    if (i < step)       { dot.style.background='#4ade80'; dot.style.width='10px'; dot.style.height='10px'; }
    else if (i===step)  { dot.style.background='#B5C413'; dot.style.width='12px'; dot.style.height='12px'; }
    else                { dot.style.background='rgba(255,255,255,0.18)'; dot.style.width='8px'; dot.style.height='8px'; }
  });
  const arrow = document.getElementById('kycArrow');
  const arrowInner = document.getElementById('kycArrowInner');
  if (step===1) {
    arrow.style.opacity='1';
    arrowInner.textContent='←';
    arrowInner.style.animation='arrow-pulse 0.7s ease-in-out infinite alternate';
  } else if (step===2) {
    arrow.style.opacity='1';
    arrowInner.textContent='→';
    arrowInner.style.animation='arrow-pulse-r 0.7s ease-in-out infinite alternate';
  } else {
    arrow.style.opacity='0';
  }
}

async function kycDetectLoop() {
  const video = document.getElementById('kycVideo');
  if (!video || video.readyState < 2 || !kycStream) return;

  const detection = await faceapi
    .detectSingleFace(video, new faceapi.TinyFaceDetectorOptions({ inputSize: 224, scoreThreshold: 0.4 }))
    .withFaceLandmarks(true);

  const oval = document.getElementById('kycOvalRing');
  if (!detection) {
    oval.style.borderColor='rgba(239,68,68,0.6)';
    kycStepFrames = 0;
    kycAnimFrame = requestAnimationFrame(kycDetectLoop);
    return;
  }

  oval.style.borderColor='rgba(181,196,19,0.85)';
  const yaw = estimateYaw(detection.landmarks); // -0.5 a +0.5

  if (kycStep === 0) {
    // Centro: yaw entre -0.12 y +0.12
    if (Math.abs(yaw) < 0.12) {
      kycStepFrames++;
      if (kycStepFrames >= FRAMES_TO_CONFIRM) updateKycStep(1);
    } else kycStepFrames = 0;

  } else if (kycStep === 1) {
    // Izquierda del usuario = derecha en el espejo (yaw positivo)
    if (yaw > 0.22) {
      kycStepFrames++;
      if (kycStepFrames >= FRAMES_TO_CONFIRM) updateKycStep(2);
    } else kycStepFrames = 0;

  } else if (kycStep === 2) {
    // Derecha del usuario = izquierda en el espejo (yaw negativo)
    if (yaw < -0.22) {
      kycStepFrames++;
      if (kycStepFrames >= FRAMES_TO_CONFIRM) {
        captureKycPhoto(video, detection.landmarks);
        return;
      }
    } else kycStepFrames = 0;
  }

  if (kycStream) kycAnimFrame = requestAnimationFrame(kycDetectLoop);
}

function captureKycPhoto(video, landmarks) {
  cancelAnimationFrame(kycAnimFrame);
  kycAnimFrame = null;

  document.getElementById('kycOvalRing').style.borderColor='#4ade80';
  document.getElementById('kycArrow').style.opacity='0';
  document.getElementById('kycCheckOverlay').style.display='flex';
  document.getElementById('kycInstructionText').textContent='¡Verificación completada!';
  document.getElementById('kycInstructionSub').textContent='Procesando tu identidad...';
  [0,1,2].forEach(i => {
    const dot = document.getElementById('dot'+i);
    dot.style.background='#4ade80'; dot.style.width='10px'; dot.style.height='10px';
  });

  const canvas = document.getElementById('kycCanvas');
  canvas.width = video.videoWidth;
  canvas.height = video.videoHeight;
  const ctx = canvas.getContext('2d');
  ctx.translate(canvas.width, 0); ctx.scale(-1, 1); // desespejear
  ctx.drawImage(video, 0, 0);
  const photoDataUrl = canvas.toDataURL('image/jpeg', 0.85);

  // Vector facial: 68 landmarks × 2 coords + normalización = 136 → pad a 256
  const pts = landmarks.positions;
  const xs = pts.map(p=>p.x), ys = pts.map(p=>p.y);
  const minX=Math.min(...xs),maxX=Math.max(...xs),minY=Math.min(...ys),maxY=Math.max(...ys);
  const w=maxX-minX||1, h=maxY-minY||1;
  const raw = [];
  pts.forEach(p => { raw.push((p.x-minX)/w); raw.push((p.y-minY)/h); });
  // Pad a 256 repitiendo valores
  while (raw.length < 256) raw.push(raw[raw.length % raw.length]);
  const vector = raw.slice(0, 256).map(v => parseFloat(v.toFixed(6)));

  setTimeout(() => {
    stopKycCamera();
    document.getElementById('kycPreviewImg').src = photoDataUrl;
    document.getElementById('kycStartSection').classList.add('hidden');
    document.getElementById('kycPreviewSection').style.display='flex';
    document.getElementById('kycPreviewSection').classList.remove('hidden');
    document.getElementById('kycBadge').classList.remove('hidden');
    document.getElementById('kycBadge').style.display='flex';
    document.getElementById('photoStatus').textContent='Verificación completada ✓';
    document.getElementById('photoStatus').style.color='#4ade80';
    faceVector = vector;
    checkFormReady();
  }, 1000);
}

function stopKycCamera() {
  if (kycAnimFrame) { cancelAnimationFrame(kycAnimFrame); kycAnimFrame = null; }
  if (kycStream) { kycStream.getTracks().forEach(t=>t.stop()); kycStream=null; }
  document.getElementById('kycModal').style.display='none';
}

async function startKyc() {
  kycStep=0; kycStepFrames=0;
  const modal = document.getElementById('kycModal');
  modal.style.display='flex';
  document.getElementById('kycCheckOverlay').style.display='none';
  document.getElementById('kycOvalRing').style.borderColor='rgba(181,196,19,0.5)';
  document.getElementById('kycArrow').style.opacity='0';
  document.getElementById('kycInstructionText').textContent='Cargando modelo...';
  document.getElementById('kycInstructionSub').textContent='Un momento por favor';
  [0,1,2].forEach(i => {
    const dot=document.getElementById('dot'+i);
    dot.style.background='rgba(255,255,255,0.18)'; dot.style.width='8px'; dot.style.height='8px';
  });

  const loaded = await initKycModel();
  if (!loaded) { modal.style.display='none'; return; }

  try {
    kycStream = await navigator.mediaDevices.getUserMedia({
      video:{ facingMode:'user', width:{ideal:640}, height:{ideal:480} }, audio:false
    });
    const video = document.getElementById('kycVideo');
    video.srcObject = kycStream;
    await new Promise(r => { video.onloadedmetadata=r; });
    video.play();
    updateKycStep(0);
    kycAnimFrame = requestAnimationFrame(kycDetectLoop);
  } catch(e) {
    stopKycCamera();
    showError('No se pudo acceder a la cámara. Verifica los permisos del navegador.');
  }
}

function setupKycUI() {
  document.getElementById('btnStartKyc')?.addEventListener('click', startKyc);
  document.getElementById('btnRetakeKyc')?.addEventListener('click', () => {
    document.getElementById('kycPreviewSection').style.display='none';
    document.getElementById('kycStartSection').classList.remove('hidden');
    document.getElementById('kycBadge').classList.add('hidden');
    document.getElementById('kycBadge').style.display='none';
    document.getElementById('photoStatus').textContent='Verificación pendiente';
    document.getElementById('photoStatus').style.color='';
    faceVector=null;
    checkFormReady();
  });
  document.getElementById('btnCloseKyc')?.addEventListener('click', stopKycCamera);
}



/* ============================================================
CAMARA
============================================================ */
// setupCameraUI reemplazado
function setupCameraUI_OLD() {
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
        showError('No se pudo acceder a la cámara. Verifica los permisos.');
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
        showError('Error procesando el rostro. Intenta de nuevo.');
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
            // Alias: si es root_ o vacío, generar sugerencia inteligente
            const rawAlias = hostData.aliasUser || '';
            const isAutoAlias = rawAlias === '' || rawAlias.startsWith('root_');
            if (isAutoAlias) {
                const n = (hostData.nameUser || '').substring(0, 2);
                const s = (hostData.lastNameUser || '').split(' ')[0];
                const suggested = (n + s).toLowerCase().replace(/[^a-z0-9]/g, '');
                const aliasField = document.getElementById('aliasUser');
                aliasField.value = suggested;
                aliasField.dataset.autosuggested = 'true';
            } else {
                document.getElementById('aliasUser').value = rawAlias;
            }
            // Teléfono: prefijo México si está vacío
            const phone = hostData.phoneUser || '';
            document.getElementById('phoneUser').value = phone || '+52 ';
        } else {
            // Usuario nuevo: poner prefijo México
            document.getElementById('phoneUser').value = '+52 ';
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