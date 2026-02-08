const express = require('express');
const router = express.Router();
const db = require('../database');
const rbacMiddleware = require('../middleware/rbac');

const ALLOWED_RESOURCES = ['employees', 'projects', 'orders'];
const REQUIRED_FIELDS = {
    employees: ['name', 'position', 'salary', 'department'],
    projects: ['name', 'assigned_to', 'status', 'budget'],
    orders: ['customer_name', 'amount', 'status', 'order_date']
};

const checkResource = (req, res, next) => {
    if (!ALLOWED_RESOURCES.includes(req.params.resource)) return res.status(404).json({ message: "Resource not found" });
    next();
};

router.get('/:resource', checkResource, (req, res, next) => {
    rbacMiddleware(req.params.resource, 'read')(req, res, next);
}, (req, res) => {
    const resource = req.params.resource;
    db.all(`SELECT * FROM ${resource}`, [], (err, rows) => {
        if (err) return res.status(500).json({ error: err.message });
        res.json(rows);
    });
});

router.post('/:resource', checkResource, (req, res, next) => {
    rbacMiddleware(req.params.resource, 'create')(req, res, next);
}, (req, res) => {
    const resource = req.params.resource;
    const data = req.body;
    const required = REQUIRED_FIELDS[resource] || [];
    const missing = required.filter(field => !data[field]);
    if (missing.length > 0) return res.status(400).json({ message: `Missing required fields: ${missing.join(', ')}` });

    const keys = Object.keys(data);
    const values = Object.values(data);
    const placeholders = keys.map(() => '?').join(', ');
    db.run(`INSERT INTO ${resource} (${keys.join(', ')}) VALUES (${placeholders})`, values, function (err) {
        if (err) return res.status(500).json({ error: err.message });
        res.json({ id: this.lastID, ...data });
    });
});

router.put('/:resource/:id', checkResource, (req, res, next) => {
    rbacMiddleware(req.params.resource, 'update')(req, res, next);
}, (req, res) => {
    const resource = req.params.resource;
    const id = req.params.id;
    const data = req.body;
    const keys = Object.keys(data);
    const values = Object.values(data);
    const setClause = keys.map(key => `${key} = ?`).join(', ');
    db.run(`UPDATE ${resource} SET ${setClause} WHERE id = ?`, [...values, id], function (err) {
        if (err) return res.status(500).json({ error: err.message });
        res.json({ message: "Updated successfully" });
    });
});

router.delete('/:resource/:id', checkResource, (req, res, next) => {
    rbacMiddleware(req.params.resource, 'delete')(req, res, next);
}, (req, res) => {
    const resource = req.params.resource;
    const id = req.params.id;
    db.run(`DELETE FROM ${resource} WHERE id = ?`, id, function (err) {
        if (err) return res.status(500).json({ error: err.message });
        res.json({ message: "Deleted successfully" });
    });
});

module.exports = router;
