// Cargar información del usuario
const currentUser = JSON.parse(sessionStorage.getItem('currentUser'));
        // Cargar información de los equipos
        let equipos = JSON.parse(sessionStorage.getItem('equipos')) || [];
        // Cargar información de los usuarios
        let usuarios= JSON.parse(sessionStorage.getItem('usuarios')) || [];


        if (!currentUser || currentUser.tipo !== 'admin_equipo') {
            window.location.href = './../registro/login.html';
        }
        
        document.getElementById('nombreUsuario').textContent = currentUser.nombre || 'Administrador de Equipo';
        
        // Cargar información del equipo y miembros
        function cargarEquipoYMiembros() {
            // Simular carga de equipo
            // En un caso real sería un fetch a la API
            /*
            fetch(`http://localhost:8001/userdata'`)
            .then(response => response.json())
            .then(equipo => {
                // Mostrar información del equipo
            });
            */
            // Buscar equipo en memoria
            const equipo = equipos.find(e => e.id === currentUser.equipoId);
            
            if (equipo) {
                document.getElementById('infoEquipo').innerHTML = `
                    <h2>${equipo.nombre}</h2>
                    <p>${equipo.descripcion || 'Sin descripción'}</p>
                    <p><strong>Creado el:</strong> ${new Date(equipo.fechaCreacion).toLocaleDateString()}</p>
                `;
                
                document.getElementById('tituloBienvenida').textContent = `Bienvenido, ${currentUser.nombre || 'Administrador'} - ${equipo.nombre}`;
                
                // Cargar miembros del equipo
                cargarMiembros(equipo);
            } else {
                document.getElementById('infoEquipo').innerHTML = '<p>No se encontró información del equipo</p>';
            }
        }
        
        function cargarMiembros(equipo) {
            // Simular carga de miembros
            // En un caso real sería un fetch a la API
            /*
            fetch(`/api/equipos/${equipo.id}/miembros`)
            .then(response => response.json())
            .then(miembros => {
                // Mostrar miembros
            });
            */
            
            if (equipo.miembros && equipo.miembros.length > 0 && usuarios.length > 0) {
                let html = '';
                
                equipo.miembros.forEach(miembroEquipo => {
                    // Buscar información del usuario
                    const usuario = usuarios.find(u => u.id === miembroEquipo.usuarioId);
                    
                    if (usuario) {
                        const rolClass = miembroEquipo.rol === 'admin' ? 'rol-admin' : 
                                        miembroEquipo.rol === 'miembro' ? 'rol-miembro' : '';
                        
                        html += `
                            <div class="miembro-item">
                                <div>
                                    <h3>${usuario.nombre}</h3>
                                    <p>${usuario.email}</p>
                                </div>
                                <div>
                                    <span class="miembro-rol ${rolClass}">${miembroEquipo.rol.toUpperCase()}</span>
                                </div>
                            </div>
                        `;
                    }
                });
                
                document.getElementById('listaMiembros').innerHTML = html;
            } else {
                document.getElementById('listaMiembros').innerHTML = '<p>No hay miembros en este equipo aún.</p>';
            }
        }
        
        cargarEquipoYMiembros();
        
        function logout() {
            sessionStorage.removeItem('currentUser');
            window.location.href = './../registro/login_register.html';
        }
