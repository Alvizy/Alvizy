// smooth page transition
document.addEventListener('click', e => {
  const a = e.target.closest('a');
  if (!a) return;
  const href = a.getAttribute('href');
  if (!href || href.startsWith('http') || href.startsWith('mailto') || href.startsWith('#') || a.target === '_blank') return;

  e.preventDefault();

  if (document.startViewTransition) {
    // modern browser
    document.startViewTransition(() => {
      window.location.href = href;
    });
  } else {
    // fallback: manual fade
    document.body.style.transition = 'opacity 200ms ease';
    document.body.style.opacity = '0';
    setTimeout(() => { window.location.href = href; }, 200);
  }
});

// fade in saat halaman load
document.addEventListener('DOMContentLoaded', () => {
  document.body.style.opacity = '0';
  document.body.style.transition = 'opacity 250ms ease';
  requestAnimationFrame(() => {
    requestAnimationFrame(() => {
      document.body.style.opacity = '1';
    });
  });
});
