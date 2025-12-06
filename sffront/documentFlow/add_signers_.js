document.addEventListener('DOMContentLoaded', () => {
  // ==============================
  // 🔹 CARGA DE DOCUMENTOS DESDE SESSION STORAGE
  // ==============================
  let folderData = JSON.parse(sessionStorage.getItem("folder")) || {};
  let docIds = Object.keys(folderData);
  let currentIndex = parseInt(sessionStorage.getItem("currentIndex") || "0");
  let currentDocId = docIds[currentIndex];
  let currentDoc = folderData[currentDocId];

  if (!currentDoc) {
    alert("No hay documentos para procesar.");
    window.location.href = "/documents";
    return;
  }

  // Mostrar información básica del documento actual
  document.getElementById("docName").value = `${currentDoc.documentName}.${currentDoc.documentExt}`;

  // ==============================
  // 🔹 VARIABLES Y FUNCIONES GLOBALES
  // ==============================
  const reviewers = [];
  const teamSelect = document.getElementById("teamSelect");
  const userSelect = document.getElementById("userSelect");

  async function toBase64(file) {
    return new Promise((resolve, reject) => {
      const reader = new FileReader();
      reader.readAsDataURL(file);
      reader.onload = () => resolve(reader.result);
      reader.onerror = error => reject(error);
    });
  }

  // ==============================
  // 🔹 CARGAR EQUIPOS DE LA INSTITUCIÓN
  // ==============================
  async function loadTeams() {
  try {
    console.log("Cargando equipos...");
    const response = await fetch("/instteams", {
      method: "GET",
      credentials: "include",
    });

    if (!response.ok) throw new Error(`HTTP ${response.status}`);
    const teams = await response.json();

    const teamSelect = document.getElementById("teamSelect");
    if (!teamSelect) {
      console.warn("⚠️ No se encontró el select de equipos.");
      return;
    }

    if (!Array.isArray(teams) || teams.length === 0) {
      teamSelect.innerHTML = `<option value="">No hay equipos disponibles</option>`;
      return;
    }

    teamSelect.innerHTML = `<option value="">Selecciona un equipo...</option>`;
    teams.forEach(team => {
      if (team.deletedAt && team.deletedAt.trim() !== "") return;
      const opt = document.createElement("option");
      opt.value = team.id;
      opt.textContent = `${team.name} (Límite: ${team.limitusers || "-"} usuarios / ${team.limitsigners || "-"} firmantes)`;
      teamSelect.appendChild(opt);
    });
    console.log("Equipos cargados correctamente:", teams);
  } catch (err) {
    console.error("Error al cargar equipos:", err);
  }
}

  // ==============================
  // 🔹 CARGAR USUARIOS DE UN EQUIPO
  // ==============================
  async function loadUsers(idTeam) {
  try {
    const response = await fetch("/teamusers", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      credentials: "include",
      body: JSON.stringify({
        idteam: idTeam,
        fields: []
      })
    });

    if (!response.ok) throw new Error(`HTTP ${response.status}`);
    const users = await response.json();
    console.log("Usuarios recibidos:", users);

    userSelect.innerHTML = `<option value="">Selecciona un usuario...</option>`;

    users.forEach(u => {
      // Mostrar solo usuarios activos
      if (u.active !== "1") return;

      const opt = document.createElement("option");
      opt.value = u.id;
      opt.textContent = `${(u.name && u.lastname) ? (u.name + " " + u.lastname) : (u.alias || "Sin nombre")} - ${u.email || "sin correo"}`;
      userSelect.appendChild(opt);
    });

    // Habilitar el select una vez cargado
    userSelect.disabled = false;

  } catch (err) {
    console.error("Error al cargar usuarios:", err);
    userSelect.innerHTML = `<option value="">Error al cargar usuarios</option>`;
  }
}
  // ==============================
  // 🔹 EVENTO CAMBIO DE EQUIPO
  // ==============================
  teamSelect.addEventListener("change", async () => {
    const idTeam = teamSelect.value;
    if (!idTeam) {
      userSelect.innerHTML = `<option value="">Selecciona un usuario...</option>`;
      return;
    }
    await loadUsers(idTeam);
  });

  // Inicializar carga de equipos al entrar
  loadTeams();

  // ==============================
  // 🔹 RENDERIZAR TABLA
  // ==============================
  function renderTable() {
    const tbody = document.getElementById('reviewersTableBody');
    tbody.innerHTML = '';

    reviewers.forEach((rev, idx) => {
      const tr = document.createElement('tr');
      tr.classList.add("border-t", "border-t-white/10");

      tr.innerHTML = `
        <td class="px-4 py-2 w-[400px] text-primary text-sm font-normal">${rev.user}</td>
        <td class="px-4 py-2 w-60 text-sm font-normal">
          <button class="toggle-role bg-white/20 text-primary w-full rounded-lg h-8 hover:bg-white/30" data-index="${idx}">
            ${rev.role === 1? "Signer":"Viewer"}
          </button>
        </td>
        <td class="px-4 py-2 w-[400px] text-secondary text-sm font-normal">
          <span class="editable-date" data-index="${idx}">${rev.due_date}</span>
        </td>
        <td class="px-4 py-2 w-[400px] text-secondary text-sm font-normal">${rev.team}</td>
        <td class="px-4 py-2 w-[400px] text-secondary text-sm font-normal">
          <span class="editable-comment" data-index="${idx}">${rev.comment}</span>
        </td>
      `;
      tbody.appendChild(tr);
    });

    // cambiar rol firmante/visor
    document.querySelectorAll('.toggle-role').forEach(btn => {
      btn.addEventListener('click', () => {
        const idx = btn.dataset.index;
        reviewers[idx].role = reviewers[idx].role === 1 ? 0 : 1;
        renderTable();
      });
    });

    // editar fecha
    document.querySelectorAll('.editable-date').forEach(span => {
      span.addEventListener('click', () => {
        const idx = span.dataset.index;
        const input = document.createElement('input');
        input.type = 'date';
        input.value = reviewers[idx].due_date;
        input.className = "bg-transparent text-sm text-secondary";

        input.addEventListener('blur', () => {
          reviewers[idx].due_date = input.value;
          renderTable();
        });
        input.addEventListener('keydown', (e) => {
          if (e.key === 'Enter') input.blur();
        });
        span.replaceWith(input);
        input.focus();
      });
    });

    // editar comentarios
    document.querySelectorAll('.editable-comment').forEach(span => {
      span.addEventListener('click', () => {
        const idx = span.dataset.index;
        const textarea = document.createElement('textarea');
        textarea.value = reviewers[idx].comment || '';
        textarea.className = "bg-transparent text-sm text-secondary w-full h-20 p-2 resize-none";
        textarea.addEventListener('blur', () => {
          reviewers[idx].comment = textarea.value;
          renderTable();
        });
        textarea.addEventListener('keydown', (e) => {
          if (e.key === 'Enter' && e.ctrlKey) textarea.blur();
        });
        span.replaceWith(textarea);
        textarea.focus();
      });
    });
  }

  // ==============================
  // 🔹 EVENTO AÑADIR REVISOR
  // ==============================
  document.getElementById('addReviewerBtn').addEventListener('click', () => {
    const user = document.getElementById('userSelect');
    const team = document.getElementById('teamSelect');
    const comment = document.getElementById('userComment');
    const today = new Date().toISOString().split('T')[0];

    if (!user.value || !team.value) return alert("Selecciona usuario y equipo.");

    reviewers.push({
      user: user.value,
      role: 1,
      due_date: today,
      team: team.value,
      comment: comment.value,
      positions: [] // Inicializar array vacío para posiciones
    });

    renderTable();
    user.value = '';
    comment.value = '';
  });

  // ==============================
  // 🔹 GUARDAR Y PASAR AL SIGUIENTE DOCUMENTO
  // ==============================
  document.getElementById('uploadBtn').addEventListener('click', async () => {
    // Guardar revisores en el documento actual
    currentDoc.reviewers = reviewers;
    folderData[currentDocId] = currentDoc;
    sessionStorage.setItem("folder", JSON.stringify(folderData));

    if (currentIndex + 1 < docIds.length) {
      // Pasar al siguiente documento
      sessionStorage.setItem("currentIndex", currentIndex + 1);
      window.location.href = "/addSigners";
    } else {
      // Todos los documentos tienen revisores, pasar a posicionamiento de firmas
      sessionStorage.removeItem("currentIndex");
      
      // Preparar datos para enviar al backend
      const inviteRequest = {};
      docIds.forEach(docId => {
        const doc = folderData[docId];
        inviteRequest[docId] = {
          idfolder: doc.idFolder, // Asegúrate de que este campo existe
          document: {
            idDocument: docId,
            activeDoc: doc.activeDoc,
            authRoleStatus: doc.authRoleStatus,
            authUseStatus: doc.authUseStatus,
            createdAtDoc: doc.createdAtDoc,
            documentExt: doc.documentExt,
            documentHash: doc.documentHash,
            documentName: doc.documentName,
            documentPath: doc.documentPath,
            documentFullName: `${doc.documentName}.${doc.documentExt}`,
            abstract: doc.abstractDoc,
            lastModifiedDoc: doc.lastModifiedDoc
          },
          reviewers: doc.reviewers || []
        };
      });

      // Guardar en sessionStorage para usar en add_signs.js
      sessionStorage.setItem("inviteRequest", JSON.stringify(inviteRequest));
      window.location.href = "/addSignatures";
    }
  });

  // ==============================
  // 🔹 PREVISUALIZACIÓN DEL DOCUMENTO PDF
  // ==============================
  const previewContainer = document.getElementById('preview-container');
  const dropMessage = document.getElementById('drop-message');
  if (dropMessage) dropMessage.style.display = 'none';

  const wrapper = document.createElement('div');
  wrapper.className = 'pdf-wrapper';
  wrapper.style.position = 'relative';
  wrapper.style.width = '100%';
  wrapper.style.height = '100%';

  const iframe = document.createElement('iframe');
  iframe.style.width = '100%';
  iframe.style.height = '100%';
  iframe.style.border = 'none';
  iframe.style.display = 'none';

  const loading = document.createElement('div');
  loading.className = "flex justify-center items-center h-full text-secondary";
  loading.innerHTML = `
      <div class="flex flex-col items-center gap-2">
          <span class="material-symbols-outlined text-4xl animate-spin">progress_activity</span>
          <p>Cargando documento...</p>
      </div>
  `;

  const overlay = document.createElement('div');
  overlay.innerText = 'Ver completo';
  overlay.style.position = 'absolute';
  overlay.style.bottom = '10px';
  overlay.style.right = '10px';
  overlay.style.background = 'rgba(0,0,0,0.6)';
  overlay.style.color = '#fff';
  overlay.style.padding = '6px 10px';
  overlay.style.borderRadius = '8px';
  overlay.style.cursor = 'pointer';
  overlay.style.fontSize = '12px';
  overlay.style.zIndex = '10';
  overlay.style.display = 'none'; // se muestra cuando el documento cargue

  wrapper.appendChild(loading);
  wrapper.appendChild(iframe);
  wrapper.appendChild(overlay);
  previewContainer.appendChild(wrapper);

  // Descargar el documento desde el backend
  (async () => {
    try {
      const response = await fetch('/downloadDoc', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          id: currentDocId,
          type: "uploaded",
          path: ""
        })
      });

      if (!response.ok) throw new Error(`HTTP ${response.status}`);

      const blob = await response.blob();
      const url = URL.createObjectURL(blob);
      iframe.src = url;
      iframe.style.display = 'block';
      loading.style.display = 'none';
      overlay.style.display = 'block';

      // abrir en pestaña nueva al pulsar "Ver completo"
      overlay.onclick = () => window.open(url, '_blank');
    } catch (err) {
      console.error('Error al cargar documento:', err);
      loading.innerHTML = `
        <div class="text-red-400 text-center">
          <span class="material-symbols-outlined text-4xl mb-2">error</span>
          <p>Error al cargar el documento.</p>
        </div>
      `;
    }
  })();
});