const express = require('express');
const router = express.Router();
const bcrypt = require('bcryptjs');
const jwt = require('jsonwebtoken');
const db = require('../database');

router.post('/login', (req, res) => {
    const { username, password } = req.body;
    db.get('SELECT * FROM users WHERE username = ?', [username], (err, user) => {
        if (err) return res.status(500).json({ error: err.message });
        if (!user) return res.status(401).json({ message: 'User not found' });
        // Check for null password (invited users)
        if (!user.password) return res.status(401).json({ message: 'Account pending activation. Please use the invitation link.' });

        if (!bcrypt.compareSync(password, user.password)) return res.status(401).json({ message: 'Invalid credentials' });

        const token = jwt.sign(
            { id: user.id, username: user.username, role_id: user.role_id },
            process.env.JWT_SECRET || 'secret_key',
            { expiresIn: '1h' }
        );
        res.json({ token, user: { id: user.id, username: user.username, role_id: user.role_id } });
    });
});

router.post('/set-password', (req, res) => {
    const { token, password } = req.body;

    db.get("SELECT * FROM users WHERE invitation_token = ?", [token], (err, user) => {
        if (err) return res.status(500).json({ error: err.message });
        if (!user) return res.status(400).json({ message: "Invalid or expired invitation token" });

        const hash = bcrypt.hashSync(password, 10);

        db.run("UPDATE users SET password = ?, invitation_token = NULL, status = 'Active' WHERE id = ?", [hash, user.id], (err) => {
            if (err) return res.status(500).json({ error: err.message });
            res.json({ message: "Password set successfully" });
        });
    });
});

router.post('/decline-invitation', (req, res) => {
    const { token } = req.body;

    db.get("SELECT * FROM users WHERE invitation_token = ?", [token], (err, user) => {
        if (err) return res.status(500).json({ error: err.message });
        if (!user) return res.status(400).json({ message: "Invalid or expired invitation token" });

        // Update status to Declined and remove token so it can't be used again
        db.run("UPDATE users SET status = 'Declined', invitation_token = NULL WHERE id = ?", [user.id], (err) => {
            if (err) return res.status(500).json({ error: err.message });
            res.json({ message: "Invitation declined" });
        });
    });
});

router.get('/permissions', (req, res) => {
    const token = req.header('Authorization')?.replace('Bearer ', '');
    if (!token) return res.status(401).json({ message: "No token" });
    try {
        const decoded = jwt.verify(token, process.env.JWT_SECRET || 'secret_key');
        db.all("SELECT * FROM permissions WHERE role_id = ?", [decoded.role_id], (err, rows) => {
            if (err) return res.status(500).json({ error: err.message });
            res.json(rows);
        });
    } catch (e) { res.status(401).json({ message: "Invalid token" }); }
});

router.get('/employees', (req, res) => {
    const token = req.header('Authorization')?.replace('Bearer ', '');
    if (!token) return res.status(401).json({ message: "No token" });
    try {
        jwt.verify(token, process.env.JWT_SECRET || 'secret_key');
        db.all("SELECT id, name FROM employees", [], (err, rows) => {
            if (err) return res.status(500).json({ error: err.message });
            res.json(rows);
        });
    } catch (e) { res.status(401).json({ message: "Invalid token" }); }
});

module.exports = router;
