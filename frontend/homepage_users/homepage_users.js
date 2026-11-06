(function () {
  const API = 'http://localhost:8080/api';
  const CART_KEY = 'user-cart';

  // ---------- Guard: harus login ----------
  const token = localStorage.getItem('token');
  const currentUser = JSON.parse(localStorage.getItem('user') || 'null');
  if (!token || !currentUser) {
    alert('Silakan masuk terlebih dahulu.');
    location.href = '../login_register/login.html';
    return;
  }

  // ---------- Guard: logout ----------
  document.querySelector('.logout')?.addEventListener('click', function (e) {
    e.preventDefault();
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    location.href = '../login_register/login.html';
  });

  // ---------- State ----------
  let data = { products: [] };
  let cart = {}; // { product_id: qty }

  // State Pagination & Search
  let currentPage = 1;
  const itemsPerPage = 8; // Tampilkan 8 produk per halaman
  let searchQuery = "";
  let selectedCategory = "";
  
  try {
    const c = localStorage.getItem(CART_KEY);
    cart = c ? JSON.parse(c) : {};
  } catch (e) { cart = {}; }

  const rupiah = new Intl.NumberFormat('id-ID', {
    style: 'currency', currency: 'IDR', maximumFractionDigits: 0
  });

  function $(id) { return document.getElementById(id); }
  function el(tag, cls, text) {
    const n = document.createElement(tag);
    if (cls) n.className = cls;
    if (text !== undefined) n.textContent = text;
    return n;
  }
  function find(id) { return data.products.find(p => p.id === Number(id)); }
  function saveCart() { try { localStorage.setItem(CART_KEY, JSON.stringify(cart)); } catch (e) {} }

  function authHeaders() {
    return {
      'Content-Type': 'application/json',
      'Authorization': 'Bearer ' + localStorage.getItem('token')
    };
  }

  // ---------- API ----------
  async function apiGet(path) {
    const res = await fetch(API + path);
    const body = await res.json().catch(() => ({}));
    if (!res.ok) throw new Error(body.error || 'Gagal memuat data');
    return body;
  }

  async function loadCategoriesForUser() {
    try {
      const res = await apiGet('/categories');
      const select = $('filter-category');
      // Hapus opsi lama (kecuali opsi pertama)
      while (select.options.length > 1) { select.remove(1); }
      
      res.data.forEach(c => {
        const opt = document.createElement('option');
        opt.value = c.name; // Kirim nama kategori sesuai API
        opt.textContent = c.name;
        select.appendChild(opt);
      });
    } catch (err) {
      console.error("Gagal load kategori", err);
    }
  }

  async function loadProducts() {
    let url = `/products?page=${currentPage}&limit=${itemsPerPage}&in_stock=true`;
    if (searchQuery) url += `&q=${encodeURIComponent(searchQuery)}`;
    if (selectedCategory) url += `&category=${encodeURIComponent(selectedCategory)}`;

    const res = await apiGet(url);
    const meta = res.meta;

    // Update Info Pagination
    $('user-page-info').textContent = `Halaman ${meta.current_page} dari ${meta.total_page || 1}`;
    $('user-btn-prev').disabled = meta.current_page <= 1;
    $('user-btn-next').disabled = meta.current_page >= meta.total_page;

    data.products = res.data;
  }

  // ---------- Cart helpers ----------
  function sanitizeCart() {
    Object.keys(cart).forEach(id => {
      const p = find(id);
      if (!p || p.stock < 1) { delete cart[id]; return; }
      cart[id] = Math.max(1, Math.min(cart[id], p.stock));
    });
  }
  function cartLines() {
    return Object.keys(cart).map(id => {
      const p = find(id);
      return { p, qty: cart[id], sub: p.price * cart[id] };
    });
  }
  function cartTotal() { return cartLines().reduce((s, l) => s + l.sub, 0); }

  function addToCart(id) {
    const p = find(id);
    if (!p || p.stock < 1) return;
    cart[id] = Math.min((cart[id] || 0) + 1, p.stock);
    saveCart(); render();
  }
  function setQty(id, qty) {
    const p = find(id);
    if (!p) return;
    qty = Math.floor(Number(qty));
    if (!qty || qty < 1) qty = 1;
    cart[id] = Math.min(qty, p.stock);
    saveCart(); render();
  }
  function removeItem(id) { delete cart[id]; saveCart(); render(); }

  // ---------- Render ----------
  function renderProducts() {
    const box = $('products');
    box.textContent = '';
    if (!data.products.length) {
      box.appendChild(el('p', 'empty', 'Tidak ada produk yang cocok dengan pencarian/filter.'));
      return;
    }
    data.products.forEach(p => {
      const card = el('article', 'product');
      card.appendChild(el('h3', '', p.name));
      card.appendChild(el('p', 'desc', p.desc || 'Tidak ada deskripsi.'));
      card.appendChild(el('div', 'price', rupiah.format(p.price)));
      card.appendChild(el('div', 'stock', p.stock > 0 ? 'Stok: ' + p.stock : 'Stok habis'));

      const btns = el('div', 'btns');
      const add = el('button', 'btn ghost', 'Tambah ke keranjang');
      add.type = 'button';
      const buy = el('button', 'btn', 'Beli sekarang');
      buy.type = 'button';
      if (p.stock < 1) { add.disabled = true; buy.disabled = true; }
      add.addEventListener('click', () => addToCart(p.id));
      buy.addEventListener('click', () => { addToCart(p.id); openCheckout(); });
      btns.appendChild(add); btns.appendChild(buy);
      card.appendChild(btns);
      box.appendChild(card);
    });
  }

  function renderCart() {
    const box = $('cart-items');
    box.textContent = '';
    const lines = cartLines();
    if (!lines.length) box.appendChild(el('p', 'empty', 'Keranjang masih kosong.'));

    lines.forEach(l => {
      const item = el('div', 'item');
      item.appendChild(el('div', 'name', l.p.name));
      item.appendChild(el('div', 'sub', rupiah.format(l.p.price) + ' / item · stok ' + l.p.stock));

      const row = el('div', 'item-row');
      const qty = el('div', 'qty');
      const minus = el('button', '', '−'); minus.type = 'button';
      minus.setAttribute('aria-label', 'Kurangi jumlah ' + l.p.name);
      const input = el('input'); input.type = 'number';
      input.min = 1; input.max = l.p.stock; input.value = l.qty;
      input.setAttribute('aria-label', 'Jumlah ' + l.p.name);
      const plus = el('button', '', '+'); plus.type = 'button';
      plus.setAttribute('aria-label', 'Tambah jumlah ' + l.p.name);
      minus.addEventListener('click', () => setQty(l.p.id, l.qty - 1));
      plus.addEventListener('click', () => setQty(l.p.id, l.qty + 1));
      input.addEventListener('change', () => setQty(l.p.id, input.value));
      qty.appendChild(minus); qty.appendChild(input); qty.appendChild(plus);

      row.appendChild(qty);
      row.appendChild(el('span', 'line-total', rupiah.format(l.sub)));
      item.appendChild(row);

      const rm = el('button', 'remove', 'Hapus dari keranjang');
      rm.type = 'button';
      rm.addEventListener('click', () => removeItem(l.p.id));
      item.appendChild(rm);
      box.appendChild(item);
    });

    $('cart-total').textContent = rupiah.format(cartTotal());
    $('checkout-btn').disabled = !lines.length;
  }

  function render() { sanitizeCart(); renderProducts(); renderCart(); }

  // ---------- Checkout ----------
  const dlg = $('checkout');

  function openCheckout() {
    const lines = cartLines();
    if (!lines.length) return;
    const ul = $('co-summary');
    ul.textContent = '';
    lines.forEach(l => {
      const li = el('li');
      li.appendChild(el('span', '', l.p.name + ' × ' + l.qty));
      li.appendChild(el('span', '', rupiah.format(l.sub)));
      ul.appendChild(li);
    });
    $('co-total').textContent = rupiah.format(cartTotal());
    $('pay-err').textContent = '';
    $('co-form-view').hidden = false;
    $('receipt-view').hidden = true;
    if (!dlg.open) dlg.showModal();
    $('buyer').focus();
  }

  $('checkout-btn').addEventListener('click', openCheckout);
  $('co-cancel').addEventListener('click', () => dlg.close());

  $('pay-form').addEventListener('submit', async function (e) {
    e.preventDefault();
    const buyer = $('buyer'), method = $('method'), err = $('pay-err');
    buyer.removeAttribute('aria-invalid');
    method.removeAttribute('aria-invalid');
    err.textContent = '';

    if (!buyer.value.trim()) {
      err.textContent = 'Nama pembeli wajib diisi.';
      buyer.setAttribute('aria-invalid', 'true'); buyer.focus(); return;
    }
    if (!method.value) {
      err.textContent = 'Pilih metode pembayaran.';
      method.setAttribute('aria-invalid', 'true'); method.focus(); return;
    }

    sanitizeCart();
    const lines = cartLines();
    if (!lines.length) { err.textContent = 'Keranjang kosong atau stok sudah habis.'; return; }

    const items = Object.keys(cart).map(id => ({ product_id: Number(id), qty: cart[id] }));

    const btn = $('pay-btn');
    btn.disabled = true; btn.textContent = 'Memproses...';

    try {
      const res = await fetch(API + '/checkout', {
        method: 'POST',
        headers: authHeaders(),
        body: JSON.stringify({
          buyer: buyer.value.trim(),
          method: method.value,
          items
        })
      });
      const body = await res.json();
      if (!res.ok) throw new Error(body.error || 'Pembayaran gagal');

      const s = body.data;
      showReceipt({
        no: s.invoice_no,
        date: new Date(s.date),
        buyer: s.buyer,
        method: s.method,
        total: s.total,
        lines: s.items.map(it => ({
          p: { name: it.name, price: it.price },
          qty: it.qty,
          sub: it.subtotal
        }))
      });

      cart = {}; saveCart();
      
      // Setelah checkout, kembali ke halaman 1 dan refresh data
      currentPage = 1; 
      await loadProducts();
      render();
      $('pay-form').reset();
    } catch (ex) {
      err.textContent = ex.message;
    } finally {
      btn.disabled = false; btn.textContent = 'Bayar';
    }
  });

  function showReceipt(r) {
    $('r-no').textContent = r.no;
    $('r-date').textContent = r.date.toLocaleString('id-ID', {
      dateStyle: 'long', timeStyle: 'short'
    });
    $('r-buyer').textContent = r.buyer;
    $('r-method').textContent = r.method;
    const body = $('r-items');
    body.textContent = '';
    r.lines.forEach(l => {
      const tr = el('tr');
      tr.appendChild(el('td', '', l.p.name));
      tr.appendChild(el('td', 'num', String(l.qty)));
      tr.appendChild(el('td', 'num', rupiah.format(l.sub)));
      body.appendChild(tr);
    });
    $('r-total').textContent = rupiah.format(r.total);
    $('co-form-view').hidden = true;
    $('receipt-view').hidden = false;
  }

  $('r-print').addEventListener('click', () => window.print());
  $('r-close').addEventListener('click', () => dlg.close());

  // ---------- Event Pencarian & Filter Kategori & Pagination ----------
  let searchTimeout;
  $('search-user').addEventListener('input', (e) => {
    clearTimeout(searchTimeout);
    searchTimeout = setTimeout(async () => {
      searchQuery = e.target.value.trim();
      currentPage = 1;
      await loadProducts();
      render();
    }, 500);
  });

  $('filter-category').addEventListener('change', async (e) => {
    selectedCategory = e.target.value;
    currentPage = 1;
    await loadProducts();
    render();
  });

  $('user-btn-prev').addEventListener('click', async () => {
    if (currentPage > 1) {
      currentPage--;
      await loadProducts();
      render();
    }
  });

  $('user-btn-next').addEventListener('click', async () => {
    currentPage++;
    await loadProducts();
    render();
  });

    // ---------- Ubah Kata Sandi ----------
  const pwDlg = $('change-pw-dialog');
  const pwForm = $('change-pw-form');
  const btnChangePw = $('btn-change-pw');

  btnChangePw.addEventListener('click', (e) => {
    e.preventDefault();
    $('old-pw').value = '';
    $('new-pw').value = '';
    $('pw-err').textContent = '';
    if (!pwDlg.open) pwDlg.showModal();
    $('old-pw').focus();
  });

  $('pw-cancel').addEventListener('click', () => pwDlg.close());

  pwForm.addEventListener('submit', async (e) => {
    e.preventDefault();
    const oldPw = $('old-pw').value;
    const newPw = $('new-pw').value;
    const err = $('pw-err');
    const btn = $('pw-submit');

    err.textContent = '';

    // Validasi Frontend
    if (!oldPw) { err.textContent = 'Kata sandi lama wajib diisi.'; return; }
    if (!newPw || newPw.length < 8) { err.textContent = 'Kata sandi baru minimal 8 karakter.'; return; }

    btn.disabled = true; btn.textContent = 'Menyimpan...';

    try {
      // currentUser.id didapat dari localStorage saat inisialisasi awal
      const res = await fetch(API + '/users/' + currentUser.id + '/password', {
        method: 'PUT',
        headers: authHeaders(),
        body: JSON.stringify({
          old_password: oldPw,
          new_password: newPw
        })
      });

      const body = await res.json();
      if (!res.ok) throw new Error(body.error || 'Gagal mengubah kata sandi');

      alert('Kata sandi berhasil diubah. Silakan login kembali dengan sandi baru.');
      pwDlg.close();
      
      // Paksa logout setelah ganti sandi
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      location.href = '../login_register/login.html';

    } catch (ex) {
      err.textContent = ex.message;
    } finally {
      btn.disabled = false; btn.textContent = 'Simpan';
    }
  });

  // ---------- Start ----------
  (async function start() {
    try {
    await loadCategoriesForUser(); // Load dropdown filter
    await loadProducts();
    render();
    } catch (err) {
      $('products').appendChild(el('p', 'empty', 'Gagal memuat produk: ' + err.message));
    }
  })();
})(); 