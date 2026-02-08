const sqlite3 = require('sqlite3').verbose();
const bcrypt = require('bcryptjs');

const db = new sqlite3.Database('./rbac.db', (err) => {
    if (err) console.error('Error opening database', err.message);
    else {
        console.log('Connected to the SQLite database.');
        initializeDatabase();
    }
});

function initializeDatabase() {
    db.serialize(() => {
        db.run(`CREATE TABLE IF NOT EXISTS roles (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT UNIQUE)`);
        db.run(`CREATE TABLE IF NOT EXISTS users (id INTEGER PRIMARY KEY AUTOINCREMENT, username TEXT UNIQUE, password TEXT, role_id INTEGER, FOREIGN KEY(role_id) REFERENCES roles(id))`);
        db.run(`CREATE TABLE IF NOT EXISTS permissions (id INTEGER PRIMARY KEY AUTOINCREMENT, role_id INTEGER, resource TEXT, action TEXT, attributes TEXT, FOREIGN KEY(role_id) REFERENCES roles(id))`);
        db.run(`CREATE TABLE IF NOT EXISTS employees (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT, position TEXT, salary INTEGER, department TEXT)`);
        db.run(`CREATE TABLE IF NOT EXISTS projects (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT, assigned_to TEXT, status TEXT, budget INTEGER)`);
        db.run(`CREATE TABLE IF NOT EXISTS orders (id INTEGER PRIMARY KEY AUTOINCREMENT, customer_name TEXT, amount INTEGER, status TEXT, order_date TEXT)`);
        seedData();
    });
}

function seedData() {
    db.get("SELECT id FROM roles WHERE name = 'Admin'", (err, row) => {
        if (!row) {
            db.run("INSERT INTO roles (name) VALUES ('Admin')", function (err) {
                if (err) return console.error(err.message);
                const adminRoleId = this.lastID;
                const passwordHash = bcrypt.hashSync('admin123', 10);
                db.run("INSERT INTO users (username, password, role_id) VALUES (?, ?, ?)",
                    ['admin', passwordHash, adminRoleId],
                    (err) => { if (err) console.error(err.message); else console.log("Admin user created"); }
                );
                const resources = ['employees', 'projects', 'orders', 'roles', 'users'];
                const actions = ['read', 'create', 'update', 'delete'];
                const stmt = db.prepare("INSERT INTO permissions (role_id, resource, action, attributes) VALUES (?, ?, ?, ?)");
                resources.forEach(res => {
                    actions.forEach(action => {
                        stmt.run(adminRoleId, res, action, '*');
                    });
                });
                stmt.finalize();
            });
            db.run("INSERT INTO roles (name) VALUES ('user')");
        }
    });
}

module.exports = db;
