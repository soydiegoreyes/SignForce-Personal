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
  // 🔹 VARIABLES GLOBALES
  // ==============================
  // Cargamos revisores existentes si los hay (para persistencia al navegar)
  let reviewers = currentDoc.reviewers || [];
  
  // Elementos DOM para búsqueda
  const searchInput = document.getElementById('searchInput');
  const searchTypeSelect = document.getElementById('searchType');
  const checkExact = document.getElementById('checkExact');
  const checkGuest = document.getElementById('checkGuest');
  const checkAlive = document.getElementById('checkAliveProof');
  const suggestionsList = document.getElementById('suggestionsList');
  const selectedUserDisplay = document.getElementById('selectedUserDisplay');
  const selectedUserName = document.getElementById('selectedUserName');
  const clearSelectionBtn = document.getElementById('clearSelectionBtn');
  const searchLabel = document.getElementById('searchLabel');

  // Estado de búsqueda
  let cachedUsers = []; 
  let selectedUserObj = null; 
  let debounceTimer = null; 

  // ==============================
  // 🔹 LOGICA DE BÚSQUEDA DE USUARIOS
  // ==============================

  checkGuest.addEventListener('change', () => {
    if (checkGuest.checked) {
      searchTypeSelect.disabled = true;
      checkExact.disabled = true;
      searchLabel.innerText = "Ingresa el email del invitado";
      searchInput.placeholder = "ejemplo@correo.com";
      searchInput.value = "";
      clearSelection();
      suggestionsList.classList.add('hidden');
    } else {
      searchTypeSelect.disabled = false;
      checkExact.disabled = false;
      searchLabel.innerText = "Escribe para buscar usuario";
      searchInput.placeholder = "Escribe al menos 3 caracteres...";
      searchInput.value = "";
      clearSelection();
    }
  });

  async function executeSearch(term) {
    if (term.length === 0) {
      cachedUsers = [];
      renderSuggestions([]);
      suggestionsList.classList.add('hidden');
      return;
    }

    if (!checkExact.checked) {
        if (term.length >= 3) {
            if (cachedUsers.length > 0 && term.length > 3) {
                filterLocalSuggestions(term);
            } else {
                await fetchUsers(term);
            }
        } else {
            suggestionsList.classList.add('hidden');
        }
    } else {
        if (term.length >= 3) {
            await fetchUsers(term);
        }
    }
  }

  searchInput.addEventListener('input', (e) => {
    const term = e.target.value.trim();
    if (checkGuest.checked) return;

    if (term.length === 0) {
      clearTimeout(debounceTimer);
      executeSearch("");
      return;
    }

    clearTimeout(debounceTimer);
    debounceTimer = setTimeout(() => {
        executeSearch(term);
    }, 600); 
  });

  searchInput.addEventListener('keydown', (e) => {
      if (e.key === 'Enter' && !checkGuest.checked) {
          e.preventDefault(); 
          clearTimeout(debounceTimer); 
          const term = searchInput.value.trim();
          executeSearch(term); 
      }
  });

  async function fetchUsers(term) {
    
    try {
      const searchType = searchTypeSelect.value;
      const isRegex = !checkExact.checked; 

      const payload = {
        params: [term.trim()],
        type: searchType,
        likeop: isRegex
      };

      const response = await fetch('/findUser', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload)
      });

      if (!response.ok) {
        if (response.status === 401) {  
            window.location.href = '/login';
        }
        throw new Error('Error en búsqueda');
      }
      
      const data = await response.json();
      
      cachedUsers = Object.entries(data).map(([id, userData]) => ({
        id: id,
        ...userData
      }));

      renderSuggestions(cachedUsers);

    } catch (err) {
      console.error("Error buscando usuarios:", err);
      cachedUsers = []; 
      renderSuggestions([]); 
    }
  }

  function filterLocalSuggestions(term) {
    const searchField = getSearchFieldFromType(searchTypeSelect.value);
    const lowerTerm = term.toLowerCase();

    const filtered = cachedUsers.filter(u => {
      const value = u[searchField];
      return value && value.toString().toLowerCase().includes(lowerTerm);
    });

    renderSuggestions(filtered);
  }

  function getSearchFieldFromType(type) {
    switch (type) {
      case 'nameUser': return 'Name';
      case 'lastNameUser': return 'LastName';
      case 'aliasUser': return 'Alias';
      case 'emailUser': return 'Email';
      case 'phoneUser': return 'Phone'; 
      default: return 'Name';
    }
  }

  function renderSuggestions(users) {
    suggestionsList.innerHTML = '';
    
    if (!users || users.length === 0) {
      suggestionsList.classList.add('hidden');
      return;
    }

    users.forEach(u => {
      const div = document.createElement('div');
      div.className = 'suggestion-item';
      
      const displayName = u.name && u.lastname ? `${u.name} ${u.lastname}` : (u.alias || "Usuario");
      const displayEmail = u.email || "Sin email";

      div.innerHTML = `
        <div class="flex justify-between items-center">
            <span class="font-bold text-sm">${displayName}</span>
            <span class="text-xs text-secondary opacity-70">${u.Role || 'User'}</span>
        </div>
        <div class="text-xs text-secondary">${displayEmail}</div>
      `;

      div.addEventListener('click', () => {
        selectUser(u);
      });

      suggestionsList.appendChild(div);
    });

    suggestionsList.classList.remove('hidden');
  }

  function selectUser(user) {
    selectedUserObj = user;
    searchInput.value = ''; 
    suggestionsList.classList.add('hidden'); 
    
    const displayName = user.name && user.lastname ? `${user.name} ${user.lastname}` : (user.alias || user.email);
    selectedUserName.textContent = displayName;
    selectedUserDisplay.classList.remove('hidden');
    searchInput.disabled = true; 
  }

  function clearSelection() {
    selectedUserObj = null;
    selectedUserDisplay.classList.add('hidden');
    searchInput.disabled = false;
    searchInput.focus();
  }

  clearSelectionBtn.addEventListener('click', clearSelection);

  // ==============================
  // 🔹 RENDERIZAR TABLA DE FIRMANTES
  // ==============================
  function renderTable() {
    const tbody = document.getElementById('reviewersTableBody');
    tbody.innerHTML = '';

    reviewers.forEach((rev, idx) => {
      const tr = document.createElement('tr');
      tr.classList.add("border-t", "border-t-white/10");

      const roleClass = rev.role === 1 
          ? "bg-green-500/20 text-green-300 border border-green-500/30" 
          : "bg-blue-500/20 text-blue-300 border border-blue-500/30";

      // Nota: rev.external mapea a IsExternal del struct Go
      const typeLabel = rev.external 
          ? '<span class="text-purple-400 text-xs border border-purple-500/30 px-2 py-0.5 rounded-full">Externo</span>' 
          : '<span class="text-blue-400 text-xs border border-blue-500/30 px-2 py-0.5 rounded-full">Interno</span>';

      tr.innerHTML = `
        <td class="px-4 py-2 w-[350px] text-primary text-sm font-normal">
            <div class="flex flex-col">
                <span class="font-bold">${rev.displayName}</span>
                <span class="text-xs text-secondary">${rev.displayEmail}</span>
            </div>
        </td>
        <td class="px-4 py-2 w-40 text-sm font-normal">
          <button class="toggle-role w-full rounded-lg px-3 py-1 text-xs font-bold transition-all ${roleClass}" data-index="${idx}">
            ${rev.role === 1 ? "Firmante" : "Visor"}
          </button>
        </td>
        <td class="px-4 py-2 w-[200px] text-secondary text-sm font-normal">
          <span class="editable-date cursor-pointer hover:text-white border-b border-dashed border-white/20 pb-0.5" data-index="${idx}">
            ${rev.due_date}
          </span>
        </td>
        <td class="px-4 py-2 w-[150px] text-secondary text-sm font-normal">
            ${typeLabel}
        </td>
        <td class="px-4 py-2 text-secondary text-sm font-normal">
          <span class="editable-comment cursor-pointer hover:text-white" data-index="${idx}">
            ${rev.comment || '<span class="italic opacity-50">Sin comentario...</span>'}
          </span>
        </td>
        <td class="px-4 py-2 w-10">
            <button class="delete-reviewer text-red-400 hover:text-red-300 transition-colors" data-index="${idx}">
                <span class="material-symbols-outlined text-lg">delete</span>
            </button>
        </td>
      `;
      tbody.appendChild(tr);
    });

    // Listeners de tabla
    document.querySelectorAll('.toggle-role').forEach(btn => {
      btn.addEventListener('click', () => {
        const idx = btn.dataset.index;
        reviewers[idx].role = reviewers[idx].role === 1 ? 0 : 1;
        renderTable();
      });
    });

    document.querySelectorAll('.delete-reviewer').forEach(btn => {
      btn.addEventListener('click', () => {
        const idx = btn.dataset.index;
        reviewers.splice(idx, 1);
        renderTable();
      });
    });

    document.querySelectorAll('.editable-date').forEach(span => {
      span.addEventListener('click', () => {
        const idx = span.dataset.index;
        const input = document.createElement('input');
        input.type = 'date';
        input.value = reviewers[idx].due_date;
        input.className = "bg-slate-800 text-white text-sm border border-slate-600 rounded p-1";

        const saveDate = () => {
             if(input.value) reviewers[idx].due_date = input.value;
             renderTable();
        };

        input.addEventListener('blur', saveDate);
        input.addEventListener('keydown', (e) => {
          if (e.key === 'Enter') input.blur();
        });
        span.replaceWith(input);
        input.focus();
      });
    });

    document.querySelectorAll('.editable-comment').forEach(span => {
      span.addEventListener('click', () => {
        const idx = span.dataset.index;
        const textarea = document.createElement('textarea');
        textarea.value = reviewers[idx].comment || '';
        textarea.className = "bg-slate-800 text-white text-sm w-full h-16 p-2 rounded resize-none border border-slate-600";
        
        const saveComment = () => {
          reviewers[idx].comment = textarea.value;
          renderTable();
        };

        textarea.addEventListener('blur', saveComment);
        textarea.addEventListener('keydown', (e) => {
          if (e.key === 'Enter' && e.ctrlKey) textarea.blur();
        });
        span.replaceWith(textarea);
        textarea.focus();
      });
    });
  }

  // Inicializar tabla si hay datos previos
  renderTable();

  // ==============================
  // 🔹 EVENTO AÑADIR REVISOR
  // ==============================
  document.getElementById('addReviewerBtn').addEventListener('click', () => {
    const comment = document.getElementById('userComment');
    const today = new Date().toISOString().split('T')[0];
    
    let newReviewer = null;

    if (checkGuest.checked) {
        // --- CASO EXTERNO (INVITADO) ---
        const guestEmail = searchInput.value.trim();
        const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
        
        if (!guestEmail || !emailRegex.test(guestEmail)) {
            alert("Por favor, introduce un correo electrónico válido para el invitado.");
            return;
        }

        newReviewer = {
            // Campos requeridos por el Struct GO
            external: true,    // json:"external"
            user: "-1",  // json:"user"
            role: 1,           // json:"role"
            alive_proof: checkAlive.checked, // json:"alive_proof"
            due_date: today,   // json:"due_date"
            comment: comment.value, // json:"comment"
            positions: [],     // json:"positions"
            
            // Campos auxiliares para UI (se ignoran en backend si no están en struct)
            displayName: "Invitado Externo",
            displayEmail: guestEmail
        };

    } else {
        // --- CASO INTERNO ---
        if (!selectedUserObj) {
            alert("Por favor, busca y selecciona un usuario de la lista.");
            return;
        }

        newReviewer = {
            // Campos requeridos por el Struct GO
            external: false,          // json:"external"
            user: selectedUserObj.id, // json:"user" (Aquí va el ID para internos)
            role: 1,                  // json:"role"
            alive_proof: checkAlive.checked, // json:"alive_proof"
            due_date: today,          // json:"due_date"
            comment: comment.value,   // json:"comment"
            positions: [],            // json:"positions"

            // Campos auxiliares para UI
            displayName: (selectedUserObj.name && selectedUserObj.lastname) ? `${selectedUserObj.name} ${selectedUserObj.lastname}` : selectedUserObj.alias,
            displayEmail: selectedUserObj.email
        };
    }

    // Evitar duplicados (por user/id)
    const exists = reviewers.some(r => r.user === newReviewer.user);
    if (exists) {
        alert("Este usuario ya ha sido añadido a la lista.");
        return;
    }

    reviewers.push(newReviewer);
    renderTable();
    
    // Limpiar campos
    checkAlive.checked=false;
    comment.value = '';
    if (checkGuest.checked) {
        searchInput.value = '';
    } else {
        clearSelection();
    }
  });

  // ==============================
  // 🔹 GUARDAR Y FINALIZAR
  // ==============================
  document.getElementById('uploadBtn').addEventListener('click', async () => {
    if (reviewers.length === 0) {
        if(!confirm("No has añadido ningún firmante. ¿Deseas continuar de todas formas?")) return;
    }

    currentDoc.reviewers = reviewers;
    folderData[currentDocId] = currentDoc;
    sessionStorage.setItem("folder", JSON.stringify(folderData));

    if (currentIndex + 1 < docIds.length) {
      sessionStorage.setItem("currentIndex", currentIndex + 1);
      window.location.href = "/addSigners";
    } else {
      sessionStorage.removeItem("currentIndex");
      
      const inviteRequest = {};
      docIds.forEach(docId => {
        const doc = folderData[docId];
        inviteRequest[docId] = {
          idfolder: doc.idFolder, 
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
          // El array de reviewers ahora contiene objetos con la estructura correcta
          // (external, user, role, due_date, comment, positions)
          reviewers: doc.reviewers || []
        };
      });

      sessionStorage.setItem("inviteRequest", JSON.stringify(inviteRequest));
      sessionStorage.removeItem("folder");
      window.location.href = "/addSignatures";
    }
  });

  // ==============================
  // 🔹 PREVISUALIZACIÓN DEL DOCUMENTO PDF
  // ==============================
  const loadingSpinner = document.getElementById('loading-spinner');
  const pdfFrame = document.getElementById('pdf-frame');
  const previewOverlay = document.getElementById('preview-overlay');
  
  (async () => {
    try {
      const docs = {idDocs: [currentDocId], type: "uploaded"}
      const response = await fetch('/downloadDoc', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(docs)
      });

      if (!response.ok) {
        if (response.status === 401) {  
            window.location.href = '/login';
        }
        throw new Error(`HTTP ${response.status}`);
      }

      const blob = await response.blob();
      const url = URL.createObjectURL(blob);
      
      pdfFrame.src = url;
      
      loadingSpinner.classList.add('hidden');
      pdfFrame.classList.remove('hidden');
      previewOverlay.classList.remove('hidden');
      previewOverlay.classList.add('flex');

      previewOverlay.onclick = () => window.open(url, '_blank');
      
    } catch (err) {
      console.error('Error al cargar documento:', err);
      loadingSpinner.innerHTML = `
        <div class="text-red-400 text-center">
          <span class="material-symbols-outlined text-4xl mb-2">error</span>
          <p>Error al cargar el documento.</p>
        </div>
      `;
    }
  })();
});