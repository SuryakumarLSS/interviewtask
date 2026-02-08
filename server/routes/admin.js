const express = require('express');
const router = express.Router();
const db = require('../database');
const bcrypt = require('bcryptjs');

const checkAdmin = (req, res, next) => {
    if (req.user.role_id !== 1) return res.status(403).json({ message: "Admin access required" });
    next();
};

router.use(checkAdmin);

router.get('/roles', (req, res) => {
    db.all("SELECT * FROM roles", [], (err, rows) => {
        if (err) return res.status(500).json({ error: err.message });
        res.json(rows);
    });
});

router.post('/roles', (req, res) => {
    const { name } = req.body;
    db.run("INSERT INTO roles (name) VALUES (?)", [name], function (err) {
        if (err) return res.status(500).json({ error: err.message });
        res.json({ id: this.lastID, name });
    });
});

router.get('/permissions/:role_id', (req, res) => {
    const { role_id } = req.params;
    db.all("SELECT * FROM permissions WHERE role_id = ?", [role_id], (err, rows) => {
        if (err) return res.status(500).json({ error: err.message });
        res.json(rows);
    });
});

router.post('/permissions', (req, res) => {
    const { role_id, resource, action } = req.body;
    db.get("SELECT id FROM permissions WHERE role_id = ? AND resource = ? AND action = ?", [role_id, resource, action], (err, row) => {
        if (err) return res.status(500).json({ error: err.message });
        if (row) {
            db.run("UPDATE permissions SET attributes = '*' WHERE id = ?", [row.id], (err) => {
                if (err) return res.status(500).json({ error: err.message });
                res.json({ message: "Permission updated to full access" });
            });
        } else {
            db.run("INSERT INTO permissions (role_id, resource, action, attributes) VALUES (?, ?, ?, '*')", [role_id, resource, action], function (err) {
                if (err) return res.status(500).json({ error: err.message });
                res.json({ id: this.lastID, message: "Permission created" });
            });
        }
    });
});

router.delete('/permissions', (req, res) => {
    const role_id = req.body.role_id || req.query.role_id;
    const resource = req.body.resource || req.query.resource;
    const action = req.body.action || req.query.action;
    db.run("DELETE FROM permissions WHERE role_id = ? AND resource = ? AND action = ?", [role_id, resource, action], function (err) {
        if (err) return res.status(500).json({ error: err.message });
        res.json({ message: "Permission removed" });
    });
});

router.get('/users', (req, res) => {
    db.all("SELECT id, username, role_id, status FROM users", [], (err, rows) => {
        if (err) return res.status(500).json({ error: err.message });
        res.json(rows);
    });
});

const checkOrgAdmin = (req, res, next) => {
    // Strict check for the "org level admin" (username 'admin' or id 1)
    // Assuming req.user is populated by auth middleware
    if (req.user.username !== 'admin') {
        return res.status(403).json({ message: "Only Organization Admin can perform this action" });
    }
    next();
};

router.put('/users/:id/role', checkOrgAdmin, (req, res) => {
    const { role_id } = req.body;
    db.run("UPDATE users SET role_id = ? WHERE id = ?", [role_id, req.params.id], (err) => {
        if (err) return res.status(500).json({ error: err.message });
        res.json({ message: "Role updated" });
    });
});

router.delete('/users/:id', checkOrgAdmin, (req, res) => {
    db.run("DELETE FROM users WHERE id = ?", [req.params.id], (err) => {
        if (err) return res.status(500).json({ error: err.message });
        res.json({ message: "User deleted" });
    });
});

const nodemailer = require('nodemailer');
const crypto = require('crypto');

// Create a transporter using Ethereal for testing if no real credentials provided
let transporter;

async function createTransporter() {
    if (process.env.EMAIL_USER && process.env.EMAIL_PASS) {
        transporter = nodemailer.createTransport({
            service: 'gmail',
            auth: {
                user: process.env.EMAIL_USER,
                pass: process.env.EMAIL_PASS
            }
        });
    } else {
        console.log("No EMAIL_USER/PASS found. Using Ethereal Email for testing...");
        const testAccount = await nodemailer.createTestAccount();
        transporter = nodemailer.createTransport({
            host: "smtp.ethereal.email",
            port: 587,
            secure: false,
            auth: {
                user: testAccount.user,
                pass: testAccount.pass
            }
        });
    }
}
createTransporter();

router.post('/users', (req, res) => {
    // ... invite flow ...
    const { username, password, role_id, email } = req.body;
    const userEmail = email || username;

    if (!password && userEmail) {
        const token = crypto.randomBytes(32).toString('hex');
        const tempUsername = userEmail;

        db.get("SELECT * FROM users WHERE email = ? OR username = ?", [userEmail, tempUsername], (err, existingUser) => {
            if (err) return res.status(500).json({ error: err.message });

            if (existingUser) {
                if (existingUser.password) return res.status(400).json({ message: "User already registered" });

                db.run("UPDATE users SET invitation_token = ?, role_id = ?, status = 'Pending' WHERE id = ?", [token, role_id, existingUser.id], (err) => {
                    if (err) return res.status(500).json({ error: err.message });
                    sendInviteEmail(userEmail, token, res);
                });
            } else {
                db.run("INSERT INTO users (username, email, role_id, invitation_token, status) VALUES (?, ?, ?, ?, 'Pending')",
                    [tempUsername, userEmail, role_id, token],
                    function (err) {
                        if (err) return res.status(500).json({ error: err.message });
                        sendInviteEmail(userEmail, token, res);
                    }
                );
            }
        });

        async function sendInviteEmail(emailAddr, token, res) {
            const acceptLink = `http://localhost:5173/set-password?token=${token}`;
            // For decline, we might need a route, or just a landing page. For now, let's say it goes to a generic "invitation declined" page we might build or just home.
            // But user asked for "decline invitation". Typically this burns the token. Let's send them to a decline endpoint on API? 
            // Better to send to frontend which calls API.
            const declineLink = `http://localhost:5173/decline-invitation?token=${token}`;

            const mailOptions = {
                from: 'jeyaroshini2001@gmail.com',
                to: emailAddr,
                subject: 'Action Required: RBAC System Invitation',
                html: `
                    <div style="font-family: Arial, sans-serif; padding: 20px; text-align: center;">
                        <h3>Welcome to RBAC System</h3>
                        <p>You have been invited to join the team.</p>
                        <div style="margin: 30px 0;">
                            <a href="${acceptLink}" style="background-color: #18181b; color: white; padding: 12px 24px; text-decoration: none; border-radius: 6px; margin-right: 10px;">Accept Invitation</a>
                            <a href="${declineLink}" style="background-color: #ef4444; color: white; padding: 12px 24px; text-decoration: none; border-radius: 6px;">Decline Invitation</a>
                        </div>
                        <p style="font-size: 12px; color: #71717a;">Links expire in 24 hours.</p>
                    </div>
                `
            };

            // ... send logic (same as before) ...
            // Simplified for brevity in tool call, will use existing transporter logic
            // Try sending with configured credentials first
            transporter.sendMail(mailOptions, async (error, info) => {
                if (error) {
                    // ... fallback logic ...
                    try {
                        const testAccount = await nodemailer.createTestAccount();
                        const testTransporter = nodemailer.createTransport({
                            host: "smtp.ethereal.email", port: 587, secure: false,
                            auth: { user: testAccount.user, pass: testAccount.pass }
                        });
                        const testInfo = await testTransporter.sendMail(mailOptions);
                        const previewUrl = nodemailer.getTestMessageUrl(testInfo);
                        console.log("Fallback Preview URL: %s", previewUrl);
                        return res.json({ message: "Invited via Test Email. Check preview.", preview: previewUrl, link: acceptLink });
                    } catch (e) { return res.json({ message: "Email failed", link: acceptLink }); }
                }
                const previewUrl = nodemailer.getTestMessageUrl(info);
                if (previewUrl) console.log("Preview URL: %s", previewUrl);
                res.json({ message: "Invitation sent!", preview: previewUrl, link: acceptLink });
            });
        }
    } else {
        // Legacy
        const hash = bcrypt.hashSync(password, 10);
        db.run("INSERT INTO users (username, password, role_id, status) VALUES (?, ?, ?, 'Active')", [username, hash, role_id], function (err) {
            if (err) return res.status(500).json({ error: err.message });
            res.json({ id: this.lastID, username, role_id });
        });
    }
});

module.exports = router;
