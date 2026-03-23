/* SignForce — Demo Data v1
   Populates all pages with realistic dummy data when API returns empty */
(function() {
  window.SF_DEMO = true;

  // ═══ Demo Documents ═══
  window.SF_DEMO_DOCS = {
    "DOC-2024-001": {
      id: "DOC-2024-001", nameDoc: "Contrato de Arrendamiento - Local 42", status: 3,
      createDate: "2026-03-15T10:30:00Z", signers: ["María González", "Carlos López"],
      tags: ["contrato", "arrendamiento"], pages: 12, fileSize: "2.4 MB"
    },
    "DOC-2024-002": {
      id: "DOC-2024-002", nameDoc: "Acuerdo de Confidencialidad NDA", status: 5,
      createDate: "2026-03-18T14:20:00Z", signers: ["Ana Ramírez"],
      tags: ["nda", "confidencialidad"], pages: 4, fileSize: "890 KB"
    },
    "DOC-2024-003": {
      id: "DOC-2024-003", nameDoc: "Poder Notarial General", status: 1,
      createDate: "2026-03-20T09:15:00Z", signers: ["Roberto Méndez", "Laura Vázquez", "Pedro Sánchez"],
      tags: ["poder", "notarial"], pages: 8, fileSize: "1.7 MB"
    },
    "DOC-2024-004": {
      id: "DOC-2024-004", nameDoc: "Contrato de Servicios Profesionales", status: 5,
      createDate: "2026-03-10T16:45:00Z", signers: ["Diana Torres"],
      tags: ["contrato", "servicios"], pages: 6, fileSize: "1.2 MB"
    },
    "DOC-2024-005": {
      id: "DOC-2024-005", nameDoc: "Acta Constitutiva - Tech Solutions SA", status: 3,
      createDate: "2026-03-22T11:00:00Z", signers: ["Fernando Ruiz", "Gabriela Flores"],
      tags: ["acta", "constitutiva"], pages: 24, fileSize: "5.1 MB"
    },
    "DOC-2024-006": {
      id: "DOC-2024-006", nameDoc: "Convenio de Terminación Laboral", status: 2,
      createDate: "2026-03-21T08:30:00Z", signers: ["Miguel Ángel Herrera"],
      tags: ["convenio", "laboral"], pages: 3, fileSize: "650 KB"
    },
    "DOC-2024-007": {
      id: "DOC-2024-007", nameDoc: "Contrato de Compraventa Inmueble", status: 5,
      createDate: "2026-03-05T13:20:00Z", signers: ["Alejandra Morales", "Juan Pérez"],
      tags: ["compraventa", "inmueble"], pages: 18, fileSize: "3.8 MB"
    },
    "DOC-2024-008": {
      id: "DOC-2024-008", nameDoc: "Carta Poder Simple", status: 4,
      createDate: "2026-03-19T17:10:00Z", signers: ["Sofía Castillo"],
      tags: ["carta", "poder"], pages: 2, fileSize: "420 KB"
    }
  };

  // Status labels
  window.SF_STATUS_MAP = {
    1: {label: "Pendiente", color: "#fbbf24", bg: "rgba(251,191,36,.1)", border: "rgba(251,191,36,.15)"},
    2: {label: "En proceso", color: "#60a5fa", bg: "rgba(96,165,250,.1)", border: "rgba(96,165,250,.15)"},
    3: {label: "Por firmar", color: "#a78bfa", bg: "rgba(167,139,250,.1)", border: "rgba(167,139,250,.15)"},
    4: {label: "Rechazado", color: "#f87171", bg: "rgba(248,113,113,.1)", border: "rgba(248,113,113,.15)"},
    5: {label: "Completado", color: "#34d399", bg: "rgba(52,211,153,.1)", border: "rgba(52,211,153,.15)"}
  };

  // ═══ Demo Folders ═══
  window.SF_DEMO_FOLDERS = [
    {id: 1, name: "Contratos 2026", docs: 12, created: "2026-01-15", color: "#B5C413", icon: "folder"},
    {id: 2, name: "NDAs Clientes", docs: 8, created: "2026-02-03", color: "#60a5fa", icon: "folder"},
    {id: 3, name: "Poderes Notariales", docs: 5, created: "2026-02-20", color: "#a78bfa", icon: "folder"},
    {id: 4, name: "Documentos Internos", docs: 15, created: "2026-01-08", color: "#34d399", icon: "folder"},
    {id: 5, name: "Recursos Humanos", docs: 22, created: "2025-11-12", color: "#fbbf24", icon: "folder"},
    {id: 6, name: "Legal / Compliance", docs: 9, created: "2026-03-01", color: "#f87171", icon: "gavel"}
  ];

  // ═══ Demo Templates ═══
  window.SF_DEMO_TEMPLATES = [
    {id: 1, name: "Contrato de Arrendamiento", category: "Inmobiliario", uses: 34, fields: 12, updated: "2026-03-10"},
    {id: 2, name: "NDA Estándar", category: "Legal", uses: 67, fields: 6, updated: "2026-03-18"},
    {id: 3, name: "Contrato de Servicios", category: "Comercial", uses: 45, fields: 10, updated: "2026-02-28"},
    {id: 4, name: "Poder Notarial", category: "Legal", uses: 12, fields: 8, updated: "2026-03-05"},
    {id: 5, name: "Carta Responsiva", category: "General", uses: 28, fields: 5, updated: "2026-03-15"},
    {id: 6, name: "Acta Constitutiva", category: "Corporativo", uses: 8, fields: 18, updated: "2026-01-20"},
    {id: 7, name: "Convenio Laboral", category: "RH", uses: 19, fields: 9, updated: "2026-03-22"},
    {id: 8, name: "Pagaré Simple", category: "Financiero", uses: 41, fields: 7, updated: "2026-03-12"}
  ];

  // ═══ Demo Keys ═══
  window.SF_DEMO_KEYS = [
    {id: 1, name: "Firma Principal", type: "e.firma (SAT)", status: "active", created: "2026-01-10", expires: "2028-01-10", lastUsed: "2026-03-22"},
    {id: 2, name: "Certificado Digital", type: "CSD", status: "active", created: "2025-06-15", expires: "2027-06-15", lastUsed: "2026-03-20"},
    {id: 3, name: "Firma Avanzada FIEL", type: "FIEL", status: "expired", created: "2023-03-01", expires: "2025-03-01", lastUsed: "2025-02-28"}
  ];

  // ═══ Demo Dashboard Stats ═══
  window.SF_DEMO_STATS = {
    totalDocs: 156,
    signed: 89,
    pending: 34,
    rejected: 12,
    inProcess: 21,
    thisMonth: 28,
    lastMonth: 31,
    users: 8,
    storage: "2.4 GB",
    storagePercent: 48
  };

  // ═══ Demo Activity ═══
  window.SF_DEMO_ACTIVITY = [
    {action: "Documento firmado", doc: "Contrato de Arrendamiento", user: "María González", time: "Hace 2 horas", icon: "draw", color: "#34d399"},
    {action: "Documento enviado", doc: "NDA - Proyecto Alfa", user: "Carlos López", time: "Hace 4 horas", icon: "send", color: "#60a5fa"},
    {action: "Nueva plantilla", doc: "Convenio Laboral v2", user: "Ana Ramírez", time: "Hace 6 horas", icon: "description", color: "#a78bfa"},
    {action: "Documento rechazado", doc: "Poder Notarial", user: "Roberto Méndez", time: "Ayer", icon: "cancel", color: "#f87171"},
    {action: "Firma completada", doc: "Contrato de Servicios", user: "Diana Torres", time: "Ayer", icon: "verified", color: "#34d399"},
    {action: "Carpeta creada", doc: "Legal / Compliance", user: "Leonardo", time: "Hace 2 días", icon: "create_new_folder", color: "#fbbf24"},
    {action: "Documento subido", doc: "Acta Constitutiva", user: "Fernando Ruiz", time: "Hace 3 días", icon: "upload_file", color: "#B5C413"}
  ];

  // ═══ Demo Reports ═══
  window.SF_DEMO_REPORTS = {
    monthly: [
      {month: "Oct", docs: 18, signed: 14},
      {month: "Nov", docs: 24, signed: 19},
      {month: "Dic", docs: 31, signed: 27},
      {month: "Ene", docs: 22, signed: 18},
      {month: "Feb", docs: 28, signed: 23},
      {month: "Mar", docs: 33, signed: 25}
    ],
    topSigners: [
      {name: "María González", count: 23},
      {name: "Carlos López", count: 18},
      {name: "Ana Ramírez", count: 15},
      {name: "Diana Torres", count: 12},
      {name: "Roberto Méndez", count: 9}
    ],
    byType: [
      {type: "Contratos", count: 45, pct: 29},
      {type: "NDAs", count: 32, pct: 21},
      {type: "Poderes", count: 24, pct: 15},
      {type: "Actas", count: 18, pct: 12},
      {type: "Convenios", count: 15, pct: 10},
      {type: "Otros", count: 22, pct: 14}
    ]
  };
})();
