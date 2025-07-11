document.getElementById('fileInput').addEventListener('change', function (event) {
  let fileURL = null;
  const previewContainer = document.getElementById('preview-container');

  const file = event.target.files[0];
  if (file && file.type === 'application/pdf') {
    document.getElementById('docName').value= file.name;
    const reader = new FileReader();
    reader.onload = function (e) {
      if (fileURL) URL.revokeObjectURL(fileURL);
      fileURL = URL.createObjectURL(file);

      previewContainer.innerHTML = '';

      // Wrapper relativo
      const wrapper = document.createElement('div');
      wrapper.style.position = 'relative';
      wrapper.style.width = '100%';
      wrapper.style.height = '100%';

      // Iframe con eventos habilitados
      const iframe = document.createElement('iframe');
      iframe.src = `${fileURL}#page=1&zoom=25%`;
      iframe.style.width = '100%';
      iframe.style.height = '100%';
      iframe.style.border = 'none';

      // Div flotante para click
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