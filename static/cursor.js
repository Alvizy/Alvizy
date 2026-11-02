// cursor smoke — nyebar lebar, fluffy
(function() {
  const canvas = document.createElement('canvas');
  canvas.style.cssText = 'position:fixed;inset:0;pointer-events:none;z-index:9999;';
  document.body.appendChild(canvas);
  const ctx = canvas.getContext('2d');

  function resize() {
    canvas.width  = window.innerWidth;
    canvas.height = window.innerHeight;
  }
  resize();
  window.addEventListener('resize', resize);

  let mouse = { x: -999, y: -999 };
  let prevMouse = { x: -999, y: -999 };
  window.addEventListener('mousemove', e => {
    prevMouse = { ...mouse };
    mouse.x = e.clientX;
    mouse.y = e.clientY;
  });

  // particle system
  const particles = [];
  let hue = 180;

  function spawnParticles() {
    const dx = mouse.x - prevMouse.x;
    const dy = mouse.y - prevMouse.y;
    const speed = Math.sqrt(dx*dx + dy*dy);

    // jangan spawn kalau mouse diam
    if (speed < 2) return;

    const count = Math.min(Math.floor(speed / 4) + 1, 5);

    for (let i = 0; i < count; i++) {
      particles.push({
        x:     mouse.x + (Math.random() - 0.5) * 8,
        y:     mouse.y + (Math.random() - 0.5) * 8,
        vx:    (Math.random() - 0.5) * 1.2,
        vy:    (Math.random() - 0.5) * 1.2 - 0.4, // sedikit naik
        r:     Math.random() * 18 + 10,            // radius besar
        alpha: Math.random() * 0.4 + 0.25,
        hue:   hue + (Math.random() - 0.5) * 40,
        decay: Math.random() * 0.012 + 0.008,
        grow:  Math.random() * 0.3 + 0.1,          // particle mengembang
      });
    }
    hue = (hue + 1.2) % 360;
  }

  function draw() {
    ctx.clearRect(0, 0, canvas.width, canvas.height);

    spawnParticles();

    for (let i = particles.length - 1; i >= 0; i--) {
      const p = particles[i];

      p.x     += p.vx;
      p.y     += p.vy;
      p.r     += p.grow;        // mengembang
      p.alpha -= p.decay;
      p.vx    *= 0.98;
      p.vy    *= 0.98;

      if (p.alpha <= 0) { particles.splice(i, 1); continue; }

      // soft blob
      const g = ctx.createRadialGradient(p.x, p.y, 0, p.x, p.y, p.r);
      g.addColorStop(0,   `hsla(${p.hue}, 100%, 70%, ${p.alpha})`);
      g.addColorStop(0.4, `hsla(${p.hue}, 100%, 60%, ${p.alpha * 0.5})`);
      g.addColorStop(1,   `hsla(${p.hue}, 100%, 50%, 0)`);

      ctx.beginPath();
      ctx.arc(p.x, p.y, p.r, 0, Math.PI * 2);
      ctx.fillStyle = g;
      ctx.fill();
    }

    requestAnimationFrame(draw);
  }

  requestAnimationFrame(draw);
})();
