document.addEventListener('DOMContentLoaded', () => {
  const reviewers = [];
  document.getElementById('uploadBtn').addEventListener('click', async () => {
    const docType = document.getElementById('docType');
    const docDesc = document.getElementById('docDesc');
    const fileInput = document.getElementById('fileInput');
    const user = document.getElementById('userSelect');
    const team = document.getElementById('teamSelect');
    if (!fileInput.files.length) {
        console.error("No file selected");
        return;
    }

    const file = fileInput.files[0];
    const base64 = await toBase64(file);
    
    const tbody = document.getElementById('reviewersTableBody');
    const rows = tbody.querySelectorAll('tr');
    const reviewers = [];

    rows.forEach(row => {
        const cells = row.querySelectorAll('td');
        const user = cells[0]?.textContent.trim();
        const role = cells[1]?.querySelector('button')?.textContent.trim();
        const dueDate = cells[2]?.textContent.trim();
        const team = cells[3]?.textContent.trim();

        if (user && role && dueDate && team) {
            reviewers.push({
                user,
                role,
                due_date: dueDate,
                team,
                status: "Not Started" // Estado inicial por defecto
            });
        }

    });
    
    const documentData = {
        type: docType.value,
        description: docDesc.value,
        fileName: file.name,
        fileSize: file.size,
        fileMime: file.type,
        fileBase64: base64.split(',')[1],
        reviewers: reviewers,
        uploadDate: new Date().toLocaleDateString('en-US', {
            year: 'numeric',
            month: 'long',
            day: 'numeric'
        }),
        uploadedBy: "Alex" // Puedes cambiar esto por el usuario actual
    };

    // Guardar en localStorage
    localStorage.setItem('currentDocument', JSON.stringify(documentData));
    // Resetear el input
    tbody.innerHTML = '';
    fileInput.value = '';
    docType.value = '';
    docDesc.value = '';
    user.value = '';
    team.value = '';
    // Redirigir a document_viewer.html
    window.location.href = 'document_viewer.html';
});

// Función para convertir archivo a base64
function toBase64(file) {
    return new Promise((resolve, reject) => {
        const reader = new FileReader();
        reader.readAsDataURL(file);
        reader.onload = () => resolve(reader.result);
        reader.onerror = error => reject(error);
    });
}

// Event listener para el input file
document.getElementById('fileInput').addEventListener('change', function (event) {
  let fileURL = null;
  const previewContainer = document.getElementById('preview-container');
  const dropMessage = document.getElementById('drop-message');
  const fileInput = document.getElementById('fileInput');

  const file = event.target.files[0];
  if (file && file.type === 'application/pdf') {
    document.getElementById('docName').value = file.name;
    const reader = new FileReader();
    
    reader.onload = function (e) {
      if (fileURL) URL.revokeObjectURL(fileURL);
      fileURL = URL.createObjectURL(file);

      // Ocultar mensaje de arrastrar
      if (dropMessage) dropMessage.style.display = 'none';

      // Limpiar solo el contenido del PDF, manteniendo el botón
      const existingWrapper = previewContainer.querySelector('.pdf-wrapper');
      if (existingWrapper) previewContainer.removeChild(existingWrapper);

      // Wrapper relativo para el PDF
      const wrapper = document.createElement('div');
      wrapper.className = 'pdf-wrapper';
      wrapper.style.position = 'relative';
      wrapper.style.width = '100%';
      wrapper.style.height = '100%';

      // Iframe con el PDF
      const iframe = document.createElement('iframe');
      iframe.src = `${fileURL}#page=1&zoom=50%`;
      iframe.style.width = '100%';
      iframe.style.height = '100%';
      iframe.style.border = 'none';

      // Div flotante para "Ver completo"
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

      overlay.onclick = () => window.open(fileURL, '_blank');

      wrapper.appendChild(iframe);
      wrapper.appendChild(overlay);
      previewContainer.appendChild(wrapper);
    };
    
    reader.readAsArrayBuffer(file);
  }
});

// Drag and drop opcional
const previewContainer = document.getElementById('preview-container');
previewContainer.addEventListener('dragover', (e) => {
  e.preventDefault();
  previewContainer.classList.add('border-blue-400');
});

previewContainer.addEventListener('dragleave', () => {
  previewContainer.classList.remove('border-blue-400');
});

previewContainer.addEventListener('drop', (e) => {
  e.preventDefault();
  previewContainer.classList.remove('border-blue-400');
  
  if (e.dataTransfer.files.length) {
    document.getElementById('fileInput').files = e.dataTransfer.files;
    const event = new Event('change');
    document.getElementById('fileInput').dispatchEvent(event);
  }
});

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
          ${rev.role}
        </button>
      </td>
      <td class="px-4 py-2 w-[400px] text-secondary text-sm font-normal">
        <span class="editable-date" data-index="${idx}">${rev.due_date}</span>
      </td>
      <td class="px-4 py-2 w-[400px] text-secondary text-sm font-normal">${rev.team}</td>
    `;
    tbody.appendChild(tr);
  });

  // Agrega listeners para cambiar rol
  document.querySelectorAll('.toggle-role').forEach(btn => {
    btn.addEventListener('click', () => {
      const idx = btn.dataset.index;
      reviewers[idx].role = reviewers[idx].role === 'Signer' ? 'Viewer' : 'Signer';
      renderTable();
    });
  });

  // Agrega listeners para cambiar fecha
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

      span.replaceWith(input);
      input.focus();
    });
  });
}

document.getElementById('addReviewerBtn').addEventListener('click', () => {
  const user = document.getElementById('userSelect');
  const team = document.getElementById('teamSelect');
  const today = new Date().toISOString().split('T')[0];

  if (!user || !team) return alert("Select both user and team.");

  reviewers.push({
    user: user.value,
    role: 'Signer',
    due_date: today,
    team: team.value
  });

  renderTable();
  user.value = '';
});
});