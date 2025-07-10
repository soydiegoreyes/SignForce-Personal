document.addEventListener('DOMContentLoaded', () => {
    const reviewers = [];
document.getElementById('uploadBtn').addEventListener('click', async () => {
    const docName = document.getElementById('docName').value;
    const docType = document.getElementById('docType').value;
    const docDesc = document.getElementById('docDesc').value;
    const fileInput = document.getElementById('fileInput');

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
                team
            });
        }
    });

    const data = {
      name: docName,
      type: docType,
      description: docDesc,
      fileName: file.name,
      fileSize: file.size,
      fileMime: file.type,
      fileBase64: base64.split(',')[1],
      reviewers: reviewers
    };

    console.log("Upload Data:", data);
});

function toBase64(file) {
    return new Promise((resolve, reject) => {
        const reader = new FileReader();
        reader.readAsDataURL(file);
        reader.onload = () => resolve(reader.result);
        reader.onerror = error => reject(error);
    });
}


function toBase64(file) {
    return new Promise((resolve, reject) => {
        const reader = new FileReader();
        reader.readAsDataURL(file);
        reader.onload = () => resolve(reader.result);
        reader.onerror = error => reject(error);
    });
}

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
  const user = document.getElementById('userSelect').value;
  const team = document.getElementById('teamSelect').value;
  const today = new Date().toISOString().split('T')[0];

  if (!user || !team) return alert("Select both user and team.");

  reviewers.push({
    user: user,
    role: 'Signer',
    due_date: today,
    team: team
  });

  renderTable();
});
});