document.addEventListener('DOMContentLoaded', () => {
  // Efecto liquid glass para la burbuja del menú
  const liquidBubble = document.getElementById('liquidBubble');
  const navItems = document.querySelectorAll('.nav-item');
  
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
      // Aquí podrías añadir lógica para cambiar de vista
    });
  });
  
  // Calcular posiciones al cargar y al redimensionar
  window.addEventListener('load', calculatePositions);
  window.addEventListener('resize', calculatePositions);
});