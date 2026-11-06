(function () {
  const API = 'http://localhost:8080/api';

  // ---------- Guard: harus admin ----------
  const token = localStorage.getItem('token');
  const me = JSON.parse(localStorage.getItem('user') || 'null');
  if (!token || !me || me.role !== 'admin') {
    alert('Anda harus masuk sebagai admin.');
    location.href = '../login_register/login.html';
    return;
  }

  // ---------- State Pagination & Search ----------
  let currentPage = 1;
  const itemsPerPage = 10;
  let searchQuery = "";
  let categoriesData = [];

  // ---------- Helpers ----------
  const rupiah = new Intl.NumberFormat('id-ID', {
    style: 'currency', currency: 'IDR', maximumFractionDigits: 0
  });

  function $(id) { return document.getElementById(id); }
  
  function tdText(text) {
    const td = document.createElement('td');
    td.textContent = text;
    return td;
  }

  function authHeaders() {
    return {
      'Content-Type': 'application/json',
      'Authorization': 'Bearer ' + localStorage.getItem('token')
    };
  }

  function fmtDate(iso) {
    const d = new Date(iso);
    return d.toLocaleDateString('id-ID', { day: 'numeric', month: 'short', year: 'numeric' });
  }

  async function apiGet(path) {
    const res = await fetch(API + path, { headers: authHeaders() });
    const body = await res.json().catch(() => ({}));
    if (!res.ok) throw new Error(body.error || 'Gagal memuat');
    return body;
  }

  function row(cells) {
    const tr = document.createElement('tr');
    cells.forEach(c => {
      const td = document.createElement('td');
      if (c.num) td.className = 'num';
      if (c.badge) {
        const b = document.createElement('span');
        b.className = 'badge';
        b.textContent = c.text;
        td.appendChild(b);
      } else {
        td.textContent = c.text;
      }
      tr.appendChild(td);
    });
    return tr;
  }

  function fill(id, rows, cols) {
    const body = $(id);
    body.textContent = '';
    if (!rows.length) {
      const tr = document.createElement('tr');
      const td = document.createElement('td');
      td.className = 'empty';
      td.colSpan = cols;
      td.textContent = 'Belum ada data.';
      tr.appendChild(td);
      body.appendChild(tr);
      return;
    }
    rows.forEach(r => body.appendChild(r));
  }

  // ---------- Data loaders ----------
  async function loadStats() {
    const s = await apiGet('/stats');
    $('st-products').textContent = s.total_products;
    $('st-users').textContent = s.total_users;
    $('st-sold').textContent = s.total_sold;
    $('st-revenue').textContent = rupiah.format(s.total_revenue);
  }

  async function loadProducts() {
    let url = `/products?page=${currentPage}&limit=${itemsPerPage}`;
    if (searchQuery) url += `&q=${encodeURIComponent(searchQuery)}`;
    
    const res = await apiGet(url);
    const meta = res.meta;

    // Update Info Pagination
    $('page-info').textContent = `Halaman ${meta.current_page} dari ${meta.total_page || 1}`;
    $('btn-prev').disabled = meta.current_page <= 1;
    $('btn-next').disabled = meta.current_page >= meta.total_page;

    fill('products-body', res.data.map(p => {
      const tr = document.createElement('tr');

      // Kolom Kategori
      tr.appendChild(tdText(p.category ? p.category.name : '-'));
      tr.appendChild(tdText(p.name));

      const tdPrice = tdText(rupiah.format(p.price));
      tdPrice.classList.add('num');
      tr.appendChild(tdPrice);

      const tdStock = tdText(String(p.stock));
      tdStock.classList.add('num');
      tr.appendChild(tdStock);

      // Kolom Aksi
      const tdAksi = document.createElement('td');
      tdAksi.className = 'aksi';

      const btnEdit = document.createElement('button');
      btnEdit.type = 'button';
      btnEdit.className = 'btn-mini';
      btnEdit.textContent = 'Edit';
      btnEdit.addEventListener('click', () => openEditProduct(p));

      const btnHapus = document.createElement('button');
      btnHapus.type = 'button';
      btnHapus.className = 'btn-mini danger';
      btnHapus.textContent = 'Hapus';
      btnHapus.addEventListener('click', () => deleteProduct(p));

      tdAksi.appendChild(btnEdit);
      tdAksi.appendChild(btnHapus);
      tr.appendChild(tdAksi);

      return tr;
    }), 5);
  }

  async function loadCategories() {
    const res = await apiGet('/categories');
    categoriesData = res.data;
    
    const pSelect = $('p-category');
    const eSelect = $('e-category');
    const currentEditVal = eSelect.value;

    pSelect.innerHTML = '<option value="">Pilih Kategori</option>';
    eSelect.innerHTML = '';

    categoriesData.forEach(c => {
      const opt1 = document.createElement('option');
      opt1.value = c.id; opt1.textContent = c.name;
      pSelect.appendChild(opt1);

      const opt2 = document.createElement('option');
      opt2.value = c.id; opt2.textContent = c.name;
      eSelect.appendChild(opt2);
    });

    if (currentEditVal) eSelect.value = currentEditVal;

    fill('categories-body', categoriesData.map(c => {
      const tr = document.createElement('tr');
      tr.appendChild(tdText(c.name));

      const tdAksi = document.createElement('td');
      tdAksi.className = 'aksi';
      
      const btnHapus = document.createElement('button');
      btnHapus.type = 'button';
      btnHapus.className = 'btn-mini danger';
      btnHapus.textContent = 'Hapus';
      btnHapus.addEventListener('click', () => deleteCategory(c));
      
      tdAksi.appendChild(btnHapus);
      tr.appendChild(tdAksi);
      return tr;
    }), 2);
  }

  async function loadUsers() {
    const res = await apiGet('/users');
    const all = res.data;
    const admins = all.filter(u => u.role === 'admin');
    const customers = all.filter(u => u.role === 'customer');

    fill('admins-body', admins.map(u => row([
      { text: u.name },
      { text: u.email },
      { text: fmtDate(u.joined) },
      { text: u.role, badge: true }
    ])), 4);

    fill('users-body', customers.map(u => row([
      { text: u.name },
      { text: u.email },
      { text: fmtDate(u.joined) },
      { text: u.role, badge: true }
    ])), 4);
  }

  async function loadSales() {
    const [summary, list] = await Promise.all([
      apiGet('/sales/summary'),
      apiGet('/sales')
    ]);

    fill('sold-body', summary.data.map(s => row([
      { text: s.product },
      { text: String(s.qty), num: true },
      { text: rupiah.format(s.revenue), num: true }
    ])), 3);

    fill('sales-body', list.data.map(s => row([
      { text: fmtDate(s.date) },
      { text: s.buyer },
      { text: s.items.map(i => i.name).join(', ') },
      { text: String(s.items.reduce((a, i) => a + i.qty, 0)), num: true },
      { text: rupiah.format(s.total), num: true }
    ])), 5);
  }

  // ---------- Edit produk ----------
  const editDlg = $('edit-product');
  const editForm = $('edit-form');
  const editFields = {
    name: $('e-name'),
    category: $('e-category'),
    desc: $('e-desc'),
    price: $('e-price'), 
    stock: $('e-stock')
  };
  const editErr = $('e-err');
  let editingId = null;

  function openEditProduct(p) {
    editingId = p.id;
    editFields.name.value = p.name;
    editFields.category.value = p.category_id;
    editFields.desc.value = p.desc || '';
    editFields.price.value = p.price;
    editFields.stock.value = p.stock;
    editErr.textContent = '';
    ['name','price','stock'].forEach(k => editFields[k].removeAttribute('aria-invalid'));
    if (!editDlg.open) editDlg.showModal();
    editFields.name.focus();
  }

  $('e-cancel').addEventListener('click', () => editDlg.close());

  editForm.addEventListener('submit', async function (e) {
    e.preventDefault();
    editErr.textContent = '';
    ['name','price','stock'].forEach(k => editFields[k].removeAttribute('aria-invalid'));

    const name = editFields.name.value.trim();
    const category_id = Number(editFields.category.value);  
    const desc = editFields.desc.value.trim();
    const price = Number(editFields.price.value);
    const stock = Number(editFields.stock.value);

    let bad = null, msg = '';
    if (!name) { bad = 'name'; msg = 'Nama produk wajib diisi.'; }
    else if (!editFields.price.value || price <= 0) { bad = 'price'; msg = 'Harga harus lebih dari 0.'; }
    else if (editFields.stock.value === '' || stock < 0 || !Number.isInteger(stock)) {
      bad = 'stock'; msg = 'Stok harus bilangan bulat 0 atau lebih.';
    }
    if (bad) {
      editErr.textContent = msg;
      editFields[bad].setAttribute('aria-invalid', 'true');
      editFields[bad].focus();
      return;
    }

    const btn = $('e-submit');
    btn.disabled = true; btn.textContent = 'Menyimpan...';

    try {
      const res = await fetch(API + '/products/' + editingId, {
        method: 'PUT',
        headers: authHeaders(),
        body: JSON.stringify({ category_id, name, desc, price, stock })
      });
      const body = await res.json();
      if (!res.ok) throw new Error(body.error || 'Gagal menyimpan');

      editDlg.close();
      await Promise.all([loadProducts(), loadStats()]);
    } catch (ex) {
      editErr.textContent = ex.message;
    } finally {
      btn.disabled = false; btn.textContent = 'Simpan perubahan';
    }
  });

  // ---------- Hapus produk ----------
  async function deleteProduct(p) {
    const yes = confirm('Hapus produk "' + p.name + '"? Tindakan ini tidak bisa dibatalkan.');
    if (!yes) return;

    try {
      const res = await fetch(API + '/products/' + p.id, {
        method: 'DELETE',
        headers: authHeaders()
      });
      const body = await res.json();
      if (!res.ok) throw new Error(body.error || 'Gagal menghapus');
      await Promise.all([loadProducts(), loadStats()]);
    } catch (ex) {
      alert('Gagal menghapus: ' + ex.message);
    }
  }

  // ---------- Hapus Kategori ----------
  async function deleteCategory(c) {
    if (!confirm(`Hapus kategori "${c.name}"?`)) return;
    try {
      const res = await fetch(API + '/categories/' + c.id, {
        method: 'DELETE',
        headers: authHeaders()
      });
      const body = await res.json();
      if (!res.ok) throw new Error(body.error || 'Gagal menghapus');
      await loadCategories();
    } catch (ex) {
      alert(ex.message);
    }
  }

  // ---------- Tambah produk ----------
  const form = $('product-form');
  const f = {
    name: $('p-name'),
    category: $('p-category'),
    price: $('p-price'),
    stock: $('p-stock'),
    desc: $('p-desc')
  };
  const err = $('p-err');
  const ok = $('p-ok');

  form.addEventListener('submit', async function (e) {
    e.preventDefault();
    err.textContent = ''; ok.textContent = '';
    ['name', 'price', 'stock', 'category'].forEach(k => f[k].removeAttribute('aria-invalid'));

    const name = f.name.value.trim();
    const price = Number(f.price.value);
    const stock = Number(f.stock.value);
    const category_id = Number(f.category.value);

    let bad = null, msg = '';
    if (!category_id) { bad = 'category'; msg = 'Pilih kategori terlebih dahulu.'; }
    else if (!name) { bad = 'name'; msg = 'Nama produk wajib diisi.'; }
    else if (!f.price.value || price <= 0) { bad = 'price'; msg = 'Harga harus lebih dari 0.'; }
    else if (f.stock.value === '' || stock < 0 || !Number.isInteger(stock)) {
      bad = 'stock'; msg = 'Stok harus berupa bilangan bulat 0 atau lebih.';
    }
    if (bad) {
      err.textContent = msg;
      f[bad].setAttribute('aria-invalid', 'true');
      f[bad].focus();
      return;
    }

    try {
      const res = await fetch(API + '/products', {
        method: 'POST',
        headers: authHeaders(),
        body: JSON.stringify({ 
          category_id,
          name, 
          price, 
          stock, 
          desc: f.desc.value.trim() 
        })
      });
      const body = await res.json();
      if (!res.ok) throw new Error(body.error || 'Gagal menyimpan');

      form.reset();
      ok.textContent = 'Produk "' + name + '" berhasil ditambahkan.';
      f.name.focus();
      await Promise.all([loadProducts(), loadStats()]);
    } catch (ex) {
      err.textContent = ex.message;
    }
  });

  // ---------- Form Tambah Kategori ----------
  $('category-form').addEventListener('submit', async (e) => {
    e.preventDefault();
    const name = $('c-name').value.trim();
    const err = $('c-err');
    err.textContent = '';
    if (!name) { err.textContent = 'Nama kategori wajib diisi.'; return; }

    try {
      const res = await fetch(API + '/categories', {
        method: 'POST',
        headers: authHeaders(),
        body: JSON.stringify({ name })
      });
      const body = await res.json();
      if (!res.ok) throw new Error(body.error || 'Gagal menambah kategori');
      
      $('category-form').reset();
      await loadCategories(); 
    } catch (ex) {
      err.textContent = ex.message;
    }
  });

  // ---------- Event Pencarian & Pagination ----------
  let searchTimeout;
  $('search-product').addEventListener('input', (e) => {
    clearTimeout(searchTimeout);
    searchTimeout = setTimeout(() => {
      searchQuery = e.target.value.trim();
      currentPage = 1; 
      loadProducts();
    }, 500); 
  });

  $('btn-prev').addEventListener('click', () => {
    if (currentPage > 1) { currentPage--; loadProducts(); }
  });

  $('btn-next').addEventListener('click', () => {
    currentPage++; loadProducts();
  });

  // ---------- Logout ----------
  const logout = document.querySelector('.logout');
  if (logout) {
    logout.addEventListener('click', function (e) {
      e.preventDefault();
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      location.href = '../login_register/login.html';
    });
  }

  // ---------- Tabs ----------
  const tabs = ['products', 'categories', 'admins', 'users', 'sales'];
  tabs.forEach(t => {
    $('tab-' + t).addEventListener('click', () => {
      tabs.forEach(x => {
        const active = x === t;
        $('tab-' + x).setAttribute('aria-selected', active);
        $('panel-' + x).hidden = !active;
      });
    });
  });

  // ---------- Start ----------
  async function refreshAll() {
    try {
      await Promise.all([loadStats(), loadCategories(), loadProducts(), loadUsers(), loadSales()]);
    } catch (err) {
      console.error(err);
      alert('Gagal memuat data: ' + err.message);
      if (String(err.message).toLowerCase().includes('token')) {
        localStorage.removeItem('token');
        localStorage.removeItem('user');
        location.href = '../login_register/login.html';
      }
    }
  }

  refreshAll();

})(); 