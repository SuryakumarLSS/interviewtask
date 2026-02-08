const db = require('../database');

module.exports = (resource, action) => {
    return (req, res, next) => {
        if (!req.user) return res.status(401).json({ message: "Unauthorized" });
        const { role_id } = req.user;

        db.get(
            "SELECT attributes FROM permissions WHERE role_id = ? AND resource = ? AND action = ?",
            [role_id, resource, action],
            (err, row) => {
                if (err) return res.status(500).json({ error: err.message });
                if (!row) return res.status(403).json({ message: "Forbidden" });
                req.allowedAttributes = row.attributes ? row.attributes.split(',') : [];
                next();
            }
        );
    };
};
