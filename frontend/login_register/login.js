(function () {
  var form = document.getElementById('login');
  var email = document.getElementById('email');
  var pw = document.getElementById('password');
  var emailErr = document.getElementById('email-err');
  var pwErr = document.getElementById('pw-err');
  var status = document.getElementById('status');
  var submit = document.getElementById('submit');
  var toggle = document.getElementById('toggle');
  var remember = document.getElementById('remember');

  const API = 'http://localhost:8080/api';

  // -------Isi ulang email jika pernah dicentang "Ingat saya"---------------
  try {
    var saved = localStorage.getItem('login-email');
    if (saved) { email.value = saved; remember.checked = true; }
  } catch (e) {}

  toggle.addEventListener('click', function () {
    var show = pw.type === 'password';
    pw.type = show ? 'text' : 'password';
    toggle.textContent = show ? 'Sembunyikan' : 'Lihat';
    toggle.setAttribute('aria-label', show ? 'Sembunyikan kata sandi' : 'Tampilkan kata sandi');
  });

  function validate() {
    var ok = true;
    emailErr.textContent = ''; pwErr.textContent = '';
    email.removeAttribute('aria-invalid'); pw.removeAttribute('aria-invalid');

    var v = email.value.trim();
    if (!v) { emailErr.textContent = 'Email wajib diisi.'; ok = false; }
    else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(v)) {
      emailErr.textContent = 'Format email belum benar. Contoh: nama@contoh.com'; ok = false;
    }
    if (!ok) email.setAttribute('aria-invalid', 'true');

    if (!pw.value) { pwErr.textContent = 'Kata sandi wajib diisi.'; pw.setAttribute('aria-invalid', 'true'); ok = false; }
    else if (pw.value.length < 8) { pwErr.textContent = 'Kata sandi minimal 8 karakter.'; pw.setAttribute('aria-invalid', 'true'); ok = false; }
    return ok;
  }

  form.addEventListener('submit', async function (e) {
    e.preventDefault();
    status.textContent = ''; status.className = 'status';

    if (!validate()) {
      (email.getAttribute('aria-invalid') ? email : pw).focus();
      return;
    }

    try {
      if (remember.checked) localStorage.setItem('login-email', email.value.trim());
      else localStorage.removeItem('login-email');
    } catch (err) {}

    submit.disabled = true; submit.textContent = 'Memproses...';

    try {
      const res = await fetch(API + '/auth/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email: email.value.trim(), password: pw.value })
      });
      const body = await res.json();
      if (!res.ok) throw new Error(body.error || 'Login gagal');

      // -------------Simpan token + user untuk request berikutnya-----------------------
      localStorage.setItem('token', body.token);
      localStorage.setItem('user', JSON.stringify(body.user));

      status.className = 'status ok';
      status.textContent = 'Berhasil masuk. Mengalihkan...';

      // --------------------Redirect sesuai role----------------------------------------
      setTimeout(function () {
        if (body.user.role === 'admin') {
          location.href = '../homepage_admin/homepage_admin.html';
        } else {
          location.href = '../homepage_users/homepage_users.html';
        }
      }, 400);
    } catch (err) {
      status.className = 'status';
      status.textContent = err.message;
      submit.disabled = false;
      submit.textContent = 'Masuk';
    }
  });

  document.getElementById('forgot').addEventListener('click', function (e) {
    e.preventDefault();
    status.textContent = 'Fitur lupa kata sandi belum tersedia.';
  });
})();