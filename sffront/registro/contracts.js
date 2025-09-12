// ====== Estado y utilidades ======
const LS_KEYS = {
  ONBOARDING_DONE: 'sf_onboarding_done',
  DOC_STATUSES: 'sf_doc_statuses',
  PROFILE: 'sf_profile',
  CREDS_META: 'sf_creds_meta'
};

const DOCS = [
  {
    id: 'aviso',
    title: 'Aviso de Privacidad y Términos y Condiciones',
    type: 'static'
  },
  {
    id: 'contrato_ts',
    title: 'Contrato de Prestación de Servicio de Sellos de Tiempo',
    type: 'dynamic'
  },
  {
    id: 'contrato_conservacion',
    title: 'Contrato de Constancia de Conservación de Mensajes de Texto',
    type: 'dynamic'
  }
];

function $(sel, root = document) { return root.querySelector(sel); }
function $all(sel, root = document) { return Array.from(root.querySelectorAll(sel)); }

function readFileAsText(file) {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.onload = () => resolve(String(reader.result));
    reader.onerror = reject;
    reader.readAsText(file);
  });
}

function getTodayISO() {
  const d = new Date();
  return d.toISOString().slice(0, 10);
}

function loadStatuses() {
  try { return JSON.parse(localStorage.getItem(LS_KEYS.DOC_STATUSES)) || {}; }
  catch { return {}; }
}
function saveStatuses(statuses) {
  localStorage.setItem(LS_KEYS.DOC_STATUSES, JSON.stringify(statuses));
}
function loadProfile() {
  try { return JSON.parse(localStorage.getItem(LS_KEYS.PROFILE)) || {}; }
  catch { return {}; }
}
function saveProfile(profile) {
  localStorage.setItem(LS_KEYS.PROFILE, JSON.stringify(profile));
}

function setOnboardingDone(meta) {
  localStorage.setItem(LS_KEYS.ONBOARDING_DONE, '1');
  if (meta) localStorage.setItem(LS_KEYS.CREDS_META, JSON.stringify(meta));
}
function isOnboardingDone() {
  return localStorage.getItem(LS_KEYS.ONBOARDING_DONE) === '1';
}

// ====== Plantillas de documento ======
const templates = {
  aviso: `
  <h2>Aviso de Privacidad</h2>
  <p>Este Aviso de Privacidad describe cómo SignForce trata tus datos personales conforme a la legislación aplicable. Tu información se utiliza exclusivamente para la prestación de servicios de firma y sellado de tiempo. Al continuar, aceptas los términos descritos.</p>
  <h3>Términos y Condiciones</h3>
  <p>El uso de la plataforma implica la aceptación plena de los presentes términos. El usuario se compromete a proporcionar información veraz y a custodiar sus credenciales de firma.</p>
  <p>Fecha: ${getTodayISO()}</p>
  `,

  contrato_ts: (p) => `
  <h2>Contrato de Prestación de Servicio de Sellos de Tiempo</h2>
  <p>Entre <strong>${p.nombreEmpresa || 'NOMBRE EMPRESA'}</strong>, con RFC <strong>${p.rfcEmpresa || 'RFC-EMPRESA'}</strong> y domicilio en <strong>${p.direccionEmpresa || 'DIRECCIÓN EMPRESA'}</strong>, en lo sucesivo <em>“La Empresa”</em>, y <strong>${p.nombreUsuario || 'NOMBRE USUARIO'}</strong>, con RFC <strong>${p.rfcUsuario || 'RFC-USUARIO'}</strong> y domicilio en <strong>${p.direccionUsuario || 'DIRECCIÓN USUARIO'}</strong>, en lo sucesivo <em>“El Usuario”</em>, celebran el presente contrato al tenor de las siguientes cláusulas:</p>
  <ol>
    <li><strong>Objeto.</strong> La Empresa proveerá el servicio de sellos de tiempo para los documentos electrónicos del Usuario.</li>
    <li><strong>Vigencia.</strong> A partir de la firma del presente y hasta su terminación por cualquiera de las partes.</li>
    <li><strong>Obligaciones.</strong> El Usuario deberá resguardar su llave privada y contraseña. La Empresa mantendrá la trazabilidad y evidencia de sellado.</li>
  </ol>
  <p>Lugar y fecha: Ciudad de México, ${getTodayISO()}.</p>
  `,

  contrato_conservacion: (p) => `
  <h2>Contrato de Constancia de Conservación de Mensajes de Texto</h2>
  <p>Comparecen por una parte <strong>${p.nombreEmpresa || 'NOMBRE EMPRESA'}</strong>, con RFC <strong>${p.rfcEmpresa || 'RFC-EMPRESA'}</strong>, y por la otra <strong>${p.nombreUsuario || 'NOMBRE USUARIO'}</strong>, RFC <strong>${p.rfcUsuario || 'RFC-USUARIO'}</strong>, para acordar lo siguiente:</p>
  <ol>
    <li><strong>Servicio.</strong> La Empresa conservará mensajes de texto proveídos por el Usuario asegurando integridad y disponibilidad.</li>
    <li><strong>Evidencia.</strong> Se generarán constancias electrónicas de conservación asociadas a cada mensaje o lote.</li>
    <li><strong>Confidencialidad.</strong> Ambas partes se obligan a no divulgar información reservada.</li>
  </ol>
  <p>Firmas electrónicas de conformidad en la fecha ${getTodayISO()}.</p>
  `
};

// ====== Render de documentos ======
function renderDocBody(doc, profile) {
  if (doc.type === 'static') return templates.aviso;
  if (doc.id === 'contrato_ts') return templates.contrato_ts(profile);
  if (doc.id === 'contrato_conservacion') return templates.contrato_conservacion(profile);
  return '<p>Documento no disponible</p>';
}

function statusBadgeText(s) {
  if (s?.status === 'signed') return 'Firmado';
  return 'Pendiente';
}
function statusBadgeClass(s) {
  if (s?.status === 'signed') return 'border:1px solid rgba(35,197,94,.4); color:#d1fae5; background:rgba(35,197,94,.15)';
  return 'border:1px solid rgba(234,179,8,.4); color:#fff8dc; background:rgba(234,179,8,.15)';
}

// ====== Simulación de endpoint ======
async function fetchStatusFromAPI() {
  // Simula latencia y respuesta desde /api/documentStatus
  await new Promise(r => setTimeout(r, 300));
  // En un escenario real, aquí harías fetch('/api/documentStatus', {credentials:'include'})
  // Retornamos vacío para que localStorage sea la fuente de verdad
  return {};
}

// ====== App init ======
document.addEventListener('DOMContentLoaded', () => {
  // Liquid bubble simple: posicionar sobre 'Inicio'
  const bubble = document.getElementById('liquidBubble');
  const mainHeader = document.getElementById('mainHeader');
  const navItems = document.querySelectorAll('.nav-item');
  function moveBubble(target) {
    const rect = target.getBoundingClientRect();
    const headerRect = mainHeader.getBoundingClientRect();
    bubble.style.left = `${rect.left - headerRect.left - 10}px`;
    bubble.style.width = `${rect.width + 20}px`;
  }
  window.addEventListener('load', () => { if (navItems[0]) moveBubble(navItems[0]); });
  navItems.forEach(item => item.addEventListener('mouseenter', () => moveBubble(item)));

  // Tema
  const themeToggle = document.getElementById('themeToggle');
  themeToggle.addEventListener('click', () => {
    document.body.classList.toggle('light-theme');
    themeToggle.textContent = document.body.classList.contains('light-theme') ? '☀️' : '🌙';
  });

  // Ruteo inicial
  if (isOnboardingDone()) {
    showDashboard();
  } else {
    showSetup();
  }

  // ====== Form credenciales ======
const credentialsForm = document.getElementById('credentialsForm');
credentialsForm.addEventListener('submit', async (e) => {
  e.preventDefault();
  console.log("Enviando a api");
  const keyFile = document.getElementById('keyFile').files[0];
  const certFile = document.getElementById('certFile').files[0];
  const keyPassword = document.getElementById('keyPassword').value;

  if (!keyFile || !certFile || !keyPassword) {
    alert('Por favor, sube la llave, el certificado y escribe la contraseña.');
    return;
  }

  try {
    // Crear FormData para enviar los archivos
    const formData = new FormData();
    formData.append('keyFile', keyFile);
    formData.append('certFile', certFile);
    formData.append('passKey', keyPassword);

    // Enviar al endpoint de Go
    const response = await fetch('/uploadKeys', {
      method: 'POST',
      body: formData,
      credentials: 'include' // Para incluir las cookies (JWT)
    });

    if (!response.ok) {
      const errorText = await response.text();
      throw new Error(`Error del servidor: ${response.status} - ${errorText}`);
    }

    const result = await response.json();
    
    if (result.success) {
      // Guardar metadata localmente
      const meta = {
        keyName: keyFile.name,
        certName: certFile.name,
        uploadedAt: new Date().toISOString()
      };
      setOnboardingDone(meta);
      //showDashboard();
      alert('Credenciales validadas y guardadas exitosamente.');
    } else {
      alert('Error: ' + result.message);
    }
  } catch (err) {
    console.error('Error subiendo archivos:', err);
    alert('Error al subir los archivos: ' + err.message);
  }
});

  // Perfil (usuario/empresa)
  const profileForm = document.getElementById('profileForm');
  profileForm.addEventListener('submit', (e) => {
    e.preventDefault();
    const profile = {
      nombreUsuario: document.getElementById('nombreUsuario').value.trim(),
      rfcUsuario: document.getElementById('rfcUsuario').value.trim(),
      direccionUsuario: document.getElementById('direccionUsuario').value.trim(),
      nombreEmpresa: document.getElementById('nombreEmpresa').value.trim(),
      rfcEmpresa: document.getElementById('rfcEmpresa').value.trim(),
      direccionEmpresa: document.getElementById('direccionEmpresa').value.trim()
    };
    saveProfile(profile);
    renderDocs(); // refrescar vistas con datos nuevos
    alert('Datos guardados.');
  });

  // Modal preview
  const previewModal = document.getElementById('previewModal');
  document.getElementById('closePreview').addEventListener('click', () => previewModal.close());
  document.getElementById('closePreview2').addEventListener('click', () => previewModal.close());

  // Acciones header
  document.getElementById('goHome').addEventListener('click', () => window.scrollTo({ top: 0, behavior: 'smooth' }));

  // Carga inicial de perfil
  const existingProfile = loadProfile();
  if (Object.keys(existingProfile).length) {
    document.getElementById('nombreUsuario').value = existingProfile.nombreUsuario || '';
    document.getElementById('rfcUsuario').value = existingProfile.rfcUsuario || '';
    document.getElementById('direccionUsuario').value = existingProfile.direccionUsuario || '';
    document.getElementById('nombreEmpresa').value = existingProfile.nombreEmpresa || '';
    document.getElementById('rfcEmpresa').value = existingProfile.rfcEmpresa || '';
    document.getElementById('direccionEmpresa').value = existingProfile.direccionEmpresa || '';
  }
});

// ====== Vistas ======
function showSetup() {
  document.getElementById('view-setup').classList.remove('hidden');
  document.getElementById('view-dashboard').classList.add('hidden');
}
function showDashboard() {
  document.getElementById('view-setup').classList.add('hidden');
  document.getElementById('view-dashboard').classList.remove('hidden');
  initDashboard();
}

async function initDashboard() {
  // Simular obtención de estatus del endpoint y mezclar con localStorage
  const apiStatuses = await fetchStatusFromAPI();
  const localStatuses = loadStatuses();
  const merged = { ...apiStatuses, ...localStatuses };
  saveStatuses(merged);
  renderDocs();
}

function renderDocs() {
  const list = document.getElementById('docsList');
  list.innerHTML = '';
  const profile = loadProfile();
  const statuses = loadStatuses();

  DOCS.forEach(doc => {
    const tpl = document.getElementById('doc-card-template');
    const node = tpl.content.cloneNode(true);
    const card = node.querySelector('.doc-card');

    card.querySelector('.doc-title').textContent = doc.title;
    const badge = card.querySelector('[data-status]');
    const meta = card.querySelector('[data-meta]');
    const body = card.querySelector('[data-doc-body]');

    const currentStatus = statuses[doc.id];
    badge.textContent = statusBadgeText(currentStatus);
    badge.setAttribute('style', statusBadgeClass(currentStatus));

    // Render cuerpo
    body.innerHTML = renderDocBody(doc, profile);

    // Botones
    const btnPreview = card.querySelector('[data-preview]');
    const btnSign = card.querySelector('[data-sign]');

    btnPreview.addEventListener('click', () => openPreview(doc.title, body.innerHTML));

    btnSign.addEventListener('click', async () => {
      // Validación mínima: que exista meta de credenciales
      const metaCreds = localStorage.getItem(LS_KEYS.CREDS_META);
      if (!metaCreds) {
        alert('Primero sube tu llave, certificado y contraseña.');
        return;
      }
      // Simular tiempo de firmado
      btnSign.disabled = true;
      btnSign.textContent = 'Firmando…';
      await new Promise(r => setTimeout(r, 700));

      const newStatuses = loadStatuses();
      newStatuses[doc.id] = { status: 'signed', timestamp: new Date().toISOString() };
      saveStatuses(newStatuses);

      // Actualizar UI
      badge.textContent = 'Firmado';
      badge.setAttribute('style', statusBadgeClass({ status: 'signed' }));
      meta.textContent = `Firmado el ${new Date().toLocaleString()}`;

      btnSign.textContent = 'Firmado';
      btnSign.disabled = true;
    });

    // Mostrar meta si ya estaba firmado
    if (currentStatus?.status === 'signed' && currentStatus.timestamp) {
      meta.textContent = `Firmado el ${new Date(currentStatus.timestamp).toLocaleString()}`;
      btnSign.textContent = 'Firmado';
      btnSign.disabled = true;
    }

    list.appendChild(node);
  });
}

function openPreview(title, html) {
  const modal = document.getElementById('previewModal');
  document.getElementById('previewTitle').textContent = title;
  document.getElementById('previewBody').innerHTML = html;
  modal.showModal();
}
