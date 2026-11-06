(function () {
  var form = document.getElementById('register');
  var fields = {
    name: document.getElementById('name'),
    email: document.getElementById('email'),
    pw: document.getElementById('password'),
    confirm: document.getElementById('confirm'),
    terms: document.getElementById('terms')
  };
  var errs = {
    name: document.getElementById('name-err'),
    email: document.getElementById('email-err'),
    pw: document.getElementById('pw-err'),
    confirm: document.getElementById('confirm-err'),
    terms: document.getElementById('terms-err')
  };
  var status = document.getElementById('status');
  var submit = document.getElementById('submit');
  var toggle = document.getElementById('toggle');

  // Tombol lihat/sembunyikan kata sandi
  toggle.addEventListener('click', function () {
    var show = fields.pw.type === 'password';
    fields.pw.type = show ? 'text' : 'password';
    fields.confirm.type = show ? 'text' : 'password';
    toggle.textContent = show ? 'Sembunyikan' : 'Lihat';
    toggle.setAttribute('aria-label', show ? 'Sembunyikan kata sandi' : 'Tampilkan kata sandi');
  });

  // Mengembalikan daftar field yang bermasalah, urut dari atas
  function validate() {
    var invalid = [];
    Object.keys(errs).forEach(function (k) {
      errs[k].textContent = '';
      fields[k].removeAttribute('aria-invalid');
    });

    function fail(key, msg) {
      errs[key].textContent = msg;
      fields[key].setAttribute('aria-invalid', 'true');
      invalid.push(key);
    }

    if (fields.name.value.trim().length < 2) fail('name', 'Nama lengkap wajib diisi.');

    var mail = fields.email.value.trim();
    if (!mail) fail('email', 'Email wajib diisi.');
    else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(mail)) fail('email', 'Format email belum benar. Contoh: nama@contoh.com');

    var p = fields.pw.value;
    if (!p) fail('pw', 'Kata sandi wajib diisi.');
    else if (p.length < 8) fail('pw', 'Kata sandi minimal 8 karakter.');
    else if (!/[A-Za-z]/.test(p) || !/[0-9]/.test(p)) fail('pw', 'Gunakan kombinasi huruf dan angka.');

    if (!fields.confirm.value) fail('confirm', 'Ulangi kata sandi Anda.');
    else if (fields.confirm.value !== p) fail('confirm', 'Kata sandi tidak sama.');

    if (!fields.terms.checked) fail('terms', 'Anda harus menyetujui syarat & ketentuan.');

    return invalid;
  }

  form.addEventListener('submit', function (e) {
    e.preventDefault();
    status.textContent = '';
    status.className = 'status';

    var invalid = validate();
    if (invalid.length) {
      fields[invalid[0]].focus();
      return;
    }

    submit.disabled = true;
    submit.textContent = 'Memproses...';

    // Demo: ganti bagian ini dengan permintaan ke server Anda, misalnya:
    // fetch('/api/register', { method: 'POST', headers: {'Content-Type': 'application/json'},
    //   body: JSON.stringify({ name: fields.name.value, email: fields.email.value, password: fields.pw.value }) })
    // setTimeout(function () {
    //   submit.disabled = false;
    //   submit.textContent = 'Daftar';
    //   status.className = 'status ok';
    //   status.textContent = 'Akun berhasil dibuat (demo). Hubungkan form ini ke server Anda.';
    //   form.reset();
    // }, 900);

    fetch('http://localhost:8080/api/auth/register', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    name: fields.name.value.trim(),
    email: fields.email.value.trim(),
    password: fields.pw.value
  })
})
.then(async res => {
  const body = await res.json();
  if (!res.ok) throw new Error(body.error || 'Registrasi gagal');
  status.className = 'status ok';
  status.textContent = 'Akun berhasil dibuat. Silakan masuk.';
  form.reset();
})
.catch(err => { status.textContent = err.message; })
.finally(() => { submit.disabled = false; submit.textContent = 'Daftar'; });
  });
})();