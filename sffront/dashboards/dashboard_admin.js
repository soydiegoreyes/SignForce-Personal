document.addEventListener('DOMContentLoaded', () => {
    // ==============================
    // 🔹 VARIABLES GLOBALES
    // ==============================
    let currentPage = 1;
    let pageSize = 10;
    let totalUsers = 0;
    let cachedUsers = [];
    let selectedUserObj = null;
    let debounceTimer = null;

    // ==============================
    // 🔹 NAVEGACIÓN ENTRE SECCIONES
    // ==============================
    const navLinks = document.querySelectorAll('[data-tab]');
    const sections = document.querySelectorAll('section[id$="-section"]');

    navLinks.forEach(link => {
        link.addEventListener('click', (e) => {
            e.preventDefault();
            const targetTab = link.dataset.tab;
            
            // Actualizar clases activas
            navLinks.forEach(l => l.classList.remove('bg-gradient-to-r', 'from-blue-500', 'to-purple-600', 'text-white'));
            navLinks.forEach(l => l.classList.add('text-secondary', 'hover:bg-white/10'));
            
            link.classList.remove('text-secondary', 'hover:bg-white/10');
            link.classList.add('bg-gradient-to-r', 'from-blue-500', 'to-purple-600', 'text-white');
            
            // Mostrar sección correspondiente
            sections.forEach(section => {
                section.classList.remove('active-section');
                section.classList.add('hidden-section');
            });
            
            const targetSection = document.getElementById(`${targetTab}-section`);
            if (targetSection) {
                targetSection.classList.remove('hidden-section');
                targetSection.classList.add('active-section');
                
                // Cargar datos específicos de la sección
                switch(targetTab) {
                    case 'dashboard':
                        loadDashboardData();
                        break;
                    case 'user-management':
                        loadUsersList();
                        break;
                    case 'document-stats':
                        loadDocumentStats();
                        break;
                    // Agregar más casos según sea necesario
                }
            }
        });
    });

    // ==============================
    // 🔹 FUNCIONES DE DATOS DEL DASHBOARD
    // ==============================
    async function loadDashboardData() {
        try {
            // Cargar estadísticas generales
            const stats = await fetchAdminStats();
            updateDashboardStats(stats);
            
            // Cargar usuarios recientes
            await loadRecentUsers();
            
            // Cargar actividad reciente (simulada)
            simulateActivityChart();
            
        } catch (error) {
            console.error('Error cargando dashboard:', error);
            showToast('Error al cargar datos del dashboard', 'error');
        }
    }

    // Función ficticia para obtener estadísticas del admin
    async function fetchAdminStats() {
        // NOTA: Esta es una función ficticia que simula la respuesta del backend
        // En producción, reemplazar con fetch real a tu API
        return new Promise((resolve) => {
            setTimeout(() => {
                resolve({
                    totalUsers: 154,
                    usersGrowth: 12.5,
                    signedDocs: 1287,
                    docsGrowth: 8.3,
                    pendingDocs: 42,
                    pendingInvites: 7,
                    activeUsers: 89
                });
            }, 500);
        });
    }

    function updateDashboardStats(stats) {
        document.getElementById('totalUsers').textContent = stats.totalUsers;
        document.getElementById('usersGrowth').textContent = `${stats.usersGrowth}%`;
        document.getElementById('signedDocs').textContent = stats.signedDocs.toLocaleString();
        document.getElementById('docsGrowth').textContent = `${stats.docsGrowth}%`;
        document.getElementById('pendingDocs').textContent = stats.pendingDocs;
        document.getElementById('pendingInvites').textContent = stats.pendingInvites;
        document.getElementById('activeUsersCount').textContent = stats.activeUsers;
    }

    async function loadRecentUsers() {
        try {
            // Simular carga de usuarios recientes
            const users = await fetchRecentUsers();
            renderRecentUsersTable(users);
            
        } catch (error) {
            console.error('Error cargando usuarios recientes:', error);
        }
    }

    // Función ficticia para obtener usuarios recientes
    async function fetchRecentUsers() {
        return new Promise((resolve) => {
            setTimeout(() => {
                resolve([
                    {
                        id: "USR001",
                        name: "Juan Pérez",
                        email: "juan.perez@empresa.com",
                        role: "Administrador",
                        joinDate: "2024-01-15",
                        status: "active"
                    },
                    {
                        id: "USR002",
                        name: "María González",
                        email: "maria.gonzalez@empresa.com",
                        role: "Usuario",
                        joinDate: "2024-01-10",
                        status: "active"
                    },
                    {
                        id: "USR003",
                        name: "Carlos López",
                        email: "carlos.lopez@empresa.com",
                        role: "Revisor",
                        joinDate: "2024-01-05",
                        status: "pending"
                    },
                    {
                        id: "USR004",
                        name: "Ana Martínez",
                        email: "ana.martinez@empresa.com",
                        role: "Usuario",
                        joinDate: "2024-01-02",
                        status: "active"
                    },
                    {
                        id: "USR005",
                        name: "Pedro Sánchez",
                        email: "pedro.sanchez@empresa.com",
                        role: "Usuario",
                        joinDate: "2023-12-28",
                        status: "inactive"
                    }
                ]);
            }, 300);
        });
    }

    function renderRecentUsersTable(users) {
        const tbody = document.getElementById('recentUsersTable');
        tbody.innerHTML = '';
        
        users.forEach(user => {
            const statusClass = getStatusClass(user.status);
            const statusText = getStatusText(user.status);
            
            const tr = document.createElement('tr');
            tr.className = 'border-b border-white/5 hover:bg-white/5';
            tr.innerHTML = `
                <td class="px-4 py-3">
                    <div class="flex items-center gap-3">
                        <div class="w-8 h-8 rounded-full bg-blue-500/20 flex items-center justify-center">
                            <span class="text-sm font-bold">${user.name.charAt(0)}</span>
                        </div>
                        <span class="font-medium">${user.name}</span>
                    </div>
                </td>
                <td class="px-4 py-3 text-secondary">${user.email}</td>
                <td class="px-4 py-3">
                    <span class="badge badge-active">${user.role}</span>
                </td>
                <td class="px-4 py-3 text-secondary">${formatDate(user.joinDate)}</td>
                <td class="px-4 py-3">
                    <span class="badge ${statusClass}">${statusText}</span>
                </td>
                <td class="px-4 py-3">
                    <button class="text-blue-400 hover:text-blue-300 transition-colors view-user-btn" data-id="${user.id}">
                        <span class="material-symbols-outlined text-sm">visibility</span>
                    </button>
                </td>
            `;
            tbody.appendChild(tr);
        });
        
        // Agregar event listeners para botones de vista
        document.querySelectorAll('.view-user-btn').forEach(btn => {
            btn.addEventListener('click', (e) => {
                const userId = e.currentTarget.dataset.id;
                showUserDetails(userId);
            });
        });
    }

    function simulateActivityChart() {
        const chartContainer = document.getElementById('activityChart');
        // En producción, aquí se integraría una librería como Chart.js
        chartContainer.innerHTML = `
            <div class="text-center">
                <div class="text-lg font-bold text-primary mb-2">Actividad de usuarios</div>
                <div class="text-secondary text-sm mb-6">Últimas 24 horas</div>
                <div class="flex items-end justify-center gap-2 h-32">
                    ${[65, 80, 45, 90, 70, 85, 60].map((height, i) => `
                        <div class="flex flex-col items-center">
                            <div class="w-6 bg-gradient-to-t from-blue-500 to-purple-600 rounded-t" style="height: ${height}%"></div>
                            <span class="text-xs text-secondary mt-2">${['00', '04', '08', '12', '16', '20', '24'][i]}:00</span>
                        </div>
                    `).join('')}
                </div>
            </div>
        `;
    }

    // ==============================
    // 🔹 GESTIÓN DE USUARIOS - BÚSQUEDA
    // ==============================
    const userSearchInput = document.getElementById('userSearchInput');
    const userSearchType = document.getElementById('userSearchType');
    const searchUserBtn = document.getElementById('searchUserBtn');
    const clearSearchBtn = document.getElementById('clearSearchBtn');
    const userSuggestionsList = document.getElementById('userSuggestionsList');
    const addUserBtn = document.getElementById('addUserBtn');

    // Búsqueda en tiempo real
    userSearchInput.addEventListener('input', (e) => {
        const term = e.target.value.trim();
        
        if (term.length === 0) {
            clearTimeout(debounceTimer);
            userSuggestionsList.classList.add('hidden');
            cachedUsers = [];
            return;
        }

        clearTimeout(debounceTimer);
        debounceTimer = setTimeout(() => {
            if (term.length >= 3) {
                searchUsers(term);
            }
        }, 500);
    });

    // Búsqueda manual con botón
    searchUserBtn.addEventListener('click', async () => {
        const term = userSearchInput.value.trim();
        if (term.length >= 3) {
            await searchUsers(term, true);
        } else {
            showToast('Ingresa al menos 3 caracteres para buscar', 'warning');
        }
    });

    clearSearchBtn.addEventListener('click', () => {
        userSearchInput.value = '';
        userSuggestionsList.classList.add('hidden');
        cachedUsers = [];
        selectedUserObj = null;
        loadUsersList(); // Recargar lista completa
    });

    // Función de búsqueda de usuarios (similar a add_signers.js)
    async function searchUsers(term, showResults = false) {
        try {
            const searchType = userSearchType.value;
            
            // NOTA: Esta es la misma estructura que en add_signers.js
            // En producción, usar el mismo endpoint /findUser
            const payload = {
                params: [term],
                type: searchType,
                likeop: true // Búsqueda flexible por defecto
            };

            // Simular respuesta del servidor
            const mockResponse = {
                "USR001": { name: "Juan", lastname: "Pérez", email: "juan.perez@empresa.com", alias: "juanp", role: "Admin" },
                "USR002": { name: "María", lastname: "González", email: "maria.gonzalez@empresa.com", alias: "mariag", role: "User" },
                "USR003": { name: "Carlos", lastname: "López", email: "carlos.lopez@empresa.com", alias: "carlosl", role: "Reviewer" }
            };

            // En producción, descomentar esto:
            /*
            const response = await fetch('/findUser', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(payload)
            });
            
            if (!response.ok) throw new Error('Error en búsqueda');
            const data = await response.json();
            */

            const data = mockResponse; // Usar mock para demo
            
            cachedUsers = Object.entries(data).map(([id, userData]) => ({
                id: id,
                ...userData
            }));

            if (showResults) {
                // Si es búsqueda manual, mostrar resultados en tabla
                renderUsersTable(cachedUsers);
            } else {
                // Si es búsqueda en tiempo real, mostrar sugerencias
                renderUserSuggestions(cachedUsers);
            }

        } catch (err) {
            console.error("Error buscando usuarios:", err);
            showToast('Error al buscar usuarios', 'error');
            cachedUsers = [];
            userSuggestionsList.classList.add('hidden');
        }
    }

    function renderUserSuggestions(users) {
        userSuggestionsList.innerHTML = '';
        
        if (!users || users.length === 0) {
            userSuggestionsList.classList.add('hidden');
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
                    <span class="text-xs text-secondary opacity-70">${u.role || 'User'}</span>
                </div>
                <div class="text-xs text-secondary">${displayEmail}</div>
            `;

            div.addEventListener('click', () => {
                selectUserForManagement(u);
            });

            userSuggestionsList.appendChild(div);
        });

        userSuggestionsList.classList.remove('hidden');
    }

    function selectUserForManagement(user) {
        selectedUserObj = user;
        userSearchInput.value = user.email;
        userSuggestionsList.classList.add('hidden');
        
        // Mostrar detalles del usuario seleccionado
        showUserDetails(user.id);
    }

    // ==============================
    // 🔹 CARGA Y RENDERIZADO DE LISTA DE USUARIOS
    // ==============================
    async function loadUsersList() {
        try {
            showLoading('usersTableBody');
            
            // NOTA: Esta función es ficticia - en producción implementar endpoint real
            const users = await fetchAllUsers(currentPage, pageSize);
            totalUsers = users.total || users.length;
            
            renderUsersTable(users.data || users);
            updatePaginationControls();
            document.getElementById('usersCount').textContent = totalUsers;
            
        } catch (error) {
            console.error('Error cargando lista de usuarios:', error);
            showToast('Error al cargar usuarios', 'error');
            showError('usersTableBody', 'No se pudo cargar la lista de usuarios');
        }
    }

    // Función ficticia para obtener todos los usuarios
    async function fetchAllUsers(page = 1, limit = 10) {
        return new Promise((resolve) => {
            setTimeout(() => {
                // Generar usuarios de ejemplo
                const users = Array.from({ length: 50 }, (_, i) => ({
                    id: `USR${String(i + 1).padStart(3, '0')}`,
                    name: `Usuario ${i + 1}`,
                    email: `usuario${i + 1}@empresa.com`,
                    role: i === 0 ? 'Admin' : (i % 3 === 0 ? 'Reviewer' : 'User'),
                    status: i % 10 === 0 ? 'inactive' : 'active',
                    lastActivity: new Date(Date.now() - Math.random() * 30 * 24 * 60 * 60 * 1000).toISOString(),
                    joinDate: new Date(Date.now() - Math.random() * 365 * 24 * 60 * 60 * 1000).toISOString()
                }));
                
                // Paginación simulada
                const start = (page - 1) * limit;
                const end = start + limit;
                const paginatedUsers = users.slice(start, end);
                
                resolve({
                    data: paginatedUsers,
                    total: users.length,
                    page: page,
                    pageSize: limit,
                    totalPages: Math.ceil(users.length / limit)
                });
            }, 800);
        });
    }

    function renderUsersTable(users) {
        const tbody = document.getElementById('usersTableBody');
        
        if (!users || users.length === 0) {
            tbody.innerHTML = `
                <tr>
                    <td colspan="7" class="px-4 py-8 text-center text-secondary">
                        <span class="material-symbols-outlined text-4xl mb-2">search_off</span>
                        <p>No se encontraron usuarios</p>
                    </td>
                </tr>
            `;
            return;
        }
        
        tbody.innerHTML = '';
        
        users.forEach(user => {
            const statusClass = getStatusClass(user.status);
            const statusText = getStatusText(user.status);
            const lastActivity = formatRelativeTime(user.lastActivity);
            
            const tr = document.createElement('tr');
            tr.className = 'border-b border-white/5 hover:bg-white/5 transition-colors';
            tr.innerHTML = `
                <td class="px-4 py-3 text-secondary text-sm">${user.id}</td>
                <td class="px-4 py-3 font-medium">${user.name}</td>
                <td class="px-4 py-3 text-secondary">${user.email}</td>
                <td class="px-4 py-3">
                    <span class="badge ${user.role === 'Admin' ? 'badge-active' : 'badge-pending'}">${user.role}</span>
                </td>
                <td class="px-4 py-3">
                    <span class="badge ${statusClass}">${statusText}</span>
                </td>
                <td class="px-4 py-3 text-secondary text-sm">${lastActivity}</td>
                <td class="px-4 py-3">
                    <div class="flex gap-2">
                        <button class="text-blue-400 hover:text-blue-300 transition-colors view-user-btn" data-id="${user.id}" title="Ver detalles">
                            <span class="material-symbols-outlined text-sm">visibility</span>
                        </button>
                        <button class="text-green-400 hover:text-green-300 transition-colors edit-user-btn" data-id="${user.id}" title="Editar">
                            <span class="material-symbols-outlined text-sm">edit</span>
                        </button>
                        <button class="text-red-400 hover:text-red-300 transition-colors delete-user-btn" data-id="${user.id}" title="Eliminar">
                            <span class="material-symbols-outlined text-sm">delete</span>
                        </button>
                    </div>
                </td>
            `;
            tbody.appendChild(tr);
        });
        
        // Agregar event listeners
        document.querySelectorAll('.view-user-btn').forEach(btn => {
            btn.addEventListener('click', (e) => {
                const userId = e.currentTarget.dataset.id;
                showUserDetails(userId);
            });
        });
        
        document.querySelectorAll('.edit-user-btn').forEach(btn => {
            btn.addEventListener('click', (e) => {
                const userId = e.currentTarget.dataset.id;
                editUser(userId);
            });
        });
        
        document.querySelectorAll('.delete-user-btn').forEach(btn => {
            btn.addEventListener('click', (e) => {
                const userId = e.currentTarget.dataset.id;
                deleteUser(userId);
            });
        });
    }

    // ==============================
    // 🔹 PAGINACIÓN
    // ==============================
    const prevPageBtn = document.getElementById('prevPageBtn');
    const nextPageBtn = document.getElementById('nextPageBtn');
    const currentPageSpan = document.getElementById('currentPage');
    const totalPagesSpan = document.getElementById('totalPages');

    prevPageBtn.addEventListener('click', () => {
        if (currentPage > 1) {
            currentPage--;
            loadUsersList();
        }
    });

    nextPageBtn.addEventListener('click', () => {
        const totalPages = Math.ceil(totalUsers / pageSize);
        if (currentPage < totalPages) {
            currentPage++;
            loadUsersList();
        }
    });

    function updatePaginationControls() {
        const totalPages = Math.ceil(totalUsers / pageSize);
        
        currentPageSpan.textContent = currentPage;
        totalPagesSpan.textContent = totalPages;
        
        prevPageBtn.disabled = currentPage <= 1;
        nextPageBtn.disabled = currentPage >= totalPages;
        
        // Actualizar clases de estilo
        prevPageBtn.classList.toggle('opacity-50', currentPage <= 1);
        nextPageBtn.classList.toggle('opacity-50', currentPage >= totalPages);
    }

    // ==============================
    // 🔹 INVITACIÓN DE USUARIOS POR EMAIL
    // ==============================
    const inviteUserModal = document.getElementById('inviteUserModal');
    const inviteUserForm = document.getElementById('inviteUserForm');
    const sendInviteBtn = document.getElementById('sendInviteBtn');
    const inviteBtnText = document.getElementById('inviteBtnText');
    const inviteLoading = document.getElementById('inviteLoading');

    // Abrir modal para invitar usuario
    addUserBtn.addEventListener('click', () => {
        inviteUserForm.reset();
        inviteUserModal.style.display = 'block';
        document.body.style.overflow = 'hidden';
    });

    // Cerrar modales
    document.querySelectorAll('.close-modal').forEach(btn => {
        btn.addEventListener('click', () => {
            inviteUserModal.style.display = 'none';
            document.getElementById('userDetailsModal').style.display = 'none';
            document.body.style.overflow = 'auto';
        });
    });

    // Cerrar modal al hacer clic fuera
    window.addEventListener('click', (e) => {
        if (e.target === inviteUserModal) {
            inviteUserModal.style.display = 'none';
            document.body.style.overflow = 'auto';
        }
        
        const userDetailsModal = document.getElementById('userDetailsModal');
        if (e.target === userDetailsModal) {
            userDetailsModal.style.display = 'none';
            document.body.style.overflow = 'auto';
        }
    });

    // Enviar invitación por email
    inviteUserForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        
        const email = document.getElementById('inviteEmail').value.trim();
        const role = document.getElementById('inviteRole').value;
        const message = document.getElementById('inviteMessage').value.trim();
        
        if (!email || !validateEmail(email)) {
            showToast('Por favor ingresa un email válido', 'warning');
            return;
        }
        
        // Mostrar estado de carga
        inviteBtnText.classList.add('hidden');
        inviteLoading.classList.remove('hidden');
        sendInviteBtn.disabled = true;
        
        try {
            // NOTA: Esta es la estructura que mencionaste para el endpoint /inviteuser
            const inviteData = {
                idInvitado: "", // Vacío como mencionaste
                emailInvitado: email,
                // teamInvitado: "" // Opcional si tienes equipos
            };
            
            // Simular envío de invitación
            const response = await sendInvitation(inviteData);
            
            if (response.message === "OK") {
                showToast(`Invitación enviada a ${email}`, 'success');
                inviteUserModal.style.display = 'none';
                document.body.style.overflow = 'auto';
                
                // Actualizar contador de invitaciones pendientes
                updatePendingInvitesCount();
                
                // Opcional: Registrar en logs de actividad
                logActivity(`Invitación enviada a ${email}`, 'invitation');
                
            } else {
                throw new Error(response.message || 'Error al enviar invitación');
            }
            
        } catch (error) {
            console.error('Error enviando invitación:', error);
            showToast(`Error: ${error.message}`, 'error');
            
        } finally {
            // Restaurar botón
            inviteBtnText.classList.remove('hidden');
            inviteLoading.classList.add('hidden');
            sendInviteBtn.disabled = false;
        }
    });

    // Función ficticia para enviar invitación
    async function sendInvitation(inviteData) {
        // NOTA: En producción, reemplazar con fetch real
        return new Promise((resolve, reject) => {
            setTimeout(() => {
                // Simular éxito o error aleatorio para demo
                const success = Math.random() > 0.2; // 80% de éxito
                
                if (success) {
                    resolve({
                        message: "OK",
                        inviteId: `INV${Date.now()}`,
                        timestamp: new Date().toISOString()
                    });
                } else {
                    reject(new Error("El usuario ya tiene una invitación pendiente"));
                }
            }, 1500);
        });
    }

    // ==============================
    // 🔹 DETALLES DEL USUARIO
    // ==============================
    async function showUserDetails(userId) {
        const modal = document.getElementById('userDetailsModal');
        const content = document.getElementById('userDetailsContent');
        
        try {
            // Mostrar loading
            content.innerHTML = `
                <div class="text-center py-8">
                    <div class="spinner mx-auto"></div>
                    <p class="mt-2 text-secondary">Cargando detalles del usuario...</p>
                </div>
            `;
            
            modal.style.display = 'block';
            document.body.style.overflow = 'hidden';
            
            // Obtener datos del usuario
            const userData = await fetchUserDetails(userId);
            renderUserDetails(userData);
            
        } catch (error) {
            console.error('Error cargando detalles:', error);
            content.innerHTML = `
                <div class="text-center py-8 text-red-400">
                    <span class="material-symbols-outlined text-4xl mb-2">error</span>
                    <p>Error al cargar detalles del usuario</p>
                </div>
            `;
        }
    }

    // Función ficticia para obtener detalles del usuario
    async function fetchUserDetails(userId) {
        return new Promise((resolve) => {
            setTimeout(() => {
                resolve({
                    id: userId,
                    name: "Juan Pérez",
                    email: "juan.perez@empresa.com",
                    role: "Administrador",
                    status: "active",
                    joinDate: "2024-01-15T10:30:00Z",
                    lastLogin: "2024-01-20T14:45:00Z",
                    documentsSigned: 45,
                    documentsPending: 3,
                    teams: ["Equipo Directivo", "Equipo Proyecto A"],
                    permissions: ["admin", "write", "read", "delete"]
                });
            }, 800);
        });
    }

    function renderUserDetails(user) {
        const content = document.getElementById('userDetailsContent');
        const statusClass = getStatusClass(user.status);
        const statusText = getStatusText(user.status);
        
        content.innerHTML = `
            <div class="space-y-6">
                <!-- Header con avatar -->
                <div class="flex items-center gap-4">
                    <div class="w-16 h-16 rounded-full bg-gradient-to-r from-blue-500 to-purple-600 flex items-center justify-center text-2xl font-bold">
                        ${user.name.charAt(0)}
                    </div>
                    <div>
                        <h4 class="text-xl font-bold">${user.name}</h4>
                        <p class="text-secondary">${user.email}</p>
                    </div>
                    <div class="ml-auto">
                        <span class="badge ${statusClass} text-sm">${statusText}</span>
                    </div>
                </div>
                
                <!-- Información básica -->
                <div class="grid grid-cols-2 gap-4">
                    <div class="glass-card p-4">
                        <p class="text-secondary text-sm">Rol</p>
                        <p class="font-medium">${user.role}</p>
                    </div>
                    <div class="glass-card p-4">
                        <p class="text-secondary text-sm">Fecha de registro</p>
                        <p class="font-medium">${formatDate(user.joinDate)}</p>
                    </div>
                    <div class="glass-card p-4">
                        <p class="text-secondary text-sm">Último acceso</p>
                        <p class="font-medium">${formatRelativeTime(user.lastLogin)}</p>
                    </div>
                    <div class="glass-card p-4">
                        <p class="text-secondary text-sm">ID Usuario</p>
                        <p class="font-medium text-sm">${user.id}</p>
                    </div>
                </div>
                
                <!-- Estadísticas -->
                <div class="glass-card p-6">
                    <h5 class="font-bold mb-4">Estadísticas</h5>
                    <div class="grid grid-cols-3 gap-4">
                        <div class="text-center">
                            <div class="text-2xl font-bold text-primary">${user.documentsSigned}</div>
                            <p class="text-secondary text-sm">Docs firmados</p>
                        </div>
                        <div class="text-center">
                            <div class="text-2xl font-bold text-yellow-400">${user.documentsPending}</div>
                            <p class="text-secondary text-sm">Docs pendientes</p>
                        </div>
                        <div class="text-center">
                            <div class="text-2xl font-bold text-green-400">${user.teams.length}</div>
                            <p class="text-secondary text-sm">Equipos</p>
                        </div>
                    </div>
                </div>
                
                <!-- Equipos -->
                <div>
                    <h5 class="font-bold mb-2">Equipos</h5>
                    <div class="flex flex-wrap gap-2">
                        ${user.teams.map(team => `
                            <span class="badge badge-pending">${team}</span>
                        `).join('')}
                    </div>
                </div>
                
                <!-- Permisos -->
                <div>
                    <h5 class="font-bold mb-2">Permisos</h5>
                    <div class="flex flex-wrap gap-2">
                        ${user.permissions.map(perm => `
                            <span class="badge badge-active">${perm}</span>
                        `).join('')}
                    </div>
                </div>
                
                <!-- Acciones -->
                <div class="flex gap-4 pt-4 border-t border-white/10">
                    <button class="btn-secondary flex-1 h-10" onclick="resendInvitation('${user.email}')">
                        <span class="material-symbols-outlined mr-2 text-sm">send</span>
                        Reenviar invitación
                    </button>
                    <button class="btn-primary flex-1 h-10">
                        <span class="material-symbols-outlined mr-2 text-sm">edit</span>
                        Editar usuario
                    </button>
                </div>
            </div>
        `;
    }

    // ==============================
    // 🔹 FUNCIONES AUXILIARES
    // ==============================
    function getStatusClass(status) {
        switch(status) {
            case 'active': return 'badge-active';
            case 'inactive': return 'badge-inactive';
            case 'pending': return 'badge-pending';
            case 'suspended': return 'badge-suspended';
            default: return 'badge-inactive';
        }
    }

    function getStatusText(status) {
        switch(status) {
            case 'active': return 'Activo';
            case 'inactive': return 'Inactivo';
            case 'pending': return 'Pendiente';
            case 'suspended': return 'Suspendido';
            default: return status;
        }
    }

    function formatDate(dateString) {
        const date = new Date(dateString);
        return date.toLocaleDateString('es-ES', {
            year: 'numeric',
            month: 'short',
            day: 'numeric'
        });
    }

    function formatRelativeTime(dateString) {
        const date = new Date(dateString);
        const now = new Date();
        const diffMs = now - date;
        const diffMins = Math.floor(diffMs / (1000 * 60));
        const diffHours = Math.floor(diffMs / (1000 * 60 * 60));
        const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24));
        
        if (diffMins < 60) {
            return `Hace ${diffMins} min`;
        } else if (diffHours < 24) {
            return `Hace ${diffHours} horas`;
        } else if (diffDays === 1) {
            return 'Ayer';
        } else if (diffDays < 7) {
            return `Hace ${diffDays} días`;
        } else {
            return formatDate(dateString);
        }
    }

    function validateEmail(email) {
        const re = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
        return re.test(email);
    }

    function showLoading(elementId) {
        const element = document.getElementById(elementId);
        if (element) {
            element.innerHTML = `
                <tr>
                    <td colspan="7" class="px-4 py-8 text-center text-secondary">
                        <div class="spinner mx-auto"></div>
                        <p class="mt-2">Cargando...</p>
                    </td>
                </tr>
            `;
        }
    }

    function showError(elementId, message) {
        const element = document.getElementById(elementId);
        if (element) {
            element.innerHTML = `
                <tr>
                    <td colspan="7" class="px-4 py-8 text-center text-secondary">
                        <span class="material-symbols-outlined text-4xl mb-2 text-red-400">error</span>
                        <p>${message}</p>
                        <button class="btn-secondary mt-4 px-4 py-2 text-sm" onclick="loadUsersList()">
                            Reintentar
                        </button>
                    </td>
                </tr>
            `;
        }
    }

    function showToast(message, type = 'info') {
        // Crear elemento toast
        const toast = document.createElement('div');
        toast.className = `fixed top-4 right-4 z-50 px-6 py-3 rounded-lg shadow-lg transform transition-all duration-300 translate-x-full`;
        
        switch(type) {
            case 'success':
                toast.classList.add('bg-green-500/20', 'border', 'border-green-500/30', 'text-green-300');
                break;
            case 'error':
                toast.classList.add('bg-red-500/20', 'border', 'border-red-500/30', 'text-red-300');
                break;
            case 'warning':
                toast.classList.add('bg-yellow-500/20', 'border', 'border-yellow-500/30', 'text-yellow-300');
                break;
            default:
                toast.classList.add('bg-blue-500/20', 'border', 'border-blue-500/30', 'text-blue-300');
        }
        
        toast.innerHTML = `
            <div class="flex items-center gap-3">
                <span class="material-symbols-outlined">
                    ${type === 'success' ? 'check_circle' : type === 'error' ? 'error' : 'info'}
                </span>
                <span>${message}</span>
            </div>
        `;
        
        document.body.appendChild(toast);
        
        // Animación de entrada
        setTimeout(() => {
            toast.classList.remove('translate-x-full');
            toast.classList.add('translate-x-0');
        }, 10);
        
        // Eliminar después de 5 segundos
        setTimeout(() => {
            toast.classList.remove('translate-x-0');
            toast.classList.add('translate-x-full');
            setTimeout(() => {
                document.body.removeChild(toast);
            }, 300);
        }, 5000);
    }

    // ==============================
    // 🔹 FUNCIONES DE ESTADÍSTICAS DE DOCUMENTOS (Placeholder)
    // ==============================
    async function loadDocumentStats() {
        // Implementar cuando se cargue la sección de estadísticas
        console.log('Cargando estadísticas de documentos...');
    }

    function logActivity(action, type = 'info') {
        // Registrar actividad en consola (en producción, enviar a backend)
        console.log(`[${type.toUpperCase()}] ${new Date().toISOString()} - ${action}`);
    }

    function updatePendingInvitesCount() {
        // Actualizar contador de invitaciones pendientes
        const current = parseInt(document.getElementById('pendingInvites').textContent);
        document.getElementById('pendingInvites').textContent = current + 1;
    }

    // Función global para reenviar invitación
    window.resendInvitation = async function(email) {
        var _ok = await sfConfirm({title:'Reenviar invitación',message:'¿Reenviar invitación a ' + email + '?',type:'info',confirmText:'Reenviar'});if(!_ok) return;
        
        try {
            showToast(`Reenviando invitación a ${email}...`, 'info');
            
            // Lógica de reenvío similar a sendInvitation
            const inviteData = {
                idInvitado: "",
                emailInvitado: email
            };
            
            const response = await sendInvitation(inviteData);
            
            if (response.message === "OK") {
                showToast(`Invitación reenviada a ${email}`, 'success');
                logActivity(`Invitación reenviada a ${email}`, 'invitation');
            }
            
        } catch (error) {
            showToast(`Error: ${error.message}`, 'error');
        }
    };

    // Función para editar usuario (placeholder)
    window.editUser = function(userId) {
        showToast(`Editando usuario ${userId}...`, 'info');
        // Implementar lógica de edición
    };

    // Función para eliminar usuario (placeholder)
    window.deleteUser = function(userId) {
        sfConfirm({title:'Eliminar usuario',message:'¿Estás seguro de eliminar al usuario ' + userId + '? Esta acción no se puede deshacer.',type:'danger',confirmText:'Eliminar',confirmClass:'sf-modal-btn-danger'}).then(function(ok){if(ok){
            showToast(`Eliminando usuario ${userId}...`, 'warning');
            // Implementar lógica de eliminación
            setTimeout(() => {
                showToast(`Usuario ${userId} eliminado`, 'success');
                loadUsersList(); // Recargar lista
            }, 1000);
        }});
    };

    // ==============================
    // 🔹 INICIALIZACIÓN
    // ==============================
    function init() {
        // Cargar datos iniciales del dashboard
        loadDashboardData();
        
        // Configurar logout
        document.getElementById('logoutBtn').addEventListener('click', async () => {
            var ok = await sfConfirm({title:'Cerrar sesión',message:'¿Seguro que deseas cerrar tu sesión actual?',type:'warn',confirmText:'Cerrar sesión',confirmClass:'sf-modal-btn-danger'}); if(ok){
                window.location.href = '/logout';
            }
        });
        
        // Configurar tema oscuro/claro si es necesario
        const themeToggle = document.getElementById('themeToggle');
        if (themeToggle) {
            themeToggle.addEventListener('click', () => {
                document.documentElement.classList.toggle('dark');
                document.documentElement.classList.toggle('light');
                localStorage.setItem('theme', 
                    document.documentElement.classList.contains('dark') ? 'dark' : 'light'
                );
            });
        }
    }

    // Inicializar aplicación
    init();
});