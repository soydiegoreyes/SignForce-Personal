document.addEventListener('DOMContentLoaded', () => {
  // Efecto liquid glass para la burbuja del menú
  const liquidBubble = document.getElementById('liquidBubble');
  
  // Verificar si hay elementos de navegación para el efecto burbuja
  const navItems = document.querySelectorAll('.nav-item');
  
  if (navItems.length > 0) {
    // Posiciones iniciales de los items del menú
    const itemPositions = {};
    
    // Calcular posiciones de los items
    function calculatePositions() {
      navItems.forEach(item => {
        const rect = item.getBoundingClientRect();
        const headerRect = document.getElementById('mainHeader').getBoundingClientRect();
        
        itemPositions[item.dataset.item] = {
          left: rect.left - headerRect.left,
          width: rect.width
        };
      });
      
      // Posicionar burbuja en el item activo (home por defecto)
      moveBubble('home');
    }
    
    // Mover la burbuja al item seleccionado
    function moveBubble(itemName) {
      const item = itemPositions[itemName];
      if (item) {
        liquidBubble.style.left = `${item.left - 10}px`;
        liquidBubble.style.width = `${item.width + 20}px`;
      }
    }
    
    // Event listeners para los items del menú
    navItems.forEach(item => {
      item.addEventListener('mouseenter', () => {
        moveBubble(item.dataset.item);
      });
      
      item.addEventListener('click', (e) => {
        e.preventDefault();
        moveBubble(item.dataset.item);
      });
    });
    
    // Calcular posiciones al cargar y al redimensionar
    window.addEventListener('load', calculatePositions);
    window.addEventListener('resize', calculatePositions);
  }

  // Animación de las tarjetas al cargar
  const actionCards = document.querySelectorAll('.action-card');
  actionCards.forEach((card, index) => {
    card.style.opacity = '0';
    card.style.transform = 'translateY(20px)';
    card.style.animation = `fadeInUp 0.5s ease-out ${index * 0.1}s forwards`;
  });

  // Añadir estilos de animación dinámicamente
  const style = document.createElement('style');
  style.textContent = `
    @keyframes fadeInUp {
      from {
        opacity: 0;
        transform: translateY(20px);
      }
      to {
        opacity: 1;
        transform: translateY(0);
      }
    }
  `;
  document.head.appendChild(style);

  let equipos = JSON.parse(sessionStorage.getItem('equipos')) || [];
  // Cargar información del usuario
  const currentUser = JSON.parse(sessionStorage.getItem('currentUser'));
  
  if (!currentUser || currentUser.tipo !== 'root') {
      window.location.href = './../registro/login.html';
  }
  
  document.getElementById('nombreUsuario').textContent = currentUser.nombre || 'Usuario ROOT';
  document.getElementById('welcomeUserName').textContent = currentUser.nombre || 'Usuario ROOT';

  const equiposEmpresa = equipos.filter(e => e.empresaId === currentUser.empresaId);
    
  if (equiposEmpresa.length > 0) {
      let html = '';
      equiposEmpresa.forEach(equipo => {
          html += `
              <div class="equipo-item">
                  <div>
                      <h3>${equipo.nombre}</h3>
                      <p>Usuario maestro: ${equipo.usuarioMaestroEmail}</p>
                  </div>
                  <div>
                      <span>Miembros: ${equipo.miembros ? equipo.miembros.length : 0}</span>
                  </div>
              </div>
          `;
      });
      document.getElementById('listaEquipos').innerHTML = html;
  }

});


async function logout() {
    sessionStorage.removeItem('currentUser');
    window.location.href = './../registro/login.html';
}